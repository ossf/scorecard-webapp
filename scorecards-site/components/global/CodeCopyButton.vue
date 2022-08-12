<template>
  <button ref="copy" class="copy">
    <IconClipboardCheck v-if="state === 'copied'" class="w-26 h-26 mx-auto" />
    <IconClipboardCopy v-else class="w-26 h-26 mx-auto" />
  </button>
</template>

<script>
import Clipboard from 'clipboard'
import IconClipboardCheck from '../../assets/icons/icon-clipboard-check.svg?inline'
import IconClipboardCopy from '../../assets/icons/icon-clipboard-copy.svg?inline'
export default {
  components: {
    IconClipboardCopy,
    IconClipboardCheck,
  },
  data() {
    return {
      state: 'init',
    }
  },
  mounted() {
    const copyCode = new Clipboard(this.$refs.copy, {
      target(trigger) {
        return trigger.previousElementSibling
      },
    })
    copyCode.on('success', (event) => {
      event.clearSelection()
      this.state = 'copied'
      window.setTimeout(() => {
        this.state = 'init'
      }, 2000)
    })
  },
}
</script>
