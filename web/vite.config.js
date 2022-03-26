import { defineConfig } from 'vite';
import path from 'path';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
import envCompatible from 'vite-plugin-env-compatible';
import { injectHtml } from 'vite-plugin-html';
import { viteCommonjs } from '@originjs/vite-plugin-commonjs';

// https://vitejs.dev/config/
export default defineConfig({
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
    strictPort: false,
    open: false, //自动打开 
    port: 4080,
    host: '0.0.0.0',
    //https://zxuqian.cn/vite-proxy-config/
    proxy: {
      '/api': {
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
    outDir: '../routers/static',
    emptyOutDir: true,//清除目标目录：https://cn.vitejs.dev/config/#build-emptyoutdir
  }
})
