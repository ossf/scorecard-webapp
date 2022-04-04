import { shallowMount } from '@vue/test-utils'
import Footer from '@/modules/Footer/Footer.vue'

describe('Footer.vue', () => {
  it('renders Footer', () => {
    const wrapper = shallowMount(Footer)
    expect(wrapper.findAll('footer')).toHaveLength(1)
  })

  it('shows a with href of /', () => {
    const wrapper = shallowMount(Footer)
    expect(wrapper.find('a').attributes('href')).toBe('/')
  })

  it('has text', () => {
    const wrapper = shallowMount(Footer)
    expect(wrapper.findAll('div').at(1).text()).toBe('© 2022 The Linux Foundation, under the terms of the Apache License 2.0.')
  })

})
