import { useAlternates, useLanguage } from 'tifa/client'
import { Popover, PopoverContent, PopoverTrigger } from 'fumadocs-ui/components/ui/popover'
import { useTranslation } from 'tifa/i18n'
import { CheckIcon, ChevronDownIcon, GlobeIcon } from './icons'

// The name of each language in its own language.
const names: Record<string, string> = { 'en-US': 'English', 'zh-CN': '简体中文' }

// Where each language starts, for a page that has no versions of its own, such as the 404 page.
const homes = [
  { language: 'en-US', url: '/' },
  { language: 'zh-CN', url: '/zh-CN' },
]

// LanguageMenu is a button with the current language that opens a list of the languages. Each
// is a plain link to the page in that language: the other language is another document, as
// docs/guides/i18n.md has it.
export function LanguageMenu() {
  const { t } = useTranslation()
  const alternates = useAlternates()
  const current = useLanguage()
  const versions = alternates.length > 1 ? alternates : homes
  return (
    <Popover>
      <PopoverTrigger
        aria-label={t('language.label')}
        className="inline-flex h-9 items-center gap-1.5 rounded-md px-2 text-sm text-fd-muted-foreground transition-colors hover:bg-fd-accent hover:text-fd-accent-foreground data-[state=open]:bg-fd-accent"
      >
        <GlobeIcon />
        <span className="max-sm:hidden">{names[current] ?? current}</span>
        <ChevronDownIcon width={12} height={12} className="max-sm:hidden" />
      </PopoverTrigger>
      <PopoverContent align="end" className="flex w-40 flex-col gap-0.5 p-1">
        {versions.map((v) => (
          <a
            key={v.language}
            href={v.url}
            hrefLang={v.language}
            lang={v.language}
            aria-current={v.language === current ? 'page' : undefined}
            className="flex items-center justify-between rounded-md px-2 py-1.5 text-sm text-fd-muted-foreground transition-colors hover:bg-fd-accent hover:text-fd-accent-foreground aria-[current=page]:text-fd-foreground"
          >
            {names[v.language] ?? v.language}
            {v.language === current && <CheckIcon />}
          </a>
        ))}
      </PopoverContent>
    </Popover>
  )
}
