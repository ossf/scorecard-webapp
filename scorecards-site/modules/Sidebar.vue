<template>
  <ul>
    <li
      v-for="link in toc"
      :key="link.id"
      :class="{ toc2: link.depth === 2, toc3: link.depth === 3 }"
      @click="tableOfContentsHeadingClick(link)"
    >
      <NuxtLink
        :class="{
          'text-orange-dark hover:text-red-600': link.id === currentlyActiveToc,
          'text-black hover:gray-900': link.id !== currentlyActiveToc,
        }"
        role="button"
        class="transition-colors duration-75 text-base mb-2 block"
        :to="`#${link.id}`"
        >{{ link.text }}</NuxtLink
      >
    </li>
  </ul>
</template>

<script>
export default {
  name: "SideBar",
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
      .querySelectorAll(".nuxt-content h2[id], .nuxt-content h3[id]")
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
  data() {
    return {
      currentlyActiveToc: "",
      observer: null,
      observerOptions: {
        root: this.$refs.nuxtContent,
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
  },
};
</script>

<style lang="scss" scoped></style>
