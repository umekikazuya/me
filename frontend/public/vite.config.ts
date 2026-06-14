import { defineConfig } from 'vite'

const apiPort = process.env.API_PORT
if (!apiPort) throw new Error('API_PORT is not set')

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: `http://localhost:${apiPort}`,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api(?=\/|$)/, ''),
      },
    },
  },
})
