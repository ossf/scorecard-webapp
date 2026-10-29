## Context

See proposal.md for motivation and the specs for required behavior. Constraints that shape the approach:

- **Netlify** deploys the site today from `scorecards-site/` (`netlify.toml`: `base = "scorecards-site/"`, `nuxt generate`, publish `dist`). DNS for scorecard.dev is managed in `ossf/scorecard-infra`, so the site stays on Netlify.
- **Required checks on `main`** (ruleset 22586671): DCO, CodeQL, and dependency-review. Web commits must be signed off (`web_commit_signoff_required`). Actions on pull requests from first-time contributors wait for maintainer approval; Netlify deploy previews do not.
- **The repository has no Go, Python, or shell code.** The only toolchain today is the Nuxt site's Node and Yarn.
- **The viewer** is a standalone `scorecards-site/static/viewer/index.html` that calls `api.scorecard.dev`; it is not Nuxt code.
- **Upstream docs** live in `ossf/scorecard`: `docs/checks/internal/checks.yaml` (embedded in the binary), `docs/faq.md`, and `docs/beginner-checks.md`. The latest release is v5.6.0. `docs/faq.md` has three bare `<placeholder>` strings in prose. `checks.yaml` does not say which checks are experimental. Its reader (`docs/checks/internal/reader.go`) decodes it into a struct.
- **GoatCounter**: the home page loads `//gc.zgo.at/count.js` for `securityscorecards.goatcounter.com`, and the footer discloses it. The viewer has no counter. Whether hits are being recorded has not been checked; the public counter is off.
- **A reference implementation exists**: a Hugo + Hextra mockup with the planned pages, data files, Checks page, guides, edit and report links, and template checks.

## Goals / Non-Goals

**Goals:**
- `main` is unchanged until cutover; `staging-00` carries the new site and its build configuration.
- One tool, Hugo, builds the site and enforces the content rules, so CI and the Netlify preview fail the same way.
- No custom scripts or workflows to maintain.

**Non-Goals:**
- Changing the viewer, or adding analytics to it.
- A visual CMS, translations, or a REST API reference generated from the OpenAPI spec.
- Pixel-matching the current design; the site follows the OpenSSF brand guidelines through the theme.

## Decisions

### Hugo with the Hextra theme
Hugo is a single binary, so builds need no npm tree; that removes the Dependabot churn that dominates the current site's history. Hextra provides the landing-page hero, docs sidebar, table of contents, search, and dark mode, renders GitHub alert syntax through a blockquote render hook, and ships prebuilt CSS.
- *Astro Starlight*: its build validates front matter completely and marks untranslated pages, but it keeps an npm tree (about 350 packages), had 49 releases in a year on a pre-1.0 version, and needs a plugin for GitHub alerts.
- *Other Hugo themes* (Docsy, hugo-book, Relearn, Geekdoc, Lotus Docs, Hinode, Doks): Docsy, Hinode, and Doks need Go or Node. hugo-book, Relearn, and Geekdoc lack a landing-page hero, and Geekdoc and Lotus Docs need shortcodes for callouts. Hextra's top author made 62% of last year's commits, the lowest share among them.
- *Stay on Nuxt*: Nuxt Content schemas do not fail the build on bad front matter, adding a page needs a catch-all route, and the dependency churn stays.

### Theme as a git submodule pinned to a release tag
`site/themes/hextra` is a submodule at `v0.13.0`. Netlify checks out submodules automatically, and Dependabot's `gitsubmodule` ecosystem proposes updates as pull requests. Builds need only Hugo and Git.
- *Hugo module*: `go.sum` checksums, but every build and local preview would need Go just for the theme.
- *Vendored copy*: 200+ committed files, and updates would need Go and be invisible to Dependabot.

### Upstream content fetched at build from a pinned release tag
`hugo.yaml` sets `params.scorecardVersion` (`v5.6.0`). Templates fetch `checks.yaml`, `faq.md`, and `beginner-checks.md` with `resources.GetRemote` from `https://raw.githubusercontent.com/ossf/scorecard/<tag>/`. A failed or missing fetch calls `errorf` with the file name. A test fetched 20 checks and the FAQ in about 1 second, and a nonexistent tag failed the build. The files are used unmodified, and the site documents the release people install. Bumping the tag is a one-line pull request, added to Scorecard's `RELEASE.md` checklist, because Dependabot cannot track a value in `hugo.yaml`.
- *Weekly sync workflow* copying files from `main` into the repository: a custom workflow, bot commits that need DCO sign-off, and pull requests opened with `GITHUB_TOKEN` do not trigger the required Actions checks. It also tracks unreleased content.
- *`ossf/scorecard` as a Hugo module*: file mounts work, but fetching it resolves the whole Go dependency graph. A test filled a 4.6 GB module cache without finishing.
- *Hosting the site in `ossf/scorecard`*: would make upstream files local, but site-only pull requests would face seven required checks (unit, integration, Docker, lint, and others) and a site that follows `main`. Not chosen; the site stays in this repository.

### Content checks inside the Hugo build
- **Front matter and data files**: template partials call `errorf` naming the file and field. In the mockup this caught a missing `title`, a non-text `description`, broken YAML, missing data fields, and a non-`https` link. Hugo turns a numeric `title` into text before templates see it, so that one case passes.
- **Inline HTML**: Goldmark's `unsafe: false` plus `hugo --panicOnWarning`. In a test, raw HTML in prose failed the build with `Raw HTML omitted while rendering "<file>"`, and HTML inside a code block passed.
- **Shortcodes and directives**: a partial scans each page's `.RawContent`. It fails on `{{<` or `{{%` anywhere, because Hugo runs shortcodes before Markdown, even inside code blocks. With fenced code blocks and code spans removed, it fails on a line starting with `:::`.
- **Missing theme**: a checkout without submodules leaves `themes/hextra` empty, and Hugo then fails with a panic about missing layouts. A root content adapter, which runs before rendering, checks for the theme and fails naming it.
- **Broken internal links and anchors**: a render-link hook resolves site-relative links with `.GetPage` and checks fragments against the target page's `.Fragments.Identifiers`, calling `errorf` on a miss.
- **Tests**: fixtures under `site/tests/` hold known-good and known-bad pages and data. A CI job builds each fixture set and asserts the expected exit status and message, so a check that stops failing is caught.

### Checks page and guides from content adapters
`site/content/docs/_content.gotmpl` creates three pages from the fetched files. The Checks page is sorted by risk then name, with counts of default and experimental checks. It reads experimental status from an `experimental` field when the pinned `checks.yaml` has one, and from `site/data/experimental-checks.yaml` until then. The FAQ and the beginner's guide get their title from front matter set by the adapter, with the upstream H1 removed, a "from `ossf/scorecard`" notice, and `editURL` pointing at the upstream file on `main`. Relative links in the guides are rewritten by the render-link hook to `https://github.com/ossf/scorecard/blob/<tag>/…`. Until the pinned release escapes the FAQ's placeholders, the adapter escapes `<word>` patterns in the guides. That is safe for both files today, because neither has placeholders inside code. Both stopgaps are deleted when the pinned release includes the upstream fixes.

### Build on `staging-00` and merge once
All work happens on `staging-00`. Its `netlify.toml` sets `[build]` to `base = "site/"`, `command = "hugo --minify --panicOnWarning"`, `publish = "public"`, and `HUGO_VERSION`. Context-specific values are placed so the same file is right before and after the merge: `HUGO_CASCADE_PARAMS_NOINDEX` under `[context.branch-deploy.environment]` and `[context.deploy-preview.environment]`, and the GoatCounter code under `[context.production.environment]`. Merging `staging-00` into `main` is the cutover.
- *Merge to `main` in steps* behind a `[context.staging-00]` block: smaller reviews, but `main` carries an unused site for weeks. Not chosen, in favor of one cutover.

### Environment-specific output through Hugo params set per Netlify context
Netlify runs every context with Hugo's `production` environment, so `hugo.IsProduction` cannot tell staging from production. Instead, each context sets Hugo configuration through environment variables. `HUGO_PARAMS_ANALYTICS_GOATCOUNTER_CODE` sets Hextra's `params.analytics.goatCounter`, which emits the same counter endpoint and script as today. `HUGO_CASCADE_PARAMS_NOINDEX` cascades a `noindex` page param to every page, and Hextra's head template then emits `noindex, nofollow` instead of `index, follow`. No custom template is needed for either.

### One pinned Hugo version, read from `netlify.toml`
`HUGO_VERSION` (Hugo extended 0.154.5, the version the mockup was verified with) is set once in `netlify.toml`. CI reads that value and installs the matching release archive, verifying its published checksum. Hextra's minimum is 0.146.0.

### Viewer, search, fonts, and videos
- **The viewer** moves once with `git mv` into `site/static/viewer/`. Reverting the cutover restores it.
- **Search**: FlexSearch 0.8.143 (the version Hextra v0.13.0 expects) is committed under `site/assets/lib/` and set through `params.search.flexsearch.js`. By default Hextra fetches it from jsDelivr at build time with no pinned hash.
- **Fonts**: Cairo and IBM Plex Sans are self-hosted from Fontsource under the SIL Open Font License.
- **Videos**: after cutover, the hero videos move to `site/assets/media/`, where Hugo publishes nothing that no template references. They stay in `scorecards-site/` until then, so a rollback still builds.

## Risks / Trade-offs

- **`--panicOnWarning` fails the build on any warning**, including deprecation warnings after a Hugo or theme upgrade → Upgrades are deliberate pull requests with a preview; fix the warning or hold the pin. `ignoreLogs` can silence a specific, understood warning.
- **Builds depend on `raw.githubusercontent.com`** → An outage fails the build loudly instead of shipping stale or empty pages. Retry the deploy; the production site keeps serving the last good build.
- **The site lags `main` until the tag is bumped** → Intended: the site documents the release people install. The `RELEASE.md` checklist item keeps the bump from being forgotten.
- **Dependabot's `gitsubmodule` updates follow the theme's branch, not its release tags** → Its pull requests propose the latest commit. Review each against Hextra's releases and retarget it to a tag, or bump the tag by hand.
- **Single dominant theme maintainer** → Pin a release tag, keep content in portable Markdown, and fork the theme if it stalls.
- **Fast-moving Hugo releases** (41 in a year; 0.146 changed the template system) → One pinned version shared by CI and Netlify.
- **`staging-00` drifts from `main`** → Rebase `staging-00` onto `main` before each review checkpoint and before cutover; `main` requires linear history.
- **A large final review** → Split the cutover into pull requests by task group, with the `netlify.toml` change in the last one.
- **Rebasing rewrites `staging-00`** → It is a staging-only branch; force-push with lease, and nobody bases other work on it.
- **Web-editor sign-off for people without write access is unverified** → Test it on the staging preview before writing that section of `CONTRIBUTING.md`; the guide covers fixing a failed DCO check either way.
- **The upstream fixes may be declined or slow** → The stopgaps keep working; they are small and listed for removal.
- **GoatCounter may not be recording today** → Port it unchanged in behavior, and confirm hits in the dashboard before cutover.

## Migration Plan

1. Build on `staging-00` in task-group order. CI runs on every push, and the branch deploy previews each one. Open the upstream pull requests to `ossf/scorecard` early so the stopgaps can be retired at the next release.
2. Before cutover: all specs pass on the staging preview; Lighthouse thresholds are met; the viewer and edit links work; the GoatCounter dashboard shows hits on production.
3. **Cutover**: rebase `staging-00` onto `main` and open pull requests from it, split by task group if reviewers ask, with the `netlify.toml` change in the last one. Merging that one switches production. **Rollback**: revert it; `scorecards-site/` is still present and unchanged.
4. After production is confirmed: move the hero videos into `site/assets/media/`, then remove `scorecards-site/`, its npm Dependabot entry, and Nuxt-only CI steps. Archive `migrate-nuxt3-vue3` without syncing its specs, then archive this change.

## Open Questions

- Who holds the GoatCounter `securityscorecards` login and can confirm hits before cutover.
- How often to upgrade Hugo after the initial pin.
