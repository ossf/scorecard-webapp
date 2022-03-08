<template>
  <div class="bg-pastel-white">
    <section
      class="flex justify-center items-center relative min-h-mobile md:min-h-threeQuarters"
    >
      <div class="mx-auto w-full md:w-4/6 text-center hero-text">
        <h1>Build better security habits, one test at a time</h1>
        <p class="subheading mt-32">
          Quickly assess open source projects for risky practices
        </p>
        <div class="flex justify-center items-center my-32">
          <button class="btn cta mx-12" @click="scrollToAnchorPoint('run-the-checks')">
            Run the checks
          </button>
          <button class="btn cta mx-12" @click="scrollToAnchorPoint('learn-more')">
            Learn more
          </button>
        </div>
      </div>
    </section>
    <section ref="homeSection" class="md:min-h-threeQuarters">
      <div class="mx-auto w-full md:w-3/4 rounded-lg overflow-hidden">
        <video
          ref="videoD"
          class="object-fit h-full w-full z-0 hidden md:block"
          autoplay
          muted
        >
          <source src="../assets/video.mp4" type="video/mp4" />
          Your browser does not support the video tag.
        </video>
        <video
          ref="videoM"
          class="object-fit h-full w-full z-0 block md:hidden px-16"
          autoplay
          muted
        >
          <source src="../assets/video-mobile.mp4" type="video/mp4" />
          Your browser does not support the video tag.
        </video>
      </div>
      <div class="my-64 text-center">
        <p class="subheading">Part of the Open Source Security Foundation</p>
        <div class="flex justify-center items-center my-16 mx-auto md:w-2/4 w-full px-32">
          <div
            class="w-6/12 md:4/12 flex justify-center md:mb-0 mb-32"
            v-for="(logo, index) in logos"
            :key="index"
          >
            <a href="https://openssf.org/" class="w-2/3 md:w-3/5 h-auto"><img class="w-full h-full" :src="logo.pathLong" /></a>
          </div>
        </div>
      </div>
    </section>
    <nuxt-content
      class="container md:flex justify-start pb-132"
      ref="nuxtContent"
      :document="page"
    />
  </div>
</template>

<script>
import { mapActions } from "vuex";

export default {
  transition(to, from) {
    if (!from) {
      return "slide-left";
    }
    const fromIndex = from.query.i;
    const toIndex = to.query.i;
    return toIndex < fromIndex ? "slide-right" : "slide-left";
  },
  components: {},

  async asyncData({ $content, params, error }) {
    // const slug = params.slug || "home";
    const page = await $content("home").fetch();
    const toc = page.toc;

    if (!page) {
      return error({ statusCode: 404, message: "Page not found" });
    }

    console.log(page);

    return {
      page,
      toc,
    };
  },
  data() {
    return {
      animation: "",
      tocs: [],
      logos: [],
      observer: null,
      observerOptions: {
        root: this.$refs.homeSection,
        rootMargin: "-50% 0% -50% 0%",
        threshold: 0,
      },
    };
  },
  head() {
    return this.page
      ? {
          title: this.page.title,
          titleTemplate: `%s · OSSF Scorecards`,
          script: [
            {
              src: "https://identity.netlify.com/v1/netlify-identity-widget.js",
            },
            {
              json: {
                "@context": "https://schema.org",
                "@type": "NewsArticle",
                mainEntityOfPage: {
                  "@type": "WebPage",
                  "@id": `${process.env.VUE_APP_FRONTEND + this.$route.path}`,
                },
                headline: `Home`,
                url: `${process.env.VUE_APP_FRONTEND + this.$route.path}`,
              },
              type: "application/ld+json",
            },
          ],
          meta: [
            { charset: "utf-8" },
            {
              name: "viewport",
              content: "width=device-width, initial-scale=1",
            },
            {
              hid: "description",
              name: "description",
              content: this.page.description,
            },
            { name: "format-detection", content: "telephone=no" },
            // Twitter Card
            {
              name: "twitter:card",
              content: process.env.VUE_APP_SITENAME,
            },
            { name: "twitter:title", content: this.page.title },
            {
              name: "twitter:description",
              content: this.page.description,
            },
            // image must be an absolute path
            {
              name: "twitter:image",
              content: "../assets/checks.png",
            },
            // Facebook OpenGraph
            { property: "og:title", content: this.page.title },
            {
              property: "og:site_name",
              content: process.env.VUE_APP_SITENAME,
            },
            { property: "og:type", content: "website" },
            {
              property: "og:image",
              content: "../assets/checks.png",
            },
            {
              property: "og:description",
              content: this.page.description,
            },
          ],
          link: [
            { rel: "icon", type: "image/x-icon", href: "/favicon.ico" },
            {
              rel: "canonical",
              href: `${process.env.VUE_APP_FRONTEND + this.$route.path}`,
            },
          ],
        }
      : null;
  },
  computed: {},

  created() {
    if (this.toc) {
      this.$nuxt.$emit("storeTocs", this.toc);
    }
  },
  beforeDestroy() {
    this.observer.disconnect();
  },
  mounted() {
    this.importAll(require.context("../assets/logos/", true, /\.svg$/));
    const videoD = this.$refs.videoD;
    const videoM = this.$refs.videoM;
    let playState = null;
    this.observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (!entry.isIntersecting) {
          videoD.pause();
          videoM.pause();
          playState = false;
        } else {
          videoD.play();
          videoM.play();
          playState = true;
        }
      });
    }, this.observerOptions);

    this.observer.observe(videoD, videoM);

    const onVisibilityChange = () => {
      if (document.hidden || !playState) {
        videoD.pause();
        videoM.pause();
      } else {
        videoD.play();
        videoM.play();
      }
    };

    document.addEventListener("visibilitychange", onVisibilityChange);
  },

  methods: {
    ...mapActions("settings", ["setHeaderColour"]),
    scrollToAnchorPoint(refName) {
      const el = document.getElementById(refName);
      el.scrollIntoView({ behavior: "smooth" });
      // this.$router.push({ hash: `#${refName}` });
    },
    importAll(r) {
      r.keys().forEach((key) => this.logos.push({ pathLong: r(key), pathShort: key }));
    },
  },
};
</script>
<style lang="scss"></style>
