---
title: Home
date: 2021-07-12T15:33:03.264Z
description: Quickly assess open source projects for risky practices
slug: home
thumbnail: /assets/checks.png
---

<div>
  <sidebar ref="sideBar" class="sticky top-60 h-400 min-w-max w-1/3 hidden md:block mr-80"></sidebar>
</div>

<section class="bg-orange prose md:prose-lg w-full">

<h2 class="h1" id="run-the-checks">Run the checks</h2>

OpenSSF Scorecard can be used in a couple of different ways:

1. Run automatically on code you own **using the GitHub Action**
2. Run manually on your (or somebody else’s) project **via the Command Line**

### Using the GitHub Action

<section class="highlight-section">

### Install time: <10 mins

Use the action to automatically scan any code updates for security vulnerabilities. Any time someone commits a change, the action will automatically check the repo and alert you (and other maintainers) if there are problems.

<details open><summary>See it in action</summary>

<iframe
      title="action video"
      width="100%"
      height="393"
      allow="autoplay"
      loop="true"
      src="assets/github-action.mp4?autoplay=1&controls=0&loop=1&mute=1">
  </iframe>

</details>

### Installation instructions

1. You need to own the repository you are installing the action to, or have admin rights to it.
2. [Authenticate](https://github.com/ossf/scorecard-action#authentication) your access to the repository with a Personal Access Token
3. Add Scorecard to your [codescanning suite](https://github.com/ossf/scorecard-action#workflow-setup) inside GitHub using the link below:
   <br/>

<a href="https://github.com/marketplace/actions/ossf-scorecard-action" class="btn cta-icon">Install the action</a>

</section>

### Using the CLI

<section class="highlight-section">

### Install time: <10mins

You can use Scorecard on the Command Line. This enables you to:

- Check someone else’s repository
- Select which checks you want to run
- Control how detailed your results are

<details open><summary>See it in action</summary>

<iframe
      title="CLI video"
      width="100%"
      height="477"
      allow="autoplay"
      loop="true"
      src="assets/cli.mp4?autoplay=1&controls=0&loop=1&mute=1">
  </iframe>

</details>

### Install and run

1. Create a [GitHub personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) with 'public_repo' scope. Store the token somewhere safe.
2. Choose a language-specific quick start below, or refer to our [detailed instructions](https://github.com/ossf/scorecard#scorecard-command-line-interface)

<code-group>

  <code-block title="Homebrew" active>

```bash
# For posix platforms, e.g. linux, mac:
export GITHUB_AUTH_TOKEN=<your access token>

# For windows:
set GITHUB_AUTH_TOKEN=<your access token>

brew install scorecard

scorecard --repo=<your choice of repo e.g. github.com/ossf-tests/scorecard-check-branch-protection-e2e>
```

  </code-block>

  <code-block title="Docker">

```bash
docker run -e GITHUB_AUTH_TOKEN=<your access token> gcr.io/openssf/scorecard:stable --repo=<your choice of repo e.g. github.com/ossf-tests/scorecard-check-branch-protection-e2e>
```

  </code-block>

  <code-block title="Nix">

```bash
export GITHUB_AUTH_TOKEN=<your access token>

nix-shell -p nixpkgs.scorecard

scorecard --repo=<your choice of repo e.g. github.com/ossf-tests/scorecard-check-branch-protection-e2e>
```

  </code-block>

</code-group>

Scorecard also has standalone binaries and other platforms troubleshooting and custom configuration available. Learn more here:

<br/>

<a href="https://github.com/ossf/scorecard#scorecard-command-line-interface" class="btn cta-icon">Detailed installation instructions</a>

</section>

<h2 class="h1" id="learn-more">Learn more</h2>

> We rely on Security Scorecards *[i.e., OpenSSF Scorecard]* to ensure we follow secure development best practices.

<div class="text-right"><cite>Appu Goundan, <a href="https://github.com/GoogleContainerTools/distroless">Distroless</a></cite></div>

### The problem

By some estimates\* 84% of all codebases have at least one vulnerability, with an average of 158 per codebase. The majority have been in the code for more than 2 years and have documented solutions available.

Even in large tech companies, the tedious process of reviewing code for vulnerabilities falls down the priority list, and there is little insight into known vulnerabilities and solutions that companies can draw on.

That’s where Security Scorecards *[i.e., OpenSSF Scorecard]* is helping. Its focus is to understand the security posture of a project and assess the risks that dependencies introduce.

\*[Open Source Security and Risk Analysis Report](https://www.synopsys.com/software-integrity/resources/analyst-reports/open-source-security-risk-analysis.html?intcmp=sig-blog-ossra1) (Synopsys, 2021)

<br/>

### What is OpenSSF Scorecard?

##### Scorecard assesses open source projects for security risks through a series of automated checks

It was created by OSS developers to help improve the health of critical projects that the community depends on.

You can use it to proactively assess and make informed decisions about accepting security risks within your codebase. You can also use the tool to evaluate other projects and dependencies, and work with maintainers to improve codebases you might want to integrate.

Scorecard help you enforce best practices that can guard against:

<section class="grid grid-cols-1 grid-rows-1 md:grid-cols-2 md:grid-cols-2">

<div>

![malicious maintainers](assets/malicious-maintainer.svg)

#### Malicious maintainers

</div>

<div>

![build system compromises](assets/compromised-build.svg)

#### Build system compromises

</div>

<div>

![source code compromises](assets/compromised-source.svg)

#### Source code compromises

</div>

<div>

![malicious packages](assets/malicious-package.svg)

#### Malicious packages

</div>

</section>

<br/>
<br/>

> It took less than 5 minutes to install. It quickly analysed the repo and identified easy ways to make the project more secure.

<div class="text-right"><cite>Priya Wadhwa, <a href="https://github.com/GoogleContainerTools/kaniko">Kaniko</a></cite></div>

<br/>

### How it works

Scorecard checks for vulnerabilities affecting different parts of the software supply chain including **source code**, **build**, **dependencies**, **testing**, and project **maintenance**.

Each automated check returns a **score out of 10** and a **risk level**. The risk level **adds a weighting** to the score, and this weighting is compiled into a single, **aggregate score**. This score helps give a sense of the overall security posture of a project.

Alongside the scores, the tool provides remediation prompts to help you **fix problems** and strengthen your development practices.

![scale of risk](assets/diagram-risks.svg)

<br/>

### The checks

##### The checks collect together security best practises and industry standards

The riskiness of each vulnerability is based on how easy it is to exploit. For example if something can be exploited via a pull request, we consider that a high risk. There are currently 18 checks made across 3 themes: holistic security practises, source code risk assessment and build process risk assessment.

You can learn more about the scoring criteria, risks, and remediation suggestions for each check in the [detailed documentation](https://github.com/ossf/scorecard/blob/main/docs/checks.md).

<div class="w-full items-center justify-center">

![What Scorecard assesses](assets/diagram-flower.svg)

</div>

#### Holistic security practises

| Code vulnerabilities                                                                          | Description                                                                              | Risk |
|-----------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------|------|
| [Vulnerabilities](https://github.com/ossf/scorecard/blob/main/docs/checks.md#vulnerabilities) | Does the project have unfixed vulnerabilities? Uses the [OSV service](https://osv.dev/). | High |

| Maintenance                                                                                                 | Description                                                                                                                                                                                                                                                                                          | Risk   |
|-------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|
| [Dependency Update Tool](https://github.com/ossf/scorecard/blob/main/docs/checks.md#dependency-update-tool) | Does the project use tools to help update its dependencies e.g. [Dependabot](https://docs.github.com/en/code-security/supply-chain-security/managing-vulnerabilities-in-your-projects-dependencies/configuring-dependabot-security-updates), [RenovateBot](https://github.com/renovatebot/renovate)? | High   |
| [Maintained](https://github.com/ossf/scorecard/blob/main/docs/checks.md#maintained)                         | Is the project maintained?                                                                                                                                                                                                                                                                           | High   |
| [Security Policy](https://github.com/ossf/scorecard/blob/main/docs/checks.md#security-policy)               | Does the project contain a [security policy](https://docs.github.com/en/free-pro-team@latest/github/managing-security-vulnerabilities/adding-a-security-policy-to-your-repository)?                                                                                                                  | Medium |
| [Licence](https://github.com/ossf/scorecard/blob/main/docs/checks.md#license)                               | Does the project declare a licence?                                                                                                                                                                                                                                                                  | Low    |
| [CII Best Practices](https://github.com/ossf/scorecard/blob/main/docs/checks.md#cii-best-practices)         | Does the project have a [CII Best Practices Badge](https://bestpractices.coreinfrastructure.org/en)?                                                                                                                                                                                                 | Low    |

| Continuous testing                                                              | Description                                                                                                                                                                                                                                                                                                                    | Risk   |
|---------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|
| [CI Tests](https://github.com/ossf/scorecard/blob/main/docs/checks.md#ci-tests) | Does the project run tests in CI, e.g. [GitHub Actions](https://docs.github.com/en/free-pro-team@latest/actions), [Prow](https://github.com/kubernetes/test-infra/tree/master/prow)?                                                                                                                                           | Low    |
| [Fuzzing](https://github.com/ossf/scorecard/blob/main/docs/checks.md#fuzzing)   | Does the project use fuzzing tools, e.g. [OSS-Fuzz](https://github.com/google/oss-fuzz)?                                                                                                                                                                                                                                       | Medium |
| [SAST](https://github.com/ossf/scorecard/blob/main/docs/checks.md#sast)         | Does the project use static code analysis tools, e.g. [CodeQL](https://docs.github.com/en/free-pro-team@latest/github/finding-security-vulnerabilities-and-errors-in-your-code/enabling-code-scanning-for-a-repository#enabling-code-scanning-using-actions), [LGTM](https://lgtm.com/), [SonarCloud](https://sonarcloud.io/)? | Medium |

#### Source risk assessment

| Name                                                                                               | Description                                                                                                                                           | Risk     |
|----------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| [Binary Artifacts](https://github.com/ossf/scorecard/blob/main/docs/checks.md#binary-artifacts)    | Is the project free of checked-in binaries?                                                                                                           | High     |
| [Branch Protection](https://github.com/ossf/scorecard/blob/main/docs/checks.md#branch-protection)  | Does the project avoid dangerous coding patterns in GitHub Actions?                                                                                   | High     |
| [Dangerous Workflow](https://github.com/ossf/scorecard/blob/main/docs/checks.md#branch-protection) | Does the project use [Branch Protection](https://docs.github.com/en/free-pro-team@latest/github/administering-a-repository/about-protected-branches)? | Critical |
| [Code Review](https://github.com/ossf/scorecard/blob/main/docs/checks.md#code-review)              | Does the project require code review before code is merged?                                                                                           | High     |
| [Contributors](https://github.com/ossf/scorecard/blob/main/docs/checks.md#contributors)            | Does the project have contributors from at least two different organizations?                                                                         | Low      |

#### Build risk assessment

| Name                                                                                                  | Description                                                                                                                                                                                                                      | Risk   |
|-------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|
| [Pinned Dependencies](https://github.com/ossf/scorecard/blob/main/docs/checks.md#pinned-dependencies) | Does the project declare and pin [dependencies](https://docs.github.com/en/free-pro-team@latest/github/visualizing-repository-data-with-graphs/about-the-dependency-graph#supported-package-ecosystems)?                         | Medium |
| [Token Permissions](https://github.com/ossf/scorecard/blob/main/docs/checks.md#token-permissions)     | Does the project declare GitHub workflow tokens as [read only](https://docs.github.com/en/actions/reference/authentication-in-a-workflow)?                                                                                       | High   |
| [Packaging](https://github.com/ossf/scorecard/blob/main/docs/checks.md#packaging)                     | Does the project build and publish official packages from CI/CD, e.g. [GitHub Publishing](https://docs.github.com/en/free-pro-team@latest/actions/guides/about-packaging-with-github-actions#workflows-for-publishing-packages)? | Medium |
| [Signed Releases](https://github.com/ossf/scorecard/blob/main/docs/checks.md#signed-releases)         | Does the project [cryptographically sign releases](https://wiki.debian.org/Creating%20signed%20GitHub%20releases)?                                                                                                               | High   |

<br/>

> Machine checkable properties are an essential part of a sound security process. That’s why we have incorporated Security Scorecards *[i.e., OpenSSF Scorecard]* into our dependency acceptance criteria.

<div class="text-right mb-60"><cite>Harvey Tuch, <a href="https://www.envoyproxy.io/">Envoy</a></cite></div>

### Use cases

##### OpenSSF Scorecard reduces the effort required to continually evaluate changing packages when maintaining a project’s supply chain

#### For individual maintainers

Scorecard is helpful as a pre-launch security checker for a new OSS project or to help to plan improvements to an existing one. If a project is well maintained, it’s more likely to be used by others instead of an alternative. It can also be used to check a new dependency being added to a project, so a maintainer can make an informed decision about the risk of doing so.

#### For an organisation

Scorecard can be included in the continuous integration/continuous deployment processes using the GitHub action and run by default on pull requests.

#### For consumers

Scorecard helps to make informed decisions about security risks and vulnerabilities. Using the public data, it is also possible to evaluate the security posture of over 1 million of the most used OSS projects.

<br/>

### About the project name

This project was initially called "Security Scorecards" but that form wasn't used consistently. In particular, the repo was named "scorecard" and so was the program. Over time people started referring to either form (singular and plural), with or without "Security", and the inconsitency became prevalent. To end this situation the decision was made to consolidate over the use of the singular form in keeping with the repo and program name, drop the "Security" part and use "OpenSSF" instead to ensure uniqueness. One should therefore refer to this project as "OpenSSF Scorecard" or "Scorecard" for short.

### Part of the OSS community

<div class="w-full md:w-3/4 inline-flex items-center gap-x-40">

<div>

![cisco](assets/logos/cisco.svg)

</div>

<div>

![datto](assets/logos/datto.svg)

</div>

<div>

![endor](assets/logos/endor.svg)

</div>

<div>

![google](assets/logos/google.svg)

</div>

<div>

& many others

</div>

</div>

OpenSSF Scorecard is being [developed and facilitated by contributors](https://github.com/ossf/scorecard/graphs/contributors) from across the OSS ecosystem.

We're part of the [Open Source Security Foundation (OpenSSF)](https://openssf.org), a cross-industry collaboration that brings together OSS security initiatives under one foundation and seeks to improve the security of OSS by building a broader community, targeted initiatives, and best practises.

OpenSSF launched Scorecard in November 2020 with the intention of auto-generating a “security score” for open source projects to help users as they decide the trust, risk, and security posture for their use case.

<br/>

### Get involved

![Open Source Security Foundation](assets/logos/openssf.svg)

Scorecard is part of the [OpenSSF Best Practices Working Group](https://github.com/ossf/wg-best-practices-os-developers).

If you want to get involved in the OpenSSF Scorecard community or have ideas you'd like to chat about, [we'd love to connect](https://github.com/ossf/scorecard#connect-with-the-scorecard-community).

</section>
