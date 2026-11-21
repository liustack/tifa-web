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
              // Without langs Shiki loads all of its languages before the first block, about 2 s
              // of every build. TestCodeLanguagesAreLoaded fails when a doc uses one missing here.
              langs: ['go', 'sh', 'tsx', 'ts', 'yaml', 'json', 'ini', 'diff', 'markdown', 'md', 'sql', 'css', 'xml'],
            },
          ],
        ],
      },
    }),
    react(),
    tailwindcss(),
  ],
})
