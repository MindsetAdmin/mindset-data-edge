# État de l'art robotique — analyse complète et brief workshop

**Document unique et faisant autorité sur le domaine robotique.** Consolide, au 2026-08-31, le brief workshop et l'analyse complète du domaine. Ne contient que du contenu à jour : les hypothèses invalidées et les cadrages périmés ont été retirés (historique dans `analysis_log.md`).

Action item de `docs/workshop.md` (Mohamed) : *« Analyser l'état de l'art en robotique et identifier les modèles technologiques implémentables dans le secteur des usines »*, dans le cadre de l'exploration du deuxième ICP (économie physique) ouvert par Jalil.

**Structure** :
- **Partie I — Brief workshop** (§0-§9) : ce qu'on présente. Deux niveaux de lecture — **Niveau 1** = ce qui se dit à voix haute, sans jargon ; **Niveau 2** = la preuve technique, à sortir si on est challengé.
- **Partie II — Analyse complète** (§10-§15) : la matière brute derrière le brief. Recherche marché, données VLA, extraction intégrale de la conférence Physical Intelligence, contacts, questions de discovery.

**Posture** : la Partie I prend position, elle n'expose pas des options à égalité.

**Cadrage** : rien dans la Partie I ne présuppose que l'architecture actuelle de Mindset Data est la réponse — la direction produit n'est pas arrêtée. Ce sont des faits sur l'état du domaine, et une recommandation sur quoi vérifier ensuite.

---
---

# PARTIE I — BRIEF WORKSHOP

## 0. Bottom line

1. **Le bon découpage n'est pas « AMR / cobots / humanoïdes », c'est la boucle de latence.** Le contexte externe ne peut entrer que dans les deux boucles lentes. Ça explique en une phrase pourquoi les humanoïdes sont durs — et ce n'est pas « parce qu'ils sont humanoïdes ».
2. **Un mur réglementaire vient de bouger (ISO 10218:2025, première révision depuis 2011).** Toute donnée qui *déclenche* un mouvement entre dans le dossier de sécurité — et depuis 2025, la cybersécurité en fait explicitement partie. Donnée qui *informe* = libre. Donnée qui *actionne* = certification.
3. **Correction honnête de notre propre hypothèse** : les fleet managers AMR de 2026 s'intègrent déjà aux MES/ERP. Notre hypothèse « le WMS ne porte pas de contexte » était trop large. Ce qui survit est plus étroit, mais plus défendable.
4. **MCP est devenu en 2026 la couche de connectivité de fait pour les agents** — y compris en robotique (50+ serveurs, bridges ROS2). La question « comment un agent récupère du contexte » a maintenant une réponse standard.
5. **Recommandation** : ne rien construire pour l'instant. Deux appels d'ingénieur suffisent à trancher, et la question à leur poser a changé depuis l'analyse du 24/08.

---

## 1. Le cadre qui manque : classer par boucle de latence, pas par forme de robot

**Niveau 1** — On parle des robots par leur apparence (mobile, bras, humanoïde). Ce n'est pas le bon critère pour savoir où de la donnée externe peut servir. Le bon critère, c'est la vitesse de la boucle de décision qu'on veut alimenter. Un robot n'est pas un système, c'est quatre boucles empilées qui tournent à des vitesses différentes — et on ne peut brancher quelque chose que sur les deux plus lentes.

**Niveau 2** — Les quatre boucles :

| Boucle | Ordre de grandeur | Ce qui y vit | Contexte externe possible ? |
|---|---|---|---|
| **Sécurité** | < 1 ms, temps réel dur | Arrêt d'urgence, limitation de vitesse, surveillance de zone | **Non.** Certifiée, figée, isolée par conception |
| **Mouvement / contrôle** | 1–100 ms | Asservissement, trajectoire, évitement | **Non.** Déterminisme requis, toute latence externe est un risque |
| **Tâche / dispatch** | 100 ms – secondes | Quelle mission, quelle destination, quelle priorité | **Oui.** C'est ici que ça se joue |
| **Planification** | secondes – minutes | Ordonnancement de flotte, allocation, réordonnancement | **Oui.** Le plus tolérant |

**Ce que ça reclasse immédiatement** : l'inférence VLA des humanoïdes tourne à 30–100 Hz — c'est la boucle *mouvement*. D'où la difficulté réelle : le problème n'est pas « les humanoïdes sont durs », c'est que **la couche où on voudrait injecter du contexte n'est pas celle qui est ouverte**. À l'inverse, un AMR reçoit ses missions dans la boucle tâche/dispatch — nativement ouverte, par conception, et déjà standardisée (§3).

**Deux faits techniques qui complètent le tableau** :
- **ROS2 + DDS tient réellement le temps réel industriel** — latence sub-10 ms démontrée à 50 Hz dans des cellules industrielles distribuées, jusqu'à **<150 µs** avec PREEMPT_RT + Fast-DDS. Le temps réel n'est donc pas l'obstacle : l'obstacle est que l'inférence VLA (30-100 Hz) occupe la boucle *mouvement*, fermée par conception.
- **Les cobots/bras** (FANUC, ABB, KUKA, Universal Robots) s'intègrent via **API REST/webhook** — pattern d'intégration standard, mais un connecteur reste à construire.

**Phrase à dire** : *« L'ordre de faisabilité technique — AMR, puis cobots, puis humanoïdes — est exactement l'inverse de l'ordre de visibilité médiatique. Et la raison n'est pas la forme du robot, c'est la boucle dans laquelle on essaie d'entrer. »*

---

## 2. Le mur réglementaire que personne ne mentionne — et il vient de bouger

**Niveau 1** — Il y a une frontière nette entre donner une information à un robot et déclencher son mouvement. La première est libre. La seconde fait entrer la donnée dans le dossier de sécurité de la machine, avec un coût de certification. Cette frontière a été **redessinée en 2025**, et la plupart des discussions sur « l'IA dans les usines » l'ignorent complètement.

**Niveau 2** — **ISO 10218-1:2025 et 10218-2:2025** — première révision majeure depuis 2011 :

- Les exigences de **sécurité fonctionnelle deviennent explicites** au lieu d'être implicites : fonctions de sécurité définies, Performance Levels imposés, validation documentée selon **ISO 13849-1**.
- **ISO/TS 15066** (robots collaboratifs) est désormais **intégrée dans ISO 10218-2:2025** — ce n'est plus une spécification technique séparée.
- Nouvelle **classification Classe I / Classe II** : les robots légers et lents relèvent d'exigences réduites.
- **La cybersécurité est désormais traitée comme une composante de la sécurité fonctionnelle.** ← le point le plus important pour nous.

**Pourquoi ce dernier point compte précisément** : dès qu'une donnée externe influence un comportement robot, elle n'est plus seulement une question d'intégration — elle touche à un dossier réglementaire qui, depuis 2025, couvre explicitement la surface cyber. Une plateforme qui pousse du contexte vers un robot doit pouvoir répondre à des questions d'intégrité et d'authentification de la donnée, pas seulement de format.

**La ligne à tenir, et elle est nette** :
- Contexte qui **informe** un dispatcher, un ordonnanceur, un opérateur, un tableau de bord → hors périmètre sécurité, déployable librement.
- Contexte qui **actionne** directement un mouvement → dossier de sécurité + désormais dossier cyber. Coût et délai d'un autre ordre.

**Phrase à dire** : *« Tout ce qui influence le mouvement d'un robot entre dans son dossier de sécurité. Depuis la révision 2025, ça inclut la cybersécurité. Donc la question n'est pas "est-ce qu'on peut envoyer du contexte à un robot" — c'est "à quelle boucle", et la réponse détermine si on parle d'un connecteur ou d'un projet de certification. »*

---

## 3. La robotique vit sa crise d'interopérabilité — maintenant, pas dans cinq ans

**Niveau 1** — Les flottes de robots de marques différentes ne savent pas travailler ensemble. Le secteur est en train de résoudre ça avec des standards, en ce moment même. C'est exactement la même histoire que la fragmentation IT/OT dans les usines, avec vingt ans de décalage.

**Niveau 2** — Trois efforts, complémentaires et non concurrents :

| Standard | Origine | Portée |
|---|---|---|
| **VDA5050** | VDA (automobile allemande) + VDMA (manutention/intralogistique) | Interface entre robot mobile et logiciel de gestion de flotte. Transporté sur **MQTT** |
| **MassRobotics AMR Interoperability Standard** | MassRobotics (US) | Partage de données de base entre robots. Explicitement non conflictuel avec VDA5050 — un robot peut porter les deux |
| **Open-RMF** | Open Source Robotics | Orchestration multi-flotte, couche au-dessus |

**Signal de maturité 2026** : l'interopérabilité VDA5050 est devenue **un argument de vente** en 2026 ; les capacités de base ont été déployées en 2025 et des acteurs annoncent la conformité complète en 2026 (ex. Seegrid). Ce n'est plus une spec théorique, c'est un critère d'achat.

**Le parallèle à faire** (angle de fond, pas analogie décorative) : VDA5050 existe pour la même raison que l'UNS/ISA-95 côté usine — parce que des systèmes hétérogènes devaient parler une langue commune sans que chacun réécrive un connecteur par partenaire. Le secteur robotique refait aujourd'hui, sur son périmètre, le chemin que le monde OT/IT a fait avant lui.

---

## 4. Correction honnête : les fleet managers parlent déjà aux MES

**C'est une correction de notre propre analyse du 24/08, à porter nous-mêmes au workshop plutôt qu'à se la faire opposer.**

**Niveau 1** — On avait supposé que les flottes AMR étaient aveugles au contexte usine, et que le WMS ne portait pas les signaux de production. C'est trop large. En 2026, les logiciels de flotte annoncent déjà de l'intégration MES/ERP.

**Niveau 2** — Ce que les fleet managers revendiquent publiquement en 2026 (KUKA, Omron, Kinexon, Fives, Zimmer, Ati Robotics) :
- connexion **MES / ERP / WMS**, propagation automatique des jobs vers la flotte en temps réel ;
- déclencheurs **« call-for-parts »** issus directement des processus de production ;
- intégration **bidirectionnelle** avec SAP, Oracle, Infor : les signaux de production entrent, les événements de mouvement matière ressortent ;
- communication avec des équipements d'infrastructure (portes, convoyeurs, machines de dépose).

**Ce qui survit à la correction — plus étroit, mais plus solide** :

L'intégration annoncée est **transactionnelle** : un ordre, une demande de pièce, un événement de mouvement. Ce n'est pas la même chose qu'un **état OT vif** — machine qui vient de tomber, dérive qualité détectée à la seconde, changement d'état à granularité fine. La distinction entre « le MES sait qu'un OF existe » et « la ligne sait qu'une machine vient de s'arrêter il y a 4 secondes » reste réelle. Mais c'est une distinction **beaucoup plus fine** que celle qu'on portait le 24/08, et elle ne se démontre pas depuis un bureau : il faut demander à quelqu'un qui exploite une flotte.

**Signal externe que la question est vive** : le programme d'IMTS 2026 comporte une session intitulée *« Architecting Interoperable AMR Ecosystems: Bridging the Gap Between Shop Floor Logistics and Machine Throughput »* — soit très exactement l'écart logistique-de-plancher ↔ débit-machine qu'on avait posé en hypothèse. *(Honnêteté : seul le titre de la session a pu être vérifié ; le contenu n'était pas récupérable — à ne pas citer comme une preuve de contenu.)*

**Phrase à dire** : *« On s'était trompés d'un cran. Les flottes parlent déjà aux MES, en transactionnel. Ce qui reste ouvert, c'est l'état machine vif — et je ne peux pas trancher ça depuis un bureau, ça se demande à un exploitant. »*

---

## 5. MCP est devenu l'interface de contexte standard — en 2026

**Niveau 1** — La question « comment un agent IA récupère du contexte externe » avait autant de réponses que d'éditeurs. Depuis cette année, elle a une réponse standard, et la robotique l'a adoptée aussi.

**Niveau 2** — MCP, ouvert par Anthropic fin 2024, est décrit en 2026 comme **la couche de connectivité de fait pour l'IA agentique**, adopté par OpenAI et Google DeepMind. Côté robotique, l'écosystème compte début 2026 **50+ serveurs MCP**, dont plusieurs **bridges ROS 2** (`ros2-mcp-server`, `wise-vision/ros2_mcp`, extensions RDE ROS 2) permettant à un modèle d'introspecter un système ROS 2 vivant — topics, nodes, services.

**Pourquoi c'est structurant** : ça déplace la question. Elle n'est plus « quel protocole propriétaire faut-il apprendre pour parler à un robot » mais « qu'est-ce qu'on a d'utile à exposer ». Le tuyau est en train de se standardiser tout seul, indépendamment de nous.

**Nuance à ne pas sauter** : les bridges ROS2/MCP existants servent surtout à *introspecter et piloter* un système ROS depuis un LLM — c'est-à-dire l'inverse du flux qui nous intéresserait (pousser du contexte usine *vers* le robot). Même protocole, direction opposée. Ne pas présenter ces 50+ serveurs comme une validation de notre cas d'usage : ils prouvent que le tuyau existe, pas que quelqu'un l'utilise dans notre sens.

---

## 6. Ce qui reste vrai sur les VLA et les humanoïdes

Résumé — le détail complet est en Partie II (§10, §11, §12, §13).

**Niveau 1** — Entraîner un robot et donner du contexte à un robot déjà déployé sont deux métiers différents. Le premier n'est pas le nôtre, et prétendre le contraire serait faux.

**Niveau 2** —
- La donnée d'entraînement VLA = observation visuelle + instruction en langage naturel + trajectoire d'action + label de succès, synchronisés. Échelle : ~20 000 heures de donnée réelle pour un modèle fondation (LingBot-VLA) ; 1,4 M d'épisodes sur 22 plateformes pour Open X-Embodiment. **Aucun rapport avec des tags OPC-UA ou des lignes SQL.**
- **Trouvaille de la conférence Physical Intelligence** (labo π0/π0.5/π0.7) : la plupart des modèles fondation robotiques n'ont **aucune mémoire** — uniquement l'observation courante. Le contexte vidéo brut est prohibitif (~500 000 tokens pour 10 s à 50 Hz sur 4 caméras). Leur solution : mémoire multi-échelles, dont un **résumé textuel compressé** injecté à côté de la vidéo, plus des métadonnées structurées.
- **Ce que ça prouve, précisément** : les modèles fondation robotiques sont architecturés pour accepter du contexte **structuré/textuel en entrée de première classe**. C'est un fait sur l'architecture de ces modèles — pas une validation qu'un fournisseur de contexte externe est attendu.
- **La distinction à ne jamais brouiller** : leur mémoire est l'historique **propre à la tâche du robot**, calculée depuis ses propres caméras. Un contexte usine externe est une source différente. Complémentaire, pas identique.
- **Calendrier** : déploiements contractuels réels à horizon **2028–2032**. Tâches réalistes aujourd'hui : manutention légère, inspection.

---

## 7. Recommandation

**Ne rien construire maintenant. Deux appels tranchent, et la question a changé.**

1. **Concentrer le fil robotique sur la boucle tâche/dispatch, explicitement.** C'est la seule ouverte sans dossier de sécurité (§1, §2). Le dire comme un choix d'ingénierie assumé, pas comme un repli.

2. **Cibler l'AMR en ligne de production, pas en entrepôt.** L'entrepôt pur est couvert : le WMS route déjà, et les fleet managers s'intègrent déjà aux MES en transactionnel (§4). Si un écart existe, il est sur l'état OT vif en contexte production.

3. **La question de discovery a changé** — celle du 24/08 est périmée par §4. Ne plus demander *« votre flotte a-t-elle de la visibilité sur le contexte usine ? »* (réponse 2026 : « oui, on est intégrés au MES »), mais :
   > *« Votre flotte reçoit des ordres du MES. Est-ce qu'elle reçoit aussi l'état machine en temps réel — et si une machine tombe maintenant, combien de temps avant que la flotte le sache et se réordonne ? »*

   C'est la seule formulation qui distingue le transactionnel du vif.

4. **Deux appels suffisent, ils sont déjà identifiés** (§14) : **Romain Desarzens** (Movu, ingénieur — faisabilité VDA5050 et ce que le fleet manager voit réellement) et **Khalil Mosrati** (exploitant — ce qui se passe vraiment quand une machine s'arrête). Rien à construire avant leurs réponses.

5. **Ne pas prioriser les humanoïdes** — marché réel mais horizon 2028-2032, et la boucle visée est fermée (§1). Veille, pas piste.

**Critère d'abandon** (explicite, posé maintenant) : si les deux appels indiquent que le fleet manager reçoit déjà l'état machine vif, ou que le réordonnancement sur événement machine n'a pas de valeur perçue — **la piste robotique tombe**. Pas de zone grise, pas de « on garde en parallèle » indéfini. Sans ce critère, « explorer deux ICP en parallèle » devient mécaniquement « porter cinq pistes indéfiniment ».

---

## 8. Ce qu'on ne sait toujours pas

À dire au workshop plutôt qu'à laisser passer pour acquis.

- **Aucun ingénieur robotique n'a été consulté.** Toute la compatibilité protocolaire est déduite de spécifications publiques, jamais testée en pratique.
- **Personne dans le monde robotique n'a exprimé ce manque.** L'opportunité est déduite d'une compatibilité (les robots acceptent du contexte structuré + du contexte existe), pas d'une demande observée. **Compatibilité ≠ demande.**
- **La tension de fond n'est pas résolue** : les tâches où un contexte riche compterait le plus (longues, multi-étapes, jugement) sont précisément celles de la boucle fermée (§1). Les plus faciles à atteindre (dispatch AMR) sont peut-être assez simples pour ne pas en avoir besoin. « Facile à construire » et « là où on apporte de la valeur » pourraient tirer dans des directions opposées.
- **Rentabilité non chiffrée.** Aucune donnée sur trajets évités, temps de dwell WIP, incidents qualité évités. Ne jamais dire « c'est rentable » sans pilote quantifié.

---

## 9. Garde-fous — à ne jamais dire

1. ❌ « On prépare la donnée d'entraînement pour l'IA robotique » — faux, la donnée VLA est d'une autre nature (§6, §11).
2. ❌ « Les 50+ serveurs MCP robotiques valident notre cas d'usage » — faux, ils vont dans la direction opposée (§5).
3. ❌ « Les flottes AMR sont aveugles au contexte usine » — trop large, corrigé en §4.
4. ❌ « C'est rentable » pour l'AMR — aucun pilote quantifié (§8).
5. ❌ Présenter la compatibilité technique comme une preuve de demande (§8).

---
---

# PARTIE II — ANALYSE COMPLÈTE

Matière brute derrière la Partie I : recherche marché, données VLA, extraction intégrale de la conférence Physical Intelligence, contacts et questions de discovery. Le contenu périmé ou invalidé a été retiré — ce document ne contient que ce qui est à jour au 31/08.

## 10. État du marché — c'est réel, pas juste du hype

*(Note de cadrage : section marché, donc plutôt du ressort de Cécilia côté workshop — conservée ici comme matière de fond, à ne pas porter comme argument business par Mohamed.)*

Le marché mondial de la robotique humanoïde est estimé à **~4,2 Md$ en 2026**, avec **8,7 Md$ de financement** levé sur le secteur jusqu'à juillet 2026. Ce n'est plus au stade prototype :

- **Schaeffler** (un des plus gros équipementiers automobiles au monde) a signé un accord de déploiement contractuel avec Humanoid (UK) : 1 000 à 2 000 robots humanoïdes à roues sur ses sites de production mondiaux d'ici 2032.
- **Hyundai Motor** va introduire des humanoïdes Boston Dynamics dans ses usines, à commencer par le site de Géorgie (USA) en 2028.
- **BMW** (Caroline du Sud) et **Japan Airlines** (tarmac de Tokyo Haneda) font déjà tourner des systèmes d'IA physique en opération réelle.

Acteurs clés : Tesla, Figure AI, Agility Robotics, Apptronik, Boston Dynamics, 1X Technologies, UBTECH, Unitree, Fourier Intelligence, Sanctuary AI.

Contexte structurel : pénurie de main d'œuvre manufacturière (2,1 M de postes non pourvus projetés aux USA d'ici 2030, gap de 425 K en 2026, Deloitte).

**Ce qui est réellement implémentable aujourd'hui (pas dans 5 ans)** : déplacement de bacs/totes en logistique/entrepôt, transfert léger de matériel entre postes, tournées d'inspection (le robot porte un capteur).

**Ce qui n'est PAS réaliste à court terme** : tâches à répétabilité sub-millimétrique, charge utile > ~10 kg, environnements certifiés ATEX/zones dangereuses, cadence de production automobile constante. À garder en tête pour ne pas survendre le calendrier.

---

## 11. Le vrai problème de donnée en robotique — et pourquoi ce n'est PAS ce que fait Mindset Data aujourd'hui

Les modèles qui pilotent ces robots sont des **VLA (Vision-Language-Action)** — l'architecture dominante en 2026, adoptée par tous les grands labs IA. Chaque exemple d'entraînement nécessite 4 éléments synchronisés : **observation visuelle + instruction en langage naturel + trajectoire d'action + label de succès/échec**.

Échelle de donnée nécessaire : un fine-tuning demande de quelques milliers à quelques centaines de milliers de démonstrations de haute qualité. Un modèle fondation, bien plus : **LingBot-VLA** a été entraîné sur **~20 000 heures** de données réelles issues de 9 configurations de robots à double bras. Le dataset **Open X-Embodiment** (utilisé pour RT-X) contient **1,4 million d'épisodes** issus de **22 plateformes robotiques** différentes (apprentissage cross-embodiment).

Contrainte additionnelle : l'inférence VLA doit tourner à **30-100 Hz** sur du matériel embarqué avec une latence bornée — un niveau de temps réel encore plus strict que le « temps réel en millisecondes » déjà identifié côté OT dans `insights_2026-08-21.md`.

**Constat important, à dire honnêtement** : ce type de donnée (démonstrations vision + action + trajectoire, collectées par téléopération ou simulation) n'a rien à voir avec ce que Mindset Data ingère aujourd'hui (valeurs de tags OPC-UA/MQTT, lignes SQL d'ERP). Prétendre que Mindset Data « prépare la donnée d'entraînement pour l'IA robotique » serait un survente — l'architecture actuelle ne touche à aucun moment de la donnée vision/multimodale de démonstration robotique.

---

## 12. Le problème de mémoire/contexte — conférence Physical Intelligence

Source directe : `docs/video-transcript.md` — conférence de **Physical Intelligence** (labo à l'origine des modèles π0/π0.5/π0.7, une référence du secteur, pas une source secondaire).

**Constat clé** : la plupart des modèles fondation robotiques SOTA n'ont **aucune mémoire** — ils n'opèrent que sur l'observation sensorielle courante. Pour des tâches longues (10-15 min, multi-étapes, ex. nettoyer une cuisine), la mémoire est indispensable pour suivre la progression. Mais passer naïvement du contexte vidéo brut est prohibitif : 10 secondes de vidéo à 50 Hz sur 4 caméras représentent déjà **~500 000 tokens**.

**Leur solution** : une mémoire à échelles de temps multiples — mémoire vidéo court-terme (~10 s) calculée efficacement, **et pour le moyen/long terme une mémoire représentée en texte** (résumé compressé de ce qui s'est passé), injectée dans le modèle aux côtés de la vidéo, avec en plus des métadonnées structurées (instruction de sous-tâche, qualité de la donnée, etc.).

**Pourquoi c'est directement pertinent** : ça prouve que les modèles fondation robotiques sont explicitement conçus pour accepter du **contexte structuré et textuel comme entrée de première classe, à côté de la vision** — pas un hack, un choix d'architecture central chez un labo de référence. C'est exactement la *forme* de ce que produit un moteur de contextualisation/KG.

**La distinction à ne pas brouiller** : la mémoire de Physical Intelligence est l'historique **propre à la tâche du robot** (qu'a-t-il déjà fait dans cette tâche), calculée en interne à partir de ses propres caméras. Un contexte opérationnel **externe** de l'usine (état machine, OF actif, alerte qualité) est une source différente, même si c'est le même *type* d'entrée (structuré/textuel). Complémentaire, pas identique.

**Autre point utile** : l'IA physique doit avoir un taux d'erreur bien plus faible que l'IA classique, parce qu'elle agit directement sur le monde physique sans qu'un humain médiatise la décision finale — ce qui renforce, avec une raison technique précise, pourquoi le « grounding »/contexte n'est pas cosmétique mais structurant pour la fiabilité. Leur conseil aux petites équipes (« partir d'un modèle pré-entraîné open source comme π0/π0.5 et fine-tuner, pas entraîner from scratch ») confirme indirectement que la question « qui entraîne le robot » a déjà une réponse standard dans l'écosystème.

---

## 13. Extraction complète — conférence Physical Intelligence

Depuis `docs/video-transcript.md`. Liste exhaustive, pas seulement les points repris en §12.

### Problèmes techniques centraux

1. **L'IA physique doit se tromper beaucoup moins souvent que l'IA classique.** Contrairement à un système de recommandation (où un humain décide en dernier ressort, donc une erreur est rattrapable), un robot agit directement sur le monde physique sans médiation — la barre de fiabilité est bien plus haute. *Repère cité* : Waymo a passé 250 000 courses autonomes hebdomadaires, preuve qu'une IA physique fiable et autonome est atteignable.

2. **Atteindre >90% de fiabilité sur une tâche de manipulation complexe (faire un espresso) est difficile.** Itérer manuellement (collecter, entraîner, évaluer, recommencer) plafonne vite — les humains se fatiguent avant d'atteindre une très haute fiabilité.
   - **Solution** : apprentissage par renforcement en boucle fermée où le système cherche lui-même où il a besoin de plus de donnée/supervision. Résultat : >90% de réussite sur l'espresso, 2× de débit en plus grâce au RL post-training.

3. **Le RL façon LLM (PPO/GRPO) ne se transpose pas tel quel — les essais physiques coûtent cher.** Un LLM peut faire des millions d'essais quasi gratuitement. Pour un robot, 1 million de trajectoires d'une tâche d'une minute représenterait **700 jours-robot**.
   - **Solution 1 — éviter les trajectoires sans issue** : si le robot part sur un mauvais chemin (ex. attrape deux boîtes collées au lieu d'une), un humain intervient en téléopération pour montrer la récupération, ou au minimum on arrête l'épisode tôt.
   - **Solution 2 — amortir le coût d'estimation sur plusieurs tâches** : au lieu de 10-50 tentatives par tâche individuelle (comme PPO/GRPO), entraîner une **fonction de valeur générale** sur des vidéos de tâches très diverses (plier une chemise, récupérer un objet dans un frigo…), qui estime la progression sur n'importe quelle tâche.

4. **La fiabilité doit tenir dans la durée, pas juste sur un essai.** Test : la politique de préparation de latte a tourné **13 heures d'affilée** pour valider une fiabilité soutenue.

5. **Un modèle entraîné sur une tâche doit généraliser à d'autres environnements réels.** Testé sur le workflow réel d'une chocolaterie (Dandelion Chocolate — construction/étiquetage/empilement de boîtes) et sur le pliage de vêtements jamais vus dans une maison jamais vue.

6. **Les robots restent plus lents que les humains, et se trompent encore.** Seulement quelques itérations d'amélioration ont été faites.
   - **Piste** : le RL post-training donne déjà 2× de gain de vitesse ; une release séparée (« RL token ») a montré des vitesses dépassant la téléopération humaine. Une partie du goulot vient de ce que les téléopérateurs humains sont eux-mêmes lents en générant la donnée.

7. **La plupart des modèles fondation robotiques n'ont aucune mémoire/contexte.** Détaillé en §12.

8. **Un modèle spécialisé par tâche n'est pas un vrai généraliste — et manque de généralisation compositionnelle** (combiner des concepts jamais vus ensemble à l'entraînement, ex. « avocat » + « chaise » façon DALL-E 2021).
   - **Solution** : entraîner un seul modèle à grande capacité sur de la donnée maximalement diverse et hétérogène (y compris de basse qualité, rollouts de RL, vidéo humaine, donnée web), **et le « prompter » avec du contexte structuré riche** : mémoire, instruction de tâche, instruction de sous-tâche, métadonnées (qualité, longueur d'épisode), et optionnellement une image de sous-objectif générée. Ce prompting détaillé est décrit comme **la clé** pour exploiter de la donnée aussi hétérogène.
   - **Résultat** : le généraliste unique (π0.7) égale ou dépasse les spécialistes fine-tunés (RL et SFT) sur les tâches testées.
   - **Généralisation compositionnelle démontrée** : (a) interaction avec un appareil quasi absent des données d'entraînement (un air fryer, 3 épisodes seulement) — ouvre, insère une patate douce, ferme, en autonomie ; (b) transfert d'une compétence (plier du linge) vers une **plateforme robotique totalement différente** n'ayant reçu aucune donnée de pliage.

9. **Diversité et prompting riche sont-ils nécessaires, ou du raffinement ?** Vérifié par ablation : retirer le sous-ensemble le plus diversifié fait chuter drastiquement la performance sur tâches non vues ; retirer 20% de donnée aléatoire (moins diversifiée) a un impact minime — **la diversité compte plus que le volume brut**. Sans prompting par métadonnées, ajouter de la donnée de basse qualité **dégrade** la performance ; avec ce prompting, la même donnée **améliore** la performance.

### Points de Q&A (stratégiques/méta)

10. **Y aura-t-il un « moment ChatGPT » pour la robotique ?** Probablement pas de la même forme — la distribution physique sera plus lente qu'un produit logiciel. Mais un niveau de capacité comparable est « vraiment à l'horizon dans les prochaines années ».
11. **Quand une petite équipe doit-elle passer d'un spécialiste à un généraliste ?** Dès le départ — partir d'un pré-entraîné open source (π0/π0.5) et fine-tuner. Exception : environnements très contraints (pas d'internet, GPU faible, ex. robot chirurgical).
12. **L'équivalent robotique de la donnée « à l'échelle d'internet » ?** L'expérience du robot lui-même (téléopération puis autonomie), complétée par de la vidéo humaine (YouTube) et des images web légendées — mais « aucun substitut à l'expérience du robot lui-même ».
13. **Les modèles généralistes seront-ils démocratisés en open source comme les LLM ?** Incertain — le coût de la donnée incarnée et du matériel pourrait rendre la robotique différente ; optimisme prudent, sans certitude.
14. **Sortie du modèle : commandes moteur brutes, ou position cible + contrôleur ?** Leurs modèles prédisent des positions cibles d'articulations (ou position 3D de la pince), un contrôleur PD atteint la cible — pas un goulot identifié à ce stade.
15. **Les robots ont-ils besoin « d'imagination » (prédire une image future) ?** π0.7 le fait et ça améliore mesurablement la performance (ex. pliage de chemise), mais le modèle reste très bon sans — utile, pas indispensable.
16. **Anecdote de généralisation émergente** : sur une tâche d'assemblage (insérer une goupille dans du papier), le robot a fait une erreur de placement puis a spontanément inversé main gauche/droite pour corriger — comportement jamais vu dans les données, signe d'une équivariance apprise implicitement.

---

## 14. Plan de reachouts — valider l'angle AMR

Objectif : valider la position avant d'investir. Deux pistes en parallèle qui valident des choses différentes.
- **Track A — vendeurs AMR** : teste l'angle B2B2B de Cécilia (aider les vendeurs de robots à prouver le ROI à leurs clients).
- **Track B — clients finaux avec flotte AMR déployée** : discovery classique.

**Rien envoyé.** Recherche + contacts identifiés seulement.

**Découverte utile** : la vue de recherche LinkedIn a un bouton direct **« Add profile(s) to lemlist »** — pas besoin d'un nouvel outil.

### Track A — Vendeurs AMR

*Attention : chercher « AMR » seul remonte des gens prénommés « Amr » (collision de mot-clé) ; utiliser « autonomous mobile robot » ou les noms de vendeurs.*

| Nom | Rôle | Localisation | Connexion |
|---|---|---|---|
| **Lucas Heraud** | Manager et Référent technique Déploiement chez **Movu Robotics** | Paris | 2nd — meilleur candidat vendeur : rôle déploiement = confronté aux problèmes d'intégration/ROI client |
| **Romain Desarzens** | Senior Robotics Software Engineer @ **Movu Robotics** — Control & Motion planning | Paris | 2nd, mutuels : Mohamed Ben Gammem, Oussama Darouez |
| **Lukasz Tomaszewski** | Robotics Support Engineer @ **Locus Robotics** | Birmingham, UK | 2nd, mutuel : Ferid Š. Sejdović |
| Benoit Boyeau | Senior Robotics Software Engineer, C++/Python/ROS | Paris | 2nd, mutuels : Pierre-Arthur O'Hara, Yassine Hermi +2 |
| Duc-Canh Nguyen | Senior Robotics Engineer — ML/Computer Vision/AMR/ROS2/SLAM | Greater Paris | 2nd, mutuel : Mazelin Hervé |

**Priorité** : Lucas Heraud (déploiement, voit les vrais points de friction client) puis **Romain Desarzens** (technique, pour valider la compatibilité MQTT/VDA5050 directement).

**Ajout (2026-08-27)** : **Arnaud Lubespere** — Chef de projet certifié PMP®/SAFe®, **Intégrateur Automatisation Logistique SAP chez Airbus**, Toulouse, 1K followers. Double intérêt : automatisation logistique concrète (proche AMR) **et** contexte Airbus — recoupe l'exemple RFQ/Airbus de `tarik.md`. Bon candidat pour tester les deux threads avec une seule personne.

### Track B — Clients finaux avec flotte AMR déployée

| Nom | Rôle | Localisation | Connexion |
|---|---|---|---|
| **Khalil Mosrati** | Responsable Robotisation & Automatisation Industrielle — Moyens Industriels & Lignes de Production, PMP/LSSGB | Paris | 2nd, mutuels : Mazelin Hervé, Clémence Savarit +2 — meilleur candidat, titre exact, 11K followers |
| Sami Aloui | Responsable Automatisme & Informatique Industriel | Lens | 2nd, mutuel : Sghaier Yosser |
| Bastien Charrier | Responsable projets automatisation chez SBPROCESS | Lyon | 2nd |
| Emmanuel Lebreton | Responsable Automatisme et Informatique Industrielle | Pays de la Loire | 2nd |

### Questions de discovery par track

**Track A — vendeurs AMR**
- Comment vos clients mesurent-ils le ROI de leur flotte aujourd'hui — uniquement le remplacement de main d'œuvre, ou qualité/throughput/sécurité sont-ils intégrés ?
- **« Votre flotte reçoit des ordres du MES. Est-ce qu'elle reçoit aussi l'état machine en temps réel — et si une machine tombe maintenant, combien de temps avant que la flotte le sache et se réordonne ? »**
- Où s'arrête votre intégration côté client aujourd'hui — WMS/MES ? Des clients demandent-ils plus de contexte que ce que vous fournissez nativement ?
- Un partenaire qui apporterait ce contexte en MQTT/VDA5050 directement compatible avec votre stack — ça résout un problème réel côté vous, ou côté client ?

**Track B — clients finaux**
- Comment votre flotte AMR décide-t-elle de ses priorités aujourd'hui ?
- **Reçoit-elle l'état machine en temps réel, ou seulement des ordres transactionnels du MES ?**
- Un incident récent où la flotte a pris une décision sous-optimale par manque de contexte usine ?
- Si la flotte savait en temps réel qu'une machine vient de tomber en panne ou qu'une urgence qualité est déclarée, qu'est-ce que ça changerait concrètement à son comportement ?

---

## 15. Check d'honnêteté — théorique ou vrai besoin ?

> Le constat du 31/08 le **renforce** plutôt qu'il ne l'affaiblit : §4 montre que l'écart supposé est encore plus étroit qu'on ne le pensait.

Question posée directement : ce que Mindset Data pourrait apporter à la robotique, c'est théorique ou un vrai besoin validé ?

**Aujourd'hui, c'est théorique — une hypothèse bien construite, pas un besoin validé.**

**Ce qui est confirmé (preuve réelle)** :
- L'interface technique existe — les modèles fondation robotiques sont explicitement conçus pour accepter du contexte structuré/textuel en entrée (§12). Fait sourcé.
- Un point de douleur business adjacent existe — les entreprises calculent mal le ROI robot, ratant 30-60% de la valeur réelle (piste de Cécilia).
- Les flottes AMR parlent MQTT/VDA5050, protocole déjà géré. Fait de faisabilité réel.

**Ce qui n'est PAS confirmé — le vrai trou** : personne dans le monde robotique n'a dit « on est bloqués par un manque de contexte usine externe ». L'opportunité a été déduite en reliant une capacité architecturale à un actif technique déjà possédé — c'est un argument de **compatibilité**, pas une preuve de **demande**. La règle du workshop (Jalil : pas de validation technique sans KPI business) s'applique ici aussi.

**Une vraie tension à ne pas lisser** : les tâches où le contexte riche compterait le plus — longues, multi-étapes, qui demandent du jugement — sont exactement la catégorie humanoïde/VLA la plus dure à intégrer (§1). Les tâches les plus faciles à implémenter — un AMR qui déplace un bac d'un point A à un point B — sont assez simples pour ne pas forcément avoir besoin d'un contexte usine riche ; un fleet manager + intégration WMS suffit peut-être déjà (et §4 montre que cette intégration existe déjà). Donc « AMR = plus facile à construire » et « contexte = où on apporte de la valeur » ne se renforcent pas forcément — seule une vraie conversation peut trancher.

**C'est exactement pourquoi §14 existe.** Tant que Romain Desarzens, Khalil Mosrati, ou quelqu'un de similaire n'a pas dit quelque chose proche de *« oui, notre flotte prend de mauvaises décisions parce qu'elle ne sait pas X »* — ça reste une hypothèse, pas une opportunité validée, et ça doit être présenté comme tel.

---
---

## Sources

**Vérifiées le 2026-08-31 (Partie I)**

- [ISO 10218-2:2025 — Robotics, Safety requirements, Part 2](https://www.iso.org/standard/73934.html)
- [ISO 10218-1:2025 — Robotics, Safety requirements, Part 1](https://www.iso.org/obp/ui/en/#!iso:std:73933:en)
- [Updated ISO 10218: Major Advancements in Industrial Robot Safety Standards — A3/Automate](https://www.automate.org/robotics/news/updated-iso-10218-major-advancements-in-industrial-robot-safety-standards-now-available)
- [ISO 10218-1:2025 — Designing Industrial Robots with Built-in Safety (WIDE Automation)](https://www.wideautomation.com/en/iso-10218%E2%80%9112025-designing-industrial-robots-with-built%E2%80%91in-safety/)
- [A guide to VDA5050 — OTTO by Rockwell Automation](https://ottomotors.com/blog/interoperability-standard-vda5050/)
- [VDA 5050, MassRobotics, Open-RMF: what is what — SYNAOS](https://www.synaos.com/en/blog/vda-5050-massrobotics-open-rmf)
- [What Is the MassRobotics AMR Interoperability Standard?](https://www.massrobotics.org/what-is-the-massrobotics-amr-interoperability-standard/)
- [Seegrid adopts VDA5050 standard](https://whserobotics.com/news/seegrid-adopts-vda5050-standard-to-boost-amr-interoperability-and-fleet-integration/)
- [AMR fleet management: AI as a driver of efficiency — KUKA](https://www.kuka.com/en-us/products/amr-autonomous-mobile-robotics/amr-fleet-management-software)
- [Multi-Vendor AMR Fleet Manager — Ati Robotics](https://www.atirobotics.ai/solutions/orchestration/fleet-manager/)
- [Centralized AMR/AGV Orchestration & VDA 5050 Support — Kinexon](https://kinexon.com/solutions/amr-agv-fleet-management)
- [IMTS 2026 Conference: Architecting Interoperable AMR Ecosystems](https://www.aerospacemanufacturinganddesign.com/news/imts-2026-conference-architecting-interoperable-amr-ecosystems-bridging-the-gap-between-shop-floor-logistics-and-machine-throughput/) *(titre seul vérifié, contenu non récupérable)*
- [MCP and Robotics: bridging AI agents and robot systems via ROS — ChatForest](https://chatforest.com/guides/mcp-robotics-ros-integration/)
- [wise-vision/ros2_mcp — MCP Server bridging AI agents into ROS 2](https://github.com/wise-vision/ros2_mcp)
- [kakimochi/ros2-mcp-server](https://github.com/kakimochi/ros2-mcp-server)
- [What's New in MCP in 2026](https://strategizeyourcareer.com/p/whats-new-in-mcp-in-2026)

**Vérifiées les 2026-08-24 / 26 (Partie II)**

- [Physical AI Is Sending Humanoid Robots to Real Factory Floors in 2026 — Memeburn](https://memeburn.com/physical-ai-is-sending-humanoid-robots-to-real-factory-floors-in-2026/)
- [State of Robotics 2026 Report — Robotics Center of Silicon Valley](https://www.roboticscenter.ai/state-of-robotics-2026)
- [Humanoid & Quadruped Robots for Manufacturing Plants 2026 — ifactoryapp](https://ifactoryapp.com/industries/manufacturing-plant/humanoid-quadruped-robots-manufacturing-plant-2026-guide)
- [VLA Models: Training Data Requirements Explained — Shaip](https://www.shaip.com/blog/vla-models-what-vision-language-action-models-need-from-training-data/)
- [Vision-Language-Action (VLA) Models 2026 — Internet Pros](https://internet-pros.com/blog/vision-language-action-models-robotics-2026/)
- [NVIDIA adds AMR fleet management tools to Isaac ROS — The Robot Report](https://www.therobotreport.com/nvidia-adds-amr-fleet-management-tools-to-isaac-ros/)
- [Open source fleet management tools for AMRs — NVIDIA Developer Blog](https://developer.nvidia.com/blog/open-source-fleet-management-tools-for-autonomous-mobile-robots/)
- [ROS2 Real-time Performance Optimization and Evaluation — Chinese Journal of Mechanical Engineering](https://cjme.springeropen.com/articles/10.1186/s10033-023-00976-5)
- [Probabilistic Latency Analysis of DDS in ROS 2 — arXiv](https://arxiv.org/pdf/2508.10413)
- Conférence Physical Intelligence : `docs/video-transcript.md` (transcript intégral conservé séparément).
