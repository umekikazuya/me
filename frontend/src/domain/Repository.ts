import { signal, type Signal } from '@preact/signals-core'

/**
 * Standard machine-readable status for data-driven operations.
 */
export type StateStatus = 'idle' | 'loading' | 'error' | 'success'

/**
 * Unified state model for ensuring consistent UI across the package.
 * Decouples the presentation layer from the specific data fetching details.
 */
export interface IState<T> {
  readonly status: StateStatus
  readonly data: T | null
  readonly error: { code: string; message: string } | null
}

/**
 * Creates an initial idle state.
 */
export function createInitialState<T>(initialData: T | null = null): IState<T> {
  return {
    status: 'idle',
    data: initialData,
    error: null,
  }
}

/**
 * Base class for all domain repositories.
 * Extends EventTarget for discrete event notifications.
 * Provides standard Signal-based state management and generation tracking.
 */
export abstract class Repository extends EventTarget {
  private _generation = 0

  /**
   * Tracks a new asynchronous operation, returning a generation ID.
   * Use this to discard stale request results.
   */
  protected nextGeneration(): number {
    this._generation += 1
    return this._generation
  }

  /**
   * Verifies if the given generation ID is still current.
   */
  protected isCurrent(gen: number): boolean {
    return this._generation === gen
  }

  /**
   * Helper to update a signal-based state.
   */
  protected updateState<T>(
    stateSignal: Signal<IState<T>>,
    patch: Partial<IState<T>>,
  ) {
    stateSignal.value = { ...stateSignal.value, ...patch }
  }

  /**
   * Dispatches a legacy generic change event for backward compatibility.
   * @deprecated Prefer granular typed events or direct Signal consumption.
   */
  protected notifyChange() {
    this.dispatchEvent(new Event('change'))
  }
}
