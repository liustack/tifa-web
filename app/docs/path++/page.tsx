import { HTML } from 'tifa/client'
import { DocsBody, DocsDescription, DocsPage, DocsTitle } from 'fumadocs-ui/layouts/docs/page'
import { page } from './route.gen'

export default function DocPage() {
  const { title, description, body, headings } = page.useLoaderData()
  const toc = headings
    .filter((h) => h.depth >= 2 && h.depth <= 3)
    .map((h) => ({ title: h.text, url: `#${h.id}`, depth: h.depth }))
  return (
    <DocsPage toc={toc}>
      <DocsTitle>{title}</DocsTitle>
      <DocsDescription>{description}</DocsDescription>
      <DocsBody>
        <HTML value={body} />
      </DocsBody>
    </DocsPage>
  )
}

export const head = ({ loaderData }: { loaderData?: { title: string; description: string } }) => {
  const title = `${loaderData?.title ?? 'Docs'} | Tifa`
  return {
    meta: [
      { title },
      { name: 'description', content: loaderData?.description },
      { property: 'og:title', content: title },
      { property: 'og:description', content: loaderData?.description },
    ],
  }
}
