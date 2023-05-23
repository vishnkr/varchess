/** @type {import('tailwindcss').Config} */
export default {
  mode:'jit',
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: 'class',
  theme: {
    extend:
    {
      opacity: {
        '10': '0.1',
      },
      fontFamily: {
        inter: ['Inter', 'sans-serif'],
      },
      gridTemplateAreas: {
        'sidebar': ['"left-panel content right-panel"']
      },
    },
  },
  plugins: [require('@tailwindcss/typography'),],
  purge: ['./src/**/*.{html,js,svelte,ts}']
};