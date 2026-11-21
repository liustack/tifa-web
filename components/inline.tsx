import { Fragment, type ReactNode } from 'react'
import { Link } from '@tanstack/react-router'

const linkClass = 'underline decoration-fd-border underline-offset-4 hover:decoration-fd-foreground'

// Inline turns **text** into bold, [text](url) into links and `text` into code. The copy of the
// site is written this way, as in the README it comes from. A link that starts with / is a page
// of the site, which keeps the language of the page. Any other link leaves the site.
export function Inline({ text }: { text: string }): ReactNode {
  const parts = text.split(/(\*\*.+?\*\*(?!\*)|\[[^\]]+\]\([^)]+\)|`[^`]+`)/)
  return (
    <>
      {parts.map((part, i) => {
        const bold = /^\*\*(.+)\*\*$/.exec(part)
        if (bold) {
          return (
            <strong key={i} className="font-semibold text-fd-foreground">
              <Inline text={bold[1]} />
            </strong>
          )
        }
        const link = /^\[([^\]]+)\]\(([^)]+)\)$/.exec(part)
        if (link) {
          const [, label, url] = link
          if (url.startsWith('/')) {
            const [to, hash] = url.split('#')
            return (
              <Link key={i} to={to as never} hash={hash} className={linkClass}>
                {label}
              </Link>
            )
          }
          return (
            <a key={i} href={url} className={linkClass}>
              {label}
            </a>
          )
        }
        if (part.startsWith('`') && part.endsWith('`') && part.length > 2) {
          return (
            <code key={i} className="rounded bg-fd-muted px-1 py-0.5 text-[0.9em]">
              {part.slice(1, -1)}
            </code>
          )
        }
        return <Fragment key={i}>{part}</Fragment>
      })}
    </>
  )
}
