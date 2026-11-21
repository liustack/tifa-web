import type { BaseLayoutProps } from 'fumadocs-ui/layouts/shared'
import { useTranslation } from 'tifa/i18n'
import { site } from './site'
import { LanguageSwitch } from './language-switch'
import logo from './brand/tifa-logo.svg'
import logoDark from './brand/tifa-logo-dark.svg'

// useBaseOptions gives the header both layouts share: the logo, the links of the site, GitHub,
// the language switch, and the search and theme switch that Fumadocs adds.
export function useBaseOptions(): BaseLayoutProps {
  const { t } = useTranslation()
  return {
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
      { type: 'main', text: t('nav.docs'), url: site.links.docs, active: 'nested-url' },
      { type: 'main', text: t('nav.examples'), url: '/examples' },
      { type: 'main', text: t('nav.changelog'), url: '/changelog' },
      { type: 'main', text: t('nav.why'), url: '/why' },
      { type: 'custom', secondary: true, children: <LanguageSwitch label={t('language.label')} /> },
    ],
  }
}
