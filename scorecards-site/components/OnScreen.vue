<script>
export default {
  props: {
    selector: { type: String, default: '' },
    once: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      items: [],
      onScreen: false,
    }
  },
  watch: {
    onScreen(onScreen) {
      if (this.once && onScreen) {
        this.$bus.on('observer.observed', (observer) => {
          for (const { el } of this.items) {
            observer.unobserve(el)
          }
        })
      }
    },
  },
  mounted() {
    this.items = [
      ...(this.selector
        ? document.querySelectorAll(this.selector)
        : [this.$el]),
    ].map((el) => ({ el, onScreen: false }))

    if (this.items.length) {
      this.$bus.on('observer.observed', (entries) => {
        entries.forEach((entry) => {
          const item = this.items.find(({ el }) => entry.target === el)

          if (item) {
            item.onScreen = entry.isIntersecting
          }
        })

        this.updateOnScreen()
      })

      this.$bus.on('observer.created', (observer) => {
        for (const { el } of this.items) {
          observer.observe(el)
        }
      })

      this.updateOnScreen()
    }
  },
  methods: {
    updateOnScreen() {
      this.onScreen = !!this.items.find((item) => item.onScreen)
    },
  },
  render() {
    return this.$slots.default?.({
      onScreen: this.onScreen,
    })
  },
}
</script>
