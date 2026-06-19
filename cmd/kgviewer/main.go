// cmd/kgviewer/main.go
// Standalone viewer for the DOMAIN knowledge graph stored in data/mindset.db.
// Unlike cmd/agent, this binary needs NO OPC-UA and NO MQTT — it just opens the
// SQLite DB and serves the graph (machines, micro-stops, causes, costs) in a browser.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/MindsetAdmin/mindset-data-edge/internal/kg"
)

// viewerHTML is the Cytoscape page, adapted to the domain GraphJSON shape
// (label / from_id / to_id / relation) served by /api/kg/domain.
const viewerHTML = `<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MindSet Data - Domain Knowledge Graph</title>
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
        .legend-color.equipment { background: #3b82f6; }
        .legend-color.event { background: #f59e0b; }
        .legend-color.cause { background: #8b5cf6; }
        .legend-color.cost { background: #10b981; }

        .footer { margin-top: auto; text-align: center; font-size: 0.65rem; color: #64748b; padding-top: 16px; }

        .main-content { flex: 1; position: relative; display: flex; }
        .graph-container { flex: 1; position: relative; background: #0f172a; }
        #cy { width: 100%; height: 100%; background: #0f172a; }
        .empty { position: absolute; inset: 0; display: none; align-items: center; justify-content: center; flex-direction: column; gap: 8px; color: #64748b; text-align: center; padding: 24px; }
        .empty.show { display: flex; }

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
            <div class="logo"><h1>🧠 MindSet Data</h1><p>Domain Knowledge Graph</p></div>
            <div class="stats"><div class="stat-item"><span class="stat-label">📦 Nœuds</span><span class="stat-value" id="node-count">-</span></div><div class="stat-item"><span class="stat-label">🔗 Relations</span><span class="stat-value" id="edge-count">-</span></div></div>
            <div class="controls"><button id="btn-fit">🔍 Fit</button><button id="btn-zoom-in">➕</button><button id="btn-zoom-out">➖</button><button id="btn-reset">🔄 Reset</button><button id="btn-refresh">⬇️ Refresh</button></div>
            <div class="legend"><h4>Légende</h4><div class="legend-item"><span class="legend-color equipment"></span><span>Equipment</span></div><div class="legend-item"><span class="legend-color event"></span><span>Event (micro-stop)</span></div><div class="legend-item"><span class="legend-color cause"></span><span>Cause</span></div><div class="legend-item"><span class="legend-color cost"></span><span>Cost</span></div></div>
            <div class="footer"><p>Powered by Cytoscape.js</p><p>data/mindset.db · read-only viewer</p></div>
        </aside>
        <main class="main-content">
            <div class="graph-container"><div id="cy"></div><div class="empty" id="empty"><div style="font-size:2rem">📭</div><div>Le knowledge graph est vide.</div><div style="font-size:0.7rem">Lance l'agent pour générer des micro-stops, puis Refresh.</div></div></div>
            <div id="details-panel" class="details-panel"><div class="details-header"><h3>📋 Détails</h3><button class="close-btn" id="close-details">×</button></div><div class="details-content" id="details-content"><p>Sélectionnez un nœud</p></div></div>
        </main>
    </div>
    <script>
        let cy = null, graphData = null;

        function updateStats() {
            if(graphData){
                document.getElementById('node-count').innerText = (graphData.nodes || []).length;
                document.getElementById('edge-count').innerText = (graphData.edges || []).length;
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
            content.innerHTML = '<div class="detail-row"><div class="detail-label">ID</div><div class="detail-value">'+nodeData.id+'</div></div><div class="detail-row"><div class="detail-label">Type</div><div class="detail-value">'+nodeData.type+'</div></div><div class="detail-row"><div class="detail-label">Label</div><div class="detail-value">'+nodeData.label+'</div></div>'+propsHtml;
            document.getElementById('details-panel').classList.add('open');
        }

        function hideDetails(){
            document.getElementById('details-panel').classList.remove('open');
        }

        function renderGraph(graph){
            graphData = graph;
            updateStats();

            const nodes = graph.nodes || [];
            const edges = graph.edges || [];
            document.getElementById('empty').classList.toggle('show', nodes.length === 0);

            const elements = [];
            for(const node of nodes){
                // Domain shape: id, type, label, properties
                let label = node.label || node.id;
                if(label.length > 22) label = label.substring(0, 20) + '..';
                elements.push({ data: { id: node.id, label: label, type: node.type, properties: node.properties }, classes: node.type });
            }
            for(const edge of edges){
                // Domain shape: id, from_id, to_id, relation
                elements.push({ data: { id: edge.id, source: edge.from_id, target: edge.to_id, label: edge.relation }, classes: edge.relation });
            }

            const styles = [
                { selector: 'node', style: {
                    'label': 'data(label)',
                    'text-valign': 'top',
                    'text-halign': 'center',
                    'text-margin-y': '8px',
                    'font-size': '8px',
                    'font-weight': '600',
                    'color': '#f1f5f9',
                    'text-wrap': 'wrap',
                    'text-max-width': '70px',
                    'width': '40px',
                    'height': '40px',
                    'border-width': '2px',
                    'border-color': '#0f172a',
                    'background-opacity': 0.9
                } },
                // Domain node types (match labels written by internal/kg/graph.go)
                { selector: 'node[type="Equipment"]', style: { 'background-color': '#3b82f6', 'shape': 'round-rectangle' } },
                { selector: 'node[type="Event"]',     style: { 'background-color': '#f59e0b', 'shape': 'ellipse' } },
                { selector: 'node[type="Cause"]',     style: { 'background-color': '#8b5cf6', 'shape': 'diamond' } },
                { selector: 'node[type="Cost"]',      style: { 'background-color': '#10b981', 'shape': 'round-rectangle' } },
                { selector: 'edge', style: {
                    'width': 1.4,
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
                { selector: 'edge[label="occurred_at"]', style: { 'line-color': '#fbbf24', 'target-arrow-color': '#fbbf24' } },
                { selector: 'edge[label="caused_by"]',   style: { 'line-color': '#a78bfa', 'target-arrow-color': '#a78bfa' } },
                { selector: 'edge[label="costs"]',       style: { 'line-color': '#34d399', 'target-arrow-color': '#34d399' } }
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

            setTimeout(() => { if(cy && nodes.length){ cy.fit(); cy.zoom(0.85); } }, 500);
        }

        async function loadGraph(){
            try{
                const resp = await fetch('/api/kg/domain');
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

func main() {
	dbPath := flag.String("db", "./data/mindset.db", "path to the mindset SQLite database")
	addr := flag.String("addr", ":8090", "HTTP listen address")
	flag.Parse()

	kgInstance, err := kg.NewKnowledgeGraph(*dbPath)
	if err != nil {
		log.Fatalf("[KGVIEWER] Failed to open KG at %s: %v", *dbPath, err)
	}
	defer kgInstance.Close()

	mux := http.NewServeMux()

	// Page d'accueil — viewer du graphe domaine
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(viewerHTML))
	})

	// API: graphe domaine (nœuds + relations) depuis SQLite
	mux.HandleFunc("/api/kg/domain", func(w http.ResponseWriter, r *http.Request) {
		graph, err := kgInstance.GetFullGraph()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		json.NewEncoder(w).Encode(graph)
	})

	// API: santé
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Printf("[KGVIEWER] Reading DB: %s", *dbPath)
	log.Printf("[KGVIEWER] Domain KG viewer at http://localhost%s", *addr)
	log.Printf("[KGVIEWER] API at http://localhost%s/api/kg/domain", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
