import { html, LitElement } from 'lit'
import { customElement } from 'lit/decorators.js'
import '../components/admin/ui/me-admin-auth-boundary.js'
import './page-admin-dashboard.js'

@customElement('page-admin-entry')
export class PageAdminEntry extends LitElement {
  render() {
    return html`
      <me-admin-auth-boundary>
        <page-admin-dashboard></page-admin-dashboard>
      </me-admin-auth-boundary>
    `
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'page-admin-entry': PageAdminEntry
  }
}
