Enhance the Pipeline Builder to provide a clear, guided, and error-proof user experience with:
- Machine & Tag Selection for OPC-UA Read
- Cost Configuration with live preview
- Duplicate Pipeline Prevention
- Knowledge Graph Grouping by tag set
- Clear Configuration Panels for each function
- Smart Error Messages with actionable guidance

1. OPC-UA Read — Machine & Tag Selection
When the user adds an opcua_read function or selects it in the trigger:
The configuration panel must display:

Element	Requirement
Machine dropdown	List of all discovered machines from OPC-UA discovery
Tag selection	Checkbox list of all tags for the selected machine
Select All / Deselect All	Buttons to quickly select/deselect all tags
Filter by type	Filter tags by data type (Boolean, Float, Int32, etc.)
Filter by name	Search/filter tags by name
Live preview	Show live values next to each tag (if agent is running)
All tags option	Option to select "All tags" from a machine

2. Cost Configuration — calculate_cost
When the user adds a calculate_cost function:

The configuration panel must display:

Element	Requirement
Hourly Rate Source	Radio buttons: Manual / From config / From tag
Manual entry	Number input field with unit (€/h)
From config	Pre-filled from agent.yaml (read-only with option to override)
From tag	Dropdown to select a tag that contains the hourly rate (e.g., from ERP)
Currency	Dropdown: EUR, USD, GBP
Live preview	Show cost for typical durations: 30s, 1min, 3min, 5min

3. Duplicate Pipeline Prevention
Rule: A pipeline cannot be saved if it has the same combination of tags and functions as an existing pipeline.


4. Knowledge Graph — Pipeline Grouping
Rule: When you build a pipeline with tags, and then add the same pipeline with different tags, the Knowledge Graph should show one pipeline node with all tags listed.


5. Clear Configuration Panel for Each Function
Every function's configuration panel must include:

Element	Requirement
Header	Function name + icon + category badge (Connector/Transform/Calculate/Condition/Output)
Description	Brief explanation of what the function does
Input fields	All required/optional parameters with labels
Help text	Tooltips or helper text below each field (ⓘ icon)
Validation	Real-time validation (required fields, format checks)
Preview	Live preview of the output (if applicable)
Examples	Example values shown as placeholders
Apply/Cancel	Buttons to apply or cancel changes

6. Smart Error Messages
Current Error: ❌ ID et nom sont requis.

Required Improvement: Error messages must include what is missing and how to fix it.

Error	New Message
Missing ID	❌ Veuillez donner un nom à votre pipeline. Le nom sera utilisé comme identifiant unique.
Missing Name	❌ Veuillez donner un titre à votre pipeline (ex: "Micro-stop Detection").
No trigger	❌ Aucun connecteur (trigger) trouvé. Veuillez ajouter un connecteur dans la zone ENTRÉE.
No output	❌ Aucune sortie trouvée. Veuillez ajouter "mqtt_publish" ou "add_to_dashboard" dans la zone SORTIE.
Missing machine_id	❌ Veuillez sélectionner une machine pour "state_machine".
Missing topic	❌ Veuillez sélectionner un topic pour "mqtt_subscribe".
Duplicate pipeline	⚠️ Une pipeline avec cette configuration existe déjà : "Micro-stop Detection". Options : [Modifier] [Nouvelle version] [Annuler]
Incompatible types	⚠️ La fonction "calculate_cost" attend une durée (en secondes), mais reçoit un booléen. Vérifiez la chaîne avant "calculate_cost".
Missing opcua tags	❌ Veuillez sélectionner au moins un tag pour "opcua_read".
7. Pipeline Naming & Save Flow
Save Flow Requirements:

Step	Action
1	User clicks "💾 Save Pipeline"
2	System validates: Name, ID, Trigger, Output, Tags, Connections
3	If validation fails → Show error message with fix instructions
4	If validation passes → Check for duplicates
5	If duplicate found → Show modal with options
6	If no duplicate → Save pipeline to YAML + Register in engine + Update KG
7	Show success notification with link to view in KG


8. Acceptance Criteria
User can select a machine and choose specific tags for OPC-UA Read

Live tag values are displayed in the configuration panel

Cost calculation shows a preview for 30s, 1min, 3min, 5min

Duplicate pipelines are blocked with a clear modal and options

Pipeline node in KG shows all associated tags as properties

Each function has a clear, structured configuration panel

Error messages include a description of the problem and a fix

Save flow validates name, ID, trigger, output before saving

Success notification includes link to view in Knowledge Graph