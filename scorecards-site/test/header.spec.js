import { shallowMount, mount } from '@vue/test-utils'
import Header from '@/modules/Header/Header.vue'

describe('Header.vue', () => {
  it('renders header', () => {
    const wrapper = shallowMount(Header)
    expect(wrapper.findAll('header')).toHaveLength(1)
  })

  it('shows a with href of /', () => {
    const wrapper = shallowMount(Header)
    expect(wrapper.find('a').attributes('href')).toBe('/')
  })

  it('shows a with title attribute', () => {
    const wrapper = shallowMount(Header)
    expect(wrapper.find('a').attributes('title')).toBe('Home')
  })

  it('shows a date in latest commit', () => {
    const wrapper = mount(Header, {
      data() {
        return {
          latestCommit: '4/2/2022, 01:15:44'
        }
      }
    })
    expect(wrapper.text()).toMatch('4/2/2022, 01:15:44')
  })

})
