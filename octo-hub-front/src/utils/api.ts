export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";

export type ApiFetchOptions = Omit<RequestInit, "body"> & {
  body?: unknown;
  baseUrl?: string;
};

export async function apiFetch<TResponse = unknown>(
  path: string,
  options: ApiFetchOptions = {}
): Promise<TResponse> {
  const { body, headers, baseUrl, ...rest } = options;
  const url = `${baseUrl ?? API_BASE_URL}${path}`;

  const init: RequestInit = {
    // Default to POST if a body is provided, otherwise keep caller's method or default GET
    method: body !== undefined ? "POST" : (rest.method ?? "GET"),
    headers: {
      "Content-Type": "application/json",
      ...(headers || {}),
    },
    ...rest,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  };

  const response = await fetch(url, init);

  const contentType = response.headers.get("content-type") || "";
  const isJson = contentType.includes("application/json");
  const data = isJson ? await response.json().catch(() => null) : await response.text();

  if (!response.ok) {
    const message = typeof data === "object" && data !== null
      ? ((data as any).message || (data as any).error || response.statusText)
      : (response.statusText || "Request failed");
    throw new Error(message);
  }

  return data as TResponse;
}

export function apiGet<TResponse = unknown>(
  path: string,
  options: Omit<ApiFetchOptions, "method" | "body"> = {}
) {
  return apiFetch<TResponse>(path, { ...options, method: "GET" });
}

export function apiPost<TResponse = unknown>(
  path: string,
  body?: unknown,
  options: Omit<ApiFetchOptions, "method" | "body"> = {}
) {
  return apiFetch<TResponse>(path, { ...options, method: "POST", body });
}


