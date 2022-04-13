import { defineConfig } from 'vite';
import path from 'path';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
import envCompatible from 'vite-plugin-env-compatible';
import { injectHtml } from 'vite-plugin-html';
import { viteCommonjs } from '@originjs/vite-plugin-commonjs';

// https://vitejs.dev/config/
export default defineConfig({

  // pdfjs-dist插件中使用了 es11 的语法，需要特殊配置：https://juejin.cn/post/6995856687106261000
  chainWebpack: config => {
    config.module.rule('pdfjs-dist').test({
      test: /\.js$/,
      include: path.join(__dirname, 'node_modules/pdfjs-dist')
    }).use('babel-loader').loader('babel-loader').options({
      presets: ['@babel/preset-env'],
      plugins: ['@babel/plugin-proposal-optional-chaining']
    })
  },

  // 静态资源基础路径 base: './' || '',
  // base: process.env.NODE_ENV === 'production' ? './' : '/',
  base: '/',
  resolve: {
    alias: [
      {
        find: /^~/,
        replacement: ''
      },
      {
        find: '@',///vite默认不支持@别名，但通过配置启用此特性
        replacement: path.resolve(__dirname, 'src')
      }
    ],
    extensions: [
      '.mjs',
      '.js',
      '.ts',
      '.jsx',
      '.tsx',
      '.json',
      '.vue'
    ]
  },
  plugins: [
    vue(),
    vueJsx(),
    viteCommonjs(),
    envCompatible(),
    injectHtml()
  ],
  server: {
    strictPort: false,//设置为 true 时，如果端口已被使用，则直接退出，而不会再进行后续端口的尝试。
    open: false, //开发服务器启动时，自动在浏览器中打开应用程序。
    port: 4080,//开发服务器端口。如果设端口已被使用，Vite 将自动尝试下一个可用端口。
    host: '0.0.0.0',//为开发服务器指定 ip 地址。 设置为 0.0.0.0 或 true 会监听所有地址，包括局域网和公共地址。
    //https://zxuqian.cn/vite-proxy-config/
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:1234/',
      },
      '/favicon.ico': {
        target: 'http://127.0.0.1:1234/',
      },
      '/raw': {
        target: 'http://127.0.0.1:1234/',
        changeOrigin: true
      },
      '/asssets': {
        target: 'http://127.0.0.1:1234/',
        changeOrigin: true
      },
      '/images': {
        target: 'http://127.0.0.1:1234/',
        changeOrigin: true
      },
      '/login': {
        target: 'http://127.0.0.1:1234/',
        changeOrigin: true
      },
      '/loginJSON': {
        target: 'http://127.0.0.1:1234/',
        changeOrigin: true
      },
      '/index.html': {
        ignorePath: true//不知道有没有用╮(￣▽￣")╭
      },
    }
  },
  build: {
    target: 'esnext',//设置最终构建的浏览器兼容目标。默认值是一个 Vite 特有的值——'modules'，这是指 支持原生 ES 模块的浏览器。另一个特殊值是 “esnext” —— 即假设有原生动态导入支持，并且将会转译得尽可能小。
    minify: 'terser',//如果 build.minify 选项为 'terser'， 'esnext' 将会强制降级为 'es2019'
    outDir: '../routers/static',
    emptyOutDir: true,//清除目标目录：https://cn.vitejs.dev/config/#build-emptyoutdir
  }
})
