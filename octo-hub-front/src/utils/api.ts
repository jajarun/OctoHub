export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";

export type ApiFetchOptions = Omit<RequestInit, "body"> & {
  body?: unknown;
  baseUrl?: string;
  skipAuth?: boolean; // 是否跳过自动添加token
};

export interface ApiResponse {  
  errcode: number;
  errmsg?: string; // 与后端保持一致
  data?: any;
}

// WebSocket相关类型定义
export interface WebSocketConnectionResponse {
  wsUrl: string;
}

// Token 管理工具
export const TokenManager = {
  getToken(): string | null {
    if (typeof window === "undefined") return null;
    return localStorage.getItem("access_token");
  },
  
  setToken(token: string): void {
    if (typeof window === "undefined") return;
    localStorage.setItem("access_token", token);
  },
  
  removeToken(): void {
    if (typeof window === "undefined") return;
    localStorage.removeItem("access_token");
  }
};

export async function apiFetch<TResponse = ApiResponse>(
  path: string,
  options: ApiFetchOptions = {}
): Promise<TResponse> {
  const { body, headers, baseUrl, skipAuth, ...rest } = options;
  const url = `${baseUrl ?? API_BASE_URL}${path}`;

  // 构建请求头
  const requestHeaders: Record<string, string> = {
    "Content-Type": "application/json",
    ...(headers as Record<string, string> || {}),
  };

  // 自动添加 JWT token，除非显式跳过或者是登录接口
  const isLoginPath = path.includes("/login");
  if (!skipAuth && !isLoginPath) {
    const token = TokenManager.getToken();
    if (token) {
      requestHeaders["Authorization"] = `Bearer ${token}`;
    }
  }

  const init: RequestInit = {
    // Default to POST if a body is provided, otherwise keep caller's method or default GET
    method: body !== undefined ? "POST" : (rest.method ?? "GET"),
    headers: requestHeaders,
    ...rest,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  };

  const response = await fetch(url, init);

  const contentType = response.headers.get("content-type") || "";
  const isJson = contentType.includes("application/json");
  const data = isJson ? await response.json().catch(() => null) : await response.text();
  if (!response.ok) {
    //状态码为401，清除token并跳转到登录页
    if (response.status === 401) {
      TokenManager.removeToken();
      window.location.href = "/login";
      return data as TResponse;
    }
    const message = typeof data === "object" && data !== null
      ? ((data as any).errmsg || (data as any).message || (data as any).error || response.statusText)
      : (response.statusText || "Request failed");
    throw new Error(message);
  }
  // 判断业务错误码
  if (typeof data === "object" && data !== null && "errcode" in data) {
    const errcode = (data as any).errcode;
    if (typeof errcode === "number" && errcode !== 0) {
      // 401 未授权，清除token并跳转到登录页
      if (errcode === 401) {
        // 检查是否在客户端环境
        if (typeof window !== "undefined") {
          TokenManager.removeToken(); // 清除无效token
          window.location.href = "/login";
          return data as TResponse; // 跳转后直接返回，避免抛出错误
        }
      }
      const message = (data as any).errmsg || (data as any).message || `请求失败，错误码：${errcode}`;
      throw new Error(message);
    }
  } else {
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


