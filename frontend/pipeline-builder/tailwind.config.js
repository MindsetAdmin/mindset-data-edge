/** @type {import('tailwindcss').Config} */
// MindSet design tokens (2026-07-01 redesign — see docs/frontend_redesign.md)
// Additive: new tokens live alongside the legacy `dark` scale so untouched
// pages keep working while redesigned pages adopt the new system.
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // ─── MindSet neutrals (Grafana-inspired dark canvas) ───────────
        canvas: '#0A0A0B',           // page background
        panel: '#131316',            // panels / cards
        'panel-alt': '#1B1B1F',      // nested panels, table row hover
        elevated: '#232329',         // dropdowns, popovers

        // ─── Borders (1px only) ───────────────────────────────────────
        'border-subtle': '#2A2A31',
        'border-strong': '#3A3A44',

        // ─── Text (opacity levels for hierarchy) ──────────────────────
        'text-primary': '#E8E8ED',   // main content, 90%
        'text-secondary': '#A8A8B2', // labels, 65%
        'text-tertiary': '#6E6E7A',  // captions, 45%
        'text-muted': '#4A4A55',     // metadata, 30%

        // ─── Brand accent (amber-warm — differentiates from blue-heavy competitors) ───
        accent: {
          DEFAULT: '#E5A445',
          muted: '#7A5620',
        },

        // ─── Semantic states (Grafana palette, muted for dark canvas) ─
        'status-running': '#4ADE80',
        'status-stopped': '#F87171',
        'status-warn': '#FBBF24',
        'status-info': '#60A5FA',
        'status-idle': '#6E6E7A',

        // ─── Legacy dark scale — kept for backward compat during migration ───
        dark: {
          50: '#f8fafc',
          100: '#f1f5f9',
          200: '#e2e8f0',
          300: '#cbd5e1',
          400: '#94a3b8',
          500: '#64748b',
          600: '#475569',
          700: '#334155',
          800: '#1e293b',
          900: '#0f172a',
          950: '#020617',
        },
      },

      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'ui-monospace', 'SFMono-Regular', 'Menlo', 'monospace'],
      },

      // MindSet type scale — 4 sizes only. Explicit px names to avoid
      // clashing with Tailwind's xs/sm/base/lg defaults.
      fontSize: {
        '11': ['11px', { lineHeight: '1.4' }],
        '13': ['13px', { lineHeight: '1.4' }],
        '15': ['15px', { lineHeight: '1.4' }],
        '20': ['20px', { lineHeight: '1.2', letterSpacing: '-0.01em' }],
      },

      transitionDuration: {
        DEFAULT: '150ms',
      },
    },
  },
  plugins: [],
}
