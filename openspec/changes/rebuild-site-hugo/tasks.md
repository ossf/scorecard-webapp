## 1. Skeleton, theme, and build pipeline

- [ ] 1.1 Create `site/` with `hugo.yaml` (`baseURL` `https://scorecard.dev/`, English only, Goldmark `unsafe: false`, top menu for Get started, Docs, Use cases, Community, and `/viewer/`) and add Hextra as a git submodule at `site/themes/hextra` pinned to `v0.13.0`; verify `hugo --source site` exits 0 with only Hugo on `PATH`.
- [ ] 1.2 Verify that a checkout without submodules makes `hugo --source site` exit non-zero with an error naming the missing theme, and record the command and output in the pull request.
- [ ] 1.3 Commit FlexSearch 0.8.143 and its Apache-2.0 license under `site/assets/lib/` and set `params.search.flexsearch.js`; verify a build with outbound network blocked exits 0, and the same build without the setting fails with `Could not retrieve FlexSearch js file`.
- [ ] 1.4 On `staging-00`, set `[build]` in `netlify.toml` to `base = "site/"`, `command = "hugo --minify"`, `publish = "public"`, and `HUGO_VERSION = "0.154.5"`; put the `noindex` Hugo param under `[context.branch-deploy.environment]` and `[context.deploy-preview.environment]`, and the GoatCounter param under `[context.production.environment]`; verify the next push serves `site/` with `noindex` at `https://staging-00--ossf-scorecard.netlify.app` and scorecard.dev is unchanged.
- [ ] 1.5 Add a `site` job to `.github/workflows/build-website.yml` that reads `HUGO_VERSION` from `netlify.toml`, downloads the matching Hugo extended release, verifies its checksum, checks out submodules, and builds `site/`; verify the job passes on the pull request and that `zizmor` reports no findings for the workflow.
- [ ] 1.6 Add a `gitsubmodule` entry to `.github/dependabot.yml` and `/site/` to `.github/CODEOWNERS` with the same owners as `/scorecards-site/`; verify the Dependabot config is accepted (the repository's Dependabot status shows no config error).
- [ ] 1.7 Add brand styling in `site/assets/css/custom.css` (OpenSSF palette, self-hosted Cairo and IBM Plex Sans with their OFL licenses), the Scorecard and OpenSSF logos, the favicon, and a footer with the OpenSSF trademark attribution and the GoatCounter disclosure; verify light and dark screenshots of the home page on the staging preview.
- [ ] 1.8 Emit a robots `noindex` meta tag only when the `noindex` param is set; verify a staging page has it and a local build without the param does not.

## 2. Content guards and their tests

- [ ] 2.1 Add a front-matter check partial that fails the build with `errorf` when a Markdown page lacks a text `title` or `description`; add fixtures under `site/tests/` (missing `title`, non-text `description`, unparseable YAML, and a clean page); verify each bad fixture fails with a message naming the file and field, and the clean one passes.
- [ ] 2.2 Add a data check partial for `data/home.yml` (headline, tagline, image, image description) and `data/testimonials.yml` (quote, name, project, `https://` link); add good and bad fixtures; verify each bad fixture fails with a message naming the file, entry, and field.
- [ ] 2.3 Add `site/.markdownlint-cli2.jsonc` enabling MD033 for `site/content/` and a `site` job in `.github/workflows/linter.yml` running a pinned `markdownlint-cli2`; verify a fixture with `<div>` fails with file and line, and the clean fixture passes.
- [ ] 2.4 Add `site/scripts/check-content.py` (Python standard library) that rejects `{{<`, `{{%`, `{%`, and `::` directives outside code blocks and spans, run it in the linter job; verify its good and bad fixtures, including a shortcode inside a code block that must pass.
- [ ] 2.5 Run `lychee` offline with fragment checking against the built `public/` in CI; verify a fixture page linking to a missing page and to a missing anchor fails, naming the source page.
- [ ] 2.6 Add `site/tests/run.sh`, which asserts that every good fixture passes and every bad fixture fails with its expected message, and run it in CI; verify that commenting out one check makes the run fail on that check's bad fixture.

## 3. Pages, data, and layouts

- [ ] 3.1 Add `data/home.yml` and a home layout with the static hero image, a lookup form posting to `/viewer/?uri=`, the three-step summary, the maintainer and consumer paths, and the testimonials, wrapping the paths in `id="run-the-checks"` and the learn-more links in `id="learn-more"`; verify on staging that submitting `github.com/ossf/scorecard` opens the viewer and that `/#run-the-checks` and `/#learn-more` scroll to those sections.
- [ ] 3.2 Add `data/testimonials.yml` with the three quotes, names, projects, and links exactly as on the current home page; verify the home page shows all three.
- [ ] 3.3 Write the Get started pages (GitHub Action, command line, authentication and advanced setup, viewer, REST API, badges) from current upstream docs: default `GITHUB_TOKEN`, a PAT only for classic branch protection, `push` and `schedule` triggers, `ghcr.io/ossf/scorecard:latest`; verify each instruction against `ossf/scorecard` and `ossf/scorecard-action` READMEs and that the content checks pass.
- [ ] 3.4 Write the Docs pages (what Scorecard is and the problem it addresses, how it works, understanding results), describing practice assessment, the Vulnerabilities check's OSV lookup, and "scans more than 1.3 million repositories a week", with no 2021 Synopsys statistic; verify with `grep` that "used by over" and "Synopsys" do not appear in `site/content/`.
- [ ] 3.5 Write the two Use cases pages: maintainers improving their own project, and consumers, including organizations, security teams, and registries using results at scale through the REST API; verify the home paths link to both.
- [ ] 3.6 Write the Community pages (project name, part of the open source community, get involved); verify they appear under Community in the navigation.
- [ ] 3.7 Add an edit and report partial on every page: "Edit this page" to `https://github.com/ossf/scorecard-webapp/edit/main/site/content/<path>`, "Edit upstream" when a page sets `editURL` and `upstream`, and "Report a problem with this page" to a prefilled new issue; verify the links on the command-line guide, the Checks page, and the FAQ.
- [ ] 3.8 Ensure each page outputs a meta description from front matter, a canonical URL under `https://scorecard.dev/`, Open Graph title and description, and that `sitemap.xml` is published; verify by inspecting the built HTML of three pages.
- [ ] 3.9 Configure GoatCounter through Hextra's `params.analytics.goatCounter`, set only by the `HUGO_PARAMS_…` variable in the production context (1.4); verify the snippet matches today's (`securityscorecards`, `gc.zgo.at/count.js`) in a local build with the variable set, and is absent from the staging preview.

## 4. Upstream content and sync

- [ ] 4.1 Add `site/scripts/sync.py` (Python standard library) that copies `checks.yaml` unchanged to `site/data/checks.yaml`, transforms the FAQ and beginner's guide into `site/content/docs/`, and copies the viewer to `site/static/viewer/index.html`; verify a run against current sources succeeds and a run with a missing source exits 1 naming the file without writing anything.
- [ ] 4.2 Add the content adapter `site/content/docs/_content.gotmpl` and `site/data/experimental-checks.yaml` (SBOM, Webhooks) that build the Checks page sorted by risk then name, with counts of default and experimental checks, how to enable experimental ones, and `editURL` to upstream `checks.yaml`; verify all 20 checks appear in order and that a fixture with an unknown risk fails the build naming the check.
- [ ] 4.3 Run the sync and commit its output; verify the viewer copy is byte-identical to its source, `/viewer/?uri=github.com/ossf/scorecard` shows results on staging, and `<github_or_gitlab>` renders as text on the FAQ.
- [ ] 4.4 Add `.github/workflows/sync-upstream.yml` (weekly cron plus `workflow_dispatch`) that runs the sync script, commits with sign-off as `github-actions[bot]`, and opens or updates one pull request with `gh`, whose body asks reviewers to check the experimental list and explains reopening it to trigger required checks; pin actions by SHA, scope `contents: write` and `pull-requests: write` to the job, and don't persist credentials; verify `zizmor` reports no findings. Scheduled and manual runs need the workflow on `main`, so the run test is in 7.2.

## 5. Contributor documentation

- [ ] 5.1 Ask someone without write access to make one change through "Edit this page" on the staging preview; record which branch the pull request targets, whether a deploy preview builds, and whether the commit carries a DCO sign-off, then close the pull request unmerged.
- [ ] 5.2 Write `CONTRIBUTING.md` covering the web-editor path, DCO sign-off as found in 5.1, how to fix a failed DCO check, and that check text is edited upstream; verify a reviewer can follow it from "Edit this page" to an open pull request.
- [ ] 5.3 Update `README.md` to describe `site/` as the website, how to build and preview it, and how upstream content is synced; verify the commands it lists work as written.

## 6. Pre-cutover verification

- [ ] 6.1 Run Lighthouse in mobile mode against the staging home page; verify performance and accessibility are each at least 0.95 with no `color-contrast` failures, and record the report in the cutover pull request.
- [ ] 6.2 Walk every spec scenario on the staging preview (viewer, anchors, edit and report links, offline search, no video requests, no `api.github.com` requests, metadata) and record the results in the cutover pull request.
- [ ] 6.3 Have a holder of the GoatCounter `securityscorecards` login confirm that recent hits are being recorded for scorecard.dev, and record the date checked.

## 7. Cutover

- [ ] 7.1 Rebase `staging-00` onto `main` and open the cutover pull requests from it, split by task group if reviewers ask, with the `netlify.toml` change in the last one; verify that pull request's deploy preview serves `site/`, and after it merges scorecard.dev serves `site/`, loads the GoatCounter counter, and has no `noindex`.
- [ ] 7.2 After production is confirmed, `git mv` both hero videos to `site/assets/media/`, then remove `scorecards-site/`, its npm Dependabot entry, and the Nuxt-only CI steps; verify CI passes, the build output contains no `.mp4` files, and a manual run of the sync workflow with no upstream change opens no pull request.

## Workflow follow-up

- Archive `migrate-nuxt3-vue3` without syncing its specs, since it describes the removed Nuxt site.
- Archive this change after cutover is confirmed, then verify the archived specs under `openspec/specs/`.
- Delete the `staging-00` branch and its branch deploy after cutover.
