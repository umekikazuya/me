import { execSync } from 'node:child_process'
import { defineConfig } from 'vite'

const git = (args: string) => {
  try {
    return execSync(`git ${args}`, { stdio: ['ignore', 'pipe', 'ignore'] })
      .toString()
      .trim()
  } catch {
    return ''
  }
}

export default defineConfig(({ command }) => {
  if (command === 'serve' && !process.env.API_PORT) {
    throw new Error('API_PORT is not set')
  }

  return {
    define: {
      __BUILD_COMMIT__: JSON.stringify(git('rev-parse --short HEAD')),
      __BUILD_COMMITTED_AT__: JSON.stringify(git('log -1 --format=%cI')),
    },
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
