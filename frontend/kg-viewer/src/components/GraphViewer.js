import cytoscape from 'cytoscape'
import dagre from 'cytoscape-dagre'
import fcose from 'cytoscape-fcose'

// Register extensions
cytoscape.use(dagre)
cytoscape.use(fcose)

export default class GraphViewer {
    constructor(containerId) {
        this.containerId = containerId
        this.cy = null
        this.onNodeSelect = null
        this.onBackgroundTap = null
    }
    
    render(graphData) {
        const container = document.getElementById(this.containerId)
        if (!container) return
        
        container.innerHTML = ''
        
        const elements = this.convertToCytoscapeElements(graphData)
        
        this.cy = cytoscape({
            container: container,
            elements: elements,
            style: this.getStyles(),
            layout: this.getLayout(),
            wheelSensitivity: 0.5,
            minZoom: 0.1,
            maxZoom: 3
        })
        
        this.attachEvents()
        
        setTimeout(() => this.fit(), 100)
    }
    
    convertToCytoscapeElements(graphData) {
        const elements = []
        
        for (const node of graphData.nodes) {
            elements.push({
                data: {
                    id: node.id,
                    label: node.name,
                    type: node.type,
                    properties: node.properties
                },
                classes: node.type
            })
        }
        
        for (const edge of graphData.edges) {
            elements.push({
                data: {
                    id: edge.id,
                    source: edge.from,
                    target: edge.to,
                    label: edge.type,
                    weight: edge.weight
                },
                classes: edge.type
            })
        }
        
        return elements
    }
    
    getLayout() {
        return {
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
        }
    }
    
    getStyles() {
        return [
            // Node base styles
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
            // Node types
            {
                selector: 'node[type="connection"]',
                style: { 'background-color': '#3b82f6', 'shape': 'rectangle' }
            },
            {
                selector: 'node[type="topic"]',
                style: { 'background-color': '#10b981', 'shape': 'ellipse' }
            },
            {
                selector: 'node[type="function"]',
                style: { 'background-color': '#f59e0b', 'shape': 'round-rectangle' }
            },
            {
                selector: 'node[type="pipeline"]',
                style: { 'background-color': '#ef4444', 'shape': 'diamond' }
            },
            {
                selector: 'node[type="dashboard"]',
                style: { 'background-color': '#8b5cf6', 'shape': 'hexagon' }
            },
            // Edge base styles
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
            // Edge types
            {
                selector: 'edge[label="publishes_to"]',
                style: { 'line-color': '#3b82f6', 'target-arrow-color': '#3b82f6' }
            },
            {
                selector: 'edge[label="subscribes_to"]',
                style: { 'line-color': '#10b981', 'target-arrow-color': '#10b981' }
            },
            {
                selector: 'edge[label="depends_on"]',
                style: { 'line-color': '#f59e0b', 'target-arrow-color': '#f59e0b' }
            },
            {
                selector: 'edge[label="produces"]',
                style: { 'line-color': '#ef4444', 'target-arrow-color': '#ef4444' }
            },
            {
                selector: 'edge[label="triggers"]',
                style: { 'line-color': '#8b5cf6', 'target-arrow-color': '#8b5cf6' }
            }
        ]
    }
    
    attachEvents() {
        if (!this.cy) return
        
        this.cy.on('tap', 'node', (evt) => {
            const node = evt.target
            const nodeData = node.data()
            if (this.onNodeSelect) {
                this.onNodeSelect(nodeData)
            }
        })
        
        this.cy.on('tap', (evt) => {
            if (evt.target === this.cy && this.onBackgroundTap) {
                this.onBackgroundTap()
            }
        })
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
    
    showLoading() {
        const container = document.getElementById(this.containerId)
        if (container) {
            container.innerHTML = `
                <div class="loading">
                    <div class="spinner"></div>
                    <p>Chargement du Knowledge Graph...</p>
                </div>
            `
        }
    }
}