<template>
  <nav>
    <ul
      class="flex items-start"
      :class="[
        navType == 'social'
          ? 'justify-start md:justify-end w-full'
          : 'justify-start flex-wrap',
      ]"
    >
      <li>
        <NuxtLink
          v-if="navType == 'header'"
          class="mr-16 lg:mr-36 flex justify-between"
          :class="[text]"
          :to="`/`"
        >
          Overview
        </NuxtLink>
      </li>
      <li v-for="(navItem, index) in navList" :key="index">
        <NuxtLink
          v-if="navType == 'header'"
          class="mr-16 lg:mr-36 flex justify-between"
          :class="[text]"
          :to="`${navItem.menuLink}`"
        >
          {{ navItem.name }}
        </NuxtLink>
        <NuxtLink
          v-else-if="navType == 'footer'"
          class="footer--nav-button mr-16 lg:mr-36 flex justify-between h text-16 leading-24 text-gray-dark"
          :to="`${navItem.footerMenuLink}`"
        >
          {{ navItem.name }}
        </NuxtLink>
        <a
          v-else
          class="footer--nav-button md:mr-30 mr-16 flex justify-between h text-16 leading-24 text-gray-dark"
          target="_blank"
          :href="`${navItem.url ? navItem.url : navItem.externalLink}`"
        >
          <span
            v-if="navType == 'social' && navItem.name == 'GitHub'"
            class="mr-16"
          >
            <GitHubLogo :class="['fill-grey']" />
          </span>
          <span
            v-if="navType == 'social' && navItem.name == 'Twitter'"
            class="mr-16"
          >
            <TwitterLogo :class="['fill-grey']" />
          </span>
          {{ navItem.name }}
        </a>
      </li>
      <ul v-if="navType == 'header'" class="flex items-start">
        <li>
          <a
            class="mr-16 lg:mr-36 flex justify-between"
            :class="[text]"
            target="_blank"
            href="https://medium.com/sigstore"
          >
            Blog
          </a>
        </li>
        <li>
          <a
            class="mr-16 lg:mr-36 flex justify-between items-center"
            :class="[text]"
            target="_blank"
            href="https://docs.sigstore.dev"
          >
            Docs
            <span class="ml-6 h-16">
              <svg
                width="15"
                height="15"
                viewBox="0 0 24 20"
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
        <li>
          <a
            :class="[text]"
            href="https://github.com/sigstore"
            target="_blank"
            title="GitHub profile link"
          >
            <GitHubLogo
              :class="[
                $route.params.slug === 'trust-security'
                  ? 'fill-white'
                  : 'fill-black',
              ]"
            />
          </a>
        </li>
      </ul>
    </ul>
  </nav>
</template>
<script lang="js" src="./Navigation.js" />
