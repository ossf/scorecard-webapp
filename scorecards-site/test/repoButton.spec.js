import { shallowMount } from '@vue/test-utils'
import RepoButton from '@/components/repo-button.vue'

describe('RepoButton.vue', () => {
  it('renders props.stars when passed', () => {
    const stars = 200
    const wrapper = shallowMount(RepoButton, {
      propsData: { stars }
    })
    expect(wrapper.props().stars).toBe(200);
  })

  it('has the repo link in the component data', () => {
    const wrapper = shallowMount(RepoButton)
    expect(wrapper.find('a').attributes('href')).toBe('https://github.com/ossf/scorecard')
  })

  it('renders the svg github logo', () => {
    const wrapper = shallowMount(RepoButton)
    expect(wrapper.find('svg')).toBeTruthy()
  })

})
