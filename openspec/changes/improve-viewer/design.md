## Context

The viewer is one 1,666-line static file, `scorecards-site/static/viewer/index.html`, that calls `api.scorecard.dev` from the browser and renders with Handlebars from a CDN. It has no build step and no tests. On `staging-00`, `rebuild-site-hugo` has already moved it to `site/static/viewer/` as a pure rename, and requires builds that need only Hugo and Git.

## Goals / Non-Goals

**Goals:** fix the four launch blockers on production first, then bring the rest of the viewer to the same standard as the rebuilt site, one shippable task group at a time.

**Non-Goals:** changing the REST API, adding a build step or framework to the viewer, translations.

## Decisions

- **Ship on `main`, not `staging-00`.** The defects are live on scorecard.dev. Fixing them on `main` reaches production now; the rebuild picks them up when `staging-00` is rebased onto `main`.
- **Pin Handlebars on the CDN with SRI first; vendor it when assets split (6.2).** Pinning plus `integrity` is a one-line fix that removes the moving dependency today. Serving a local copy fits the later Content Security Policy work and avoids the third-party origin, but is a larger change.
- **Keep the viewer build-free.** Hashed asset names (6.2) are produced by hand or by a script committed with the change, not a bundler, so the rebuild's Hugo-and-Git build rule still holds after the move.
- **CI checks run in `build-website.yml` with shell tools.** The repository has no Go or Node tooling that survives cutover, so the `@latest` check (1.2) and the `checks.yaml` comparison (4.2) use `grep` and `yq`, pinned in the workflow.

## Risks / Trade-offs

- Every badge links to the viewer. A regression breaks badge links across the ecosystem, so each task group checks `/viewer/?uri=github.com/ossf/scorecard` on its deploy preview before merge.
- The rebase carries edits across `staging-00`'s rename, but new files that 6.2 adds under `scorecards-site/static/viewer/` stay at the old path. Rebase `staging-00` after 6.2 lands and move those files to `site/static/viewer/` in the same pass.

## Open Questions

- Should viewer reports be indexed by search engines (6.1)? This is for the maintainers and the website owner; the task records whatever they decide.
