---
title: Viewer
description: Look up results for projects in the weekly scan.
weight: 4
---

The [Scorecard viewer](/viewer/) shows results for projects that Scorecard
scans every week. Link to any project with:

```text
https://scorecard.dev/viewer/?uri=<github_or_gitlab>.com/<owner>/<repository>
```

For example:

- [github.com/ossf/scorecard](/viewer/?uri=github.com/ossf/scorecard)
- [gitlab.com/fdroid/fdroidclient](/viewer/?uri=gitlab.com/fdroid/fdroidclient)

For projects that aren't in the weekly scan, use the
[command line](/get-started/cli/). To add a project to the scan, follow
[`ossf/scorecard-infra`'s contributing guide](https://github.com/ossf/scorecard-infra/blob/main/CONTRIBUTING.md).
