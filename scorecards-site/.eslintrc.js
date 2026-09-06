module.exports = {
  root: true,
  env: {
    browser: true,
    node: true,
    'vue/setup-compiler-macros': true,
  },
  parserOptions: {
    parser: '@babel/eslint-parser',
    requireConfigFile: false,
  },
  extends: ['@nuxtjs', 'plugin:nuxt/recommended', 'prettier'],
  plugins: [],
  // Nuxt 3 auto-imports (framework composables + our own composables/utils)
  globals: {
    defineNuxtConfig: 'readonly',
    defineNuxtPlugin: 'readonly',
    defineContentConfig: 'readonly',
    defineCollection: 'readonly',
    useState: 'readonly',
    useHead: 'readonly',
    useAsyncData: 'readonly',
    useRuntimeConfig: 'readonly',
    useRoute: 'readonly',
    useNuxtApp: 'readonly',
    createError: 'readonly',
    queryCollection: 'readonly',
    useHeaderColour: 'readonly',
    setHeaderColour: 'readonly',
    flattenToc: 'readonly',
  },
  // add your custom rules here
  rules: {
    camelcase: 0,
    'dot-notation': 0,
    snakecase: 0,
    'vue/multi-word-component-names': 'off',
  },
  overrides: [
    {
      files: ['pages/*', 'modules/*'],
      rules: {
        'vue/multi-word-component-names': 'off',
      },
    },
  ],
}
