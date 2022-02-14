<template>
  <nav class="block md:hidden navigation--mobile w-full" :class="[bg]">
    <div class="container">
      <div class="py-30 flex flex-wrap justify-between items-start">
        <a href="/" title="Home link" class="logo z-20"
          ><Logo
            :class="[
              $route.params.slug === 'trust-security'
                ? 'logo-white'
                : 'logo-black',
            ]"
        /></a>
        <NavCloseButton
          class="block md:hidden w-32 h-32 z-20 relative"
          :class="[
            $route.params.slug === 'trust-security'
              ? 'fill-white'
              : 'fill-black',
          ]"
          @click="closeNavDrawer"
        />
      </div>
    </div>
    <ul
      class="md:flex flex-wrap items-start"
      :class="[navType == 'social' ? 'justify-end' : 'justify-start']"
    >
      <li class="mb-48">
        <NuxtLink
          class="footer--nav-button mr-36 flex justify-between text-19 hover:text-purple-light h"
          :class="[text]"
          :to="`/`"
        >
          Overview
        </NuxtLink>
      </li>
      <li v-for="(navItem, index) in navList" :key="index" class="mb-48">
        <NuxtLink
          v-if="navType == 'header'"
          class="footer--nav-button mr-36 flex justify-between text-19 hover:text-purple-light h"
          :class="[text]"
          :to="`${navItem.menuLink}`"
        >
          {{ navItem.name }}
        </NuxtLink>
        <NuxtLink
          v-else-if="navType == 'footer'"
          class="footer--nav-button mr-36 flex justify-between text-19 hover:text-purple-light h"
          :class="[text]"
          :to="`${navItem.footerMenuLink}`"
        >
          {{ navItem.name }}
        </NuxtLink>
        <a
          v-else
          class="footer--nav-button mr-8 flex justify-between text-19 hover:text-purple-light h"
          :class="[text]"
          :href="`${navItem.url}`"
        >
          <!-- TODO: refactor with dynamic namescase to use the assets folder -->
          <span
            v-if="navType == 'social' && navItem.name == 'GitHub'"
            class="mx-6"
          >
            <GitHubLogo :class="['fill-grey']" />
          </span>
          <span
            v-if="navType == 'social' && navItem.name == 'Twitter'"
            class="mx-6"
          >
            <TwitterLogo :class="['fill-grey']" />
          </span>
          {{ navItem.name }}
        </a>
      </li>
      <li>
        <a
          class="footer--nav-button mr-36 flex justify-between text-19 hover:text-purple-light h items-center"
          :class="[text]"
          target="_blank"
          href="https://docs.sigstore.dev"
        >
          Docs
          <span class="ml-6 h-20">
            <svg
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M18 13V19C18 19.5304 17.7893 20.0391 17.4142 20.4142C17.0391 20.7893 16.5304 21 16 21H5C4.46957 21 3.96086 20.7893 3.58579 20.4142C3.21071 20.0391 3 19.5304 3 19V8C3 7.46957 3.21071 6.96086 3.58579 6.58579C3.96086 6.21071 4.46957 6 5 6H11"
                :stroke="[
                  $route.params.slug === 'trust-security'
                    ? '#ffffff'
                    : '#444444',
                ]"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M15 3H21V9"
                :stroke="[
                  $route.params.slug === 'trust-security'
                    ? '#ffffff'
                    : '#444444',
                ]"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M10 14L21 3"
                :stroke="[
                  $route.params.slug === 'trust-security'
                    ? '#ffffff'
                    : '#444444',
                ]"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </span>
        </a>
      </li>
    </ul>
    <ul
      class="w-full p-60 border-t-2"
      :class="[bg != 'bg-purple-dark' ? 'border-gray-dark' : 'border-white']"
    >
      <li
        v-for="(navItem, index) in socialLinks"
        :key="index"
        :class="[index === socialLinks.length - 1 ? 'mb-0' : 'mb-48']"
      >
        <a
          class="footer--nav-button mr-8 flex justify-between text-19 hover:text-purple-light h"
          :class="[text]"
          :href="`${navItem.url}`"
        >
          {{ navItem.name }}
        </a>
      </li>
    </ul>
  </nav>
</template>
<script lang="js" src="./MobileNavigation.js" />
<style lang="scss" src="./MobileNavigation.scss" />
