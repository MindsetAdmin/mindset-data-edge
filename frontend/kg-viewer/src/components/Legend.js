export default class Legend {
    constructor(containerId) {
        this.containerId = containerId
    }
    
    show() {
        const container = document.getElementById(this.containerId)
        if (!container) return
        
        const legendHtml = `
            <div class="logo">
                <h1>📊 MindSet Data</h1>
                <p>Knowledge Graph Explorer</p>
            </div>
            
            <div class="legend">
                <h4>Légende</h4>
                <div class="legend-item">
                    <span class="legend-color connection"></span>
                    <span>🔌 Connection</span>
                </div>
                <div class="legend-item">
                    <span class="legend-color topic"></span>
                    <span>📡 Topic</span>
                </div>
                <div class="legend-item">
                    <span class="legend-color function"></span>
                    <span>⚙️ Function</span>
                </div>
                <div class="legend-item">
                    <span class="legend-color pipeline"></span>
                    <span>🔧 Pipeline</span>
                </div>
                <div class="legend-item">
                    <span class="legend-color dashboard"></span>
                    <span>📊 Dashboard</span>
                </div>
            </div>
            
            <div class="footer">
                <p>Powered by Cytoscape.js</p>
                <p>MindSet Data - v1.0</p>
            </div>
        `
        
        container.innerHTML = legendHtml
    }
}