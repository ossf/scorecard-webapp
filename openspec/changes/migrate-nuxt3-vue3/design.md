## Context

`scorecards-site/` (the OpenSSF Scorecard marketing/docs site) is a single-page Nuxt app: one route (`pages/index.vue`), one layout, ~1,800 lines of `.vue`/`.js` across `components/`, `modules/` (Header/Footer), `store/`, and `plugins/`, plus one `@nuxt/content` v1 markdown document (`content/home.md`, 350 lines) that renders most of the page body.

Dependabot's automated bump (`nuxt` `2.18.1` → `3.21.7`, PR #980) only ever touches the version string it's targeting — it left `vue@^2.6.14`, `vue-server-renderer`, `vue-template-compiler`, `@nuxt/content@^1.15.1`, and every Nuxt-2-only module in place. `package.json` on the PR branch is consequently in a **broken transitional state**, not a working Nuxt 2 app plus one pending bump. A code-level audit (see proposal) additionally found:
- `vuex` is imported (`Header.js`, `pages/index.vue`) but **not declared in `package.json` at all** — it only worked because Nuxt 2 bundled Vuex transitively. This must be fixed regardless of migration.
- Several Nuxt-2-only dependencies are already dead code with zero runtime usage: `@nuxtjs/axios`, `@nuxtjs/proxy`, `@nuxtjs/redirect-module`, `vue-awesome-swiper`/`swiper`, `vue-intersect`, `vue-intersect-directive`.
- The app's cross-component communication is a homegrown event bus riding on Nuxt 2's `$nuxt` root instance (`$nuxt.$on/$off/$emit`), used in ~13 call sites across 5 files for nav toggle, TOC scroll-spy, and intersection-observer relaying. Nuxt 3 has no `$nuxt` event bus.
- The Vuex store is mostly dead (`store/index.js` is commented out) except one small `settings` module (~30 lines: header background/text colour).

## Goals / Non-Goals

**Goals:**
- Get `scorecards-site/` building, linting, and deploying cleanly on Nuxt 3.21.7 / Vue 3, with no regression in current site behavior or visual appearance.
- Reconcile `package.json`/`yarn.lock` into an internally-consistent Nuxt-3-compatible dependency set in one atomic change (not incremental patches on top of the broken intermediate state).
- Remove dependencies that are already dead code as part of the cleanup (lower migration surface, no separate follow-up needed).
- Resolve the `dependency-review` `esbuild` advisory as a byproduct of landing on current, patched Nuxt 3/Vite versions.

**Non-Goals:**
- No visual redesign or content changes to the site.
- No Composition API rewrite — components keep the Options API style already in use; only patterns Vue 3 actually removed/renamed are touched.
- No migration off `@nuxt/content` to a different content system — staying within the `@nuxt/content` family (v1 → v3).
- No adoption of new state-management libraries beyond what's needed to replace the one live Vuex module.

## Decisions

1. **Dependency reconciliation is atomic.** All package.json changes (removals, replacements, version bumps) land in one commit/tasks-group before any source code is touched, so there's never a half-migrated `yarn install` state checked in.

2. **Drop dead Nuxt-2-only dependencies outright** rather than migrating them: `@nuxtjs/axios` (no `$axios` call sites), `@nuxtjs/proxy` (registered with an empty `proxy: []` config), `@nuxtjs/redirect-module` (no redirect rules configured anywhere, and `netlify.toml` has no `[[redirects]]` either), `vue-awesome-swiper`/`swiper` (registered globally in `plugins/components.client.js` but zero `<swiper>`/`<Swiper>` usage in any template), `vue-intersect`/`vue-intersect-directive` (never imported; all intersection-observer logic is hand-rolled with the native API already).
   - *Alternative considered*: migrate each to its Nuxt-3-compatible equivalent for parity. Rejected — porting dead code adds migration risk and review surface for zero behavioral benefit.

3. **Replace `@nuxtjs/svg` with `vite-svg-loader`.** Nuxt 3 builds with Vite by default; `vite-svg-loader` supports the same "import SVG as a Vue component" pattern the app already uses (`import Icon from '@/assets/icons/x.svg?inline'` → drop the `?inline` suffix per `vite-svg-loader`'s convention, or configure it to accept the existing suffix).
   - *Alternative considered*: `nuxt-svgo`. Either works; `vite-svg-loader` is chosen for being a thinner, more widely-used Vite-native option with less config surface.

4. **Replace the webpack-only `require.context()` call in `pages/index.vue`** (dynamic import of `assets/logos/*.svg`) with Vite's `import.meta.glob()`. There is no Vite equivalent of `require.context`; this is a mandatory (not optional) change for the build to succeed at all.

5. **Migrate `@nuxt/content` v1 → v3**, not v2. v3 is current and actively maintained; skipping v2 avoids a second migration later. Concretely:
   - `asyncData({ $content }) => $content('home').fetch()` → `useAsyncData` + `queryCollection`/`queryContent` (finalize exact API during implementation against installed v3 version).
   - `<nuxt-content :document="page" />` → `<ContentRenderer :value="page" />`.
   - Custom `highlighter` (highlight.js, `nord` theme) and `rehypePlugins: [['rehype-add-classes', {table: 'table'}]]` → re-implemented via v3's markdown build config. v3 defaults to Shiki; if Shiki cannot reproduce the exact `nord` highlight.js theme and the `rehype-add-classes` table-class behavior through supported config, custom markdown transform hooks are used to preserve current visual output exactly (this is flagged as a risk below, not assumed solved up front).

6. **Replace the Vuex `settings` module with Nuxt 3's built-in `useState()` composable**, not Pinia. The only live state is `{ bg, textColor }` plus one action (`setHeaderColour`) — `useState('headerColour', () => ({ bg: null, textColor: null }))` plus a plain function to mutate it covers this with zero new dependencies. The dead `store/index.js`, `store/state.js`, `store/getters.js`, `store/mutation-types.js` are deleted outright (already non-functional/commented out).
   - *Alternative considered*: Pinia (Nuxt's official recommended Vuex replacement). Rejected as disproportionate — Pinia is designed for real multi-module stores; this app has one tiny piece of shared state, and adding a state-management library for it would be over-engineering.

7. **Replace the `$nuxt` event bus with a `mitt`-based bus provided via a Nuxt plugin** (`defineNuxtPlugin` returning `provide: { bus: mitt() }`, or attached to `nuxtApp`). This is the closest 1:1 behavioral replacement for the existing `$emit`/`$on`/`$off` call sites (`openNavigation`, `storeTocs`, `observer.observed`, `observer.created`, `setActiveToc`) with the least risk of subtly changing event timing/ordering.
   - *Alternative considered*: refactor to `provide`/`inject` or lift state to a shared composable. Rejected for this pass — larger blast radius for a behavior-preservation migration; can be a follow-up cleanup once the app is stable on Nuxt 3.

8. **Vue 3 API renames/removals are fixed mechanically, one-to-one, preserving Options API style**: `beforeDestroy` → `beforeUnmount` (`Sidebar.vue`, `layouts/default.vue`, `pages/index.vue`, `Header.js`); `this.$scopedSlots.default({...})` → `this.$slots.default?.({...})` (`OnScreen.vue`); `Vue.extend(...)` + `new X().$mount()` → `createApp(X).mount(el)` (`pages/index.vue`'s imperative copy-button injection into content-rendered DOM — behavior preserved as-is, not redesigned, per non-goals); `Vue.directive`/`Vue.use` at plugin scope → `nuxtApp.vueApp.directive`/`.use` inside `defineNuxtPlugin`; directive's `bind` hook → `mounted`; the Vue-2-only VNode introspection in `CodeGroup.vue` (`slot.componentOptions`, `slot.elm`) → rewritten against Vue 3's VNode shape (`vnode.type`, `vnode.props`, `vnode.el`); the `capitalize` `filters:` entry in `layouts/default.vue` → converted to a method (Vue 3 removed `filters` entirely), after confirming actual template usage.

9. **`process.env.VUE_APP_*` reads move into `runtimeConfig.public`**, accessed via `useRuntimeConfig()`. Vite does not statically replace arbitrary `process.env.*` the way webpack did; the existing env var names (`VUE_APP_FRONTEND`, `VUE_APP_SITENAME`) are kept as the config keys/env var names to avoid touching Netlify env var configuration outside this repo.

10. **`nuxt.config.js` is rewritten with `defineNuxtConfig`**: `head` → `app.head`; `buildModules` merged into the single `modules` array; `@nuxtjs/eslint-module` is dropped from the module list entirely (the `lint` script already runs `eslint` standalone via CLI — the buildModule only added dev-server integration, which isn't required for `yarn lint`/CI to pass); plugin entries drop the Nuxt-2-only `mode`/`ssr` keys in favor of `.client.js` filename suffixes (already the case for `components.client.js`; `prism.js` gains a `.client.js` rename); `generate.fallback` and `target: 'static'` are reviewed against Nuxt 3's Nitro-based static generation config.

11. **`@nuxtjs/sitemap`, `@nuxtjs/google-fonts`, `@nuxtjs/tailwindcss` are bumped to their current Nuxt-3-compatible majors**, keeping the same modules rather than switching to manual Vite config, to minimize config-shape churn. Each module's config block is updated to match its new major version's shape (exact diffs resolved during implementation against each module's current docs).

12. **Netlify build output path is validated empirically, not assumed.** `netlify.toml` currently publishes `dist` (Nuxt 2's `nuxt generate` output). Nuxt 3.21.7's actual `nuxt generate` output location is confirmed by running the command locally before deciding whether `netlify.toml`'s `publish` value needs to change.

13. **Work continues on the existing `dependabot/npm_and_yarn/scorecards-site/nuxt-3.21.7` branch** (already checked out, already the branch backing PR #980), rather than opening a new branch — this keeps the fix and the triggering dependency bump in the same PR, which is what CI is currently failing on.

## Risks / Trade-offs

- **`@nuxt/content` v1→v3 is the largest single risk.** The query API, `<nuxt-content>` component, and markdown-highlighting pipeline all changed across two majors. → Mitigation: treat the custom highlighter/rehype-classes behavior as a spike within the implementation tasks; verify rendered output (syntax highlighting theme, table classes) visually against the current production site before considering this task done.
- **Netlify output-directory mismatch could silently break the deploy preview** (build succeeds locally but Netlify serves an empty/stale directory). → Mitigation: don't rely on local `yarn generate` success alone; push and inspect the actual Netlify deploy preview output before merging.
- **`mitt`-based event bus changes timing/lifecycle subtly** for user-visible interactive behavior (TOC scroll-spy highlighting, mobile nav toggle). → Mitigation: manual browser QA of these specific interactions (per CLAUDE.md's requirement to test UI changes in a real browser, not just via build/lint success).
- **`CodeGroup.vue`'s Vue-2 VNode introspection rewrite** changes how tabbed code blocks are constructed from slot content. → Mitigation: dedicated manual check of `content/home.md`'s code-group-tabbed sections after migration.
- **`esbuild` advisory (low severity, Windows dev-server-only, GHSA-g7r4-m6w7-qqqr)** doesn't affect the production static build, but `dependency-review` will still fail if the final lockfile pins a vulnerable transitive version. → Mitigation: re-run `dependency-review` (or `yarn why esbuild`) after all version bumps land; if still flagged, force-resolve `esbuild` to a patched version via `resolutions`.
- **`prism.js` plugin may be dead code** (no confirmed invocation site found) predating the switch to `highlight.js` for content syntax highlighting. → Not removing it as part of this migration (out of scope per non-goals: no unrelated cleanup beyond what's needed for compatibility) — it is ported (rename to `.client.js`, fix any Vue-2-only API calls if present) rather than deleted, and flagged as a candidate for a separate follow-up cleanup PR.

## Migration Plan

1. Reconcile `package.json`/`yarn.lock`: bump `vue` to `^3.x`; remove `vue-server-renderer`, `vue-template-compiler`; add `vuex` is *not* re-added (superseded by `useState`, decision 6); remove dead deps (decision 2); replace `@nuxtjs/svg` with `vite-svg-loader`; bump `@nuxt/content`, `@nuxtjs/sitemap`, `@nuxtjs/google-fonts`, `@nuxtjs/tailwindcss` to Nuxt-3-compatible majors; drop `@nuxtjs/eslint-module`; add `mitt`.
2. Rewrite `nuxt.config.js` (decision 10) and `postcss.config.js` if the new `@nuxtjs/tailwindcss` major requires one.
3. Replace the Vuex store with `useState` (decision 6); delete dead `store/*` files.
4. Add the `mitt` event-bus plugin (decision 7); update all ~13 `$nuxt.$on/$off/$emit` call sites to use it.
5. Migrate `@nuxt/content` usage end-to-end (decision 5): `asyncData`/`$content` → `useAsyncData`/content-v3 query API, `<nuxt-content>` → `<ContentRenderer>`, highlighter/rehype config.
6. Fix remaining Vue 3 API changes (decision 8): lifecycle hook renames, `$scopedSlots`, `Vue.extend`/imperative mount, custom directive hook rename, `CodeGroup.vue` VNode shape, `filters` removal.
7. Replace SVG imports and `require.context` (decisions 3–4).
8. Move `process.env.VUE_APP_*` reads to `runtimeConfig.public` (decision 9).
9. Update `package.json` scripts (`start` script's `nuxt start` doesn't exist in Nuxt 3) and validate/update `netlify.toml`'s `publish` path (decision 12).
10. Verify locally: `yarn install`, `yarn lint`, `yarn build`, `yarn generate`; then run the dev server and manually check in a browser: nav toggle, TOC scroll-spy, code-copy buttons, code-group tabs, header colour-on-scroll, content rendering/highlighting fidelity against the current live site.
11. Push to the `dependabot/npm_and_yarn/scorecards-site/nuxt-3.21.7` branch; confirm all PR #980 checks (`build`, `lint`, `dependency-review`, Netlify deploy preview) go green; inspect the actual deploy preview URL, not just the build log.

Rollback: since all work lands on the existing PR branch (not `main`), rollback is simply not merging the PR / reverting the branch to its current dependabot-only state — no production impact until merged.

## Open Questions

- Exact `@nuxt/content` v3 config shape for reproducing the current `highlight.js` + `nord` theme + `rehype-add-classes` table styling — resolve via spike during implementation (task-level, not blocking proposal approval).
- Whether Nuxt 3.21.7's `nuxt generate` truly outputs to `dist/` (some Nuxt 3 versions/presets differ) — resolve empirically (decision 12) before touching `netlify.toml`.
- Whether `prism.js` is genuinely dead code — worth a final grep pass during implementation; if confirmed dead, note it for a follow-up removal PR rather than deleting now.
