
1. Cost Configuration — calculate_cost
When the user adds a calculate_cost function:

The configuration panel must display:

Element	Requirement
Hourly Rate Source	Radio buttons: Manual / From config / From tag
Manual entry	Number input field with unit (€/h)
From config	Pre-filled from agent.yaml (read-only with option to override)
From tag	Dropdown to select a tag that contains the hourly rate (e.g., from ERP)
Currency	Dropdown: EUR, USD, GBP
Live preview	Show cost for typical durations: 30s, 1min, 3min, 5min

2. Duplicate Pipeline Prevention
Rule: A pipeline cannot be saved if it has the same combination of tags and functions as an existing pipeline.


3. Knowledge Graph — Pipeline Grouping
Rule: When you build a pipeline with tags, and then add the same pipeline with different tags, the Knowledge Graph should show one pipeline node with all tags listed.

4. Output Functions — Single Connection Point
Current Behavior
Output functions (mqtt_publish, add_to_dashboard) currently have 2 connection points (input + output).

Required Behavior
Output functions must have only 1 connection point (input only).

Function	Current	Required
mqtt_publish	Input + Output	Input only (no output port)
add_to_dashboard	Input + Output	Input only (no output port)

5. Cost Calculation — Excel/CSV File Upload

Required Behavior
Users can upload an Excel or CSV file containing cost data with multiple rates:

Product	Cost per Hour (€/h)	Cost per Unit (€)
Product_A	85.00	2.50
Product_B	95.00	3.20
Product_C	75.00	1.80
Default	85.00	2.00

6. Dashboard: Remove "Pareto des causes" Section

Remove the Pareto section entirely from the Dashboard.

7. Dashboard: Fix add_to_dashboard
Data/events published by add_to_dashboard must appear immediately on the Dashboard.
