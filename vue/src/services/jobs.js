async function parseResponse(response) {
  const contentType = response.headers.get("content-type") || "";
  const body = contentType.includes("application/json")
    ? await response.json()
    : await response.text();

  if (!response.ok) {
    const message = typeof body === "object" ? body.error : body;
    throw new Error(message || `HTTP ${response.status}`);
  }
  return body;
}

export async function startJob(operation, options = {}) {
  const response = await fetch(`/jobs/start/${operation}`, {
    method: "POST",
    ...options,
  });
  return parseResponse(response);
}

export async function getJob(id) {
  const response = await fetch(`/jobs/${encodeURIComponent(id)}`);
  return parseResponse(response);
}

export async function listJobs(type = "", limit = 10) {
  const params = new URLSearchParams({ limit: String(limit) });
  if (type) params.set("type", type);
  const response = await fetch(`/jobs?${params}`);
  return parseResponse(response);
}

export async function cancelJob(id) {
  const response = await fetch(`/jobs/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
  return parseResponse(response);
}

export function isTerminalJob(job) {
  return ["succeeded", "failed", "cancelled"].includes(job?.status);
}
