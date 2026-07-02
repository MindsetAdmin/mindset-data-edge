// Pipeline loading helpers — converts a saved pipeline to a ReactFlow graph,
// stripping node configs so the user reconfigures inputs/outputs each session.

import { pipelineToFlow } from './pipelineMapping';
import { defaultFunctionConfig } from './functionDefaults';
import { defaultConfigFor } from './connectorTemplates';

// Load only the function chain from a saved pipeline: preserves node types and
// connections but seeds each node with its default config rather than the saved
// one. Trigger config is also reset. The user must reconfigure before saving.
export function pipelineToFlowChainOnly(p) {
  const flow = pipelineToFlow(p);
  flow.nodes = flow.nodes.map((n) => {
    if (n.type === 'pipelineNode') {
      return { ...n, data: { ...n.data, config: defaultFunctionConfig(n.data.function) } };
    }
    if (n.type === 'triggerNode') {
      return { ...n, data: { ...n.data, config: defaultConfigFor(n.data.function) } };
    }
    return n; // zone nodes, etc.
  });
  return flow;
}
