import { mount } from '@vue/test-utils'
import Sidebar from '@/components/global/Sidebar.vue'

describe('Sidebar.vue', () => {

    beforeEach(() => {
        // IntersectionObserver isn't available in test environment
        const mockIntersectionObserver = jest.fn();
        mockIntersectionObserver.mockReturnValue({
          observe: () => null,
          unobserve: () => null,
          disconnect: () => null
        });
        window.IntersectionObserver = mockIntersectionObserver;
    });

    const navLinks = {
        toc:[
            {
                depth:1,
                id:"run",
                text:"Run",
            },
            {
                depth:2,
                id:"the",
                text:"the",
            },
            {
                depth:3,
                id:"checks",
                text:"checks",
            },
        ]
    }

  it('does render a nav', async () => {
    const wrapper = mount(Sidebar, {
        mocks: {
            $content: () => {},
            $nuxt: {
              $on: () => {}
            }
        },
        stubs:{
            NuxtLink: true
        }
    })
    await wrapper.setData({ navList: navLinks })
    expect(wrapper.findAll('nav')).toHaveLength(1)
  })

  it('does render a ul', async () => {
    const wrapper = mount(Sidebar, {
        mocks: {
            $content: () => {},
            $nuxt: {
              $on: () => {}
            }
        },
        stubs:{
            NuxtLink: true
        }
    })
    await wrapper.setData({ navList: navLinks })
    expect(wrapper.findAll('ul')).toHaveLength(1)
  })


  it('does render a ul with 3 li', async () => {
    const wrapper = mount(Sidebar, {
        mocks: {
            $content: () => {},
            $nuxt: {
              $on: () => {}
            }
        },
        stubs:{
            NuxtLink: true
        }
    })
    await wrapper.setData({ navList: navLinks })
    expect(wrapper.findAll('li')).toHaveLength(3)
  })

  it('does render li with text Run', async () => {
    const wrapper = mount(Sidebar, {
        mocks: {
            $content: () => {},
            $nuxt: {
              $on: () => {}
            }
        },
        stubs:{
            NuxtLink: true
        }
    })
    await wrapper.setData({ navList: navLinks })
    expect(wrapper.findAll('li').at(0).text()).toBe('Run')
  })
})
