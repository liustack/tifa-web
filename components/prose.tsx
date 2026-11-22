import { HTML, type HTMLProps } from 'tifa/client'
import { DocsBody } from 'fumadocs-ui/layouts/docs/page'
import { useTranslation } from 'tifa/i18n'

const noteClass = 'mb-6 rounded-lg border border-fd-border bg-fd-card px-4 py-3 text-sm'

// DocNote says what is wrong with the language of a doc: its body is English on the Chinese
// site, or a translation made before the English doc changed. The link goes to the English page,
// which is another document, so it is a plain <a>.
export function DocNote({ english, outdated, englishUrl }: { english?: boolean; outdated?: boolean; englishUrl?: string }) {
  const { t } = useTranslation()
  if (english) return <p className={noteClass}>{t('docs.englishOnly')}</p>
  if (outdated) {
    return (
      <p className={noteClass}>
        {t('docs.outdated')}{' '}
        <a className="font-medium underline underline-offset-4" href={englishUrl}>
          {t('docs.outdatedLink')}
        </a>
      </p>
    )
  }
  return null
}

// Prose is the text of a page of Markdown, in the style of a doc. english and outdated add the
// note that says so.
export function Prose({ body, english, outdated, englishUrl }: { body: HTMLProps['value']; english?: boolean; outdated?: boolean; englishUrl?: string }) {
  return (
    <>
      <DocNote english={english} outdated={outdated} englishUrl={englishUrl} />
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
