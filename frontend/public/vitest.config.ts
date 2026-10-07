import { defineConfig } from 'vitest/config'

export default defineConfig({
  define: {
    __BUILD_COMMIT__: JSON.stringify(''),
    __BUILD_COMMITTED_AT__: JSON.stringify(''),
  },
  test: {
    environment: 'happy-dom',
    setupFiles: ['./src/test/setup.ts'],
    reporters: process.env.GITHUB_ACTIONS
      ? ['dot', 'github-actions']
      : ['default'],
    coverage: {
      provider: 'v8',
      include: ['src/**/*.{ts,tsx}'],
      reporter: ['text', 'json-summary'],
      reportsDirectory: './coverage',
      reportOnFailure: true,
    },
  },
})
