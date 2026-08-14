import { defineConfig } from 'vite'

export default defineConfig(({ command }) => {
  if (command === 'serve' && !process.env.API_PORT) {
    throw new Error('API_PORT is not set')
  }

  return {
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
