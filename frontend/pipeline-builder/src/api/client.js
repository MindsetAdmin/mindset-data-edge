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

export async function deletePipeline(id) {
    const response = await fetch(`${API_BASE}/pipelines/${encodeURIComponent(id)}`, {
        method: 'DELETE',
    });
    if (!response.ok) {
        const body = await response.text();
        throw new Error(body || `Failed to delete pipeline: ${response.statusText}`);
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

// Live active production per machine, from a validated ERP work-order
// mapping (Entry 120) — each fact carries equipment_id when the IT-side
// work_center resolved against a real OT Equipment node, "" otherwise.
export async function fetchActiveProduction(workCenter = '') {
    const qs = workCenter ? `?work_center=${encodeURIComponent(workCenter)}` : '';
    const response = await fetch(`${API_BASE}/production/active${qs}`);
    if (!response.ok) throw new Error(`Failed to fetch active production: ${response.statusText}`);
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
// selections: [{ node_id, mode: 'raw'|'isa95'|'both', area?, work_center?, work_unit?, tag_name? }]
// The 4 optional fields (Entry 124) correct the auto-computed ISA-95 mapping
// shown after discover — leave any of them blank to keep the mapper's guess
// for just that field.
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

// Unified KG endpoint. category: 'business' | 'platform' | 'all' (default 'all').
// Legacy alias mapping is handled server-side, but new code should prefer this.
export async function fetchKG(category = 'all') {
    const response = await fetch(`${API_BASE}/kg?category=${category}`);
    if (!response.ok) {
        throw new Error(`Failed to fetch KG (${category}): ${response.statusText}`);
    }
    return response.json();
}

// Legacy shim — kept for backward compatibility. Maps 'technical'→'platform' and
// 'domain'→'business'. New code should call fetchKG(category) directly.
export async function fetchKnowledgeGraph(kind = 'technical') {
    const category = kind === 'technical' ? 'platform' : kind === 'domain' ? 'business' : kind;
    return fetchKG(category);
}

// --- KG structural bootstrap validation (v0 — docs/analysis_log.md Entries 95/96) --

// fetchPendingKGNodes returns business-category nodes auto-generated from OPC-UA
// discovery that haven't been human-validated yet.
export async function fetchPendingKGNodes() {
    const response = await fetch(`${API_BASE}/kg/pending`);
    if (!response.ok) throw new Error(`Failed to fetch pending KG nodes: ${response.statusText}`);
    return response.json(); // { nodes: [...], total }
}

export async function validateKGNode(id) {
    const response = await fetch(`${API_BASE}/kg/pending/${encodeURIComponent(id)}/validate`, { method: 'POST' });
    if (!response.ok) throw new Error(`Failed to validate node: ${response.statusText}`);
    return response.json();
}

export async function rejectKGNode(id) {
    const response = await fetch(`${API_BASE}/kg/pending/${encodeURIComponent(id)}/reject`, { method: 'POST' });
    if (!response.ok) throw new Error(`Failed to reject node: ${response.statusText}`);
    return response.json();
}

// --- SQL Connections (internal/connections) ---------------------------------

export async function fetchConnections() {
    const response = await fetch(`${API_BASE}/connections`);
    if (!response.ok) throw new Error(`Failed to fetch connections: ${response.statusText}`);
    return response.json(); // { connections: [...], total }
}

// cfg: { id, name, driver, host, port, database, username, password_env, tls,
//        read_timeout_seconds, write_timeout_seconds, max_open_conns,
//        max_idle_conns, conn_max_lifetime_seconds }
export async function createConnection(cfg) {
    const response = await fetch(`${API_BASE}/connections`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(cfg),
    });
    if (!response.ok) {
        const text = await response.text();
        throw new Error(text || `Failed to save connection: ${response.statusText}`);
    }
    return response.json();
}

// testConnection always resolves — a failed health check comes back as
// {ok:false, error} rather than an HTTP error.
export async function testConnection(id) {
    const response = await fetch(`${API_BASE}/connections/${encodeURIComponent(id)}/test`, { method: 'POST' });
    if (!response.ok) throw new Error(`Failed to test connection: ${response.statusText}`);
    return response.json(); // { ok, latency_ms, read_only?, error? }
}

// fetchConnectionDatabases browses every database + table visible to this
// connection's user in one call (scoped by the account's real MySQL grants —
// not necessarily every database on the server).
export async function fetchConnectionDatabases(id) {
    const response = await fetch(`${API_BASE}/connections/${encodeURIComponent(id)}/databases`);
    const body = await response.json().catch(() => ({}));
    if (!response.ok) throw new Error(body.error || `Failed to list databases: ${response.statusText}`);
    return body; // { databases: [{ name, tables: [{ name, columns: [...] }] }], total }
}

// preview: { query, params, limit } — runs through the same guards as
// sql_query, capped server-side at 5 rows.
export async function previewConnection(id, preview) {
    const response = await fetch(`${API_BASE}/connections/${encodeURIComponent(id)}/preview`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(preview),
    });
    if (!response.ok) {
        const text = await response.text();
        throw new Error(text || `Failed to preview: ${response.statusText}`);
    }
    return response.json(); // { rows, canonical, canonical_type, row_count, query_ms }
}

export async function deleteConnection(id) {
    const response = await fetch(`${API_BASE}/connections/${encodeURIComponent(id)}`, { method: 'DELETE' });
    if (!response.ok) {
        const text = await response.text();
        throw new Error(text || `Failed to delete connection: ${response.statusText}`);
    }
    return response.json();
}