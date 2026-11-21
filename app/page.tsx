import { Link } from '@tanstack/react-router'
import { CodeBlock, Pre } from 'fumadocs-ui/components/codeblock'
import { HomeLayout } from 'fumadocs-ui/layouts/home'
import { useBaseOptions } from '@/components/layout-options'
import { Inline } from '@/components/inline'
import { site } from '@/components/site'
import icon from '@/components/brand/tifa-icon.svg'

// Column widths, so that the long answers of the last column get the room.
const widths = ["w-[26%]", "w-[20%]", ""]

// On a narrow screen a row is a stack, and each answer names its column above itself.
const cell =
  "text-fd-muted-foreground max-sm:block max-sm:py-0 max-sm:pt-2 max-sm:before:mb-0.5 max-sm:before:block max-sm:before:text-xs max-sm:before:font-medium max-sm:before:uppercase max-sm:before:tracking-wide max-sm:before:text-fd-foreground/60 max-sm:before:content-[attr(data-label)]"

export default function HomePage() {
  const options = useBaseOptions()
  return (
    <HomeLayout {...options}>
      <main className="mx-auto w-full max-w-3xl flex-1 px-4 pb-24 pt-16 sm:pt-24">
        <section>
          <img src={icon} alt="" width={72} height={72} className="mb-8 size-[72px]" />
          <h1 className="text-balance text-4xl font-semibold tracking-tight sm:text-5xl">{site.tagline}</h1>
          <div className="mt-8 space-y-5 text-lg leading-8 text-fd-muted-foreground">
            {site.intro.map((p) => (
              <p key={p}>
                <Inline text={p} />
              </p>
            ))}
          </div>
          <div className="mt-10 flex flex-wrap gap-3">
            <Link
              to={site.links.docs}
              className="rounded-lg bg-fd-primary px-5 py-2.5 text-sm font-medium text-white dark:text-black"
            >
              Docs
            </Link>
            <a
              href={site.links.github}
              className="rounded-lg border border-fd-border px-5 py-2.5 text-sm font-medium hover:bg-fd-muted"
            >
              GitHub
            </a>
          </div>
        </section>

        <section className="mt-20">
          <p className="mb-2 font-semibold sm:hidden">{site.vps.header[0]}</p>
          <table className="w-full border-collapse text-left text-sm leading-6">
            <thead className="max-sm:hidden">
              <tr className="border-b border-fd-border">
                {site.vps.header.map((h, i) => (
                  <th key={h} className={'py-3 pr-6 font-semibold ' + widths[i]}>
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {site.vps.rows.map(([what, usual, tifa]) => (
                <tr key={what} className="border-b border-fd-border align-top max-sm:block max-sm:py-4">
                  <th scope="row" className="py-3 pr-6 text-left font-medium max-sm:block max-sm:py-0">
                    {what}
                  </th>
                  <td
                    data-label={site.vps.header[1]}
                    className={cell + ' py-3 pr-6'}
                  >
                    <Inline text={usual} />
                  </td>
                  <td data-label={site.vps.header[2]} className={cell + ' py-3'}>
                    <Inline text={tifa} />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>

        <section className="mt-20">
          <h2 className="text-2xl font-semibold tracking-tight">{site.quickstart.title}</h2>
          <p className="mt-4 leading-7 text-fd-muted-foreground">
            <Inline text={site.quickstart.note} />
          </p>
          <CodeBlock className="mt-6">
            <Pre className="p-4">
              <code>{site.quickstart.commands}</code>
            </Pre>
          </CodeBlock>
        </section>

        <footer className="mt-20 flex flex-wrap gap-x-6 gap-y-2 border-t border-fd-border pt-6 text-sm text-fd-muted-foreground">
          <Link to={site.links.docs} className="hover:text-fd-foreground">
            Docs
          </Link>
          <a href={site.links.github} className="hover:text-fd-foreground">
            GitHub
          </a>
          <a href={site.links.x} className="hover:text-fd-foreground">
            X
          </a>
        </footer>
      </main>
    </HomeLayout>
  )
}
