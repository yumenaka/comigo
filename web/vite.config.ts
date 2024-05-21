import { defineConfig } from "vite";
import path from "path";
import vue from "@vitejs/plugin-vue";
import vueJsx from "@vitejs/plugin-vue-jsx";
import envCompatible from "vite-plugin-env-compatible";
//vite-plugin-html:一个针对 index.html，提供压缩和基于 ejs 模板功能的 vite 插件。通过搭配 .env 文件，可以在开发或构建项目时，对 index.html 注入动态数据，例如替换网站标题。
//Vite 应用的 title 默认是写死的，如需要替换成实际的 title 值需要安装 vite-plugin-html 插件，然后通过 ejs 模板注入变量。 https://juejin.cn/post/6988704825450397709
// import { createHtmlPlugin } from 'vite-plugin-html'
import { viteCommonjs } from "@originjs/vite-plugin-commonjs";
//旧版浏览器支持插件，自动生成旧版块和相应的ES语言功能polyfills。需要安装 npm add -D terser
import legacy from "@vitejs/plugin-legacy";
// VueDevtools():  https://v2ex.com/t/939478#reply9
import VueDevtools from 'vite-plugin-vue-devtools'
// https://vitejs.dev/config/
export default defineConfig({
  // 静态资源基础路径 base: './' || '',
  // base: process.env.NODE_ENV === 'production' ? './' : '/',
  base: "/",
  resolve: {
    alias: [
      {
        find: /^~/,
        replacement: "",
      },
      {
        find: "@", ///vite默认不支持@别名，但可以通过配置启用此特性
        replacement: path.resolve(__dirname, "src"),
      },
      //解决i18n警告
      {
        find: "vue-i18n",
        replacement: "vue-i18n/dist/vue-i18n.cjs.js",
      },
    ],
    extensions: [".mjs", ".js", ".ts", ".jsx", ".tsx", ".json", ".vue"],
  },
  plugins: [
    vue(),
    vueJsx(),
    viteCommonjs(),
    envCompatible(),
    // injectHtml(),
    VueDevtools(),
    legacy({
      targets: ["defaults", "not IE 11"],
    }),
  ],
  server: {
    strictPort: false, //设置为 true 时，如果端口已被使用，则直接退出，而不会再进行后续端口的尝试。
    open: "/#/", //开发服务器启动时，自动在浏览器中打开应用程序。(false  或 'index.html')
    port: 4080, //开发服务器端口。如果设端口已被使用，Vite 将自动尝试下一个可用端口。
    host: "0.0.0.0", //为开发服务器指定 ip 地址。 设置为 0.0.0.0 或 true 会监听所有地址，包括局域网和公共地址。
    //https://zxuqian.cn/vite-proxy-config/
    proxy: {
      //正则表达式：https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Guide/Regular_Expressions
      "^/api/(server_info|raw|book_infos|get_book|get_file|config.toml|comigo.reg|qrcode.png|redirect|upload|form|raw|login|logout|form|group_info_filter|parent_book_info).*":
        {
          target: "http://127.0.0.1:1234/",
          // 允许跨域:是否改写 origin，设置为 true 之后，就会把请求 API header 中的 origin，改成跟 target 里边的域名一样
          changeOrigin: true,
        },
      //在线测试正则表达式，测试的时候不需要表示开头的 ^
      // https://tool.chinaz.com/regex
      //  127.0.0.1:1234/api/ws/api/ws
      "^/api/ws.*": {
        target: "ws://127.0.0.1:1234/api/ws",
        // 允许跨域:是否改写 origin，设置为 true 之后，就会把请求 API header 中的 origin，改成跟 target 里边的域名一样
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, "").replace(/^\/ws/, ""),
      },
      "/favicon.ico": {
        target: "http://127.0.0.1:1234/",
      },
      "/raw": {
        target: "http://127.0.0.1:1234/",
        changeOrigin: true,
      },
      "/asssets": {
        target: "http://127.0.0.1:1234/",
        changeOrigin: true,
      },
      "/images": {
        target: "http://127.0.0.1:1234/",
        changeOrigin: true,
      },
      "/index.html": {
        ignorePath: true, //不知道有没有用╮(￣▽￣")╭
      },
    },
  },
  build: {
    // target: 'esnext',//设置最终构建的浏览器兼容目标。默认值是一个 Vite 特有的值——'modules'，这是指 支持原生 ES 模块的浏览器。另一个特殊值是 “esnext” —— 即假设有原生动态导入支持，并且将会转译得尽可能小。
    // minify: 'terser',//如果 build.minify 选项为 'terser'， 'esnext' 将会强制降级为 'es2019' need to install it (npm add -D terser)
    outDir: "../routers/vue_static",
    emptyOutDir: true, //清除目标目录：
    chunkSizeWarningLimit: 1500, //规定触发警告的 chunk 大小。（以 kbs 为单位）https://cn.vitejs.dev/config/#build-emptyoutdir
  },
});
