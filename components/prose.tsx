import { HTML, type HTMLProps } from 'tifa/client'
import { DocsBody } from 'fumadocs-ui/layouts/docs/page'
import { useTranslation } from 'tifa/i18n'

const noteClass = 'mb-6 rounded-lg border border-fd-border bg-fd-card px-4 py-3 text-sm'

// DocNote says that the body of a page is English on the Chinese site.
export function DocNote({ english }: { english?: boolean }) {
  const { t } = useTranslation()
  if (english) return <p className={noteClass}>{t('docs.englishOnly')}</p>
  return null
}

// Prose is the text of a page of Markdown, in the style of a doc. english adds the note
// that says so.
export function Prose({ body, english }: { body: HTMLProps['value']; english?: boolean }) {
  return (
    <>
      <DocNote english={english} />
      <DocsBody>
        <HTML value={body} />
      </DocsBody>
    </>
  )
}

// ProseTitle is the title of a page of Markdown.
export function ProseTitle({ children }: { children: string }) {
  return <h1 className="mb-8 text-balance text-3xl font-semibold tracking-tight sm:text-4xl">{children}</h1>
}

// proseHead is the head of a page of Markdown: its title and its description.
export function proseHead(data?: { title: string; description: string }) {
  if (!data) return {}
  const title = `${data.title} | Tifa`
  return {
    meta: [
      { title },
      { name: 'description', content: data.description },
      { property: 'og:title', content: title },
      { property: 'og:description', content: data.description },
    ],
  }
}
