the title is MindSet Data in place of MindSet Data — Pipeline Studio

Si on clique sur pipline pour chrger et on clique sur un deuxieme, le premier reste et ne disparaitre pas

le pipline est juste le coeur du pipline (y a pas l'entré et la sortie)

Au lieu d'avoir une fonction kg_save que l'utilisateur doit ajouter,
le Knowledge Graph s'enrichit AUTOMATIQUEMENT à chaque événement.
L'utilisateur n'a rien à faire → le KG est toujours à jour.

L'utilisateur ne crée PAS de nouveaux connecteurs, topics ou fonctions.
Il CHOISIT parmi ceux qui sont déjà disponibles.
Pour chaque sélection :
1. Une liste des options disponibles est affichée
2. L'utilisateur peut rechercher dans la liste
3. L'utilisateur sélectionne une option
4. La configuration est pré-remplie automatiquement


 Les différents "Pickers" dans l'UI
1. Sélecteur de Connecteur (MQTT Subscribe)
text
┌─────────────────────────────────────────────────────────────────────────────────────┐
│  📡 MQTT Subscribe - Choisir le broker et le topic                                 │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  🔌 Broker MQTT :                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │  🔍 Rechercher un broker...                                                │   │
│  │                                                                             │   │
│  │  ○ tcp://localhost:1883              [Connecté]  (défaut)                  │   │
│  │  ○ tcp://192.168.1.100:1883          [Connecté]                            │   │
│  │  ○ tcp://factory-mqtt.internal:1883  [Déconnecté]                          │   │
│  │  ○ tcp://broker.aws.amazon.com:1883  [Déconnecté]                          │   │
│  │                                                                             │   │
│  │  [Ajouter un nouveau broker] → Ouvre un formulaire                         │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                     │
│  📡 Topics disponibles sur ce broker :                                             │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │  🔍 Rechercher un topic...                                                 │   │
│  │                                                                             │   │
│  │  mindest/raw/ns=3;i=1009       ☑ Sélectionné   (6 msg/s)                   │   │
│  │  mindest/raw/ns=3;i=1011       ○              (4 msg/s)                    │   │
│  │  mindest/raw/ns=3;i=1014       ○              (2 msg/s)                    │   │
│  │  mindest/raw/#                 ○              (Tous les tags)              │   │
│  │  mindest/site/local-test/...   ○              (1 msg/s)                    │   │
│  │  mindest/events/micro-stop     ○              (0.5 msg/s)                  │   │
│  │                                                                             │   │
│  │  [S'abonner à un topic personnalisé] → Champ texte avec validation        │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                     │
│  [Valider]  [Annuler]                                                               │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
2. Sélecteur de Fonction
text
┌─────────────────────────────────────────────────────────────────────────────────────┐
│  ⚙️ Ajouter une fonction                                                            │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  🔍 Rechercher une fonction...                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │  🔍 filter                                                          [🔍]   │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                     │
│  📋 Fonctions disponibles :                                                         │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │  🔌 Connecteurs                                                             │   │
│  │  ┌──────────────────────────────────────────────────────────────────────┐   │   │
│  │  │  📡 mqtt_subscribe    S'abonne à MQTT                [Ajouter]     │   │   │
│  │  └──────────────────────────────────────────────────────────────────────┘   │   │
│  │  ⚙️ Transforms                                                              │   │
│  │  ┌──────────────────────────────────────────────────────────────────────┐   │   │
│  │  │  🔍 filter            Filtre les données               [Ajouter]     │   │   │
│  │  │  🔄 state_machine     Détecte transitions              [Ajouter]     │   │   │
│  │  │  🗺️ uns_mapper        Normalise ISA-95                 [Ajouter]     │   │   │
│  │  └──────────────────────────────────────────────────────────────────────┘   │   │
│  │  📊 Calculs                                                                │   │
│  │  ┌──────────────────────────────────────────────────────────────────────┐   │   │
│  │  │  ⏱️ duration          Calcule une durée                  [Ajouter]  │   │   │
│  │  │  💰 cost              Calcule un coût                    [Ajouter]  │   │   │
│  │  └──────────────────────────────────────────────────────────────────────┘   │   │
│  │  🚦 Conditions                                                              │   │
│  │  ┌──────────────────────────────────────────────────────────────────────┐   │   │
│  │  │  🚦 threshold         Vérifie un seuil                    [Ajouter]  │   │   │
│  │  └──────────────────────────────────────────────────────────────────────┘   │   │
│  │  📤 Outputs                                                                 │   │
│  │  ┌──────────────────────────────────────────────────────────────────────┐   │   │
│  │  │  📤 mqtt_publish      Publie sur MQTT                  [Ajouter]    │   │   │
│  │  │  💾 kg_save           Sauvegarde dans le KG             [Ajouter]    │   │   │
│  │  └──────────────────────────────────────────────────────────────────────┘   │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                     │
│  [Fermer]                                                                           │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
3. Sélecteur de Pipeline Pré-définie
text
┌─────────────────────────────────────────────────────────────────────────────────────┐
│  📋 Charger une pipeline pré-définie                                                │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  🔍 Rechercher une pipeline...                                                    │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │  🔍 micro                                                          [🔍]   │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                     │
│  📋 Pipelines disponibles :                                                         │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │  🔴 Micro-stop Detection     │  Détecte les micro-arrêts                    │   │
│  │  ────────────────────────────│───────────────────────────────────────────────│   │
│  │  Nœuds : 6                   │  Dépend de : opcua_read, state_machine,     │   │
│  │  Version : 1.0               │  duration, threshold, mqtt_publish, kg_save  │   │
│  │  [Charger]                   │                                             │   │
│  ├─────────────────────────────┼───────────────────────────────────────────────┤   │
│  │  🔄 OPC-UA → UNS             │  Normalise en ISA-95                         │   │
│  │  ────────────────────────────│───────────────────────────────────────────────│   │
│  │  Nœuds : 3                   │  Dépend de : mqtt_subscribe, uns_mapper,    │   │
│  │  Version : 1.0               │  mqtt_publish                                │   │
│  │  [Charger]                   │                                             │   │
│  ├─────────────────────────────┼───────────────────────────────────────────────┤   │
│  │  💰 Cost Calculation         │  Calcule le coût des micro-stops             │   │
│  │  ────────────────────────────│───────────────────────────────────────────────│   │
│  │  Nœuds : 4                   │  Dépend de : mqtt_subscribe, cost,          │   │
│  │  Version : 1.0               │  mqtt_publish, kg_save                      │   │
│  │  [Charger]                   │                                             │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                     │
│  [Fermer]                                                                           │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
4. Sélecteur de Topic (pour MQTT Publish)
text
┌─────────────────────────────────────────────────────────────────────────────────────┐
│  📤 MQTT Publish - Choisir le topic de destination                                  │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  📡 Topics ISA-95 disponibles :                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │  🔍 Rechercher un topic...                                                 │   │
│  │                                                                             │   │
│  │  🌐 Site level                                                              │   │
│  │  └── mindset/site/local-test/temperature   ○  (0 msg)                      │   │
│  │  └── mindset/site/local-test/status        ○  (0 msg)                      │   │
│  │                                                                             │   │
│  │  🏭 Area level                                                              │   │
│  │  └── mindset/site/local-test/area1/temperature  ○  (2 msg/s)               │   │
│  │  └── mindset/site/local-test/area1/status       ○  (1 msg/s)               │   │
│  │                                                                             │   │
│  │  ⚙️ Work Center level                                                        │   │
│  │  └── mindset/site/local-test/area1/machine1/temperature  ☑ Sélectionné     │   │
│  │  └── mindset/site/local-test/area1/machine1/status        ○                │   │
│  │  └── mindset/site/local-test/area1/machine2/temperature   ○                │   │
│  │                                                                             │   │
│  │  📡 Events level                                                            │   │
│  │  └── mindset/events/micro-stop                     ○                       │   │
│  │  └── mindset/events/status-change                  ○                       │   │
│  │                                                                             │   │
│  │  [Topic personnalisé] → Validation format ISA-95                           │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                     │
│  [Valider]  [Annuler]                                                               │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
5. Sélecteur de Machine (pour State Machine)
text
┌─────────────────────────────────────────────────────────────────────────────────────┐
│  🔄 State Machine - Choisir la machine                                              │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  🔍 Rechercher une machine...                                                     │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │  🔍 machine1                                                       [🔍]   │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                     │
│  🏭 Machines disponibles :                                                          │
│  ┌─────────────────────────────────────────────────────────────────────────────┐   │
│  │  ☑ machine1   (2 tags)       Status: Running   Temp: 23.5°C               │   │
│  │  ○ machine2   (3 tags)       Status: Stopped   Temp: 21.0°C               │   │
│  │  ○ machine3   (1 tag)        Status: Running   Temp: 22.0°C               │   │
│  │  ○ machine4   (0 tag)        Status: Unknown                               │   │
│  │                                                                             │   │
│  │  [Ajouter une machine manuellement] → Champ texte avec validation          │   │
│  └─────────────────────────────────────────────────────────────────────────────┘   │
│                                                                                     │
│  [Valider]  [Annuler]                                                               │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘

Le dashboard affiche des valeurs reel de chaque data or event qui est lié au dashboard (avec un graphe d'historique)

Le KG affiche reelement ce qu'on a 

dans compose ajoute une option pour supprimer un noeud


Fonctions Disponibles — MindSet Data
🔌 1. mqtt_subscribe (Connecteur)
Description : S'abonne à un topic MQTT et reçoit les messages en temps réel.

Type : Source (0 entrée → 1 sortie)

Configuration :

Paramètre	Type	Description	Exemple
broker	string	URL du broker MQTT	tcp://localhost:1883
topic	string	Topic auquel s'abonner	mindset/raw/ns=3;i=1014
qos	number	Qualité de service (0, 1, 2)	1
Sortie : Message MQTT reçu

json
{
  "node_id": "ns=3;i=1014",
  "name": "machine1.status",
  "value": true,
  "data_type": "Boolean",
  "timestamp_ms": 1705419450718
}
⚙️ 2. filter (Transform)
Description : Filtre les données selon une condition (ex: garder les valeurs > 0).

Type : Transform (1 entrée → 1 sortie)

Configuration :

Paramètre	Type	Description	Exemple
field	string	Champ à filtrer	"value"
operator	string	Opérateur (eq, ne, gt, lt, contains)	"gt"
value	any	Valeur de comparaison	0
Sortie : Données filtrées (ou rien si ne passe pas le filtre)

Exemple : Garde uniquement les valeurs > 0

⚙️ 3. state_machine (Transform)
Description : Détecte les transitions d'état (Run ↔ Stop) et déclenche sur les arrêts.

Type : Transform (1 entrée → 1 sortie)

Configuration :

Paramètre	Type	Description	Exemple
machine_id	string	Identifiant de la machine	"machine1"
Sortie : Transition détectée

json
{
  "from": true,
  "to": false,
  "timestamp": "2026-06-15T14:32:05Z",
  "duration_seconds": 45.0,
  "work_center": "machine1"
}
Comportement :

true → false : Machine arrêtée (début du stop)

false → true : Machine redémarrée (durée calculée)

⚙️ 4. uns_mapper (Transform)
Description : Normalise les tags OPC-UA en topics ISA-95 structurés.

Type : Transform (1 entrée → 1 sortie)

Configuration :

Paramètre	Type	Description	Exemple
site_id	string	Identifiant du site	"local-test"
area	string	Zone	"area1"
Sortie : Topic UNS structuré

json
{
  "topic": "local-test/area1/machine1/temperature",
  "full_topic": "mindset/site/local-test/area1/machine1/temperature",
  "node": {
    "site": "local-test",
    "area": "area1",
    "work_center": "machine1",
    "tag_name": "temperature",
    "unit": "celsius"
  }
}
Normalisation :

temp → temperature

presion → pressure

stat → status

spd → speed

📊 5. duration (Calculate)
Description : Calcule la durée entre deux événements (start/stop).

Type : Calculate (1 entrée → 1 sortie)

Configuration : Aucune

Sortie : Durée en secondes

json
{
  "start_time": "2026-06-15T14:32:05Z",
  "end_time": "2026-06-15T14:32:50Z",
  "seconds": 45.0,
  "minutes": 0.75,
  "work_center": "machine1"
}
📊 6. cost (Calculate)
Description : Calcule le coût en euros d'un événement.

Type : Calculate (1 entrée → 1 sortie)

Configuration :

Paramètre	Type	Description	Exemple
hourly_rate	number	Coût horaire en €/h	85.0
currency	string	Devise	"EUR"
Sortie : Coût en euros

json
{
  "duration_seconds": 45.0,
  "duration_minutes": 0.75,
  "cost_per_minute": 1.416,
  "total_cost_eur": 63.75,
  "currency": "EUR",
  "work_center": "machine1"
}
Formule : cost = (duration_seconds / 3600) × hourly_rate

🚦 7. threshold (Condition)
Description : Vérifie si une durée est entre un minimum et un maximum (micro-stop).

Type : Condition (1 entrée → 1 sortie)

Configuration :

Paramètre	Type	Description	Exemple
min	number	Durée minimale en secondes	30
max	number	Durée maximale en secondes	180
Sortie : Événement si la valeur est dans l'intervalle

json
{
  "value": 45.0,
  "is_in_range": true,
  "is_micro_stop": true,
  "work_center": "machine1"
}
Règle :

Si 30s < valeur < 180s → Micro-stop 🟢

Si valeur < 30s → Ignoré (bruit)

Si valeur > 180s → Arrêt majeur ⛔

📤 8. mqtt_publish (Output)
Description : Publie des données sur un topic MQTT.

Type : Output (1 entrée → 0 sortie)

Configuration :

Paramètre	Type	Description	Exemple
topic	string	Topic de publication	"mindset/events/micro-stop"
qos	number	Qualité de service	1
retained	boolean	Message retenu	false
Sortie : Aucune (publie sur MQTT)

💡 Knowledge Graph (Automatique)
Le Knowledge Graph s'enrichit automatiquement à chaque événement publié sur mindset/events/#. L'utilisateur n'a rien à configurer.

Ce qui est ajouté automatiquement :

Nœud Equipment (machine)

Nœud Event (micro-stop)

Relation occurred_at

Nœud Cause (si présente)

Nœud Cost (si présent)

🏗️ Pipelines Disponibles
🔴 Pipeline 1 : Micro-stop Detection
Description : Détecte les micro-arrêts (30s-3min) et publie les événements.

Nœuds :

text
mqtt_subscribe → state_machine → duration → threshold → mqtt_publish
Composition :

Étape	Fonction	Rôle
1	mqtt_subscribe	S'abonne au topic de la machine (ex: mindset/raw/ns=3;i=1014)
2	state_machine	Détecte les transitions Run ↔ Stop
3	duration	Calcule la durée de l'arrêt
4	threshold	Vérifie si la durée est entre 30s et 180s
5	mqtt_publish	Publie sur mindset/events/micro-stop
YAML :

yaml
id: "pipeline_microstop_detection"
name: "Micro-stop Detection"
trigger:
  type: "mqtt"
  function: "mqtt_subscribe"
  config:
    broker: "tcp://localhost:1883"
    topic: "mindset/raw/ns=3;i=1014"
    qos: 1
nodes:
  - id: "state_machine"
    function: "state_machine"
    config:
      machine_id: "machine1"
  - id: "duration"
    function: "duration"
  - id: "threshold"
    function: "threshold"
    config:
      min: 30
      max: 180
  - id: "mqtt_publish"
    function: "mqtt_publish"
    config:
      topic: "mindset/events/micro-stop"
      qos: 1
🔄 Pipeline 2 : OPC-UA → UNS (ISA-95)
Description : Transforme les données OPC-UA brutes en UNS ISA-95 structuré.

Nœuds :

text
mqtt_subscribe → uns_mapper → mqtt_publish
Composition :

Étape	Fonction	Rôle
1	mqtt_subscribe	S'abonne aux données brutes (mindset/raw/#)
2	uns_mapper	Normalise le tag en ISA-95 (ex: machine1.temp → site/area1/machine1/temperature)
3	mqtt_publish	Publie sur mindset/site/#
YAML :

yaml
id: "pipeline_opcua_to_uns"
name: "OPC-UA to UNS"
trigger:
  type: "mqtt"
  function: "mqtt_subscribe"
  config:
    broker: "tcp://localhost:1883"
    topic: "mindset/raw/#"
    qos: 1
nodes:
  - id: "uns_mapper"
    function: "uns_mapper"
    config:
      site_id: "local-test"
      area: "area1"
  - id: "mqtt_publish"
    function: "mqtt_publish"
    config:
      topic: "mindset/site/local-test/area1"
      qos: 1
💰 Pipeline 3 : Cost Calculation
Description : Calcule le coût des micro-stops en euros.

Nœuds :

text
mqtt_subscribe → cost → mqtt_publish
Composition :

Étape	Fonction	Rôle
1	mqtt_subscribe	S'abonne aux micro-stops (mindset/events/micro-stop)
2	cost	Calcule le coût en euros (85€/h)
3	mqtt_publish	Publie sur mindset/events/micro-stop-cost
YAML :

yaml
id: "pipeline_cost_calculation"
name: "Cost Calculation"
trigger:
  type: "mqtt"
  function: "mqtt_subscribe"
  config:
    broker: "tcp://localhost:1883"
    topic: "mindset/events/micro-stop"
    qos: 1
nodes:
  - id: "cost"
    function: "cost"
    config:
      hourly_rate: 85.0
      currency: "EUR"
  - id: "mqtt_publish"
    function: "mqtt_publish"
    config:
      topic: "mindset/events/micro-stop-cost"
      qos: 1
