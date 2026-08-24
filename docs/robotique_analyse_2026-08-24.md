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

## 3. Où Mindset Data pourrait réellement jouer un rôle — pas l'entraînement, le contexte d'exécution

La question posée en réunion (`workshop.md`, ligne 47) est : "le rôle de Mindset Data pourrait consister à préparer les données pour rendre les futurs modèles d'IA robotique **opérationnels**." Il y a une distinction cruciale entre deux problèmes différents, et un seul des deux correspond à ce que le produit fait déjà :

- **Entraîner le modèle du robot** (générer les démonstrations VLA) — pas le métier de Mindset Data, aucune brique existante ne s'en approche.
- **Donner à un robot déjà entraîné/déployé le contexte opérationnel temps réel pour agir correctement dans une usine donnée** — c'est exactement le moteur déjà construit (contextualisation ISA-95, graphe de connaissance, `kg_active_production`, détection Run/Stop). Un robot qui doit "aller chercher la pièce X" a besoin de savoir : quel OF est actif, quel est l'état de la machine, y a-t-il un problème qualité signalé — la même donnée contextualisée que le produit fournit déjà à un dashboard humain ou à un agent MCP, juste consommée par un système de contrôle robotique au lieu d'un humain.

C'est le même principe déjà établi pour le fil supply chain (`tarik.md` §1bis) : le rôle de l'IA/robot est borné et l'infrastructure de contextualisation reste inchangée — seul le consommateur final change. Positionnement honnête et défendable : **Mindset Data ne construit pas le cerveau du robot, il pourrait fournir les yeux et les oreilles contextualisées sur l'état de l'usine** dans laquelle ce robot opère.

---

## 4. Recommandation pour le workshop

- Ne pas positionner Mindset Data comme un acteur de la donnée d'entraînement robotique (VLA/démonstrations) — ce serait un claim non défendable techniquement.
- Le positionnement défendable est celui du "grounding" temps réel : contexte opérationnel structuré pour un robot déjà opérationnel, en réutilisant le moteur existant, pas un nouveau produit.
- Le calendrier réaliste du marché (déploiements Schaeffler/Hyundai à horizon 2028-2032, tâches limitées à la manutention légère/inspection aujourd'hui) suggère que ce deuxième ICP est un pari sur 2-3 ans, pas une opportunité de vente immédiate — cohérent avec la décision du workshop de le garder en exploration parallèle plutôt que d'investir dessus à fond maintenant.

Sources :
- [Physical AI Is Sending Humanoid Robots to Real Factory Floors in 2026 — Memeburn](https://memeburn.com/physical-ai-is-sending-humanoid-robots-to-real-factory-floors-in-2026/)
- [State of Robotics 2026 Report — Robotics Center of Silicon Valley](https://www.roboticscenter.ai/state-of-robotics-2026)
- [Humanoid & Quadruped Robots for Manufacturing Plants 2026 — ifactoryapp](https://ifactoryapp.com/industries/manufacturing-plant/humanoid-quadruped-robots-manufacturing-plant-2026-guide)
- [VLA Models: Training Data Requirements Explained — Shaip](https://www.shaip.com/blog/vla-models-what-vision-language-action-models-need-from-training-data/)
- [Vision-Language-Action (VLA) Models 2026 — Internet Pros](https://internet-pros.com/blog/vision-language-action-models-robotics-2026/)
