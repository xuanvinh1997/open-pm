import type { Config } from 'tailwindcss'
import forms from '@tailwindcss/forms'
import typography from '@tailwindcss/typography'

export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        brand: {
          50: '#EBF5FF',
          100: '#D6EBFF',
          200: '#ADD6FF',
          300: '#85C2FF',
          400: '#5CADFF',
          500: '#3B93E6',
          600: '#2B7ACC',
          700: '#006399',
          800: '#004D75',
          900: '#003752',
          950: '#002133',
        },
        custom: {
          background: {
            100: 'var(--color-background-100)',
            90: 'var(--color-background-90)',
            80: 'var(--color-background-80)',
          },
          text: {
            100: 'var(--color-text-100)',
            200: 'var(--color-text-200)',
            300: 'var(--color-text-300)',
            400: 'var(--color-text-400)',
          },
          border: {
            100: 'var(--color-border-100)',
            200: 'var(--color-border-200)',
            300: 'var(--color-border-300)',
          },
          sidebar: {
            background: 'var(--color-background-80)',
            border: 'var(--color-border-200)',
          },
        },
        priority: {
          urgent: '#EF4444',
          high: '#F97316',
          medium: '#EAB308',
          low: '#3B82F6',
          none: '#A3A3A3',
        },
        state: {
          backlog: '#A3A3A3',
          unstarted: '#80838D',
          started: '#F59E0B',
          completed: '#16A34A',
          cancelled: '#EF4444',
          triage: '#8B5CF6',
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
      fontSize: {
        '2xs': ['0.625rem', { lineHeight: '0.875rem' }],
      },
      borderRadius: {
        DEFAULT: '0.375rem',
        lg: '0.75rem',
        xl: '1rem',
        '2xl': '1.25rem',
        '3xl': '1.5rem',
      },
      boxShadow: {
        'custom-xs': '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        'custom-sm': '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1)',
        'custom-md': '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1)',
        'custom-lg': '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)',
      },
      spacing: {
        '4.5': '1.125rem',
        '13': '3.25rem',
        '15': '3.75rem',
      },
      animation: {
        'slide-in-up': 'slideInUp 0.15s ease-out',
        'fade-in': 'fadeIn 0.15s ease-out',
      },
      keyframes: {
        slideInUp: {
          '0%': { transform: 'translateY(4px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
      },
    },
  },
  plugins: [forms, typography],
} satisfies Config
