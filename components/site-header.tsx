import { useState, type ReactNode } from 'react'
import { Link, useRouterState } from '@tanstack/react-router'
import { FullSearchTrigger, SearchTrigger } from 'fumadocs-ui/layouts/shared/slots/search-trigger'
import { useTranslation } from 'tifa/i18n'
import { CloseIcon, GitHubIcon, MenuIcon } from './icons'
import { LanguageMenu } from './language-menu'
import { site } from './site'
import logo from './brand/tifa-logo.svg'
import logoDark from './brand/tifa-logo-dark.svg'

// The links of the header, by the key of their name in the language files. A link is active on its
// own address and below it, and Docs on every doc.
const links = [
  { key: 'nav.docs', to: site.links.docs, base: '/docs' },
  { key: 'nav.examples', to: '/examples', base: '/examples' },
  { key: 'nav.changelog', to: '/changelog', base: '/changelog' },
  { key: 'nav.why', to: '/why', base: '/why' },
] as const

// useActive tells which link the address is on. The address has the prefix of the language.
function useActive() {
  const path = useRouterState({ select: (s) => s.location.pathname }).replace(/^\/zh-CN(?=\/|$)/, '')
  return (base: string) => path === base || path.startsWith(base + '/')
}

const ghost = 'inline-flex size-9 items-center justify-center rounded-md text-fd-muted-foreground transition-colors hover:bg-fd-accent hover:text-fd-accent-foreground'

// Nav is the links of the site, in the header on a wide screen and in the menu on a narrow one.
export function Nav({ className, linkClass }: { className?: string; linkClass: (active: boolean) => string }) {
  const { t } = useTranslation()
  const active = useActive()
  return (
    <nav className={className}>
      {links.map((l) => (
        <Link key={l.to} to={l.to as never} className={linkClass(active(l.base))}>
          {t(l.key)}
        </Link>
      ))}
    </nav>
  )
}

export const wideLink = (active: boolean) =>
  'rounded-md px-3 py-2 text-sm transition-colors ' + (active ? 'text-fd-foreground' : 'text-fd-muted-foreground hover:text-fd-foreground')

export const menuLink = (active: boolean) =>
  'rounded-md px-3 py-2 text-base ' + (active ? 'bg-fd-accent text-fd-foreground' : 'text-fd-muted-foreground hover:bg-fd-accent')

// SiteHeader is the one header of the site: the logo, the links, search, the language menu and
// GitHub. On a narrow screen the links go behind one menu button. menu is that button: a page
// with a tree of its own, the docs, passes the one that opens the tree, which holds the links
// too. Without it the button opens a panel of the links.
export function SiteHeader({ menu }: { menu?: ReactNode }) {
  const { t } = useTranslation()
  const [open, setOpen] = useState(false)
  const path = useRouterState({ select: (s) => s.location.pathname })
  const [seen, setSeen] = useState(path)
  if (seen !== path) {
    setSeen(path)
    setOpen(false)
  }
  return (
    <>
      <header className="sticky top-0 z-40 h-16 border-b border-fd-border bg-fd-background/80 backdrop-blur-sm">
        <div className="mx-auto flex h-full max-w-[97rem] items-center gap-2 px-4 md:px-6">
          <Link to="/" className="me-4 inline-flex shrink-0 items-center" aria-label={site.name}>
            <img src={logo} alt="" height={28} className="h-7 w-auto dark:hidden" />
            <img src={logoDark} alt="" height={28} className="hidden h-7 w-auto dark:block" />
          </Link>
          <Nav className="flex items-center gap-1 max-md:hidden" linkClass={wideLink} />
          <div className="ms-auto flex items-center gap-1">
            <FullSearchTrigger hideIfDisabled className="me-1 h-9 w-56 max-lg:hidden" />
            <SearchTrigger hideIfDisabled className="size-9 lg:hidden" />
            <LanguageMenu />
            <a href={site.links.github} aria-label="GitHub" className={ghost}>
              <GitHubIcon />
            </a>
            <div className="md:hidden">
              {menu ?? (
                <button type="button" aria-label={t('nav.menu')} aria-expanded={open} onClick={() => setOpen(!open)} className={ghost}>
                  {open ? <CloseIcon /> : <MenuIcon />}
                </button>
              )}
            </div>
          </div>
        </div>
      </header>
      {!menu && open && (
        <div className="border-b border-fd-border bg-fd-background px-4 py-3 md:hidden">
          <Nav className="flex flex-col gap-1" linkClass={menuLink} />
        </div>
      )}
    </>
  )
}
