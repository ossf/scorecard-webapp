<template>
  <a :href="url" :title="`Link to ${title}`" class="github-action-cta">
    <slot />
  </a>
</template>

<script>
export default {
  name: "CodeBlock",
  props: {
    title: {
      type: String,
      required: true,
    },
    url: {
      type: String,
      required: true,
    },
  },
  mounted() {
    if (this.$parent && this.$parent.loadTabs) {
      this.$parent.loadTabs();
    }
  },
};
</script>

<style lang="scss" scoped>
.github-action-cta {
  display: none;

  & .nuxt-content-highlight > pre {
    border-radius: 0;
    border-top: 2px solid #302825;
    background-color: black;
    padding: 1.1em;
    & > code {
      padding: 0;
    }
  }
}
</style>
