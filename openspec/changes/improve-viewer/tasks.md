## 1. Pin Handlebars

- [ ] 1.1 In `scorecards-site/static/viewer/index.html`, replace `handlebars@latest` with `handlebars@4.7.10` and add an `integrity` attribute (SHA-384 of the file at that URL) and `crossorigin="anonymous"`; verify on the deploy preview that `/viewer/?uri=github.com/ossf/scorecard` renders results and the browser console shows no integrity error.
- [ ] 1.2 Add a `build-website.yml` step that fails when any file under `scorecards-site/static/viewer/` contains `@latest` or a cross-origin `<script>` without `integrity`; verify it fails on a branch that restores `@latest` and passes on this one.

## 2. Accessibility blockers

- [ ] 2.1 Replace each check row's `onclick` cell with a `<button aria-expanded aria-controls>` that toggles the details; verify Tab reaches every check, Enter and Space toggle it, and `aria-expanded` changes.
- [ ] 2.2 Remove `outline: none` from focus styles and add a visible focus indicator with at least 3:1 contrast; verify every interactive element shows it when focused by keyboard.
- [ ] 2.3 Show the sort control at every width (drop `d-none d-md-inline-block`); verify it is visible and usable at 320 CSS pixels.
- [ ] 2.4 Change the risk badge colors or text color so each badge reaches 4.5:1 and keeps its text label; verify with a contrast calculator and record the four ratios in the pull request.

## 3. Safe rendering and result states

- [ ] 3.1 Percent-encode each `uri` segment and `commit` with `encodeURIComponent` when building the API URL; verify `?uri=github.com/org/repo%3Fx%3D1` requests `.../repo%3Fx%3D1`.
- [ ] 3.2 Remove `Handlebars.SafeString` from the error path and render error codes and messages through escaped template fields; verify a mocked error message containing `<b>x</b>` appears as text.
- [ ] 3.3 Replace the catch-all error with distinct messages and next steps for not found, invalid repository, 429, 5xx, network failure or timeout, and a non-JSON body, announced through an `aria-live` region; verify each state by pointing the viewer at a mock or by blocking the request in the browser's developer tools.
- [ ] 3.4 Show `-1` scores as "Unavailable" and not-applicable checks as "Not applicable", each with the check's reason; verify on a repository with an inconclusive check.
- [ ] 3.5 Keep `apiURL`, expanded checks, and focus across `onSortChange`; verify by expanding a check, re-sorting, and confirming it stays expanded and the API link remains.
- [ ] 3.6 Respect `prefers-reduced-motion` for every animation and scroll; verify with reduced motion turned on in the browser's rendering emulation.

## 4. Structure and metadata

- [ ] 4.1 Replace the report's layout table with headings and a list (or a data table with header cells); verify with a screen reader's headings list that each check is a navigable item.
- [ ] 4.2 Move check names, risk levels, and docs links into one list, fix the duplicate `Pinned-Dependencies` entry, and add the missing checks (including SBOM); add a `build-website.yml` step that fails on a duplicate or on a name mismatch with `checks.yaml` at the pinned `ossf/scorecard` tag; verify it fails on a branch with a duplicate entry.
- [ ] 4.3 Set `document.title` and the meta description to the repository when a report loads; verify the title for `github.com/ossf/scorecard`.

## 5. Filters and default order

- [ ] 5.1 Add failing, risk-level, and unavailable filters stored in URL parameters, and default to highest-risk failing checks first; verify a filtered URL opened in a new tab shows the same filter and that clearing restores all checks.

## 6. Indexing, assets, and export

- [ ] 6.1 Record the maintainers' indexing decision for viewer reports in `scorecards-site/static/viewer/README.md` and apply it to the robots meta tag, canonical link, and sitemap; verify the three agree on the deploy preview.
- [ ] 6.2 Move inline styles, scripts, and templates into content-hashed files under `scorecards-site/static/viewer/` and remove inline event attributes; verify the HTML has no inline `<script>` body or `on*` attribute and the page still renders on the deploy preview.
- [ ] 6.3 Add a print stylesheet and a download that includes repository, commit, scan date, Scorecard version, and every check's score and reason; verify the printed page and the downloaded file both contain these and read correctly in grayscale.

## Workflow follow-up

- Rebase `staging-00` onto `main` after each task group lands, and after 6.2 move its new files to `site/static/viewer/`.
- Archive this change after its last task group ships, syncing `site-viewer` into `openspec/specs/`.
