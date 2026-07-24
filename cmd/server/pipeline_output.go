// cmd/server/pipeline_output.go
// Auto-publish for pipeline output — replaces the old explicit mqtt_publish
// node (docs/analysis_log.md Entry 119). A pipeline no longer needs a node
// whose only job is calling mqtt_publish; the engine's caller (this file)
// publishes the declared Output node's result automatically after a
// successful run.
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/MindsetAdmin/mindset-data-edge/internal/pipeline"
)

// publishPipelineOutput publishes p's declared Output node's result to MQTT,
// once, after a successful run. No-op if the pipeline didn't succeed, didn't
// declare an Output node, that node didn't itself succeed, or no MQTT client
// is connected.
//
// Topic resolution: p.OutputTopic if set — required for pipelines whose topic
// name is load-bearing (chained into another pipeline's trigger topic, or
// consumed by internal/kg/subscriber.go's hardcoded mindset/events/micro-stop
// subscription) — otherwise an auto-derived default
// (mindset/pipelines/<id>/output).
func (s *server) publishPipelineOutput(p *pipeline.Pipeline, result *pipeline.ExecutionResult) {
	if s.mqttClient == nil || result.Status != pipeline.StatusSuccess || p.Output == "" {
		return
	}

	var output interface{}
	found := false
	for _, n := range result.Nodes {
		if n.NodeID == p.Output && n.Status == pipeline.StatusSuccess {
			output = n.Output
			found = true
			break
		}
	}
	if !found {
		return
	}

	topic := p.OutputTopic
	if topic == "" {
		topic = fmt.Sprintf("mindset/pipelines/%s/output", p.ID)
	}

	data, err := json.Marshal(output)
	if err != nil {
		log.Printf("[API] Pipeline %q: failed to marshal output for auto-publish: %v", p.ID, err)
		return
	}

	token := s.mqttClient.Publish(topic, 1, false, data)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Printf("[API] Pipeline %q: auto-publish to %q failed: %v", p.ID, topic, err)
		return
	}
	log.Printf("[API] Pipeline %q auto-published output -> %s", p.ID, topic)
}
