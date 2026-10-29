---
title: Evaluate a dependency
description: For people, organizations, and platforms deciding whether to depend on a project, or showing results to others.
weight: 2
---

Scorecard helps you judge the risks a dependency brings, so you can accept
them, choose an alternative, or work with its maintainers.

- **One project:** look it up in the [viewer](/viewer/), or run the
  [command line](/get-started/cli/) if it isn't in the weekly scan.
- **Many projects:** query the [REST API](/get-started/rest-api/) from your
  own tooling.
- **Policies:** check the specific practices you require, rather than an
  aggregate score. See [Understanding results](/docs/understanding-results/).

## At scale

Organizations, security teams, and package registries or platforms can use
Scorecard results across many projects. [deps.dev](https://deps.dev) is one
example.

- Get results from the [REST API](/get-started/rest-api/). The data is
  licensed under the [CDLA Permissive 2.0](https://cdla.dev/permissive-2-0).
- Link each result to the viewer at
  `https://scorecard.dev/viewer/?uri=<platform>/<owner>/<repository>`, so
  people can see the details. See [Viewer](/get-started/viewer/).
