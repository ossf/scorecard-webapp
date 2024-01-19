# DNS

Our DNS entries are delegated to zones the maintainers have configured in [Google Cloud DNS](https://cloud.google.com/dns).

We have one zone for scorecard.dev and one for securityscorecards.dev.
The NS values can be found by going to a DNS zone and clicking on "Registrar Setup".

This will look something like:
* ns-cloud-b1.googledomains.com.
* ns-cloud-b2.googledomains.com.
* ns-cloud-b3.googledomains.com.
* ns-cloud-b4.googledomains.com.

## Docs

The doc portion of the website is managed through Netlify. Follow their 
[instructions](https://docs.netlify.com/domains-https/custom-domains/configure-external-dns/)
for setting up external DNS.

For HTTPS, we use Netlify's [managed certificates](https://docs.netlify.com/domains-https/https-ssl/#netlify-managed-certificates).
They auto-generate and renew Let's Encrypt certs for the site.
We don't currently set CAA records, as we also have certificates for the API hosted on Google Cloud.

## API

The API portion is hosted on [Cloud Run](https://cloud.google.com/run).
We map our domain name to the service via Cloud Run domain mapping.
Follow [these instructions](https://cloud.google.com/run/docs/mapping-custom-domains#map) to do the mapping.