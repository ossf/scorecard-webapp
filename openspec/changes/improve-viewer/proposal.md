## Why

The results viewer at `/viewer/` is the page every Scorecard badge links to, and it has the site's most serious defects. A review of the `staging-00` deploy (2026-10-29) found four launch blockers in it: check rows are clickable table cells that keyboards can't reach, focus outlines are removed and sorting is hidden on small screens, the low and medium risk badges have white text at 1.73:1 and 2.51:1 contrast, and Handlebars loads from a CDN as `handlebars@latest` with no integrity hash. The same file serves scorecard.dev today, so these defects are live in production now, not only on the rebuild.

`rebuild-site-hugo` moves the viewer into `site/` unchanged and left its code out of scope. This change brings the viewer in, on `main`, so fixes ship to production without waiting for cutover.

## What Changes

- Load Handlebars from a pinned version (4.7.10) with a Subresource Integrity hash, instead of `@latest`. This is the first fix and ships on its own.
- Make the report keyboard- and screen-reader-accessible: check rows become buttons that expand their details, focus stays visible, sorting works at every width, motion respects `prefers-reduced-motion`, and the report uses semantic structure rather than a layout table.
- Raise the risk badges to at least 4.5:1 contrast.
- Escape everything the API returns before rendering it, including error codes and messages, and percent-encode each part of the repository path before building the API URL.
- Show distinct, plain-language states for loading, not found, rate limited, server error, offline, and malformed responses, each with a next step. Show unavailable (`-1`) and not-applicable scores as words, not `?`.
- Keep the API link, expanded checks, and focus when the report is re-sorted.
- Take check risk levels and documentation links from one list that matches the pinned `ossf/scorecard` release's `checks.yaml`, with a CI check that fails on a duplicate, missing, or unknown check.
- Set the page title and description to the repository being viewed.
- Add filters (failing, risk, unavailable) and default to highest-risk failures first.
- Decide whether viewer reports are indexed, and apply that decision consistently (robots, canonical, sitemap).
- Move the inline CSS, scripts, and templates into versioned files under `viewer/`, cached by hash, so a Content Security Policy can drop inline allowances.
- Add a print stylesheet and a downloadable report that keeps repository, commit, scan date, and Scorecard version.

## Capabilities

### New Capabilities
- `site-viewer`: How the results viewer at `/viewer/` loads, renders, and presents a repository's Scorecard results: its dependencies, accessibility, safe rendering of API data, result states, check metadata, and report export.

### Modified Capabilities
- (none — no specs exist in `openspec/specs/`)

## Impact

- **Changed code**: `scorecards-site/static/viewer/` on `main`. The viewer stays a static page with no build step.
- **New CI check**: compares the viewer's check list with the pinned release's `checks.yaml`.
- **`rebuild-site-hugo`**: `staging-00` already moved the viewer to `site/static/viewer/` as a pure rename. Rebasing `staging-00` onto `main` (its task 7.1) carries these fixes across the rename, and its "moved without modification" requirement then compares against `main`'s updated copy. Its Out of scope line points here.
- **Production**: each task group can ship to scorecard.dev on its own. Badge links (`/viewer/?uri=<repository>`) keep working throughout.
- **Out of scope**: the REST API on `api.scorecard.dev`, and translations.
