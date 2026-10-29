---
title: REST API
description: Query pre-calculated results for open source projects.
weight: 5
---

The [REST API](https://api.scorecard.dev) returns results from the weekly
scan and from projects that publish their results through the GitHub Action.

```sh
curl https://api.scorecard.dev/projects/github.com/ossf/scorecard
```

> [!NOTE]
> Results from the weekly scan leave out the `CI-Tests`, `Contributors`, and
> `Dependency-Update-Tool` checks because they're costly to run at scale.

To make your project available, set `publish_results: true` in the
[GitHub Action](/get-started/github-action/). API data is licensed under the
[CDLA Permissive 2.0](https://cdla.dev/permissive-2-0).
