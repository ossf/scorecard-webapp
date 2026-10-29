## Purpose

Defines how the scorecard.dev site in `site/` is built, previewed, and deployed, and the compatibility, analytics, performance, and accessibility guarantees every build must keep.

## ADDED Requirements

### Requirement: Build needs only Hugo and Git
Building the site from a checkout SHALL require only Git and the pinned Hugo extended binary. A build SHALL NOT require Node.js, npm, Go, Python, Dart Sass, or PostCSS. Network access during the build SHALL be limited to fetching the pinned `ossf/scorecard` release files.

#### Scenario: Build with only Hugo and Git on PATH
- **WHEN** the repository is cloned with submodules and the site is built with a `PATH` containing only Hugo, Git, and system utilities
- **THEN** the build exits successfully and writes the complete site to its output directory

### Requirement: A broken build fails instead of publishing a broken site
A build that is missing its theme, a data file, or a pinned upstream file SHALL exit with a non-zero status and an error that names what is missing. It SHALL NOT exit successfully with an empty, unstyled, or partial site.

#### Scenario: Theme submodule not checked out
- **WHEN** the site is built from a checkout without submodules
- **THEN** the build exits non-zero and reports the missing theme

#### Scenario: Pinned upstream tag does not exist
- **WHEN** the pinned `ossf/scorecard` tag names a release that does not exist
- **THEN** the build exits non-zero and names the file it could not fetch

### Requirement: Versions are pinned
The Hugo version, the theme version, and the `ossf/scorecard` release the site takes content from SHALL each be pinned. Netlify builds and CI builds SHALL use the same Hugo version. A theme update SHALL arrive as an automated pull request.

#### Scenario: CI and Netlify agree
- **WHEN** CI builds the site
- **THEN** it uses the Hugo version set in `netlify.toml`

#### Scenario: Theme update
- **WHEN** the theme publishes a new release
- **THEN** an automated pull request proposes the new pinned version, and nothing changes without it

### Requirement: Staging preview of the new site
Pushes to the `staging-00` branch SHALL publish the new site from `site/` as a Netlify branch deploy. Branch deploys and deploy previews SHALL mark every page `noindex`; production SHALL NOT.

#### Scenario: Push to staging-00
- **WHEN** a commit is pushed to `staging-00`
- **THEN** `https://staging-00--ossf-scorecard.netlify.app` serves the new site built from that commit's `site/` directory

#### Scenario: Staging is not indexed
- **WHEN** a search engine fetches any staging page
- **THEN** the page carries a robots `noindex` directive

#### Scenario: Production is indexable after cutover
- **WHEN** a search engine fetches a page on scorecard.dev after cutover
- **THEN** the page carries no robots `noindex` directive

### Requirement: Production is unchanged until cutover
Until cutover, `main` SHALL keep building `scorecards-site/` for production and for pull-request deploy previews. Work on `staging-00` SHALL NOT change what scorecard.dev serves.

#### Scenario: Push to staging-00 before cutover
- **WHEN** commits that change `site/` and `netlify.toml` are pushed to `staging-00`
- **THEN** scorecard.dev still serves the site built from `scorecards-site/`

### Requirement: Cutover switches production to the new site
Merging the `netlify.toml` change from `staging-00` into `main` SHALL switch production to `site/`, with analytics on and no `noindex`, without further edits. After production is confirmed, `scorecards-site/` and the CI jobs and dependency-update entries that only apply to it SHALL be removed.

#### Scenario: Production after cutover
- **WHEN** the cutover pull request from `staging-00` is merged to `main`
- **THEN** scorecard.dev serves the site built from `site/`, loads the GoatCounter counter, and pull-request deploy previews build `site/`

#### Scenario: Old site removed
- **WHEN** production has been confirmed after cutover
- **THEN** the repository no longer contains `scorecards-site/` or CI steps and dependency-update entries that only applied to it

### Requirement: Viewer is served unchanged
The results viewer SHALL be served at `/viewer/` from the existing `index.html`, moved into `site/` without modification, and `/viewer/?uri=<repository>` SHALL keep working, because every Scorecard badge links there.

#### Scenario: Badge link
- **WHEN** a visitor opens `/viewer/?uri=github.com/ossf/scorecard` on a preview or production deploy
- **THEN** the viewer loads and shows that repository's results

#### Scenario: Viewer content unchanged
- **WHEN** the moved viewer file is compared with `scorecards-site/static/viewer/index.html` before the move
- **THEN** the two are byte-identical

### Requirement: Analytics run only in production, with disclosure
The production site SHALL load the GoatCounter counter for the `securityscorecards` site. Branch deploys and deploy previews SHALL NOT load it. Every page SHALL show the footer disclosure that GoatCounter is used.

#### Scenario: Production page
- **WHEN** a visitor loads a page on scorecard.dev after cutover
- **THEN** the page loads the GoatCounter script configured for `securityscorecards`

#### Scenario: Preview page
- **WHEN** a visitor loads a page on the staging branch deploy or a deploy preview
- **THEN** the page does not load the GoatCounter script, and the footer disclosure is still shown

### Requirement: Pages meet performance and accessibility thresholds
On a mobile Lighthouse run of the home page, performance and accessibility SHALL each score at least 0.90. Pages SHALL NOT load video files or make browser-side requests to `api.github.com`.

#### Scenario: Lighthouse on the preview
- **WHEN** Lighthouse runs in mobile mode against the home page of a preview deploy
- **THEN** performance and accessibility each score at least 0.90

#### Scenario: No video or GitHub API calls
- **WHEN** any page loads
- **THEN** no video file is requested and no request goes to `api.github.com`

### Requirement: Retained video assets are not published
The existing hero videos SHALL be kept in the repository under `site/` for possible reuse, and SHALL NOT appear in the build output.

#### Scenario: Build output
- **WHEN** the site is built after the videos have moved into `site/`
- **THEN** the output directory contains no `.mp4` files, and the repository still contains both hero videos

### Requirement: Search works without third-party hosts
Site search SHALL work from files committed to the repository and served by the site. Neither the build nor the browser SHALL fetch the search library from a CDN.

#### Scenario: Searching a page
- **WHEN** a visitor searches for "Pinned-Dependencies"
- **THEN** the Checks page appears in the results, and no request for the search library went to a third-party host

### Requirement: Pages carry search and social metadata
Every page SHALL include a title, a meta description from its front matter, a canonical URL on scorecard.dev, and Open Graph tags. The site SHALL publish a sitemap.

#### Scenario: Page metadata
- **WHEN** any page's HTML is inspected
- **THEN** it has a `<title>`, a meta description matching the page's `description`, a canonical link under `https://scorecard.dev/`, and Open Graph title and description tags
