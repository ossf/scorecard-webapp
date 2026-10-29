## 1. Skeleton, theme, and build pipeline

- [x] 1.1 On `staging-00`, create `site/` with `hugo.yaml` (`baseURL` `https://scorecard.dev/`, English only, Goldmark `unsafe: false`, top menu for Get started, Docs, Use cases, Community, and `/viewer/`, `params.scorecardVersion: v5.6.0`) and add Hextra as a git submodule at `site/themes/hextra` pinned to `v0.13.0`; verify `hugo --source site --panicOnWarning` exits 0 with a `PATH` holding only Hugo, Git, and system utilities.
- [x] 1.2 Add a root content adapter that fails naming `themes/hextra` when the theme is missing; verify that a checkout without submodules makes the build exit non-zero with that error, and record the command and output in the commit message.
- [x] 1.3 Commit FlexSearch 0.8.143 and its Apache-2.0 license under `site/assets/lib/` and set `params.search.flexsearch.js`; verify the built pages load `/lib/flexsearch.bundle.min.*.js` and a build with `cdn.jsdelivr.net` unreachable still exits 0.
- [ ] 1.4 Set `[build]` in `netlify.toml` to `base = "site/"`, `command = "hugo --minify --panicOnWarning"`, `publish = "public"`, and `HUGO_VERSION = "0.154.5"`. Put `HUGO_CASCADE_PARAMS_NOINDEX` under `[context.branch-deploy.environment]` and `[context.deploy-preview.environment]`, and the GoatCounter param under `[context.production.environment]`; verify the next push serves `site/` with `noindex` at `https://staging-00--ossf-scorecard.netlify.app` and scorecard.dev is unchanged.
- [ ] 1.5 Add a `site` job to `.github/workflows/build-website.yml` that reads `HUGO_VERSION` from `netlify.toml`, downloads the matching Hugo extended release, verifies its checksum, checks out submodules, and runs `hugo --source site --minify --panicOnWarning`; verify the job passes on push and `zizmor` reports no findings for the workflow.
- [ ] 1.6 Add a `gitsubmodule` entry to `.github/dependabot.yml` and `/site/` to `.github/CODEOWNERS` with the same owners as `/scorecards-site/`; verify the repository's Dependabot status shows no config error.
- [ ] 1.7 Add brand styling in `site/assets/css/custom.css` (OpenSSF palette, self-hosted Cairo and IBM Plex Sans with their OFL licenses), the Scorecard and OpenSSF logos, the favicon, and a footer with the OpenSSF trademark attribution and the GoatCounter disclosure; verify light and dark screenshots of the home page on the staging preview.
- [x] 1.8 Rely on Hextra's robots tag, which reads the cascaded `noindex` page param; verify a staging page has `noindex` and a local build without the variable has `index, follow`.
- [ ] 1.9 `git mv scorecards-site/static/viewer site/static/viewer`; verify Git records a 100% rename and `/viewer/?uri=github.com/ossf/scorecard` shows results on the staging preview.

## 2. Content guards and their tests

- [x] 2.1 Add a front-matter check partial that calls `errorf` when a Markdown page lacks a text `title` or `description`; add fixtures under `site/tests/` (missing `title`, non-text `description`, unparseable YAML, clean page); verify each bad fixture fails naming the file and field, and the clean one passes.
- [x] 2.2 Add a data check partial for `data/home.yml` (headline, tagline, image, image description) and `data/testimonials.yml` (quote, name, project, `https://` link); add good and bad fixtures; verify each bad fixture fails naming the file, entry, and field.
- [x] 2.3 Add a fixture with `<div>` in prose and one with HTML inside a fenced code block; verify the first fails with `Raw HTML omitted while rendering` and its file name under `--panicOnWarning`, and the second passes.
- [x] 2.4 Add a partial that scans `.RawContent` and calls `errorf` on `{{<` or `{{%` anywhere, and, with fenced code blocks and code spans removed, on a line starting with `:::`; verify fixtures with a shortcode and a directive fail naming the file, and a directive inside a code block passes.
- [x] 2.5 Add a render-link hook that resolves site-relative links with `.GetPage` and checks fragments against the target's `.Fragments.Identifiers`, calling `errorf` on a miss; verify fixtures linking to a missing page and to a missing anchor fail naming the source page, and external links are untouched.
- [x] 2.6 Add a CI job, with inline workflow steps and no script files, that builds each fixture set and asserts the expected exit status and message; verify that disabling one check makes the job fail on that check's bad fixture.

## 3. Pages, data, and layouts

- [ ] 3.1 Add `data/home.yml` and a home layout with the static hero image, a lookup form posting to `/viewer/?uri=`, the three-step summary, the maintainer and consumer paths, and the testimonials; verify on staging that submitting `github.com/ossf/scorecard` opens the viewer.
- [x] 3.2 Add `data/testimonials.yml` with the three quotes, names, projects, and links exactly as on the current home page; verify the home page shows all three.
- [x] 3.3 Write the Get started pages (GitHub Action, command line, authentication and advanced setup, viewer, REST API, badges) from current upstream docs: default `GITHUB_TOKEN`, a PAT only for classic branch protection, `push` and `schedule` triggers, `ghcr.io/ossf/scorecard:latest`; verify each instruction against the `ossf/scorecard` and `ossf/scorecard-action` READMEs and that the build passes.
- [x] 3.4 Write the Docs pages (what Scorecard is and the problem it addresses, how it works, understanding results), describing practice assessment, the Vulnerabilities check's OSV lookup, and "scans more than 1.3 million repositories a week", with no 2021 Synopsys statistic; verify with `grep` that "used by over" and "Synopsys" do not appear in `site/content/`.
- [x] 3.5 Write the two Use cases pages: maintainers improving their own project, and consumers, including organizations, security teams, and registries using results at scale through the REST API; verify the home paths link to both.
- [x] 3.6 Write the Community pages (project name, part of the open source community, get involved); verify they appear under Community in the navigation.
- [x] 3.7 Add an edit and report partial on every page: "Edit this page" to `https://github.com/ossf/scorecard-webapp/edit/main/site/content/<path>`, "Edit upstream" when a page sets `editURL` and `upstream`, and "Report a problem with this page" to a prefilled new issue; verify the links on the command-line guide, the Checks page, and the FAQ.
- [x] 3.8 Ensure each page outputs a meta description from front matter, a canonical URL under `https://scorecard.dev/`, Open Graph title and description, and that `sitemap.xml` is published; verify by inspecting the built HTML of three pages.
- [ ] 3.9 Configure GoatCounter through Hextra's `params.analytics.goatCounter`, set only by the `HUGO_PARAMS_…` variable in the production context (1.4); verify the snippet matches today's (`securityscorecards`, `gc.zgo.at/count.js`) in a local build with the variable set, and is absent from the staging preview.

## 4. Upstream content

- [x] 4.1 Add a partial that fetches a file from `https://raw.githubusercontent.com/ossf/scorecard/<scorecardVersion>/` with `resources.GetRemote` and calls `errorf` with the file name when the fetch fails or returns nothing; verify a build at `v5.6.0` succeeds and a build with a nonexistent tag fails naming the file.
- [x] 4.2 Add the content adapter `site/content/docs/_content.gotmpl` and `site/data/experimental-checks.yaml` (SBOM, Webhooks) that build the Checks page from the fetched `checks.yaml`. Sort by risk then name, show counts of default and experimental checks, read an `experimental` field when present, and set `editURL` to upstream `checks.yaml`. Allow a test-only param to read `checks.yaml` from a fixture; verify all 20 checks appear in order with 18 default and 2 experimental, and a fixture with an unknown risk fails naming the check.
- [x] 4.3 Extend the adapter to render the FAQ and beginner's guide from the fetched files: drop the upstream H1, set title, description, the "from `ossf/scorecard`" notice, and `editURL`; escape `<word>` placeholders as a stopgap; rewrite relative links in the render-link hook to `https://github.com/ossf/scorecard/blob/<tag>/…`; verify `<github_or_gitlab>` renders as text on the FAQ and relative links open GitHub at `v5.6.0`.
- [ ] 4.4 Draft a pull request to `ossf/scorecard` that escapes the three placeholders in `docs/faq.md`, and open it with a maintainer's go-ahead; verify GitHub renders the placeholders as text in the PR's file view.
- [ ] 4.5 Draft a pull request to `ossf/scorecard` that adds `experimental: true` to SBOM and Webhooks in `checks.yaml` and the field to the reader in `docs/checks/internal/reader.go`, and open it with a maintainer's go-ahead; verify `go test ./docs/...` passes and `docs/checks.md` regenerates without unexpected changes.
- [ ] 4.6 Draft a pull request to `ossf/scorecard` adding "bump `scorecardVersion` in `ossf/scorecard-webapp`'s `site/hugo.yaml`" to `RELEASE.md`, and open it with a maintainer's go-ahead; verify the item sits in the post-release steps.

## 5. Contributor documentation

- [ ] 5.1 Ask someone without write access to make one change through "Edit this page" on the staging preview; record which branch the pull request targets, whether a deploy preview builds, and whether the commit carries a DCO sign-off, then close the pull request unmerged.
- [ ] 5.2 Write `CONTRIBUTING.md` covering the web-editor path, DCO sign-off as found in 5.1, how to fix a failed DCO check, and that check text and the guides are edited in `ossf/scorecard`; verify a reviewer can follow it from "Edit this page" to an open pull request.
- [x] 5.3 Update `README.md` to describe `site/` as the website, how to build and preview it, and how upstream content is pinned; verify the commands it lists work as written.

## 6. Pre-cutover verification

- [ ] 6.1 Run Lighthouse in mobile mode against the staging home page; verify performance and accessibility each score at least 0.90, and record the report in the cutover pull request.
- [ ] 6.2 Walk every spec scenario on the staging preview (viewer, edit and report links, search, no video requests, no `api.github.com` requests, metadata, build-failure cases) and record the results in the cutover pull request.
- [ ] 6.3 Have a holder of the GoatCounter `securityscorecards` login confirm that recent hits are being recorded for scorecard.dev, and record the date checked.

## 7. Cutover

- [ ] 7.1 Rebase `staging-00` onto `main` and open the cutover pull requests from it, split by task group if reviewers ask, with the `netlify.toml` change in the last one; verify that pull request's deploy preview serves `site/`, and after it merges scorecard.dev serves `site/`, loads the GoatCounter counter, and has no `noindex`.
- [ ] 7.2 After production is confirmed, `git mv` both hero videos to `site/assets/media/`, then remove `scorecards-site/`, its npm Dependabot entry, and the Nuxt-only CI steps; verify CI passes and the build output contains no `.mp4` files.

## Workflow follow-up

- When a pinned release includes the upstream fixes from 4.4 and 4.5, delete the placeholder escape and `site/data/experimental-checks.yaml`, and verify the FAQ and Checks page still render correctly.
- Archive `migrate-nuxt3-vue3` without syncing its specs, since it describes the removed Nuxt site.
- Archive this change after cutover is confirmed, then verify the archived specs under `openspec/specs/`.
- Delete the `staging-00` branch and its branch deploy after cutover.
