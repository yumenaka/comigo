/** @type {import('tailwindcss').Config} */
export default {
  // darkMode: "class", // or 'media'
  mode: 'jit',
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    // https://cloud.tencent.com/developer/article/1967312
    extend: {
      // colors: {
      //   amber: colors.amber,
      //   lime: colors.lime,
      //   rose: colors.rose,
      //   orange: colors.orange,
      // },
    },
    // backgroundColor: {
    //   //utilities like `bg-base` and `bg-primary`
    //   base: 'var(--color-base)',
    //   'off-base': 'var(--color-off-base)',
    //   primary: 'var(--color-primary)',
    //   secondary: 'var(--color-secondary)',
    //   muted: 'var(--color-text-muted)',
    // },
    // textColor: {
    //   //like `text-base` and `text-primary`
    //   base: 'var(--color-text-base)',
    //   muted: 'var(--color-text-muted)',
    //   'muted-hover': 'var(--color-text-muted-hover)',
    //   primary: 'var(--color-primary)',
    //   secondary: 'var(--color-secondary)',
    // },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@headlessui/tailwindcss')
  ],
}

