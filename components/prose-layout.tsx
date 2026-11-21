import type { ReactNode } from 'react'
import { HomeLayout } from 'fumadocs-ui/layouts/home'
import { useBaseOptions } from './layout-options'
import { SiteFooter } from './site-footer'

// ProseLayout is the header of the site around a column of text: a page of Markdown such as Why.
export function ProseLayout({ children }: { children: ReactNode }) {
  const options = useBaseOptions()
  return (
    <HomeLayout {...options}>
      <main className="mx-auto w-full max-w-3xl flex-1 px-4 pb-24 pt-12 sm:pt-16">{children}</main>
      <SiteFooter />
    </HomeLayout>
  )
}
