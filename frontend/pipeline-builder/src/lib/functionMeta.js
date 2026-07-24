// Shared visual metadata for function/node types.
// NOTE: Tailwind classes are written as full literal strings so the JIT compiler
// picks them up (do not build class names dynamically from fragments).

import i18n from '../i18n';

const STYLE = {
  connector: { border: 'border-blue-500/60', card: 'border-blue-500/30 bg-blue-500/10', chip: 'bg-blue-500/15 text-blue-300', icon: '🔌', labelKey: 'category.connectors' },
  transform: { border: 'border-orange-500/60', card: 'border-orange-500/30 bg-orange-500/10', chip: 'bg-orange-500/15 text-orange-300', icon: '⚙️', labelKey: 'category.transforms' },
  calculate: { border: 'border-green-500/60', card: 'border-green-500/30 bg-green-500/10', chip: 'bg-green-500/15 text-green-300', icon: '📊', labelKey: 'category.calculates' },
  condition: { border: 'border-purple-500/60', card: 'border-purple-500/30 bg-purple-500/10', chip: 'bg-purple-500/15 text-purple-300', icon: '🚦', labelKey: 'category.conditions' },
  output: { border: 'border-red-500/60', card: 'border-red-500/30 bg-red-500/10', chip: 'bg-red-500/15 text-red-300', icon: '📤', labelKey: 'category.outputs' },
};

const FALLBACK = { border: 'border-gray-500/60', card: 'border-gray-500/30 bg-gray-500/10', chip: 'bg-gray-500/15 text-gray-300', icon: '📦', labelKey: 'category.other' };

export function typeStyle(type) {
  const s = STYLE[type] || FALLBACK;
  return { ...s, label: i18n.t(s.labelKey) };
}

export function typeIcon(type) {
  return typeStyle(type).icon;
}

export function getCategory(type) {
  return typeStyle(type).label;
}
