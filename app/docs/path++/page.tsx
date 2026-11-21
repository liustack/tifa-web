import { HTML } from 'tifa/client'
import { i18n, useTranslation } from 'tifa/i18n'
import { DocsBody, DocsDescription, DocsPage, DocsTitle, MarkdownCopyButton, ViewOptionsPopover } from 'fumadocs-ui/layouts/docs/page'
import { page } from './route.gen'

export default function DocPage() {
  const { title, description, body, headings, english, markdown, source } = page.useLoaderData()
  const { t } = useTranslation()
  const toc = headings
    .filter((h) => h.depth >= 2 && h.depth <= 3)
    .map((h) => ({ title: h.text, url: `#${h.id}`, depth: h.depth }))
  return (
    <DocsPage toc={toc}>
      <DocsTitle>{title}</DocsTitle>
      <DocsDescription>{description}</DocsDescription>
      <div className="mb-8 flex flex-row flex-wrap items-center gap-2 border-b border-fd-border pb-6">
        <MarkdownCopyButton markdownUrl={markdown}>{t('docs.copyPage')}</MarkdownCopyButton>
        <ViewOptionsPopover markdownUrl={markdown} githubUrl={source} />
      </div>
      {english && <p className="mb-6 rounded-lg border border-fd-border bg-fd-card px-4 py-3 text-sm">{t('docs.englishOnly')}</p>}
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
