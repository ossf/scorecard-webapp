# OpenSSF Scorecard API and website

[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/ossf/scorecard-webapp/badge)](https://api.securityscorecards.dev/projects/github.com/ossf/scorecard-webapp)
[![Netlify Status](https://api.netlify.com/api/v1/badges/d631bbe2-0e67-48ae-81a7-d7015195c9fd/deploy-status)](https://app.netlify.com/sites/ossf-scorecard/deploys)

## scorecard-webapp

Code for https://securityscorecards.dev
([`./scorecards-site`](./scorecards-site)) and
https://api.securityscorecards.dev ([`./app`](./app)).

The site is deployed on Netlify and the deployment configuration is in
[netlify.toml](./netlify.toml). Any changes committed to
[netlify.toml](./netlify.toml) and [scorecards-site/](./scorecards-site) on
`main` branch gets automatically deployed to production. So please make sure to
review deploy previews when making changes to the site.

The API uses [OpenAPI](https://www.openapis.org/) spec and
[go-swagger](https://goswagger.io/) to auto-generate server and client code. Any
changes committed to [openapi.yaml](./openapi.yaml) on the `main` branch gets
deployed to the staging site only. To make changes to the production API, a new
Git tag needs to be generated which will auto deploy the latest tag to
production.

The documentation for local development can be found [here](/scorecards-site/README.md)

## Release process

### GitHub release

Cut a release for the project via the GitHub UX or by pushing a new tag.

### TODO: automate these steps

Any updates made to [openapi.yaml](./openapi.yaml) needs to be deployed onto
Google Cloud Endpoints. To do that, follow these steps:

```
$ gcloud auth login
$ gcloud endpoints services deploy openapi.yaml --project openssf --quiet --format=json > /tmp/gcloud.json
$ wget https://raw.githubusercontent.com/GoogleCloudPlatform/esp-v2/master/docker/serverless/gcloud_build_image \
   --output-document=/tmp/gcloud_build_image
$ chmod +x /tmp/gcloud_build_image
$ /tmp/gcloud_build_image -c $(cat /tmp/gcloud.json | jq -r .serviceConfig.id) \
   -s $(cat /tmp/gcloud.json | jq -r .serviceConfig.name) \
   -p openssf -z us
$ gcloud run deploy scorecard-endpoints-prod \
   --image=<image-from-above-step> \
   --project=openssf
   # For region prompt, choose us-central1.
```
