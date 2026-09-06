const config = {
  rootMargin: '0px',
  threshold: [0.2, 0.6],
}

const animateOnScrollObserver = new IntersectionObserver(function (
  entries,
  animateOnScrollObserver,
) {
  entries.forEach((entry) => {
    if (entry.isIntersecting) {
      if (entry.intersectionRatio > 0.1) {
        entry.target.classList.add('enter')
        animateOnScrollObserver.unobserve(entry.target)

        if (entry.target.getBoundingClientRect().top < 0) {
          animateOnScrollObserver.unobserve(entry.target)
        }
      }
    }
  })
}, config)

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.directive('animate-on-scroll', {
    mounted: (el) => {
      el.classList.add('before-enter')
      animateOnScrollObserver.observe(el)
    },
  })
})
