import { Fragment, type ReactNode } from 'react'

// Inline turns [text](url) into links and `text` into code. The copy of the site is written
// this way, as in the README it comes from.
export function Inline({ text }: { text: string }): ReactNode {
  const parts = text.split(/(\[[^\]]+\]\([^)]+\)|`[^`]+`)/)
  return (
    <>
      {parts.map((part, i) => {
        const link = /^\[([^\]]+)\]\(([^)]+)\)$/.exec(part)
        if (link) {
          return (
            <a key={i} href={link[2]} className="underline decoration-fd-border underline-offset-4 hover:decoration-fd-foreground">
              {link[1]}
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
