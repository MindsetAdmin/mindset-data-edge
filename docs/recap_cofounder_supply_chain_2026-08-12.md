# Récap pour Cécilia — Cas d'usage Supply Chain (2026-08-12)

Résumé de tout ce qu'on a creusé sur `docs/tarik.md` depuis l'échange avec Tariq — pensé pour être lisible sans détail technique d'abord, avec la version technique et la question sécurité/certifications ensuite. Rien de ce qui suit n'est construit — c'est le plan, pas l'état actuel du produit.

---

## 1. Les deux cas d'usage, en une phrase chacun

**Cas d'usage 1 — Sélection fournisseur :** Airbus a plusieurs fournisseurs candidats pour un nouveau projet et doit choisir lequel — on l'aide à comparer, avec des données, plutôt qu'au feeling.

**Cas d'usage 2 — Monitoring continu :** suivre en continu les signaux de risque sur des fournisseurs déjà en place, pour alerter avant qu'un problème n'impacte une livraison.

**Lequel démarrer en premier n'est pas une règle fixe — ça dépend du prospect en face.** Le cas d'usage 1 est le point d'entrée le plus simple *si* le prospect a un nouveau projet de sourcing en cours au moment où on le contacte (c'est cet événement qui crée l'incitation). Mais un nouveau projet de sourcing n'arrive pas tout le temps — la plupart des entreprises n'ouvrent pas un RFQ en continu, et on n'aura pas forcément beaucoup de projets de ce type disponibles chez les prospects qu'on approche. Le cas d'usage 2 n'a besoin d'aucun événement déclencheur : n'importe quelle entreprise a déjà des fournisseurs en place, disponibles dès aujourd'hui — pas besoin d'attendre qu'un nouveau sourcing se présente.

**Nuance à garder en tête** : ça ne supprime pas totalement le problème d'incitation identifié dans l'échange avec Tariq. Le Palier 0 (donnée publique, voir §3) fonctionne aussi bien pour les deux cas d'usage, sans dépendre d'un nouveau projet — c'est le point de démarrage universellement disponible. Ce qui reste plus facile avec un nouveau projet de sourcing (cas d'usage 1), c'est d'aller au-delà du Palier 0 : la compétition entre candidats crée naturellement l'incitation à déclarer des données précises (Palier 1). Sur un fournisseur déjà en place sans nouveau contrat en jeu (cas d'usage 2 direct), obtenir plus que le Palier 0 redevient plus difficile — le problème d'incitation d'origine.

**Critère pratique pour choisir par où démarrer avec un prospect donné** : a-t-il un nouveau sourcing en cours ? Si oui, cas d'usage 1. Si non, cas d'usage 2 en Palier 0 reste immédiatement actionnable, quitte à rester à ce niveau de donnée plus longtemps.

---

## 2. Comment ça marche — version non technique

**Le problème qu'on résout** : aujourd'hui, un acheteur choisit un fournisseur en grande partie sur ce que le fournisseur déclare lui-même dans sa réponse à l'appel d'offres — capacité, délais, qualité. Personne ne croise ça avec une source indépendante avant de signer.

**Ce qu'on ajoute** : on va chercher, en parallèle, des signaux publics sur chaque fournisseur candidat (santé financière, certifications qualité à jour, mentions presse négatives — rachat, faillite, litige) et on les confronte à ce que le fournisseur a déclaré. Le résultat : un classement des candidats avec, pour chacun, le détail de pourquoi il est classé où il est — pas juste un chiffre sorti d'une boîte noire.

**Pourquoi un fournisseur accepterait de jouer le jeu** : parce qu'il est déjà en compétition pour obtenir le contrat — donner des données précises et vérifiables, c'est ce qui l'aide à gagner l'appel d'offres. On ne demande rien de nouveau qui n'existe pas déjà dans un processus RFQ classique ; on le rend juste plus fiable.

**Pour le cas d'usage 2** : même principe, mais en continu au lieu d'une comparaison ponctuelle — si un signal de risque dépasse un certain seuil (retard qui s'allonge, certification qui expire), on alerte avant que ça touche la livraison. Ça marche aussi bien sur un fournisseur choisi via le cas d'usage 1 (données déjà là, rien à redémarrer) que sur un fournisseur déjà en place chez le client, sans lien avec un cas d'usage 1 — dans ce second cas, on démarre en Palier 0 (donnée publique) comme point d'entrée, voir §1.

---

## 3. Comment ça marche — version technique (résumé)

Détail complet dans `docs/tarik.md` §1/§1bis/§1ter. Les points clés :

- **Même moteur que la plateforme actuelle** (usine OT/IT) — pas un nouveau produit. Ingestion en lecture seule, normalisation avec score de confiance, graphe de connaissance contextualisé, appliqués à la donnée fournisseur au lieu de la donnée machine.
- **Trois niveaux de donnée, pas un accès uniforme à tout** :
  - **Palier 0** — donnée publique, aucune coopération fournisseur nécessaire (santé financière, certifications, presse).
  - **Palier 1** — donnée déclarée par le fournisseur lui-même dans sa réponse RFQ, croisée avec le Palier 0.
  - **Palier 2** — donnée plus profonde (carnet de commandes, finances détaillées), uniquement une fois la relation installée (voir `tarik.md` pour le détail du déclencheur et la question des contrats).
- **Le score est une formule transparente, pas un modèle IA qui "décide"** — pondération par palier de confiance + vérification croisée, traçable jusqu'à chaque donnée source. Aucun poids caché, aucune donnée d'entraînement.
- **Où l'IA intervient réellement** : lire des documents non structurés pour en extraire des champs (ex. un PDF de certification), et permettre à quelqu'un d'interroger le graphe en langage naturel pour comprendre un score — jamais pour décider du classement lui-même.

---

## 4. Sécurité — ce qu'il faut, comment l'obtenir, et en combien de temps

Point important à ne pas mélanger : **ISO 27001 et SOC 2 sont de vraies certifications**, avec un auditeur externe accrédité — on ne les "obtient" pas soi-même, il faut passer un audit. **SIG Lite et CAIQ ne sont pas des certifications** : ce sont des questionnaires d'auto-évaluation qu'on remplit nous-mêmes et qu'on transmet directement au client qui les demande (pas d'organisme certificateur, pas d'audit externe) — donc beaucoup plus rapides, mais aussi moins "prouvés" qu'un audit indépendant.

| Palier / cas d'usage | Ce qui est probablement demandé | C'est quoi, concrètement | Comment l'obtenir | Délai réaliste |
|---|---|---|---|---|
| **Palier 0** (RFQ, donnée publique seule) | Rien | — | — | Immédiat, pilote possible dès aujourd'hui |
| **Palier 1** (RFQ, donnée déclarée réelle) | Questionnaire de sécurité — **SIG Lite** ou **CAIQ (Lite)**, selon ce que demande le client | Auto-évaluation déclarative de nos pratiques de sécurité (~125 questions chacun — SIG Lite = 126, CAIQ Lite = 124), pas un audit tiers | On le remplit nous-mêmes et on l'envoie directement au client (Airbus). Le CAIQ peut aussi être publié gratuitement sur le registre public de la Cloud Security Alliance (CSA STAR, niveau auto-évaluation) pour donner de la visibilité sans attendre qu'on nous le demande | Quelques jours à ~1 semaine si notre documentation de sécurité interne est déjà à jour ; plus long la toute première fois, le temps de la rédiger |
| **Palier 1, si le client veut une preuve indépendante** | **SOC 2 Type 1** (étape intermédiaire, pas obligatoire mais un signal fort) | Vrai audit, mais à un instant T (pas sur une période comme le Type 2) | Faire appel à un cabinet d'audit accrédité (ex. via une plateforme comme Vanta/Drata pour préparer, puis un auditeur indépendant certifie) | ~2-3 mois |
| **Cas d'usage 2** (monitoring continu, relation de données persistante) | **ISO 27001** | Vrai audit indépendant, certificat valable 3 ans avec audits de surveillance annuels | Mise en place d'un SMSI (système de management de la sécurité de l'information), audit interne, puis audit de certification par un organisme accrédité | **6-18 mois** en réalité — à démarrer dès qu'un cas d'usage 2 réel se profile, pas au moment de la négociation |

**Ce qu'on peut faire dès maintenant, sans attendre aucune certification** : préparer un DPA (accord de traitement de données) et un NDA type — ça ne nécessite aucun audit, juste un travail juridique, et c'est ce qui couvre réellement le Palier 1 le plus tôt possible (voir `tarik.md` §2 et §"Ce qu'il faut dire honnêtement").

**Recommandation pratique** : commencer par SIG Lite ou CAIQ (rapide, gratuit, suffisant pour un premier pilote Palier 0/1), lancer le SOC 2 Type 1 si un client sérieux le demande explicitement en Palier 1, et ne démarrer ISO 27001 que quand le cas d'usage 2 devient concret avec un acteur du niveau Airbus — pas avant, vu le délai.

---

*Source complète, avec tout le détail (modèle technique, comparaison marché, sources Palier 0) : `docs/tarik.md`.*

Sources (recherche web 2026-08-12) :
- [SIG Lite vs CAIQ: Which Vendor Questionnaire to Use?](https://www.shieldrisk.ai/blog/sig-lite-vs-caiq/)
- [SIG questionnaire guide — Vanta](https://www.vanta.com/collection/trust/sig-questionnaire)
- [CAIQ vs SIG Assessment — HyperComply](https://www.hypercomply.com/blog/caiq-vs-sig)
- [CAIQ vs. SIG Questionnaires — Bitsight](https://www.bitsight.com/blog/caiq-vs-sig-top-questionnaires-vendor-risk-assessment)
