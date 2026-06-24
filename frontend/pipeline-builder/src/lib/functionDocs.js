// Human-friendly documentation per function: a description and, per config field,
// a label / help text / example. Drives the guided configuration panels.

export const FUNCTION_DOCS = {
  mqtt_subscribe: {
    description: "S'abonne à un topic MQTT — c'est le déclencheur (trigger) du pipeline.",
    fields: {
      topic: { label: 'Topic', help: 'Topic MQTT à écouter.', example: 'mindset/raw/ns=3;i=1014' },
      qos: { label: 'QoS', help: 'Qualité de service MQTT (0, 1 ou 2).', example: '1' },
    },
  },
  opcua_read: {
    description: 'Lit une ou plusieurs valeurs depuis un serveur OPC-UA.',
    fields: {
      node_id: { label: 'Node ID', help: 'Identifiant du nœud OPC-UA à lire.', example: 'ns=3;i=1011' },
      endpoint: { label: 'Endpoint', help: 'URL du serveur OPC-UA.', example: 'opc.tcp://localhost:53530/...' },
      timeout: { label: 'Timeout (ms)', help: 'Délai max de lecture en millisecondes.', example: '5000' },
    },
  },
  modbus_read: { description: 'Lit depuis un équipement Modbus (démo).', fields: {} },
  sql_query: { description: 'Interroge une base SQL (démo).', fields: {} },
  state_machine: {
    description: 'Détecte les transitions Run ↔ Stop et déclenche sur les arrêts.',
    fields: {
      machine_id: { label: 'Machine', help: 'Machine à surveiller (depuis la découverte OPC-UA).', example: 'machine1' },
    },
  },
  uns_mapper: {
    description: 'Normalise un tag OPC-UA en topic ISA-95 (mindset/site/...).',
    fields: {
      site_id: { label: 'Site', help: 'Identifiant du site.', example: 'local-test' },
      area: { label: 'Zone', help: 'Zone de production.', example: 'area1' },
    },
  },
  filter: {
    description: 'Garde ou rejette les données selon une condition.',
    fields: {
      field: { label: 'Champ', help: 'Champ sur lequel filtrer.', example: 'value' },
      operator: { label: 'Opérateur', help: 'eq, ne, gt, lt, contains.', example: 'gt' },
      value: { label: 'Valeur', help: 'Valeur de comparaison.', example: '0' },
    },
  },
  calculate_duration: {
    description: "Calcule la durée entre le début et la fin d'un arrêt.",
    fields: {},
  },
  calculate_cost: {
    description: "Calcule le coût en euros d'un arrêt à partir de sa durée.",
    fields: {
      hourly_rate: { label: 'Coût horaire (€/h)', help: 'Coût de la machine par heure.', example: '85' },
      currency: { label: 'Devise', help: 'EUR, USD ou GBP.', example: 'EUR' },
    },
  },
  threshold: {
    description: 'Vérifie si une durée est dans un intervalle (fenêtre micro-arrêt).',
    fields: {
      min: { label: 'Min (s)', help: 'Durée minimale en secondes.', example: '30' },
      max: { label: 'Max (s)', help: 'Durée maximale en secondes.', example: '180' },
    },
  },
  mqtt_publish: {
    description: 'Publie le résultat sur un topic MQTT.',
    fields: {
      topic: { label: 'Topic', help: 'Topic de publication.', example: 'mindset/events/micro-stop' },
      qos: { label: 'QoS', help: 'Qualité de service MQTT.', example: '1' },
    },
  },
  add_to_dashboard: {
    description: "Épingle la donnée ou l'événement sur le Dashboard.",
    fields: {
      label: { label: 'Libellé', help: 'Nom affiché du widget.', example: 'Coût micro-stop' },
      kind: { label: 'Type', help: '"value" pour un nombre, "event" pour un événement.', example: 'value' },
    },
  },
};

export function functionDoc(name) {
  return FUNCTION_DOCS[name] || { description: '', fields: {} };
}

export function fieldDoc(fnName, key) {
  return functionDoc(fnName).fields[key] || { label: key, help: '', example: '' };
}
