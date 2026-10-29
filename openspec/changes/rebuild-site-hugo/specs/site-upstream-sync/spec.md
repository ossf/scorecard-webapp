## Purpose

Defines how upstream sources — check definitions, user guides, and the viewer — reach the site as committed, reviewable copies, so the site stays accurate without depending on live fetches at build time.

## ADDED Requirements

### Requirement: Scheduled sync opens a reviewable pull request
A workflow SHALL run weekly and on manual trigger. It SHALL copy `docs/checks/internal/checks.yaml`, `docs/faq.md`, and `docs/beginner-checks.md` from `ossf/scorecard` at `main`, and the viewer from `scorecards-site/static/viewer/`. When the copies differ from what is committed, it SHALL open or update a single pull request.

#### Scenario: Upstream changed a check
- **WHEN** the weekly run finds that upstream `checks.yaml` differs from the committed copy
- **THEN** the workflow opens a pull request with the updated copy, or updates its existing open sync pull request

#### Scenario: Nothing changed
- **WHEN** the run finds every copy identical to what is committed
- **THEN** no pull request is opened or updated

### Requirement: Builds never fetch upstream content
Site builds SHALL use only the committed copies of upstream content. Changes from upstream SHALL reach production only by merging a sync pull request.

#### Scenario: Upstream changes between syncs
- **WHEN** upstream edits `checks.yaml` and the site is rebuilt before the next sync
- **THEN** the rebuilt site still shows the committed copy

### Requirement: Check data is copied unchanged
The synced `checks.yaml` SHALL be a byte-for-byte copy of the upstream file, and the Checks page SHALL be generated from it at build time.

#### Scenario: Comparing copies
- **WHEN** the committed copy is compared with upstream `checks.yaml` at the commit the sync used
- **THEN** they are identical

### Requirement: A missing source fails the sync
If any source file is missing or cannot be fetched, the workflow SHALL fail, and it SHALL NOT open a pull request that deletes or empties a committed copy.

#### Scenario: Upstream file renamed
- **WHEN** upstream renames `docs/beginner-checks.md`
- **THEN** the workflow run fails with an error naming the missing file, and no pull request is opened

### Requirement: Sync pull requests prompt the experimental-check review
Each sync pull request SHALL remind reviewers to check whether the site's list of experimental checks still matches upstream.

#### Scenario: New check upstream
- **WHEN** a sync pull request adds a check to `checks.yaml`
- **THEN** its description asks reviewers to confirm whether the new check is experimental

### Requirement: The sync workflow follows least privilege
The workflow SHALL pin every action to a full commit SHA, grant only the permissions it needs to push a branch and open a pull request, not persist checkout credentials, and pass `zizmor` with no findings.

#### Scenario: Workflow audit
- **WHEN** `zizmor` runs against the sync workflow
- **THEN** it reports no findings
