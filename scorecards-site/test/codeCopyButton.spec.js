import { shallowMount, mount } from '@vue/test-utils'
import CodeCopyButton from '@/components/global/CodeCopyButton.vue'
import IconClipboardCheck from "@/assets/icons/icon-clipboard-check.svg?inline";
import IconClipboardCopy from "@/assets/icons/icon-clipboard-copy.svg?inline";

describe('CodeCopyButton.vue', () => {

    it('is a button', () => {

        const wrapper = mount(CodeCopyButton)

        const button = wrapper.find('button')

        expect(button).toBeTruthy()

    })

    it('renders the svg github logo', () => {

        const wrapper = mount(CodeCopyButton)

        const button = wrapper.find('button')

        expect(wrapper.find('svg')).toBeTruthy()

    })
    
})
