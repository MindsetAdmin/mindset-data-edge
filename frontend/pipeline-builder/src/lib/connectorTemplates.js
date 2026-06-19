// Default config templates for connectors, so picking a connector pre-fills
// sensible fields. Keys/values mirror the existing pipeline YAML examples
// (config/pipelines/*.yaml) and the handler config lookups in
// internal/functions/connectors/*.go.

export const CONNECTOR_TEMPLATES = {
  mqtt_subscribe: { topic: 'mindset/events/status-change', qos: 1 },
  opcua_read: { node_id: 'ns=3;i=1001', timeout: 5000 },
  modbus_read: { host: '127.0.0.1', port: 502, unit_id: 1, register: 40001 },
  sql_query: { dsn: 'postgres://localhost/mindset', query: 'SELECT * FROM events' },
};

// Trigger "type" string that goes in pipeline.Trigger.type for a given connector.
export const TRIGGER_TYPE_BY_CONNECTOR = {
  mqtt_subscribe: 'mqtt',
  opcua_read: 'opcua',
  modbus_read: 'modbus',
  sql_query: 'sql',
};

export function defaultConfigFor(fnName) {
  return { ...(CONNECTOR_TEMPLATES[fnName] || {}) };
}

export function triggerTypeFor(fnName) {
  return TRIGGER_TYPE_BY_CONNECTOR[fnName] || 'mqtt';
}
