import addClasses from 'rehype-add-classes';
import remarkGemoji from 'remark-gemoji'
import { getHighlighter } from 'shiki'
import highlightjs from 'highlight.js'
export default {
  target: 'static',
  ssr: false,
  // Global page headers: https://go.nuxtjs.dev/config-head
  head: {
    title: 'goost-scorecards',
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
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
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

  // Global CSS: https://go.nuxtjs.dev/config-css
  css: ['@/assets/css/base','highlight.js/styles/nord.css'],

  // Plugins to run before rendering page: https://go.nuxtjs.dev/config-plugins
  plugins: [
    { src: '~plugins/components.client' },
    { src: '~plugins/prism', mode: 'client', ssr: false }
  ],

  // Auto import components: https://go.nuxtjs.dev/config-components
  components: true,

  // Modules for dev and build (recommended): https://go.nuxtjs.dev/config-modules
  buildModules: [
    // https://go.nuxtjs.dev/eslint
    '@nuxtjs/eslint-module',
    // https://go.nuxtjs.dev/tailwindcss
    '@nuxtjs/tailwindcss',

    '@nuxtjs/google-fonts',
  ],

  // Modules: https://go.nuxtjs.dev/config-modules
  modules: [
    // https://go.nuxtjs.dev/axios
    '@nuxtjs/axios',
    // https://go.nuxtjs.dev/content
    '@nuxt/content',

    '@nuxtjs/svg',

    '@nuxtjs/dotenv',

    '@nuxtjs/redirect-module',

    '@nuxtjs/sitemap',

    '@nuxtjs/proxy',
  ],

  proxy: [
    // // Proxies /foo to http://example.com/foo
    // 'http://example.com/foo',
    // // Proxies /api/books/*/**.json to http://example.com:8000
    // 'http://example.com:8000/api/books/*/**.json',
    // // You can also pass more options
    // [ 'http://example.com/foo', { ws: false } ]
  ],

  sitemap: {
    path: '/sitemap.xml',
    hostname: process.env.VUE_APP_FRONTEND,
    generate: true,
    cacheTime: 86400,
    trailingSlash: true,
  },

  content: {
    markdown: {
      highlighter(rawCode, lang) {
        const highlightedCode = highlightjs.highlight(rawCode, { language: lang }).value

        // We need to create a wrapper, because
        // the returned code from highlight.js
        // is only the highlighted code.
        return `<pre><code class="hljs ${lang}">${highlightedCode}</code></pre>`
      },
      rehypePlugins: [
        ['rehype-add-classes', { table: 'table' }]
      ],
      // prism: {
      //   theme: 'prism-themes/themes/prism-material-oceanic.css'
      // },
      // async highlighter() {
      //   const highlighter = await getHighlighter({

      //     theme: 'nord'
      //   })
      //   return (rawCode, lang) => {
      //     return highlighter.codeToHtml(rawCode, lang)
      //   }
      // }
    },
    fullTextSearchFields: ['title', 'description', 'slug', 'text'],
  },

  googleFonts: {
    families: {
      'Public Sans': [400,600],
      'IBM Plex Mono': [400,500],
    },
    display: 'swap' // 'auto' | 'block' | 'swap' | 'fallback' | 'optional'
  },

  // Axios module configuration: https://go.nuxtjs.dev/config-axios
  axios: {
    // Workaround to avoid enforcing hard-coded localhost:3000: https://github.com/nuxt-community/axios-module/issues/308
    baseURL: '/',
  },

  generate: {
    fallback: true,
  },

  // Build Configuration: https://go.nuxtjs.dev/config-build
  build: {},
}