/**
 * Google Analytics (gtag.js) utility.
 *
 * Initializes GA and provides a page_view tracking helper.
 * No-ops when VITE_GA_MEASUREMENT_ID is not set (e.g. local dev).
 */

declare global {
  interface Window {
    dataLayer: unknown[]
    gtag: (...args: unknown[]) => void
  }
}

const GA_ID = import.meta.env.VITE_GA_MEASUREMENT_ID

/**
 * Injects the gtag.js script and initializes GA.
 * Safe to call multiple times – subsequent calls are ignored.
 */
export function initGA(): void {
  if (!GA_ID || document.getElementById('ga-script')) return

  const script = document.createElement('script')
  script.id = 'ga-script'
  script.src = `https://www.googletagmanager.com/gtag/js?id=${GA_ID}`
  script.async = true
  document.head.appendChild(script)

  window.dataLayer = window.dataLayer ?? []
  window.gtag = function gtag() {
    // biome-ignore lint/complexity/noArguments: Google Analytics requires the arguments object
    window.dataLayer.push(arguments)
  }
  window.gtag('js', new Date())
  window.gtag('config', GA_ID, {
    // SPA のため自動ページビュー送信は無効にして手動で送る
    send_page_view: false,
  })
}

/**
 * Sends a page_view event to GA.
 * @param path - The path to track (e.g. '/articles')
 */
export function trackPageView(path: string): void {
  if (!GA_ID || typeof window.gtag !== 'function') return
  window.gtag('event', 'page_view', {
    page_location: `${window.location.origin}${path}`,
    page_path: path,
  })
}
