---
title: What is Scorecard?
description: The problem Scorecard addresses, and what it does and doesn't do.
weight: 1
---

Open source projects depend on other open source projects, and each
dependency brings the risks of how it is built and maintained. A project
without code review, with unpinned dependencies, or with workflows that run
untrusted code is easier to compromise, and so is everything that depends
on it.

Scorecard is an automated tool that assesses a project's security practices
through a set of [checks](/docs/checks/), and gives each check a score from
0 to 10. It helps:

- **Maintainers** find the practices to improve in their own projects.
- **Consumers** judge the risks a dependency brings, and decide whether to
  accept them, choose an alternative, or work with its maintainers.

Scorecard scans more than 1.3 million repositories a week, and publishes the
results in the [viewer](/viewer/) and the [REST API](/get-started/rest-api/).

## What Scorecard doesn't do

- **It doesn't scan your code for new vulnerabilities.** It assesses
  practices. The Vulnerabilities check reports known vulnerabilities in the
  project and its dependencies from the [OSV](https://osv.dev) database.
- **It isn't a requirement every project must meet.** The checks are
  heuristics, with false positives and false negatives, and an aggregate
  score says nothing about which practices a project follows. See
  [Understanding results](/docs/understanding-results/).
