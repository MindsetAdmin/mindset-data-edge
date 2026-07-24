import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

// Connector gallery — big tiles, no grouping, no config preview (2026-07-21).
// Replaces the old /connect + /connections pair: this is purely a showcase;
// clicking a tile navigates straight to wherever that connector is actually
// configured. Set expanded from docs/mindset.md §5 "Protocols & Connectors"
// (the real, ranked roadmap), not invented — `implemented: true` connectors
// are wired up today; the rest are shown honestly as not-yet-built, never
// linked anywhere fake.
const CONNECTOR_TILES = [
  // Implemented today
  { name: 'opcua_read', label: 'OPC-UA', icon: '🛰️', implemented: true, to: '/connect/opcua' },
  { name: 'sql_query', label: 'MySQL', icon: '🐬', implemented: true, to: '/connectors/sql' },
  { name: 'mqtt_subscribe', label: 'MQTT', icon: '📶', implemented: true, to: '/connectors/mqtt' },
  { name: 'modbus_read', label: 'Modbus TCP', icon: '🔧', implemented: true, to: null }, // registered but a metadata-only stub — errors if run
  // V1 roadmap — SQL multi-dialect (docs/mindset.md §5, docs/decisions.md)
  { name: 'postgresql', label: 'PostgreSQL', icon: '🐘', implemented: false, to: null },
  { name: 'mssql', label: 'MSSQL', icon: '🗃️', implemented: false, to: null },
  // V1.5 roadmap (docs/mindset.md §5)
  { name: 's7', label: 'Siemens S7', icon: '🏭', implemented: false, to: null },
  { name: 'rest_api', label: 'REST API', icon: '🌐', implemented: false, to: null },
  { name: 'ftp', label: 'FTP / SFTP', icon: '📁', implemented: false, to: null },
  // V2+ / other ranked protocols (docs/mindset.md §5)
  { name: 'ignition', label: 'Ignition', icon: '🧩', implemented: false, to: null },
  { name: 'influxdb', label: 'InfluxDB', icon: '📈', implemented: false, to: null },
  { name: 'mongodb', label: 'MongoDB', icon: '🍃', implemented: false, to: null },
];

export default function ConnectorsPage() {
  const navigate = useNavigate();
  const { t } = useTranslation();

  return (
    <div className="h-full overflow-y-auto p-8">
      <div className="max-w-5xl mx-auto">
        <h2 className="text-2xl font-semibold text-white mb-1">🔌 {t('connectors.title')}</h2>
        <p className="text-dark-400 text-sm mb-8">
          {t('connectors.subtitle')}
        </p>

        <div className="grid grid-cols-3 sm:grid-cols-4 gap-5">
          {CONNECTOR_TILES.map((tile) => (
            <button
              key={tile.name}
              onClick={() => tile.to && navigate(tile.to)}
              disabled={!tile.to}
              title={tile.implemented ? tile.label : `${tile.label} — ${t('connectors.comingSoon')}`}
              className={`flex flex-col items-center justify-center gap-2.5 rounded-xl p-6 border transition ${
                tile.to
                  ? 'border-dark-700 bg-dark-900 hover:border-blue-500/60 hover:bg-dark-800 hover:-translate-y-0.5 cursor-pointer'
                  : 'border-dark-800 bg-dark-900/50 opacity-40 cursor-not-allowed'
              }`}
            >
              <span className="text-5xl leading-none">{tile.icon}</span>
              <span className="text-sm text-dark-200 font-medium text-center leading-tight">{tile.label}</span>
              {!tile.implemented && (
                <span className="text-[10px] text-dark-500 uppercase tracking-wide">{t('connectors.comingSoon')}</span>
              )}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
