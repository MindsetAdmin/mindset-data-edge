import { useState, useEffect } from 'react';
import { fetchFunctions } from './api/client';
import './App.css';

function App() {
    const [functions, setFunctions] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        loadFunctions();
    }, []);

    const loadFunctions = async () => {
        try {
            setLoading(true);
            const data = await fetchFunctions();
            setFunctions(data.functions || []);
            setError(null);
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-dark-950 text-white">
            <header className="bg-dark-900 border-b border-dark-700 p-4">
                <h1 className="text-2xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
                    🧩 MindSet Data — Pipeline Builder
                </h1>
                <p className="text-sm text-dark-400">Construisez vos pipelines par glisser-déposer</p>
            </header>

            <div className="flex h-[calc(100vh-80px)]">
                {/* Palette */}
                <aside className="w-80 bg-dark-900 border-r border-dark-700 p-4 overflow-y-auto">
                    <h2 className="text-sm font-semibold text-dark-400 uppercase tracking-wider mb-4">
                        📦 Palette de composants
                    </h2>
                    
                    {loading && (
                        <div className="text-center text-dark-400 py-8">
                            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto mb-4"></div>
                            <p>Chargement des fonctions...</p>
                        </div>
                    )}
                    
                    {error && (
                        <div className="bg-red-500/20 border border-red-500/50 rounded-lg p-4 text-red-400 text-sm">
                            ❌ {error}
                        </div>
                    )}
                    
                    {!loading && !error && functions.length === 0 && (
                        <div className="text-center text-dark-400 py-8">
                            <p>Aucune fonction disponible</p>
                        </div>
                    )}
                    
                    {!loading && !error && functions.length > 0 && (
                        <div className="space-y-6">
                            {groupFunctionsByCategory(functions).map(([category, items]) => (
                                <div key={category}>
                                    <h3 className="text-sm font-medium text-dark-300 mb-2">{category}</h3>
                                    <div className="space-y-2">
                                        {items.map((fn) => (
                                            <FunctionCard key={fn.name} function={fn} />
                                        ))}
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </aside>

                {/* Canvas */}
                <main className="flex-1 bg-dark-950 relative">
                    <div className="flex items-center justify-center h-full text-dark-500">
                        <div className="text-center">
                            <p className="text-4xl mb-4">🖱️</p>
                            <p className="text-lg">Glissez une fonction depuis la palette</p>
                            <p className="text-sm text-dark-600">ou sélectionnez un connecteur pour commencer</p>
                        </div>
                    </div>
                </main>
            </div>
        </div>
    );
}

function FunctionCard({ function: fn }) {
    const colorMap = {
        'Connecteurs': 'border-blue-500/30 bg-blue-500/10',
        'Transformations': 'border-orange-500/30 bg-orange-500/10',
        'Calculs': 'border-green-500/30 bg-green-500/10',
        'Conditions': 'border-purple-500/30 bg-purple-500/10',
        'Sorties': 'border-red-500/30 bg-red-500/10',
    };
    
    const iconMap = {
        'Connecteurs': '🔌',
        'Transformations': '⚙️',
        'Calculs': '📊',
        'Conditions': '🚦',
        'Sorties': '📤',
    };

    const category = getCategory(fn.type);
    const colorClass = colorMap[category] || 'border-gray-500/30 bg-gray-500/10';
    const icon = iconMap[category] || '📦';

    return (
        <div className={`border rounded-lg p-3 cursor-grab hover:bg-dark-800 transition ${colorClass}`}>
            <div className="flex items-center gap-2">
                <span className="text-lg">{icon}</span>
                <div>
                    <div className="font-medium text-sm text-white">{fn.name}</div>
                    <div className="text-xs text-dark-400">{fn.description || 'Aucune description'}</div>
                </div>
            </div>
        </div>
    );
}

function groupFunctionsByCategory(functions) {
    const groups = new Map();
    functions.forEach((fn) => {
        const category = getCategory(fn.type);
        if (!groups.has(category)) {
            groups.set(category, []);
        }
        groups.get(category).push(fn);
    });
    return Array.from(groups.entries());
}

function getCategory(type) {
    const map = {
        'connector': 'Connecteurs',
        'transform': 'Transformations',
        'calculate': 'Calculs',
        'condition': 'Conditions',
        'output': 'Sorties',
    };
    return map[type] || 'Autres';
}

export default App;