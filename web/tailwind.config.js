//配置文档：https://www.tailwindcss.cn/docs/configuration
module.exports = {
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],//配置 purge 选项指定所有的 pages 和 components 文件，使得 Tailwind 可以在生产构建中对未使用的样式进行摇树优化。
  darkMode: "media", //false or 'media' or 'class'
  content: [],
  theme: {
    extend: {},
  },
  plugins: [],
}
