---
title: Command line
description: Run Scorecard on any repository, including ones you don't own.
weight: 2
---

Use the command-line tool to scan projects you don't own, or projects that
aren't in the weekly scan.

## Install

| Method | Command |
| --- | --- |
| Homebrew (macOS or Linux) | `brew install scorecard` |
| Docker | `docker pull ghcr.io/ossf/scorecard:latest` |
| Nix | `nix-shell -p nixpkgs.scorecard` |
| Release binary | Download from the [latest release](https://github.com/ossf/scorecard/releases/latest) |

Scorecard supports macOS and Linux. Windows may have issues.

## Authenticate

GitHub rate-limits unauthenticated requests. Create a
[classic personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic)
with the `public_repo` scope and set it in your environment:

```sh
export GITHUB_AUTH_TOKEN=<your access token>
```

## Run a scan

```sh
scorecard --repo=github.com/ossf/scorecard
```

With Docker:

```sh
docker run -e GITHUB_AUTH_TOKEN=token ghcr.io/ossf/scorecard:latest --show-details --repo=https://github.com/ossf/scorecard
```

Useful options:

- `--show-details` explains why each check scored the way it did.
- `--checks=CI-Tests,Code-Review` runs only the checks you list.
- `--format=json` prints machine-readable results.
- `--npm=angular` (also `--pypi`, `--rubygems`, `--nuget`) finds the source
  repository for a package. It can't be combined with `--repo`.

## GitLab

Create a [GitLab access token](https://gitlab.com/-/profile/personal_access_tokens)
with `read_api`, `read_user`, and `read_repository`, then:

```sh
export GITLAB_AUTH_TOKEN=glpat-xxxx
scorecard --repo gitlab.com/<org>/<project>
```

For self-hosted GitLab at a path such as `foo.com/bar/`, also set
`GL_HOST=foo.com/bar`.

## GitHub Enterprise Server

```sh
export GH_HOST=github.corp.com
export GITHUB_AUTH_TOKEN=token
scorecard --repo=github.corp.com/org/repo
```
