/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,jsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: "#E8F7F2",
          100: "#C4EBDD",
          200: "#9FDFC9",
          300: "#7AD3B5",
          400: "#56C7A1",
          500: "#1D9E75",
          600: "#198D68",
          700: "#157B5A",
          800: "#116A4D",
          900: "#0D5840",
        },
        accent: {
          50: "#FBEFEA",
          100: "#F6D8CD",
          200: "#F1C1AF",
          300: "#ECAA92",
          400: "#E79375",
          500: "#D85A30",
          600: "#C24F2A",
          700: "#AB4625",
          800: "#943C20",
          900: "#7D331B",
        },
      },
      borderRadius: {
        card: "12px",
        button: "8px",
      },
      fontFamily: {
        sans: ["Inter", "ui-sans-serif", "system-ui", "sans-serif"],
      },
    },
  },
  plugins: [],
};