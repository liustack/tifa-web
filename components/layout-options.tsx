import type { BaseLayoutProps } from 'fumadocs-ui/layouts/shared'
import { site } from './site'
import logo from './brand/tifa-logo.svg'
import logoDark from './brand/tifa-logo-dark.svg'

function XIcon() {
  return (
    <svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
      <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
    </svg>
  )
}

// The header both layouts share: the logo, and the links to GitHub and X next to the search
// and the theme switch.
export const baseOptions: BaseLayoutProps = {
  nav: {
    title: (
      <>
        <img src={logo} alt="Tifa" height={28} className="h-7 w-auto dark:hidden" />
        <img src={logoDark} alt="Tifa" height={28} className="hidden h-7 w-auto dark:block" />
      </>
    ),
    url: '/',
  },
  githubUrl: site.links.github,
  links: [
    { type: 'icon', text: 'X', label: 'X', url: site.links.x, icon: <XIcon />, external: true },
  ],
}
