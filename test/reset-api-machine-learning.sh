#!/usr/bin/env bash

set -euo pipefail

readonly USERNAMEAPI="userapi"
readonly PASSWORDAPI="fortinet123!"
readonly VDOMAPI="root"
readonly FWBMGTIP="10.163.7.21"
readonly FWBMGTPORT="443"

readonly BASE_URL="https://${FWBMGTIP}:${FWBMGTPORT}/api/v2.0/machine_learning"
readonly TOKEN_DATA="{\"username\":\"${USERNAMEAPI}\",\"password\":\"${PASSWORDAPI}\",\"vdom\":\"${VDOMAPI}\"}"
readonly TOKEN="$(printf '%s' "${TOKEN_DATA}" | base64 | tr -d '\n')"

if ! command -v curl >/dev/null 2>&1; then
  echo "Erreur : curl est requis." >&2
  exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
  echo "Erreur : jq est requis pour extraire les rule_id de la réponse." >&2
  exit 1
fi

echo "Récupération des règles API Learning depuis ${FWBMGTIP}:${FWBMGTPORT}..."

policy_rules="$(
  curl --silent --show-error --fail-with-body --insecure \
    --request GET \
    --header "Authorization: ${TOKEN}" \
    --header "Accept: application/json" \
    "${BASE_URL}/api_learning_policy.get_policy_rule"
)"

if ! jq -e 'type == "object"' >/dev/null <<<"${policy_rules}"; then
  echo "Erreur : le serveur n'a pas retourné une réponse JSON valide." >&2
  exit 1
fi

mapfile -t rules < <(
  jq -r '.results[]?.rule[]? | [.id, .["domain-name"]] | @tsv' <<<"${policy_rules}"
)

if ((${#rules[@]} == 0)); then
  echo "Aucune règle API Learning trouvée."
  exit 0
fi

echo "${#rules[@]} règle(s) trouvée(s)."

for rule in "${rules[@]}"; do
  IFS=$'\t' read -r rule_id domain_name <<<"${rule}"

  echo "Reset de ${domain_name:-domaine inconnu} (rule_id=${rule_id})..."

  curl --silent --show-error --fail-with-body --insecure \
    --request POST \
    --header "Authorization: ${TOKEN}" \
    --header "Content-Type: application/json" \
    "${BASE_URL}/api_learning_policy.refreshdomain?rule_id=${rule_id}"

  echo
done

echo "Reset Machine Learning terminé."
