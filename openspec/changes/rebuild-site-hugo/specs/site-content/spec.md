## Purpose

Defines what the site's content is made of — plain Markdown pages, YAML data files, a page generated from upstream check data, and mirrored upstream guides — and how it is organized, so non-code contributors can edit it safely.

## ADDED Requirements

### Requirement: Pages are plain Markdown
Content pages SHALL be GitHub-flavored Markdown. Outside code blocks and code spans, pages SHALL NOT contain inline HTML, MDX, shortcodes, or directive syntax. Layout, spacing, and color SHALL come from templates, not content.

#### Scenario: Ordinary page
- **WHEN** a contributor opens any page under `site/content/` in GitHub's web editor
- **THEN** it contains only Markdown and front matter, and GitHub's preview shows its text, headings, lists, links, tables, and code blocks

#### Scenario: Placeholder text in prose
- **WHEN** a page's prose contains a placeholder such as `<github_or_gitlab>`
- **THEN** it is written so that it renders as visible text, not as an HTML tag

### Requirement: Front matter has a title and description
Each content page SHALL have front matter with a text `title` and a text `description`, and no other field is required of contributors.

#### Scenario: Complete front matter
- **WHEN** a page has a text `title` and a text `description`
- **THEN** the page builds, and its title and description appear in the page heading, the browser title, and the meta description

### Requirement: Callouts use GitHub's alert syntax
Callouts SHALL be written in GitHub's alert syntax (`> [!NOTE]`, `> [!TIP]`, `> [!IMPORTANT]`, `> [!WARNING]`, `> [!CAUTION]`) and SHALL render on the site as styled callouts.

#### Scenario: Note callout
- **WHEN** a page contains `> [!NOTE]` followed by quoted text
- **THEN** the site shows a styled note callout without the literal `[!NOTE]`, and GitHub's preview shows a note alert

### Requirement: Structured content lives in data files
The home page hero (headline, tagline, image, image description, buttons) and the testimonials (quote, name, project, link) SHALL be stored in YAML data files and rendered by templates, so they can be changed without editing layout.

#### Scenario: Change the hero headline
- **WHEN** a contributor edits the headline in the home data file
- **THEN** the home page shows the new headline, and no template or page file changed

#### Scenario: Add a testimonial
- **WHEN** a contributor adds an entry with a quote, name, project, and `https://` link to the testimonials file
- **THEN** the home page shows the new testimonial

### Requirement: Site map
The site SHALL have these sections, each page being one Markdown file: Home; Get started (GitHub Action, command line, authentication and advanced setup, viewer, REST API, badges); Docs (what Scorecard is, how it works, checks, understanding results, FAQ, beginner's guide); Use cases (two pages); Community (project name, part of the open source community, get involved).

#### Scenario: Top navigation
- **WHEN** a visitor opens any page
- **THEN** the top navigation links to Get started, Docs, Use cases, Community, and the viewer

#### Scenario: "How it works" has one home
- **WHEN** a visitor looks for how Scorecard works
- **THEN** exactly one page covers it, under Docs

### Requirement: Home page
The home page SHALL show the headline and tagline, a three-step "how it works" summary, a repository lookup box that opens `/viewer/?uri=<input>`, one path for maintainers ("protect your repository") and one for consumers ("evaluate a dependency"), and the testimonials. Its hero SHALL use a static image, not video.

#### Scenario: Look up a repository
- **WHEN** a visitor types `github.com/ossf/scorecard` into the lookup box and submits
- **THEN** the browser opens `/viewer/?uri=github.com/ossf/scorecard`

#### Scenario: Paths by reader
- **WHEN** a visitor reads the home page
- **THEN** it links to the maintainer page and the consumer page under Use cases

### Requirement: Two use-case pages
Use cases SHALL have one page for maintainers who want to improve their own project, and one for consumers evaluating dependencies. The consumer page SHALL also cover organizations, security teams, and registries or platforms that use results at scale.

#### Scenario: Registry operator
- **WHEN** someone who runs a package registry opens Use cases
- **THEN** the consumer page explains how to get results at scale from the REST API and how to link to the viewer

### Requirement: Checks page is generated from upstream data
The Checks page SHALL be generated at build time from an unmodified copy of upstream `docs/checks/internal/checks.yaml`. It SHALL list every check sorted by risk (Critical, High, Medium, Low), then by name, with its risk, summary, description, and remediation. A check with an unknown risk SHALL fail the build.

#### Scenario: Every upstream check appears
- **WHEN** the synced `checks.yaml` defines 20 checks
- **THEN** the Checks page lists all 20 with each check's risk and remediation, and its count separates checks that run by default from experimental ones

#### Scenario: Sort order
- **WHEN** the Checks page is rendered
- **THEN** all Critical checks come before High, High before Medium, and Medium before Low, alphabetical within each level

#### Scenario: Unknown risk
- **WHEN** a check in the synced file has a risk other than Critical, High, Medium, or Low
- **THEN** the build fails and names the check and its risk

### Requirement: Experimental checks are labeled
Checks that upstream excludes from default runs SHALL be marked as experimental on the Checks page, with how to enable them. Until `checks.yaml` carries that information, the site SHALL keep a short override list of experimental checks.

#### Scenario: SBOM and Webhooks
- **WHEN** the Checks page lists SBOM and Webhooks
- **THEN** each is marked experimental and says how to enable it

### Requirement: Mirrored upstream guides
The FAQ and the beginner's guide SHALL be copies of upstream `docs/faq.md` and `docs/beginner-checks.md`. Each SHALL say it is copied from upstream and link to the original. Relative links SHALL point to the files on GitHub, and placeholders in prose SHALL render as text.

#### Scenario: FAQ page
- **WHEN** a visitor opens the FAQ
- **THEN** a notice says the page is copied from `ossf/scorecard`'s `docs/faq.md`, and `<github_or_gitlab>` in its text renders as visible text

#### Scenario: Relative link in a guide
- **WHEN** a mirrored guide links to another upstream document with a relative path
- **THEN** the link opens that document on GitHub

### Requirement: Getting-started content matches upstream
Get started pages SHALL match current upstream instructions. The GitHub Action guide SHALL use the default `GITHUB_TOKEN`, with a personal access token only for classic branch protection, and list the supported triggers. Docker instructions SHALL use `ghcr.io/ossf/scorecard:latest`.

#### Scenario: Action setup
- **WHEN** a maintainer follows the GitHub Action guide
- **THEN** it does not require a personal access token, and it lists `push` and `schedule` as the supported triggers

#### Scenario: Docker image
- **WHEN** a reader copies the Docker command
- **THEN** it uses `ghcr.io/ossf/scorecard:latest`

### Requirement: Claims about Scorecard are accurate
Site copy SHALL describe Scorecard as assessing security practices. It SHALL NOT claim Scorecard scans code for vulnerabilities or protects code. Usage figures SHALL say how many repositories Scorecard scans, not how many use it. The site SHALL NOT cite the 2021 Synopsys statistic.

#### Scenario: Headline and value statement
- **WHEN** a visitor reads the home page
- **THEN** it describes checking security practices, and any scale figure says Scorecard scans more than 1.3 million repositories a week

### Requirement: English at the root URLs
English content SHALL be served at the root URLs (for example `/get-started/`), and the site SHALL be English-only until translations are added.

#### Scenario: URL layout
- **WHEN** a visitor opens `/docs/checks/`
- **THEN** the English Checks page is served, with no language prefix in the URL
