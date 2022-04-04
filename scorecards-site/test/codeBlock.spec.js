import { shallowMount } from '@vue/test-utils'
import CodeBlock from '@/components/global/CodeBlock.vue'

describe('CodeBlock.vue', () => {
  it('has code', () => {
    const wrapper = shallowMount(CodeBlock)
    expect(wrapper.html()).toBe('<div class="theme-code-block"></div>')
  })

  it('shows a title', () => {
    const wrapper = shallowMount(CodeBlock, {
        propsData: {
          title: 'Homebrew'
        }
    })
    expect(wrapper.props().title).toBe('Homebrew');
  })

  it('should be active', () => {
    const wrapper = shallowMount(CodeBlock, {
        propsData: {
          active: true
        }
    })
    expect(wrapper.props().active).toBe(true);
  })

})
