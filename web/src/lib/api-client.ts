import { env } from '@/config/env'

// тело ошибки с бэкенда, как httpx.ErrorResponse в catalog
type ErrorBody = {
  slug: string
  msg: string
}

// ошибка ответа API: статус и slug, чтобы UI мог отличить card-not-found от остального
export class ApiError extends Error {
  status: number
  slug: string

  constructor(status: number, slug: string, msg: string) {
    super(msg)
    this.name = 'ApiError'
    this.status = status
    this.slug = slug
  }
}

// query-параметры; undefined пропускаем, чтобы необязательные поля не уходили как "undefined"
type QueryParams = Record<string, string | number | undefined>

function buildUrl(path: string, params: QueryParams = {}): string {
  const url = new URL(path, env.API_URL)
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined) url.searchParams.set(key, String(value))
  }
  return url.toString()
}

// общий запрос: базовый URL, проверка статуса и JSON в одном месте для всех фич
async function request<T>(path: string, init?: RequestInit & { params?: QueryParams }): Promise<T> {
  const res = await fetch(buildUrl(path, init?.params), init)

  // fetch не бросает на 4xx/5xx, статус проверяем сами
  if (!res.ok) {
    // тело может оказаться не JSON (например, ответ прокси), тогда берём заглушку
    const body = (await res.json().catch(() => null)) as Partial<ErrorBody> | null
    throw new ApiError(
      res.status,
      body?.slug ?? 'unknown-error',
      body?.msg ?? `Request failed: ${String(res.status)}`,
    )
  }

  // json() возвращает any: as доверяет форме из контракта, в рантайме не проверяет
  return (await res.json()) as T
}

export const api = {
  get<T>(path: string, params?: QueryParams): Promise<T> {
    return request<T>(path, { method: 'GET', params })
  },
}
