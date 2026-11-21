import { Outlet } from '@tanstack/react-router'
import { DocsLayout } from 'fumadocs-ui/layouts/docs'
import { baseOptions } from '@/components/layout-options'
import { layout } from './route.gen'

export default function DocsFrame() {
  const tree = layout.useLoaderData()
  return (
    <DocsLayout {...baseOptions} tree={tree as never}>
      <Outlet />
    </DocsLayout>
  )
}
