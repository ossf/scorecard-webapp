# OpenSSF Scorecard website

[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/ossf/scorecard-webapp/badge)](https://api.scorecard.dev/projects/github.com/ossf/scorecard-webapp)
[![Netlify Status](https://api.netlify.com/api/v1/badges/d631bbe2-0e67-48ae-81a7-d7015195c9fd/deploy-status)](https://app.netlify.com/sites/ossf-scorecard/deploys)

## scorecard-webapp

Code for [`https://scorecard.dev`](/scorecards-site).

The site is deployed on Netlify and the deployment configuration is in [netlify.toml](/netlify.toml).

Any changes committed to the Netlify configuration and [scorecards-site/](/scorecards-site) on
`main` branch gets automatically deployed to production, so please make sure to review deploy
previews when making changes to the site.

The documentation for local development can be found [here](/scorecards-site/README.md).

The site's viewer fetches results from the OpenSSF Scorecard results API at
https://api.scorecard.dev.

The API source code has moved to [`ossf/scorecard-infra`](https://github.com/ossf/scorecard-infra/tree/main/api) — report API issues there.
