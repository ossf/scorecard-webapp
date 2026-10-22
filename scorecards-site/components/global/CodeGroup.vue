<template>
  <ClientOnly>
    <div class="theme-code-group">
      <div class="theme-code-group__nav">
        <ul class="theme-code-group__ul">
          <li
            v-for="(tab, i) in codeTabs"
            :key="tab.id"
            class="theme-code-group__li"
          >
            <button
              class="theme-code-group__nav-tab"
              :class="{
                'theme-code-group__nav-tab-active': i === activeCodeTabIndex,
              }"
              @click="changeCodeTab(i)"
            >
              {{ tab.title }}
            </button>
          </li>
        </ul>
      </div>
      <slot />
      <pre v-if="codeTabs.length < 1" class="pre-blank">
// Make sure to add code blocks to your code group</pre>
    </div>
  </ClientOnly>
</template>

<script>
export default {
  name: 'CodeGroup',
  provide() {
    return { codeGroup: this }
  },
  data() {
    return {
      codeTabs: [],
      activeCodeTabIndex: 0,
    }
  },
  methods: {
    changeCodeTab(index) {
      this.activeCodeTabIndex = index
    },
    // Each CodeBlock calls this once and keeps the returned id. Titles can
    // repeat, so blocks match on the id. The last block marked `active` opens
    // first. With none marked, the first tab opens.
    registerTab(title, active) {
      const id = this.codeTabs.length
      this.codeTabs.push({ id, title })
      if (active) {
        this.activeCodeTabIndex = id
      }
      return id
    },
  },
}
</script>

<style lang="scss" scoped>
.theme-code-group {
}
.theme-code-group__nav {
  background-color: black;
  border-top-left-radius: 6px;
  border-top-right-radius: 6px;
  border-bottom: 2px solid #302825;
}
.theme-code-group__ul {
  margin: auto 0;
  padding-left: 0;
  display: inline-flex;
  list-style: none;
}
.theme-code-group__li {
  padding: 5px 0;
  margin-left: 1.1em;
  &:before {
    display: none;
  }
}
.theme-code-group__nav-tab {
  border: 0;
  padding: 4px 10px;
  cursor: pointer;
  background: transparent;
  border-radius: 4px;
  font-family: 'Public Sans';
  font-style: normal;
  font-weight: normal;
  font-size: 18px;
  color: #feece3;
  transition: all 0.4s ease-in-out;
}
.theme-code-group__nav-tab-active {
  background: #302825;
  transition: all 0.4s ease-in-out;
}
.pre-blank {
  color: #42b983;
}
</style>
