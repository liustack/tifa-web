import { Outlet } from '@tanstack/react-router'
import { DocsLayout } from 'fumadocs-ui/layouts/docs'
import { useTranslation } from 'tifa/i18n'
import { useBaseOptions } from '@/components/layout-options'
import { layout } from './route.gen'

export default function DocsFrame() {
  const { tree, version } = layout.useLoaderData()
  const { t } = useTranslation()
  const options = useBaseOptions()
  return (
    <DocsLayout
      {...options}
      tree={tree as never}
      sidebar={{
        banner: (
          <p className="rounded-lg border border-fd-border bg-fd-card px-3 py-2 text-sm text-fd-muted-foreground">
            {t('docs.version', { version })}
          </p>
        ),
      }}
    >
      <Outlet />
    </DocsLayout>
  )
}
