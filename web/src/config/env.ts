import { z } from 'zod'

// схема переменных окружения; всё, что читает import.meta.env, описываем здесь
const envSchema = z.object({
  VITE_APP_API_URL: z.url(),
})

// проверяем при импорте модуля: без нужных переменных приложение падает сразу с понятной ошибкой
const parsed = envSchema.safeParse(import.meta.env)
if (!parsed.success) {
  throw new Error(`Invalid env (web/.env):\n${z.prettifyError(parsed.error)}`)
}

// дальше в коде env.API_URL, а не import.meta.env напрямую
export const env = {
  API_URL: parsed.data.VITE_APP_API_URL,
}
