import { Outlet } from '@tanstack/react-router'
import { TifaError, useLanguage } from 'tifa/client'
import { I18nProvider, i18n, useTranslation } from 'tifa/i18n'
import { HomeLayout } from 'fumadocs-ui/layouts/home'
import { RootProvider } from 'fumadocs-ui/provider/tanstack'
import { Link } from '@tanstack/react-router'
import { useBaseOptions } from '@/components/layout-options'
import { site } from '@/components/site'
import icon from '@/components/brand/tifa-icon.svg'
import './styles.css'

// Head tags of every page, also of the 404 page, that do not depend on the language. The title
// and the description are in head below, and a page's head can set its own.
export const documentHead = {
  meta: [
    { property: 'og:site_name', content: site.name },
    { property: 'og:type', content: 'website' },
    { property: 'og:image', content: `${site.url}/tifa-og.png` },
    { name: 'twitter:card', content: 'summary_large_image' },
    { name: 'twitter:site', content: '@liustack' },
  ],
  links: [{ rel: 'icon', href: icon, type: 'image/svg+xml' }],
}

// The title and description of a page that sets none of its own.
export const head = () => {
  const title = `${site.name} | ${i18n.t('hero.tagline')}`
  const description = i18n.t('hero.intro.0')
  return {
    meta: [
      { title },
      { name: 'description', content: description },
      { property: 'og:title', content: title },
      { property: 'og:description', content: description },
    ],
  }
}

function Providers({ children }: { children: React.ReactNode }) {
  const { t } = useTranslation()
  const language = useLanguage()
  // Fumadocs finds its own texts (search, theme, copy buttons) by their English wording. English
  // needs no table.
  const translations = t('fumadocs', { returnObjects: true }) as Record<string, string>
  return (
    // The search address is written out because Fumadocs builds the default from Vite's base,
    // which is /__tifa/vite/ under tifa dev. In Chinese it is the route under /zh-CN, so that
    // search answers in that language.
    <RootProvider
      i18n={{ locale: language, translations }}
      search={{ options: { api: language === 'en-US' ? '/api/search' : `/${language}/api/search` } }}
    >
      {children}
    </RootProvider>
  )
}

export default function Root() {
  return (
    <I18nProvider>
      <Providers>
        <Outlet />
      </Providers>
    </I18nProvider>
  )
}

// Error shows inside this Layout when a page fails to load or render and has no Error of its
// own, and for a URL that matches no page, as a TifaError with status 404.
export function ErrorBoundary({ error }: { error: unknown }) {
  const { t } = useTranslation()
  const options = useBaseOptions()
  const notFound = error instanceof TifaError && error.status === 404
  return (
    <HomeLayout {...options}>
      <main className="mx-auto flex w-full max-w-xl flex-1 flex-col justify-center gap-4 px-4 py-24">
        <h1 className="text-2xl font-semibold">{notFound ? t('error.notFound') : t('error.failed')}</h1>
        <p>
          <Link to="/" className="underline underline-offset-4">
            {t('error.home')}
          </Link>
        </p>
      </main>
    </HomeLayout>
  )
}
