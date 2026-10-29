## Why

scorecard.dev is hard to keep accurate and hard for non-code contributors to edit. Its content is one 350-line `scorecards-site/content/home.md` that mixes Markdown with 50 inline HTML elements, Tailwind classes, and Vue components. The hero text is hard-coded in `pages/index.vue`. Upkeep is dominated by the framework: 15 of the last 35 commits to `scorecards-site/` came from Dependabot, and the lockfile resolves about 1,250 packages. The Nuxt 2 → 3 bump forced a full migration in August 2026. The checks table is copied by hand from upstream (`ossf/scorecard`) and omits the two experimental checks that `checks.yaml` defines.

Rebuilding the informational pages as plain Markdown on Hugo with the Hextra theme removes the npm dependency tree. Contributors can change wording through GitHub's web editor, and the checks page is generated from upstream data. A working mockup with the same pages and data builds with Hugo alone in about 0.2 seconds.

## What Changes

- Add a new site in `site/`, built with Hugo extended and the Hextra theme (a git submodule pinned to v0.13.0). It is developed on the `staging-00` branch and previewed through that branch's Netlify branch deploy; `main` is unchanged until the site is complete.
- Content is GitHub-flavored Markdown with `title` and `description` front matter only. Callouts use GitHub's alert syntax (`> [!NOTE]`). Pages contain no inline HTML, MDX, or shortcodes. Hero, testimonials, and similar structured content live in YAML data files.
- Restructure the site from one long page into sections (Home, Get started, Docs, Use cases, Community) following the user stories in #960. The home page offers one path for maintainers and one for consumers. The viewer at `/viewer/` is served unchanged.
- Generate the Checks page at build time from an unmodified copy of upstream `checks.yaml`, sorted from Critical to Low risk. Mirror upstream's FAQ and beginner's guide with "Edit upstream" links.
- Add a scheduled workflow that copies `checks.yaml`, the mirrored guides, and the viewer from their sources and opens a pull request when they change.
- Every page gets "Edit this page" and "Report a problem with this page" links. CI fails a pull request on inline HTML or shortcodes in content, missing or malformed front matter, malformed data files, or broken links.
- Add `CONTRIBUTING.md` covering the web-editor path and DCO sign-off, based on a test of the web editor by someone without write access.
- Replace the two hero videos with a static image; keep the video files in the repository as unpublished assets.
- Port the GoatCounter analytics snippet and its footer disclosure, emitted only on the production site.
- Preview the new site through a Netlify branch deploy of `staging-00`, configured by that branch's own `netlify.toml`.
- **BREAKING**: Cutover is merging `staging-00` into `main`: the `netlify.toml` it carries builds `site/` for production. After production is confirmed, `scorecards-site/` and its Nuxt-only CI and Dependabot entries are removed. Anchors `/#run-the-checks` and `/#learn-more` keep landing on matching content.

## Capabilities

### New Capabilities
- `site-platform`: How the site is built, previewed, and deployed — Hugo and theme versions, the toolchain a build may need, Netlify contexts, the unchanged viewer, analytics, legacy anchors, and the performance and accessibility thresholds the site must meet.
- `site-content`: The content model — what contributors may write in pages, the data files, the site map, the generated Checks page, and the mirrored upstream guides.
- `site-contribution`: How people change the site — edit and report links on every page, the CI checks that guard content, and the contributor documentation.
- `site-upstream-sync`: How upstream sources reach the site — the scheduled workflow that copies them and opens a reviewable pull request.

### Modified Capabilities
- (none — no specs exist in `openspec/specs/`)

## Impact

- **New code**: `site/` (Hugo config, layouts, content, data, static assets, the Hextra submodule, a local FlexSearch copy, the sync script), `CONTRIBUTING.md`, a sync workflow under `.github/workflows/`.
- **Changed config**: `netlify.toml` on `staging-00` (`[build]` builds `site/`; `noindex` for branch deploys and deploy previews; GoatCounter for production), reaching `main` at cutover; `.github/dependabot.yml` (`gitsubmodule` ecosystem), `.github/CODEOWNERS` (`/site/`), `.gitmodules`, `build-website.yml` and `linter.yml` (jobs for `site/`).
- **Removed at cutover**: `scorecards-site/`, its npm Dependabot entry, and the Nuxt-only CI steps.
- **In-flight change**: `migrate-nuxt3-vue3` specifies `scorecards-site-platform` for the Nuxt site. Its specs have not been synced. After cutover it describes removed code, so archive it without syncing its specs.
- **Out of scope**: the viewer's own code, the REST API and its docs on `api.scorecard.dev`, translations (English only for now), and a visual CMS.
