export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";

export type ApiFetchOptions = Omit<RequestInit, "body"> & {
  body?: unknown;
  baseUrl?: string;
};

export interface ApiResponse {  
  errcode: number;
  message?: string;
  data?: any;
}

export async function apiFetch<TResponse = ApiResponse>(
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
  // 判断业务错误码
  if (typeof data === "object" && data !== null && "errcode" in data) {
    const errcode = (data as any).errcode;
    if (typeof errcode === "number" && errcode !== 0) {
      // 401 未授权，跳转到登录页
      if (errcode === 401) {
        // 检查是否在客户端环境
        if (typeof window !== "undefined") {
          window.location.href = "/login/login";
          return data as TResponse; // 跳转后直接返回，避免抛出错误
        }
      }
      const message = (data as any).message || `请求失败，错误码：${errcode}`;
      throw new Error(message);
    }
  }else{
    console.log(data);
    throw new Error('请求失败, 未知错误');
  }
  return data as TResponse;
}

export function apiGet<TResponse = ApiResponse>(
  path: string,
  options: Omit<ApiFetchOptions, "method" | "body"> = {}
) {
  return apiFetch<TResponse>(path, { ...options, method: "GET" });
}

export function apiPost<TResponse = ApiResponse>(
  path: string,
  body?: unknown,
  options: Omit<ApiFetchOptions, "method" | "body"> = {}
) {
  return apiFetch<TResponse>(path, { ...options, method: "POST", body });
}


