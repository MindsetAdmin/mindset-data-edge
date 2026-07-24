import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import OpcuaConnectionPanel from '../components/OpcuaConnectionPanel';
import OpcuaTagSelector from '../components/OpcuaTagSelector';
import { useStudioStore } from '../store/studioStore';

// Dynamic OPC-UA flow: connect to a user-specified server, discover its tags,
// choose per-tag data routing (Raw / ISA-95 / Both), then jump to Compose with
// the opcua_read connector applied to the trigger.
export default function OpcuaConnectPage() {
  const { t } = useTranslation();
  const [connected, setConnected] = useState(false);
  const setOpcuaSelections = useStudioStore((s) => s.setOpcuaSelections);
  const navigate = useNavigate();

  const handleApplied = (selections) => {
    setOpcuaSelections(selections);
    navigate('/compose');
  };

  return (
    <div className="h-full overflow-y-auto p-6">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center gap-2 mb-1">
          <Link to="/connectors" className="text-dark-400 hover:text-white text-sm">← {t('nav.connectors')}</Link>
        </div>
        <h2 className="text-xl font-semibold text-white mb-1">🔌 {t('opcuaConnect.title')}</h2>
        <p className="text-dark-400 text-sm mb-6">
          {t('opcuaConnect.subtitlePre')} <span className="text-blue-300">ISA-95</span>
          {t('opcuaConnect.subtitleMid')} <span className="text-dark-300">{t('opcuaConnect.raw')}</span>
          {t('opcuaConnect.subtitlePost')}
        </p>

        <div className="space-y-5">
          <OpcuaConnectionPanel onConnected={() => setConnected(true)} />

          {connected ? (
            <OpcuaTagSelector onApplied={handleApplied} />
          ) : (
            <div className="border border-dark-700 bg-dark-900 rounded-lg p-4 text-center text-dark-500 text-sm">
              {t('opcuaConnect.connectHint')}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
