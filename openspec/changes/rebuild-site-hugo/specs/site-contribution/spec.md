## Purpose

Defines how people change the site: links on every page that start an edit or a problem report, CI checks that keep content within the content rules, and documentation that explains the path, including DCO sign-off.

## ADDED Requirements

### Requirement: Every page links to its source for editing
Every page SHALL show an edit link. For a page written on the site, it SHALL open that page's file under `site/content/` on `main` in GitHub's web editor. For the generated Checks page and the mirrored guides, it SHALL be labeled "Edit upstream" and open the upstream source file.

#### Scenario: Site page
- **WHEN** a visitor selects "Edit this page" on the command-line guide
- **THEN** GitHub's web editor opens that guide's Markdown file in `ossf/scorecard-webapp` on `main`

#### Scenario: Generated page
- **WHEN** a visitor selects "Edit upstream" on the Checks page
- **THEN** GitHub's web editor opens `docs/checks/internal/checks.yaml` in `ossf/scorecard`

#### Scenario: Mirrored guide
- **WHEN** a visitor selects "Edit upstream" on the FAQ
- **THEN** GitHub's web editor opens `docs/faq.md` in `ossf/scorecard`

### Requirement: Every page links to a problem report
Every page SHALL show a "Report a problem with this page" link that opens a new issue in `ossf/scorecard-webapp`, with the page's title and URL filled in.

#### Scenario: Reporting a problem
- **WHEN** a visitor selects "Report a problem with this page" on the FAQ
- **THEN** a new-issue form opens with a title naming the FAQ and a body containing the page URL

### Requirement: CI rejects markup and syntax outside the content rules
A pull request SHALL fail CI when a page under `site/content/` contains inline HTML, shortcodes, or directive syntax outside code blocks and code spans. Pages that follow the rules SHALL pass.

#### Scenario: Inline HTML
- **WHEN** a pull request adds `<div>` to a page
- **THEN** CI fails and names the file and line

#### Scenario: Shortcode
- **WHEN** a pull request adds `{{< callout >}}` to a page outside a code block
- **THEN** CI fails and names the file and line

#### Scenario: Clean page
- **WHEN** a pull request changes only plain Markdown prose
- **THEN** the content checks pass

### Requirement: CI rejects bad front matter and data
A pull request SHALL fail CI when a page's front matter does not parse, lacks a `title` or `description`, or has a non-text `description`, or when a data-file entry lacks a required field or has a testimonial link not starting with `https://`. The error SHALL name the file and the field.

#### Scenario: Missing title
- **WHEN** a pull request removes `title` from a page's front matter
- **THEN** CI fails with an error naming the page and `title`

#### Scenario: Broken YAML
- **WHEN** a pull request leaves a page's front matter unparseable
- **THEN** CI fails and names the page

#### Scenario: Testimonial missing its link
- **WHEN** a pull request adds a testimonial without a `link`
- **THEN** CI fails with an error naming the testimonials file, the entry, and `link`

### Requirement: Content checks are tested both ways
Each content and data check SHALL have at least one passing and one failing test case that run in CI, so a check that silently stops failing is caught.

#### Scenario: Check stops detecting errors
- **WHEN** a change causes a check to accept its known-bad input
- **THEN** CI fails on that check's failing test case

### Requirement: CI rejects broken links
A pull request SHALL fail CI when the built site contains an internal link to a page or anchor that does not exist.

#### Scenario: Link to a removed page
- **WHEN** a pull request deletes a page that another page links to
- **THEN** CI fails and names the page with the broken link

### Requirement: The new site is built on every pull request
Every pull request and push SHALL build `site/` in CI with the pinned generator, and fail if the build fails.

#### Scenario: Pull request touching site/
- **WHEN** a pull request changes a file under `site/`
- **THEN** CI builds `site/` and reports the result as a check on the pull request

### Requirement: Contributor guide
The repository SHALL have a `CONTRIBUTING.md` that explains how to change a page through GitHub's web editor, how DCO sign-off works for web edits, and how to fix a failed DCO check. What it says about web-editor sign-off for people without write access SHALL reflect a recorded test.

#### Scenario: First-time contributor
- **WHEN** someone without write access follows `CONTRIBUTING.md` to fix a typo through "Edit this page"
- **THEN** they can open a pull request that passes the required DCO check, or follow its steps to fix the check if it fails

#### Scenario: Grounded in a test
- **WHEN** a reviewer asks what `CONTRIBUTING.md` claims about web-editor sign-off
- **THEN** the claim matches the recorded result of a web edit made by someone without write access
