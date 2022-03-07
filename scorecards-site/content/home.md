---
title: Home
date: 2021-07-12T15:33:03.264Z
description: Homepage
slug: home
thumbnail: /img/icon.png
---
<sidebar ref="sideBar" class="sticky top-100 h-400 w-1/3 hidden md:block"></sidebar>

<section class="bg-orange prose md:prose-lg w-full">

<h2 class="h1" id="run-the-checks">Run the checks</h2>

## Using the GitHub Action

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

<section class="highlight-section">
<h3>Using Github Action</h3>

<details><summary>A summary title </summary>

*Bold:* Something here

*Bold:* Needs to be bold. <sup>another</sup>

*Example:* text here.

> TODO() 
> something
> here?

</details>
</section>


<section class="highlight-section">
<h3>Using CLI</h3>

<code-group>
<code-block title="Bash" active>

```bash
yarn create vuepress-site [optionalDirectoryName]
```

</code-block>

<code-block title="Homebrew">

```bash
npx create-vuepress-site [optionalDirectoryName]
```

</code-block>

<code-block title="Docker">

```bash
npx create-vuepress-site [optionalDirectoryName]
```

</code-block>
</code-group>
</section>

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
          if(id === 'run-the-checks'){
            alert(id);
          }
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

<h2 class="h1" id="learn-more">Learn more</h2>

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
