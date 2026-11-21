import { useAlternates, useLanguage } from 'tifa/client'

// The name of each language in its own language, shown in the switch.
const names: Record<string, string> = { 'en-US': 'English', 'zh-CN': '简体中文' }

// LanguageSwitch links to the page in each of its languages. The other language is another
// document, so these are plain links, as docs/guides/i18n.md has it. The page with one version
// and the 404 page have no alternates, so the switch shows nothing there.
export function LanguageSwitch({ label }: { label: string }) {
  const versions = useAlternates()
  const current = useLanguage()
  if (versions.length < 2) return null
  return (
    <nav aria-label={label} className="flex items-center gap-1 text-sm">
      {versions.map((v) => (
        <a
          key={v.language}
          href={v.url}
          hrefLang={v.language}
          lang={v.language}
          aria-current={v.language === current ? 'page' : undefined}
          className={
            v.language === current
              ? 'rounded-md px-2 py-1 font-medium text-fd-foreground'
              : 'rounded-md px-2 py-1 text-fd-muted-foreground hover:bg-fd-accent hover:text-fd-accent-foreground'
          }
        >
          {names[v.language] ?? v.language}
        </a>
      ))}
    </nav>
  )
}
