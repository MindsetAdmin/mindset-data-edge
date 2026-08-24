# Fiabilité, gouvernance & sécurité — points à mettre en avant (2026-08-02)

Liste complète, mise à jour, pour rassurer clients et investisseurs sur la gestion de la donnée — accès humains, systèmes, agents IA. Structurée en deux niveaux volontairement : **ce qui est vrai aujourd'hui** (vérifié contre le code réel) vs. **ce qui est roadmap**. Vu le pattern déjà identifié dans nos propres docs de pitch (citation McKinsey non sourcée, architecture décrite dans `mindset.md` jamais construite — voir `docs/analysis_log.md` Entrées 137-141), je préfère qu'on ne répète pas l'erreur ici, surtout sur un sujet où un acheteur pharma/cosmétique posera des questions précises.

---

## 1. Contrôle d'accès & gestion d'identité

**Aujourd'hui — à ne PAS présenter comme acquis :** c'est le point le plus important à ne pas survendre. Il n'existe pas encore de couche d'authentification/autorisation sur l'API ni sur l'UI (CORS ouvert `*`, pas de login, pas de gestion de rôles). C'est un vrai vide, pas juste "pas fini" — si un client sécurité-sensible pose la question directement, la réponse honnête est "sur la roadmap V1.5 / motion entreprise", pas "c'est géré".

**Ce qu'on a déjà, en hygiène des identifiants (différent du contrôle d'accès utilisateur) :**
- Les identifiants de connexion SQL (ERP/MES) ne sont jamais en clair dans la config — indirection par variable d'environnement (`password_env`), résolue seulement à l'ouverture de connexion.
- Le connecteur SQL est strictement en lecture seule, requêtes paramétrées — protection structurelle contre l'injection, indépendante de toute couche d'auth.

**Vrai gap de sécurité côté OT, à connaître avant qu'un client technique le trouve lui-même :** les modes sécurisés OPC-UA (`Sign` / `SignAndEncrypt`, avec certificat client) ne sont pas encore câblés — seul le mode `None` fonctionne aujourd'hui. Sur un vrai déploiement client, c'est un point d'attention réel, pas cosmétique.

**Roadmap à nommer explicitement :** RBAC, SSO, audit log signé — décrits dans nos docs de planification, pas implémentés.

---

## 2. Observabilité

**Vrai aujourd'hui**, et sous-estimé comme argument :
- `GET /api/health` — liveness basique.
- `GET /api/stats` — compteurs, uptime, statut de connexion au broker MQTT.
- `GET /api/topics` — topics actifs en temps réel, débit (msg/s) par topic, catégorie, statut `broker_connected`.
- Suivi d'état par machine (Running/Stopped + historique des transitions) — visibilité opérationnelle continue, pas juste un snapshot.
- Logging structuré (`log.Printf`, sortie stderr) — vérifié en conditions réelles que le flux stdout du transport MCP stdio ne contient jamais de log parasite, seulement le protocole. Discipline réelle, pas une affirmation en l'air.

**Roadmap, à nommer si la question vient :** pas de stack de métriques/traces centralisée (type Prometheus/Grafana/OpenTelemetry), pas d'alerting automatique sur échec de pipeline, pas de stockage de logs interrogeable à froid.

---

## 3. Intégrité & gouvernance de la donnée

**Vrai aujourd'hui — le cœur de l'argumentaire fiabilité :**
- Aucune donnée générée automatiquement n'est jamais silencieusement considérée comme fiable : chaque mapping auto-généré porte un score de confiance explicite (heuristique, pas boîte noire). En dessous du seuil, la donnée reste `pending` — visible comme telle dans l'interface (anneau pointillé), validée ou rejetée explicitement par un humain avant qu'un agent IA ou un dashboard ne la traite comme un fait.
- Discipline anti-fabrication : le système refuse d'inventer un chiffre non auditable (ex. une pénalité de retard non contractualisée dans la donnée source) — il signale l'urgence plutôt que d'inventer un montant. C'est une garde-fou produit, pas juste une promesse marketing.
- Les calculs critiques qui combinent plusieurs signaux (coût + urgence de livraison, par exemple) sont calculés **une seule fois, côté serveur, de façon déterministe** — jamais laissés à un agent IA pour les recalculer à chaque requête. Résultat : la même question donne toujours la même réponse, quel que soit l'agent qui la pose. C'est une propriété de fiabilité que peu de plateformes "AI-native" ont réellement.

---

## 4. Diffusion sécurisée vers les agents IA

**Vrai aujourd'hui :**
- Le serveur MCP exposé aux agents (Claude, Copilot, etc.) est **strictement en lecture seule** de bout en bout — un agent peut interroger le graphe de connaissance, jamais le modifier ni déclencher une action. Un agent qui "part en vrille" ne peut pas corrompre la donnée.
- Chaque outil MCP a une description qui précise explicitement ce qu'il ne répond PAS (ex. le système dit clairement qu'il ne peut pas répondre à des questions historiques sur un produit donné, plutôt que de laisser l'agent halluciner une réponse) — conçu pour empêcher la sur-affirmation de capacité côté IA, pas juste pour documenter.
- Les réponses sont "groundées" par construction : ce sont des fonctions Go déterministes qui interrogent des lignes réelles de la base, pas des résumés générés par un LLM — donc pas de risque d'hallucination sur la donnée elle-même.

---

## 5. Résilience opérationnelle

**Vrai aujourd'hui, et idéalement raconté comme un retour d'expérience réel, pas une promesse abstraite :**
- Un vrai incident en prod (le dashboard se figeait silencieusement après un certain temps) a été diagnostiqué et corrigé structurellement : l'écriture en base a été découplée du chemin de données temps réel, pour qu'une écriture lente ne puisse plus jamais geler tout le pipeline live.
- Les composants (serveur API, agent terrain) sont couplés de façon lâche — ils ne s'appellent jamais directement, seulement via un bus de messages (MQTT) et une base partagée. Un composant qui tombe n'entraîne pas l'autre avec lui.
- Chaque abonné MQTT a un identifiant client distinct — corrige une classe de bug où deux composants pouvaient silencieusement s'évincer l'un l'autre sans erreur visible.

---

## Résumé — ce qu'on peut dire sans risque vs. ce qu'il faut nommer comme roadmap

| Dire "c'est fait" | Dire "c'est sur la roadmap" |
|---|---|
| Lecture seule sur toutes les sources (PLC/SCADA/ERP/MES) | RBAC / SSO / gestion des rôles |
| Agents IA en lecture seule sur le graphe | Chiffrement au repos |
| Validation humaine obligatoire sous le seuil de confiance | Audit log signé |
| Pas de chiffres inventés — le système signale plutôt que d'halluciner | Stack de métriques/alerting centralisée |
| Calculs critiques déterministes, cohérents entre agents | Modes OPC-UA sécurisés (Sign/SignAndEncrypt) — actuellement `None` seulement |
| Observabilité basique réelle (health, stats, débit par topic, état machine) | Conformité formelle (ISO 27001, GAMP 5) |
| Résilience testée sur un vrai incident de prod, pas juste en théorie | — |

---

**Non fait** : ce document n'a pas encore été confronté à Djamil/Jalil ni à un vrai brief investisseur — c'est la matière brute, à ajuster selon le format final (deck, FAQ sécurité, one-pager).
