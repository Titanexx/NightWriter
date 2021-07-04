const colors = require('tailwindcss/colors')

module.exports = {
	purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
	darkMode: 'class', // or 'media' or 'class'
	theme: {
		colors: {
            transparent:'transparent',
            current:'currentColor',
			gray:{
				0: 'var(--color-gray-0)',
				100: 'var(--color-gray-100)',
				200: 'var(--color-gray-200)',
				300: 'var(--color-gray-300)',
				400: 'var(--color-gray-400)',
				500: 'var(--color-gray-500)',
				600: 'var(--color-gray-600)',
				700: 'var(--color-gray-700)',
				800: 'var(--color-gray-800)',
				900: 'var(--color-gray-900)',
				1000: 'var(--color-gray-1000)',
				1100: 'var(--color-gray-1100)',
			},
            black: 'var(--color-black)',
            white: 'var(--color-white)',
            smoke: {
                0: 'rgba(0, 0, 0, 0)',
                100: 'rgba(0, 0, 0, 0.1)',
                250: 'rgba(0, 0, 0, 0.25)',
                400: 'rgba(0, 0, 0, 0.4)',
                500: 'rgba(0, 0, 0, 0.5)',
                600: 'rgba(0, 0, 0, 0.6)',
                750: 'rgba(0, 0, 0, 0.75)',
                900: 'rgba(0, 0, 0, 0.9)',
                1000: 'rgba(0, 0, 0, 1)',
            },
            green:colors.emerald,
            yellow:colors.amber,
            red:colors.red,
            indigo:colors.indigo,
		},
		extend: {
			textColor: {
				primary: 'var(--color-text-primary)',
				secondary: 'var(--color-text-secondary)',
			},
			backgroundColor: {
				primary: 'var(--color-bg-primary)',
				secondary: 'var(--color-bg-secondary)',
			},
			fontFamily: {
				display: 'var(--font-display)',
				body: 'var(--font-body)',
			},
            animation: {
                'btn-text': 'max-width 1s linear',
            }
		},
	},
	variants: {
		extend: {},
	},
	plugins: [
		require('@tailwindcss/forms'),
		// ...
	],
}
