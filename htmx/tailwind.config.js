/** @type {import('tailwindcss').Config} */
module.exports = {
  // 监视这些文件的变化，然后编译CSS
  content: [
    '**/*.{html,templ}',
    './node_modules/flowbite/**/*.js',
    './node_modules/@alpinejs/persist/dist/cdn.js',
    './router/assets/scripts.js',
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('flowbite/plugin')
  ],
}
