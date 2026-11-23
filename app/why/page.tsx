import { Prose, ProseTitle, proseHead } from '@/components/prose'
import { ProseLayout } from '@/components/prose-layout'
import { page } from './route.gen'

export default function WhyPage() {
  const { title, body, philosophy } = page.useLoaderData()
  return (
    <ProseLayout>
      <ProseTitle>{title}</ProseTitle>
      <Prose body={body} />
      <h2 className="mb-6 mt-14 text-2xl font-semibold tracking-tight">{philosophy.title}</h2>
      <Prose body={philosophy.body} />
    </ProseLayout>
  )
}

export const head = ({ loaderData }: { loaderData?: { title: string; description: string } }) => proseHead(loaderData)
