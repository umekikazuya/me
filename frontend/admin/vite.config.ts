import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(({ command }) => {
  if (command === 'serve' && !process.env.API_PORT) {
    throw new Error('API_PORT is not set')
  }

  return {
    plugins: [react()],
    base: '/',
    server: {
      proxy: {
        '/api': {
          target: `http://localhost:${process.env.API_PORT || 8080}`,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api(?=\/|$)/, ''),
        },
      },
    },
  }
})
