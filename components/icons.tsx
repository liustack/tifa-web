import type { SVGProps } from 'react'

// Small line icons for the header and the footer, drawn in the current text color.
function Icon(props: SVGProps<SVGSVGElement>) {
  return (
    <svg
      viewBox="0 0 24 24"
      width={16}
      height={16}
      fill="none"
      stroke="currentColor"
      strokeWidth={1.75}
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
      {...props}
    />
  )
}

export const GlobeIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon {...p}>
    <circle cx="12" cy="12" r="9" />
    <path d="M3 12h18M12 3c2.5 2.6 3.8 5.6 3.8 9s-1.3 6.4-3.8 9c-2.5-2.6-3.8-5.6-3.8-9S9.5 5.6 12 3Z" />
  </Icon>
)
export const ChevronDownIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon {...p}>
    <path d="m6 9 6 6 6-6" />
  </Icon>
)
export const CheckIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon {...p}>
    <path d="m5 12.5 4.5 4.5L19 7.5" />
  </Icon>
)
export const MenuIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon {...p}>
    <path d="M4 7h16M4 12h16M4 17h16" />
  </Icon>
)
export const CloseIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon {...p}>
    <path d="M6 6l12 12M18 6 6 18" />
  </Icon>
)
export const SunIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon {...p}>
    <circle cx="12" cy="12" r="4" />
    <path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" />
  </Icon>
)
export const MoonIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon {...p}>
    <path d="M20 14.5A8 8 0 0 1 9.5 4a8 8 0 1 0 10.5 10.5Z" />
  </Icon>
)
export const MonitorIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon {...p}>
    <rect x="3" y="4" width="18" height="12" rx="2" />
    <path d="M8 20h8M12 16v4" />
  </Icon>
)
export const GitHubIcon = (p: SVGProps<SVGSVGElement>) => (
  <svg viewBox="0 0 24 24" width={18} height={18} fill="currentColor" aria-hidden="true" {...p}>
    <path d="M12 .3a12 12 0 0 0-3.8 23.4c.6.1.8-.3.8-.6v-2c-3.3.7-4-1.6-4-1.6-.6-1.4-1.4-1.8-1.4-1.8-1-.7.1-.7.1-.7 1.2.1 1.8 1.2 1.8 1.2 1 1.8 2.8 1.3 3.5 1 .1-.8.4-1.3.7-1.6-2.7-.3-5.5-1.3-5.5-5.9 0-1.3.5-2.4 1.2-3.2-.1-.3-.5-1.5.1-3.2 0 0 1-.3 3.3 1.2a11.5 11.5 0 0 1 6 0c2.3-1.5 3.3-1.2 3.3-1.2.7 1.7.2 2.9.1 3.2.8.8 1.2 1.9 1.2 3.2 0 4.6-2.8 5.6-5.5 5.9.4.4.8 1.1.8 2.2v3.3c0 .3.2.7.8.6A12 12 0 0 0 12 .3Z" />
  </svg>
)

export const CopyIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon strokeWidth={2} {...p}>
    <rect x="9" y="9" width="12" height="12" rx="2" />
    <path d="M5 15V5a2 2 0 0 1 2-2h10" />
  </Icon>
)
export const FileIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon strokeWidth={2} {...p}>
    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
    <path d="M14 2v6h6" />
  </Icon>
)
export const ArrowRightIcon = (p: SVGProps<SVGSVGElement>) => (
  <svg viewBox="0 0 40 40" fill="none" stroke="currentColor" strokeWidth={1.6} strokeLinecap="round" strokeLinejoin="round" aria-hidden="true" {...p}>
    <path d="M6 20h26M24 12l8 8-8 8" />
  </svg>
)
export const BranchIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon strokeWidth={2} {...p}>
    <circle cx="6" cy="6" r="2.5" />
    <circle cx="6" cy="18" r="2.5" />
    <circle cx="18" cy="9" r="2.5" />
    <path d="M6 8.5v7M18 11.5c0 3-6 2-11 4.5" />
  </Icon>
)
export const DownloadIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon strokeWidth={2} {...p}>
    <path d="M12 3v12M7 10l5 5 5-5M5 21h14" />
  </Icon>
)
export const ServerIcon = (p: SVGProps<SVGSVGElement>) => (
  <Icon strokeWidth={2} {...p}>
    <rect x="3" y="4" width="18" height="7" rx="1.5" />
    <rect x="3" y="13" width="18" height="7" rx="1.5" />
    <path d="M7 7.5h.01M7 16.5h.01" />
  </Icon>
)
