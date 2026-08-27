# Analyse robotique — état de l'art et modèles implémentables (2026-08-24)

Action item de `docs/workshop.md` (Mohamed) : "Analyser l'état de l'art en robotique et identifier les modèles technologiques implémentables dans le secteur des usines" — en préparation du workshop dans 15 jours, dans le cadre de l'exploration du deuxième ICP/proposition de valeur (économie physique) mentionné par Jalil.

Sources vérifiées par recherche web (2026-08-24), pas des estimations.

---

## 1. État du marché — c'est réel, pas juste du hype

Le marché mondial de la robotique humanoïde est estimé à **~4,2 Md$ en 2026**, avec **8,7 Md$ de financement** levé sur le secteur jusqu'à juillet 2026. Ce n'est plus au stade prototype :

- **Schaeffler** (un des plus gros équipementiers automobiles au monde) a signé un accord de déploiement contractuel avec Humanoid (UK) : 1 000 à 2 000 robots humanoïdes à roues sur ses sites de production mondiaux d'ici 2032.
- **Hyundai Motor** va introduire des humanoïdes Boston Dynamics dans ses usines, à commencer par le site de Géorgie (USA) en 2028.
- **BMW** (Caroline du Sud) et **Japan Airlines** (tarmac de Tokyo Haneda) font déjà tourner des systèmes d'IA physique en opération réelle.

Acteurs clés : Tesla, Figure AI, Agility Robotics, Apptronik, Boston Dynamics, 1X Technologies, UBTECH, Unitree, Fourier Intelligence, Sanctuary AI.

**Ce qui est réellement implémentable aujourd'hui (pas dans 5 ans)** : déplacement de bacs/totes en logistique/entrepôt, transfert léger de matériel entre postes, tournées d'inspection (le robot porte un capteur). **Ce qui n'est PAS réaliste à court terme** : tâches à répétabilité sub-millimétrique, charge utile > ~10 kg, environnements certifiés ATEX/zones dangereuses, cadence de production automobile constante. À garder en tête pour ne pas survendre le calendrier au workshop.

---

## 2. Le vrai problème de donnée en robotique — et pourquoi ce n'est PAS ce que fait Mindset Data aujourd'hui

Les modèles qui pilotent ces robots sont des **VLA (Vision-Language-Action)** — l'architecture dominante en 2026, adoptée par tous les grands labs IA. Chaque exemple d'entraînement nécessite 4 éléments synchronisés : **observation visuelle + instruction en langage naturel + trajectoire d'action + label de succès/échec**.

Échelle de donnée nécessaire : un fine-tuning demande de quelques milliers à quelques centaines de milliers de démonstrations de haute qualité. Un modèle fondation, bien plus : par exemple LingBot-VLA a été entraîné sur ~20 000 heures de données réelles issues de 9 configurations de robots à double bras. Le dataset Open X-Embodiment (utilisé pour RT-X) contient 1,4 million d'épisodes issus de 22 plateformes robotiques différentes (apprentissage cross-embodiment).

Contrainte additionnelle : l'inférence VLA doit tourner à **30-100 Hz** sur du matériel embarqué avec une latence bornée — un niveau de temps réel encore plus strict que le "temps réel en millisecondes" déjà identifié côté OT dans `insights_2026-08-21.md`.

**Constat important, à dire honnêtement au workshop** : ce type de donnée (démonstrations vision + action + trajectoire, collectées par téléopération ou simulation) n'a rien à voir avec ce que Mindset Data ingère aujourd'hui (valeurs de tags OPC-UA/MQTT, lignes SQL d'ERP). Prétendre que Mindset Data "prépare la donnée d'entraînement pour l'IA robotique" serait un survente — l'architecture actuelle ne touche à aucun moment de la donnée vision/multimodale de démonstration robotique.

---

## 2bis. Le problème de mémoire/contexte — source directe : conférence Physical Intelligence (2026-08-26)

Transcript apporté par l'utilisateur (`docs/video-transcript.md`) — conférence de **Physical Intelligence** (labo à l'origine des modèles π0/π0.5/π0.7, une référence du secteur, pas une source secondaire). Apporte une précision technique plus fine que la recherche web seule sur où Mindset Data pourrait réellement s'insérer.

**Constat clé de la conférence** : la plupart des modèles fondation robotiques SOTA n'ont **aucune mémoire** — ils n'opèrent que sur l'observation sensorielle courante. Pour des tâches longues (10-15 min, multi-étapes, ex. nettoyer une cuisine), la mémoire est indispensable pour suivre la progression. Mais passer naïvement du contexte vidéo brut est prohibitif : 10 secondes de vidéo à 50Hz sur 4 caméras représentent déjà ~500 000 tokens.

**Leur solution** : une mémoire à échelles de temps multiples — mémoire vidéo court-terme (~10s), calculée efficacement, **et pour le moyen/long terme, une mémoire représentée en texte** (résumé compressé de ce qui s'est passé), injectée dans le modèle aux côtés de la vidéo, avec en plus des métadonnées structurées (instruction de sous-tâche, qualité de la donnée, etc.).

**Pourquoi c'est directement pertinent, précisément formulé** : ça prouve que les modèles fondation robotiques sont explicitement conçus pour accepter du **contexte structuré et textuel comme entrée de première classe, à côté de la vision** — pas un hack, un choix d'architecture central chez un labo de référence. C'est exactement la *forme* de ce que produit un moteur de contextualisation/KG.

**La distinction à ne pas brouiller** : la mémoire de Physical Intelligence est l'historique **propre à la tâche du robot** (qu'a-t-il déjà fait dans cette tâche), calculée en interne à partir de ses propres caméras. Ce que Mindset Data pourrait apporter, c'est le contexte opérationnel **externe** de l'usine (état machine, OF actif, alerte qualité) — une source différente, même si c'est le même *type* d'entrée (structuré/textuel, pas un flux de capteurs brut). Complémentaire, pas identique — à formuler précisément, pas comme un raccourci.

**Autre point utile** : la conférence souligne que l'IA physique doit avoir un taux d'erreur bien plus faible que l'IA classique (recommandation, etc.), parce qu'elle agit directement sur le monde physique sans qu'un humain médiatise la décision finale — ce qui renforce, avec une raison technique précise, pourquoi le "grounding"/contexte n'est pas cosmétique mais structurant pour la fiabilité. Et leur conseil aux petites équipes ("partir d'un modèle pré-entraîné open source comme π0/π0.5 et fine-tuner, pas entraîner from scratch") confirme indirectement que la question "qui entraîne le robot" a déjà une réponse standard dans l'écosystème — ce qui renforce que le terrain de Mindset Data (contexte, pas entraînement) est le bon axe, pas un lot de consolation.

---

## 3. Où Mindset Data pourrait réellement jouer un rôle — pas l'entraînement, le contexte d'exécution

La question posée en réunion (`workshop.md`, ligne 47) est : "le rôle de Mindset Data pourrait consister à préparer les données pour rendre les futurs modèles d'IA robotique **opérationnels**." Il y a une distinction cruciale entre deux problèmes différents, et un seul des deux correspond à ce que le produit fait déjà :

- **Entraîner le modèle du robot** (générer les démonstrations VLA) — pas le métier de Mindset Data, aucune brique existante ne s'en approche.
- **Donner à un robot déjà entraîné/déployé le contexte opérationnel temps réel pour agir correctement dans une usine donnée** — c'est exactement le moteur déjà construit (contextualisation ISA-95, graphe de connaissance, `kg_active_production`, détection Run/Stop). Un robot qui doit "aller chercher la pièce X" a besoin de savoir : quel OF est actif, quel est l'état de la machine, y a-t-il un problème qualité signalé — la même donnée contextualisée que le produit fournit déjà à un dashboard humain ou à un agent MCP, juste consommée par un système de contrôle robotique au lieu d'un humain.

C'est le même principe déjà établi pour le fil supply chain (`tarik.md` §1bis) : le rôle de l'IA/robot est borné et l'infrastructure de contextualisation reste inchangée — seul le consommateur final change. Positionnement honnête et défendable : **Mindset Data ne construit pas le cerveau du robot, il pourrait fournir les yeux et les oreilles contextualisées sur l'état de l'usine** dans laquelle ce robot opère.

---

## 3bis. Modèles technologiques implémentables, classés par maturité ET par compatibilité avec l'architecture existante (2026-08-26)

Approfondissement demandé sur la partie strictement technique de l'action item — Jalil demandait explicitement "identifier les modèles technologiques **implémentables**", pas seulement l'état de l'art le plus visible (les humanoïdes). Vérifié via recherche web (2026-08-26) : les protocoles d'intégration réels changent radicalement le classement de maturité.

**Tier 1 — AMR (robots mobiles autonomes, ex. MiR, Locus Robotics) : le plus implémentable, et le plus proche de l'architecture Mindset Data déjà construite.**
La communication de flotte AMR utilise **VDA5050**, un standard ouvert du secteur, transmis sur **MQTT**. Mindset Data parle déjà MQTT nativement (`mqtt_subscribe`, `internal/mqtt.Publisher`) — l'intégration d'une flotte AMR n'est pas un nouveau connecteur à inventer, c'est une extension directe du connecteur existant. C'est de loin le chemin technique le plus court de toute cette analyse.

**Tier 2 — Cobots/bras robotiques (FANUC, ABB, KUKA, Universal Robots) : implémentable, nécessite un connecteur à construire mais sur un pattern standard.**
Intégration via **API REST et webhooks** — un pattern déjà bien maîtrisé côté architecture (même famille que `sql_query`), mais **aucun connecteur REST n'existe aujourd'hui dans le code** (voir la correction déjà faite sur l'email L2S — ne pas re-claim ce qui n'est pas construit). Risque technique faible, contrairement à SAP/ODP-RFC, mais c'est un vrai travail à chiffrer, pas un acquis.

**Tier 3 — Humanoïdes/VLA : le moins mûr, et le moins compatible avec l'architecture actuelle.**
Confirmé (recherche web) : **ROS2 + DDS peut tenir des contraintes temps réel industrielles réelles** — latence sub-10ms démontrée à 50Hz dans des cellules industrielles distribuées, jusqu'à <150µs avec PREEMPT_RT + Fast-DDS. Le problème n'est donc pas que le temps réel soit impossible techniquement — c'est que l'inférence VLA (30-100Hz, §2) est un régime de latence différent de ce que le KG/dashboard actuel de Mindset Data cible aujourd'hui (interrogation via API/WebSocket, pas un bus temps réel sub-milliseconde). Brancher du contexte Mindset Data sur une boucle de contrôle robot humanoïde demanderait un pont dédié basse latence, pas juste exposer le KG existant tel quel.

**Conséquence pratique pour le classement "implémentable"** : l'ordre par facilité technique réelle est AMR (déjà quasi-acquis) > cobots (connecteur REST à construire, risque faible) > humanoïdes (pont temps réel dédié à concevoir, risque le plus élevé) — l'inverse de l'ordre de visibilité médiatique, où les humanoïdes dominent alors qu'ils sont le pari le plus lointain des trois, techniquement comme commercialement (§1).

Sources (2026-08-26) :
- [NVIDIA adds AMR fleet management tools to Isaac ROS — The Robot Report](https://www.therobotreport.com/nvidia-adds-amr-fleet-management-tools-to-isaac-ros/)
- [Open source fleet management tools for AMRs — NVIDIA Developer Blog](https://developer.nvidia.com/blog/open-source-fleet-management-tools-for-autonomous-mobile-robots/)
- [ROS2 Real-time Performance Optimization and Evaluation — Chinese Journal of Mechanical Engineering](https://cjme.springeropen.com/articles/10.1186/s10033-023-00976-5)
- [Probabilistic Latency Analysis of DDS in ROS 2 — arXiv](https://arxiv.org/pdf/2508.10413)

## 4. Recommandation pour le workshop

- Ne pas positionner Mindset Data comme un acteur de la donnée d'entraînement robotique (VLA/démonstrations) — ce serait un claim non défendable techniquement.
- Le positionnement défendable est celui du "grounding" temps réel : contexte opérationnel structuré pour un robot déjà opérationnel, en réutilisant le moteur existant, pas un nouveau produit.
- Le calendrier réaliste du marché (déploiements Schaeffler/Hyundai à horizon 2028-2032, tâches limitées à la manutention légère/inspection aujourd'hui) suggère que ce deuxième ICP est un pari sur 2-3 ans, pas une opportunité de vente immédiate — cohérent avec la décision du workshop de le garder en exploration parallèle plutôt que d'investir dessus à fond maintenant.
- Si le workshop veut un point d'entrée technique concret plutôt que le pari humanoïde à 2-3 ans, **les AMR (§3bis) sont le vrai point de départ implémentable** — protocole (MQTT/VDA5050) déjà parlé nativement par Mindset Data, contrairement aux humanoïdes qui dominent la couverture médiatique mais restent le pari technique le plus lointain des trois catégories analysées.

---

## 5. Extraction complète — problèmes et solutions de la conférence Physical Intelligence (2026-08-26)

Depuis `docs/video-transcript.md` (conférence Physical Intelligence, labo à l'origine de π0/π0.5/π0.7). Liste exhaustive, pas seulement les points repris en §2bis.

### Problèmes techniques centraux de la conférence

1. **L'IA physique doit se tromper beaucoup moins souvent que l'IA classique.** Contrairement à un système de recommandation (où un humain décide en dernier ressort, donc une erreur du modèle est rattrapable), un robot agit directement sur le monde physique sans médiation humaine — la barre de fiabilité est donc bien plus haute. *Repère cité* : Waymo a passé 250 000 courses autonomes hebdomadaires, preuve qu'une IA physique fiable et autonome est atteignable.

2. **Atteindre une fiabilité élevée (>90%) sur une tâche de manipulation complexe (faire un espresso) est difficile.** Itérer manuellement (collecter, entraîner, évaluer, recommencer) fonctionne mais plafonne vite — les humains se fatiguent avant d'atteindre une très haute fiabilité.
   - **Solution** : apprentissage par renforcement en boucle fermée où le système cherche lui-même où il a besoin de plus de donnée/supervision, plutôt qu'un humain qui itère à la main. Résultat : >90% de réussite sur l'espresso, 2x de débit en plus grâce au RL post-training.

3. **Le RL façon LLM (PPO/GRPO) ne se transpose pas tel quel à la robotique — les essais physiques coûtent cher.** Un LLM peut faire des millions d'essais quasi gratuitement (juste du calcul). Pour un robot, 1 million de trajectoires d'une tâche d'une minute représenterait 700 jours-robot — du matériel réel, du temps réel.
   - **Solution 1 — éviter les trajectoires sans issue** : si le robot part sur un mauvais chemin (ex. attrape deux boîtes collées au lieu d'une), un humain intervient en téléopération pour montrer la récupération, ou au minimum on arrête l'épisode tôt — pour ne pas gaspiller du temps robot sur de la donnée inutile.
   - **Solution 2 — amortir le coût d'estimation sur plusieurs tâches, pas une seule** : au lieu de faire 10-50 tentatives par tâche individuelle pour estimer ce qui est une bonne/mauvaise tentative (comme le fait PPO/GRPO), entraîner une **fonction de valeur générale** sur des vidéos de tâches très diverses (plier une chemise, récupérer un objet dans un frigo...), qui estime la progression sur n'importe quelle tâche — réduit le nombre de tentatives nécessaires par tâche spécifique.

4. **La fiabilité doit tenir dans la durée, pas juste sur un essai.** Test : la politique de préparation de latte a tourné 13 heures d'affilée pour valider une fiabilité soutenue, pas un coup de chance ponctuel.

5. **Un modèle entraîné sur une tâche doit généraliser à d'autres environnements/workflows réels, pas rester un tour de passe-passe.** Testé sur le workflow réel d'une chocolaterie (Dandelion Chocolate — construction/étiquetage/empilement de boîtes en carton) et sur le pliage de vêtements jamais vus dans une maison jamais vue.

6. **Les robots restent aujourd'hui plus lents que les humains, et se trompent encore.** Seulement quelques itérations d'amélioration ont été faites — plus d'itérations amélioreraient probablement encore la fiabilité et la vitesse.
   - **Piste de solution** : le RL post-training donne déjà un gain de vitesse mesurable (2x) en plus de la fiabilité ; une release séparée ("RL token") a montré des vitesses dépassant la téléopération humaine. Une partie du goulot d'étranglement vient du fait que les téléopérateurs humains sont eux-mêmes lents en générant la donnée d'entraînement.

7. **La plupart des modèles fondation robotiques n'ont aucune mémoire/contexte.** Ils n'opèrent que sur l'observation sensorielle courante — donc incapables de suivre la progression d'une tâche longue et multi-étapes. Passer du contexte vidéo brut est prohibitif (~500 000 tokens pour 10 secondes de vidéo à 50Hz sur 4 caméras).
   - **Solution** (déjà détaillée en §2bis) : mémoire à échelles de temps multiples — vidéo court-terme calculée efficacement, résumé textuel compressé pour le moyen/long terme. Permet une tâche autonome de 10-15 minutes, multi-étapes, non répétitive (nettoyer une cuisine).

8. **Un modèle spécialisé par tâche (fine-tuné) n'est pas un vrai robot généraliste — et manque de généralisation compositionnelle** (capacité à combiner des concepts/compétences jamais vus ensemble à l'entraînement, ex. "avocat" + "chaise" façon DALL-E 2021).
   - **Solution** : entraîner un seul modèle à grande capacité sur de la donnée maximalement diverse et hétérogène (y compris de la donnée de basse qualité, des rollouts de RL, de la vidéo humaine, de la donnée web), **et le "prompter" avec du contexte structuré riche** : mémoire, instruction de tâche, instruction de sous-tâche, métadonnées (qualité de la donnée, longueur d'épisode), et optionnellement une image de sous-objectif générée ("à quoi ça devrait ressembler dans quelques secondes"). Ce prompting détaillé est décrit comme **la clé** pour exploiter de la donnée aussi hétérogène.
   - **Résultat** : le modèle généraliste unique (π0.7) égale ou dépasse les modèles spécialistes fine-tunés (RL et SFT) sur les tâches testées — objectif "out of the box" atteint.
   - **Généralisation compositionnelle démontrée** : (a) interaction avec un appareil quasi absent des données d'entraînement (un air fryer, seulement 3 épisodes) — ouvre, insère une patate douce, ferme, en autonomie ; (b) transfert d'une compétence (plier du linge) vers une **plateforme robotique totalement différente** n'ayant reçu aucune donnée de pliage — a quand même réussi à plier.

9. **Est-ce que la donnée diverse et le prompting riche sont vraiment nécessaires, ou juste un raffinement ?** Vérifié par ablation : retirer le sous-ensemble le plus diversifié des données fait chuter drastiquement la performance sur des tâches non vues ; retirer 20% de donnée aléatoire (moins diversifiée) a un impact minime — la diversité compte plus que le volume brut. Sans prompting par métadonnées, ajouter de la donnée de basse qualité **dégrade** la performance ; avec ce prompting, la même donnée de basse qualité **améliore** la performance — le contexte riche est ce qui permet d'exploiter de la donnée autrement bruitée.

### Points soulevés en Q&A (stratégiques/méta, pertinence variable pour Mindset Data)

10. **Y aura-t-il un "moment ChatGPT" pour la robotique ?** Probablement pas de la même forme — la distribution physique (il faut un robot réel) sera plus lente qu'un produit logiciel. Mais le niveau de capacité comparable à ChatGPT est "vraiment à l'horizon dans les prochaines années."

11. **Quand une petite équipe doit-elle passer d'un modèle spécialiste à un modèle généraliste ?** Dès le départ — partir d'un modèle pré-entraîné open source (π0/π0.5) et fine-tuner, plutôt que reconstruire from scratch. Exception : environnements très contraints (pas d'internet, GPU faible, ex. robot chirurgical).

12. **L'équivalent robotique de la donnée "à l'échelle d'internet" ?** L'expérience du robot lui-même (téléopération, puis de plus en plus d'expérience autonome), complétée par de la vidéo humaine (YouTube) et des images web légendées — mais "aucun substitut à l'expérience du robot lui-même" (regarder Federer jouer ne rend pas meilleur au tennis).

13. **Les modèles robotiques généralistes seront-ils démocratisés en open source comme les LLM, ou le coût de la donnée/matériel concentrera-t-il les meilleurs modèles dans quelques labos ?** Incertain — le coût réel de la donnée incarnée et du matériel pourrait rendre la robotique différente des LLM sur ce point, mais de gros datasets et modèles puissants ont déjà été open-sourcés ; optimisme prudent sur une vraie communauté open source, sans certitude que ça se passera comme pour les LLM.

14. **Sortie du modèle : commandes moteur brutes, ou position cible + contrôleur ?** Leurs modèles prédisent des positions cibles d'articulations (ou position 3D de la pince), un contrôleur PD se charge d'atteindre la cible — fonctionne bien, n'est pas un goulot d'étranglement identifié à ce stade.

15. **Les robots ont-ils besoin "d'imagination" (prédire une image future) pour être utiles ?** Le modèle π0.7 le fait et ça améliore mesurablement la performance (ex. pliage de chemise), mais le modèle reste très bon même sans — utile, pas strictement indispensable à ce stade.

16. **Anecdote de généralisation émergente** : sur une tâche d'assemblage (insérer une goupille dans du papier), le robot a fait une erreur de placement puis a spontanément inversé main gauche/droite pour corriger — un comportement jamais vu dans les données d'entraînement ni de pré-entraînement, signe d'une équivariance gauche/droite apprise implicitement.

---

## 6. Notre réponse — "notre solution est-elle utile pour la robotique ?" (2026-08-26)

Synthèse de la réponse donnée dans la conversation, à conserver comme position officielle pour le workshop :

**Oui — pour une partie précise de la stack robotique, pas pour l'ensemble.**

- **Utile pour** : donner au robot le contexte opérationnel temps réel nécessaire pour agir correctement dans une usine donnée — quel OF est actif, état machine, alertes qualité. C'est exactement ce que fait déjà le KG + la contextualisation ISA-95 + le score de confiance pour un dashboard ou un agent MCP ; le robot n'est qu'un consommateur différent de la même sortie. La conférence Physical Intelligence (§2bis) confirme que les modèles fondation robotiques sont explicitement conçus pour accepter ce type de contexte structuré/textuel en entrée.
- **Pas utile pour** : entraîner le robot lui-même. La donnée d'entraînement VLA (vision + instruction + trajectoire + label de succès, collectée par téléopération/simulation) est un type de donnée complètement différent — rien dans l'architecture actuelle n'y touche, et prétendre le contraire serait le survente qu'on évite depuis le début de cette analyse.
- **Quelle catégorie de robot, concrètement** : les AMR sont le vrai point d'entrée court terme — protocole MQTT/VDA5050 déjà parlé nativement par Mindset Data (§3bis), quasi une extension plutôt qu'une nouvelle brique. Les cobots/bras nécessitent un connecteur REST (à construire, mais risque faible, pattern standard). Les humanoïdes restent les plus difficiles : techniquement faisable (ROS2/DDS tient le temps réel), mais l'inférence VLA (30-100Hz) est un régime de latence différent de l'architecture KG/dashboard actuelle — nécessiterait un pont basse latence dédié.
- **Calendrier** : marché réel, mais horizon 2-3 ans+ pour les humanoïdes spécifiquement (déploiement Schaeffler jusqu'en 2032, Hyundai à partir de 2028) — AMR/cobots sont matures **aujourd'hui**, juste moins visibles médiatiquement que les humanoïdes.

**Conclusion** : oui, utile de façon réelle et défendable — pour la couche de contexte, sur les catégories de robots matures (AMR en premier), pas comme un pari sur la donnée d'entraînement, et pas prioritairement comme un pitch humanoïde court terme.

Sources :
- [Physical AI Is Sending Humanoid Robots to Real Factory Floors in 2026 — Memeburn](https://memeburn.com/physical-ai-is-sending-humanoid-robots-to-real-factory-floors-in-2026/)
- [State of Robotics 2026 Report — Robotics Center of Silicon Valley](https://www.roboticscenter.ai/state-of-robotics-2026)
- [Humanoid & Quadruped Robots for Manufacturing Plants 2026 — ifactoryapp](https://ifactoryapp.com/industries/manufacturing-plant/humanoid-quadruped-robots-manufacturing-plant-2026-guide)
- [VLA Models: Training Data Requirements Explained — Shaip](https://www.shaip.com/blog/vla-models-what-vision-language-action-models-need-from-training-data/)
- [Vision-Language-Action (VLA) Models 2026 — Internet Pros](https://internet-pros.com/blog/vision-language-action-models-robotics-2026/)

---

## 7. Plan de reachouts — valider l'angle AMR (2026-08-26)

Objectif : valider la position du §6 (AMR = point d'entrée le plus implémentable) avant d'investir dessus. Deux pistes en parallèle, elles valident des choses différentes :
- **Track A — vendeurs AMR** : teste l'angle B2B2B de Cécilia (aider les vendeurs de robots à prouver le ROI à leurs clients via la couche contexte).
- **Track B — clients finaux avec une flotte AMR déjà déployée** : discovery classique, teste si un contexte usine plus riche améliorerait réellement leurs décisions de flotte.

**Rien envoyé.** Recherche + contacts identifiés seulement.

**Découverte utile** : la vue de recherche LinkedIn a un bouton direct **"Add profile(s) to lemlist"** — ces profils peuvent être poussés directement dans lemlist (déjà utilisé pour la campagne outreach existante), pas besoin d'un nouvel outil.

### Track A — Vendeurs AMR (angle partenariat/canal)

Trouvés via recherche LinkedIn ciblée (2026-08-26) — attention, chercher "AMR" seul remonte des gens prénommés "Amr" (collision de mot-clé), utiliser "autonomous mobile robot" ou les noms de vendeurs directement.

| Nom | Rôle | Localisation | Connexion |
|---|---|---|---|
| **Lucas Heraud** | Manager et Référent technique Déploiement chez **Movu Robotics** — Robotique, Gestion de projet, Commissioning | Paris | 2nd — meilleur candidat côté vendeur : rôle déploiement = confronté directement aux problèmes d'intégration/ROI client |
| **Romain Desarzens** | Senior Robotics Software Engineer @ **Movu Robotics** — Control & Motion planning | Paris | 2nd, mutuels : Mohamed Ben Gammem, Oussama Darouez |
| **Lukasz Tomaszewski** | Robotics Support Engineer at **Locus Robotics** | Birmingham, UK | 2nd, mutuel : Ferid Š. Sejdović |
| Benoit Boyeau | Senior Robotics Software Engineer, C++/Python/ROS | Paris | 2nd, mutuels : Pierre-Arthur O'Hara, Yassine Hermi +2 |
| Duc-Canh Nguyen | Senior Robotics Engineer — ML/Computer Vision/AMR/ROS2/SLAM | Greater Paris | 2nd, mutuel : Mazelin Hervé |

**Priorité** : Lucas Heraud (déploiement, donc voit les vrais points de friction client) puis Romain Desarzens (technique, pour valider la compatibilité MQTT/VDA5050 directement avec quelqu'un qui construit ces systèmes).

### Track B — Clients finaux avec flotte AMR déployée

| Nom | Rôle | Localisation | Connexion |
|---|---|---|---|
| **Khalil Mosrati** | Responsable Robotisation & Automatisation Industrielle — Moyens Industriels & Lignes de Production, PMP/LSSGB | Paris | 2nd, mutuels : Mazelin Hervé, Clémence Savarit +2 — meilleur candidat, titre exact, 11K followers (profil actif) |
| Sami Aloui | Responsable Automatisme & Informatique Industriel | Lens | 2nd, mutuel : Sghaier Yosser |
| Bastien Charrier | Responsable projets automatisation chez SBPROCESS | Lyon | 2nd |
| Emmanuel Lebreton | Responsable Automatisme et Informatique Industrielle | Pays de la Loire | 2nd |

### Questions de discovery adaptées par track

**Track A — vendeurs AMR**
- Comment vos clients mesurent-ils le ROI de leur flotte aujourd'hui — uniquement le remplacement de main d'œuvre, ou est-ce que qualité/throughput/sécurité sont intégrés ?
- Est-ce que votre flotte a de la visibilité sur le contexte usine au-delà de ce que vos propres capteurs voient (état machine, OF en cours, alerte qualité) ?
- Où s'arrête votre intégration côté client aujourd'hui — WMS/MES ? Y a-t-il des clients qui demandent plus de contexte que ce que vous fournissez nativement ?
- Un partenaire qui apporterait ce contexte en MQTT/VDA5050 directement compatible avec votre stack — est-ce que ça résout un problème réel côté vous, ou côté client ?

**Track B — clients finaux**
- Comment votre flotte AMR décide-t-elle de ses priorités aujourd'hui ?
- A-t-elle une visibilité sur l'état machine ou l'OF en cours, ou opère-t-elle indépendamment de ce contexte ?
- Un incident récent où la flotte a pris une décision sous-optimale par manque de contexte usine ?
- Si la flotte savait en temps réel qu'une machine vient de tomber en panne ou qu'une urgence qualité est déclarée, qu'est-ce que ça changerait concrètement à son comportement ?

### Prochaine étape

Ajouter ces profils dans lemlist (bouton déjà disponible depuis la vue de recherche), préparer un message court par track (pas un pitch produit, une question de discovery — même posture que `docs/discovery_questions_cto_2026-08-25.md` et `docs/call_oss_venture.md`).

---

## 8. Théorique ou vrai besoin du monde physique ? — check d'honnêteté (2026-08-26)

Question posée directement : ce que Mindset Data pourrait apporter à la robotique, c'est théorique ou un vrai besoin validé ? Réponse honnête, à ne pas laisser se diluer dans l'élan de l'analyse précédente.

**Aujourd'hui, c'est théorique — une hypothèse bien construite, pas un besoin validé.** À distinguer précisément :

**Ce qui est confirmé (preuve réelle)** :
- L'interface technique existe — les modèles fondation robotiques sont explicitement conçus pour accepter du contexte structuré/textuel en entrée (architecture Physical Intelligence elle-même, §2bis). Fait sourcé.
- Un vrai point de douleur business adjacent existe — les entreprises calculent mal le ROI robot, ratant 30-60% de la valeur réelle (piste de Cécilia, vérifiée). Réel aussi.
- Les flottes AMR parlent MQTT/VDA5050, un protocole que Mindset Data gère déjà. Fait de faisabilité d'intégration réel.

**Ce qui n'est PAS confirmé — c'est le vrai trou** : personne dans le monde physique/robotique n'a dit "on est bloqués par un manque de contexte usine externe." L'opportunité a été déduite en reliant une capacité architecturale (les robots peuvent accepter du contexte) à un actif technique déjà possédé par Mindset Data (le contexte) — c'est un argument de **compatibilité**, pas une preuve de **demande**. Deux choses différentes, et la règle du workshop (Jalil : pas de validation technique sans KPI business) s'applique ici aussi.

**Une vraie tension à ne pas lisser** : les tâches où le contexte riche compterait le plus — tâches longues, multi-étapes, qui demandent du jugement — sont exactement la catégorie humanoïde/VLA la plus dure à intégrer aujourd'hui (§3bis, écart de latence). Les tâches les plus faciles à implémenter maintenant — un AMR qui déplace un bac d'un point A à un point B — sont assez simples pour ne pas forcément avoir besoin d'un contexte usine riche ; un fleet manager + intégration WMS suffit peut-être déjà. Donc "AMR = plus facile à construire" et "contexte = où Mindset Data apporte de la valeur" ne se renforcent pas forcément l'un l'autre — ils pourraient tirer dans des directions différentes, et seule une vraie conversation (pas plus de recherche documentaire) peut trancher.

**C'est exactement pourquoi le §7 existe.** Le plan de reachouts n'est pas une formalité — c'est le seul moyen de transformer ceci d'"cohérent en théorie" à "réellement nécessaire." Tant que Lucas Heraud, Khalil Mosrati, ou quelqu'un de similaire n'a pas dit quelque chose proche de *"oui, notre flotte prend de mauvaises décisions parce qu'elle ne sait pas X"* — ça reste une hypothèse, pas une opportunité validée, et ça doit être présenté comme tel au workshop.
