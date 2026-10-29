---
title: Badges
description: Show your Scorecard results in your README.
weight: 6
---

With `publish_results: true` in the [GitHub Action](/get-started/github-action/),
you can show a badge that updates each time the Action runs. Add this to your README,
replacing `{owner}` and `{repo}`:

```md
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/{owner}/{repo}/badge)](https://scorecard.dev/viewer/?uri=github.com/{owner}/{repo})
```
