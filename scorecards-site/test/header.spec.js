import { shallowMount, mount } from '@vue/test-utils'
import repoData from './scorecardrepo.json'
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
          latestCommit: '4/2/2022, 01:15:44',
          stars: 200
        }
      }
    })
    expect(wrapper.text()).toMatch('4/2/2022, 01:15:44')
    expect(wrapper.vm.stars).toEqual(200)
    expect(wrapper.vm.latestCommit).toEqual('4/2/2022, 01:15:44')
    expect(wrapper.find('span').text()).toMatch('200')
  })

  it('works with async', async () => {
    const wrapper = mount(Header, {
      data() {
        return {
          stars: null
        }
      }
    })
    const response = repoData;
    const data = await response;
    wrapper.vm.stars = data.stargazers_count;

    await expect(response).not.toBeNull()
  })

})
