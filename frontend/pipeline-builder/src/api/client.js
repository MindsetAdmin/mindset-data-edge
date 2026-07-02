const API_BASE = '/api';

export async function fetchFunctions(type = null) {
    const url = type ? `${API_BASE}/functions?type=${type}` : `${API_BASE}/functions`;
    const response = await fetch(url);
    if (!response.ok) {
        throw new Error(`Failed to fetch functions: ${response.statusText}`);
    }
    return response.json();
}

export async function fetchConnectors() {
    const response = await fetch(`${API_BASE}/connectors`);
    if (!response.ok) {
        throw new Error(`Failed to fetch connectors: ${response.statusText}`);
    }
    return response.json();
}

export async function createPipeline(pipelineData) {
    const response = await fetch(`${API_BASE}/pipelines`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(pipelineData),
    });
    if (!response.ok) {
        throw new Error(`Failed to create pipeline: ${response.statusText}`);
    }
    return response.json();
}

export async function fetchPipelines() {
    const response = await fetch(`${API_BASE}/pipelines`);
    if (!response.ok) {
        throw new Error(`Failed to fetch pipelines: ${response.statusText}`);
    }
    return response.json();
}

export async function fetchExamplePipelines() {
    const response = await fetch(`${API_BASE}/pipelines/examples`);
    if (!response.ok) {
        throw new Error(`Failed to fetch example pipelines: ${response.statusText}`);
    }
    return response.json();
}

export async function runPipeline(id) {
    const response = await fetch(`${API_BASE}/pipelines/${encodeURIComponent(id)}/run`, {
        method: 'POST',
    });
    if (!response.ok) {
        const body = await response.text();
        throw new Error(body || `Failed to run pipeline: ${response.statusText}`);
    }
    return response.json();
}

export async function fetchTags() {
    const response = await fetch(`${API_BASE}/tags`);
    if (!response.ok) {
        throw new Error(`Failed to fetch tags: ${response.statusText}`);
    }
    return response.json();
}

export async function fetchMachines() {
    const response = await fetch(`${API_BASE}/machines`);
    if (!response.ok) throw new Error(`Failed to fetch machines: ${response.statusText}`);
    return response.json();
}

export async function fetchTopics() {
    const response = await fetch(`${API_BASE}/topics`);
    if (!response.ok) throw new Error(`Failed to fetch topics: ${response.statusText}`);
    return response.json();
}

export async function fetchConfig() {
    const response = await fetch(`${API_BASE}/config`);
    if (!response.ok) throw new Error(`Failed to fetch config: ${response.statusText}`);
    return response.json();
}

export async function fetchDashboardPins() {
    const response = await fetch(`${API_BASE}/dashboard/pins`);
    if (!response.ok) throw new Error(`Failed to fetch dashboard pins: ${response.statusText}`);
    return response.json();
}

export async function fetchStats() {
    const response = await fetch(`${API_BASE}/stats`);
    if (!response.ok) {
        throw new Error(`Failed to fetch stats: ${response.statusText}`);
    }
    return response.json();
}

// --- Dynamic OPC-UA control plane -----------------------------------------

// opcuaConnect connects the server to a user-specified OPC-UA endpoint.
// cfg: { endpoint, security_mode, security_policy, username, password, session_timeout }
export async function opcuaConnect(cfg) {
    const response = await fetch(`${API_BASE}/opcua/connect`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(cfg),
    });
    const body = await response.json().catch(() => ({}));
    if (!response.ok) {
        throw new Error(body.error || `Failed to connect: ${response.statusText}`);
    }
    return body;
}

// opcuaDiscover browses the connected server and returns its tags.
export async function opcuaDiscover() {
    const response = await fetch(`${API_BASE}/opcua/discover`);
    const body = await response.json().catch(() => ({}));
    if (!response.ok) {
        throw new Error(body.error || `Failed to discover tags: ${response.statusText}`);
    }
    return body;
}

// opcuaSubscribe starts monitoring the selected tags.
// selections: [{ node_id, mode: 'raw'|'isa95'|'both' }]
export async function opcuaSubscribe(selections) {
    const response = await fetch(`${API_BASE}/opcua/subscribe`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ selections }),
    });
    const body = await response.json().catch(() => ({}));
    if (!response.ok) {
        throw new Error(body.error || `Failed to subscribe: ${response.statusText}`);
    }
    return body;
}

// opcuaDisconnect closes the current OPC-UA session.
export async function opcuaDisconnect() {
    const response = await fetch(`${API_BASE}/opcua/disconnect`, { method: 'POST' });
    if (!response.ok) throw new Error(`Failed to disconnect: ${response.statusText}`);
    return response.json();
}

// opcuaStatus returns the current connection status.
export async function opcuaStatus() {
    const response = await fetch(`${API_BASE}/opcua/status`);
    if (!response.ok) throw new Error(`Failed to fetch OPC-UA status: ${response.statusText}`);
    return response.json();
}

// fetchOpcuaSelections returns the current per-tag routing with ISA-95 mapping.
// Used by the builder to restrict function field pickers to isa95/both tags.
export async function fetchOpcuaSelections() {
    const response = await fetch(`${API_BASE}/opcua/selections`);
    if (!response.ok) throw new Error(`Failed to fetch OPC-UA selections: ${response.statusText}`);
    return response.json();
}

// kind: 'technical' | 'domain'
export async function fetchKnowledgeGraph(kind = 'technical') {
    const response = await fetch(`${API_BASE}/kg/${kind}`);
    if (!response.ok) {
        throw new Error(`Failed to fetch ${kind} graph: ${response.statusText}`);
    }
    return response.json();
}