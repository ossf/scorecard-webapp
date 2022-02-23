<template>
  <div>
    <!-- <component
      :is="section.type"
      v-for="(section, index) in page.sections"
      :key="index"
      v-bind="section"
    /> -->

    <nuxt-content
      class="prose prose-sm sm:prose lg:prose-lg xl:prose-2xl mx-auto"
      ref="nuxtContent"
      :document="page"
    />
  </div>
</template>

<script>
import { mapActions } from "vuex";
// import Sidebar from "@/modules/Sidebar.vue";
export default {
  components: {
    // Sidebar,
  },

  transition: "page",

  async asyncData({ $content, params, error }) {
    const slug = params.slug || "home";
    const page = await $content(slug).fetch();
    const toc = page.toc;

    if (!page) {
      return error({ statusCode: 404, message: "Page not found" });
    }

    return {
      page,
      toc,
    };
  },
  data() {
    return {
      animation: "",
      tocs: [],
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
              content: this.page.description,
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
              content: "",
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
  mounted() {
    // const firstBanner = this.page.sections.filter(
    //   (p) => p.type === 'textBanner'
    // )[0]
    // const textBannerCards = this.page.sections.filter(
    //   (p) => p.type === 'textBannerWithcards'
    // )
    // if (firstBanner) {
    //   this.setHeaderColour({
    //     bg: firstBanner.bgColour,
    //     text:
    //       this.$route.params.slug === 'trust-security'
    //         ? 'text-white'
    //         : 'text-gray-dark',
    //   })
    // } else {
    //   this.setHeaderColour({
    //     bg: textBannerCards[0].bgColour,
    //     text:
    //       this.$route.params.slug === 'trust-security'
    //         ? 'text-white'
    //         : 'text-gray-dark',
    //   })
    // }
  },

  methods: {
    ...mapActions("settings", ["setHeaderColour"]),
  },
};
</script>
<style lang="scss"></style>
