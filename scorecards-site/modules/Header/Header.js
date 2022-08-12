import { mapGetters } from 'vuex'
import CommitData from '@/components/CommitData.vue'
import Logo from '@/assets/icons/logo.svg?inline'
import RepoButton from '@/components/RepoButton.vue'

export default {
  name: 'Header',
  components: {
    Logo,
    RepoButton,
    CommitData,
  },
  data: () => ({
    globalHeader: null,
    globalHeaderMenu: null,
    scrollPos: '',
    isScrolling: false,
    apiURL:
      'https://api.github.com/repos/ossf/scorecard/commits?per_page=3&sha=',
    branches: ['main'],
    currentBranch: 'main',
    commits: null,
    stars: null,
    latestCommit: null,
    navOpen: {
      type: Boolean,
      default: false,
    },
  }),

  computed: mapGetters({
    bg: 'settings/bg',
    text: 'settings/textColor',
  }),

  props: {
    navigation: Array,
    socialLinks: Array,
  },

  methods: {
    openNavDrawer() {
      this.$nuxt.$emit('openNavigation', true)
    },
    async fetchData() {
      // TODO: store this is state/cache so we do not have to load every time
      const options = {
        year: 'numeric',
        month: 'numeric',
        day: 'numeric',
        hour: 'numeric',
        minute: 'numeric',
        second: 'numeric',
        hour12: false,
        timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
      }
      const response = await fetch(this.apiURL)
      const data = await response.json()
      const d = data[0].commit.committer.date
      this.latestCommit = new Intl.DateTimeFormat('en-US', options).format(
        new Date(d)
      )
    },
    async getTotalCommits(owner, repo) {
      // TODO: store this is state/cache so we do not have to load every time
      const url = `https://api.github.com/repos/${owner}/${repo}`

      const response = await fetch(url)
      const data = await response.json()
      this.stars = data.stargazers_count
    },
  },

  created() {
    this.getTotalCommits('ossf', 'scorecard')
    this.fetchData()
  },

  mounted() {
    window.addEventListener('scroll', this.getScrollPos)
    this.getTotalCommits('ossf', 'scorecard')
    this.fetchData()
  },

  beforeDestroy() {
    window.removeEventListener('scroll', this.getScrollPos)
  },
}
