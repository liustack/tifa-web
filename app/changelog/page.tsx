import { Prose, ProseTitle, proseHead } from '@/components/prose'
import { ProseLayout } from '@/components/prose-layout'
import { page } from './route.gen'

export default function ChangelogPage() {
  const { title, body, english } = page.useLoaderData()
  return (
    <ProseLayout>
      <ProseTitle>{title}</ProseTitle>
      <Prose body={body} english={english} />
    </ProseLayout>
  )
}

export const head = ({ loaderData }: { loaderData?: { title: string; description: string } }) => proseHead(loaderData)
