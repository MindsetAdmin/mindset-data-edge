// Human-friendly documentation per function: a description and, per config field,
// a label / help text / example. Drives the guided configuration panels.
// Descriptions/labels/help resolve through i18next (not the React hook, since this
// is plain data used outside components) — examples are literal values, not translated.

import i18n from '../i18n';

const FIELD_SHAPE = {
  mqtt_subscribe: {
    topic: { example: 'mindset/raw/ns=3;i=1014' },
    qos: { example: '1' },
  },
  opcua_read: {
    node_id: { example: 'ns=3;i=1011' },
    endpoint: { example: 'opc.tcp://localhost:53530/...' },
    timeout: { example: '5000' },
  },
  modbus_read: {},
  sql_query: {}, // fields are rendered by SqlConfigPanel, not the generic list
  state_machine: {
    machine_id: { example: 'machine1' },
  },
  filter: {
    field: { example: 'value' },
    operator: { example: 'gt' },
    value: { example: '0' },
  },
  calculate_duration: {},
  calculate_cost: {
    hourly_rate: { example: '85' },
    currency: { example: 'EUR' },
  },
  threshold: {
    min: { example: '30' },
    max: { example: '180' },
  },
  add_to_dashboard: {
    label: { example: 'Micro-stop cost' },
    kind: { example: 'value' },
  },
};

export function functionDoc(name) {
  const shape = FIELD_SHAPE[name];
  if (!shape) return { description: '', fields: {} };
  const fields = {};
  for (const key of Object.keys(shape)) {
    fields[key] = {
      label: i18n.t(`functionDocs.${name}.fields.${key}.label`),
      help: i18n.t(`functionDocs.${name}.fields.${key}.help`),
      example: shape[key].example,
    };
  }
  return { description: i18n.t(`functionDocs.${name}.description`), fields };
}

export function fieldDoc(fnName, key) {
  return functionDoc(fnName).fields[key] || { label: key, help: '', example: '' };
}
