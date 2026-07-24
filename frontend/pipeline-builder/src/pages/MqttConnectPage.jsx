import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { fetchTopics } from '../api/client';

// MQTT configuration page (2026-07-21) — mirrors the OPC-UA and SQL connector
// pages: real broker status + live topics, not just a shortcut into Compose.
// The mqtt_subscribe trigger's actual topic/QoS binding still happens on the
// trigger node inside Compose (NodeConfigPanel already has a connector picker
// there) — this page is where you check the broker and pick which topic to
// use before you go configure that trigger, not a forced jump past it.
const DEFAULT_TOPIC = 'mindset/events/status-change';
const DEFAULT_QOS = 1;

export default function MqttConnectPage() {
  const { t } = useTranslation();
  const [status, setStatus] = useState(null); // { broker, broker_connected, topics, total }
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [topic, setTopic] = useState(DEFAULT_TOPIC);
  const [qos, setQos] = useState(DEFAULT_QOS);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      setLoading(true);
      try {
        const data = await fetchTopics();
        if (!cancelled) {
          setStatus(data);
          setError(null);
        }
      } catch (err) {
        if (!cancelled) setError(err.message);
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <div className="h-full overflow-y-auto p-6">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center gap-2 mb-1">
          <Link to="/connectors" className="text-dark-400 hover:text-white text-sm">← {t('nav.connectors')}</Link>
        </div>
        <h2 className="text-xl font-semibold text-white mb-1">📶 {t('mqtt.title')}</h2>
        <p className="text-dark-400 text-sm mb-6">
          {t('mqtt.subtitlePre')} <code className="text-blue-300">config/agent.yaml</code>
          {t('mqtt.subtitleMid')} <code className="text-blue-300">mqtt_subscribe</code>
          {t('mqtt.subtitlePost')}
        </p>

        {error && (
          <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-3 text-red-400 text-sm mb-4">
            ❌ {error}
          </div>
        )}

        {/* Broker status */}
        <h3 className="text-sm font-semibold text-dark-300 uppercase tracking-wider mb-2">Broker</h3>
        <div className="border border-dark-700 bg-dark-900 rounded-lg p-4 mb-8 flex items-center gap-3">
          <span className="text-xl">📡</span>
          <div className="min-w-0 flex-1">
            <div className="font-medium text-white font-mono text-sm">
              {loading ? t('common.loading') : status?.broker || t('mqtt.unknownBroker')}
            </div>
          </div>
          <span
            className={`text-xs px-2 py-0.5 rounded-full border font-mono ${
              status?.broker_connected
                ? 'bg-green-500/20 text-green-400 border-green-500/40'
                : 'bg-dark-700 text-dark-300 border-dark-600'
            }`}
          >
            {status?.broker_connected ? t('mqtt.connected') : t('mqtt.disconnected')}
          </span>
        </div>

        {/* Live topics — the MQTT equivalent of "discover" */}
        <h3 className="text-sm font-semibold text-dark-300 uppercase tracking-wider mb-2">
          {t('mqtt.activeTopics', { count: status?.total ?? 0 })}
        </h3>
        <div className="space-y-2 mb-8">
          {!loading && (status?.topics || []).length === 0 && (
            <p className="text-dark-500 text-sm">
              {t('mqtt.noMessages')}
            </p>
          )}
          {(status?.topics || []).map((t) => (
            <button
              key={t.topic}
              onClick={() => setTopic(t.topic)}
              className="w-full text-left border border-dark-700 bg-dark-900 hover:border-blue-500/50 rounded-lg p-3 flex items-center gap-3 transition"
            >
              <span className="font-mono text-sm text-white flex-1 truncate">{t.topic}</span>
              <span className="text-xs text-dark-500 uppercase">{t.category}</span>
              <span className="text-xs text-dark-400 font-mono">{t.rate_per_sec?.toFixed(1) ?? '0.0'} msg/s</span>
            </button>
          ))}
        </div>

        {/* Trigger config preview — applied manually inside Compose, not auto-jumped to */}
        <h3 className="text-sm font-semibold text-dark-300 uppercase tracking-wider mb-2">
          {t('mqtt.triggerConfig')} <code className="text-blue-300 normal-case">mqtt_subscribe</code>
        </h3>
        <div className="border border-blue-500/30 bg-blue-500/5 rounded-lg p-4">
          <label className="block mb-3">
            <span className="text-xs text-dark-400 block mb-1">Topic</span>
            <input
              type="text"
              value={topic}
              onChange={(e) => setTopic(e.target.value)}
              className="w-full bg-dark-800 border border-dark-600 rounded-md px-2 py-1.5 text-sm text-white font-mono"
            />
          </label>
          <label className="block mb-1">
            <span className="text-xs text-dark-400 block mb-1">QoS</span>
            <select
              value={qos}
              onChange={(e) => setQos(Number(e.target.value))}
              className="w-full bg-dark-800 border border-dark-600 rounded-md px-2 py-1.5 text-sm text-white"
            >
              <option value={0}>{t('mqtt.qos0')}</option>
              <option value={1}>{t('mqtt.qos1')}</option>
              <option value={2}>{t('mqtt.qos2')}</option>
            </select>
          </label>
          <p className="text-xs text-dark-500 mt-3">
            {t('mqtt.triggerHintPre')} <Link to="/compose" className="text-blue-300 hover:underline">Compose</Link>
            {t('mqtt.triggerHintMid')} <code className="text-blue-300">mqtt_subscribe</code>
            {t('mqtt.triggerHintPost')}
          </p>
        </div>
      </div>
    </div>
  );
}
