import { useEffect, useState } from 'react'
import { HTML } from 'tifa/client'
import { i18n, useTranslation } from 'tifa/i18n'
import { DocsBody, DocsDescription, DocsPage, DocsTitle, MarkdownCopyButton, ViewOptionsPopover } from 'fumadocs-ui/layouts/docs/page'
import { DocNote } from '@/components/prose'
import { site } from '@/components/site'
import { page } from './route.gen'

// The page actions of Fumadocs put the base of the build (/__tifa/vite/ under tifa dev) before a
// path that starts with /, and leave an address with a scheme alone. So they get the address in
// full: the site's own at first, which is what the server renders, and the origin the browser
// is on once the page is there.
function useOrigin() {
  const [origin, setOrigin] = useState(site.url)
  useEffect(() => setOrigin(window.location.origin), [])
  return origin
}

export default function DocPage() {
  const { title, description, body, headings, english, outdated, englishUrl, markdown, source } = page.useLoaderData()
  const { t } = useTranslation()
  const origin = useOrigin()
  const toc = headings
    .filter((h) => h.depth >= 2 && h.depth <= 3)
    .map((h) => ({ title: h.text, url: `#${h.id}`, depth: h.depth }))
  return (
    <DocsPage toc={toc}>
      <DocsTitle>{title}</DocsTitle>
      <DocsDescription>{description}</DocsDescription>
      <div className="mb-8 flex flex-row flex-wrap items-center gap-2 border-b border-fd-border pb-6">
        <MarkdownCopyButton markdownUrl={origin + markdown}>{t('docs.copyPage')}</MarkdownCopyButton>
        <ViewOptionsPopover markdownUrl={origin + markdown} githubUrl={source} />
      </div>
      <DocNote english={english} outdated={outdated} englishUrl={englishUrl} />
      <DocsBody>
        <HTML value={body} />
      </DocsBody>
    </DocsPage>
  )
}

export const head = ({ loaderData }: { loaderData?: { title: string; description: string } }) => {
  const title = `${loaderData?.title ?? i18n.t('nav.docs')} | Tifa`
  return {
    meta: [
      { title },
      { name: 'description', content: loaderData?.description },
      { property: 'og:title', content: title },
      { property: 'og:description', content: loaderData?.description },
    ],
  }
}
