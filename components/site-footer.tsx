import { Link } from '@tanstack/react-router'
import { useTranslation } from 'tifa/i18n'
import { site } from './site'

// SiteFooter is the links at the bottom of every page that has the header of the site.
export function SiteFooter({ wide }: { wide?: boolean }) {
  const { t } = useTranslation()
  const link = 'hover:text-fd-foreground'
  return (
    <footer className={`mx-auto mb-12 flex w-full ${wide ? 'max-w-5xl' : 'max-w-3xl'} flex-wrap gap-x-6 gap-y-2 border-t border-fd-border px-4 pt-6 text-sm text-fd-muted-foreground`}>
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
    </footer>
  )
}
