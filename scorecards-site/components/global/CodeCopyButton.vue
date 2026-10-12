<template>
  <button ref="copy" type="button" class="copy" aria-label="Copy code">
    <IconClipboardCheck v-if="state === 'copied'" class="w-26 h-26 mx-auto" />
    <IconClipboardCopy v-else class="w-26 h-26 mx-auto" />
  </button>
  <!-- The icon swap is visual only. Screen readers read this status instead. -->
  <span class="sr-only" role="status">{{
    state === 'copied' ? 'Code copied' : ''
  }}</span>
</template>

<script>
import Clipboard from 'clipboard'
import IconClipboardCheck from '../../assets/icons/icon-clipboard-check.svg'
import IconClipboardCopy from '../../assets/icons/icon-clipboard-copy.svg'
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
