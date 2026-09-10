# scorecards-site

## Pre-requisites

- [NVM](https://github.com/nvm-sh/nvm)
- NPM
- [Yarn](https://classic.yarnpkg.com/lang/en/docs/install/)

## Build Setup

```bash
# Get the supported Node version
$ nvm install 22

# install dependencies
$ yarn install

# serve with hot reload at localhost:3000
$ yarn dev

# build for production and launch server
$ yarn build
$ yarn start

# generate static project
$ yarn generate
```

For detailed explanation on how things work, check out the
[documentation](https://nuxtjs.org).

## E2E testing

In order to run the e2e tests we use Playwright, but we **do not include it as a project dependency**.  
To run tests locally, you must install Playwright manually and avoid committing these changes.

### Temporary local installation

```sh
yarn add -D @playwright/test chromatic @chromatic-com/playwright
```

These changes are temporary and should not be committed.

For reference you can follow the steps in the[e2e testing pipeline](../.github/workflows/playwright.yml)

### Run the tests

Basic test run:

```sh
yarn playwright test
```

Run tests with visual (HTML) reporting:

```sh
yarn playwright test --reporter=html
```

Open the HTML report:

```sh
yarn playwright show-report
```

### Testing approach

Test files are located in the `tests-e2e` directory.
They are used to perform basic interaction tests such as:

- clicking buttons
- navigation and URL redirections
- rendering checks
- simple behavior flows

To prevent visual regressions, we generate screenshots at key test steps and compare them against stored reference files.

Snapshots are located in `tests-e2e/*-snapshots`

If differences are detected, the HTML report will show a diff view and the tests will fail.

### Important

The E2E tests will **build the website and serve it locally during execution**.  
Before running the tests:

- stop any local development server you may have running
- ensure that `http://localhost:3000` is **not in use**

If this port is busy, the tests will fail.

### Clean up

Once you are done with the local tests:

- Remove all temporarily installed Playwright dependencies: `rm -rf node_modules`
- Discard the `devDependencies` modifications in `package.json` and `yarn.lock`.
- Reinstall your project dependencies:`yarn`

This restores your environment to a clean state.

## Special Directories

You can create the following extra directories, some of which have special
behaviors. Only `pages` is required; you can delete them if you don't want to
use their functionality.

### `assets`

The assets directory contains your uncompiled assets such as Stylus or Sass
files, images, or fonts.

More information about the usage of this directory in
[the documentation](https://nuxtjs.org/docs/2.x/directory-structure/assets).

### `components`

The components directory contains your Vue.js components. Components make up the
different parts of your page and can be reused and imported into your pages,
layouts and even other components.

More information about the usage of this directory in
[the documentation](https://nuxtjs.org/docs/2.x/directory-structure/components).

### `layouts`

Layouts are a great help when you want to change the look and feel of your Nuxt
app, whether you want to include a sidebar or have distinct layouts for mobile
and desktop.

More information about the usage of this directory in
[the documentation](https://nuxtjs.org/docs/2.x/directory-structure/layouts).

### `pages`

This directory contains your application views and routes. Nuxt will read all
the `*.vue` files inside this directory and setup Vue Router automatically.

More information about the usage of this directory in
[the documentation](https://nuxtjs.org/docs/2.x/get-started/routing).

### `plugins`

The plugins directory contains JavaScript plugins that you want to run before
instantiating the root Vue.js Application. This is the place to add Vue plugins
and to inject functions or constants. Every time you need to use `Vue.use()`,
you should create a file in `plugins/` and add its path to plugins in
`nuxt.config.js`.

More information about the usage of this directory in
[the documentation](https://nuxtjs.org/docs/2.x/directory-structure/plugins).

### `static`

This directory contains your static files. Each file inside this directory is
mapped to `/`.

Example: `/static/robots.txt` is mapped as `/robots.txt`.

More information about the usage of this directory in
[the documentation](https://nuxtjs.org/docs/2.x/directory-structure/static).

### `store`

This directory contains your Vuex store files. Creating a file in this directory
automatically activates Vuex.

More information about the usage of this directory in
[the documentation](https://nuxtjs.org/docs/2.x/directory-structure/store).
