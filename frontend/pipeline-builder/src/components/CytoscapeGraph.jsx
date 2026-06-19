import { useEffect, useRef } from 'react';
import cytoscape from 'cytoscape';

// Thin React wrapper around Cytoscape. Re-initializes when elements/stylesheet/
// layout change; node taps are reported via a ref so the callback identity
// doesn't force a re-init.
export default function CytoscapeGraph({ elements, stylesheet, layout, onNodeSelect }) {
  const containerRef = useRef(null);
  const selectRef = useRef(onNodeSelect);
  selectRef.current = onNodeSelect;

  useEffect(() => {
    if (!containerRef.current) return undefined;

    const cy = cytoscape({
      container: containerRef.current,
      elements,
      style: stylesheet,
      layout,
      minZoom: 0.2,
      maxZoom: 2.5,
    });

    cy.on('tap', 'node', (e) => selectRef.current && selectRef.current(e.target.data()));
    cy.on('tap', (e) => {
      if (e.target === cy) selectRef.current && selectRef.current(null);
    });

    const fit = setTimeout(() => cy.fit(undefined, 40), 300);

    return () => {
      clearTimeout(fit);
      cy.destroy();
    };
  }, [elements, stylesheet, layout]);

  return <div ref={containerRef} className="w-full h-full" />;
}
