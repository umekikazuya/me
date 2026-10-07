/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_ORIGIN?: string
  readonly VITE_GA_MEASUREMENT_ID?: string
}

declare const __BUILD_COMMIT__: string
declare const __BUILD_COMMITTED_AT__: string
