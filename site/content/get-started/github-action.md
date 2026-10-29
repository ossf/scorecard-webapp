---
title: GitHub Action
description: Run Scorecard automatically on repositories you own.
weight: 1
---

The easiest way to use Scorecard on GitHub projects you own is the
[Scorecard GitHub Action](https://github.com/ossf/scorecard-action). Alerts
appear in the repository's **Security** tab.

## Triggers

The Action supports the `push` and `schedule` triggers, on the default
branch only. The `pull_request` and `workflow_dispatch` triggers are
experimental.

## Install

1. On your repository's main page, select **Security**, then **Code scanning**.
2. Select **Add tool**, or **Configure scanning tool** and then
   **Explore workflows**.
3. Search for "OSSF", choose **OSSF Scorecard**, and select **Configure**.
4. Commit the workflow.

To add the workflow to many repositories at once, use the
[installer](https://github.com/ossf/scorecard-action/blob/main/docs/installer.md).

## Authentication

The workflow's default `GITHUB_TOKEN` is enough, and it's the recommended
setup. You only need a fine-grained personal access token if you use
GitHub's classic branch protection and want it reflected in your results.
See [Authentication and advanced setup](/get-started/advanced/).

## Publish your results

Set `publish_results: true` in the workflow to make your results available
through the [REST API](/get-started/rest-api/) and to display a
[badge](/get-started/badges/).
