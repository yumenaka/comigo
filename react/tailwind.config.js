/** @type {import('tailwindcss').Config} */
export default {
  // darkMode: "class", // or 'media'
  mode: 'jit',
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    // 需要自定义tailwindcss，做颜色模板的时候，可以尝试这些：
    // https://github.com/crswll/tailwindcss-theme-swapper
    // https://github.com/innocenzi/tailwindcss-theming
    // https://github.com/aniftyco/awesome-tailwindcss
    extend: {
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
  ],
}

