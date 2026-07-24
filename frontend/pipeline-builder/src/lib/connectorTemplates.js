// Default config templates for connectors, so picking a connector pre-fills
// sensible fields. Keys/values mirror the existing pipeline YAML examples
// (config/pipelines/*.yaml) and the handler config lookups in
// internal/functions/connectors/*.go.

export const CONNECTOR_TEMPLATES = {
  mqtt_subscribe: { topic: 'mindset/events/status-change', qos: 1 },
  opcua_read: { node_id: 'ns=3;i=1001', timeout: 5000 },
  modbus_read: { host: '127.0.0.1', port: 502, unit_id: 1, register: 40001 },
  sql_query: { connection_id: '', query: 'SELECT * FROM work_orders LIMIT 10', timeout_seconds: 30, limit: 100 },
};

// Unlike opcua_read/mqtt_subscribe/modbus_read (trigger-only — they start a
// pipeline), sql_query is an enrichment connector: it's meant to run mid-pipeline
// too (e.g. mqtt_subscribe -> sql_query, config/pipelines/examples/of_enrichment.yaml
// — the result auto-publishes, no mqtt_publish node needed, see Entry 119). The
// pipeline engine already executes any node type in dependency order — this is
// purely what BuilderPage's Palette/picker allow dragging into the CŒUR band.
export const MID_PIPELINE_CONNECTORS = new Set(['sql_query']);

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
