import svgLoader from 'vite-svg-loader'

export default defineNuxtConfig({
  ssr: true,

  // Nuxt 3 renamed the public assets convention directory from `static/` to
  // `public/`; keep serving the existing `static/` directory unchanged.
  dir: {
    public: 'static',
  },

  // Global page headers: https://nuxt.com/docs/api/nuxt-config#head
  app: {
    head: {
      title: 'OpenSSF Scorecard',
      htmlAttrs: {
        lang: 'en',
      },
      meta: [
        { charset: 'utf-8' },
        {
          name: 'viewport',
          content: 'width=device-width, initial-scale=1, user-scalable=no',
        },
        { name: 'format-detection', content: 'telephone=no' },
        { name: 'msapplication-TileColor', content: '#da532c' },
        {
          hid: 'description',
          name: 'description',
          content: 'Quickly assess open source projects for risky practices',
        },
        {
          hid: 'keywords',
          name: 'keywords',
          content:
            'scorecards, scorecard, openssf, slsa, sigstore, security, vulnerabilities, cve, supply chain, supply-chain',
        },
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.png' },
        { rel: 'mask-icon', href: '/safari-pinned-tab.svg', color: '#5bbad5' },
        {
          rel: 'icon',
          type: 'image/png',
          sizes: '16x16',
          href: '/favicon-16x16.png',
        },
        {
          rel: 'icon',
          type: 'image/png',
          sizes: '32x32',
          href: '/favicon-32x32.png',
        },
        {
          rel: 'apple-touch-icon',
          sizes: '180x180',
          href: '/apple-touch-icon.png',
        },
      ],
    },
  },

  // Global CSS: https://nuxt.com/docs/api/nuxt-config#css
  css: ['@/assets/css/base.scss'],

  // Auto import components: https://nuxt.com/docs/api/nuxt-config#components
  components: true,

  // Modules: https://nuxt.com/docs/api/nuxt-config#modules
  modules: [
    '@nuxtjs/tailwindcss',
    '@nuxtjs/google-fonts',
    '@nuxt/content',
    '@nuxtjs/sitemap',
  ],

  site: {
    url: process.env.VUE_APP_FRONTEND || 'http://localhost:3000',
  },

  sitemap: {
    cacheMaxAgeSeconds: 86400,
  },

  content: {
    experimental: {
      // Use Node's built-in `node:sqlite` instead of the `better-sqlite3`
      // native addon, which ships prebuilt binaries per platform/Node ABI
      // and can fail to load on build hosts it wasn't prebuilt for.
      sqliteConnector: 'native',
    },
    build: {
      markdown: {
        toc: {
          depth: 3,
          searchDepth: 3,
        },
        highlight: {
          theme: 'nord',
        },
        rehypePlugins: {
          'rehype-add-classes': { options: { table: 'table' } },
        },
      },
    },
  },

  googleFonts: {
    families: {
      'Public Sans': [400, 600, 700],
      'DM Mono': [400, 500],
    },
    display: 'swap', // 'auto' | 'block' | 'swap' | 'fallback' | 'optional'
  },

  runtimeConfig: {
    public: {
      frontendUrl: process.env.VUE_APP_FRONTEND || 'http://localhost:3000',
      siteName: process.env.VUE_APP_SITENAME || 'OpenSSF Scorecard',
    },
  },

  nitro: {
    prerender: {
      failOnError: true,
    },
  },

  vite: {
    plugins: [svgLoader()],
  },

  // Build Configuration: https://nuxt.com/docs/api/nuxt-config#build
  build: {},
})
