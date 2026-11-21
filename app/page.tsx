import { useEffect, type ReactNode } from 'react'
import { Link } from '@tanstack/react-router'
import { useTranslation } from 'tifa/i18n'
import { CopyButton } from '@/components/copy-button'
import { ArrowRightIcon, BranchIcon, CheckIcon, DownloadIcon, FileIcon, GitHubIcon, ServerIcon } from '@/components/icons'
import { Inline } from '@/components/inline'
import { pageTsx, routeGo } from '@/components/landing-code'
import { SiteShell } from '@/components/site-shell'
import { deployTerminal, heroTerminal, numbers, site } from '@/components/site'
import { Terminal } from '@/components/terminal'
import { page } from './route.gen'
import './landing.css'

// The words of the page are in app/_locales, under hero and home. t returns an array or an object
// for a key that holds one.
function useList<T>(key: string): T {
  const { t } = useTranslation()
  return t(key, { returnObjects: true }) as T
}

// useReveal fades in what is below the fold as it comes into view. The page is complete without
// it: the server writes everything visible, and nothing is hidden until the browser has the page
// and the visitor has not asked for less motion.
function useReveal() {
  useEffect(() => {
    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches || !('IntersectionObserver' in window)) return
    const io = new IntersectionObserver(
      (entries) => {
        for (const e of entries) {
          if (e.isIntersecting) {
            e.target.classList.add('in')
            io.unobserve(e.target)
          }
        }
      },
      { rootMargin: '0px 0px -8% 0px' },
    )
    for (const el of document.querySelectorAll<HTMLElement>('.lp [data-reveal]')) {
      if (el.getBoundingClientRect().top > window.innerHeight) {
        el.classList.add('rv')
        io.observe(el)
      }
    }
    return () => io.disconnect()
  }, [])
}

// headlineLines breaks the headline after its first comma, the full-width one of the Chinese
// headline too, so that it never breaks inside a word. Each line is one unit (see .h1-line).
function headlineLines(text: string): string[] {
  const at = text.search(/[，,]/)
  return at < 0 ? [text] : [text.slice(0, at + 1), text.slice(at + 1).trim()]
}

function SectionHead({ rail, children, no }: { rail: ReactNode; children: ReactNode; no: string }) {
  return (
    <div className="sec-head" data-reveal>
      <div className="sec-no">
        <b>{no}</b>
        {rail}
      </div>
      <div>{children}</div>
    </div>
  )
}

function Install({ margin }: { margin?: boolean }) {
  return (
    <div className="install" style={margin === false ? { marginTop: 0 } : undefined}>
      <code>
        <span className="dollar">$</span>
        {site.newApp}
      </code>
      <CopyButton text={site.newApp} />
    </div>
  )
}

export default function HomePage() {
  const { version } = page.useLoaderData()
  const { t } = useTranslation()
  useReveal()
  const stack = useList<{ rows: string[][] }>('home.stack')
  const cost = useList<string[]>('home.numbers.cost')
  const notes = useList<string[][]>('home.pages.notes')
  const stages = useList<{ node: string; where: string; text: string }[]>('home.deploy.stages')
  const facts = useList<string[]>('home.deploy.facts')
  const story = useList<string[]>('home.story.text')
  const board = useList<string[][]>('home.status.rows')
  const stageIcons = [<BranchIcon key="b" />, <CheckIcon key="c" />, <DownloadIcon key="d" />, <ServerIcon key="s" />]
  const columns = [t('home.numbers.tifa'), t('home.numbers.tanstack'), t('home.numbers.next')]
  return (
    <SiteShell>
      <main className="lp flex-1">
        <section className="hero">
          <div className="wrap hero-grid">
            <div>
              <p className="eyebrow">{t('hero.eyebrow')}</p>
              <h1>
                {headlineLines(t('hero.tagline')).map((line, _, lines) =>
                  lines.length > 1 ? (
                    <span key={line} className="h1-line">
                      {line}
                    </span>
                  ) : (
                    line
                  ),
                )}
              </h1>
              <p className="sub">
                <Inline text={t('hero.sub')} />
              </p>
              <Install />
              <div className="actions">
                <Link to={site.links.docs} className="btn btn-primary">
                  {t('hero.start')}
                </Link>
                <a className="btn btn-quiet" href={site.links.github}>
                  <GitHubIcon width={17} height={17} />
                  {t('hero.github')}
                </a>
              </div>
              <p className="status-line">
                <span className="tag">{version}</span>
                {t('hero.status')}
              </p>
            </div>
            <Terminal title="~/code" label={t('hero.terminal')} lines={heroTerminal} type cursor />
          </div>
        </section>

        <section className="section" id="stack">
          <div className="wrap">
            <SectionHead no="01" rail={t('home.stack.rail')}>
              <h2 className="h2">{t('home.stack.title')}</h2>
              <p className="lead">{t('home.stack.lead')}</p>
            </SectionHead>
            <div className="servers" data-reveal>
              <div className="server">
                <div className="server-head">
                  <b>{t('home.stack.usual')}</b>
                  <span>{t('home.stack.vps')}</span>
                </div>
                <div className="slots">
                  {stack.rows.map(([role, usual], i) => (
                    <div key={usual} className={'chip' + (i === stack.rows.length - 1 ? ' lib' : '')}>
                      <div className="n">{usual}</div>
                      <div className="r">{role}</div>
                    </div>
                  ))}
                </div>
              </div>
              <div className="arrow-col" aria-hidden="true">
                <ArrowRightIcon className="pulse" />
              </div>
              <div className="server after">
                <div className="server-head">
                  <b>{t('home.stack.with')}</b>
                  <span>{t('home.stack.vps')}</span>
                </div>
                <div className="binary">
                  <div className="binary-top">
                    <FileIcon />
                    {t('home.stack.binary')}
                    <span>{t('home.stack.one')}</span>
                  </div>
                  <ul>
                    {stack.rows.map(([, usual, tifa]) => (
                      <li key={usual}>
                        <span className="k">{usual}</span>
                        <span>
                          <Inline text={tifa} />
                        </span>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section className="section" id="numbers">
          <div className="wrap">
            <SectionHead no="02" rail={t('home.numbers.rail')}>
              <h2 className="h2">{t('home.numbers.title')}</h2>
              <p className="lead">{t('home.numbers.lead')}</p>
            </SectionHead>
            <div className="spec" data-reveal>
              {numbers.map((n) => {
                const max = Math.max(...n.values)
                return (
                  <div className="spec-row" key={n.key}>
                    <div className="label">
                      {t(`home.numbers.${n.key}`)}
                      <small>{t(`home.numbers.${n.key === 'views' ? 'views' : n.key}_hint`)}</small>
                    </div>
                    <div className="num">
                      {n.big}
                      {n.unit && <span>{n.unit}</span>}
                    </div>
                    <div className="cmp">
                      {n.values.map((v, i) => (
                        <div key={i} className={i === 0 ? 'us' : undefined}>
                          <span>{columns[i]}</span>
                          <span className="bar">
                            <i style={{ ['--w' as string]: `${Math.round((v / max) * 1000) / 10}%` }} />
                          </span>
                          <span className="v">{n.text[i]}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )
              })}
            </div>
            <div className="spec-foot" data-reveal>
              <div>
                <p className="caption">
                  <Inline text={t('home.numbers.caption')} />
                </p>
                <p style={{ margin: '20px 0 0' }}>
                  <Link to={site.links.performance as never} className="link-arrow">
                    {t('home.numbers.method')}
                  </Link>
                </p>
              </div>
              <div className="cost">
                <h3>{t('home.numbers.costTitle')}</h3>
                {cost.map((p) => (
                  <p key={p}>
                    <Inline text={p} />
                  </p>
                ))}
              </div>
            </div>
          </div>
        </section>

        <section className="section" id="code">
          <div className="wrap">
            <SectionHead no="03" rail={t('home.pages.rail')}>
              <h2 className="h2">{t('home.pages.title')}</h2>
              <p className="lead">
                <Inline text={t('home.pages.lead')} />
              </p>
            </SectionHead>
            <div className="editor win" data-reveal>
              <div className="tabs">
                <div>
                  <span className="dot" />
                  app/posts/post+/route.go<span className="lang">Go</span>
                </div>
                <div>
                  <span className="dot" />
                  app/posts/post+/page.tsx<span className="lang">{t('home.pages.react')}</span>
                </div>
              </div>
              <div className="split">
                <pre>
                  <code dangerouslySetInnerHTML={{ __html: routeGo }} />
                </pre>
                <pre data-file="app/posts/post+/page.tsx">
                  <code dangerouslySetInnerHTML={{ __html: pageTsx }} />
                </pre>
              </div>
            </div>
            <div className="under" data-reveal>
              <span className="legend">
                <i />
                {t('home.pages.legend')}
              </span>
              <Link to={site.links.docs} className="link-arrow">
                {t('home.pages.tutorial')}
              </Link>
            </div>
            <div className="notes" data-reveal>
              {notes.map(([title, text]) => (
                <div key={title}>
                  <h3>{title}</h3>
                  <p>
                    <Inline text={text} />
                  </p>
                </div>
              ))}
            </div>
          </div>
        </section>

        <section className="section" id="deploy">
          <div className="wrap">
            <SectionHead no="04" rail={t('home.deploy.rail')}>
              <h2 className="h2">{t('home.deploy.title')}</h2>
              <p className="lead">
                <Inline text={t('home.deploy.lead')} />
              </p>
            </SectionHead>
            <div className="pipe" data-reveal>
              {stages.map((s, i) => (
                <div key={s.node} className={'stage' + (i === stages.length - 1 ? ' live' : '')}>
                  <div className="node">
                    {stageIcons[i]}
                    {s.node}
                  </div>
                  {i < stages.length - 1 && <span className="conn" />}
                  <span className="st">{s.where}</span>
                  <p>
                    <Inline text={s.text} />
                  </p>
                </div>
              ))}
            </div>
            <div className="deploy-grid" data-reveal>
              <Terminal title="tifa-web" label={t('home.deploy.terminal')} lines={deployTerminal} />
              <ul className="facts">
                {facts.map((f) => (
                  <li key={f}>
                    <Inline text={f} />
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </section>

        <section className="section" id="story">
          <div className="wrap">
            <SectionHead no="05" rail={t('home.story.rail')}>
              <h2 className="h2">{t('home.story.title')}</h2>
            </SectionHead>
            <div className="story-body" data-reveal>
              <p>{story[0]}</p>
              <p className="punch">{story[1]}</p>
              <p>{story[2]}</p>
              <div className="foot">
                <span>{t('home.story.from')}</span>
                <Link to={site.links.why as never} className="link-arrow">
                  {t('home.story.more')}
                </Link>
              </div>
            </div>
          </div>
        </section>

        <section className="section" id="status">
          <div className="wrap">
            <SectionHead no="06" rail={`${version} · ${t('home.status.rail')}`}>
              <h2 className="h2">{t('home.status.title')}</h2>
            </SectionHead>
            <div className="board" data-reveal>
              {board.map(([tone, title, text]) => (
                <div key={title} className={'board-row' + (tone ? ' ' + tone : '')}>
                  <h3>
                    <i />
                    {title}
                  </h3>
                  <p>
                    <Inline text={text} />
                  </p>
                </div>
              ))}
            </div>
          </div>
        </section>

        <section className="closing">
          <div className="wrap closing-grid">
            <div>
              <h2 className="h2">{t('home.start.title')}</h2>
              <p className="req">{t('home.start.req')}</p>
            </div>
            <div>
              <Install margin={false} />
              <div className="actions">
                <Link to={site.links.docs} className="btn btn-primary">
                  {t('hero.start')}
                </Link>
              </div>
            </div>
          </div>
        </section>
      </main>
    </SiteShell>
  )
}
