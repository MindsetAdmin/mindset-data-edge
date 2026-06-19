// cmd/agent/server.go - version avec cache
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MindsetAdmin/mindset-data-edge/internal/functions"
	"github.com/MindsetAdmin/mindset-data-edge/internal/kg"
	"github.com/MindsetAdmin/mindset-data-edge/internal/pipeline"
)

// kgDashboardHTML est le template HTML pour le KG Viewer (version finale)
const kgDashboardHTML = `<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MindSet Data - Knowledge Graph Explorer</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #0f172a; color: #e2e8f0; }
        .app-container { display: flex; height: 100vh; width: 100vw; overflow: hidden; }
        
        .sidebar { width: 260px; background: #1e293b; border-right: 1px solid #334155; display: flex; flex-direction: column; padding: 20px; gap: 20px; overflow-y: auto; }
        .logo h1 { font-size: 1.3rem; font-weight: bold; background: linear-gradient(135deg, #38bdf8 0%, #818cf8 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; margin-bottom: 4px; }
        .logo p { font-size: 0.7rem; color: #94a3b8; }
        .stats { background: #0f172a; border-radius: 10px; padding: 12px; display: flex; justify-content: space-between; border: 1px solid #334155; }
        .stat-item { text-align: center; }
        .stat-label { font-size: 0.65rem; color: #94a3b8; display: block; }
        .stat-value { font-size: 1.5rem; font-weight: bold; color: #38bdf8; }
        .controls { display: flex; flex-wrap: wrap; gap: 6px; }
        .controls button { background: #334155; border: none; color: #e2e8f0; padding: 6px 10px; border-radius: 6px; cursor: pointer; transition: all 0.2s; font-size: 0.8rem; flex: 1; }
        .controls button:hover { background: #475569; }
        
        .legend h4 { font-size: 0.8rem; margin-bottom: 10px; color: #94a3b8; }
        .legend-item { display: flex; align-items: center; gap: 8px; margin-bottom: 6px; font-size: 0.75rem; }
        .legend-color { width: 12px; height: 12px; border-radius: 3px; border: 1px solid #1e293b; }
        .legend-color.connection { background: #3b82f6; }
        .legend-color.topic { background: #10b981; }
        .legend-color.function { background: #f59e0b; }
        .legend-color.pipeline { background: #ef4444; }
        .legend-color.dashboard { background: #8b5cf6; }
        
        .footer { margin-top: auto; text-align: center; font-size: 0.65rem; color: #64748b; padding-top: 16px; }
        
        .main-content { flex: 1; position: relative; display: flex; }
        .graph-container { flex: 1; position: relative; background: #0f172a; }
        #cy { width: 100%; height: 100%; background: #0f172a; }
        
        .details-panel { position: absolute; right: 0; top: 0; width: 300px; height: 100%; background: #1e293b; border-left: 1px solid #334155; transform: translateX(100%); transition: transform 0.3s ease; display: flex; flex-direction: column; z-index: 10; }
        .details-panel.open { transform: translateX(0); }
        .details-header { display: flex; justify-content: space-between; align-items: center; padding: 14px; border-bottom: 1px solid #334155; }
        .details-header h3 { font-size: 0.9rem; margin: 0; color: #38bdf8; }
        .close-btn { background: none; border: none; color: #94a3b8; font-size: 1.3rem; cursor: pointer; }
        .close-btn:hover { color: #e2e8f0; }
        .details-content { flex: 1; padding: 14px; overflow-y: auto; font-size: 0.8rem; }
        .detail-row { margin-bottom: 10px; padding-bottom: 6px; border-bottom: 1px solid #334155; }
        .detail-label { color: #94a3b8; font-size: 0.65rem; text-transform: uppercase; margin-bottom: 3px; }
        .detail-value { font-family: monospace; word-break: break-all; color: #e2e8f0; font-size: 0.75rem; }
        
        @media (max-width: 768px) { .sidebar { width: 220px; } .details-panel { width: 260px; } }
    </style>
    <script src="https://unpkg.com/cytoscape@3.28.1/dist/cytoscape.min.js"></script>
</head>
<body>
    <div class="app-container">
        <aside class="sidebar">
            <div class="logo"><h1>📊 MindSet Data</h1><p>Knowledge Graph Explorer</p></div>
            <div class="stats"><div class="stat-item"><span class="stat-label">📦 Nœuds</span><span class="stat-value" id="node-count">-</span></div><div class="stat-item"><span class="stat-label">🔗 Relations</span><span class="stat-value" id="edge-count">-</span></div></div>
            <div class="controls"><button id="btn-fit">🔍 Fit</button><button id="btn-zoom-in">➕</button><button id="btn-zoom-out">➖</button><button id="btn-reset">🔄 Reset</button><button id="btn-refresh">⬇️ Refresh</button></div>
            <div class="legend"><h4>Légende</h4><div class="legend-item"><span class="legend-color connection"></span><span>Connection</span></div><div class="legend-item"><span class="legend-color topic"></span><span>Topic MQTT</span></div><div class="legend-item"><span class="legend-color function"></span><span>Function</span></div><div class="legend-item"><span class="legend-color pipeline"></span><span>Pipeline</span></div><div class="legend-item"><span class="legend-color dashboard"></span><span>Dashboard</span></div></div>
            <div class="footer"><p>Powered by Cytoscape.js</p><p>MindSet Data - v1.0</p></div>
        </aside>
        <main class="main-content">
            <div class="graph-container"><div id="cy"></div></div>
            <div id="details-panel" class="details-panel"><div class="details-header"><h3>📋 Détails</h3><button class="close-btn" id="close-details">×</button></div><div class="details-content" id="details-content"><p>Sélectionnez un nœud</p></div></div>
        </main>
    </div>
    <script>
        let cy = null, graphData = null;
        
        function updateStats() { 
            if(graphData){ 
                document.getElementById('node-count').innerText = graphData.nodes.length; 
                document.getElementById('edge-count').innerText = graphData.edges.length; 
            } 
        }
        
        function showDetails(nodeData) { 
            const content = document.getElementById('details-content'); 
            let propsHtml = ''; 
            if(nodeData.properties){ 
                for(const [k,v] of Object.entries(nodeData.properties)){ 
                    let displayValue = v;
                    if(typeof v === 'object') displayValue = JSON.stringify(v);
                    propsHtml += '<div class="detail-row"><div class="detail-label">'+k+'</div><div class="detail-value">'+displayValue+'</div></div>'; 
                } 
            } 
            content.innerHTML = '<div class="detail-row"><div class="detail-label">ID</div><div class="detail-value">'+nodeData.id+'</div></div><div class="detail-row"><div class="detail-label">Type</div><div class="detail-value">'+nodeData.type+'</div></div><div class="detail-row"><div class="detail-label">Nom</div><div class="detail-value">'+nodeData.label+'</div></div>'+propsHtml; 
            document.getElementById('details-panel').classList.add('open'); 
        }
        
        function hideDetails(){ 
            document.getElementById('details-panel').classList.remove('open'); 
        }
        
        function renderGraph(graph){ 
            graphData = graph; 
            updateStats(); 
            
            const elements = []; 
            for(const node of graph.nodes){ 
                // Nettoyer le label
                let label = node.name;
                if(label.length > 20) label = label.substring(0, 18) + '..';
                elements.push({ data: { id: node.id, label: label, type: node.type, properties: node.properties }, classes: node.type }); 
            } 
            for(const edge of graph.edges){ 
                elements.push({ data: { id: edge.id, source: edge.from, target: edge.to, label: edge.type }, classes: edge.type }); 
            }
            
            const styles = [ 
                // Nœuds - petits, texte en dessous
                { selector: 'node', style: { 
                    'label': 'data(label)', 
                    'text-valign': 'top',
                    'text-halign': 'center',
                    'text-margin-y': '8px',
                    'font-size': '8px', 
                    'font-weight': '600', 
                    'color': '#f1f5f9', 
                    'text-wrap': 'wrap',
                    'text-max-width': '55px',
                    'width': '36px',
                    'height': '36px',
                    'border-width': '2px',
                    'border-color': '#0f172a',
                    'background-opacity': 0.9
                } },
                // Connection - bleu doux
                { selector: 'node[type="connection"]', style: { 'background-color': '#0ea5e9', 'shape': 'rectangle' } },
                // Topic - vert doux
                { selector: 'node[type="topic"]', style: { 'background-color': '#34d399', 'shape': 'ellipse' } },
                // Function - orange doux
                { selector: 'node[type="function"]', style: { 'background-color': '#fbbf24', 'shape': 'round-rectangle' } },
                // Pipeline - rouge doux
                { selector: 'node[type="pipeline"]', style: { 'background-color': '#f87171', 'shape': 'round-rectangle' } },
                // Dashboard - violet doux
                { selector: 'node[type="dashboard"]', style: { 'background-color': '#a78bfa', 'shape': 'round-rectangle' } },
                // Relations
                { selector: 'edge', style: { 
                    'width': 1.2, 
                    'line-color': '#64748b', 
                    'target-arrow-color': '#64748b', 
                    'target-arrow-shape': 'triangle',
                    'curve-style': 'bezier', 
                    'label': 'data(label)', 
                    'font-size': '6px', 
                    'color': '#94a3b8', 
                    'text-rotation': 'autorotate',
                    'text-background-color': '#1e293b',
                    'text-background-opacity': 0.8,
                    'text-background-padding': '2px',
                    'text-margin-y': '-4px'
                } },
                // Couleurs des relations
                { selector: 'edge[label="publishes_to"]', style: { 'line-color': '#60a5fa', 'target-arrow-color': '#60a5fa' } },
                { selector: 'edge[label="subscribes_to"]', style: { 'line-color': '#34d399', 'target-arrow-color': '#34d399' } },
                { selector: 'edge[label="depends_on"]', style: { 'line-color': '#fbbf24', 'target-arrow-color': '#fbbf24' } },
                { selector: 'edge[label="produces"]', style: { 'line-color': '#f87171', 'target-arrow-color': '#f87171' } },
                { selector: 'edge[label="triggers"]', style: { 'line-color': '#a78bfa', 'target-arrow-color': '#a78bfa' } }
            ];
            
            cy = cytoscape({ 
                container: document.getElementById('cy'), 
                elements: elements, 
                style: styles, 
                layout: { 
                    name: 'breadthfirst',
                    directed: true,
                    padding: 40,
                    spacingFactor: 1.8,
                    animate: true,
                    animationDuration: 600
                }, 
                minZoom: 0.2, 
                maxZoom: 2.5 
            });
            
            cy.on('tap', 'node', (evt) => { showDetails(evt.target.data()); });
            cy.on('tap', (evt) => { if(evt.target === cy) hideDetails(); });
            
            setTimeout(() => { if(cy){ cy.fit(); cy.zoom(0.85); } }, 500);
        }
        
        async function loadGraph(){ 
            try{ 
                const resp = await fetch('/api/kg/graph'); 
                const data = await resp.json(); 
                renderGraph(data); 
            } catch(e){ console.error('Error loading graph:', e); } 
        }
        
        document.getElementById('btn-fit').addEventListener('click', () => cy && cy.fit());
        document.getElementById('btn-zoom-in').addEventListener('click', () => cy && cy.zoom(cy.zoom()*1.2));
        document.getElementById('btn-zoom-out').addEventListener('click', () => cy && cy.zoom(cy.zoom()*0.8));
        document.getElementById('btn-reset').addEventListener('click', () => { if(cy){ cy.reset(); cy.fit(); } });
        document.getElementById('btn-refresh').addEventListener('click', () => loadGraph());
        document.getElementById('close-details').addEventListener('click', () => hideDetails());
        
        loadGraph();
    </script>
</body>
</html>`

// startServer démarre le serveur HTTP unifié avec cache
func startServer(kgInstance *kg.KnowledgeGraph, pipelineEngine *pipeline.Engine, funcRegistry *functions.Registry) {
	mux := http.NewServeMux()

	// Page d'accueil - KG Viewer
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(kgDashboardHTML))
			return
		}
		http.NotFound(w, r)
	})

	// API: retourne le graphe en JSON (avec cache)
	mux.HandleFunc("/api/kg/graph", func(w http.ResponseWriter, r *http.Request) {
		if kgInstance == nil || pipelineEngine == nil {
			http.Error(w, "KG or Pipeline not available", 503)
			return
		}

		// Utiliser le cache avec invalidation automatique
		start := time.Now()
		techGraph, err := kgInstance.GetTechnicalGraphWithCache(pipelineEngine.GetRegistry())
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		elapsed := time.Since(start)
		log.Printf("[API] /api/kg/graph served in %v (%d nodes, %d edges)",
			elapsed, len(techGraph.Nodes), len(techGraph.Edges))

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		json.NewEncoder(w).Encode(techGraph)
	})

	// API: santé
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// API: statistiques
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		stats := map[string]interface{}{"pipelines": 0, "nodes": 0, "edges": 0}
		if pipelineEngine != nil {
			stats["pipelines"] = len(pipelineEngine.ListPipelines())
		}
		if kgInstance != nil && pipelineEngine != nil {
			techGraph, _ := kgInstance.GetTechnicalGraphWithCache(pipelineEngine.GetRegistry())
			if techGraph != nil {
				stats["nodes"] = len(techGraph.Nodes)
				stats["edges"] = len(techGraph.Edges)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	// API: purge du cache (pour forcer un rebuild)
	mux.HandleFunc("/api/kg/purge-cache", func(w http.ResponseWriter, r *http.Request) {
		if kgInstance != nil {
			kgInstance.PurgeCache()
			log.Printf("[API] Cache purged")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "cache purged"})
			return
		}
		http.Error(w, "KG not available", 503)
	})

	// API: liste toutes les fonctions
	mux.HandleFunc("/api/functions", func(w http.ResponseWriter, r *http.Request) {
		if funcRegistry == nil {
			http.Error(w, "Functions registry not available", 503)
			return
		}

		// Vérifier le paramètre type
		typeParam := r.URL.Query().Get("type")
		var fnList []*functions.FunctionInfo

		if typeParam != "" {
			var fnType functions.FunctionType
			switch typeParam {
			case "connector":
				fnType = functions.TypeConnector
			case "transform":
				fnType = functions.TypeTransform
			case "calculate":
				fnType = functions.TypeCalculate
			case "condition":
				fnType = functions.TypeCondition
			case "output":
				fnType = functions.TypeOutput
			default:
				http.Error(w, "Invalid type parameter", 400)
				return
			}
			fnList = funcRegistry.ListFunctionsByType(fnType)
		} else {
			fnList = funcRegistry.ListFunctions()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"functions": fnList,
			"total":     len(fnList),
		})
	})

	// API: liste les connecteurs (alias pour functions?type=connector)
	mux.HandleFunc("/api/connectors", func(w http.ResponseWriter, r *http.Request) {
		if funcRegistry == nil {
			http.Error(w, "Functions registry not available", 503)
			return
		}

		connectors := funcRegistry.ListFunctionsByType(functions.TypeConnector)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"connectors": connectors,
			"total":      len(connectors),
		})
	})

	// Démarrer le serveur
	go func() {
		log.Printf("[HTTP] Server starting on :8080")
		log.Printf("[HTTP] KG Dashboard available at http://localhost:8080")
		log.Printf("[HTTP] API available at http://localhost:8080/api/")
		log.Printf("[HTTP] Cache invalidation: 5 minutes")
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Printf("[HTTP] Server error: %v", err)
		}
	}()
}
