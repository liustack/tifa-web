import { useState } from 'react'
import { useTranslation } from 'tifa/i18n'
import { CheckIcon, CopyIcon } from './icons'

// CopyButton puts a command on the clipboard and shows a check for a moment.
export function CopyButton({ text }: { text: string }) {
  const { t } = useTranslation()
  const [copied, setCopied] = useState(false)
  const copy = () => {
    const done = () => {
      setCopied(true)
      setTimeout(() => setCopied(false), 1600)
    }
    navigator.clipboard.writeText(text).then(done, done)
  }
  return (
    <button type="button" className="copy" data-copied={copied ? '' : undefined} aria-label={t('hero.copy')} onClick={copy}>
      {copied ? <CheckIcon /> : <CopyIcon />}
    </button>
  )
}
