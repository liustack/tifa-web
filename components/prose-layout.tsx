import type { ReactNode } from 'react'
import { SiteShell } from './site-shell'

// ProseLayout is the header and the footer of the site around a column of text: a page of
// Markdown such as Why.
export function ProseLayout({ children }: { children: ReactNode }) {
  return (
    <SiteShell>
      <main className="mx-auto w-full max-w-3xl flex-1 px-4 pb-24 pt-12 sm:pt-16">{children}</main>
    </SiteShell>
  )
}
