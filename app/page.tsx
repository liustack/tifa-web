import { Link } from '@tanstack/react-router'
import { CodeBlock, Pre } from 'fumadocs-ui/components/codeblock'
import { HomeLayout } from 'fumadocs-ui/layouts/home'
import { HTML } from 'tifa/client'
import { useTranslation } from 'tifa/i18n'
import { useBaseOptions } from '@/components/layout-options'
import { Inline } from '@/components/inline'
import { SiteFooter } from '@/components/site-footer'
import { referenceRows, site } from '@/components/site'
import icon from '@/components/brand/tifa-icon.svg'
import { page } from './route.gen'

// The words of the page are in app/_locales, under home and hero. t returns an array or an object
// for a key that holds one.
function useList<T>(key: string): T {
  const { t } = useTranslation()
  return t(key, { returnObjects: true }) as T
}

// Column widths of the table of the VPS, so that the long answers of the last column get the room.
const widths = ['w-[26%]', 'w-[20%]', '']

// On a narrow screen a row is a stack, and each answer names its column above itself.
const cell =
  'text-fd-muted-foreground max-sm:block max-sm:py-0 max-sm:pt-2 max-sm:before:mb-0.5 max-sm:before:block max-sm:before:text-xs max-sm:before:font-medium max-sm:before:uppercase max-sm:before:tracking-wide max-sm:before:text-fd-foreground/60 max-sm:before:content-[attr(data-label)]'

const h2 = 'text-2xl font-semibold tracking-tight'
const button = 'rounded-lg px-5 py-2.5 text-sm font-medium'

export default function HomePage() {
  const { version, demo } = page.useLoaderData()
  const { t } = useTranslation()
  const options = useBaseOptions()
  const intro = useList<string[]>('hero.intro')
  const vps = useList<{ header: string[]; rows: string[][] }>('home.vps')
  const numbers = useList<{ columns: string[]; rows: Record<string, string>; cost: string[] }>('home.numbers')
  const help = useList<string[]>('home.status.help')
  const doors = useList<{ to: string; label: string; text: string }[]>('home.doors.items')
  const why = useList<string[]>('home.why.text')
  return (
    <HomeLayout {...options}>
      <main className="mx-auto w-full max-w-5xl flex-1 px-4 pb-24 pt-16 sm:pt-24">
        <section className="max-w-3xl">
          <img src={icon} alt="" width={72} height={72} className="mb-8 size-[72px]" />
          <h1 className="text-balance text-4xl font-semibold tracking-tight sm:text-5xl">{t('hero.tagline')}</h1>
          <div className="mt-8 space-y-5 text-lg leading-8 text-fd-muted-foreground">
            {intro.map((p) => (
              <p key={p}>{p}</p>
            ))}
          </div>
          <p className="mt-8 inline-flex flex-wrap items-center gap-x-2 gap-y-1 rounded-2xl border border-fd-border bg-fd-card px-3 py-1.5 text-sm text-fd-muted-foreground">
            <span className="rounded-full bg-fd-primary/10 px-2 py-0.5 font-mono text-xs text-fd-primary">{version}</span>
            <span>{t('hero.status')}</span>
          </p>
          <CodeBlock className="mt-6" title={t('hero.command')}>
            <Pre className="p-4">
              <code>{site.newApp}</code>
            </Pre>
          </CodeBlock>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link to={site.links.docs} className={`${button} bg-fd-primary text-white dark:text-black`}>
              {t('hero.start')}
            </Link>
            <a href={site.links.github} className={`${button} border border-fd-border hover:bg-fd-muted`}>
              {t('hero.github')}
            </a>
          </div>
        </section>

        <section className="mt-24">
          <h2 className={h2}>{t('home.demo.title')}</h2>
          <p className="mt-4 leading-7 text-fd-muted-foreground">{t('home.demo.intro')}</p>
          <div className="mt-6 grid min-w-0 gap-x-6 lg:grid-cols-2">
            {demo?.map((f) => (
              <figure key={f.path} className="min-w-0">
                <figcaption className="mt-2 font-mono text-sm">
                  <a href={f.url} className="text-fd-muted-foreground hover:text-fd-foreground">
                    {f.path}
                  </a>
                </figcaption>
                <HTML value={f.body} className="[&_pre.shiki]:mt-2" />
              </figure>
            ))}
          </div>
          <p className="text-sm text-fd-muted-foreground">
            <Inline text={t('home.demo.from')} />
          </p>
        </section>

        <section className="mt-24">
          <p className="mb-2 font-semibold sm:hidden">{vps.header[0]}</p>
          <table className="w-full border-collapse text-left text-sm leading-6">
            <thead className="max-sm:hidden">
              <tr className="border-b border-fd-border">
                {vps.header.map((h, i) => (
                  <th key={h} className={'py-3 pr-6 font-semibold ' + widths[i]}>
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {vps.rows.map(([what, usual, tifa]) => (
                <tr key={what} className="border-b border-fd-border align-top max-sm:block max-sm:py-4">
                  <th scope="row" className="py-3 pr-6 text-left font-medium max-sm:block max-sm:py-0">
                    {what}
                  </th>
                  <td data-label={vps.header[1]} className={cell + ' py-3 pr-6'}>
                    <Inline text={usual} />
                  </td>
                  <td data-label={vps.header[2]} className={cell + ' py-3'}>
                    <Inline text={tifa} />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>

        <section className="mt-24">
          <h2 className={h2}>{t('home.numbers.title')}</h2>
          <div className="mt-4 max-w-3xl space-y-4 leading-7 text-fd-muted-foreground">
            <p>
              <Inline text={t('home.numbers.memo')} />
            </p>
            <p>{t('home.numbers.visits')}</p>
            <p>{t('home.numbers.caption')}</p>
          </div>
          <div className="mt-4 overflow-x-auto">
            <table className="w-full border-collapse text-left text-sm leading-6">
              <thead className="max-sm:hidden">
                <tr className="border-b border-fd-border">
                  <th className="py-3 pr-4" />
                  {numbers.columns.map((h) => (
                    <th key={h} className="px-4 py-3 text-right font-semibold">
                      {h}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {site.throughput.map(([key, ...values]) => (
                  <tr key={key} className={'border-b border-fd-border max-sm:block max-sm:py-3 ' + (referenceRows.has(key) ? 'text-fd-muted-foreground' : 'font-medium')}>
                    <th scope="row" className="py-3 pr-4 text-left font-[inherit] max-sm:block max-sm:py-1">
                      {numbers.rows[key]}
                    </th>
                    {values.map((v, i) => (
                      <td
                        key={i}
                        data-label={numbers.columns[i]}
                        className="px-4 py-3 text-right font-mono tabular-nums max-sm:flex max-sm:justify-between max-sm:px-0 max-sm:py-0.5 max-sm:before:font-sans max-sm:before:text-fd-muted-foreground max-sm:before:content-[attr(data-label)]"
                      >
                        {v.toFixed(1)}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="mt-4 max-w-3xl space-y-4 text-sm leading-6 text-fd-muted-foreground">
            <p>{t('home.numbers.definition')}</p>
            <p>
              <Inline text={t('home.numbers.method')} />
            </p>
          </div>
          <div className="mt-6 max-w-3xl space-y-3 rounded-lg border border-fd-border bg-fd-card p-4 text-sm leading-6">
            {numbers.cost.map((p) => (
              <p key={p}>
                <Inline text={p} />
              </p>
            ))}
          </div>
          <p className="mt-4 flex flex-wrap gap-x-6 gap-y-2 text-sm">
            <Link to={'/docs/reference/performance' as never} className="underline underline-offset-4">
              {t('home.numbers.performance')}
            </Link>
            <a href={t('home.numbers.capacityUrl')} className="underline underline-offset-4">
              {t('home.numbers.capacity')}
            </a>
          </p>
        </section>

        <section className="mt-24 max-w-3xl">
          <h2 className={h2}>{t('home.status.title')}</h2>
          <p className="mt-4 border-l-2 border-fd-primary pl-4 leading-7">{t('home.status.note')}</p>
          <p className="mt-4 leading-7 text-fd-muted-foreground">{t('home.status.policy')}</p>
          <ul className="mt-3 list-disc space-y-2 pl-6 leading-7 text-fd-muted-foreground">
            {help.map((p) => (
              <li key={p}>
                <Inline text={p} />
              </li>
            ))}
          </ul>
        </section>

        <section className="mt-24">
          <h2 className={h2}>{t('home.doors.title')}</h2>
          <ul className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {doors.map((d) => {
              const [to, hash] = d.to.split('#')
              return (
                <li key={d.to}>
                  <Link
                    to={to as never}
                    hash={hash}
                    className="flex h-full flex-col gap-2 rounded-lg border border-fd-border p-4 transition-colors hover:bg-fd-muted"
                  >
                    <span className="font-medium">{d.label}</span>
                    <span className="text-sm leading-6 text-fd-muted-foreground">
                      <Inline text={d.text} />
                    </span>
                  </Link>
                </li>
              )
            })}
          </ul>
        </section>

        <section className="mt-24 max-w-3xl">
          <h2 className={h2}>{t('home.why.title')}</h2>
          <div className="mt-4 space-y-4 leading-7 text-fd-muted-foreground">
            {why.map((p) => (
              <p key={p}>{p}</p>
            ))}
          </div>
          <p className="mt-4">
            <Link to="/why" className="underline underline-offset-4">
              {t('home.why.more')}
            </Link>
          </p>
        </section>

        <section className="mt-24 max-w-3xl">
          <h2 className={h2}>{t('home.start.title')}</h2>
          <p className="mt-4 leading-7 text-fd-muted-foreground">
            <Inline text={t('home.start.note')} />
          </p>
          <CodeBlock className="mt-6">
            <Pre className="p-4">
              <code>{site.quickstart}</code>
            </Pre>
          </CodeBlock>
        </section>
      </main>
      <SiteFooter wide />
    </HomeLayout>
  )
}
