# Questions & Answers

_Format : chaque question de Mohamed / Cécilia est suivie d'une réponse. Chaque paire (Q → R) donne une nouvelle entrée dans `docs/analysis_log.md`._

---

## Q1 — 2026-07-05 — KG + MCP + agent externe (pour le call de Cécilia du 06/07/2026)

demain ma cofounder a un call avec un startup qui build des agents IA qui aide l'automatisation des tâches dans les usines.
voici ce qu'elle m'a demandé :

> [7:06 pm, 05/07/2026] Cecilia: Hey !
> [7:06 pm, 05/07/2026] Cecilia: Petite question sur le knowledge graph aujourd'hui
> [7:07 pm, 05/07/2026] Cecilia: Let's say a few pipelines are built for a certain purpose that we agreed with a client
> [7:08 pm, 05/07/2026] Cecilia: Si dessus on rajoute un MCP, l'agent a les liens qu'il faut pour raisonner dessus right ?
> [7:09 pm, 05/07/2026] Cecilia: Ofc si tu les bonnes pipelines ont été built
> [7:09 pm, 05/07/2026] Cecilia: Je demande car c'est important pour le call de demain
> [7:09 pm, 05/07/2026] Cecilia: Ils ont travaillé déjà avec un acteur qui avait un knowledge graph, pour sur eux ils puissent directement se brancher dessus
> [7:11 pm, 05/07/2026] Cecilia: Je vais poser la question pour savoir quelle type de données et de liaison / comment le knowledge graph etait build chez eux, ça peut être utile
> [7:11 pm, 05/07/2026] Cecilia: Mais surtout pour déterminer si la manière dont le notre est build et se build, est compatible avec eux
> [7:11 pm, 05/07/2026] Cecilia: Si dessus on rajoute un MCP, l'agent a les liens qu'il faut pour raisonner dessus right ?

Répond à ça

---

### Réponse — pour le call du 06/07/2026

**Contexte :** ils ont déjà travaillé avec un KG externe et veulent se brancher directement sur le nôtre.

#### Réponse courte

**Oui — à condition que 3 choses soient vraies :**

1. **Les pipelines client produisent des événements enrichis** (batch, cause, coût, équipement, opérateur…) — pas juste des points de télémétrie bruts. Sans enrichissement, le KG est un stream de valeurs, pas un graphe de contexte.
2. **Le MCP expose au minimum** : (a) le schéma du KG (types de nœuds + relations), (b) une traversée du graphe, (c) l'accès aux valeurs live des tags.
3. **L'agent parle une ontologie compatible avec la nôtre**, OU on lui fournit une couche de mapping.

Sans les 3, l'agent peut lire nos données mais ne saura pas les interpréter dans le contexte usine.

#### Le mécanisme en 4 lignes

```
Pipelines (client)
   → événements enrichis (Event, Cause, Cost, Equipment liés)
      → KG (nœuds + relations persistés)
         → MCP (expose schéma + traversée + live)
            → Agent raisonne
```

Le **KG répond à "quoi et pourquoi"**. Le **MCP donne l'accès à l'agent**. L'agent raisonne — mais uniquement dans les limites de ce que le KG et le MCP savent exposer.

#### Ce qui déterminera la compatibilité (les 3 axes)

**1. Ontologie**

Notre KG business a 4 types principaux : **Equipment · Event · Cause · Cost**.
Notre KG platform (topologie technique) a : **Pipeline · Function · Topic · Connection · Dashboard**.
Relations principales : `emits`, `caused_by`, `costs`, `runs_on`, `feeds_into`.

- Si leur agent attend une **ontologie ISA-95** (Site > Area > WorkCenter > WorkUnit) → compatible, on le mappe naturellement (c'est déjà notre structure sous-jacente).
- Si leur agent attend une **ontologie custom** → il faut voir laquelle et construire un connecteur.
- Si leur agent est **ontology-agnostic** (LLM qui lit le schéma dynamiquement) → parfait, aucun mapping requis.

**2. Accès live vs. historique**

Notre KG contient les **événements passés** + la **topologie**.
Les **valeurs live** des tags OPC-UA vivent dans le TagRegistry (accessible via `/api/tags` + WebSocket `/api/ws`).

Question critique : si leur agent doit répondre à *"quelle est la vitesse de la Ligne 3 en ce moment ?"*, il lui faut l'accès aux DEUX :
- KG (contexte : quelle ligne, quel batch, quel opérateur, quel coût par minute)
- TagRegistry (valeur live)

Notre MCP devra exposer les deux. Ce n'est pas juste un serveur "graphe" — c'est un serveur "graphe + valeurs".

**3. Protocole d'accès**

Aujourd'hui notre KG est exposé via **REST** (`/api/kg?category=business|platform|all`).

Comment leur agent consomme habituellement un KG ?
- **MCP natif** (le standard actuel — c'est exactement ce qu'on prévoit) → parfait, on livre un serveur MCP avec les bons tools.
- **GraphQL** → possible, on peut wrapper.
- **SPARQL** (KG sémantique classique) → possible mais gros travail, notre KG n'est pas RDF.
- **REST custom** → OK, mais moins scalable.

#### Questions à leur poser demain — copie-colle

1. **Ontologie** — *"Vous parlez quel vocabulaire côté agent ? ISA-95 ? Une ontologie custom ? Un standard IEC ?"*

2. **Interface** — *"Comment votre agent consomme le KG aujourd'hui ? MCP natif, GraphQL, SPARQL, REST custom ?"*

3. **Live vs. batch** — *"L'agent a besoin de valeurs live (état actuel de la ligne), ou uniquement d'historique (analyse post-mortem) ?"*

4. **Types de raisonnement** — *"Vos agents font quoi typiquement : traversée de graphe, semantic search, causal reasoning, action triggering ? Ça détermine ce qu'on doit exposer en priorité côté MCP."*

5. **Prior KG** — *"Le KG précédent avec lequel vous vous êtes branchés, il exposait quoi comme schéma / endpoints / protocoles ? Vous avez pu extraire tout ce que vous vouliez, ou il vous manquait des choses ?"*
   → **C'est la question la plus utile.** Elle te dit exactement ce qu'ils veulent qu'on livre sans avoir à leur expliquer notre archi.

6. **Latence** — *"Quelle latence l'agent tolère pour une réponse ? < 1 s ? < 10 s ? Ça influence si on doit tout mettre en cache côté MCP ou requêter le KG en live."*

7. **Cible client** — *"Les usines où vos agents tournent aujourd'hui, elles sont plutôt cloud-first ou on-prem ? Parce qu'on tourne sur la box du client — data ne sort jamais."*
   → Filtre les prospects incompatibles dès le call.

#### Angle stratégique — comment se positionner

Pour eux, MindSet = la couche de contexte industriel qu'ils n'ont pas à construire eux-mêmes. Ils se concentrent sur les agents (leur moat), on se concentre sur le KG + accès live + protocoles OT (notre moat).

Pour nous, leur agent = une feature "automatisation" qu'on n'a pas à construire nous-mêmes. On reste la source de vérité contextuelle sur le terrain.

**Phrase de positionnement à sortir dans le call :**

> *"On construit la couche de contexte au bord (edge) que les agents IA doivent avoir pour donner des réponses utiles en usine. Vous construisez les agents. On peut se brancher — voyons comment."*

#### À NE PAS révéler dans le call

Ces trois éléments sont nos moats — ils restent confidentiels tant qu'on n'a pas un accord signé :

1. **Le mécanisme du Fuzzy Join OF/batch** (moat #2 — comment on relie les événements OT au contexte batch/produit/OF de l'ERP en temps réel)
2. **Notre thèse "MCP au bord = source de vérité IA"** (moat #3 — c'est notre positionnement long terme, ne pas le tipper à un partenaire potentiel avant contrat)
3. **Le modèle Impact Engine** (moat #1 — calcul de coût multi-facteur, pas juste durée × taux horaire)

Tu peux parler du KG, des pipelines, du MCP en général — mais pas de ce QUI rend chacun unique chez nous.

#### Ce qu'on peut proposer si le call avance bien

**Un pilote 2 semaines**, sans engagement :

| Semaine | Nous | Eux |
|---|---|---|
| S1 | On déploie MindSet sur 1 site test (ou notre simulateur), on expose un MCP au-dessus du KG + TagRegistry | Ils identifient 5 questions types que leur agent doit pouvoir répondre |
| S2 | On adapte le MCP en fonction du feedback | Ils branchent leur agent sur notre MCP, testent les 5 questions |
| Fin | Réunion de retour : est-ce qu'il y a un fit ? | Décision partenariat / non |

Pas de contrat, pas d'engagement. Juste un jour de setup + une semaine d'essai + une réunion.

**Bénéfice pour nous** : on apprend concrètement ce qu'un agent tiers attend d'un KG industriel — c'est de la R&D pour notre MCP, gratuite.

#### TL;DR pour le call

> **Oui**, l'agent aura de quoi raisonner sur nos pipelines — à condition que le MCP expose (a) le KG (b) les valeurs live (c) une traversée, et que leur ontologie soit compatible ou mappée.
>
> **Priorité #1 du call :** comprendre leur ontologie et leur pattern d'accès habituel. La question 5 (*"le KG précédent, il exposait quoi ?"*) te donne 80 % de la réponse sans avoir à révéler notre archi.
>
> **Objectif secondaire :** proposer un pilote 2 semaines sans engagement pour valider le fit techniquement avant de parler contrat.

_Journalisé dans `docs/analysis_log.md` → Entry 57._

---
