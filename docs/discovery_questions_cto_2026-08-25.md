# Questions de discovery — entretiens CTO/tech (2026-08-25)

Prépare l'action item de `docs/workshop.md` (Mohamed) : "Discovery technique : réaliser des entretiens de découverte avec des CTO pour explorer les besoins en données et préparer la stratégie technique."

Construit sur la guidance déjà obtenue de Geneviève (`docs/call_oss_venture.md`) : ne pas pitcher la plateforme, faire un diagnostic. Combine son cadrage IT/DSI avec l'exigence KPI-first du workshop (Jalil) et le maintien des deux ICP en parallèle (usine OT/IT vs économie physique/robotique).

**Rappel de posture** : l'interlocuteur CTO/IT est un champion technique, pas forcément l'Economic Buyer (P&L) — identifier qui décide du budget reste une question ouverte à la fin de l'entretien, pas une hypothèse à faire à l'avance.

---

## A. Ouverture — douleur réelle, pas de pitch

- *"Pourquoi n'avez-vous pas réussi à valoriser toutes vos bases de données jusqu'ici ? Comment vos opérateurs enregistrent-ils la donnée terrain ?"* — question déjà validée par Geneviève, meilleur point de départ.
- Racontez-moi le dernier projet data/IT lancé en usine — comment ça s'est passé ?

## B. Fragmentation IT/OT — teste l'hypothèse silos (podcast, `insights_2026-08-21.md`), sans l'affirmer

- Qui décide côté IT vs côté production/ingénierie sur un projet touchant la donnée machine — mêmes personnes, ou équipes séparées ?
- Un projet a-t-il déjà été bloqué parce que l'IT et le terrain ne parlaient pas le même langage ?
- Où s'arrête la donnée aujourd'hui — reste-t-elle sur l'automate, ou remonte-t-elle jusqu'à l'ERP/le cloud ?

## C. Ce qui a déjà été tenté — teste l'hypothèse "modèle SAP/consultants cassé"

- Déjà lancé une transformation digitale avec un intégrateur/ERP ? Comment ça s'est terminé ?
- Qu'est-ce qui a fait que ça n'a pas tenu ses promesses ?
- Aujourd'hui, connecter une nouvelle source de donnée terrain prendrait combien de temps, réalistement ?

## D. Impact business / KPI — obligatoire (exigence de Jalil, `workshop.md`)

- Quel est le coût réel de ne pas avoir cette donnée en temps réel — qualité, arrêts machine, rebut ?
- Une décision prise trop tard récemment — ça a coûté quoi, concrètement ?
- Top 3 des priorités d'investissement de l'année — la donnée en fait partie ?

## E. Économie physique / robotique — deuxième piste, en parallèle

- La robotique/l'automatisation avancée, dans vos plans à 2-3 ans ?
- Si un robot devait comprendre le contexte de production demain, qu'est-ce qui lui manquerait le plus ?

## F. Proof points — ce qu'il faut obtenir avant de partir (Geneviève)

- Un accès réel (pas une promesse) à une source de donnée qu'on pourrait connecter pour un test rapide ?
- Quelqu'un en interne qui porterait ce sujet si ça avançait ?

## G. Clôture

- Qu'est-ce qui vous ferait dire "on teste" plutôt que "on verra plus tard" ?

---

## Noms identifiés pour ces questions (2026-08-27)

Recherche LinkedIn ciblée — cette liste n'existait pas avant, les questions ci-dessus n'avaient aucun contact attaché.

| Nom | Rôle | Localisation | Pourquoi |
|---|---|---|---|
| **Stéphane Jaud** | CTO — Directeur Technique & Innovation chez **VLAD** | Tours | Titre exact "CTO," profil industriel — le meilleur match direct pour ces questions |
| **Frédéric Kieffer** | DSI en ETI, PME et Association, orienté métier / CTO / Directeur de Projets / Stratégie IT / Transformation digitale / Dette technique / Data Intelligence | Paris | Langage quasi identique au ciblage ETI/PME du produit — "orienté métier" recoupe directement l'angle Geneviève (pas d'IT pour l'IT) |
| **Christophe Fournel** | Expert Data & Digital Transformation, accompagnement stratégique CDO/CIO/CMO et Dirigeants PME/ETI, 25 ans en architecture de données/ETL | Paris | Profil conseil plutôt qu'opérationnel — utile pour challenger l'architecture plutôt que pour un pilote direct |

**Non retenus** : plusieurs profils "CTO" trouvés dans la même recherche (Jean-Yves Cambier, Ludovic Coralie, Loic Caroli) sont positionnés sur du conseil/IA générique ou de l'infogérance IT pure — pas assez spécifiquement industriel/OT pour ces questions précises.

**Rien envoyé** — même posture que le reste des recherches de contacts (voir `robotique_etat_art_workshop_2026-08-31.md` §16).
