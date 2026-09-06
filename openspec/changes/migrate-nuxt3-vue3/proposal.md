## Why

Dependabot opened ossf/scorecard-webapp#980 bumping `nuxt` from `2.18.1` to `3.21.7` in `scorecards-site/`. Nuxt 3 is a ground-up rewrite that requires Vue 3, and the site still runs Vue 2 (`vue@^2.6.14`) with several Nuxt-2-only modules (`@nuxtjs/axios`, `@nuxtjs/proxy`, `@nuxtjs/redirect-module`, `@nuxtjs/svg`, `@nuxt/content@1`, `@nuxtjs/eslint-module@3`) and a Vue-2-only carousel dependency (`vue-awesome-swiper`). As a result PR #980 fails `build` (crashes inside `@nuxtjs/axios`'s Nuxt-2-only module hook), `lint` (flags real Options API usage — `beforeDestroy`, `$scopedSlots` — as deprecated under Vue 3 rules), `dependency-review` (a new transitive `esbuild` advisory pulled in via Vite), and the Netlify deploy preview (downstream of the broken build).

Merging the dependency bump as-is is not possible without migrating the whole app to Nuxt 3 / Vue 3. This change proposes doing that migration so the dependency bump (or an equivalent) can land, keeping the site on a maintained framework version.

## What Changes

- **BREAKING**: Upgrade `vue` 2.6 → 3.x, `vue-server-renderer`/`vue-template-compiler` are removed (Nuxt 3 uses Vite/Vue 3 SFC compiler internally).
- **BREAKING**: Replace `@nuxtjs/axios` with native `$fetch`/`useFetch` (Nuxt 3's built-in Ofetch-based data fetching); remove the module and its config from `nuxt.config.js`.
- **BREAKING**: Remove `@nuxtjs/proxy` (Nuxt-2-only dev-server proxy module, currently unconfigured/unused in `nuxt.config.js`); if a dev proxy is still needed, use Nitro's built-in dev proxy or Vite server proxy config instead.
- **BREAKING**: Replace `@nuxtjs/redirect-module` with Nitro route rules (`routeRules` in `nuxt.config.js`) or a Nuxt 3-compatible redirect module.
- **BREAKING**: Replace `@nuxtjs/svg` (webpack-loader-based) with a Vite-compatible SVG solution (e.g. `vite-svg-loader` or native `~icons`/asset imports), since Nuxt 3 builds with Vite by default.
- **BREAKING**: Upgrade `@nuxt/content` v1 → v3 (config shape, `<nuxt-content>` component, and query API all changed between major versions).
- **BREAKING**: Upgrade `@nuxtjs/eslint-module` → Nuxt 3's ESLint integration (`@nuxt/eslint-module` or standalone `eslint` + `@nuxt/eslint-config`), and fix the Options API patterns the newer `eslint-plugin-vue` now flags (`beforeDestroy` → `beforeUnmount`, `$scopedSlots` → `$slots`).
- **BREAKING**: Replace `vue-awesome-swiper` (Vue-2-only wrapper) with Swiper's native Vue 3 integration (`swiper/vue`) or an equivalent Vue-3-compatible carousel.
- Update `@nuxtjs/google-fonts`, `@nuxtjs/sitemap`, `@nuxtjs/tailwindcss` to their Nuxt-3-compatible major versions.
- Rewrite `nuxt.config.js` using `defineNuxtConfig`, moving `head` → `app.head`, `buildModules`/`modules` → unified `modules`, and reviewing `generate`/`target` options against Nuxt 3's `nitro`/`ssr` config.
- Audit and update all `.vue` files for Options API patterns incompatible with or deprecated under Vue 3 (lifecycle hook renames, `$scopedSlots` → `$slots`, `<no-ssr>` → `<client-only>`, event-bus patterns, filters removal, etc.), preserving existing Options API style (no forced Composition API rewrite) except where Vue 3 requires it.
- Verify the Netlify build output path and command still match Nuxt 3's output structure (`.output/public` for `nuxt generate`, vs. Nuxt 2's `dist/`).
- Resolve the `esbuild@0.27.7` dependency-review advisory (upgrade the transitive dependency once Vite/Nuxt versions are finalized, or confirm it's remediated upstream).

## Capabilities

### New Capabilities
- `scorecards-site-platform`: The build, rendering, and deployment behavior of the OpenSSF Scorecard marketing/docs site (`scorecards-site/`) — static generation, Markdown content rendering, client-side data fetching, redirects, sitemap generation, and asset handling — expressed as framework-agnostic requirements so the underlying stack (currently Nuxt 2/Vue 2, target Nuxt 3/Vue 3) can be verified against the same behavioral contract before and after migration.

### Modified Capabilities
- (none — no existing specs predate this change)

## Impact

- **Affected code**: All of `scorecards-site/` — `nuxt.config.js`, `pages/`, `layouts/`, `components/`, `plugins/`, `content/` (if present), `assets/`, `static/`, ESLint/Prettier config, `package.json`/`yarn.lock`.
- **Affected dependencies**: `nuxt`, `vue`, `vue-server-renderer`, `vue-template-compiler`, `@nuxtjs/axios`, `@nuxtjs/proxy`, `@nuxtjs/redirect-module`, `@nuxtjs/svg`, `@nuxt/content`, `@nuxtjs/eslint-module`, `@nuxtjs/eslint-config`, `vue-awesome-swiper`, `swiper`, `@nuxtjs/google-fonts`, `@nuxtjs/sitemap`, `@nuxtjs/tailwindcss`, `@tailwindcss/typography`, `sass-loader`.
- **Affected CI/deploy**: GitHub Actions `build`, `lint`, `dependency-review` workflows; Netlify build/deploy-preview config (build command and publish directory).
- **Out of scope**: No changes to `main.go` / Go backend services, or to visual design/content of the site — this is a like-for-like framework migration, not a redesign.
