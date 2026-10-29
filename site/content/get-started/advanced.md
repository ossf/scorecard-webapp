---
title: Authentication and advanced setup
description: GitHub Apps, multiple tokens, and verifying releases.
weight: 3
---

## Use several tokens

Separate tokens with commas and Scorecard uses them in turn:

```sh
export GITHUB_AUTH_TOKEN=<token1>,<token2>
```

## Use a GitHub App

For higher rate limits, use a
[GitHub App installation](https://docs.github.com/en/developers/apps/building-github-apps/creating-a-github-app)
and set:

```sh
export GITHUB_APP_KEY_PATH=<path to the key file>
export GITHUB_APP_INSTALLATION_ID=<installation id>
export GITHUB_APP_ID=<app id>
```

## Verify a release

Releases carry [SLSA](https://slsa.dev) provenance. Install
[`slsa-verifier`](https://github.com/slsa-framework/slsa-verifier#installation),
download `multiple.intoto.jsonl` from the same release, and run:

```sh
slsa-verifier verify-artifact --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/ossf/scorecard --source-tag <the-tag> <the-archive>
```
