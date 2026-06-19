export default class NodeDetails {
    constructor(containerId) {
        this.containerId = containerId
        this.panel = null
    }
    
    show(nodeData) {
        let panel = document.getElementById(this.containerId)
        if (!panel) return
        
        const propertiesHtml = this.renderProperties(nodeData.properties)
        
        panel.innerHTML = `
            <div class="details-panel open" id="details-inner">
                <div class="details-header">
                    <h3>📋 Détails du nœud</h3>
                    <button class="close-btn" id="details-close">×</button>
                </div>
                <div class="details-content">
                    <div class="detail-row">
                        <div class="detail-label">ID</div>
                        <div class="detail-value">${this.escapeHtml(nodeData.id)}</div>
                    </div>
                    <div class="detail-row">
                        <div class="detail-label">Type</div>
                        <div class="detail-value">${this.getTypeLabel(nodeData.type)}</div>
                    </div>
                    <div class="detail-row">
                        <div class="detail-label">Nom</div>
                        <div class="detail-value">${this.escapeHtml(nodeData.label)}</div>
                    </div>
                    ${propertiesHtml}
                </div>
            </div>
        `
        
        // Attach close event
        const closeBtn = document.getElementById('details-close')
        if (closeBtn) {
            closeBtn.addEventListener('click', () => this.hide())
        }
    }
    
    renderProperties(properties) {
        if (!properties || Object.keys(properties).length === 0) {
            return '<div class="detail-row">Aucune propriété</div>'
        }
        
        let html = ''
        for (const [key, value] of Object.entries(properties)) {
            const displayValue = typeof value === 'object' 
                ? JSON.stringify(value) 
                : value
            html += `
                <div class="detail-row">
                    <div class="detail-label">${this.escapeHtml(key)}</div>
                    <div class="detail-value">${this.escapeHtml(String(displayValue))}</div>
                </div>
            `
        }
        return html
    }
    
    getTypeLabel(type) {
        const labels = {
            'connection': '🔌 Connection',
            'topic': '📡 Topic',
            'function': '⚙️ Function',
            'pipeline': '🔧 Pipeline',
            'dashboard': '📊 Dashboard'
        }
        return labels[type] || type
    }
    
    hide() {
        const panel = document.getElementById(this.containerId)
        if (panel) {
            panel.innerHTML = ''
        }
    }
    
    escapeHtml(str) {
        if (!str) return ''
        return str
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#39;')
    }
}