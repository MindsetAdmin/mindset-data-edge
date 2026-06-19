// Shared visual metadata for function/node types.
// NOTE: Tailwind classes are written as full literal strings so the JIT compiler
// picks them up (do not build class names dynamically from fragments).

const STYLE = {
  connector: { border: 'border-blue-500/60', card: 'border-blue-500/30 bg-blue-500/10', chip: 'bg-blue-500/15 text-blue-300', icon: '🔌', label: 'Connecteurs' },
  transform: { border: 'border-orange-500/60', card: 'border-orange-500/30 bg-orange-500/10', chip: 'bg-orange-500/15 text-orange-300', icon: '⚙️', label: 'Transformations' },
  calculate: { border: 'border-green-500/60', card: 'border-green-500/30 bg-green-500/10', chip: 'bg-green-500/15 text-green-300', icon: '📊', label: 'Calculs' },
  condition: { border: 'border-purple-500/60', card: 'border-purple-500/30 bg-purple-500/10', chip: 'bg-purple-500/15 text-purple-300', icon: '🚦', label: 'Conditions' },
  output: { border: 'border-red-500/60', card: 'border-red-500/30 bg-red-500/10', chip: 'bg-red-500/15 text-red-300', icon: '📤', label: 'Sorties' },
};

const FALLBACK = { border: 'border-gray-500/60', card: 'border-gray-500/30 bg-gray-500/10', chip: 'bg-gray-500/15 text-gray-300', icon: '📦', label: 'Autres' };

export function typeStyle(type) {
  return STYLE[type] || FALLBACK;
}

export function typeIcon(type) {
  return typeStyle(type).icon;
}

export function getCategory(type) {
  return typeStyle(type).label;
}
