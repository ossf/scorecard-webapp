<template>
  <div class="bg-pastel-white">
    <section
      class="flex justify-center items-center relative min-h-mobile md:min-h-threeQuarters"
    >
      <div class="pt-20 pb-32 text-22" style="background-color:white;">
        Take the <a href="https://forms.gle/aqxZwmVQzWJkNuso8">OpenSSF Scorecard User Survey</a>
      </div>
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
    <nuxt-content
      ref="nuxtContent"
      class="container md:flex justify-start pb-132"
      :document="page"
    />
  </div>
</template>

<script>
import { mapActions } from 'vuex'
import Vue from 'vue'
import CodeCopyButton from '../components/global/CodeCopyButton'

export default {
  components: {},
  transition(to, from) {
    if (!from) {
      return 'slide-left'
    }
    const fromIndex = from.query.i
    const toIndex = to.query.i
    return toIndex < fromIndex ? 'slide-right' : 'slide-left'
  },

  async asyncData({ $content, params, error }) {
    // const slug = params.slug || "home";
    const page = await $content('home').fetch()
    const toc = page.toc

    if (!page) {
      return error({ statusCode: 404, message: 'Page not found' })
    }

    return {
      page,
      toc,
    }
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
  head() {
    return this.page
      ? {
          title: this.page.title,
          titleTemplate: `OpenSSF Scorecard`,
          script: [
            {
              src: 'https://identity.netlify.com/v1/netlify-identity-widget.js',
            },
            {
              vmid: 'home',
              defer: true,
              // Changed after script load
              callback: () => {
                this.isGoatCounterLoaded = true
              },
              src: '//gc.zgo.at/count.js',
              'data-goatcounter':
                'https://securityscorecards.goatcounter.com/count',
            },
            {
              json: {
                '@context': 'https://schema.org',
                '@type': 'NewsArticle',
                mainEntityOfPage: {
                  '@type': 'WebPage',
                  '@id': `${process.env.VUE_APP_FRONTEND + this.$route.path}`,
                },
                headline: `Home`,
                url: `${process.env.VUE_APP_FRONTEND + this.$route.path}`,
              },
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
              hid: 'description',
              name: 'description',
              content: this.page.description,
            },
            { name: 'format-detection', content: 'telephone=no' },
            // Twitter Card
            {
              name: 'twitter:card',
              content: process.env.VUE_APP_SITENAME,
            },
            { name: 'twitter:title', content: this.page.title },
            {
              name: 'twitter:description',
              content: this.page.description,
            },
            // image must be an absolute path
            {
              name: 'twitter:image',
              content: '../assets/checks.png',
            },
            // Facebook OpenGraph
            { property: 'og:title', content: this.page.title },
            {
              property: 'og:site_name',
              content: process.env.VUE_APP_SITENAME,
            },
            { property: 'og:type', content: 'website' },
            {
              property: 'og:image',
              content: '../assets/checks.png',
            },
            {
              property: 'og:description',
              content: this.page.description,
            },
          ],
          link: [
            { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
            {
              rel: 'canonical',
              href: `${process.env.VUE_APP_FRONTEND + this.$route.path}`,
            },
          ],
        }
      : null
  },
  computed: {},

  created() {
    if (this.toc) {
      this.$nuxt.$emit('storeTocs', this.toc)
    }
  },
  beforeDestroy() {
    this.observer.disconnect()
  },
  mounted() {
    this.importAll(require.context('../assets/logos/', true, /\.svg$/))

    setTimeout(() => {
      const blocks = document.getElementsByClassName('nuxt-content-highlight')
      for (const block of blocks) {
        const CopyButton = Vue.extend(CodeCopyButton)
        const component = new CopyButton().$mount()
        block.appendChild(component.$el)
      }
    }, 100)
  },

  methods: {
    ...mapActions('settings', ['setHeaderColour']),
    scrollToAnchorPoint(refName) {
      const el = document.getElementById(refName)
      el.scrollIntoView({ behavior: 'smooth' })
      // this.$router.push({ hash: `#${refName}` });
    },
    importAll(r) {
      r.keys().forEach((key) =>
        this.logos.push({ pathLong: r(key), pathShort: key })
      )
    },
  },
}
</script>
<style lang="scss"></style>
