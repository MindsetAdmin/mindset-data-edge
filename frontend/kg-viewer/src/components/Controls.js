export default class Controls {
    constructor(containerId) {
        this.containerId = containerId
        this.onFit = null
        this.onZoomIn = null
        this.onZoomOut = null
        this.onReset = null
        this.onRefresh = null
    }
    
    show() {
        const container = document.getElementById(this.containerId)
        if (!container) return
        
        const controlsHtml = `
            <div class="controls">
                <button id="ctrl-fit" title="Ajuster la vue">🔍 Fit</button>
                <button id="ctrl-zoom-in" title="Zoom +">➕</button>
                <button id="ctrl-zoom-out" title="Zoom -">➖</button>
                <button id="ctrl-reset" title="Réinitialiser">🔄 Reset</button>
                <button id="ctrl-refresh" title="Rafraîchir">⬇️ Refresh</button>
            </div>
        `
        
        const existingControls = container.querySelector('.controls')
        if (existingControls) {
            existingControls.outerHTML = controlsHtml
        } else {
            const logo = container.querySelector('.logo')
            if (logo) {
                logo.insertAdjacentHTML('afterend', controlsHtml)
            } else {
                container.insertAdjacentHTML('beforeend', controlsHtml)
            }
        }
        
        this.attachEvents()
    }
    
    attachEvents() {
        const fitBtn = document.getElementById('ctrl-fit')
        const zoomInBtn = document.getElementById('ctrl-zoom-in')
        const zoomOutBtn = document.getElementById('ctrl-zoom-out')
        const resetBtn = document.getElementById('ctrl-reset')
        const refreshBtn = document.getElementById('ctrl-refresh')
        
        if (fitBtn) fitBtn.addEventListener('click', () => this.onFit?.())
        if (zoomInBtn) zoomInBtn.addEventListener('click', () => this.onZoomIn?.())
        if (zoomOutBtn) zoomOutBtn.addEventListener('click', () => this.onZoomOut?.())
        if (resetBtn) resetBtn.addEventListener('click', () => this.onReset?.())
        if (refreshBtn) refreshBtn.addEventListener('click', () => this.onRefresh?.())
    }
}