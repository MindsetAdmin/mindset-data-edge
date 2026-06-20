Objective
Bridge the gap between the Edge Agent (backend, Go) and the Frontend (React/Next.js) so that the user sees their real machines, tags, and topics — not static mock data.

📊 Current State vs Target State
Current (Mock data)
text
Frontend shows:
- Static list of machines (machine1, machine2)
- Predefined topics (mindset/raw/ns=3;i=1011)
- Hardcoded pipeline templates
- No real-time data sync
Target (Real data)
text
Frontend shows:
- ✅ Actual machines discovered by OPC-UA
- ✅ Live topics from MQTT broker
- ✅ Real-time tag values
- ✅ Actual pipelines from config/pipelines/
- ✅ Live events from mindset/events/#


The Frontend should feel like you are manipulating the backend directly — every selection in the UI maps to real data flowing through the Edge Agent. 




1. 🔌 OPC-UA Read
What the user sees in Frontend:

List of all discovered machines with their live tags

Each tag shows: Node ID, Name, Data Type, Current Value

Configuration panel with endpoint, node selection, timeout, security mode

Frontend UI Element	What It Displays	Where It Comes From (Edge)
List of machines	Machine1, Machine2, Simulation...	discovery/opcua.go → BrowseNodeTree() → Extracts WorkCenter names from OPC-UA tags
Tags per machine	temperature (Float), status (Boolean), pressure (Float)...	discovery/opcua.go → BrowseNodeTree() → All child tags under each WorkCenter
Node ID	ns=3;i=1011, ns=3;i=1014...	discovery/opcua.go → BrowseNodeTree() → ref.NodeID.NodeID.String()
Data Type	Float, Boolean, Int32, Double...	discovery/opcua.go → readNodeValue() → AttributeIDDataType
Current Value (Live)	23.5, true, 42...	discovery/opcua.go → Subscribe() → Real-time value changes
Timestamp	2026-06-20T14:32:05Z	discovery/opcua.go → readNodeValue() → result.Value.SourceTimestamp
Configuration Panel:

Parameter	Description	Edge Source
Endpoint	OPC-UA Server URL	config/agent.yaml → opcua.endpoint
Node ID	Tag to read	discovery/opcua.go → Selected from BrowseNodeTree()
Timeout	Milliseconds	config/agent.yaml → opcua.timeout (or default 5000ms)
Security Mode	None / Sign / SignAndEncrypt	config/agent.yaml → opcua.security_mode
2. 📡 MQTT Subscribe (Connector)
What the user sees in Frontend:

List of all available MQTT topics (raw, site, events)

Broker connection status

Message rate per topic (msg/s)

Frontend UI Element	What It Displays	Where It Comes From (Edge)
Broker list	tcp://localhost:1883 (pre-configured)	config/agent.yaml → mqtt.broker
Raw topics	mindset/raw/ns=3;i=1009, mindset/raw/ns=3;i=1011...	mqtt/publisher.go → PublishRaw() → Dynamically created from OPC-UA tags
Site topics	mindset/site/local-test/area1/machine1/temperature...	uns/contextualizer.go → Start() → Published after UNS mapping
Event topics	mindset/events/micro-stop, mindset/events/status-change...	rules/engine.go → publishStatusEvent() → Published after detection
Message rate	6 msg/s, 4 msg/s, 0.5 msg/s...	mqtt/publisher.go → Calculated from publish frequency
Broker status	Connected / Disconnected	mqtt/publisher.go → client.IsConnected()
Configuration Panel:

Parameter	Description	Edge Source
Broker URL	MQTT broker address	config/agent.yaml → mqtt.broker
Topic	Topic to subscribe to	List from mqtt/publisher.go (raw/site/events)
QoS	0, 1, 2	Documentation / Default: 1
Client ID	MQTT client identifier	Auto-generated from mqtt/publisher.go → clientID
3. ⚙️ Filter (Transform)
What the user sees in Frontend:

Field selector (which field to filter)

Operator selector (eq, ne, gt, lt, contains)

Value input (what to compare against)

Frontend UI Element	What It Displays	Where It Comes From (Edge)
Available fields	value, name, data_type, node_id...	functions/transforms/filter.go → From incoming JSON structure
Operators	eq (equal), ne (not equal), gt (greater than), lt (less than), contains	functions/transforms/filter.go → Supported operators in code
Dynamic value type	Number input (for gt/lt), Text input (for contains), Boolean (for eq/ne)	functions/transforms/filter.go → Determined by operator and field type
Configuration Panel:

Parameter	Description	Edge Source
Field	Which field to filter on	functions/transforms/filter.go → From JSON structure
Operator	Comparison operator	functions/transforms/filter.go → Supported operators
Value	Value to compare against	User input, validated by functions/transforms/filter.go
4. 🔄 State Machine (Transform)
What the user sees in Frontend:

List of machines with their current status (Running/Stopped)

Transition history (Run → Stop → Run)

Duration since last transition

Frontend UI Element	What It Displays	Where It Comes From (Edge)
Machine list	Machine1 (Running), Machine2 (Stopped), Machine3 (Running)...	discovery/opcua.go → BrowseNodeTree() → Extracts WorkCenter names with "status" tags
Current status	Running ✅ / Stopped ❌	rules/engine.go → stateStore.Get() → Current state from StateStore
Status timestamp	Last change at 14:32:05	rules/engine.go → stateStore.Get() → Timestamp from StateStore
Transition history	Run → Stop (14:32:05), Stop → Run (14:33:05)...	rules/engine.go → stateStore.GetHistory() → History from StateStore
Duration	45 seconds	rules/engine.go → Calculated from transition timestamps
Configuration Panel:

Parameter	Description	Edge Source
Machine ID	Which machine to monitor	discovery/opcua.go → From BrowseNodeTree() (list of WorkCenters)
Initial State	Starting state (default: false)	functions/transforms/state_machine.go → Default parameter
5. 🗺️ UNS Mapper (Transform)
What the user sees in Frontend:

Site ID configuration

Area configuration

Tag normalization preview (machine1.temp → temperature)

Generated ISA-95 topic preview

Frontend UI Element	What It Displays	Where It Comes From (Edge)
Site ID	local-test (pre-filled)	config/agent.yaml → site.id
Area	area1 (pre-filled)	config/agent.yaml → Default "area1"
Normalized tag name	temperature (from temp), pressure (from presion), status (from stat)	uns/mapper.go → normalizeTagName() → Table of abbreviations
Inferred unit	celsius, bar, rpm, ""	uns/mapper.go → inferUnit() → From normalized name
Generated topic	mindset/site/local-test/area1/machine1/temperature	uns/mapper.go → UNSNode.FullTopic()
Tag description	Temperature sensor on machine1	uns/mapper.go → buildDescription()
Configuration Panel:

Parameter	Description	Edge Source
Site ID	Factory site identifier	config/agent.yaml → site.id
Area	Production area	config/agent.yaml → site.area (or default "area1")
Custom Tag Mapping	Add custom tag normalizations (optional)	User input → stored in uns/mapper.go at runtime
6. ⏱️ Duration (Calculate)
What the user sees in Frontend:

Duration calculator (start/stop automatic)

Result in seconds and minutes

Frontend UI Element	What It Displays	Where It Comes From (Edge)
Start time	2026-06-20T14:32:05Z	functions/calculates/duration.go → Recorded on first event
End time	2026-06-20T14:32:50Z	functions/calculates/duration.go → Recorded on second event
Duration (seconds)	45.0s	functions/calculates/duration.go → Calculated from start/end
Duration (minutes)	0.75min	functions/calculates/duration.go → Calculated from seconds
Configuration Panel:

Parameter	Description	Edge Source
None	No configuration parameters	Function is automatic
7. 💰 Cost (Calculate)
What the user sees in Frontend:

Cost per minute (auto-calculated from hourly rate)

Total cost in euros

Currency selector

Frontend UI Element	What It Displays	Where It Comes From (Edge)
Hourly rate	85.00 €/h (pre-filled)	config/agent.yaml → cost.hourly_cost
Cost per minute	1.42 €/min	functions/calculates/cost.go → hourly_rate / 60
Total cost	63.75 €	functions/calculates/cost.go → cost_per_minute × duration_minutes
Currency	EUR (pre-filled)	config/agent.yaml → cost.currency
Configuration Panel:

Parameter	Description	Edge Source
Hourly Rate	Cost per hour in €	config/agent.yaml → cost.hourly_cost
Currency	EUR, USD, GBP	config/agent.yaml → cost.currency
8. 🚦 Threshold (Condition)
What the user sees in Frontend:

Min value (default: 30 seconds)

Max value (default: 180 seconds)

Result (true/false) → Is it a micro-stop?

Frontend UI Element	What It Displays	Where It Comes From (Edge)
Min value	30 (pre-filled)	functions/conditions/threshold.go → Default min
Max value	180 (pre-filled)	functions/conditions/threshold.go → Default max
Current value	45.0s	functions/calculates/duration.go → Output from duration function
Result	✅ True (Is micro-stop)	functions/conditions/threshold.go → min < value < max
Configuration Panel:

Parameter	Description	Edge Source
Min	Minimum duration in seconds	functions/conditions/threshold.go → Default 30
Max	Maximum duration in seconds	functions/conditions/threshold.go → Default 180
9. 📤 MQTT Publish (Output)
What the user sees in Frontend:

Topic selector (ISA-95 or events)

QoS selection

Retained checkbox

Frontend UI Element	What It Displays	Where It Comes From (Edge)
Available topics	mindset/site/#, mindset/events/#, mindset/events/micro-stop...	mqtt/publisher.go → Known topics (from UNS mapper + rules engine)
ISA-95 topics	mindset/site/local-test/area1/machine1/temperature...	uns/mapper.go → Generated from UNS mapping
Event topics	mindset/events/micro-stop, mindset/events/status-change...	rules/engine.go → Published events
QoS	0, 1, 2	Documentation → Default 1
Retained	True / False	Documentation → Default false
Configuration Panel:

Parameter	Description	Edge Source
Topic	MQTT topic to publish to	List from mqtt/publisher.go (site/events)
QoS	Quality of Service	Documentation (default 1)
Retained	Keep last message	Documentation (default false)