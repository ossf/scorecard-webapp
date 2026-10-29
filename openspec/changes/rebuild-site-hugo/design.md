## Context

See proposal.md for motivation and the specs for required behavior. Constraints that shape the approach:

- **Netlify** deploys the site today from `scorecards-site/` (`netlify.toml`: `base = "scorecards-site/"`, `nuxt generate`, publish `dist`). DNS for scorecard.dev is managed in `ossf/scorecard-infra`, so the site stays on Netlify.
- **Required checks on `main`** (ruleset 22586671): DCO, CodeQL, and dependency-review. Web commits must be signed off (`web_commit_signoff_required`).
- **The viewer** is a standalone `scorecards-site/static/viewer/index.html` that calls `api.scorecard.dev`; it is not Nuxt code.
- **GoatCounter**: the home page loads `//gc.zgo.at/count.js` for `securityscorecards.goatcounter.com`, and the footer discloses it. The viewer has no counter. Whether hits are being recorded has not been checked; the public counter is off.
- **Anchors** `#run-the-checks` and `#learn-more` are section IDs in today's `home.md`, used by its sidebar.
- **A reference implementation exists**: a Hugo + Hextra mockup with the planned pages, data files, generated Checks page, mirrored guides, edit and report links, and build-time checks. It builds in about 0.2 seconds with only Hugo on `PATH` and the network blocked.

## Goals / Non-Goals

**Goals:**
- `main` is unchanged until cutover; `staging-00` carries the new site and its build configuration.
- Every guard on content runs inside the Hugo build or in CI, never only on a contributor's machine.
- Each step lands as a small, reviewable pull request.

**Non-Goals:**
- Changing the viewer, or adding analytics to it.
- A visual CMS, translations, or a REST API reference generated from the OpenAPI spec.
- Pixel-matching the current design; the site follows the OpenSSF brand guidelines through the theme.

## Decisions

### Hugo with the Hextra theme
Hugo is a single binary, so builds need no npm tree; that removes the Dependabot churn that dominates the current site's history. Hextra provides the landing-page hero, docs sidebar, table of contents, search, and dark mode. It renders GitHub alert syntax through a blockquote render hook, and its CSS is prebuilt.
- *Astro Starlight*: its build validates front matter completely, and it marks untranslated pages. But it keeps an npm tree (about 350 packages), had 49 releases in a year on a pre-1.0 version, and needs a plugin for GitHub alerts.
- *Other Hugo themes* (Docsy, hugo-book, Relearn, Geekdoc, Lotus Docs, Hinode, Doks): each lacked something Hextra has. Docsy, Hinode, and Doks need Go or Node. hugo-book, Relearn, and Geekdoc lack a landing-page hero, and Geekdoc and Lotus Docs need shortcodes for callouts. Relearn, Docsy, Geekdoc, and hugo-book require Hugo 0.158 or newer. Hextra's top author made 62% of last year's commits, the lowest share among them.
- *Stay on Nuxt*: Nuxt Content schemas do not fail the build on bad front matter, adding a page needs a catch-all route, and the dependency churn stays.

### Theme as a git submodule pinned to a release tag
`site/themes/hextra` is a submodule at `v0.13.0`. Netlify checks out submodules automatically, and Dependabot's `gitsubmodule` ecosystem proposes updates as pull requests.
- *Hugo module*: tracks versions in `go.mod` but needs Go at build time.
- *Vendored copy*: builds with Hugo alone, but 200+ committed files, and updates need Go and are invisible to Dependabot.

### One pinned Hugo version, read from `netlify.toml`
`HUGO_VERSION` (Hugo extended 0.154.5, the version the mockup was verified with) is set once in `netlify.toml`. CI reads that value and installs the matching release archive, verifying its published checksum, so Netlify and CI cannot drift. Hextra's minimum is 0.146.0. Upgrades are deliberate pull requests that change the one value.

### Build on `staging-00` and merge once
All work happens on `staging-00`. Its `netlify.toml` sets `[build]` to `base = "site/"`, `command = "hugo --minify"`, `publish = "public"`, and `HUGO_VERSION`, so its Netlify branch deploy previews every push. Context-specific values are placed so the same file is right before and after the merge: `noindex` under `[context.branch-deploy.environment]` and `[context.deploy-preview.environment]`, and the GoatCounter code under `[context.production.environment]`. Merging `staging-00` into `main` is the cutover. Production then builds `site/` with analytics on and no `noindex`, with no further edit.
- *Merge to `main` in steps* behind a `[context.staging-00]` block: smaller reviews, but `main` carries an unused site for weeks. Not chosen, in favor of one cutover.
- *Second Netlify site*: continuous previews, but it needs a new site created in Netlify.

### Environment-specific output through Hugo params set per Netlify context
Netlify runs every context with Hugo's `production` environment, so `hugo.IsProduction` cannot tell staging from production. Instead, the production context sets the GoatCounter site code as a Hugo parameter through an environment variable (`HUGO_PARAMS_…`), and the branch-deploy and deploy-preview contexts set a `noindex` parameter the same way. Templates emit the counter only when the code is present, and the robots `noindex` tag only when the flag is set. Hextra already supports `params.analytics.goatCounter`; its output is checked against today's snippet.

### Content checks: build-time templates plus two CI linters
- **Front matter and data files** are checked in Hugo templates with `errorf`, so a bad page fails the Netlify build and CI with a message naming the file and field. In the mockup this caught a missing `title`, a non-text `description`, broken YAML, missing data fields, a non-`https` link, missing `checks.yaml`, and an unknown risk. Hugo turns a numeric `title` into text before templates see it, so that one case passes.
- **Inline HTML** is caught by `markdownlint-cli2` with MD033 enabled for `site/content/`. CI pins the version, and `site/` gets its own config, because `scorecards-site/.markdownlint.yaml` disables MD033.
- **Shortcodes and directives** (`{{<`, `{{%`, `{%`, `::`) outside code are caught by a small Python standard-library script.
- **Broken links** are caught by running `lychee` offline against the built `public/` directory, with fragment checking.
- Each check has known-good and known-bad fixtures under `site/tests/`. A CI job asserts that the good fixtures pass and that each bad fixture fails with its expected message.

### Checks page from a content adapter
`site/content/docs/_content.gotmpl` builds the Checks page from `site/data/checks.yaml`, sorting by risk then name. No generated Markdown is committed, and the page's headings appear in the table of contents. Experimental checks come from a small `site/data/experimental-checks.yaml` override list, because `checks.yaml` does not mark them.

### Guides transformed at sync time, not build time
The sync script (Python standard library) copies upstream guides and transforms them: it drops the title and hand-written table of contents, points relative links at GitHub, escapes `<placeholder>` text outside code, and adds the "copied from" notice and `editURL`. Doing this in Go templates would need code-fence-aware escaping, which is fragile. The committed result is reviewable in the sync pull request.

### Sync workflow opens pull requests with the repository token
A scheduled workflow (weekly plus `workflow_dispatch`) runs the sync script. GitHub runs scheduled workflows only from the default branch, and manual runs also need the file on `main`, so the workflow first runs after cutover; until then the script runs by hand on `staging-00`. The workflow signs off its commit as `github-actions[bot]`, and opens or updates one pull request with `gh`. Permissions are `contents: write` and `pull-requests: write` on that job only. Actions are pinned by SHA, checkout credentials are not persisted, and `zizmor` must pass.

### Legacy anchors in the home template
The home layout wraps its getting-started paths in an element with `id="run-the-checks"` and its learn-more links in `id="learn-more"`. Fragments never reach the server, so redirects cannot handle them; the IDs live in the template, not in content.

### Search, fonts, and assets served locally
FlexSearch 0.8.143 (the version Hextra v0.13.0 expects) is committed under `site/assets/lib/` and set through `params.search.flexsearch.js`. By default Hextra fetches it from jsDelivr at build time, and an offline build fails. Cairo and IBM Plex Sans are self-hosted from Fontsource under the SIL Open Font License. After cutover, the hero videos move to `site/assets/media/`, where Hugo publishes nothing that no template references; they stay in `scorecards-site/` until then so a rollback still builds.

## Risks / Trade-offs

- **Single dominant theme maintainer** → Pin a release tag, keep content in portable Markdown, and fork the theme if it stalls. Templates the site overrides stay small.
- **Fast-moving Hugo releases** (41 in a year; 0.146 changed the template system) → One pinned version shared by CI and Netlify, upgraded deliberately with a preview.
- **`staging-00` drifts from `main`** (Dependabot and workflow changes keep landing there) → Rebase `staging-00` onto `main` before each review checkpoint and before cutover; `main` requires linear history, so the cutover pull request must apply cleanly.
- **A large final review** → Split the cutover into pull requests by task group, with the `netlify.toml` change in the last one. Reviewers can follow `staging-00` along the way through its preview and the compare view.
- **Rebasing rewrites `staging-00`** → It is a staging-only branch; force-push with lease, and nobody bases other work on it.
- **Pull requests opened with `GITHUB_TOKEN` do not trigger Actions workflows**, so the required CodeQL and dependency-review checks will not report on sync pull requests → A maintainer closes and reopens the pull request to trigger them, documented in the workflow. A GitHub App token would remove this step (see Open Questions).
- **Web-editor sign-off for people without write access is unverified** → Test it on the staging preview before writing that section of `CONTRIBUTING.md`; the guide covers fixing a failed DCO check either way.
- **Lighthouse and contrast thresholds are untested with the brand colors** → Measure on the first staging deploy, and adjust colors in `site/assets/css/custom.css` before content work.
- **Both directories exist on `main` between cutover and clean-up** → Remove `scorecards-site/` as soon as production is confirmed; until then the README says `site/` is live, and npm Dependabot entries stay scoped to `scorecards-site/`.
- **GoatCounter may not be recording today** → Port it unchanged in behavior, and confirm hits in the dashboard before cutover.

## Migration Plan

1. Build on `staging-00` in task-group order: skeleton, theme, `netlify.toml`, and CI; content checks with fixtures; pages and data; Checks page, sync script, and mirrored guides; sync workflow; `CONTRIBUTING.md`. CI runs on every push, and the branch deploy previews each one.
2. Before cutover: all specs pass on the staging preview; Lighthouse and accessibility thresholds are met; the viewer, anchors, and edit links work; the GoatCounter dashboard shows hits on production.
3. **Cutover**: rebase `staging-00` onto `main` and open pull requests from it, split by task group if reviewers ask, with the `netlify.toml` change in the last one. Merging that one switches production. **Rollback**: revert it; `scorecards-site/` is still present and unchanged.
4. After production is confirmed: move the hero videos into `site/assets/media/`, then remove `scorecards-site/`, its npm Dependabot entry, and Nuxt-only CI steps. Archive `migrate-nuxt3-vue3` without syncing its specs, then archive this change.

## Open Questions

- Whether sync pull requests should use a GitHub App token, so required Actions checks run without a manual close and reopen. This needs an app installed on the repository.
- Who holds the GoatCounter `securityscorecards` login and can confirm hits before cutover.
- How often to upgrade Hugo after the initial pin.
