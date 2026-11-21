import { useTranslation } from 'tifa/i18n'
import { MonitorIcon, MoonIcon, SunIcon } from './icons'
import { useTheme } from './theme'

const options = [
  { value: 'system', Icon: MonitorIcon },
  { value: 'light', Icon: SunIcon },
  { value: 'dark', Icon: MoonIcon },
] as const

// ThemeSwitch is the three small buttons of the footer: the system's theme, light and dark. The
// choice is kept by ThemeProvider in localStorage. The server does not know it, so the selected
// button shows once the page is in the browser.
export function ThemeSwitch() {
  const { t } = useTranslation()
  const { theme, setTheme } = useTheme()
  return (
    <div role="radiogroup" aria-label={t('footer.theme')} className="inline-flex rounded-full border border-fd-border p-0.5">
      {options.map(({ value, Icon }) => {
        const selected = theme === value
        return (
          <button
            key={value}
            type="button"
            role="radio"
            aria-checked={selected}
            aria-label={t(`footer.${value}`)}
            onClick={() => setTheme(value)}
            className={
              'inline-flex size-6 items-center justify-center rounded-full transition-colors ' +
              (selected ? 'bg-fd-accent text-fd-foreground' : 'text-fd-muted-foreground hover:text-fd-foreground')
            }
          >
            <Icon width={14} height={14} />
          </button>
        )
      })}
    </div>
  )
}
