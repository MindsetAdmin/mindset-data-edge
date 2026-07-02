import { create } from 'zustand';

// Cross-page intents: Connect/Pipelines pages set these, the Compose page
// consumes them on mount and clears them.
export const useStudioStore = create((set) => ({
  pendingConnector: null, // a connector FunctionInfo to apply to the trigger node
  pipelineToLoad: null, // a full pipeline object to load onto the canvas

  // OPC-UA selections the user applied (node_id -> mode). The builder uses these
  // to offer only function-eligible (isa95/both) tags in function field pickers.
  opcuaSelections: [],

  // Canvas state persisted within a browser session (in-memory Zustand — resets
  // on page refresh, survives tab switches). null = use fresh initial state.
  canvasNodes: null,
  canvasEdges: null,
  canvasMeta: { id: '', name: '', description: '' },

  // Machines referenced in state_machine nodes of the current session's pipeline.
  // The Dashboard uses this to show only the machines the user has actively selected.
  selectedMachines: [],

  selectConnector: (fn) => set({ pendingConnector: fn }),
  requestLoadPipeline: (pipeline) => set({ pipelineToLoad: pipeline }),
  setOpcuaSelections: (selections) => set({ opcuaSelections: selections }),
  clearPending: () => set({ pendingConnector: null, pipelineToLoad: null }),

  saveCanvasState: (nodes, edges, meta) => set({ canvasNodes: nodes, canvasEdges: edges, canvasMeta: meta }),
  // clearCanvas also resets selectedMachines so a "New" pipeline starts fully clean.
  clearCanvas: () => set({
    canvasNodes: null,
    canvasEdges: null,
    canvasMeta: { id: '', name: '', description: '' },
    selectedMachines: [],
  }),
  addSelectedMachine: (id) => set((s) => ({
    selectedMachines: s.selectedMachines.includes(id) ? s.selectedMachines : [...s.selectedMachines, id],
  })),
}));
