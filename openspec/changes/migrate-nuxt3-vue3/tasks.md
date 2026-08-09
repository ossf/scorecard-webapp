## 1. Dependency reconciliation

- [ ] 1.1 Bump `vue` to `^3.x`; remove `vue-server-renderer` and `vue-template-compiler` from `package.json`
- [ ] 1.2 Remove dead Nuxt-2-only dependencies: `@nuxtjs/axios`, `@nuxtjs/proxy`, `@nuxtjs/redirect-module`, `vue-awesome-swiper`, `swiper`, `vue-intersect`, `vue-intersect-directive`
- [ ] 1.3 Replace `@nuxtjs/svg` with `vite-svg-loader` in `package.json`
- [ ] 1.4 Bump `@nuxt/content` to v3, `@nuxtjs/sitemap`, `@nuxtjs/google-fonts`, `@nuxtjs/tailwindcss` to their current Nuxt-3-compatible majors
- [ ] 1.5 Remove `@nuxtjs/eslint-module` from `package.json`
- [ ] 1.6 Add `mitt` as a dependency
- [ ] 1.7 Remove the `resolutions` block entries that only existed to patch transitive deps of now-removed packages (`@nuxtjs/svg` sub-resolutions); keep any still-relevant ones (e.g. `glob-parent` if still pulled in by a remaining dependency)
- [ ] 1.8 Run `yarn install` and confirm the lockfile resolves with no unresolvable peer conflicts; regenerate `yarn.lock`

## 2. Nuxt config rewrite

- [ ] 2.1 Rewrite `nuxt.config.js` using `defineNuxtConfig`
- [ ] 2.2 Move `head` config to `app.head`
- [ ] 2.3 Merge `buildModules` into the single `modules` array, dropping `@nuxtjs/eslint-module`
- [ ] 2.4 Update plugin entries to drop the Nuxt-2-only `mode`/`ssr` keys; rename `plugins/prism.js` to `plugins/prism.client.js`
- [ ] 2.5 Review `target: 'static'` and `generate.fallback` against Nuxt 3's Nitro-based static generation config; update or remove as needed
- [ ] 2.6 Update `sitemap`, `googleFonts` config blocks to match their new module major versions' config shape
- [ ] 2.7 Add `runtimeConfig.public` entries for `VUE_APP_FRONTEND` and `VUE_APP_SITENAME`
- [ ] 2.8 Create `postcss.config.js` if required by the updated `@nuxtjs/tailwindcss` major version

## 3. State management migration

- [ ] 3.1 Implement a `useState`-based composable replacing the Vuex `settings` module (`bg`, `textColor` state; `setHeaderColour` action-equivalent)
- [ ] 3.2 Update `modules/Header/Header.js` to use the new composable instead of `mapGetters('settings/...')`
- [ ] 3.3 Update `pages/index.vue` to use the new composable instead of `mapActions('settings', [...])`
- [ ] 3.4 Delete `store/index.js`, `store/state.js`, `store/getters.js`, `store/mutation-types.js`, `store/settings.js`, `store/README.md`

## 4. Event bus migration

- [ ] 4.1 Create a Nuxt plugin providing a `mitt()`-based event bus (e.g. `plugins/bus.js`, `provide: { bus: mitt() }`)
- [ ] 4.2 Replace all `this.$nuxt.$emit(...)`/`$on(...)`/`$off(...)` call sites with the new bus in: `components/OnScreen.vue`, `components/global/Sidebar.vue`, `layouts/default.vue`, `pages/index.vue`, `modules/Header/Header.js`
- [ ] 4.3 Remove the dead commented-out `window.$nuxt.$emit('setActiveToc', ...)` code in `plugins/components.client.js`

## 5. `@nuxt/content` v1 → v3 migration

- [ ] 5.1 Replace `asyncData({ $content }) => $content('home').fetch()` in `pages/index.vue` with the v3 equivalent (`useAsyncData` + content-v3 query API)
- [ ] 5.2 Replace `this.$content('home').only([...]).fetch()` in `components/global/Sidebar.vue` with the v3 equivalent
- [ ] 5.3 Replace `<nuxt-content :document="page" />` in `pages/index.vue` with `<ContentRenderer :value="page" />`
- [ ] 5.4 Port the custom `highlight.js` + `nord`-theme highlighter and `rehype-add-classes` table-class behavior to `@nuxt/content` v3's markdown build config; verify visually against current production rendering
- [ ] 5.5 Update `content.fullTextSearchFields` / search config to v3's equivalent, if the site's search UI depends on it
- [ ] 5.6 Verify `content/home.md`'s custom in-content components (`<sidebar>`, `<code-group>`, `<code-block>`) still resolve and render correctly under v3

## 6. Vue 3 API fixes

- [ ] 6.1 Rename `beforeDestroy` → `beforeUnmount` in `components/global/Sidebar.vue`, `layouts/default.vue`, `pages/index.vue`, `modules/Header/Header.js`
- [ ] 6.2 Replace `this.$scopedSlots.default({...})` with `this.$slots.default?.({...})` in `components/OnScreen.vue`
- [ ] 6.3 Rewrite `components/global/CodeGroup.vue`'s slot VNode introspection for Vue 3's VNode shape (replace `slot.componentOptions`/`slot.elm` usage with `vnode.type`/`vnode.props`/`vnode.el`)
- [ ] 6.4 Replace `Vue.extend(CodeCopyButton)` + `new CopyButton().$mount()` in `pages/index.vue` with `createApp(CodeCopyButton).mount(el)`
- [ ] 6.5 Update `plugins/components.client.js`: remove the dead Swiper registration (`Vue.use(getAwesomeSwiper(...))` and its imports); convert `Vue.directive('animate-on-scroll', { bind: ... })` to a Nuxt plugin registering `nuxtApp.vueApp.directive('animate-on-scroll', { mounted: ... })`
- [ ] 6.6 Convert the `capitalize` `filters:` entry in `layouts/default.vue` to a method or computed property, and update any template usage (`| capitalize`) accordingly

## 7. SVG and asset loading

- [ ] 7.1 Update SVG imports using the `@nuxtjs/svg` `?inline` suffix (`components/RepoButton.vue`, `components/global/CodeCopyButton.vue`, `modules/Footer/Footer.js`, `modules/Header/Header.js`) to `vite-svg-loader`'s import convention
- [ ] 7.2 Replace the `require.context('../assets/logos/', true, /\.svg$/)` call in `pages/index.vue` with `import.meta.glob('../assets/logos/**/*.svg')`

## 8. Build, scripts, and deploy config

- [ ] 8.1 Update the `start` script in `package.json` (Nuxt 3 has no `nuxt start`); replace with a static-preview-compatible command
- [ ] 8.2 Run `yarn generate` locally and determine Nuxt 3.21.7's actual static output directory
- [ ] 8.3 Update `netlify.toml`'s `publish` path if it no longer matches the actual output directory from 8.2
- [ ] 8.4 Add `.output` to `jsconfig.json`'s `exclude` list alongside `.nuxt`/`dist`

## 9. Verification

- [ ] 9.1 Run `yarn install --frozen-lockfile` and confirm it succeeds
- [ ] 9.2 Run `yarn lint` and confirm zero errors, including zero deprecated-Vue-2-API violations
- [ ] 9.3 Run `yarn build` and confirm it exits successfully
- [ ] 9.4 Run `yarn generate` and confirm it exits successfully and produces the expected static output
- [ ] 9.5 Start the dev server and manually verify in a browser: home page content renders (including all `content/home.md` sections), TOC scroll-spy, mobile nav toggle, code-copy buttons, code-group tabs, header background color change on scroll, canonical/OG/Twitter meta tags present in page source
- [ ] 9.6 Confirm the `dependency-review` advisory on `esbuild` is resolved (re-check `yarn why esbuild` or re-run the workflow); force-resolve via `resolutions` if still flagged
- [ ] 9.7 Push to the `dependabot/npm_and_yarn/scorecards-site/nuxt-3.21.7` branch and confirm all PR #980 checks pass: `build`, `lint`, `dependency-review`, Netlify deploy preview
- [ ] 9.8 Open the actual Netlify deploy preview URL and visually confirm the site renders correctly (not just that the build log succeeded)
