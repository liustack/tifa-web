import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'

export type Theme = 'system' | 'light' | 'dark'

const key = 'theme'

// themeScript runs in the head, before the first paint, so that a page never shows the wrong theme
// first: dark unless the visitor chose otherwise. It is written into the document and not
// rendered by React, so a page that React renders in the browser (an error page) does not render
// a script.
export const themeScript = `try{var t=localStorage.getItem('${key}')||'dark';var d=t==='dark'||(t==='system'&&matchMedia('(prefers-color-scheme: dark)').matches);document.documentElement.classList.toggle('dark',d);document.documentElement.style.colorScheme=d?'dark':'light'}catch(e){document.documentElement.classList.add('dark')}`

function apply(theme: Theme) {
  const dark = theme === 'dark' || (theme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)
  document.documentElement.classList.toggle('dark', dark)
  document.documentElement.style.colorScheme = dark ? 'dark' : 'light'
}

function read(): Theme {
  try {
    const stored = localStorage.getItem(key)
    return stored === 'system' || stored === 'light' || stored === 'dark' ? stored : 'dark'
  } catch {
    return 'dark'
  }
}

const Context = createContext<{ theme: Theme | undefined; setTheme: (theme: Theme) => void }>({ theme: undefined, setTheme: () => {} })

// ThemeProvider keeps the choice of the visitor: system, light or dark, in localStorage. The
// server does not know it, so theme is undefined until the browser has read it, and what shows the
// choice waits for that.
export function ThemeProvider({ children }: { children: ReactNode }) {
  const [theme, setState] = useState<Theme | undefined>(undefined)
  useEffect(() => setState(read()), [])
  useEffect(() => {
    if (theme !== 'system') return
    const media = window.matchMedia('(prefers-color-scheme: dark)')
    const change = () => apply('system')
    media.addEventListener('change', change)
    return () => media.removeEventListener('change', change)
  }, [theme])
  useEffect(() => {
    const sync = (e: StorageEvent) => {
      if (e.key !== key) return
      const next = read()
      setState(next)
      apply(next)
    }
    window.addEventListener('storage', sync)
    return () => window.removeEventListener('storage', sync)
  }, [])
  const setTheme = useCallback((next: Theme) => {
    try {
      localStorage.setItem(key, next)
    } catch {
      // The choice then lasts for this page only.
    }
    setState(next)
    apply(next)
  }, [])
  const value = useMemo(() => ({ theme, setTheme }), [theme, setTheme])
  return <Context value={value}>{children}</Context>
}

export function useTheme() {
  return useContext(Context)
}
