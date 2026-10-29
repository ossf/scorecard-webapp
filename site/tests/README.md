# Content check fixtures

Each directory is a fixture set for one content check. CI
(`.github/workflows/build-website.yml`, job `site-checks`) builds every set
and asserts that `clean` passes and each other set fails with its expected
message. Directories with `content/` replace the site's content; directories
with `data/` replace its data and use `clean/content`. `checks/` holds
`checks.yaml` fixtures for the Checks page.
