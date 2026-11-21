import { useEffect, useState } from 'react'
import type { TerminalLine } from './site'

// A line of a terminal as the markup of landing.css has it.
function Line({ line, shown }: { line: TerminalLine; shown?: number }) {
  const text = shown === undefined ? line.text : line.text.slice(0, shown)
  switch (line.kind) {
    case 'cmd':
      return (
        <>
          <span className="t-p">$ </span>
          <span className="t-cmd">{text}</span>
        </>
      )
    case 'ok':
      return <span className="t-ok">{text}</span>
    case 'dim':
      return <span className="t-dim">{text}</span>
    default:
      return <>{text}</>
  }
}

// Terminal shows lines of a terminal. The server writes all of them, so the page is complete
// without script. With type, the browser then takes them away and types the commands out again,
// unless the visitor asked for less motion.
export function Terminal({ title, label, lines, type = false, cursor = false }: { title: string; label: string; lines: TerminalLine[]; type?: boolean; cursor?: boolean }) {
  // null shows every line. Otherwise the lines before n are whole and line n has c characters.
  const [view, setView] = useState<{ n: number; c: number } | null>(null)
  useEffect(() => {
    if (!type || window.matchMedia('(prefers-reduced-motion: reduce)').matches) return
    let stopped = false
    let timer: ReturnType<typeof setTimeout>
    let i = 0
    const step = () => {
      if (stopped) return
      if (i >= lines.length) return
      const line = lines[i]
      if (line.kind === 'cmd') {
        let k = 0
        setView({ n: i, c: 0 })
        const typeChar = () => {
          if (stopped) return
          k++
          setView({ n: i, c: k })
          if (k < line.text.length) {
            timer = setTimeout(typeChar, 18 + Math.random() * 30)
          } else {
            i++
            timer = setTimeout(step, 420)
          }
        }
        typeChar()
      } else {
        i++
        setView({ n: i, c: 0 })
        timer = setTimeout(step, 90)
      }
    }
    setView({ n: 0, c: 0 })
    timer = setTimeout(step, 500)
    return () => {
      stopped = true
      clearTimeout(timer)
    }
  }, [type, lines])
  const whole = view ? lines.slice(0, view.n) : lines
  const partial = view && view.n < lines.length && lines[view.n].kind === 'cmd' ? lines[view.n] : null
  return (
    <div className="win" role="img" aria-label={label}>
      <div className="win-head">
        <span className="dots">
          <i />
          <i />
          <i />
        </span>
        <span className="title">{title}</span>
      </div>
      <pre>
        {whole.map((line, i) => (
          <span key={i}>
            <Line line={line} />
            {'\n'}
          </span>
        ))}
        {partial && view && <Line line={partial} shown={view.c} />}
        {cursor && <span className="cursor" />}
      </pre>
    </div>
  )
}
