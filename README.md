# OpenSSF Scorecard website

[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/ossf/scorecard-webapp/badge)](https://api.scorecard.dev/projects/github.com/ossf/scorecard-webapp)
[![Netlify Status](https://api.netlify.com/api/v1/badges/d631bbe2-0e67-48ae-81a7-d7015195c9fd/deploy-status)](https://app.netlify.com/sites/ossf-scorecard/deploys)

## scorecard-webapp

Code for [`https://scorecard.dev`](https://scorecard.dev). The website is in
[`site/`](/site): plain Markdown pages built with [Hugo](https://gohugo.io)
and the [Hextra](https://github.com/imfing/hextra) theme.

The site is deployed on Netlify and the deployment configuration is in
[netlify.toml](/netlify.toml). Any changes committed to the Netlify
configuration and `site/` on the `main` branch are deployed to production,
so please review the deploy preview when making changes to the site.

The site's viewer (`site/static/viewer/`) fetches results from the OpenSSF
Scorecard results API at <https://api.scorecard.dev>.

The API source code has moved to [`ossf/scorecard-infra`](https://github.com/ossf/scorecard-infra/tree/main/api) — report API issues there.

## Build and preview the site

You need Git and [Hugo extended](https://gohugo.io/installation/), at the
version set as `HUGO_VERSION` in [netlify.toml](/netlify.toml). Nothing
else is needed: no Node.js, Go, or Python.

```sh
git clone --recurse-submodules https://github.com/ossf/scorecard-webapp
cd scorecard-webapp
hugo server --source site                 # preview at http://localhost:1313
hugo --source site --minify --panicOnWarning   # build to site/public/
```

In an existing clone without the theme, run `git submodule update --init`.

## Content from `ossf/scorecard`

The Checks page, the FAQ, and the beginner's guide are built from files in
[`ossf/scorecard`](https://github.com/ossf/scorecard) at the release set as
`params.scorecardVersion` in [site/hugo.yaml](/site/hugo.yaml). Hugo fetches
them when it builds the site, so building needs network access to
`raw.githubusercontent.com`. To show a newer release, change that value.
To change the text of a check or a guide, edit it in `ossf/scorecard`.

## Content checks

The build fails on raw HTML, shortcodes, or directives in a page, missing
front matter or data fields, and links to missing pages or headings. Known
good and bad examples of each are in [`site/tests/`](/site/tests), and CI
checks that each bad example still fails.
