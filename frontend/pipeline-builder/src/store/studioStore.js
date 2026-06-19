import { create } from 'zustand';

// Cross-page intents: Connect/Pipelines pages set these, the Compose page
// consumes them on mount and clears them.
export const useStudioStore = create((set) => ({
  pendingConnector: null, // a connector FunctionInfo to apply to the trigger node
  pipelineToLoad: null, // a pipeline id to load onto the canvas

  selectConnector: (fn) => set({ pendingConnector: fn }),
  requestLoadPipeline: (id) => set({ pipelineToLoad: id }),
  clearPending: () => set({ pendingConnector: null, pipelineToLoad: null }),
}));
