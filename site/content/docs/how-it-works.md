---
title: How it works
description: What Scorecard checks, and how the scores are calculated.
weight: 2
---

Scorecard is an automated tool that checks a project for a set of security
practices, such as code review, pinned dependencies, and branch protection.
It assesses practices; it doesn't scan your code for new vulnerabilities.
See [What is Scorecard?](/docs/what-is-scorecard/) and the
[FAQ](/docs/faq/#what-is-the-difference-between-scorecard-and-other-code-scanning-tools).

## Where results come from

- **The GitHub Action** runs on a repository you own, on pushes to the
  default branch and on a schedule.
- **The command line** runs on any repository, whenever you choose.
- **The weekly scan** covers more than 1.3 million repositories a week. Its results
  are in the [viewer](/viewer/) and the [REST API](/get-started/rest-api/).

## Scores

Each [check](/docs/checks/) returns a score from 0 to 10, where 10 is best.
The aggregate score is an average of the checks, weighted by risk:

| Risk | Weight |
| --- | --- |
| Critical | 10 |
| High | 7.5 |
| Medium | 5 |
| Low | 2.5 |
