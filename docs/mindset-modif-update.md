Objective
Transform the Dashboard's add_to_dashboard feature into a fully functional, interactive widget system where users can:

See live data from add_to_dashboard as interactive widgets

Visualize data with graphs (time series, bar charts, gauges)

Select what to display via a widget selector/manager

Customize each widget (size, chart type, refresh rate)

🔍 Current Issue
What I See Now
text
📌 Widgets épinglés
● live
Coût micro-stop
63.75
09:50:07

Mon widget
{"kind":"value","label":"Mon widget"}
15:39:20

Température machine1
23.7
01:01:05
Problems
Issue	Description
No graphs	Only raw text/numbers displayed
No chart selection	Cannot choose chart type (line, bar, gauge, etc.)
No data selection	Cannot choose which data to display
Raw JSON visible	Shows {"kind":"value","label":"Mon widget"} instead of parsed data
No widget management	Cannot add, remove, or rearrange widgets
No time range	Cannot see historical data, only latest value
📊 Required Solution
User Experience Goals
text
┌─────────────────────────────────────────────────────────────────────────────────────┐
│  📌 Widgets épinglés                                                  [+ Ajouter]  │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  ┌───────────────────────────────────────────────────────────────────────────────┐ │
│  │  Widget Selector (choose what to display)                                     │ │
│  │  ┌─────────────────────────────────────────────────────────────────────────┐ │ │
│  │  │  ○ Coût micro-stop    ● Température machine1    ○ Status machine1      │ │ │
│  │  │  ○ Vitesse moteur     ○ Pression               ○ Compteur              │ │ │
│  │  └─────────────────────────────────────────────────────────────────────────┘ │ │
│  │  Chart Type: [Line ▼]  Time Range: [Last Hour ▼]  Refresh: [5s ▼]            │ │
│  └───────────────────────────────────────────────────────────────────────────────┘ │
│                                                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │  📊 Température machine1                                      [✕] [⚙️]    │   │
│  │  ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │  │  25 │    ╭╮                                                       │   │   │
│  │  │  24 │   ╭╯╰╮    ╭╮                                                │   │   │
│  │  │  23 │  ╭╯  ╰╮  ╭╯╰╮                                              │   │   │
│  │  │  22 │ ╭╯    ╰╮╭╯  ╰╮                                             │   │   │
│  │  │     ├──┼──┼──┼──┼──┼──┤                                          │   │   │
│  │  │     14:00  14:05  14:10  14:15                                    │   │   │
│  │  └─────────────────────────────────────────────────────────────────────┘   │   │
│  │  Dernière valeur: 23.7°C  |  Min: 22.1°C  |  Max: 24.8°C  |  Moy: 23.4°C │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                     │
│  ┌─────────────────────────────────────────┐  ┌─────────────────────────────────┐   │
│  │  💰 Coût micro-stop                     │  │  🏭 Status machine1             │   │
│  │  ┌───────────────────────────────────┐  │  │  ┌───────────────────────────┐  │   │
│  │  │  80 │                           │  │  │  │  🟢 Running                │  │   │
│  │  │  60 │    ╭──╮                   │  │  │  │  Dernier changement:        │  │   │
│  │  │  40 │   ╭╯  ╰╮                  │  │  │  │  14:32:05                  │  │   │
│  │  │  20 │  ╭╯    ╰──╮              │  │  │  │  Durée: 45s                │  │   │
│  │  │     ├──┼──┼──┼──┼──┤           │  │  │  └───────────────────────────┘  │   │
│  │  │     14:00  14:05  14:10        │  │  │                                  │   │
│  │  └───────────────────────────────────┘  │  │  [Détails]                    │   │
│  │  Dernier: 63.75€  |  Total: 247.50€   │  │  └─────────────────────────────────┘   │
│  └─────────────────────────────────────────┘                                      │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
🎯 Acceptance Criteria
Users can add widgets from available data sources

Each widget shows a graph/chart (Line, Bar, Gauge, Value, Status)

Users can select chart type when adding a widget

Users can choose time range (1m, 5m, 1h, 4h, 24h)

Users can set refresh rate (1s, 5s, 10s, 30s)

Widgets show statistics (Min, Max, Average, Last, Count)

Widgets are live (update via WebSocket)

Widgets are persistent (saved in localStorage)

Users can close widgets (✕ button)

Users can configure widgets (⚙️ button)

Raw JSON is never displayed (only parsed data)

Dashboard has a clean, professional look