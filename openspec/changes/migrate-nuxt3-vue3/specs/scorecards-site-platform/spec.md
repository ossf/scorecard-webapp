## ADDED Requirements

### Requirement: Production build succeeds
The `scorecards-site` build pipeline SHALL produce a static production build via `yarn build` and `yarn generate` with zero errors, using a mutually-compatible set of framework dependencies (no mixed Nuxt 2/Vue 2 and Nuxt 3/Vue 3 packages).

#### Scenario: Clean install and build
- **WHEN** `yarn install --frozen-lockfile` followed by `yarn build` is run against the committed `package.json`/`yarn.lock`
- **THEN** both commands exit successfully with no module-resolution or module-initialization errors

#### Scenario: Static generation succeeds
- **WHEN** `yarn generate` is run after a successful build
- **THEN** it exits successfully and produces a static output directory containing the rendered home page

### Requirement: Lint passes without disabling checks
The `scorecards-site` codebase SHALL pass `yarn lint` (ESLint + Prettier) without any Vue-3-deprecation errors, and without suppressing the underlying ESLint rules that catch them.

#### Scenario: No deprecated Vue 2 API usage remains
- **WHEN** `yarn lint:js` runs against all `.vue` and `.js` files in `scorecards-site/`
- **THEN** it reports zero errors, including zero `vue/no-deprecated-dollar-scopedslots-api` and zero `vue/no-deprecated-destroyed-lifecycle` violations

### Requirement: Dependency vulnerability review passes
The `scorecards-site` dependency tree SHALL contain no packages that fail the repository's configured `dependency-review` vulnerability/license thresholds.

#### Scenario: No flagged advisories in the final lockfile
- **WHEN** the `dependency-review` GitHub Action runs against the updated `yarn.lock`
- **THEN** it reports no vulnerable-package findings (including no `esbuild` advisory) and no OpenSSF Scorecard threshold violations

### Requirement: Home page content renders with feature parity
The site SHALL render the same page content, navigation, and interactive behavior as the pre-migration Nuxt 2 site, sourced from the same `content/home.md` document.

#### Scenario: Markdown content renders
- **WHEN** a user loads the home page
- **THEN** the content from `content/home.md` renders, including custom in-content components (`<sidebar>`, `<code-group>`, `<code-block>`), syntax-highlighted code blocks, and table formatting equivalent to the pre-migration output

#### Scenario: Table of contents scroll-spy behavior
- **WHEN** a user scrolls the home page
- **THEN** the sidebar table-of-contents highlights the currently-visible section, matching pre-migration behavior

#### Scenario: Mobile navigation toggle
- **WHEN** a user opens/closes the mobile navigation menu
- **THEN** the navigation panel shows/hides correctly, matching pre-migration behavior

#### Scenario: Code copy buttons function
- **WHEN** a user clicks a "copy" button on a rendered code block
- **THEN** the code block's contents are copied to the clipboard, matching pre-migration behavior

#### Scenario: Header background color changes on scroll
- **WHEN** a user scrolls past the hero section
- **THEN** the page header's background/text color updates, matching pre-migration behavior

### Requirement: SEO and social metadata are preserved
Page metadata (title, description, canonical URL, Open Graph/Twitter tags, JSON-LD) SHALL continue to be generated from the same environment-derived values as before migration.

#### Scenario: Canonical and social meta tags present
- **WHEN** the home page is server-rendered or statically generated
- **THEN** the output HTML includes a canonical link, Open Graph/Twitter meta tags, and JSON-LD structured data populated from the site's configured frontend URL and site name

### Requirement: Deploy artifact matches Netlify's configured publish path
The static output produced by the site's build command SHALL be written to the directory Netlify is configured to publish.

#### Scenario: Netlify deploy preview serves the built site
- **WHEN** the configured Netlify build command runs for a pull request
- **THEN** the resulting deploy preview serves the generated home page (not a 404 or empty directory listing)
