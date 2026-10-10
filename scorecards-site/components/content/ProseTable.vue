<template>
  <div
    ref="scroll"
    class="table-scroll"
    :tabindex="overflows ? 0 : undefined"
    :role="overflows ? 'region' : undefined"
    :aria-label="overflows ? label : undefined"
  >
    <table v-bind="$attrs">
      <slot />
    </table>
  </div>
</template>

<script>
const firstHeader = (wrapper) => {
  const th = wrapper.querySelector('th')
  return th ? th.textContent.trim() : ''
}

// Name the table after its first header cell, such as "Maintenance". The source
// and build risk tables both start with "Name", so when another table has the
// same first header, use the heading above the table instead.
function tableName(wrapper) {
  const name = firstHeader(wrapper)
  const shared = [...document.querySelectorAll('.table-scroll')].some(
    (other) => other !== wrapper && firstHeader(other) === name,
  )
  if (name && !shared) return name
  let el = wrapper.previousElementSibling
  while (el && !/^H[2-6]$/.test(el.tagName)) el = el.previousElementSibling
  return el ? el.textContent.trim() : name || 'Check'
}

export default {
  name: 'ProseTable',
  inheritAttrs: false,
  data() {
    return {
      overflows: false,
      label: '',
    }
  },
  mounted() {
    const scroll = this.$refs.scroll
    // Only a table that scrolls sideways becomes a Tab stop. Then keyboard
    // users can focus it and scroll with the arrow keys. Wide screens get no
    // extra Tab stops.
    this.resizeObserver = new ResizeObserver(() => {
      this.overflows = scroll.scrollWidth > scroll.clientWidth
      // Name it on first overflow, once every table on the page is rendered.
      if (this.overflows && !this.label) {
        this.label = `${tableName(scroll)} table`
      }
    })
    this.resizeObserver.observe(scroll)
    this.resizeObserver.observe(scroll.querySelector('table'))
  },
  beforeUnmount() {
    this.resizeObserver.disconnect()
  },
}
</script>
