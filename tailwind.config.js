/** @type {import('tailwindcss').Config} */
module.exports = {
  // 监视这些文件的变化，然后编译CSS
  content: [
    './**/*.{html,templ}',
    '!./node_modules/**/*',
    './router/script/scripts.js',
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms')({ strategy: 'class' }),
    require('@tailwindcss/typography'),
    require('flowbite/plugin')
  ],
}
