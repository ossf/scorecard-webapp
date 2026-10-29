---
title: Understanding results
description: How to read check scores and the aggregate score.
weight: 4
---

## Start with the checks, not the aggregate

The aggregate score tells you nothing about which practices a repository
follows. Different combinations of checks can produce the same score, and
scores change as checks are added or refined. Read the individual
[checks](/docs/checks/) that matter to you.

## Checks are heuristics

Every check is a heuristic, so there are false positives and false
negatives. Scorecard isn't meant to be a requirement that every project must
meet. Use the results to decide what to improve, or which risks to accept.

## See why a check scored the way it did

Run the [command line](/get-started/cli/) with `--show-details`, or open a
project in the [viewer](/viewer/) and expand a check.

## Missing checks in weekly-scan results

Results from the weekly scan leave out `CI-Tests`, `Contributors`, and
`Dependency-Update-Tool`. Run the command line to get all checks.
