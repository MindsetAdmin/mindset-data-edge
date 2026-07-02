import { useState, useEffect } from 'react';
import { opcuaConnect, opcuaDisconnect, opcuaStatus, fetchConfig } from '../api/client';

const SECURITY_MODES = ['None', 'Sign', 'SignAndEncrypt'];
const SECURITY_POLICIES = ['None', 'Basic256Sha256', 'Basic256', 'Basic128Rsa15'];

// Status badge colors per connection state.
function StatusBadge({ status }) {
  const map = {
    connected: 'bg-green-500/20 text-green-400 border-green-500/40',
    connecting: 'bg-yellow-500/20 text-yellow-400 border-yellow-500/40',
    error: 'bg-red-500/20 text-red-400 border-red-500/40',
    disconnected: 'bg-dark-700 text-dark-300 border-dark-600',
  };
  const cls = map[status] || map.disconnected;
  return (
    <span className={`text-xs px-2 py-0.5 rounded-full border font-mono ${cls}`}>
      {status || 'disconnected'}
    </span>
  );
}

// OPC-UA connection form. Connecting doubles as "Test Connection": a successful
// connect returns connected status; a failure surfaces the server error inline.
export default function OpcuaConnectionPanel({ onConnected }) {
  const [form, setForm] = useState({
    endpoint: '',
    security_mode: 'None',
    security_policy: 'None',
    username: '',
    password: '',
    session_timeout: 300, // generous so the session survives the browse→select gap
  });
  const [status, setStatus] = useState({ status: 'disconnected', connected: false });
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState(null);

  // Prefill from agent.yaml defaults + reflect any existing session.
  useEffect(() => {
    (async () => {
      try {
        const cfg = await fetchConfig();
        if (cfg?.opcua) {
          setForm((f) => ({
            ...f,
            endpoint: cfg.opcua.endpoint || f.endpoint,
            security_mode: cfg.opcua.security_mode || f.security_mode,
            security_policy: cfg.opcua.security_policy || f.security_policy,
          }));
        }
      } catch {
        /* config is best-effort prefill */
      }
      try {
        const st = await opcuaStatus();
        setStatus(st);
        if (st.connected) onConnected?.(st);
      } catch {
        /* ignore */
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const set = (k, v) => setForm((f) => ({ ...f, [k]: v }));

  const handleConnect = async () => {
    setBusy(true);
    setError(null);
    try {
      const st = await opcuaConnect({ ...form, session_timeout: Number(form.session_timeout) || 60 });
      setStatus(st);
      if (st.connected) onConnected?.(st);
      else setError(st.error || 'Connexion échouée.');
    } catch (err) {
      setError(err.message);
      setStatus({ status: 'error', connected: false });
    } finally {
      setBusy(false);
    }
  };

  const handleDisconnect = async () => {
    setBusy(true);
    setError(null);
    try {
      const st = await opcuaDisconnect();
      setStatus(st);
    } catch (err) {
      setError(err.message);
    } finally {
      setBusy(false);
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
          className="w-full bg-dark-800 border border-dark-600 rounded-md px-2 py-1.5 text-sm text-white font-mono"
        />
      )}
    </label>
  );

  return (
    <div className="border border-blue-500/30 bg-blue-500/5 rounded-lg p-4">
      <div className="flex items-center gap-2 mb-3">
        <span className="text-xl">🔌</span>
        <h3 className="font-medium text-white">Connexion OPC-UA</h3>
        <span className="ml-auto"><StatusBadge status={status.status} /></span>
      </div>

      {field('Endpoint (URL)', 'endpoint')}
      <div className="grid grid-cols-2 gap-2">
        {field('Mode de sécurité', 'security_mode', 'text', SECURITY_MODES)}
        {field('Politique de sécurité', 'security_policy', 'text', SECURITY_POLICIES)}
      </div>
      <div className="grid grid-cols-2 gap-2">
        {field('Utilisateur (optionnel)', 'username')}
        {field('Mot de passe (optionnel)', 'password', 'password')}
      </div>
      {field('Timeout de session (s)', 'session_timeout', 'number')}

      {form.security_mode !== 'None' && (
        <p className="text-[11px] text-yellow-400/80 mb-2">
          ⚠️ Les modes Sign / SignAndEncrypt nécessitent un certificat client (non encore configuré). Utilisez « None » pour l'instant.
        </p>
      )}

      {error && (
        <div className="bg-red-500/15 border border-red-500/40 rounded-md p-2 text-red-400 text-xs mb-3">
          ❌ {error}
        </div>
      )}

      <div className="flex gap-2">
        <button
          onClick={handleConnect}
          disabled={busy || !form.endpoint.trim()}
          className="flex-1 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white text-sm px-3 py-1.5 rounded-md transition"
        >
          {busy ? '…' : status.connected ? 'Reconnecter' : 'Connecter'}
        </button>
        {status.connected && (
          <button
            onClick={handleDisconnect}
            disabled={busy}
            className="bg-dark-700 hover:bg-dark-600 text-white text-sm px-3 py-1.5 rounded-md transition"
          >
            Déconnecter
          </button>
        )}
      </div>
    </div>
  );
}
