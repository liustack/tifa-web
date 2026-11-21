import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import rehypeShiki from '@shikijs/rehype'
import tifa from './.tifa/vite.js'

export default defineConfig({
  plugins: [
    tifa({
      markdown: {
        // Code blocks are colored at build time, in a light and a dark theme. The colors are
        // CSS variables (--shiki-light, --shiki-dark) that Fumadocs' stylesheet picks by theme.
        rehypePlugins: [
          [
            rehypeShiki,
            {
              themes: { light: 'github-light', dark: 'github-dark' },
              defaultColor: false,
              fallbackLanguage: 'text',
            },
          ],
        ],
      },
    }),
    react(),
    tailwindcss(),
  ],
})
