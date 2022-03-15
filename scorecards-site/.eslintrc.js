module.exports = {
  root: true,
  env: {
    browser: true,
    node: true,
  },
  parserOptions: {
    parser: '@babel/eslint-parser',
    requireConfigFile: false,
  },
  extends: ['@nuxtjs', 'plugin:nuxt/recommended', 'prettier'],
  plugins: [],
  // add your custom rules here
  rules: {
    camelcase: 0,
    'dot-notation': 0,
    snakecase: 0,
    'vue/multi-word-component-names': 'off',
  },
  overrides: [
      {
        files: ['pages/*','modules/*'],
        rules: {
          'vue/multi-word-component-names': 'off'
        }
      }
    ]
}
