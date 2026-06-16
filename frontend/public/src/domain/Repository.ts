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
 * Extends EventTarget for discrete event notifications (Observer pattern).
 * Dispatches 'change' on state mutation; specific events on domain actions.
 */
export abstract class Repository extends EventTarget {
  private _generation = 0

  protected nextGeneration(): number {
    this._generation += 1
    return this._generation
  }

  protected isCurrent(gen: number): boolean {
    return this._generation === gen
  }

  protected emitChange(): void {
    this.dispatchEvent(new CustomEvent('change'))
  }
}
