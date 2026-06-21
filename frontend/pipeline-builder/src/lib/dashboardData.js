// Transforms the domain Knowledge Graph (/api/kg/domain) into dashboard data:
// events joined with their cause + cost, today/yesterday buckets, and a Pareto.

export function buildEvents(domain) {
  const nodes = domain?.nodes || [];
  const edges = domain?.edges || [];
  const byId = {};
  nodes.forEach((n) => (byId[n.id] = n));

  const causeByEvent = {};
  const costByEvent = {};
  edges.forEach((e) => {
    if (e.relation === 'caused_by') causeByEvent[e.from_id] = byId[e.to_id]?.label;
    if (e.relation === 'costs') costByEvent[e.from_id] = byId[e.to_id]?.properties?.amount_eur;
  });

  return nodes
    .filter((n) => n.type === 'Event')
    .map((n) => ({
      id: n.id,
      workCenter: n.properties?.work_center || '—',
      duration: Number(n.properties?.duration_seconds || 0),
      createdAt: n.created_at,
      cause: causeByEvent[n.id] || null,
      cost: costByEvent[n.id] ?? null,
    }))
    .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
}

export function effectiveCost(ev, hourly) {
  return ev.cost != null ? ev.cost : (ev.duration / 3600) * hourly;
}

function sameDay(a, b) {
  return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate();
}

// Split events into today / yesterday by created_at.
export function splitDays(events) {
  const now = new Date();
  const yest = new Date(now);
  yest.setDate(now.getDate() - 1);
  const today = [];
  const yesterday = [];
  events.forEach((e) => {
    if (!e.createdAt) return;
    const d = new Date(e.createdAt);
    if (sameDay(d, now)) today.push(e);
    else if (sameDay(d, yest)) yesterday.push(e);
  });
  return { today, yesterday };
}

// Percentage change (today vs yesterday). null when no yesterday baseline.
export function deltaPct(todayVal, yesterdayVal) {
  if (!yesterdayVal) return null;
  return ((todayVal - yesterdayVal) / yesterdayVal) * 100;
}

// Pareto of causes: count per cause, sorted desc, with cumulative %.
export function paretoCauses(events) {
  const counts = {};
  events.forEach((e) => {
    const c = e.cause || 'Non catégorisé';
    counts[c] = (counts[c] || 0) + 1;
  });
  const total = events.length || 1;
  const sorted = Object.entries(counts).sort((a, b) => b[1] - a[1]);
  let cum = 0;
  return sorted.map(([cause, count]) => {
    cum += count;
    return { cause, count, cumulative: Math.round((cum / total) * 100) };
  });
}
