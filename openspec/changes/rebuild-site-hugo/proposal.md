## Why

scorecard.dev is hard to keep accurate and hard for non-code contributors to edit. Its content is one 350-line `scorecards-site/content/home.md` that mixes Markdown with 50 inline HTML elements, Tailwind classes, and Vue components. The hero text is hard-coded in `pages/index.vue`. Upkeep is dominated by the framework: 15 of the last 35 commits to `scorecards-site/` came from Dependabot, and the lockfile resolves about 1,250 packages. The Nuxt 2 → 3 bump forced a full migration in August 2026. The checks table is copied by hand from upstream (`ossf/scorecard`) and omits the two experimental checks that `checks.yaml` defines.

Rebuilding the site as plain Markdown on Hugo with the Hextra theme removes the npm dependency tree. Contributors can change wording through GitHub's web editor, and the Checks page and user guides come straight from an `ossf/scorecard` release. A working mockup with the same pages and data builds with Hugo alone in about 0.2 seconds.

## What Changes

- Add a new site in `site/`, built with Hugo extended and the Hextra theme (a git submodule pinned to v0.13.0). Builds need only Hugo and Git; there is no Node, Go, Python, or Sass. The site is developed on the `staging-00` branch and previewed through that branch's Netlify branch deploy; `main` is unchanged until the site is complete.
- Content is GitHub-flavored Markdown with `title` and `description` front matter only. Callouts use GitHub's alert syntax (`> [!NOTE]`). Pages contain no inline HTML, MDX, or shortcodes. Hero, testimonials, and similar structured content live in YAML data files.
- Restructure the site from one long page into sections (Home, Get started, Docs, Use cases, Community) following the user stories in #960. The home page offers one path for maintainers and one for consumers.
- Generate the Checks page and render the FAQ and beginner's guide at build time from files fetched from a pinned `ossf/scorecard` release tag, used unmodified. Bumping the tag is added to Scorecard's release checklist.
- Every page gets "Edit this page" and "Report a problem with this page" links. The Hugo build itself fails on inline HTML or shortcodes in content, missing or malformed front matter, malformed data files, missing upstream files, and broken internal links.
- Add `CONTRIBUTING.md` covering the web-editor path and DCO sign-off, based on a test of the web editor by someone without write access.
- Move the viewer once into `site/static/viewer/`, unchanged.
- Replace the two hero videos with a static image; keep the video files in the repository as unpublished assets.
- Port the GoatCounter analytics snippet and its footer disclosure, emitted only on the production site.
- Propose two fixes upstream in `ossf/scorecard`: escape the bare placeholders in `docs/faq.md`, and mark experimental checks in `checks.yaml`. The site carries a stopgap for each until the pinned release includes it.
- **BREAKING**: Cutover is merging `staging-00` into `main`: the `netlify.toml` it carries builds `site/` for production. After production is confirmed, `scorecards-site/` and its Nuxt-only CI and Dependabot entries are removed.

## Capabilities

### New Capabilities
- `site-platform`: How the site is built, previewed, and deployed — the toolchain a build may need, pinned versions, Netlify contexts, the viewer, analytics, search, metadata, and the performance and accessibility thresholds the site must meet.
- `site-content`: The content model — what contributors may write in pages, the data files, the site map, and the content taken from a pinned `ossf/scorecard` release (the generated Checks page and the user guides).
- `site-contribution`: How people change the site — edit and report links on every page, the build-time checks that guard content, and the contributor documentation.

### Modified Capabilities
- (none — no specs exist in `openspec/specs/`)

## Impact

- **New code**: `site/` (Hugo config, layouts, content adapters, content, data, static assets, the viewer, the Hextra submodule, a local FlexSearch copy, test fixtures) and `CONTRIBUTING.md`.
- **Changed config**: `netlify.toml` on `staging-00` (`[build]` builds `site/`; `noindex` for branch deploys and deploy previews; GoatCounter for production), reaching `main` at cutover; `.gitmodules`; `.github/dependabot.yml` (`gitsubmodule` ecosystem); `.github/CODEOWNERS` (`/site/`); `build-website.yml` (a job that builds `site/`).
- **Upstream (`ossf/scorecard`)**: pull requests to escape placeholders in `docs/faq.md`, add an `experimental` field to `checks.yaml` and its reader, and add the site tag bump to `RELEASE.md`.
- **Removed after cutover**: `scorecards-site/`, its npm Dependabot entry, and the Nuxt-only CI steps.
- **In-flight change**: `migrate-nuxt3-vue3` specifies `scorecards-site-platform` for the Nuxt site. Its specs have not been synced. After cutover it describes removed code, so archive it without syncing its specs.
- **Out of scope**: the viewer's own code, the REST API and its docs on `api.scorecard.dev`, translations (English only for now), and a visual CMS.
