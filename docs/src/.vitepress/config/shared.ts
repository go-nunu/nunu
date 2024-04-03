import { defineConfig } from 'vitepress'
import { search as zhSearch } from './zh'
import { search as enSearch } from './en'

export const shared = defineConfig({
  title: 'Nunu',

  lastUpdated: true,
  cleanUrls: true,
  metaChunk: true,

  markdown: {
    math: true,
    codeTransformers: [
      // We use `[!!code` in demo to prevent transformation, here we revert it back.
      {
        postprocess(code) {
          return code.replace(/\[\!\!code/g, '[!code')
        }
      }
    ]
  },

  sitemap: {
    hostname: 'https://go-nunu.dev',
    transformItems(items) {
      return items.filter((item) => !item.url.includes('migration'))
    }
  },

  /* prettier-ignore */
  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/Nunu-logo-mini.svg' }],
    ['link', { rel: 'icon', type: 'image/png', href: '/Nunu-logo-mini.png' }],
    ['meta', { name: 'theme-color', content: '#5f67ee' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:locale', content: 'en' }],
    ['meta', { property: 'og:title', content: 'Nunu | ' }],
    ['meta', { property: 'og:site_name', content: 'VitePress' }],
    ['meta', { property: 'og:image', content: 'https://go-nunu.dev/nunu-og.jpg' }],
    ['meta', { property: 'og:url', content: 'https://go-nunu.dev/' }],
    // ['script', { src: 'https://cdn.usefathom.com/script.js', 'data-site': 'AZBRSFGG', 'data-spa': 'auto', defer: '' }]
  ],

  themeConfig: {
    // logo: { src: '/nunu-logo-mini.svg', width: 24, height: 24 },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/go-nunu/nunu' }
    ],

    search: {
      provider: 'local',
      options: {
        // appId: '',
        // apiKey: '',
        // indexName: '',
        locales: { ...zhSearch, ...enSearch }
      }
    },

    // carbonAds: { code: '', placement: '' }
  }
})
