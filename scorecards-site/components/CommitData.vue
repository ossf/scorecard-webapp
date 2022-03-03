<template>
  <div class="justify-around items-center">
    <p class="subheading px-12">project status: <strong>alpha</strong></p>
    <p class="subheading px-12">last commit: {{ latestCommit }}</p>
  </div>
</template>

<script>
export default {
  name: "Test",
  created() {},
  data() {
    return {};
  },
  props: {
    latestCommit: null,
  },
  methods: {},
};
</script>
