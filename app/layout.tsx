import { Link, Outlet } from '@tanstack/react-router'
import { TifaError } from 'tifa/client'
import { RootProvider } from 'fumadocs-ui/provider/tanstack'
import { site } from '@/components/site'
import icon from '@/components/brand/tifa-icon.svg'
import './styles.css'

export const htmlAttributes = { lang: 'en' }

// Head tags of every page, also of the 404 page. A page's head can set its own title and
// description.
export const documentHead = {
  meta: [
    { title: `${site.name} | ${site.tagline}` },
    { name: 'description', content: site.intro[0] },
    { property: 'og:site_name', content: site.name },
    { property: 'og:type', content: 'website' },
    { property: 'og:title', content: `${site.name} | ${site.tagline}` },
    { property: 'og:description', content: site.intro[0] },
    { property: 'og:image', content: `${site.url}/tifa-og.png` },
    { name: 'twitter:card', content: 'summary_large_image' },
    { name: 'twitter:site', content: '@liustack' },
  ],
  links: [{ rel: 'icon', href: icon, type: 'image/svg+xml' }],
}

export default function Root() {
  return (
    // The address is written out because Fumadocs builds the default from Vite's base, which
    // is /__tifa/vite/ under tifa dev.
    <RootProvider search={{ options: { api: '/api/search' } }}>
      <Outlet />
    </RootProvider>
  )
}

// Error shows inside this Layout when a page fails to load or render and has no Error of its
// own, and for a URL that matches no page, as a TifaError with status 404.
export function ErrorBoundary({ error }: { error: unknown }) {
  const notFound = error instanceof TifaError && error.status === 404
  return (
    <main className="mx-auto flex min-h-screen max-w-xl flex-col justify-center gap-4 p-8">
      <h1 className="text-2xl font-semibold">{notFound ? 'Page not found' : 'Something went wrong'}</h1>
      <p>
        <Link to="/" className="underline underline-offset-4">
          Back to the home page
        </Link>
      </p>
    </main>
  )
}
