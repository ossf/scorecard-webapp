import mitt from 'mitt'

export default defineNuxtPlugin(() => {
  return {
    provide: {
      bus: mitt(),
    },
  }
})
