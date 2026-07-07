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

// groupByMachine — aggregate events per work_center for a per-machine dashboard
// breakdown (2026-07-02). Returns rows with the same shape as the KPI cards but
// scoped to a single machine, plus a joined view of the machines list (so
// machines with no events still appear as "idle / 0 stops").
//
//   events: [{ workCenter, duration, cost, createdAt }]
//   machines: [{ work_center, state?: {running?, history?}, tags[] }]
//   hourly: default €/h for cost fallback
//   shiftSeconds: reference window for availability (default 8h)
export function groupByMachine(events, machines, hourly, shiftSeconds = 8 * 3600) {
  const { today, yesterday } = splitDays(events);
  const byMachineToday = bucketize(today);
  const byMachineYest = bucketize(yesterday);

  const wcs = new Set();
  (machines || []).forEach((m) => wcs.add(m.work_center));
  Object.keys(byMachineToday).forEach((wc) => wcs.add(wc));

  const machinesByWc = {};
  (machines || []).forEach((m) => (machinesByWc[m.work_center] = m));

  const rows = [];
  wcs.forEach((wc) => {
    if (wc === '(autres)' || !wc) return;
    const t = byMachineToday[wc] || { count: 0, downtime: 0, cost: 0 };
    const y = byMachineYest[wc] || { count: 0, downtime: 0, cost: 0 };
    const totalCostToday = t.cost || (t.downtime / 3600) * hourly;
    const totalCostYest = y.cost || (y.downtime / 3600) * hourly;
    const availability =
      Math.max(0, Math.min(1, 1 - t.downtime / shiftSeconds)) * 100;
    const m = machinesByWc[wc];
    rows.push({
      workCenter: wc,
      state: m?.state || null,
      running: m?.state?.running ?? null,
      tagCount: m?.tags?.length || 0,
      stopsToday: t.count,
      stopsYest: y.count,
      downtimeToday: t.downtime,
      downtimeYest: y.downtime,
      costToday: totalCostToday,
      costYest: totalCostYest,
      availability,
    });
  });
  rows.sort((a, b) => b.costToday - a.costToday || a.workCenter.localeCompare(b.workCenter));
  return rows;

  function bucketize(list) {
    const out = {};
    list.forEach((e) => {
      const wc = e.workCenter || '—';
      if (!out[wc]) out[wc] = { count: 0, downtime: 0, cost: 0 };
      out[wc].count++;
      out[wc].downtime += e.duration || 0;
      if (e.cost != null) out[wc].cost += e.cost;
    });
    // Fill missing cost using duration × hourly (fallback).
    Object.values(out).forEach((agg) => {
      if (agg.cost === 0 && agg.downtime > 0) {
        agg.cost = (agg.downtime / 3600) * hourly;
      }
    });
    return out;
  }
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
