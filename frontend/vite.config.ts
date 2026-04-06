import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig(() => {
  return {
    server: {
      proxy: {
        '/api': {
          target: `http://localhost:${process.env.BACKEND_PORT ?? 8080}/`,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ''),
        },
      },
    },
  }
})
