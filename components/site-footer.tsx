import { Link } from '@tanstack/react-router'
import { useTranslation } from 'tifa/i18n'
import { site } from './site'
import { ThemeSwitch } from './theme-switch'

// SiteFooter is the links at the bottom of every page, with the theme switch at its end.
export function SiteFooter() {
  const { t } = useTranslation()
  const link = 'hover:text-fd-foreground'
  return (
    <footer className="border-t border-fd-border">
      <div className="mx-auto flex w-full max-w-[97rem] flex-wrap items-center gap-x-6 gap-y-3 px-4 py-6 text-sm text-fd-muted-foreground md:px-6">
        <Link to={site.links.docs} className={link}>
          {t('nav.docs')}
        </Link>
        <Link to="/examples" className={link}>
          {t('nav.examples')}
        </Link>
        <Link to="/changelog" className={link}>
          {t('nav.changelog')}
        </Link>
        <Link to="/why" className={link}>
          {t('nav.why')}
        </Link>
        <a href={site.links.github} className={link}>
          GitHub
        </a>
        <a href={site.links.x} className={link}>
          X
        </a>
        <a href={site.links.source} className={link}>
          {t('footer.source')}
        </a>
        <span>{t('footer.license')}</span>
        <div className="ms-auto">
          <ThemeSwitch />
        </div>
      </div>
    </footer>
  )
}
