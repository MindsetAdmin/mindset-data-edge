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