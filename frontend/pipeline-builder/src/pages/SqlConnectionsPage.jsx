import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { fetchConnections, createConnection, testConnection, deleteConnection, fetchConnectionDatabases } from '../api/client';

const DEFAULT_FORM = {
  id: '',
  name: '',
  driver: 'mysql', // V1a supports only mysql
  host: '',
  port: 3306,
  database: '',
  username: '',
  password_env: '',
  tls: 'false',
  read_timeout_seconds: 30,
  write_timeout_seconds: 10,
  max_open_conns: 5,
  max_idle_conns: 2,
  conn_max_lifetime_seconds: 300,
};

// Status badge for a connection row. "unknown" until a Test (or a sql_query
// execution) has actually opened the pool at least once this server run.
function StatusBadge({ status, readOnly }) {
  const { t } = useTranslation();
  if (status !== 'ok') {
    return (
      <span className="text-xs px-2 py-0.5 rounded-full border font-mono bg-dark-700 text-dark-300 border-dark-600">
        {t('sqlConnections.notTested')}
      </span>
    );
  }
  return readOnly ? (
    <span className="text-xs px-2 py-0.5 rounded-full border font-mono bg-green-500/20 text-green-400 border-green-500/40">
      {t('sqlConnections.readOnly')}
    </span>
  ) : (
    <span className="text-xs px-2 py-0.5 rounded-full border font-mono bg-yellow-500/20 text-yellow-400 border-yellow-500/40">
      ⚠️ {t('sqlConnections.writeAccess')}
    </span>
  );
}

// SQL connections page — reached from the MySQL tile on /connectors. List,
// create, test, and delete connections used by the sql_query connector
// (internal/connections). Extracted from the old combined ConnectionsPage
// (2026-07-21) when that page became the pure connector gallery.
export default function SqlConnectionsPage() {
  const { t } = useTranslation();
  const [connections, setConnections] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [form, setForm] = useState(DEFAULT_FORM);
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [saving, setSaving] = useState(false);
  const [testResults, setTestResults] = useState({}); // id -> {ok, latency_ms, read_only, error}
  const [testingId, setTestingId] = useState(null);
  const [databases, setDatabases] = useState({}); // id -> [{ name, tables: [{ name, columns }] }]
  const [expandedTables, setExpandedTables] = useState({}); // "id/dbName/tableName" -> bool

  const load = async () => {
    setLoading(true);
    try {
      const data = await fetchConnections();
      setConnections(data.connections || []);
      setError(null);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const set = (k, v) => setForm((f) => ({ ...f, [k]: v }));

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    setError(null);
    try {
      const cfg = {
        ...form,
        port: Number(form.port) || 3306,
        read_timeout_seconds: Number(form.read_timeout_seconds) || 30,
        write_timeout_seconds: Number(form.write_timeout_seconds) || 10,
        max_open_conns: Number(form.max_open_conns) || 5,
        max_idle_conns: Number(form.max_idle_conns) || 2,
        conn_max_lifetime_seconds: Number(form.conn_max_lifetime_seconds) || 300,
      };
      await createConnection(cfg);
      setForm(DEFAULT_FORM);
      await load();
    } catch (err) {
      setError(err.message);
    } finally {
      setSaving(false);
    }
  };

  const handleTest = async (id) => {
    setTestingId(id);
    try {
      const result = await testConnection(id);
      setTestResults((r) => ({ ...r, [id]: result }));
      await load(); // refresh read_only/status badges from the registry
      // Connecting successfully doubles as "show me what's in there" — browse
      // every database + table visible to this connection right away, no
      // separate click needed. Best-effort: a browse failure shouldn't hide
      // the test result that already succeeded.
      if (result.ok) {
        try {
          const data = await fetchConnectionDatabases(id);
          setDatabases((d) => ({ ...d, [id]: data.databases || [] }));
        } catch {
          /* browse is a bonus, not required for Test to have succeeded */
        }
      }
    } catch (err) {
      setTestResults((r) => ({ ...r, [id]: { ok: false, error: err.message } }));
    } finally {
      setTestingId(null);
    }
  };

  const toggleTable = (key) =>
    setExpandedTables((e) => ({ ...e, [key]: !e[key] }));

  const handleDelete = async (id) => {
    try {
      await deleteConnection(id);
      await load();
    } catch (err) {
      setError(err.message);
    }
  };

  const field = (label, key, type = 'text', opts = null) => (
    <label className="block mb-3">
      <span className="text-xs text-dark-400 block mb-1">{label}</span>
      {opts ? (
        <select
          value={form[key]}
          onChange={(e) => set(key, e.target.value)}
          className="w-full bg-dark-800 border border-dark-600 rounded-md px-2 py-1.5 text-sm text-white"
        >
          {opts.map((o) => (
            <option key={o} value={o}>{o}</option>
          ))}
        </select>
      ) : (
        <input
          type={type}
          value={form[key]}
          onChange={(e) => set(key, e.target.value)}
          required={['id', 'name', 'host', 'database', 'username', 'password_env'].includes(key)}
          className="w-full bg-dark-800 border border-dark-600 rounded-md px-2 py-1.5 text-sm text-white font-mono"
        />
      )}
    </label>
  );

  return (
    <div className="h-full overflow-y-auto p-6">
      <div className="max-w-4xl mx-auto">
        <h2 className="text-xl font-semibold text-white mb-1">🐬 {t('sqlConnections.title')}</h2>
        <p className="text-dark-400 text-sm mb-6">
          {t('sqlConnections.subtitle')}
        </p>

        {error && (
          <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-3 text-red-400 text-sm mb-4">
            ❌ {error}
          </div>
        )}

        <div className="text-xs text-dark-500 mb-2">{t('sqlConnections.countConfigured', { count: connections.length })}</div>
        <div className="space-y-3 mb-8">
          {loading && <p className="text-dark-400 text-sm">{t('common.loading')}</p>}
          {!loading && connections.length === 0 && (
            <p className="text-dark-500 text-sm">{t('sqlConnections.none')}</p>
          )}
          {connections.map((c) => {
            const result = testResults[c.id];
            return (
              <div key={c.id} className="border border-dark-700 bg-dark-900 rounded-lg p-4">
                <div className="flex items-center gap-3">
                  <span className="text-xl">🗄️</span>
                  <div className="min-w-0 flex-1">
                    <div className="font-medium text-white">{c.name} <span className="text-dark-500 font-mono text-xs">({c.id})</span></div>
                    <div className="text-xs text-dark-400 font-mono">
                      {c.driver}://{c.username}@{c.host}:{c.port}/{c.database} {c.tls !== 'false' ? `(tls:${c.tls})` : ''}
                    </div>
                  </div>
                  <StatusBadge status={c.status} readOnly={c.read_only} />
                  <button
                    onClick={() => handleTest(c.id)}
                    disabled={testingId === c.id}
                    className="bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white text-sm px-3 py-1.5 rounded-md transition"
                  >
                    {testingId === c.id ? '…' : t('sqlConnections.connect')}
                  </button>
                  <button
                    onClick={() => handleDelete(c.id)}
                    className="bg-dark-700 hover:bg-red-500/30 text-dark-300 hover:text-red-300 text-sm px-3 py-1.5 rounded-md transition"
                  >
                    {t('common.delete')}
                  </button>
                </div>
                {result && (
                  <div className={`text-xs mt-2 ${result.ok ? 'text-green-400' : 'text-red-400'}`}>
                    {result.ok
                      ? `✅ OK / ${result.latency_ms}ms / ${result.read_only ? t('sqlConnections.readOnly') : '⚠️ ' + t('sqlConnections.writeAccess')}`
                      : `❌ ${result.error}`}
                  </div>
                )}
                {/* Databases + tables visible to this connection — populated
                    automatically once Connecter succeeds. */}
                {databases[c.id] && (
                  <div className="mt-3 border-t border-dark-700 pt-3 space-y-2">
                    {databases[c.id].length === 0 ? (
                      <p className="text-xs text-dark-500">{t('sqlConnections.noDatabaseVisible')}</p>
                    ) : (
                      databases[c.id].map((db) => (
                        <div key={db.name}>
                          <div className="text-xs font-medium text-dark-200 mb-1">
                            🗃️ {db.name} <span className="text-dark-500">{t('sqlConnections.tableCount', { count: db.tables.length })}</span>
                          </div>
                          <div className="pl-4 space-y-0.5">
                            {db.tables.map((t2) => {
                              const key = `${c.id}/${db.name}/${t2.name}`;
                              const open = !!expandedTables[key];
                              return (
                                <div key={t2.name}>
                                  <button
                                    onClick={() => toggleTable(key)}
                                    className="text-xs text-dark-300 hover:text-white font-mono flex items-center gap-1"
                                  >
                                    <span className="text-dark-500">{open ? '▾' : '▸'}</span>
                                    {t2.name}
                                    <span className="text-dark-500">({t2.columns.length})</span>
                                  </button>
                                  {open && (
                                    <div className="pl-5 text-[11px] text-dark-500 font-mono">
                                      {t2.columns.map((col) => (
                                        <div key={col.name}>
                                          {col.is_key && <span className="text-status-warn">🔑 </span>}
                                          {col.name} <span className="text-dark-600">{col.data_type}</span>
                                        </div>
                                      ))}
                                    </div>
                                  )}
                                </div>
                              );
                            })}
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                )}
              </div>
            );
          })}
        </div>

        <h3 className="text-sm font-semibold text-dark-300 uppercase tracking-wider mb-2">{t('sqlConnections.newConnection')}</h3>
        <form onSubmit={handleSave} className="border border-blue-500/30 bg-blue-500/5 rounded-lg p-4">
          <div className="grid grid-cols-2 gap-2">
            {field(t('sqlConnections.fieldId'), 'id')}
            {field(t('sqlConnections.fieldName'), 'name')}
          </div>
          <div className="grid grid-cols-2 gap-2">
            {field(t('sqlConnections.fieldHost'), 'host')}
            {field(t('sqlConnections.fieldPort'), 'port', 'number')}
          </div>
          <div className="grid grid-cols-2 gap-2">
            {field(t('sqlConnections.fieldDatabase'), 'database')}
            {field(t('sqlConnections.fieldUsername'), 'username')}
          </div>
          <div className="grid grid-cols-2 gap-2">
            {field(t('sqlConnections.fieldPasswordEnv'), 'password_env')}
            {field('TLS', 'tls', 'text', ['false', 'true', 'skip-verify'])}
          </div>

          <button
            type="button"
            onClick={() => setShowAdvanced((v) => !v)}
            className="text-xs text-dark-400 hover:text-white mb-2"
          >
            {showAdvanced ? '▾' : '▸'} {t('sqlConnections.advancedOptions')}
          </button>
          {showAdvanced && (
            <div className="grid grid-cols-2 gap-2 mb-2">
              {field(t('sqlConnections.fieldReadTimeout'), 'read_timeout_seconds', 'number')}
              {field(t('sqlConnections.fieldWriteTimeout'), 'write_timeout_seconds', 'number')}
              {field(t('sqlConnections.fieldMaxOpenConns'), 'max_open_conns', 'number')}
              {field(t('sqlConnections.fieldMaxIdleConns'), 'max_idle_conns', 'number')}
              {field(t('sqlConnections.fieldMaxLifetime'), 'conn_max_lifetime_seconds', 'number')}
            </div>
          )}

          <button
            type="submit"
            disabled={saving}
            className="mt-2 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white text-sm px-3 py-1.5 rounded-md transition"
          >
            {saving ? t('common.saving') : t('common.save')}
          </button>
        </form>
      </div>
    </div>
  );
}
