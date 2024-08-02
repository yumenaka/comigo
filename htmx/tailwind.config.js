/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['**/*.{html,templ}'],
  // theme: {
  //   extend: {},
  // },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('daisyui'),
  ],
  daisyui: {
    themes: [
      "light",
      "dark",
      "dracula",
      "retro",
      "cupcake",
      "cyberpunk",
      "valentine",
      "aqua",
      "coffee",
      "nord",
    ],
  },
}
