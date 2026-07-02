[plugin:vite:import-analysis] Failed to resolve import "lucide-react" from "src/components/DashboardWidgets.jsx". Does the file exist?
C:/Users/khena/Desktop/MINDSET/Project/mindset-data-edge/frontend/pipeline-builder/src/components/DashboardWidgets.jsx:3:90
1  |  import { useState, useEffect, useRef } from "react";
2  |  import { ResponsiveContainer, LineChart, Line, BarChart, Bar, XAxis, YAxis, Tooltip, CartesianGrid } from "recharts";
3  |  import { Settings, X, Plus, LineChart as LineIcon, BarChart3, Gauge, Hash, Factory } from "lucide-react";
   |                                                                                             ^
4  |  import { useLiveSocket } from "../lib/useLiveSocket";
5  |  import { fetchDashboardPins } from "../api/client";