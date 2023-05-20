/** @type {import('tailwindcss').Config} */
export default {
  mode:'jit',
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: 'class',
  theme: {
    extend:
    {
      fontFamily: {
        inter: ['Inter', 'sans-serif']
      },
      gridTemplateAreas: {
        'sidebar': ['"left-panel content right-panel"']
      },
    },
  },
  plugins: [],
  purge: ['./src/**/*.{html,js,svelte,ts}']
};