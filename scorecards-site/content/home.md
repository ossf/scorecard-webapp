---
title: Home
date: 2021-07-12T15:33:03.264Z
description: Homepage
slug: home
thumbnail: /img/icon.png
---
<sidebar ref="sideBar" class="sticky top-100 h-400 w-1/3 hidden md:block"></sidebar>

<section class="bg-orange prose md:prose-lg w-full">

## Run the checks


## First Heading run checks

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum

<code-group>
<code-block title="GitHub Action">

```bash
yarn create vuepress-site [optionalDirectoryName]
```

</code-block>

<code-block title="CLI" active>

```bash
npx create-vuepress-site [optionalDirectoryName]
```

</code-block>
</code-group>



## Sub Heading One run checks

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum

| Syntax      | Description |
| ----------- | ----------- |
| Header      | Title       |
| Paragraph   | Text        |

#### fourth Heading run checks

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum

<code-block title="Vue js" single active>

```javascript
<script>
mounted() {
    this.observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        const id = entry.target.getAttribute("id");
        if (entry.isIntersecting) {
          this.currentlyActiveToc = id;
        }
      });
    }, this.observerOptions);

    // Track all sections that have an `id` applied
    document
      .querySelectorAll(
        ".nuxt-content h1[id], .nuxt-content h2[id], .nuxt-content h3[id], .nuxt-content h4[id]"
      )
      .forEach((section) => {
        this.observer.observe(section);
      });
  },
</script>
```

</code-block>

## Learn more

### The importance of Security Scorecards
### Keeping ahead of security attacks 
### What is Security Scorecards
### Make better security decisions
### How Security Scorecards works
### Who it’s for

#### For individual maintainers


#### For an organisation
#### For consumers

### About the checks

#### Code vulnerabilities

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum

#### Maintenance
#### Continuous testing
#### Source integrity
#### Build integrity
#### Dependency integrity

### A trusted partner
### Get involved
### For the OS community


</section>
