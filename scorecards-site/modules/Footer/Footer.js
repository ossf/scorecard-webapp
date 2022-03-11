import Logo from '@/assets/icons/logo.svg?inline'
export default {
  name: 'FooterModule',
  components: {
    Logo,
  },
  data: () => ({
    globalFooter: null,
  }),

  computed: {},

  watch: {
    $route() {},
  },

  props: {
    navigation: Array,
    socialLinks: Array,
  },

  methods: {},

  mounted() {},
}
