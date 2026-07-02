import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import OpcuaConnectionPanel from '../components/OpcuaConnectionPanel';
import OpcuaTagSelector from '../components/OpcuaTagSelector';
import { useStudioStore } from '../store/studioStore';

// Dynamic OPC-UA flow: connect to a user-specified server, discover its tags,
// choose per-tag data routing (Raw / ISA-95 / Both), then jump to Compose with
// the opcua_read connector applied to the trigger.
export default function OpcuaConnectPage() {
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
          <Link to="/connect" className="text-dark-400 hover:text-white text-sm">← Connecteurs</Link>
        </div>
        <h2 className="text-xl font-semibold text-white mb-1">🔌 Connexion OPC-UA dynamique</h2>
        <p className="text-dark-400 text-sm mb-6">
          Configurez la connexion, découvrez les tags, puis choisissez comment chaque
          donnée circule. Les données <span className="text-blue-300">ISA-95</span> sont
          utilisables dans toutes les fonctions ; les données <span className="text-dark-300">brutes</span> sont
          stockées uniquement.
        </p>

        <div className="space-y-5">
          <OpcuaConnectionPanel onConnected={() => setConnected(true)} />

          {connected ? (
            <OpcuaTagSelector onApplied={handleApplied} />
          ) : (
            <div className="border border-dark-700 bg-dark-900 rounded-lg p-4 text-center text-dark-500 text-sm">
              Connectez-vous à un serveur OPC-UA pour découvrir ses tags.
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
