## Purpose

Defines how the results viewer at `/viewer/` loads, renders, and presents a repository's Scorecard results, and the dependency, accessibility, and safety guarantees it must keep.

## ADDED Requirements

### Requirement: Badge links keep working
`/viewer/?uri=<platform>/<org>/<repo>` SHALL load and show that repository's results, with an optional `commit` parameter, because every Scorecard badge links there.

#### Scenario: Badge link
- **WHEN** a visitor opens `/viewer/?uri=github.com/ossf/scorecard`
- **THEN** the viewer shows that repository's latest results

### Requirement: Third-party scripts are pinned and integrity-checked
Every script the viewer loads from another origin SHALL name an exact version and carry a Subresource Integrity hash. The viewer SHALL NOT load `@latest` or any other moving tag.

#### Scenario: Pinned Handlebars
- **WHEN** the viewer page is fetched
- **THEN** its Handlebars `<script>` names version 4.7.10 and has an `integrity` attribute that matches the file served at that URL

#### Scenario: No moving tags
- **WHEN** the viewer's files are searched for `@latest` or a script URL without a version
- **THEN** none is found

### Requirement: The report works with a keyboard and a screen reader
Every control in the report SHALL be reachable and operable with a keyboard, SHALL show a visible focus indicator, and SHALL expose its role and state to assistive technology. Each check's expand control SHALL be a button with `aria-expanded`. Sorting SHALL be available at every viewport width. The report SHALL use headings and lists or a data table with headers, not a layout table.

#### Scenario: Expand a check by keyboard
- **WHEN** a keyboard user tabs to a check and presses Enter or Space
- **THEN** the check's details expand, `aria-expanded` becomes `true`, and focus stays on the control

#### Scenario: Visible focus
- **WHEN** any interactive element in the viewer receives keyboard focus
- **THEN** a focus indicator with at least 3:1 contrast against its background is visible

#### Scenario: Sort on a small screen
- **WHEN** the viewer is shown at 320 CSS pixels wide
- **THEN** the sort control is visible and usable

### Requirement: Text meets contrast minimums
Text in the viewer, including risk badges, SHALL have at least 4.5:1 contrast with its background, and information conveyed by color SHALL also be conveyed by text.

#### Scenario: Risk badges
- **WHEN** the low, medium, high, and critical badges are measured
- **THEN** each has at least 4.5:1 contrast and shows its risk level as text

### Requirement: Motion respects user preference
Animations and smooth scrolling SHALL be disabled when the visitor's system requests reduced motion.

#### Scenario: Reduced motion
- **WHEN** a visitor with `prefers-reduced-motion: reduce` loads a report and expands a check
- **THEN** no animation or smooth scroll runs

### Requirement: API data is escaped and API paths are encoded
Every value the viewer renders from an API response, including error codes and messages, SHALL be HTML-escaped. Each segment of the repository path and the `commit` value SHALL be percent-encoded before the API URL is built.

#### Scenario: Markup in an error message
- **WHEN** the API returns an error whose message contains `<b>x</b>`
- **THEN** the viewer shows the literal text `<b>x</b>` and adds no element to the page

#### Scenario: Reserved characters in the repository
- **WHEN** a visitor opens `/viewer/?uri=github.com/org/repo%3Fx%3D1`
- **THEN** the request goes to a URL where `repo?x=1` is one encoded path segment

### Requirement: Each result state is distinct and actionable
The viewer SHALL show a distinct plain-language message, with a next step, for each of: loading, report not found, invalid repository, rate limited (429), server error (5xx), network failure or timeout, and a malformed response. Status changes SHALL be announced to screen readers. A score of `-1` SHALL be shown as "Unavailable" and a not-applicable check as "Not applicable", each with its reason.

#### Scenario: Report not found
- **WHEN** the API returns 404 for a repository
- **THEN** the viewer says no report exists and links to how to run Scorecard on it

#### Scenario: Malformed response
- **WHEN** the API returns a non-JSON body with a 400 status
- **THEN** the viewer shows the invalid-request message, not "Status code: undefined"

#### Scenario: Unavailable score
- **WHEN** a check's score is `-1`
- **THEN** the viewer shows "Unavailable" and the check's reason instead of `?`

### Requirement: Re-sorting keeps the report's state
Re-sorting or filtering the report SHALL keep the API link, the set of expanded checks, and focus on the control that was used.

#### Scenario: Sort after expanding a check
- **WHEN** a visitor expands a check and then changes the sort order
- **THEN** that check is still expanded, the API link is still shown, and focus is on the sort control

### Requirement: Check metadata has one source
The viewer's check names, risk levels, and documentation links SHALL come from one list. A CI check SHALL fail when that list has a duplicate entry or differs in check names from `checks.yaml` in the `ossf/scorecard` release the site pins.

#### Scenario: Duplicate entry
- **WHEN** a check appears twice in the list
- **THEN** the CI check fails and names the duplicate

#### Scenario: Missing check
- **WHEN** the pinned `checks.yaml` defines a check the list lacks
- **THEN** the CI check fails and names the missing check

### Requirement: The page describes the repository shown
When a report loads, the page title and description SHALL name the repository.

#### Scenario: Title
- **WHEN** the report for `github.com/ossf/scorecard` loads
- **THEN** the page title contains `github.com/ossf/scorecard`

### Requirement: Results can be filtered and default to what needs action
The viewer SHALL offer filters for failing checks, risk level, and unavailable checks, reflected in the URL so a filtered view can be shared. The default order SHALL put the highest-risk failing checks first.

#### Scenario: Shared filter
- **WHEN** a visitor applies the failing-checks filter and opens the resulting URL in a new tab
- **THEN** the same filter is applied

### Requirement: Indexing follows one documented decision
Whether search engines index viewer reports SHALL be decided once and documented, and the viewer's robots directives, canonical URL, and sitemap entries SHALL all follow that decision.

#### Scenario: Consistent signals
- **WHEN** a report page's robots meta tag, canonical link, and the sitemap are compared
- **THEN** they agree with the documented decision

### Requirement: Viewer assets are versioned files
The viewer's styles, scripts, and templates SHALL be separate files under `viewer/` whose names change when their content changes, so they can be cached long-term and a Content Security Policy needs no inline-script allowance.

#### Scenario: No inline script
- **WHEN** the viewer's HTML is fetched
- **THEN** it contains no inline `<script>` body or `on*` event attribute

### Requirement: Reports can be printed and downloaded
The viewer SHALL provide a print layout and a downloadable report. Both SHALL include the repository, commit, scan date, Scorecard version, and every check's score and reason, and SHALL NOT rely on color alone.

#### Scenario: Download
- **WHEN** a visitor downloads the report for a repository
- **THEN** the file names the repository, commit, scan date, and Scorecard version and lists every check with its score and reason
