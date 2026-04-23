/**
 * Base class for all domain repositories.
 * Extends EventTarget to allow framework-agnostic state change notifications.
 */
export abstract class Repository extends EventTarget {
  /**
   * Dispatches a 'change' event to notify observers of a state mutation.
   */
  protected notifyChange() {
    this.dispatchEvent(new Event('change'))
  }
}
