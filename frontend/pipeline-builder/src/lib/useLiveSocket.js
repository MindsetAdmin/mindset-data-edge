import { useEffect, useRef, useState } from 'react';

// Connects to the server's WebSocket (/api/ws) for real-time push and calls
// onMessage({type,data}) for each message. Auto-reconnects with backoff.
// Returns the live connection status.
export function useLiveSocket(onMessage) {
  const [connected, setConnected] = useState(false);
  const cbRef = useRef(onMessage);
  cbRef.current = onMessage;

  useEffect(() => {
    let ws = null;
    let retry = null;
    let attempts = 0;
    let closed = false;

    const url = `${location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/api/ws`;

    const connect = () => {
      ws = new WebSocket(url);
      ws.onopen = () => {
        attempts = 0;
        setConnected(true);
      };
      ws.onmessage = (e) => {
        try {
          const msg = JSON.parse(e.data);
          cbRef.current && cbRef.current(msg);
        } catch {
          /* ignore malformed */
        }
      };
      ws.onclose = () => {
        setConnected(false);
        if (closed) return;
        attempts += 1;
        const delay = Math.min(1000 * attempts, 5000);
        retry = setTimeout(connect, delay);
      };
      ws.onerror = () => ws && ws.close();
    };

    connect();
    return () => {
      closed = true;
      clearTimeout(retry);
      if (ws) ws.close();
    };
  }, []);

  return connected;
}
