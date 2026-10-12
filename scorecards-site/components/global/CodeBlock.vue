<template>
  <div
    class="theme-code-block"
    :class="{ 'theme-code-block__active': isActive }"
  >
    <slot />
  </div>
</template>

<script>
export default {
  name: 'CodeBlock',
  inject: {
    codeGroup: { default: null },
  },
  props: {
    title: {
      type: String,
      required: true,
    },
    active: {
      type: Boolean,
      default: false,
    },
    single: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      tabId: null,
    }
  },
  computed: {
    isActive() {
      return this.codeGroup
        ? this.codeGroup.activeCodeTabIndex === this.tabId
        : this.active
    },
  },
  created() {
    if (this.codeGroup) {
      this.tabId = this.codeGroup.registerTab(this.title, this.active)
    }
  },
}
</script>

<style lang="scss" scoped>
.theme-code-block {
  display: none;

  // ProsePre renders this <pre> as slot content, so the scoped style needs
  // :deep(). The [class*='language-'] part beats the grey code style in
  // base.scss.
  & :deep(.nuxt-content-highlight > pre[class*='language-']) {
    border-bottom-left-radius: 0.5em;
    border-bottom-right-radius: 0.5em;
    border-top: 2px solid #302825;
    background-color: black;
    font-size: 16px;
    line-height: 1.75;
    padding: 1.1em;
    margin-top: 0;
    margin-bottom: 0;
    & > code {
      padding: 0;
    }
  }
}
.theme-code-block__active {
  display: block;
}
.theme-code-block > pre {
  background-color: orange;
}
</style>
