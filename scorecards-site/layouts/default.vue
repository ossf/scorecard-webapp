<!-- eslint-disable vue/attribute-hyphenation -->
<template>
  <div v-if="loading">loading...</div>
  <div v-else>
    <AppHeader />
    <main>
      <slot class="w-full" />
    </main>
    <AppFooter :navigation="footerNavLinks" :socialLinks="footerSocialLinks" />
  </div>
</template>

<script>
import AppHeader from '@/modules/Header/Header.vue'
import AppFooter from '@/modules/Footer/Footer.vue'

export default {
  name: 'MainLayout',
  components: {
    AppHeader,
    AppFooter,
  },

  data: () => ({
    scrollPos: '',
    isScrolling: false,
    loading: false,
    headerNavLinks: null,
    footerNavLinks: null,
    footerSocialLinks: null,
    observer: undefined,
    commits: null,
    tocList: [],
    navList: [],
    mobileNavOpen: false,
  }),

  computed: {
    scrollPosX() {
      return window.scrollY
    },
  },

  watch: {
    $route(to, from) {
      this.$bus.emit('openNavigation', false)
      this.$bus.on('storeTocs', (payload) => {
        this.tocList = payload
      })
    },
  },

  mounted() {
    window.addEventListener('scroll', this.getScrollPos)

    this.observer = new IntersectionObserver((entries) => {
      this.$bus.emit('observer.observed', entries)
    })

    this.$bus.emit('observer.created', this.observer)

    this.$bus.on('openNavigation', (payload) => {
      this.mobileNavOpen = payload
    })
  },

  created() {
    this.$bus.on('storeTocs', (payload) => {
      this.tocList = payload
    })
  },

  beforeUnmount() {
    window.removeEventListener('scroll', this.getScrollPos)
    this.$bus.off('openNavigation')
    this.$bus.off('storeTocs')
  },

  methods: {
    getScrollPos() {
      if (window.scrollY > 0) {
        this.isScrolling = true
        this.scrollPos = window.scrollY
      } else {
        this.isScrolling = false
      }
    },
  },
}
</script>
<style lang="scss">
.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition-duration: 0.5s;
  transition-property: height, opacity, transform;
  transition-timing-function: cubic-bezier(0.55, 0, 0.1, 1);
  overflow: hidden;
}
.slide-left-enter,
.slide-right-leave-active {
  opacity: 0;
  transform: translateX(2em);
}
.slide-left-leave-active,
.slide-right-enter {
  opacity: 0;
  transform: translateX(-2em);
}
</style>
<!-- eslint-enable -->
