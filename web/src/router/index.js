import Vue from 'vue'
import VueRouter from 'vue-router'
import Scroll from '../views/ScrollTemplate.vue'

Vue.use(VueRouter)

  const routes = [
  {
    path: '/',
    name: 'Scrool',
    component: Scroll
  },
  {
    path: '/sketch',
    name: 'Sketch',
    component: () => import( '../views/SketchTemplate.vue')
  }
]

const router = new VueRouter({
  routes
})

export default router
