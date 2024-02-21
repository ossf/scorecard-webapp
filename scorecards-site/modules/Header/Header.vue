<template>
  <header
    ref="header"
    class="py-30 w-full z-50 bg-pastel-white"
    :class="[bg, isScrolling ? 'scrolled' : '']"
  >
    <div class="container">
      <div class="flex flex-wrap justify-between items-start">
        <!-- Anchor tag wrapping the image for the logo -->
        <a href="/" title="Home" class="logo">
          <!-- Image tag for the SVG logo -->
          <img src="./logo.svg" alt="Company Logo" class="logo-image"
            :class="[
              $route.params.slug === 'trust-security'
                ? 'logo-white'
                : 'logo-black',
            ]"
          />
        </a>
        <RepoButton :stars="stars" :commits="commits" />
      </div>
    </div>
  </header>
</template>
<script lang="js" src="./Header.js" />
<style lang="scss" src="./Header.scss" />
