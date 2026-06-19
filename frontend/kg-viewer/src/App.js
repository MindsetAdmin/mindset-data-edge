// src/App.js
import cytoscape from 'cytoscape'

export default class App {
    constructor() {
        this.cy = null
        this.currentGraph = null
        this.selectedNode = null
    }
    
    async init() {
        this.setupEventListeners()
        await this.loadGraph()
    }
    
    setupEventListeners() {
        // Contrôles
        document.getElementById('btn-fit').addEventListener('click', () => this.fit())
        document.getElementById('btn-zoom-in').addEventListener('click', () => this.zoomIn())
        document.getElementById('btn-zoom-out').addEventListener('click', () => this.zoomOut())
        document.getElementById('btn-reset').addEventListener('click', () => this.reset())
        document.getElementById('btn-refresh').addEventListener('click', () => this.refresh())
        document.getElementById('close-details').addEventListener('click', () => this.closeDetails())
    }
    
    async loadGraph() {
        this.showLoading()
        
        try {
            // Charger le graphe depuis l'API
            const response = await fetch('/api/kg/technical')
            const graph = await response.json()
            this.currentGraph = graph
            
            this.updateStats(graph)
            this.renderGraph(graph)
        } catch (error) {
            console.error('Failed to load graph:', error)
            // Fallback: utiliser des données mock
            this.loadMockData()
        } finally {
            this.hideLoading()
        }
    }
    
    loadMockData() {
        // Données mock pour développement
        const mockGraph = {
            nodes: [
                { id: "conn_opcua", type: "connection", name: "OPC-UA Connection", properties: { endpoint: "opc.tcp://localhost:4840" } },
                { id: "topic_raw", type: "topic", name: "mindset/raw/#", properties: { qos: 1 } },
                { id: "pipeline_opcua_to_uns", type: "pipeline", name: "OPC-UA to UNS", properties: { enabled: true } },
                { id: "function_uns_mapper", type: "function", name: "uns_mapper", properties: { type: "transform" } },
                { id: "topic_site", type: "topic", name: "mindset/site/#", properties: { qos: 1 } }
            ],
            edges: [
                { id: "e1", from: "conn_opcua", to: "topic_raw", type: "publishes_to", weight: 1 },
                { id: "e2", from: "topic_raw", to: "pipeline_opcua_to_uns", type: "triggers", weight: 1 },
                { id: "e3", from: "pipeline_opcua_to_uns", to: "function_uns_mapper", type: "depends_on", weight: 1 },
                { id: "e4", from: "pipeline_opcua_to_uns", to: "topic_site", type: "produces", weight: 1 }
            ]
        }
        this.currentGraph = mockGraph
        this.updateStats(mockGraph)
        this.renderGraph(mockGraph)
    }
    
    renderGraph(graph) {
        // Convertir les données pour Cytoscape
        const elements = []
        
        // Ajouter les nœuds
        for (const node of graph.nodes) {
            elements.push({
                data: {
                    id: node.id,
                    label: node.name,
                    type: node.type,
                    properties: JSON.stringify(node.properties, null, 2)
                },
                classes: node.type
            })
        }
        
        // Ajouter les relations
        for (const edge of graph.edges) {
            elements.push({
                data: {
                    id: edge.id,
                    source: edge.from,
                    target: edge.to,
                    label: edge.type
                },
                classes: edge.type
            })
        }
        
        // Initialiser Cytoscape
        const container = document.getElementById('cy')
        container.innerHTML = ''
        
        this.cy = cytoscape({
            container: container,
            elements: elements,
            style: this.getStyles(),
            layout: {
                name: 'fcose',
                quality: 'proof',
                nodeRepulsion: 4500,
                idealEdgeLength: 100,
                edgeElasticity: 0.45,
                nestingFactor: 0.1,
                gravity: 0.25,
                numIter: 2500,
                tile: true,
                animate: true,
                animationDuration: 500
            },
            wheelSensitivity: 0.5,
            minZoom: 0.1,
            maxZoom: 3
        })
        
        // Événements
        this.cy.on('tap', 'node', (evt) => {
            const node = evt.target
            this.showNodeDetails(node.data())
        })
        
        this.cy.on('tap', (evt) => {
            if (evt.target === this.cy) {
                this.closeDetails()
            }
        })
        
        // Ajuster la vue
        setTimeout(() => this.fit(), 100)
    }
    
    getStyles() {
        return [
            // Styles des nœuds
            {
                selector: 'node',
                style: {
                    'label': 'data(label)',
                    'text-valign': 'center',
                    'text-halign': 'center',
                    'font-size': '11px',
                    'font-weight': 'bold',
                    'color': '#fff',
                    'width': '60px',
                    'height': '60px',
                    'border-width': '2px',
                    'border-color': '#fff'
                }
            },
            {
                selector: 'node[type="connection"]',
                style: {
                    'background-color': '#3b82f6',
                    'shape': 'rectangle'
                }
            },
            {
                selector: 'node[type="topic"]',
                style: {
                    'background-color': '#10b981',
                    'shape': 'ellipse'
                }
            },
            {
                selector: 'node[type="function"]',
                style: {
                    'background-color': '#f59e0b',
                    'shape': 'round-rectangle'
                }
            },
            {
                selector: 'node[type="pipeline"]',
                style: {
                    'background-color': '#ef4444',
                    'shape': 'diamond'
                }
            },
            {
                selector: 'node[type="dashboard"]',
                style: {
                    'background-color': '#8b5cf6',
                    'shape': 'hexagon'
                }
            },
            // Styles des relations
            {
                selector: 'edge',
                style: {
                    'width': 2,
                    'line-color': '#6b7280',
                    'target-arrow-color': '#6b7280',
                    'target-arrow-shape': 'triangle',
                    'curve-style': 'bezier',
                    'label': 'data(label)',
                    'font-size': '9px',
                    'color': '#8892b0',
                    'text-rotation': 'autorotate'
                }
            },
            {
                selector: 'edge[label="publishes_to"]',
                style: {
                    'line-color': '#3b82f6',
                    'target-arrow-color': '#3b82f6'
                }
            },
            {
                selector: 'edge[label="subscribes_to"]',
                style: {
                    'line-color': '#10b981',
                    'target-arrow-color': '#10b981'
                }
            },
            {
                selector: 'edge[label="depends_on"]',
                style: {
                    'line-color': '#f59e0b',
                    'target-arrow-color': '#f59e0b'
                }
            },
            {
                selector: 'edge[label="produces"]',
                style: {
                    'line-color': '#ef4444',
                    'target-arrow-color': '#ef4444'
                }
            },
            {
                selector: 'edge[label="triggers"]',
                style: {
                    'line-color': '#8b5cf6',
                    'target-arrow-color': '#8b5cf6'
                }
            }
        ]
    }
    
    showNodeDetails(nodeData) {
        const panel = document.getElementById('details-panel')
        const content = document.getElementById('details-content')
        
        let propertiesHtml = ''
        if (nodeData.properties && nodeData.properties !== 'undefined') {
            try {
                const props = JSON.parse(nodeData.properties)
                for (const [key, value] of Object.entries(props)) {
                    propertiesHtml += `
                        <div class="detail-row">
                            <div class="detail-label">${key}</div>
                            <div class="detail-value">${value}</div>
                        </div>
                    `
                }
            } catch (e) {
                propertiesHtml = '<div class="detail-row">Aucune propriété</div>'
            }
        }
        
        content.innerHTML = `
            <div class="detail-row">
                <div class="detail-label">ID</div>
                <div class="detail-value">${nodeData.id}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Type</div>
                <div class="detail-value">${nodeData.type}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Nom</div>
                <div class="detail-value">${nodeData.label}</div>
            </div>
            ${propertiesHtml}
        `
        
        panel.classList.add('open')
        this.selectedNode = nodeData
    }
    
    closeDetails() {
        const panel = document.getElementById('details-panel')
        panel.classList.remove('open')
        this.selectedNode = null
    }
    
    updateStats(graph) {
        document.getElementById('node-count').innerText = graph.nodes.length
        document.getElementById('edge-count').innerText = graph.edges.length
    }
    
    fit() {
        if (this.cy) {
            this.cy.fit()
            this.cy.zoom(0.8)
        }
    }
    
    zoomIn() {
        if (this.cy) {
            this.cy.zoom(this.cy.zoom() * 1.2)
        }
    }
    
    zoomOut() {
        if (this.cy) {
            this.cy.zoom(this.cy.zoom() * 0.8)
        }
    }
    
    reset() {
        if (this.cy) {
            this.cy.reset()
            this.cy.fit()
        }
    }
    
    async refresh() {
        await this.loadGraph()
    }
    
    showLoading() {
        const container = document.getElementById('cy')
        container.innerHTML = `
            <div class="loading">
                <div class="spinner"></div>
                <p>Chargement du Knowledge Graph...</p>
            </div>
        `
    }
    
    hideLoading() {
        // Le contenu sera remplacé par Cytoscape
    }
}