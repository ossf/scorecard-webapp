<template>
  <nav>
    <ul v-if="navList" class="pl-6">
      <li
        v-for="link in navList.toc"
        :key="link.id"
        :class="{
          'pl-18': link.depth === 3,
          'parent-li':
            link.id === 'run-the-checks' ||
            (link.id === 'learn-more' && link.depth === 2),
          'pl-24': link.depth === 2,
        }"
        @click="tableOfContentsHeadingClick(link)"
      >
        <NuxtLink
          :class="{
            'text-orange-dark': link.id === currentlyActiveToc,
            'text-black hover:gray-900': link.id !== currentlyActiveToc,
            'nav-parent': link.id === 'run-the-checks' || link.id === 'learn-more',
          }"
          role="button"
          class="transition-colors duration-500 text-base mb-2 block toc-item"
          :to="`#${link.id}`"
          >{{ link.text }}</NuxtLink
        >
      </li>
    </ul>
  </nav>
</template>

<script>
export default {
  name: "SideBar",
  components: {},
  computed: {},
  mounted() {
    this.observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        const id = entry.target.getAttribute("id");
        if (entry.isIntersecting) {
          if (id === "video-section") {
            this.showLogo = !this.showLogo;
          }
          this.currentlyActiveToc = id;
        }
      });
    }, this.observerOptions);

    // Track all sections that have an `id` applied
    document
      .querySelectorAll(
        "#video-section, .nuxt-content h1[id], .nuxt-content h2[id], .nuxt-content h3[id], .nuxt-content h4[id]"
      )
      .forEach((section) => {
        this.observer.observe(section);
      });
  },
  beforeDestroy() {
    this.observer.disconnect();
  },
  created() {
    this.$nuxt.$on("setActiveToc", (id) => {
      this.currentlyActiveToc = id;
    });
  },
  this.getNavLinks();
  data() {
    return {
      currentlyActiveToc: "",
      showLogo: false,
      observer: null,
      navList: null,
      observerOptions: {
        root: this.$refs.nuxtContent,
        rootMargin: "-50% 0% -50% 0%",
        threshold: 0,
      },
    };
  },
  props: {
    toc: {
      default: null,
      type: Array,
    },
  },
  methods: {
    tableOfContentsHeadingClick(link) {
      this.currentlyActiveToc = link.id;
    },
    async getNavLinks() {
      const globalData = await this.$content("home")
        .only(["title", "slug", "toc"])
        .fetch();
      if (globalData) {
        this.navList = globalData;
      }
    },
  },
};
</script>

<style lang="scss" scoped></style>
