import { useSyncExternalStore } from 'react'

/**
 * React ↔ Signal ブリッジフック。
 *
 * Domain 層の Signal（@preact/signals-core）を React の
 * ライフサイクルに安全に接続する。
 *
 * @example
 * ```tsx
 * const status = useSignal(authRepo.status) // Signal<string>
 * ```
 */
export function useSignal<T>(signal: {
  value: T
  subscribe: (cb: () => void) => () => void
}): T {
  return useSyncExternalStore(
    (cb) => signal.subscribe(cb),
    () => signal.value,
    () => signal.value,
  )
}
