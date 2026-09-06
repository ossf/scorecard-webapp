<template>
  <div class="bg-pastel-white">
    <section
      class="flex justify-center items-center relative min-h-mobile md:min-h-threeQuarters"
    >
      <div class="mx-auto w-full md:w-4/6 text-center hero-text">
        <h1>
          Build better security habits,<br />
          one test at a time
        </h1>
        <div class="pt-20 pb-32 text-22">
          Quickly assess open source projects for risky practices
        </div>
        <div class="flex justify-center items-center my-32">
          <button
            class="btn cta mx-12"
            @click="scrollToAnchorPoint('run-the-checks')"
          >
            Run the checks
          </button>
          <button
            class="btn cta mx-12"
            @click="scrollToAnchorPoint('learn-more')"
          >
            Learn more
          </button>
        </div>
      </div>
    </section>
    <section
      id="video-section"
      ref="homeSection"
      class="md:min-h-threeQuarters"
    >
      <div class="mx-auto w-full md:w-3/4 rounded-lg overflow-hidden bg-black">
        <video
          ref="videoD"
          class="object-fit h-full w-full z-0 hidden md:block"
          autoplay
          loop
          muted
        >
          <source src="../assets/hero-video.mp4" type="video/mp4" />
          Your browser does not support the video tag.
        </video>
        <video
          ref="videoM"
          class="object-fit h-full w-full z-0 block md:hidden px-16"
          autoplay
          loop
          muted
        >
          <source src="../assets/hero-video-mobile.mp4" type="video/mp4" />
          Your browser does not support the video tag.
        </video>
      </div>
      <div class="my-64 text-center">
        <p class="subheading">Part of the Open Source Security Foundation</p>
        <div
          class="flex justify-center items-center my-16 mx-auto md:w-2/4 w-full px-32"
        >
          <div
            v-for="(logo, index) in logos"
            :key="index"
            class="w-6/12 md:4/12 flex justify-center md:mb-0 mb-32"
          >
            <img
              class="w-2/3 md:w-3/5 h-auto"
              :alt="`Logo ${index}`"
              :src="logo.pathLong"
            />
          </div>
        </div>
      </div>
    </section>
    <ContentRenderer
      v-if="page"
      ref="nuxtContent"
      class="nuxt-content container md:flex justify-start pb-132"
      :value="page"
    />
  </div>
</template>

<script>
import { computed, createApp } from 'vue'
import CodeCopyButton from '../components/global/CodeCopyButton'

const logoModules = import.meta.glob('../assets/logos/**/*.svg', {
  eager: true,
  query: '?url',
  import: 'default',
})

export default {
  components: {},

  async setup() {
    const nuxtApp = useNuxtApp()
    const config = useRuntimeConfig()
    const route = useRoute()

    const { data: page } = await useAsyncData('home', () =>
      queryCollection('content').path('/home').first(),
    )

    return nuxtApp.runWithContext(() => {
      if (!page.value) {
        throw createError({ statusCode: 404, statusMessage: 'Page not found' })
      }

      const toc = computed(() => flattenToc(page.value?.body?.toc?.links))

      useHead(() => ({
        title: page.value.title,
        titleTemplate: `OpenSSF Scorecard`,
        script: [
          {
            src: 'https://identity.netlify.com/v1/netlify-identity-widget.js',
          },
          {
            key: 'home',
            defer: true,
            src: '//gc.zgo.at/count.js',
            'data-goatcounter':
              'https://securityscorecards.goatcounter.com/count',
          },
          {
            innerHTML: JSON.stringify({
              '@context': 'https://schema.org',
              '@type': 'NewsArticle',
              mainEntityOfPage: {
                '@type': 'WebPage',
                '@id': `${config.public.frontendUrl}${route.path}`,
              },
              headline: `Home`,
              url: `${config.public.frontendUrl}${route.path}`,
            }),
            type: 'application/ld+json',
          },
        ],
        meta: [
          { charset: 'utf-8' },
          {
            name: 'viewport',
            content: 'width=device-width, initial-scale=1',
          },
          {
            name: 'description',
            content: page.value.description,
          },
          { name: 'format-detection', content: 'telephone=no' },
          // Twitter Card
          {
            name: 'twitter:card',
            content: config.public.siteName,
          },
          { name: 'twitter:title', content: page.value.title },
          {
            name: 'twitter:description',
            content: page.value.description,
          },
          // image must be an absolute path
          {
            name: 'twitter:image',
            content: '../assets/checks.png',
          },
          // Facebook OpenGraph
          { property: 'og:title', content: page.value.title },
          {
            property: 'og:site_name',
            content: config.public.siteName,
          },
          { property: 'og:type', content: 'website' },
          {
            property: 'og:image',
            content: '../assets/checks.png',
          },
          {
            property: 'og:description',
            content: page.value.description,
          },
        ],
        link: [
          { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
          {
            rel: 'canonical',
            href: `${config.public.frontendUrl}${route.path}`,
          },
        ],
      }))

      return { page, toc }
    })
  },

  data() {
    return {
      animation: '',
      tocs: [],
      logos: [],
      observer: null,
      isGoatCounterLoaded: false,
      observerOptions: {
        root: this.$refs.homeSection,
        rootMargin: '-50% 0% -50% 0%',
        threshold: 0,
      },
    }
  },
  computed: {},

  created() {
    if (this.toc) {
      this.$bus.emit('storeTocs', this.toc)
    }
  },
  beforeUnmount() {
    if (this.observer) {
      this.observer.disconnect()
    }
  },
  mounted() {
    this.importAll(logoModules)

    setTimeout(() => {
      const blocks = document.getElementsByClassName('nuxt-content-highlight')
      for (const block of blocks) {
        const mountEl = document.createElement('div')
        block.appendChild(mountEl)
        createApp(CodeCopyButton).mount(mountEl)
      }
    }, 100)
  },

  methods: {
    scrollToAnchorPoint(refName) {
      const el = document.getElementById(refName)
      el.scrollIntoView({ behavior: 'smooth' })
      // this.$router.push({ hash: `#${refName}` });
    },
    importAll(modules) {
      Object.entries(modules).forEach(([path, url]) =>
        this.logos.push({ pathLong: url, pathShort: path }),
      )
    },
  },
}
</script>
<style lang="scss"></style>
