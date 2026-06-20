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

export async function fetchStats() {
    const response = await fetch(`${API_BASE}/stats`);
    if (!response.ok) {
        throw new Error(`Failed to fetch stats: ${response.statusText}`);
    }
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