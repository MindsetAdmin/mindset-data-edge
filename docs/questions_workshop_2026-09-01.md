# Les 9 questions fondamentales — préparation workshop (2026-09-01)

Document unique regroupant **toutes les questions de fond** du workshop. Remplace et consolide `questions_fondamentales`, `solutions_build_buy_conseil` et `scalabilite_productisation` (supprimés — intégralité du contenu reprise ici).

Ce ne sont pas des questions de pitch : ce sont celles qu'un CTO ou un responsable industriel pose en premier, et une réponse floue sur l'une d'elles tue la conversation.

| # | Question | Partie |
|---|---|---|
| **Q1** | Pourquoi ce problème n'est-il pas résolu ? | I — Le problème |
| **Q2** | Que font les solutions actuelles ? | II — Le paysage |
| **Q3** | Build ou buy ? Pourquoi nous ? | II |
| **Q4** | Pourquoi pas des consultants ? | II |
| **Q5** | Pourquoi ça ne passe pas à l'échelle ? | III — Scalabilité |
| **Q6** | Comment produitiser la dernière couche ? | III |
| **Q7** | Que manque-t-il aux agents IA ? | III |
| **Q8** | Que manque-t-il aux robots ? | III |
| **Q9** | Ils ont un data lake — pourquoi nous ? | III |

**Format** : **Niveau 1** = la réponse courte, à voix haute. **Niveau 2** = la substance technique et les chiffres, si on est challengé.

**Statut** : tout ce qui est chiffré est sourcé (sources en fin de document). Ce que la recherche documentaire ne peut pas trancher est marqué **[À VALIDER]** et regroupé en fin de document.

> **⚠️ Périmètre : piste OT/IT usine uniquement** (l'offre initiale — UNS / ISA-95 / OPC-UA), sauf Q8.
> Les réponses **diffèrent selon la piste** et les paysages concurrentiels n'ont rien à voir : HighByte / Litmus / UMH pour l'OT/IT · Resilinc / Interos / Everstream / Prewave pour la supply chain (`tarik.md` §4) · quasiment personne pour la robotique (`robotique_etat_art_workshop_2026-08-31.md`). Ne pas réutiliser ces réponses pour une autre piste sans les refaire.

---
---

# §0. Les deux découvertes qui changent les réponses

**À lire avant tout le reste.** Ces deux constats, établis les 31/08 et 01/09, modifient plusieurs réponses ci-dessous — et l'un d'eux **invalide une revendication écrite dans `docs/mindset.md`**. Il vaut infiniment mieux les porter nous-mêmes au workshop que de les découvrir dans un appel client.

### A. HighByte a livré avant nous — ce que ça invalide

**HighByte Intelligence Hub version 4.2 (juillet 2025) embarque un serveur MCP industriel *au bord*, et génère des instances de modèle depuis un espace d'adressage OPC UA à l'aide d'un LLM, avec validation humaine avant enregistrement.**

Précisément, et vérifié :

| Ce que fait HighByte 4.2 | Ce qu'on revendiquait |
|---|---|
| **Serveur MCP industriel embarqué**, positionné explicitement « at the edge », exposant les pipelines de données comme *tools* aux agents IA (descriptions + paramètres). Annoncé comme **le premier serveur MCP industriel** | *« MindSet is the only edge MCP »* (`mindset.md` L589) — **factuellement faux** |
| Bouton **« AI Generate Instances »** : un LLM parcourt l'espace d'adressage OPC UA, trouve les instances, **validées avant enregistrement** | « auto-dérivation du modèle depuis la découverte, avec porte de validation » — **même mécanisme, approche différente** |
| Contextualisation assistée par LLM avec connexions natives Bedrock, Azure OpenAI, Gemini, OpenAI **et LLM locaux** | La souveraineté « la donnée ne sort pas » n'est donc **pas** un différenciateur automatique |
| **17 500 $ / site / an** (Professional, mono-site), passé à **18 500 $** en 2026 ; Enterprise sur devis. MCP Services inclus dans **toutes** les licences. Essai gratuit 30 jours | Nous n'avons pas de prix public |

**Ce que ça invalide dans nos docs** : `docs/mindset.md` L589 (« MindSet is the only edge MCP ») et le Moat #5 en L1262-1263 (« Cognite a MCP mais cloud-only… UMH n'en a pas ») — l'inventaire concurrentiel y est incomplet, HighByte n'y figure pas alors qu'il a livré avant nous.

### Ce que ça ne change pas

À poser aussi clairement, sinon la lecture devient défaitiste à tort :

1. **Le marché est validé, pas fermé.** Qu'un acteur sérieux ait construit exactement ça, l'ait tarifé à 18,5 k$/site/an et le vende, prouve que le problème est réel et budgété. C'est l'inverse d'une mauvaise nouvelle sur la *demande*.
2. **La catégorie est jeune et en croissance** — Cybus, HighByte, Litmus, Soffico, UMH cités ensemble en forte dynamique commerciale, la plupart des acteurs ayant moins de dix ans.
3. **Un différenciateur de mécanisme survit, et il est précis** : HighByte génère les instances **avec un LLM** ; notre score de confiance est une **formule pondérée déterministe**. Ce n'est pas « mieux » dans l'absolu — c'est *auditable*. Dans un contexte où il faut justifier pourquoi un nœud a été accepté (pharma régulée, qualification fournisseur, défense), « le LLM a proposé » et « voici la formule, les poids et les entrées » ne se défendent pas de la même façon.
4. **La juridiction reste un axe réel.** HighByte est américain. L'argument souveraineté ne peut plus être « la donnée ne sort pas » (leurs LLM locaux le permettent aussi) — il se resserre sur **la juridiction du fournisseur** et l'exposition CLOUD Act au niveau de la relation contractuelle. Plus étroit, toujours vrai, et décisif pour secteur public / défense / pharma régulée en France.

### Conséquence directe sur le pitch

> **La phrase « personne ne fait ce qu'on fait » est morte. Elle doit être retirée de tous les supports.**

La question honnête n'est plus *« sommes-nous les seuls ? »* mais *« pourquoi nous plutôt que HighByte ? »* — et §Q3 y répond sans bluffer.

---

### B. Ils vendent, et passent quand même par des intégrateurs

> *« HighByte, Litmus et les autres ont déjà tenté de produitiser cette couche. Ils existent, ils vendent, et ils passent quand même massivement par des intégrateurs. Ce n'est pas un hasard — c'est le signe que la dernière couche de spécificité résiste à la standardisation. »*

**Vérifié le 01/09 en consultant directement les pages**, pas seulement des résultats de recherche. Quatre éléments, classés par solidité — parce qu'ils ne se valent pas :

**1. La preuve la plus forte vient des intégrateurs eux-mêmes, pas des éditeurs.** La **CSIA** (Control System Integrators Association — l'association professionnelle des intégrateurs, indépendante de tout éditeur) décrit le travail d'un déploiement UNS en revendiquant explicitement la couche sémantique comme son métier :

> *« Normalizing inconsistent data from various machine vendors. »*
> *« Successful UNS deployments require publishing **meaningful, contextualized** data from PLCs, DCSs, SCADA systems, smart devices and legacy equipment, **not just raw tags**. »*
> *« Integrators understand how equipment truly functions. **We know how ISA-95 models connect to real plant hierarchies and how to create naming and data structures that mirror actual production environments.** »*

C'est exactement la « dernière couche de spécificité » — et ce sont les intégrateurs qui la revendiquent comme leur valeur propre. Une association professionnelle qui décrit son propre métier est une source bien plus solide qu'une page marketing d'éditeur.

**2. HighByte le formule dans les mêmes termes.** Son offre Industrial Data Fabric est commercialisée via **Deloitte, Infosys, Cyient, TensorIoT**, décrits dans son propre billet comme : *« These boots on the ground are the last mile between concept and reality. »* Sa page partenaires distingue trois catégories (technologie, distributeurs, intégrateurs) et attribue aux intégrateurs le *« design and deployment »*, le conseil et l'intégration. Litmus passe aussi par des intégrateurs (GFT), et HighByte a noué un partenariat avec Siemens.

**3. Signal faible, mais cohérent : HighByte n'internalise pas ce travail.** Sa page carrières n'affichait le 01/09 **qu'un seul poste ouvert — un Account Executive**, donc commercial. Aucun poste d'ingénieur d'implémentation, d'architecte solution ou de services professionnels. Un éditeur qui voudrait absorber le dernier segment recruterait pour le faire ; celui-ci recrute pour vendre. *À pondérer : instantané sur une petite société, et leur page indique que d'autres recrutements peuvent exister hors annonces.*

### Le contre-exemple qu'il faut porter aussi — UMH

**UMH affirme exactement l'inverse**, et l'honnêteté impose de le dire : *« First machine connection takes about 90 seconds »*, *« 18 min new machine connected with templates »*, *« <5 days from idea to live use case »*, une console de gestion *« accessible to OT engineers without writing code »* — et **aucune mention d'intégrateurs ni de services professionnels** sur leur page UNS.

**Ce qui sauve quand même la thèse, et c'est une distinction réelle, pas un sauvetage** : ces chiffres mesurent le temps de **connexion**, pas le temps de **modélisation sémantique**. Connecter une machine en 90 secondes ne dit rien sur le temps qu'il faut pour se mettre d'accord sur ce que ses tags signifient. UMH ne quantifie nulle part l'effort de définition des standards de données, de cartographie de l'existant, ni de conduite du changement multi-sites — et annonce par ailleurs **4 à 6 semaines** pour un pilote de production multi-machines, ce qui n'est pas rien.

### Ce que ça prouve — et ce que ça ne prouve pas

✅ **Prouvé** : la couche sémantique (mapping ISA-95 → hiérarchie réelle du site, conventions de nommage, normalisation inter-fournisseurs) est aujourd'hui **du travail humain facturé**, revendiqué comme tel par les intégrateurs eux-mêmes.

⚠️ **Non prouvé** : que les éditeurs passent par des intégrateurs **parce que** ce segment résiste à la standardisation. Le canal intégrateur est aussi la façon normale de vendre du logiciel d'entreprise — couverture géographique, langue, relation client existante, transfert de risque. Salesforce et SAP passent par Deloitte sans que personne n'en conclue que leur produit est inachevé. **Le recours aux intégrateurs est cohérent avec l'hypothèse, il ne la démontre pas.**

**À dire ainsi au workshop** : *« la couche sémantique est du travail humain facturé, c'est établi. Que ce soit irréductible, ça reste à démontrer — et c'est précisément ce qu'on propose d'attaquer. »*

---

---
---

# PARTIE I — LE PROBLÈME

## Q1. Pourquoi ce problème n'est-il pas résolu ?

**Niveau 1** — Parce que le problème que tout le monde a essayé de résoudre n'est pas le vrai problème. On a standardisé le **transport** de la donnée — c'est fait, ça marche, OPC-UA existe depuis 2008. Ce qui n'a jamais été résolu, c'est le **sens** : savoir que `ns=2;s=Channel1.Device1.Tag47 = 72.3` est la température du four 3 de la ligne 2. Ce mapping-là a toujours été fait à la main, plante par plante, et c'est là que les projets meurent.

**Niveau 2** — Sept raisons techniques, dans l'ordre où elles mordent :

### 1. Le transport est résolu, la sémantique ne l'est pas

Nuance importante à porter correctement, parce qu'un bon technicien corrigera : **OPC-UA ne se limite pas au transport.** Il spécifie un modèle d'information (structure, comportement, sémantique), et les **Companion Specifications** existent précisément pour standardiser la sémantique par domaine métier.

Mais dans la pratique : les structures propriétaires et les conventions de nommage incohérentes ne passent pas à l'échelle. Quand des équipements de plusieurs fournisseurs arrivent chacun avec ses propres structures, conventions de nommage, unités de mesure et méthodes d'intégration, **sans modèle de données OPC UA standardisé, les ingénieurs ne peuvent pas compter sur des structures prévisibles**. Les Companion Specs et NodeSets personnalisés commencent seulement à être davantage utilisés — c'est un chantier en cours, pas un problème réglé.

**Formulation juste** : *« la couche sémantique existe dans la norme ; elle n'existe pas dans le parc installé. »*

### 2. Le brownfield — on ne peut pas remplacer, seulement s'adapter

Un équipement industriel a un cycle de vie de 20 à 30 ans. Un automate installé en 1998 tourne encore et produit. Personne ne remplace un parc pour un projet data. Toute solution doit donc composer avec l'existant : **N connecteurs, pas un standard**. C'est la raison la plus structurelle, et elle est physique — elle tient à l'amortissement du capital, pas à un mauvais choix d'architecture.

### 3. La modélisation est le vrai coût, et c'est mesuré

**Le poste le plus lourd d'une implémentation UNS est la modélisation hiérarchique des actifs : 40 à 60 % de l'effort total.** Pas la connectique, pas le broker, pas le stockage — le fait de décider quelle donnée correspond à quel équipement, dans quelle hiérarchie.

C'est le chiffre le plus important de tout ce document : il dit que le goulot est exactement là où presque personne n'automatise.

### 4. Le mapping manuel ne passe pas à l'échelle — et il périme

Une usine expose typiquement 10 000 à 100 000 tags. Les mapper à la main, c'est des mois d'intégrateur, ce qui coûte souvent plus cher que le logiciel lui-même. Et le résultat **se périme** : on ajoute une machine, on renomme un tag, un automate est remplacé — le modèle diverge silencieusement du terrain.

Conséquence économique directe : le coût d'intégration dépasse la valeur d'un cas d'usage isolé. C'est pour ça que les projets s'arrêtent au pilote — pas parce que la techno ne marche pas.

### 5. Personne ne possède le problème horizontal

Structurellement, l'IT possède les niveaux 4-5, l'OT possède les niveaux 0-2, et **le niveau 3 n'appartient à personne**. Les budgets suivent les silos. Celui qui ressent la douleur (l'ingénieur, le responsable de site) ne tient pas le budget ; celui qui tient le budget (DSI, C-suite) ne ressent pas la douleur. Problème d'incitations divisées classique — aucun acteur n'a à la fois le mandat et la motivation.

Le podcast (`insights_2026-08-21.md`) formule la même chose côté achat : **quatre silos qui n'achètent pas de la même façon** — C-suite (stratégie), directeurs de site (outils de décision), DSI (sécurité, conformité, TCO), ingénieurs (équipement OEM sur budget CapEx). Siemens ne vend pas au C-suite, il vend aux ingénieurs. Le logiciel choisi « en haut » ne correspond pas à ce dont le terrain a besoin.

### 6. Les régimes de temps réel sont incompatibles

*« Le temps réel en OT se compte en millisecondes, en IT en secondes ou minutes. »* Une seconde sur le terrain peut valoir très cher ; côté IT, une seconde — voire une heure — reste une dépense justifiable. Une solution architecturée à la vitesse et à la logique IT **ne peut pas** servir un besoin OT. Ce n'est pas un défaut d'implémentation, c'est une incompatibilité de conception.

Corollaire utile : l'ERP est un **system of record** (il enregistre), pas un **system of ownership** (le système réellement autoritaire au moment où la donnée compte). Forcer toute la donnée usine à transiter par l'ERP avant d'être exploitable, c'est se tromper de propriétaire.

### 7. La fragmentation est rationnelle pour les fournisseurs

À dire avec précaution, mais c'est réel : l'interopérabilité réduit les coûts de changement de fournisseur, donc la marge. Aucun grand acteur de l'automatisation n'a d'intérêt économique direct à rendre son parc trivialement interchangeable. La fragmentation n'est pas un accident historique qu'on n'aurait pas encore eu le temps de corriger.

### Et ce qui a changé récemment — la vraie réponse à « pourquoi maintenant ? »

C'est la partie la plus importante de Q1, parce que « le problème existe depuis toujours » appelle immédiatement « alors pourquoi vous, maintenant ? ». Trois choses ont bougé, et aucune n'était vraie il y a trois ans :

1. **Le mapping sémantique est devenu automatisable.** Avant les LLM, faire correspondre `Usine_Paris.L2.M3.temp_four` à un modèle ISA-95 demandait des règles écrites à la main, par plante. C'est précisément le poste à 40-60 % de l'effort. C'est la première fois qu'il est attaquable autrement qu'à la main.
2. **Le compute a migré au bord.** Faire tourner de la contextualisation sérieuse sur un boîtier de plancher d'usine est devenu banal.
3. **L'IA agentique a rendu le problème visible.** La poussée pour déployer de l'IA générative et agentique sur le plancher d'usine **a exposé les limites de la télémétrie brute** — un agent branché sur des tags sans contexte ne produit rien d'utile. Le problème sémantique, longtemps ressenti comme un inconfort d'intégration, est devenu un blocage visible avec un budget en face.

**Phrase à dire** : *« Le problème n'a pas changé. Ce qui a changé, c'est que la partie la plus coûteuse — établir le sens de la donnée — vient de devenir automatisable, et qu'en même temps l'IA a rendu son absence impossible à ignorer. »*

---

---
---

# PARTIE II — LE PAYSAGE ET LA DÉFENSE CONCURRENTIELLE

## Q2. Que font les solutions actuelles ?

**Niveau 1** — Six familles. Cinq résolvent une couche adjacente et laissent la modélisation au client. **Une famille — l'Industrial DataOps — attaque exactement notre problème, et un de ses acteurs a livré avant nous.**

**Niveau 2** —

| Catégorie | Exemples | Ce qu'elles font bien | Où ça s'arrête |
|---|---|---|---|
| **Historians** | OSIsoft PI, AVEVA Historian | Stockage time-series massif, fiable, éprouvé | Un tag reste un tag — peu de sémantique. Coûteux, piloté par l'IT |
| **Éditeurs SCADA / MES** | Siemens, Rockwell, AVEVA, Opcenter | Excellents dans leur couche verticale, connaissance métier réelle | Stack verticale : aucun intérêt économique à faciliter la sortie de leur périmètre |
| **iPaaS / intégration IT** | MuleSoft, Boomi | Connecteurs réutilisables, orchestration mature | Monde IT : mauvais protocoles, mauvais régime de latence |
| **Industrial DataOps / UNS** | **HighByte**, Litmus, Cybus, United Manufacturing Hub, Soffico | **Le champ de bataille réel.** Modélisation, transformation, gouvernance. Litmus : 250+ connecteurs. HighByte : MCP embarqué + génération d'instances assistée par LLM | C'est ici qu'il faut se différencier finement — voir §0 et le détail ci-dessous |
| **Plateformes cloud industrielles** | Cognite (Leader IDC MarketScape 2026), Palantir Foundry, Databricks | Très puissantes, écosystème, contextualisation avancée | Cloud-first, coûteuses, déploiement long, contextualisation largement en prestation |
| **Intégrateurs / cabinets** | Accenture, Capgemini, intégrateurs locaux | Font le travail, connaissent le site | Livrent un projet, pas un produit — Q4 |

### Le détail qui compte — les trois acteurs à connaître nommément

**HighByte** *(le plus proche — à traiter comme la référence)*
Produit commercial, licence **18 500 $/site/an** en 2026 (Professional mono-site), Enterprise sur devis. Modèles et pipelines illimités, broker MQTT embarqué, client UNS, serveur REST, haute disponibilité, intégration PI — **et MCP Services inclus dans toutes les licences**. Depuis la 4.2 : serveur MCP industriel au bord, génération d'instances par LLM depuis l'espace OPC UA, intégration Git, OpenTelemetry, connecteurs Databricks et TimescaleDB. Essai gratuit 30 jours.
→ *C'est le produit contre lequel nous serons comparés. Le connaître précisément est non négociable.*

**United Manufacturing Hub** *(la pression par le bas)*
**Édition communautaire open-source et gratuite** ; le cœur est un chart Helm pour Kubernetes, plus une console de gestion en SaaS managé. Positionné pour les équipes « engineering-led » qui veulent auto-héberger et étendre.
→ *Le risque n'est pas qu'ils gagnent le deal, c'est qu'ils fixent le prix plancher à zéro et transforment la conversation en « pourquoi payer ? ». La réponse est l'exigence Kubernetes et le coût de possession — pas le dénigrement d'un projet open source sérieux.*

**Litmus** *(la couverture)*
**250+ connecteurs prêts à l'emploi**, plateforme edge, forte croissance.
→ *Terrain sur lequel nous ne devons pas nous battre. Ne jamais argumenter l'étendue du catalogue.*

### Position honnête à porter au workshop

Nous ne sommes **pas** différenciés sur le mécanisme. Nous sommes potentiellement différenciés sur :
- **l'approche** (déterministe et auditable vs assistée par LLM),
- **la juridiction** (européenne vs américaine),
- **la simplicité de déploiement** (binaire unique vs produit d'entreprise, et vs Kubernetes côté UMH),
- **la réconciliation OT↔IT native** — **[À VALIDER]**, je n'ai pas vérifié si HighByte fait de la résolution d'entité équipement↔ERP.

Trois de ces quatre points sont plausibles mais non vérifiés. **C'est le vrai travail restant avant de construire un argumentaire de vente.**

---

---

## Q3. Build ou buy ?

**Niveau 1** — La question n'est plus binaire. Elle est à **trois** branches : construire soi-même, acheter un produit du marché (HighByte à 18,5 k$/site/an, ou UMH gratuit), ou nous. Prétendre que « build vs buy » se résume à « nous vs vos ingénieurs » est faux depuis §0.

**Niveau 2** —

### Ce qui est facile à construire — le concéder immédiatement

Lire des tags OPC-UA, publier sur MQTT, brancher un Grafana : un bon ingénieur automatisme fait ça vite. **Le dire franchement.** Prétendre le contraire face à un technicien fait décrocher la conversation, parce qu'il sait que c'est faux.

### La longue traîne — où va réellement le temps

- **Cas particuliers protocolaires** : modes de sécurité OPC-UA et certificats, limites de souscription, profils partiellement implémentés, bizarreries par fournisseur.
- **Le mapping sémantique qui survit au changement** — le poste à **40-60 % de l'effort d'implémentation UNS**, et il ne s'agit pas de le faire une fois mais de le maintenir vrai quand le parc bouge.
- **Le workflow de validation** : personne ne valide 10 000 nœuds à la main. Sans porte de confiance, l'auto-génération est un déversement de bruit. *(Noter : HighByte a exactement la même conviction — leurs instances générées sont « validated before saving ».)*
- **Robustesse réseau** : store-and-forward, coupure, reprise.
- **Multi-site et versionnement du modèle.**
- **La maintenance** : qui corrige quand la personne qui l'a écrit est partie ?

### Les chiffres du build

- **La maintenance et le support représentent 65 à 85 % du coût de possession** sur la durée de vie d'un logiciel industriel.
- Ordres de grandeur observés : **850 K$ à 1,6 M$ d'investissement initial, 18 à 24 mois avant production, puis 480 K$ à 950 K$/an.** Sur un cas de parité fonctionnelle plus ambitieux : **13 M$ (144 mois-ingénieur) + 4 M$/an**.
- **Dépassement courant de 40 à 60 %** par rapport à l'estimation initiale.

### Le risque humain — le vrai tueur

Les systèmes internes **deviennent orphelins quand l'ingénieur part**. Il faut maintenir les intégrations quand l'ERP se met à jour, patcher la sécurité, ajouter le champ qu'un client réclame, et reconstituer la connaissance perdue. S'y ajoute un problème de recrutement réel : **les ingénieurs seniors évitent les équipes d'outils internes hérités.**

### La comparaison à trois branches, avec le prix comme ancre

| | Build interne | HighByte | UMH (communautaire) |
|---|---|---|---|
| **Coût d'entrée** | 850 K$-1,6 M$, 18-24 mois | **18 500 $/site/an**, essai 30 j | 0 € de licence |
| **Coût récurrent** | 480-950 K$/an | licence | infrastructure + compétence Kubernetes |
| **Qui maintient** | un ingénieur, jusqu'à son départ | l'éditeur | vous |
| **Risque principal** | orphelinat, dépassement 40-60 % | dépendance éditeur, juridiction US | charge d'exploitation K8s, pas de support contractuel |

**L'argument qui tient face à un technicien** : à 18 500 $/an, un produit du marché coûte moins qu'**un mois** d'ingénieur senior chargé. Le débat « build vs buy » sur ce périmètre est économiquement tranché — et il l'est *contre le build*, pas contre nous. **Notre concurrent réel n'est pas l'équipe interne du client, c'est HighByte.**

### Quand « build » reste le bon choix — le dire crédibilise tout le reste

- **Un site, un cas d'usage, quelques dizaines de tags** — construire est probablement correct. Ne pas essayer de vendre là.
- **Une équipe logicielle interne de 50+ personnes**, capable d'absorber la maintenance.
- **Un besoin qui est le différenciateur concurrentiel du client** — il doit le posséder.

**Le critère qui tranche** : le nombre de cas d'usage prévus sur trois ans. Un → build. Cinq, multi-sites, parc mouvant → le coût de possession bascule.

### Pourquoi nous — version honnête, post-§0

Ne revendiquer que ce qui est construit, testé, et qui survit à une comparaison avec HighByte :

1. **Scoring déterministe, pas LLM.** Notre score de confiance est une formule pondérée traçable : entrées nommées, poids fixes, résultat reproductible. HighByte génère ses instances via un LLM. Pour un acheteur qui devra justifier une décision (pharma régulée, défense, qualification fournisseur), « voici la formule » et « le LLM a proposé » n'ont pas la même valeur. **C'est notre meilleur argument technique restant, et il est réel.**
2. **Juridiction européenne + pas d'édition hyperscaler.** HighByte est américain. Pour le secteur public FR, la défense et la pharma régulée, l'exposition CLOUD Act au niveau du fournisseur est un critère d'exclusion, indépendamment de la localisation de la donnée.
3. **Simplicité de déploiement.** Binaire unique, sans Kubernetes (vs UMH), sans produit d'entreprise à opérer. **[À VALIDER]** face à HighByte, dont je ne connais pas la complexité d'installation réelle.
4. **Réconciliation OT↔IT native** — liens équipement OT ↔ données ERP calculés et persistés. **[À VALIDER]** : non vérifié chez HighByte.

**Ce qu'il ne faut plus jamais dire** : « nous sommes les seuls à avoir MCP au bord », « personne ne fait ça », « l'auto-génération du modèle est unique ». Les trois sont faux depuis §0.

---

---

## Q4. Pourquoi ne pas simplement payer des consultants ?

**Niveau 1** — Parce qu'un consultant livre un **projet**, pas un **produit**. Quand il part, rien ne se capitalise : le cas d'usage suivant est un nouveau devis. Le coût croît linéairement avec les cas d'usage ; un produit les amortit.

**Niveau 2** —

### Les chiffres

- **Intégration sur mesure : 50 K$ à 500 K$+ par projet** en 2026. Intégration ERP : 80 K$ à 300 K$. Programme multi-systèmes : 200 K$ à 1,2 M$.
- **Maintenance annuelle : 20 à 35 % du coût de build initial.**
- **Dépassement de 40 à 60 %** par rapport à l'estimation.
- Le modèle intégrateur convient aux intégrations ponctuelles, très spécifiques et profondes, **mais est mal adapté à un parc d'intégrations qui grandit et change fréquemment** — parce qu'il repose sur du développement spécifique plutôt que sur des composants réutilisables.

**La comparaison qui frappe** : *une seule* intégration sur mesure au bas de la fourchette (50 K$) coûte davantage que **deux ans et demi** de licence HighByte. Et elle ne couvre qu'un cas d'usage.

### Le problème structurel, pas seulement le prix

Les grands cabinets vendent **« du message et de la promesse stratégique »** au C-suite. Leurs plans à cinq ans sont directionnels, pas exécutables — *« planning is essential, plans are worthless »*. Une fois le plan livré, l'IT et les ingénieurs internes doivent le découper et aller chercher eux-mêmes les logiciels : **deux processus d'achat séparés pour ce qui devrait être un seul système cohérent.**

Point plus dur : les consultants **captent la confiance** des ingénieurs — ceux qui connaissent réellement le terrain — au lieu de la construire. La relation technique qui compte se fait d'ingénieur à ingénieur.

Ligne citable : ***« Mindset and [similar solutions] sell technology, they are not consulting. »***

### Le contexte d'échec

70 % des transformations digitales échouent (analyse bibliométrique 2025) ; 69 % (McKinsey) ; 35 % seulement atteignent leurs objectifs (BCG, 850+ entreprises) ; 85 % ne dépassent jamais le pilote (Gartner). SAP/ERP : 55-75 % d'échec, dépassement moyen de 215 %.

**Précaution** : donner la fourchette et la source, jamais un « 80 % » agrégé.

### Ce que les consultants font mieux que nous — à dire

Conduite du changement, refonte de processus, alignement organisationnel, gestion politique d'un projet multi-directions. **Nous ne faisons rien de tout ça.** Et le rappel de Daouda vaut ici : *« un outil seul ne suffit pas, il faut toujours un humain pour surveiller et corriger »* — le problème est autant humain et procédural que technique.

### Le bon cadrage : canal, pas concurrent

La formulation n'est pas « au lieu d'un consultant » mais **« ce que le consultant installe »**. Un intégrateur qui déploie un produit livre en semaines ce qu'il livrerait en mois en spécifique, et son client garde quelque chose de maintenable après son départ.

*Nuance à assumer : c'est exactement l'argument que HighByte tient à ses propres intégrateurs. Le canal n'est pas un espace vide non plus.*

---

---
---

# PARTIE III — SCALABILITÉ ET PRODUCTISATION

## Q5. Pourquoi ça ne passe pas à l'échelle

**Niveau 1** — Parce que ce qui coûte n'est pas de connecter, c'est de **s'engager sur le sens**. Connecter un automate est générique et se produitise très bien — c'est déjà fait. Décider que `Ligne2.M3.tmp_four` est la température du four de la ligne 2, que ce four est le goulot, et qu'au-dessus de 180 °C c'est une dérive — ça, c'est spécifique à un site, et ça vit dans la tête de trois personnes.

**Niveau 2** — La décomposition qui explique tout :

| Couche | Générique ? | Produitisée aujourd'hui ? |
|---|---|---|
| Transport (protocoles, connectivité) | Oui | **Oui**, largement résolu |
| Stockage, transformation, gouvernance | Oui | **Oui** (HighByte, Litmus, historiens) |
| **Engagement sémantique propre au site** | **Non** | **Non — c'est le dernier segment, confié à l'intégrateur** |
| Logique métier (ce qui compte comme arrêt, quel est le goulot) | Non | Non, et probablement jamais entièrement |

**L'équation de scalabilité, et c'est la formulation à retenir :**

> Aujourd'hui, l'effort humain est proportionnel au **volume** — O(n) sur les tags côté OT, et jusqu'à **O(n_OT × n_IT)** sur la jonction OT/IT (cf. Q6). Une usine expose 10 000 à 100 000 tags, la modélisation hiérarchique représente **40 à 60 % de l'effort** d'une implémentation UNS, donc l'effort croît avec la taille du parc. C'est structurellement non scalable.
>
> **Pour que ça scale, l'effort humain doit devenir proportionnel au nombre d'éléments réellement ambigus — O(u), avec u ≪ n.**

Tout le reste de ce document est une réponse à : *comment fait-on passer n à u, techniquement.*

**Deuxième raison, aussi importante et souvent oubliée** : même un modèle parfait **se périme**. On ajoute une machine, on renomme un tag, on remplace un automate — le modèle diverge silencieusement du terrain. Sans mécanisme de détection de dérive, chaque modèle pourrit, et l'intégrateur doit revenir. **C'est ce qui transforme un produit en prestation récurrente.**

---

---

## Q6. Comment produitiser la dernière couche — la réponse technique

**Niveau 1** — On ne peut pas standardiser **le contenu** : chaque site est réellement différent, ce n'est pas une excuse. Mais on peut standardiser **le processus de convergence vers ce contenu**. C'est une affirmation complètement différente, et c'est celle qui est produitisable.

### D'abord : de quelle couche parle-t-on ? (correction du 02/09)

**Erreur à ne pas reproduire** : réduire la convergence OT/IT au mapping des tags. C'est la moitié OT du problème, et sa couche la plus superficielle. Il y a **cinq couches**, et une seule est réellement *de la convergence* :

| Couche | Quoi | Nature | État |
|---|---|---|---|
| **1. Structure OT** | tags → équipement → hiérarchie site | **OT pur**, pas de la convergence | Auto-dérivée depuis la découverte, scorée |
| **2. Structure IT** | tables ERP/MES → objets canoniques (OF, produit) | **IT pur**, pas de la convergence | Même mécanisme, construit côté SQL |
| **3. La jonction OT↔IT** | `Machine1` (OPC-UA) = `machine1` (colonne ERP) = `M-001` (MES) | **C'est ÇA, la convergence** | Correspondance exacte normalisée — **pas floue** |
| **4. Alignement temporel** | événement OT à 14:32:07 vs saisie ERP en fin de poste | Convergence | Lecture de l'état ERP plutôt que jointure sur horodatage |
| **5. Sémantique métier** | ce qui compte comme arrêt, quel est le goulot | Ni l'un ni l'autre | **Irréductible**, reste humain |

**Pourquoi la couche 3 est le vrai sujet** : c'est elle qui produit le sens métier. Sans elle, un arrêt machine reste un arrêt machine. Avec elle, c'est *un arrêt sur l'OF 4412, produit X, client Y, livraison dans 3 jours* — et c'est seulement à ce moment-là que la donnée vaut quelque chose (cf. Q9, le terme `t_comprendre`).

**Et l'argument de scalabilité y est plus fort, pas plus faible :**

- Couche 1 : effort en **O(n)** — n tags à qualifier.
- Couche 3 : effort en **O(n_OT × n_IT)** dans le pire cas — chaque entité OT à rapprocher de chaque entité IT candidate. **Combinatoirement pire.**

Donc la porte de confiance n'est pas un confort à la jonction, elle y est **indispensable**. Et les priors inter-sites y sont plus puissants : ce qui se répète d'un site à l'autre, ce ne sont pas les noms, ce sont les **motifs d'écart** entre nommage OT et nommage ERP (l'ERP en codes courts majuscules, l'OT en texte libre hiérarchique, etc.).

**La limite honnête, et elle est ici — pas au niveau des tags** : la résolution d'entité fonctionne aujourd'hui par correspondance **exacte** normalisée (insensible à la casse). Donc ça marche quand les noms se ressemblent, et **pas du tout quand ils ne se ressemblent pas** — c'est-à-dire le cas fréquent en usine réelle. **C'est là que l'intégrateur gagne encore aujourd'hui**, parce que lui sait que `M-001` c'est la ligne 2.

**Niveau 2 — cinq mécanismes, dans l'ordre où ils réduisent l'effort.** Ils s'appliquent aux couches 1, 2 **et 3** — c'est à la couche 3 qu'ils comptent le plus :

### 1. Auto-dérivation depuis la découverte, pas depuis l'entretien
L'intégrateur part d'une page blanche et interroge les gens. Un produit part de **ce qui est réellement sur le réseau** : parcourir l'espace d'adressage, extraire la structure implicite des noms de tags, proposer une hiérarchie candidate. On ne supprime pas l'humain — on lui donne un brouillon à corriger au lieu d'un formulaire à remplir.
*Effet : supprime le coût de démarrage, pas encore le coût de volume.*

### 2. Score de confiance par nœud — le vrai levier O(n) → O(u)
**C'est le mécanisme central.** Chaque proposition reçoit un score. Au-dessus d'un seuil, elle est acceptée automatiquement. En dessous, et **seulement en dessous**, elle passe en revue humaine.
L'humain ne valide plus 10 000 nœuds : il valide les 200 sur lesquels la machine est réellement incertaine. **L'effort devient proportionnel à l'ambiguïté, pas au volume** — c'est la définition mathématique du passage à l'échelle sur ce problème.
*Condition de crédibilité : le score doit être calibré. Un score qui dit 0,9 doit avoir raison ~90 % du temps, sinon le seuil ne veut rien dire et l'auto-acceptation devient du bruit accepté silencieusement — pire que pas d'automatisation du tout.*

### 3. Capturer la correction comme donnée, pas comme conversation
Quand l'intégrateur corrige un mapping, la connaissance repart avec lui. Quand un produit capture la correction, elle devient **une donnée étiquetée** : « sur ce site, `tmp_` signifie température ; `L2` est une ligne, pas un local ». Chaque correction enrichit le modèle de priors.
*C'est ici que la prestation devient un actif.*

### 4. Détection de dérive — ce qui empêche le modèle de pourrir
Le système doit détecter qu'il a divergé du terrain : tag apparu, tag disparu, distribution de valeurs qui change, nœud qui n'émet plus. Sans ça, retour à l'intégrateur tous les 18 mois.
*C'est la différence structurelle entre un projet et un produit : un projet livre un état, un produit maintient un état.*

### 5. Priors inter-sites — le seul mécanisme qu'un intégrateur ne peut pas répliquer
**La réponse de fond à « comment produitiser ce qui résiste à la standardisation ».**

Les conventions de nommage sont différentes d'un site à l'autre, mais elles ne sont **pas aléatoires** : elles se regroupent par secteur, par stack constructeur (un site Siemens ne nomme pas comme un site Rockwell), et par intégrateur ayant câblé le site. Ces régularités sont apprenables.

Une correction faite sur le site A améliore la proposition initiale sur le site B. **L'effort marginal par site décroît à mesure que le parc de clients grandit** — alors que pour un intégrateur, chaque mission repart de zéro parce que l'apprentissage reste dans la tête d'un consultant.

> C'est la seule asymétrie structurelle réelle entre un produit et une prestation. **Ce n'est pas « on est plus rapides » — c'est « on s'améliore, eux se répètent ».**

*Contrainte de confidentialité, déjà résolue par ailleurs* : on ne partage jamais de donnée client entre sites. On partage des **priors agrégés** — des motifs, pas des données. C'est exactement le mécanisme IMDS déjà retenu pour la supply chain (`tarik.md`) : remonter un signal agrégé sans divulguer l'identité derrière.

### La limite honnête — à dire, sinon on se fait démonter

**Une partie de la spécificité est irréductible.** La logique métier propre au site — ce qui compte comme un arrêt, quelle ligne est le goulot, ce que « bon » veut dire pour ce produit — ne sera probablement jamais auto-dérivable, parce que ce n'est pas dans la donnée : c'est une décision d'organisation.

**L'objectif défendable n'est pas de supprimer le dernier segment, c'est de le faire passer de mois à jours.** Et si on y arrive, l'intégrateur n'est pas éliminé — il fait cinq fois plus de sites par an. Ce qui rejoint exactement la position de Q4 : **canal, pas concurrent.**

---

---

### Est-ce seulement faisable ? — évaluation d'ingénierie (02/09)

Question posée directement : *« est-ce qu'on peut faire ce travail, techniquement ? »* Réponse mécanisme par mécanisme, sans arrondir.

| Mécanisme | Faisable ? | Effort réaliste | Blocage |
|---|---|---|---|
| **Mesurer la calibration du score** | Oui, **immédiatement** | 2-3 jours | Aucun |
| **Résolution d'entité floue** (chaînes + structure) | Oui | 1-2 semaines | Aucun |
| **Rapprochement comportemental** | Oui | 2-3 semaines | Besoin de données OT+IT corrélées — disponibles en simu |
| **Détection de dérive** | Oui | 1-2 semaines | Aucun |
| **Priors inter-sites** | Techniquement oui | — | **Zéro deuxième site. Démarrage à froid.** |

**Tout sauf le dernier : 6 à 8 semaines pour un ingénieur.** Ce n'est pas un pari technologique — la résolution d'entité est un domaine mature (modèle Fellegi-Sunter, 1969 ; outils actuels type Splink/Zingg).

#### Le mécanisme qui vaut vraiment le coup — le comportement, pas le nom

Le rapprochement par les **noms** (Levenshtein, tokens, position hiérarchique) est facile mais plafonne : quand `M-001` et `Ligne2_Four` ne partagent aucun caractère, aucune similarité textuelle ne sauve.

Ce qui marche là où le nom échoue : **la corrélation comportementale.** L'OT fournit les transitions Run/Stop d'une machine, horodatées à la seconde. L'IT fournit l'avancement des quantités d'un OF sur un centre de charge. Si les périodes de production de la machine X corrèlent avec l'avancement du centre de charge Y, **c'est la même machine — et ça ne dépend d'aucun nom.** Techniquement : corrélation croisée entre deux flux d'événements.

**Pourquoi c'est structurellement défendable** : ça exige d'avoir les séries temporelles OT **et** les données transactionnelles IT dans le même système, **dans la durée**. Un outil qui ne fait du mapping qu'au moment de la configuration ne peut pas le faire, par construction. L'avantage vient de la position, pas de l'astuce.

*Deux détails qui changent la précision* : traiter le rapprochement comme un **problème d'affectation globale** (Hongrois) plutôt que des paires gloutonnes, et exploiter la cardinalité (une machine OT = au plus un centre de charge).

#### Protocole du test de calibration — à exécuter avant d'en parler

**Pourquoi celui-là d'abord** : si le score n'est pas calibré, tout l'argument O(n) → O(u) tombe, et l'auto-acceptation devient du bruit accepté silencieusement — pire que pas d'automatisation. C'est aussi le test le moins cher de tous.

- **Données** : serveur Prosys (hiérarchie connue) + ERP simulé (`fake_erp`, centres de charge connus). **On maîtrise la vérité terrain des deux côtés** — c'est ce qui rend le test possible sans client.
- **Procédure** : lancer la découverte → récupérer chaque mapping proposé avec son score → comparer à la vérité connue.
- **Mesures** :
  1. **Diagramme de fiabilité** — regrouper les propositions par tranche de score (0,5-0,6 / 0,6-0,7 …) et mesurer le taux d'exactitude réel de chaque tranche.
  2. **Erreur de calibration (ECE)** — écart moyen entre score annoncé et exactitude constatée.
  3. **Précision au seuil d'auto-acceptation** (`AutoAcceptThreshold`, 0,7 aujourd'hui) — quelle proportion de ce qui passe automatiquement est réellement juste.
  4. **Taux de revue humaine `u/n`** — quelle fraction descend sous le seuil.
- **Critères de réussite** :
  - Précision au seuil **≥ 0,95** — sinon la porte laisse passer des erreurs en silence, ce qui est le pire des cas.
  - Écart |score − exactitude| **≤ 0,10** par tranche.
  - **`u/n` ≤ 20 %** — au-delà, le mécanisme ne fait pas gagner grand-chose et l'argument de scalabilité est faible même s'il est honnête.
- **Si ça échoue** : ce n'est pas fatal, c'est informatif — soit les heuristiques sont à revoir, soit le seuil est mal placé. Mais **il faut le savoir avant le workshop**, pas après.

#### Les trois risques honnêtes

1. **La calibration n'est pas mesurée** — traité ci-dessus.
2. **L'approche LLM de HighByte n'est pas naïve.** Un LLM à qui on donne noms de tags et valeurs de colonnes ERP fait probablement du rapprochement flou correctement, sans cette machinerie. L'argument « déterministe et auditable » ne vaut que si la précision est **au moins comparable**. Si leur LLM est nettement meilleur, l'auditabilité ne suffira pas. **Non testé.**
3. **Les priors inter-sites ne sont pas démontrables avant 5 à 10 sites** — techniquement faisables, commercialement lointains.

#### Conclusion

> **Le risque n'est pas technique. On sait construire ça.** Ce qu'on ne sait pas, c'est si quelqu'un l'achète, et si on le fait mieux que celui qui vend déjà à 18 500 $/site — et **aucune de ces deux réponses ne s'obtient en codant.** D'où la priorité donnée à la calibration (3 jours), à la comparaison HighByte (essai gratuit) et à l'outreach, avant tout développement.

---

## Q7. Que manque-t-il aux agents IA ?

**Niveau 1** — Le tuyau existe déjà : MCP est devenu en 2026 la couche de connectivité de fait, y compris en industrie — HighByte en embarque un. **Ce qui manque n'est pas le transport, c'est le contenu.** Un agent branché sur des tags bruts peut *récupérer*, il ne peut pas *raisonner*.

**Niveau 2** — Cinq manques précis, chacun bloquant une classe de question :

| Ce qui manque | Sans ça, l'agent ne peut pas répondre à… |
|---|---|
| **Sémantique** — `4001:Val` = quoi, quelle unité, quel équipement | « la température du four 3 est-elle anormale ? » |
| **Résolution d'entité** — la même machine s'appelle `M1` dans l'OT, `Machine_01` dans le MES, `WC-001` dans l'ERP | « quel OF tournait quand cette machine s'est arrêtée ? » |
| **Topologie et relations causales** — quelle machine alimente quelle autre | « pourquoi la ligne 2 s'est-elle arrêtée ? » |
| **Alignement temporel** — événements OT à la seconde vs enregistrements IT saisis en fin de poste | « cet arrêt a-t-il causé ce rebut ? » |
| **Provenance et confiance** — d'où vient cette valeur, est-elle encore vraie | « peut-on décider sur cette base ? » |

**Le point à faire passer, et il est plus dur qu'il n'en a l'air** : un agent sans contexte n'est pas *inutile*, il est **dangereux** — il produit une réponse confiante et fausse sur une usine. C'est le même argument que Physical Intelligence formule pour les robots : l'IA physique doit se tromper beaucoup moins souvent que l'IA classique, parce qu'aucun humain ne médie la décision finale.

**Formulation** : *« MCP a réglé comment un agent se branche. Personne n'a réglé sur quoi il se branche. »*

---

---

## Q8. Que manque-t-il aux robots ?

Traité en détail dans `robotique_etat_art_workshop_2026-08-31.md`. En une ligne pour ce document :

**Pas de la donnée d'entraînement — du contexte opérationnel externe, injecté dans la boucle tâche/dispatch.** La donnée VLA (vision + langage + action + trajectoire) est d'une autre nature et hors périmètre. Ce qui manque à un robot *déjà déployé*, c'est l'état machine vif, l'OF actif, l'alerte qualité — et seule la boucle tâche/dispatch est ouverte pour les recevoir (les boucles sécurité et mouvement sont fermées par conception, et depuis ISO 10218:2025 toute donnée qui actionne un mouvement entre dans le dossier de sécurité **et** cyber).

**Rappel d'honnêteté** : personne dans le monde robotique n'a exprimé ce manque. C'est une déduction de compatibilité, pas une demande observée.

---

---

## Q9. « Ils ont déjà un data lake, et toute donnée est un actif »

C'est **l'objection la plus dangereuse** de la liste, parce qu'elle vient souvent de quelqu'un qui vient de dépenser plusieurs millions.

**Niveau 1** — Un data lake stocke ; il ne signifie pas. Il vous dit ce qui a été enregistré, pas ce que ça veut dire ni si c'est encore vrai. **Avoir un lac ne résout pas la question du sens — il la reporte**, et il la rend plus visible parce qu'on voit maintenant des téraoctets inutilisés.

**Niveau 2** —

**Le constat documenté** : les data lakes industriels dégénèrent régulièrement en **« data swamps »** — environnements ingérables et peu fiables où de la donnée OT brute s'accumule sans structure ni gouvernance, *« où des téraoctets de données capteur restent inutilisés parce que personne ne sait ce qu'elles signifient »*. La formulation la plus nette du problème, à citer telle quelle :

> ***« un tag automate nommé `4001:Val` n'apporte aucune information ; sans métadonnées (actif, équipe, produit), la donnée est du bruit. »***

Un lac qui ne stocke que du brut non transformé apporte peu de valeur relative : la donnée, extraite et stockée à grands frais, devient **inutilisable pour quiconque en dehors de l'équipe du projet lac elle-même.**

**Ce qu'un lac n'a structurellement pas** : la sémantique, la résolution d'entité entre systèmes, la topologie, la fraîcheur (un lac est batch, les décisions OT se jouent en secondes), et la provenance/confiance.

**L'argument économique, et c'est le plus fort** :

> Dans une architecture lac, **le coût de reconstruction du contexte est payé par chaque consommateur, à chaque fois** — chaque analyste, chaque tableau de bord, chaque modèle, chaque agent refait le travail de comprendre ce que les colonnes veulent dire.
> Une couche de contexte le paie **une fois**, en amont.

**Le bon positionnement, non frontal** : nous ne sommes **pas un concurrent du lac**, nous sommes **en amont**. On ne demande pas de le remplacer — on rend exploitable ce qui y atterrit. Un acheteur qui vient d'investir dans un lac n'a pas besoin d'entendre qu'il s'est trompé ; il a besoin d'entendre pourquoi son lac ne produit pas encore de valeur.

### « Toute donnée est un actif » — le recadrage

**Non. La donnée brute est un passif** : elle coûte à stocker, elle élargit la surface de conformité et de sécurité, et elle ne vaut rien tant qu'elle n'est pas interprétée. **L'actif, c'est la décision qu'elle permet.** C'est exactement la règle déjà validée deux fois en call (Geneviève, puis Cécilia) : le management ne paie pas pour « nettoyer la donnée », il paie pour un résultat chiffré.

### Pourquoi ça accélère la décision — la décomposition qui rend l'argument mesurable

Latence de décision = **t_détecter + t_comprendre + t_décider + t_agir**

- **t_détecter** est déjà court : le capteur remonte la valeur.
- **t_comprendre** — de quelle machine parle-t-on, quel produit tournait, est-ce anormal, qui est concerné — **est le terme dominant, et il est humain aujourd'hui** : quelqu'un ouvre trois systèmes et réconcilie à la main.
- **t_décider** et **t_agir** dépendent de l'organisation, pas de nous.

**La couche de contexte n'attaque qu'un terme — mais c'est celui qui domine.** Et c'est mesurable : le temps entre l'événement machine et la décision informée, avant et après. C'est précisément le type de KPI que Jalil exige, et il se mesure sur un pilote sans avoir à croire qui que ce soit sur parole.

**Pour l'entraînement de modèles** — le contexte n'est pas un accélérateur, c'est une condition : les étiquettes *sont* du contexte. Un flux de capteurs sans étiquette n'entraîne rien de supervisé.

---

---
---

# Ce qu'il faut valider

Ce que la recherche documentaire ne peut pas trancher. À poser en diagnostic, pas en pitch (cadrage Geneviève), sans annoncer la réponse attendue. Le plan d'outreach et les contacts sont dans `outreach_validation_2026-09-01.md` et `prospects_workshop_2026-09-01.xlsx`.

### Sur Q1 — le problème est-il là où on le dit ?

À poser en diagnostic, pas en pitch (cadrage Geneviève), sans annoncer la réponse attendue. *(Les questions de validation concurrentielle et build-vs-buy sont dans le document dédié.)*

- Quand vous avez dû connecter une nouvelle source de donnée terrain, qu'est-ce qui a réellement pris le plus de temps — brancher le protocole, ou se mettre d'accord sur ce que chaque tag voulait dire ?
- Combien de tags exposés sur votre parc, à peu près ? Qui sait ce qu'ils signifient — un document, une personne, personne ?
- Quand une machine est remplacée ou renommée, que devient le mapping existant ?
- Un projet a-t-il déjà été bloqué parce que l'IT et le terrain ne parlaient pas le même langage ?
- Où s'arrête la donnée aujourd'hui — reste-t-elle sur l'automate, ou remonte-t-elle jusqu'à l'ERP ?

---

### Sur Q2-Q4 — vérification concurrentielle et terrain

**Priorité 1 — vérification concurrentielle (faisable seul, quelques heures)**
1. Installer l'essai gratuit 30 jours de HighByte et **utiliser « AI Generate Instances » contre notre propre serveur Prosys.** C'est la comparaison la plus directe possible : même serveur, même arborescence, leur mécanisme contre le nôtre. Rien ne remplace ça.
2. HighByte fait-il de la **résolution d'entité OT↔IT** (équipement ↔ enregistrement ERP) ? Non vérifié.
3. Quelle est la complexité réelle d'installation de HighByte ? Notre argument « simplicité » en dépend entièrement.

**Priorité 2 — à poser aux techniciens d'usine**
- Connaissez-vous HighByte, Litmus, UMH ? Évalués ? Écartés pourquoi ?
- Un budget de l'ordre de 18 k$/site/an pour cette couche, c'est dans quelle catégorie chez vous — déjà provisionné, ou impensable ?
- *(la question la plus révélatrice)* — **Un outil interne important est-il devenu difficile à faire évoluer parce que la personne qui l'avait écrit est partie ?**
- Quand vous avez connecté une nouvelle source terrain, qu'est-ce qui a pris le plus de temps : brancher le protocole, ou se mettre d'accord sur ce que chaque tag signifiait ?
- Combien de cas d'usage espérez-vous servir avec cette donnée sur trois ans — un, ou plusieurs ?

---

### Sur Q5-Q9 — par ordre de coût croissant

1. **Le score est-il calibré ?** Mesurable seul, aujourd'hui, sur les données Prosys : quand le score dit 0,8, a-t-il raison 80 % du temps ? Sans cette calibration, tout le raisonnement O(n)→O(u) du §Q6.2 est une intention. **C'est le test le plus important et le moins cher.**
2. **Comparaison directe HighByte** — leur essai gratuit 30 j, « AI Generate Instances » contre notre serveur Prosys. Combien de nœuds chacun propose, combien sont justes, combien passent en revue humaine.
3. **Les priors inter-sites tiennent-ils ?** Non testable avant d'avoir deux sites réels. **À ne pas présenter comme acquis** — c'est le mécanisme le plus important du §Q6 et le moins prouvé.
4. **En call terrain** : « combien de temps entre le moment où une machine s'arrête et le moment où quelqu'un sait *pourquoi* ? » — mesure directe de t_comprendre.
5. **En call terrain** : « vous avez un data lake ? Qu'est-ce qui en sort réellement aujourd'hui ? »

---

---
---

# Garde-fous — à ne jamais dire

Chacun est une affirmation tentante, fausse, et qui coûterait la crédibilité devant un interlocuteur technique.

### Sur le problème (Q1)
1. ❌ **« OPC-UA ne fait que du transport »** — faux, et un bon technicien le sait. OPC-UA porte un modèle d'information et des Companion Specifications. Le point juste : **la couche sémantique existe dans la norme, pas dans le parc installé.**
2. ❌ **Citer « 80 % d'échec »** comme un chiffre unique agrégé — donner la fourchette réelle selon la source (69-85 %), et la source.
3. ❌ **Présenter la fragmentation comme un simple retard historique** — elle est en partie rationnelle économiquement pour les fournisseurs d'automatisation.

### Sur la concurrence (Q2-Q4)
4. ❌ **« Nous sommes les seuls à avoir MCP au bord »** — faux depuis juillet 2025 (HighByte 4.2). **À corriger dans `docs/mindset.md` L589 et L1262-1263.**
5. ❌ **« Personne ne fait ce qu'on fait »** / « l'auto-génération du modèle est unique » — faux.
6. ❌ **« La donnée ne sort pas de chez vous » comme différenciateur** — HighByte supporte les LLM locaux. L'argument se resserre sur la juridiction du fournisseur.
7. ❌ **« C'est trop difficile à construire vous-mêmes »** — faux pour la v1, et le dire fait décrocher. L'argument est le coût de possession.
8. ❌ **« Les consultants ne servent à rien »** — faux, et insultant pour un client qui en emploie.
9. ❌ **Revendiquer la couverture de connecteurs** — Litmus annonce 250+.

### Sur la scalabilité (Q5-Q9)
10. ❌ **« On supprime l'intégrateur »** — faux, et cela ferme la porte du meilleur canal de distribution. On réduit son dernier segment ; il fait plus de sites.
11. ❌ **« Chaque site est différent, donc personne ne peut le produitiser »** — c'est l'argument des intégrateurs, et il est faux : on ne standardise pas le contenu, on standardise le processus de convergence.
11bis. ❌ **« Ils passent par des intégrateurs, donc le dernier segment est irréductible »** — raccourci causal. Le canal intégrateur est aussi la façon normale de vendre du logiciel d'entreprise (Salesforce et SAP passent par Deloitte). Ce qui est établi, c'est que la couche sémantique est du travail humain facturé ; son irréductibilité reste à démontrer. Voir §0.B.
11ter. ❌ **Ignorer UMH quand on décrit l'effort d'implémentation** — ils annoncent 90 secondes pour connecter une machine et aucun recours aux intégrateurs. La réponse juste n'est pas de les écarter : c'est de distinguer temps de **connexion** et temps de **modélisation sémantique**, que personne ne chiffre.
12. ❌ **« Notre score de confiance résout le problème »** tant que sa calibration n'est pas mesurée. Un score non calibré est du bruit auto-accepté — pire que pas d'automatisation.
13. ❌ **« Les priors inter-sites nous donnent un effet réseau »** — mécanisme non testé, aucun deuxième site. À présenter comme une thèse, pas un acquis.
14. ❌ **« Votre data lake ne sert à rien »** — frontal, et souvent adressé à la personne qui l'a financé. Formuler en amont/complémentaire.
15. ❌ **« Toute donnée est un actif »** — c'est l'inverse : la donnée brute est un passif, l'actif est la décision.
15bis. ❌ **Réduire la convergence OT/IT au mapping des tags** — c'est la moitié OT du problème et sa couche la plus superficielle. La convergence, c'est la **jonction** entre entité OT et enregistrement IT (Q6, couche 3), et c'est là que la difficulté et la valeur se trouvent toutes les deux.
15ter. ❌ **Laisser croire que la résolution d'entité est résolue** — elle est aujourd'hui en correspondance exacte normalisée, donc inopérante quand les noms OT et ERP ne se ressemblent pas. C'est la vraie dernière couche de spécificité.
16. ❌ **Promettre l'automatisation totale de la modélisation** — la logique métier propre au site restera humaine.

---
---

# Sources

Vérifiées les 2026-08-31 et 2026-09-01.

- [Unified Namespace (UNS) Architecture: The Definitive 2026 Guide — Anexee](https://www.anexee.com/blog/unified-namespace-uns-architecture-industrial-2026) — modélisation hiérarchique = 40-60 % de l'effort
- [How UNS Prepares Manufacturing Data for AI — IIoT World](https://www.iiot-world.com/smart-manufacturing/unified-namespace-manufacturing-2026/)
- [Making OPC UA Models Easier: Companion Specs and Custom NodeSets — OPC Connect](https://opcconnect.opcfoundation.org/2025/12/making-opc-ua-models-easier/) — structures propriétaires, nommage incohérent
- [OPC UA Part 1: Overview and Concepts — OPC Foundation](https://reference.opcfoundation.org/specs/OPC-10000-1/4) — modèle d'information, message, communication, conformité
- [OPC UA Interoperability for Industrie 4.0 and IoT — OPC Foundation (PDF)](https://opcfoundation.org/wp-content/uploads/2026/01/OPC-UA-Interoperability-For-Industrie4-and-IoT-EN.pdf)

- [Top 12 industrial technology trends, Hannover Messe 2026 — IoT Analytics](https://iot-analytics.com/top-industrial-technology-trends/) — l'IA agentique a exposé les limites de la télémétrie brute
- [Can industrial DataOps steer projects away from disaster? — IoT Analytics](https://iot-analytics.com/industrial-dataops/) — marché naissant
- [HighByte Releases Industrial MCP Server for Agentic AI — communiqué HighByte](https://www.highbyte.com/news/press-releases/highbyte-releases-industrial-mcp-server-for-agentic-ai)
- [HighByte Intelligence Hub 4.2 Empowers Agentic and AI-Assisted DataOps at the Edge — DBTA](https://www.dbta.com/Editorial/News-Flashes/HighByte-Intelligence-Hub-42-Empowers-Agentic-and-AI-Assisted-DataOps-at-the-Edge-170387.aspx)
- [Industrial MCP Services — HighByte](https://www.highbyte.com/intelligence-hub/industrial-mcp-services)
- [Release notes version 4.2 — HighByte](https://www.highbyte.com/resources/release-notes/version-4-2) — « AI Generate Instances », validation avant enregistrement
- [Pricing and Licensing Options — HighByte](https://www.highbyte.com/pricing) — 17 500 $ → 18 500 $/site/an, MCP inclus, essai 30 j
- [Connections — HighByte](https://www.highbyte.com/intelligence-hub/connections) — browse OPC UA
- [Pricing and Offering — United Manufacturing Hub](https://www.umh.app/pricing)
- [This open-source bet is paying off as UMH takes on industrial giants — Tech.eu](https://tech.eu/2026/03/11/this-open-source-bet-is-paying-off-as-united-manufacturing-hub-takes-on-industrial-giants/)
- [All You Need To Know About Industrial DataOps — Litmus](https://litmus.io/blog/all-you-need-to-know-about-industrial-dataops) — 250+ connecteurs
- [Cognite, Leader IDC MarketScape Industrial DataOps 2026](https://www.cognite.com/en/company/newsroom/cognite-positioned-as-a-leader-in-idc-marketscape-for-worldwide-industrial-dataops-platforms)

- [Build vs Buy Industrial Software 2026 — Elisa IndustrIQ](https://www.elisaindustriq.com/resources/blog/build-vs-buy-industrial-software-weighing-costs-risks-and-benefits-in-2026) — maintenance 65-85 % du coût de possession
- [Build vs Buy: Smart Manufacturing Technology — Instrumental](https://instrumental.com/build-better-handbook/build-vs-buy-a-guide-to-implementing-smart-manufacturing-technology) — 13 M$ / 144 mois-ingénieur / 4 M$ par an
- [Build or Buy — Sustainment](https://www.sustainment.com/post/build-or-buy-how-supply-chain-teams-should-consider-software) — orphelinat, recrutement, 850 K$-1,6 M$
- [True Cost of Enterprise Integrations: SI vs iPaaS vs Managed — OneIO](https://www.oneio.cloud/blog/true-cost-of-enterprise-integrations) — 50 K$-500 K$+, maintenance 20-35 %, dépassement 40-60 %
- [AWS Industrial Data Fabric Launch Partners — HighByte](https://www.highbyte.com/blog/industrial-data-fabric-launches-with-an-ecosystem) — Deloitte, Infosys, Cyient, TensorIoT ; « the last mile between concept and reality »
- [Technology & Channel Partners — HighByte](https://www.highbyte.com/company/partners)
- [Siemens and HighByte Partner on Industrial Data Operations to Scale Industrial AI — Automation.com](https://www.automation.com/article/siemens-highbyte-partner-industrial-data-operations-scale-industrial-ai)
- [GFT × Litmus — Powering Industrial AI at Scale with Edge Data](https://litmus.io/partners/partner-gft)

- [Industrial Data Contextualization: Why raw OT Data often fails — CTS](https://www.group-cts.com/en/blog/ii/industrial-data-contextualization-manufacturing) — « data swamp », `4001:Val` n'apporte aucune information
- [Clean up your data lakes: Data contextualization in manufacturing — Cognite](https://www.cognite.com/en/resources/blog/clean-up-your-data-lakes-data-contextualization-in-the-manufacturing-industry) — donnée brute inutilisable hors de l'équipe du lac
- [Manufacturing Data Platform: Data Lake to Data Lakehouse — ifactoryapp](https://ifactoryapp.com/industries/manufacturing-plant/manufacturing-data-platform-data-lake)


**Canal intégrateur — pages consultées directement le 2026-09-01** *(et non via des extraits de moteur de recherche)*
- [Beyond ISA-95: How Unified Namespace Solves Manufacturing's Data Silo Problem — CSIA / Automation World](https://www.automationworld.com/communication/article/55355934/control-system-integrators-association-csia-beyond-isa-95-how-unified-namespace-solves-manufacturings-data-silo-problem) — **la source la plus solide** : l'association des intégrateurs revendique elle-même la normalisation et le mapping ISA-95 → hiérarchie réelle
- [Technology & Channel Partners — HighByte](https://www.highbyte.com/company/partners) — 3 catégories de partenaires ; intégrateurs = « design and deployment », conseil, intégration
- [AWS Industrial Data Fabric Launch Partners — HighByte](https://www.highbyte.com/blog/industrial-data-fabric-launches-with-an-ecosystem) — Deloitte, Infosys, Cyient, TensorIoT ; « These boots on the ground are the last mile between concept and reality »
- [Careers — HighByte](https://www.highbyte.com/company/careers) — au 01/09 : un seul poste ouvert, Account Executive ; aucun poste d'implémentation
- [Unified Namespace for Manufacturing Data — UMH](https://www.umh.app/solutions/unified-namespace) — **contre-exemple** : 90 s pour connecter une machine, pilote 4-6 semaines, aucune mention d'intégrateurs

**Sources internes**
- `docs/robotique_etat_art_workshop_2026-08-31.md` — Q8, boucles de latence, ISO 10218:2025
- `docs/tarik.md` — supply chain ; mécanisme IMDS (signal agrégé sans divulgation d'identité), §4 paysage concurrentiel
- `docs/insights_2026-08-21.md` — les 4 silos, OT/IT ms vs secondes, system of record vs ownership, « sell technology, not consulting »
- `docs/insights.xlsx` — statistiques d'échec (70 / 69 / 35 / 85 %, SAP 55-75 %, dépassement 215 %) avec sources nommées
- `docs/outreach_validation_2026-09-01.md` + `docs/prospects_workshop_2026-09-01.xlsx` — plan de validation terrain
