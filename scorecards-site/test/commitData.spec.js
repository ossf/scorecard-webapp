import { shallowMount } from '@vue/test-utils'
import CommitData from '@/components/commit-data.vue'

describe('commit-data.vue', () => {
  it('has p with data', () => {
    const wrapper = shallowMount(CommitData)
    expect(wrapper.find('p')).toBeTruthy()
  })

  it('does render a string with date', () => {
    const wrapper = shallowMount(CommitData, {
        propsData: {
            latestCommit: '3/28/2022, 20:11:47'
        }
    })
    expect(wrapper.props().latestCommit).toBe('3/28/2022, 20:11:47');
  })
})
