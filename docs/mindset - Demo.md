MindSet Data - Pipeline Studio: Complete Implementation Plan
🎯 Project Overview
Build a visual pipeline builder interface (inspired by MaestroHub) that allows users to:

Select a connector (OPC-UA, MQTT, Modbus, SQL)

Build pipelines by either:

Loading pre-defined pipelines (1-click)

Drag & dropping individual functions onto a canvas

Visualize results through:

Real-time dashboards

Interactive Knowledge Graph showing relationships between Connectors, Functions, Topics, Pipelines, and Dashboards

📊 Core Architecture
Application Layout
text
┌──────────────────────────────────────────────────────────────────────────────────────┐
│  🧩 MINDSET DATA - Studio                                                            │
├──────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                      │
│  📊 Overview    🔌 Connect    ⚙️ Compose    📡 Pipelines    📊 Dashboards    📈 KG   │
│                                                                                      │
├──────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                      │
│  ┌─────────────── PALETTE (gauche) ──────────────────┐  ┌─── CANVAS (centre) ──────┐│
│  │                                                    │  │                          ││
│  │  🎯 CONNECTEURS                                    │  │  ┌────────────────────┐  ││
│  │  ┌──────────────────────────────┐                  │  │  │                    │  ││
│  │  │  🔌 OPC-UA     [Sélectionner] │                  │  │  │   PIPELINE VISUEL  │  ││
│  │  │  📡 MQTT       [Sélectionner] │                  │  │  │                    │  ││
│  │  │  🔌 Modbus     [Sélectionner] │                  │  │  │   ┌──┐  ┌──┐  ┌──┐ │  ││
│  │  │  🗄️ SQL        [Sélectionner] │                  │  │  │   │OP│──│FI│──│MA│ │  ││
│  │  └──────────────────────────────┘                  │  │  │   └──┘  └──┘  └──┘ │  ││
│  │                                                    │  │  │                    │  ││
│  │  ⚙️ FONCTIONS                                      │  │  │   ┌──┐  ┌──┐      │  ││
│  │  ┌──────────────────────────────┐                  │  │  │   │DU│──│TH│      │  ││
│  │  │  🔍 Filter                    │                  │  │  │   └──┘  └──┘      │  ││
│  │  │  🔄 State Machine             │                  │  │  │                    │  ││
│  │  │  🗺️ UNS Mapper                │                  │  │  │   ┌──┐  ┌──┐      │  ││
│  │  │  ⏱️ Duration                  │                  │  │  │   │MQ│──│KG│      │  ││
│  │  │  💰 Cost                      │                  │  │  │   └──┘  └──┘      │  ││
│  │  │  🚦 Threshold                 │                  │  │  │                    │  ││
│  │  │  📤 MQTT Publish              │                  │  │  └────────────────────┘  ││
│  │  │  💾 KG Save                   │                  │  │                          ││
│  │  └──────────────────────────────┘                  │  │  [▶️ Run] [💾 Save]     ││
│  │                                                    │  └──────────────────────────┘│
│  │  📋 PIPELINES PRÉ-DÉFINIES                         │                            │
│  │  ┌──────────────────────────────┐                  │  ┌─── PANEL CONFIG (droite) ─┐│
│  │  │  🔴 Micro-stop Detection     │                  │  │                           ││
│  │  │  🔄 OPC-UA → UNS             │                  │  │  Nœud: State Machine      ││
│  │  │  💰 Cost Calculation         │                  │  │  Machine ID: [machine1   ] ││
│  │  └──────────────────────────────┘                  │  │                           ││
│  │                                                    │  │  [Appliquer]              ││
│  └────────────────────────────────────────────────────┘  └───────────────────────────┘│
│                                                                                      │
└──────────────────────────────────────────────────────────────────────────────────────┘
🔄 User Workflow (3 Steps)
Step 1: Select Connector
User selects a connector from the palette (OPC-UA, MQTT, Modbus, SQL)

Connector appears as the first node on the canvas

Configuration form opens (endpoint, port, parameters)

Step 2: Build Pipeline (Fusion of Functions + Pre-defined)
User can build their pipeline in TWO ways:

Option A: Load Pre-defined Pipeline (Drag and drop)

text
📋 PIPELINES PRÉ-DÉFINIES
┌──────────────────────────────────────────────────┐
│  🔴 Micro-stop Detection    [Charger]  ← Click  │
└──────────────────────────────────────────────────┘

→ Complete pipeline appears on canvas with all nodes configured
→ 6 nodes: OPC-UA → State Machine → Duration → Threshold → MQTT Publish → KG Save
→ User can modify any node
Option B: Build Manually (Drag & Drop)

text
⚙️ FONCTIONS
┌──────────────────────────────────────────────────┐
│  🔍 Filter      [Glisser]  ← Drag to canvas     │
│  🔄 State Mach.  [Glisser]                       │
│  ⏱️ Duration    [Glisser]                       │
└──────────────────────────────────────────────────┘

→ Functions appear as nodes on canvas
→ User connects them manually
→ Each node is configurable
Step 3: Visualize & Dashboard
▶️ Run Pipeline → Real-time execution

📊 Knowledge Graph → View relationships between all components

📈 Dashboard → View real-time metrics

📊 Knowledge Graph Node Types
Node Type	Color	Icon	Description
Connection	🔵 Blue	🔌	Data source/destination (OPC-UA, MQTT, Modbus, SQL)
Function	🟠 Orange	⚙️	Action element (Filter, State Machine, UNS Mapper, Duration, Cost, Threshold, MQTT Publish, KG Save)
Topic	🟢 Green	📡	MQTT channel (mindset/raw/#, mindset/site/#, mindset/events/#)
Pipeline	🔴 Red	🔧	Suite of functions (shows relationships, NOT internal composition)
Dashboard	🟣 Purple	📊	Data visualization
Important: Pipeline Node Behavior
The Pipeline node in KG shows ONLY external relationships:

Which Connections it uses

Which Topics it consumes

Which Topics it produces

Which Dashboards use it

It does NOT show its internal composition (the functions inside)

📁 Pre-defined Pipelines
Pipeline	Nodes	Description
Micro-stop Detection	OPC-UA → State Machine → Duration → Threshold → MQTT Publish → KG Save	Detects micro-stops (30s-3min)
OPC-UA → UNS	OPC-UA → UNS Mapper → MQTT Publish	Transforms OPC-UA to UNS ISA-95
Cost Calculation	MQTT Subscribe → Cost → MQTT Publish → KG Save	Calculates micro-stop costs in euros
📡 Available Functions
Category	Function	Description
Connector	opcua_read	Read from OPC-UA
Connector	mqtt_subscribe	Subscribe to MQTT topic
Transform	filter	Filter data by condition
Transform	state_machine	Detect Run/Stop transitions
Transform	uns_mapper	Map to UNS ISA-95
Calculate	duration	Calculate duration
Calculate	cost	Calculate cost in euros
Condition	threshold	Check if value is between min/max
Output	mqtt_publish	Publish to MQTT
Output	kg_save	Save to Knowledge Graph

🎨 Frontend Architecture

- Framework: React + Vite
- Pipeline Canvas: React Flow (XY Flow)
- Knowledge Graph: Cytoscape.js (existing)
- Drag & Drop: @dnd-kit/core
- Styling: Tailwind CSS (same theme as kg-viewer)
- State Management: Zustand
- HTTP Client: Axios

🎯 User Interactions
On Pipeline Canvas (Build View)
Action	Result
Click on connector in palette	Adds connector node to canvas
Drag function from palette	Adds function node to canvas
Click "Load Pipeline"	Pre-defined pipeline appears with all nodes
Click on a node	Opens configuration panel
Drag a node	Repositions the node
Hover on connection	Shows relationship type
Double-click on node	Zooms on node
On Knowledge Graph (KG View)
Action	Result
Click on a node	Shows details in side panel
Hover on a node	Highlights relationships
Double-click on a node	Zooms on node and neighbors
Filter "Connectors"	Shows only connections
Filter "Pipelines"	Shows only pipelines
Click on Pipeline	Shows: "This pipeline depends on N functions"
Filters for KG
All

Connectors

Functions

Topics

Pipelines

Dashboards


🎬 Demo Scenario
text
1. User opens Pipeline Studio
   → Sees palette with connectors, functions, and pre-defined pipelines

2. User selects "OPC-UA" connector
   → Connector appears on canvas, configuration form opens

3. User clicks "Load Micro-stop Detection" pipeline
   → Complete pipeline appears: OPC-UA → State Machine → Duration → Threshold → MQTT Publish → KG Save

4. User clicks on "State Machine" node
   → Configuration panel opens on right
   → User changes Machine ID to "machine2"
   → Clicks "Apply"

5. User clicks "Run Pipeline"
   → Pipeline executes in real-time

6. User clicks "Dashboard" tab
   → Sees real-time metrics: 3 micro-stops detected, 63.75€ cost

7. User clicks "Knowledge Graph" tab
   → Sees all relationships:
   - 🔵 OPC-UA → 📡 mindset/raw/# → 🔧 Pipeline → ⚙️ Functions → 📡 events
   - Pipeline node shows dependencies: "Depends on 6 functions"

8. User can navigate freely between views:
   Connect → Compose → Pipelines → Dashboards → KG

🎯 Success Criteria
User can select a connector (OPC-UA, MQTT, Modbus, SQL)

User can load a pre-defined pipeline (1-click)

User can build a pipeline manually (drag & drop functions)

User can configure each node (parameters)

User can run a pipeline and see results

User can see real-time dashboard

User can see Knowledge Graph with 5 node types (Connectors, Functions, Topics, Pipelines, Dashboards)

Pipeline node shows ONLY external relationships

Smooth navigation between views (Connect → Compose → Pipelines → Dashboards → KG)

Responsive design with dark theme (consistent with kg-viewer)