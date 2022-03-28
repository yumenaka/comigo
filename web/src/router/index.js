//两种历史模式的区别：https://router.vuejs.org/zh/guide/essentials/history-mode.html
// import { createRouter, createWebHistory } from 'vue-router'
import { createRouter, createWebHashHistory } from 'vue-router'

// 1. 定义路由组件.
// 也可以从其他文件导入

//vite默认不支持@别名，但通过配置vite.config.js启用
// import App from '@/App.vue'
import ScrollMode from '@/views/ScrollMode.vue'
import FlipMode from '@/views/FlipMode.vue'
import BookShelf from '@/views/BookShelf.vue'
import NotFound from '@/views/NotFound404.vue'

//https://localhost/#/s/bookid


const routes = [
  {
    path: '/',//  以 / 开头的嵌套路径将被视为根路径。这允许你利用组件嵌套，而不必使用嵌套的 URL。
    component: BookShelf,
    name: 'BookShelf',
  },
  {
    path: '/child_shelf/:id',//  以 / 开头的嵌套路径将被视为根路径。这允许你利用组件嵌套，而不必使用嵌套的 URL。
    component: BookShelf,
    name: 'ChildBookShelf',
  },
  {
    path: '/scroll/:id',
    component: ScrollMode,
    name: 'ScrollMode',
  },
  {
    path: '/flip/:id',
    component: FlipMode,
    name: 'FlipMode',
  },
  {
    path: '/about',
    component: () => import('@/views/About.vue'),
    name: 'About',
  },
  // 将匹配所有内容并将其放在 `$route.params.pathMatch` 下
  { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFound },
]

const router = createRouter({
  // HTML5 模式：用户在浏览器中直接访问 https://example.com/user/id，会得到一个 404 错误。
  // 所以需要在服务器上添加一个简单的回退路由。如果 URL 不匹配任何静态资源，它应提供与你的应用程序中的 index.html 相同的页面。
  // history: createWebHistory(process.env.BASE_URL),

  // 为了简单起见，目前使用 hash 模式。它在内部传递的实际 URL 之前使用了一个哈希字符（#）。
  // 由于这部分 URL 从未被发送到服务器，所以它不需要在服务器层面上进行任何特殊处理。
  history: createWebHashHistory(),
  routes
})

export default router
