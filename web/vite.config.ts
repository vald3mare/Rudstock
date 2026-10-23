import { fileURLToPath, URL } from 'node:url'

import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    // тот же алиас, что в tsconfig.app.json: TypeScript проверяет импорты, Vite их собирает
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
