import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import App from '@/app'

import './index.css'

// явная проверка вместо `!`: без #root в index.html падаем с понятной ошибкой
const root = document.getElementById('root')
if (!root) throw new Error('No #root element in index.html')

createRoot(root).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
