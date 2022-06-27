<template>
  <div class="justify-around items-center">
    <p class="subheading px-16 text-16">
      project status: <strong>alpha</strong>
    </p>
    <p class="subheading px-16 text-16">last commit: {{ latestCommit }}</p>
  </div>
</template>

<script>
export default {
  name: 'Test',
  props: {
    latestCommit: {
      type: String,
      default: null,
    },
  },
  data() {
    return {}
  },
  created() {},
  methods: {},
}
</script>
