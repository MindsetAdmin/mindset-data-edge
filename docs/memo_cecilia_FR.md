# Mémo pour Cécilia — Réponses à tes questions sur l'archi

**De :** Mohamed
**Date :** 28 juin 2026
**Sujet :** Réponses à ton message sur l'archi on-premise vs cloud, et la différenciation pitch (notamment Cognite)
**Référence :** Toutes les sources sont dans `docs/MindSet_Competitive_Analysis_v2_2.xlsx` (9 onglets) et `docs/decisions.md` (25 décisions verrouillées)

---

## ⚠️ Avant tout — 2 corrections à apporter dans notre pitch

Ton message contient 2 formulations qu'on doit corriger **avant d'en reparler à un investisseur ou un partenaire** — sinon on se fera challenger.

| Ce que tu disais | Ce qu'il faut dire à la place | Pourquoi |
|---|---|---|
| « *On pourrait mettre en avant MCP natif et le fait qu'on ait AI et cloud agnostique* » | **« MCP natif au edge + EU-souverain + IA flexible (local par défaut, remote en option avec disclosure) »** | On a **explicitement rejeté** le « cloud-agnostique » : on refuse AWS / Azure / GCP jusqu'en 2029, par design. Et on n'est pas « AI-agnostic » : on tourne sur Phi-3 local par défaut, avec opt-in remote LLM. Le pitch « agnostique » casse notre moat souverain. |
| « *Cognite ne supporte pas MCP nativement* » | **« Cognite a ajouté MCP en 2026 (via Function Apps). Notre différenciation, c'est que leur MCP tourne dans leur cloud (les données partent du site), alors que notre MCP tourne au edge (les données restent dans l'usine). »** | Cognite a rattrapé sur MCP. Dire qu'ils ne l'ont pas nous fait passer pour pas à jour. Le bon angle, c'est **edge MCP vs cloud MCP**. |

---

## Réponses à tes 13 questions

### Question 1 — Architecture on-premise vs cloud

**Réponse :** On a 3 éditions au catalogue. Voir détail plus bas et **onglet 4 de l'Excel** (« 3 Editions »).

### Question 2 — Différenciation vs Cognite

**Réponse courte :** Cognite est cloud-mandatory pour le grand compte (oil & gas, utilities, contrats 6-7 chiffres). Nous, on est edge-first pour l'ETI manufacturière européenne (Plant Manager qui signe <30k€/site). Pas les mêmes clients, pas les mêmes deals — Cognite sert de **référence catégorielle pour les investisseurs**, pas de concurrent en compétition réelle.

Détail complet : **onglet 2 (Comp Matrix)** + **onglet 9 (Technical Diff)** de l'Excel.

### Question 3 — Cognite et MCP

**Réponse (avec correction de ton hypothèse) :** Cognite **a** ajouté MCP en 2026 via leur endpoint Function Apps. Mais ce MCP tourne dans LEUR cloud — les données partent du site. Notre MCP tourne sur le edge (le serveur du client), les données ne quittent jamais l'usine. C'est **« edge MCP » qui est notre différenciateur**, pas « MCP natif ».

### Question 4 — On met en avant MCP natif + AI/cloud-agnostique ?

**Réponse :**
- **MCP natif edge** → OUI, c'est verrouillé V1 (décision dans `decisions.md`)
- **AI agnostique** → NON. On est local-par-défaut (Phi-3 via Ollama), avec opt-in remote LLM **avec un warning explicite** quand le client active le mode remote (« vos données sortent du réseau / de l'UE »). C'est notre Option B (compromis souveraineté/flexibilité).
- **Cloud agnostique** → NON. EU-souverain par design, pas de hyperscaler avant 2029. C'est notre **vrai moat réglementaire** pour défense, secteur public, pharma régulée.

### Question 5 — Qu'est-ce qui est exactement on-premise vs cloud ?

**Réponse synthétique :**

**EDGE (on-premise) — TOUT le coeur du produit :**
- Connecteurs OPC-UA / Modbus / S7 / MQTT / SQL ERP
- Découverte automatique + classification des tags (Phi-3 local)
- Contextualisation UNS ISA-95
- Moteur de règles (micro-stop, énergie, OEE)
- **Fuzzy Join OT/IT basé sur l'état des OF** (pas sur un sliding window — voir correction plus bas)
- Calcul de coût en €
- Buffer SQLite 7-15 jours
- Dashboard local React
- Alertes locales (SMTP / Slack directs)
- **Serveur MCP** — au edge en V1
- **Agent IA Ad-hoc Analyst** — tourne sur Phi-3 local

**CLOUD — uniquement le strict nécessaire :**
- Agrégation cross-sites (V1.5+, données déjà transformées, jamais brutes)
- Dashboard multi-sites / accès distant (pour le CEO/CFO qui regarde de chez lui)
- API de gestion des sites (auth, clés API, licences)
- Backup chiffré du KG
- Monitor heartbeat (détecter qu'un edge est mort)

**Détail complet : onglet 8 de l'Excel (« Edge vs Cloud Map ») — 24 composants ligne par ligne avec la latence et le « pourquoi cet emplacement ».**

### Question 6 — Full on-premise vs hybride ?

**Réponse :** On a 3 éditions :

| Édition | Cloud ? | Pour qui |
|---|---|---|
| **On-Premise** | Zéro cloud, mono-site uniquement | Défense, secteur public, pharma sensible |
| **Hybrid (défaut)** | Cloud Scaleway FR / OVH FR (MindSet hosted) | ETI commerciale — l'offre standard |
| **Self-Hosted** | Cloud du client (Hetzner / IONOS / T-Systems / Bleu / on-prem K8s) | Multi-sites avec relation cloud EU existante |

**Pas d'édition Hyperscaler. AWS / Azure / GCP exclus par design — y compris leurs régions EU (US CLOUD Act).**

Détail : **onglet 4 de l'Excel** (« 3 Editions »).

### Question 7 — Cognite est quasi full cloud, confirme ?

**Réponse :** Confirmé. Cognite = thin extractor au edge + tout le reste dans leur cloud (KG, contextualisation, AI Atlas, dashboards, pipelines, MCP). Si leur cloud tombe, le client perd presque tout. Voir **onglet 9 (Technical Diff)** ligne « Failure mode if cloud unreachable ».

### Question 8 — Le reste au cloud (KG, contextualisation, agents IA, dashboards, pipelines) — confirme ?

**Réponse :** Confirmé pour Cognite. Pour nous, c'est l'INVERSE : tous ces composants tournent au edge par défaut. Voir comparaison ligne par ligne onglet 9.

### Question 9 — Tableau visualisable

**Réponse :** C'est exactement ce qu'est l'Excel `MindSet_Competitive_Analysis_v2_2.xlsx` (24 KB, 9 onglets). Onglet 4 répond spécifiquement à ta question scenario par scenario.

### Question 10 — Estimation des coûts (présentation Bleu en attente)

**Réponse partielle :** Le modèle 3 éditions est verrouillé, mais les **coûts précis par édition ne sont pas encore dans l'Excel** (gap identifié). Estimation rapide :

| Édition | Coût cloud / mois / site |
|---|---|
| On-Premise | **0 €** — zéro cloud |
| Hybrid (Scaleway FR) | **~15 €** au V0 (PLAY2-NANO + Managed Postgres + Object Storage) |
| Self-Hosted | Variable selon le cloud client. Indicatif : Hetzner CX21 ~5 €, Scaleway VPS 7-15 €, Bleu = à confirmer après ta réunion |

**Sur Bleu** : c'est positionné comme « cloud souverain FR » mais utilise la tech Microsoft Azure sous juridiction FR (JV Orange + Capgemini). C'est un **target Self-Hosted valable** pour les clients qui veulent du Azure-style avec garanties FR. Tarification à voir après ta présentation.

Si tu veux une **vraie page de coûts dans l'Excel** (onglet 10), je peux la faire en 15 min — dis-moi.

### Question 11 — Où est stocké le KG ?

**Réponse :** **Sur le serveur du client.** SQLite local (`data/mindset.db`). Le KG ne quitte jamais le réseau client en édition On-Premise. En Hybrid/Self-Hosted, **des snapshots agrégés** (jamais le KG brut) peuvent être répliqués vers le cloud pour le backup et l'agrégation multi-sites. Mais le master du KG vit toujours au edge.

Voir **onglet 8 de l'Excel** ligne « Knowledge Graph ».

### Question 12 — Où tourne le serveur MCP ?

**Réponse :** **Au edge, par défaut.** Embedded dans `cmd/server`, écoute sur `localhost:5000`. Les agents IA externes (Claude Desktop, Copilot, notre agent natif) s'y connectent depuis le réseau de l'usine. Un relai cloud MCP est prévu en V1.5+ pour les cas d'accès distant sans VPN — opt-in seulement.

Voir **onglet 9 (Technical Diff)** ligne « MCP server location ».

### Question 13 — Ce qui est sûr vs encore à déterminer

**Réponse :** C'est exactement la structure des onglets 5 et 6 de l'Excel.

- **Onglet 5 « Locked Decisions »** : 25 décisions verrouillées avec Rationale + Alternatives rejetées + Date. C'est ce dont on est SÛRS.
- **Onglet 6 « Open Questions »** : 9 questions encore ouvertes avec les options sur la table + qui décide + quand. Ex : Pricing, 4ème template, profil 2ème hire, hyperscaler 2029, etc.

Les décisions les plus structurantes verrouillées récemment :
- License **propriétaire fermée 2 ans** (réversion de Apache 2.0 — pour protection commerciale early-stage, ré-évalué 2028)
- **Pas de hyperscaler avant 2029** (réversion possible pour scaling international)
- 3 éditions : On-Premise / Hybrid (défaut) / Self-Hosted
- AI : Phi-3 local + opt-in remote avec disclosure (Option B)
- ERP V1 : connecteur SQL multi-dialectes (PostgreSQL + MSSQL + MySQL)
- Fuzzy Join : **basé sur l'état des OF** (pas sliding window — réversion technique majeure)
- Agent IA V1 : **Ad-hoc Analyst seul** (pas les 13 du catalogue, juste celui-là)
- Tribal Knowledge : **moat ships en V1** via dropdown + free text (chatbot Phi-3 = V2 polish, pas le moat)

---

## L'archi en 1 paragraphe (pour le pitch deck)

> *MindSet Data, c'est la plateforme industrielle edge native par IA pour les ETI manufacturières européennes. Une commande Docker installe un binaire Go unique sur un PC du client ; en 48h, la plateforme auto-découvre les équipements OT, contextualise la donnée en UNS ISA-95, attribue chaque événement OT à son Ordre de Fabrication actif via la lecture d'état ERP (robuste au décalage horaire de plusieurs heures, là où les jointures temporelles cassent), et permet à n'importe quel agent IA compatible MCP (Claude, Copilot, notre agent natif Ad-hoc Analyst) d'interroger l'usine directement — sans qu'aucune donnée brute ne quitte le réseau du client. V1 ship avec 3 templates prêts à l'emploi (micro-stop, gaspillage énergétique, OEE/TRS) ; les clients et leurs agents IA construisent les use cases suivants. 3 éditions (On-Premise / Hybrid / Self-Hosted), toutes EU-souveraines, JAMAIS sur hyperscaler US. Single-vendor, pas de frais par tag, pas de middleware Kepware-style. Le client possède son site fingerprint cumulatif.*

---

## La démo killer (à mettre en première slide)

**OEE réel vs déclaré** — le calcul que personne d'autre ne fait correctement en mid-market :

> *« Ton OEE déclaré, c'est 88%. On a mesuré chaque micro-arrêt sur ta Ligne 1 cette semaine — ton OEE RÉEL, c'est 74%. Le delta de 14 points = 1h04 de downtime caché par semaine = X€/semaine. Voici le Pareto des causes — top 3 fixes récupèrent Y€. »*

**Le gap LUI-MÊME est la value proposition** — c'est l'argent que le Plant Manager perd sans le savoir.

Mécanique détaillée : onglet 1 de l'Excel, section « HOW WE DETECT REAL OEE vs DECLARED OEE » (mise en gold).

---

## Prochaines actions de mon côté (si tu valides)

1. **Construire l'onglet 10 « Cost Estimation »** dans l'Excel — pour ta réunion Bleu (~15 min)
2. **Mettre à jour `docs/mindset.md`** (le doc vision de 1257 lignes) pour qu'il reflète les 25 décisions verrouillées — housekeeping, ne bloque pas le pitch (~45 min)
3. **Décision pricing** — toujours ouverte, doit être tranchée avant le deck investisseur (workshop d'1h)
4. **Ajouter Braincube FR** au comp matrix — c'est notre rival direct le plus proche en FR

Dis-moi ce que tu veux que je traite en priorité.

---

*Toutes les sources et détails sont dans le repo `mindset-data-edge` :*
- *`docs/MindSet_Competitive_Analysis_v2_2.xlsx`* — Analyse complète (9 onglets)
- *`docs/decisions.md`* — 25 décisions verrouillées
- *`docs/analysis_log.md`* — Log complet du raisonnement (Entries 1-15)
- *`docs/mindset.md`* — Vision originale (à mettre à jour)
