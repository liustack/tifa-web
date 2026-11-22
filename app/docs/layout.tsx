import type { ComponentProps } from 'react'
import { Outlet } from '@tanstack/react-router'
import { DocsLayout } from 'fumadocs-ui/layouts/docs'
import { Container } from 'fumadocs-ui/layouts/docs/slots/container'
import { SidebarTrigger } from 'fumadocs-ui/components/sidebar/base'
import { useTranslation } from 'tifa/i18n'
import { MenuIcon } from '@/components/icons'
import { SiteFooter } from '@/components/site-footer'
import { Nav, SiteHeader, menuLink } from '@/components/site-header'
import { layout } from './route.gen'

// The header is the header of every page of the site, above the grid of the docs, so that the
// sidebar and the table of contents start below it (--fd-banner-height). On a narrow screen its
// menu button opens the sidebar, which holds the links of the site above the tree of the docs.
function DocsContainer(props: ComponentProps<'div'>) {
  const { t } = useTranslation()
  return (
    <>
      <SiteHeader
        wide
        menu={
          <SidebarTrigger
            aria-label={t('nav.menu')}
            className="inline-flex size-9 items-center justify-center rounded-md text-fd-muted-foreground transition-colors hover:bg-fd-accent hover:text-fd-accent-foreground"
          >
            <MenuIcon />
          </SidebarTrigger>
        }
      />
      <Container
        {...props}
        style={{
          ...props.style,
          ['--fd-banner-height' as string]: '4rem',
          // The sidebar sits at the left edge. The page and the table of contents are centered in
          // the space to its right, the page at a readable width.
          gridTemplate: `"sidebar . header toc toc"
"sidebar . toc-popover toc toc"
"sidebar . tabs toc toc"
"sidebar . main toc toc" 1fr / var(--fd-sidebar-col) minmax(0, 1fr) minmax(0, 50rem) var(--fd-toc-width) minmax(0, 1fr)`,
        }}
      />
    </>
  )
}

// The sidebar has no logo, search or collapse button of its own, since the header has them.
const none = () => null

export default function DocsFrame() {
  const { tree, version } = layout.useLoaderData()
  const { t } = useTranslation()
  return (
    <>
      <DocsLayout
        tree={tree as never}
        nav={{ enabled: false }}
        tabs={false}
        slots={{ container: DocsContainer, navTitle: none, searchTrigger: false, themeSwitch: false, languageSelect: false }}
        sidebar={{
          collapsible: false,
          banner: (
            <>
              <Nav className="flex flex-col gap-1 border-b border-fd-border pb-3 md:hidden" linkClass={menuLink} />
              <p className="rounded-lg border border-fd-border bg-fd-card px-3 py-2 text-sm text-fd-muted-foreground">
                {t('docs.version', { version })}
              </p>
            </>
          ),
        }}
      >
        <Outlet />
      </DocsLayout>
      <SiteFooter />
    </>
  )
}
