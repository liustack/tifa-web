import { Outlet } from '@tanstack/react-router'
import { DocsLayout } from 'fumadocs-ui/layouts/docs'
import { useBaseOptions } from '@/components/layout-options'
import { layout } from './route.gen'

export default function DocsFrame() {
  const tree = layout.useLoaderData()
  const options = useBaseOptions()
  return (
    <DocsLayout {...options} tree={tree as never}>
      <Outlet />
    </DocsLayout>
  )
}
