// Default config seeded when a function node is added, so its configuration
// fields (and the relevant pickers) appear immediately. Mirrors the parameters
// documented in docs/mindset-modif-update.md.
const FUNCTION_DEFAULTS = {
  state_machine: { machine_id: '' },
  threshold: { min: 30, max: 180 },
  calculate_cost: { hourly_rate: 85, currency: 'EUR' },
  calculate_duration: {},
  filter: { field: 'value', operator: 'gt', value: 0 },
  add_to_dashboard: { label: 'Mon widget', kind: 'value' },
};

export function defaultFunctionConfig(name) {
  return { ...(FUNCTION_DEFAULTS[name] || {}) };
}
