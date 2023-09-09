import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  base: "/admin",
  plugins: [react()],
  build: {
    outDir: "../routers/admin",
    emptyOutDir: true, //清除目标目录：aaaaaaaaacece
    chunkSizeWarningLimit: 1500, //规定触发警告的 chunk 大小。（以 kbs 为单位）https://cn.vitejs.dev/config/#build-emptyoutdir
  },
  server: {
    strictPort: false, //设置为 true 时，如果端口已被使用，则直接退出，而不会再进行后续端口的尝试。
    open: "/admin", //开发服务器启动时，自动在浏览器中打开应用程序。(false  或 'index.html')
    port: 4090, //开发服务器端口。如果设端口已被使用，Vite 将自动尝试下一个可用端口。
    host: "0.0.0.0", //为开发服务器指定 ip 地址。 设置为 0.0.0.0 或 true 会监听所有地址，包括局域网和公共地址。
    //https://zxuqian.cn/vite-proxy-config/
    proxy: {
      //正则表达式：https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Guide/Regular_Expressions
      "^/api/(getstatus|getlist|getbook|getfile|config.toml|config.json|config_update|post_config|comigo.reg|qrcode.png|redirect|upload|form|raw|login|logout|form).*":
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
});
