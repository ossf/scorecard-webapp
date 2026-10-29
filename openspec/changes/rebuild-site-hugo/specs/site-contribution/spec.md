## Purpose

Defines how people change the site: links on every page that start an edit or a problem report, build-time checks that keep content within the content rules, and documentation that explains the path, including DCO sign-off.

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

### Requirement: The build rejects markup and syntax outside the content rules
The site build SHALL fail when a page under `site/content/` contains inline HTML or directive syntax outside code blocks and code spans, or shortcode syntax anywhere, so both CI and the Netlify preview fail. Pages that follow the rules SHALL pass.

#### Scenario: Inline HTML
- **WHEN** a change adds `<div>` to a page's prose
- **THEN** the build fails and names the file

#### Scenario: Shortcode
- **WHEN** a change adds `{{< callout >}}` to a page outside a code block
- **THEN** the build fails and names the file and the shortcode

#### Scenario: Clean page
- **WHEN** a change touches only plain Markdown prose, including HTML or directive syntax shown inside a code block
- **THEN** the build passes

### Requirement: The build rejects bad front matter and data
The site build SHALL fail when a page's front matter does not parse, lacks a `title` or `description`, or has a non-text `description`, or when a data-file entry lacks a required field or has a testimonial link not starting with `https://`. The error SHALL name the file and the field.

#### Scenario: Missing title
- **WHEN** a change removes `title` from a page's front matter
- **THEN** the build fails with an error naming the page and `title`

#### Scenario: Broken YAML
- **WHEN** a change leaves a page's front matter unparseable
- **THEN** the build fails and names the page

#### Scenario: Testimonial missing its link
- **WHEN** a change adds a testimonial without a `link`
- **THEN** the build fails with an error naming the testimonials file, the entry, and `link`

### Requirement: Content checks are tested both ways
Each content and data check SHALL have at least one passing and one failing test case that run in CI, so a check that silently stops failing is caught.

#### Scenario: Check stops detecting errors
- **WHEN** a change causes a check to accept its known-bad input
- **THEN** CI fails on that check's failing test case

### Requirement: The build rejects broken internal links
The site build SHALL fail when a page links to a site page or an anchor on a site page that does not exist.

#### Scenario: Link to a removed page
- **WHEN** a change deletes a page that another page links to
- **THEN** the build fails and names the page with the broken link

### Requirement: The new site is built on every push and pull request
Every push and pull request SHALL build `site/` in CI with the pinned Hugo version, and the CI job SHALL fail if the build fails.

#### Scenario: Pull request touching site/
- **WHEN** a commit that changes a file under `site/` is pushed
- **THEN** CI builds `site/` and reports the result as a check

### Requirement: Contributor guide
The repository SHALL have a `CONTRIBUTING.md` that explains how to change a page through GitHub's web editor, how DCO sign-off works for web edits, and how to fix a failed DCO check. What it says about web-editor sign-off for people without write access SHALL reflect a recorded test.

#### Scenario: First-time contributor
- **WHEN** someone without write access follows `CONTRIBUTING.md` to fix a typo through "Edit this page"
- **THEN** they can open a pull request that passes the required DCO check, or follow its steps to fix the check if it fails

#### Scenario: Grounded in a test
- **WHEN** a reviewer asks what `CONTRIBUTING.md` claims about web-editor sign-off
- **THEN** the claim matches the recorded result of a web edit made by someone without write access
