import { createRouter, createWebHashHistory } from 'vue-router'

// 1. 定义路由组件.
// 也可以从其他文件导入
// import Home from '../views/Home.vue'

import FlipMode from '@/components/FlipMode.vue'

import ScrollMode from '@/components/ScrollMode.vue'


const routes = [
  {
    path: '/',
    name: 'FlipMode',
    component: FlipMode
  },
  {
    path: '/scroll',
    name: 'ScrollMode',
    component: ScrollMode
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  }
]

const router = createRouter({
  // history: createWebHistory(process.env.BASE_URL),
    // 4. 内部提供了 history 模式的实现。为了简单起见，我们在这里使用 hash 模式。
  history: createWebHashHistory(),
  routes
})

export default router
