# OpenSSF Scorecard website

[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/ossf/scorecard-webapp/badge)](https://api.securityscorecards.dev/projects/github.com/ossf/scorecard-webapp)
[![Netlify Status](https://api.netlify.com/api/v1/badges/d631bbe2-0e67-48ae-81a7-d7015195c9fd/deploy-status)](https://app.netlify.com/sites/ossf-scorecard/deploys)

## scorecard-webapp

Code for https://securityscorecards.dev ([`./scorecards-site`](./scorecards-site)).

The site is deployed on Netlify and the deployment configuration is in
[netlify.toml](./netlify.toml). Any changes committed to
[netlify.toml](./netlify.toml) and [scorecards-site/](./scorecards-site) on
`main` branch gets automatically deployed to production. So please make sure to
review deploy previews when making changes to the site. The documentation for 
local development can be found [here](/scorecards-site/README.md)

The site's viewer fetches results from the OpenSSF Scorecard results API at
https://api.securityscorecards.dev. The API's source has moved to
[`ossf/scorecard-infra`](https://github.com/ossf/scorecard-infra/tree/main/api) —
report API issues there.
