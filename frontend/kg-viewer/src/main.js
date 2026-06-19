// src/main.js
import cytoscape from 'cytoscape'
import dagre from 'cytoscape-dagre'
import fcose from 'cytoscape-fcose'
import App from './App'

// Register extensions
cytoscape.use(dagre)
cytoscape.use(fcose)

// Initialize app
const app = new App()
app.init()