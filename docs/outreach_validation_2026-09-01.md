# Outreach LinkedIn — valider les hypothèses sur le terrain (2026-09-01)

Objectif : transformer les hypothèses des documents workshop en réponses réelles. Rien n'a encore été envoyé à ce jour — les 13 contacts de `prospects_workshop_2026-09-01.xlsx` sont identifiés, jamais contactés.

**Posture, non négociable** (cadrage Geneviève, `call_oss_venture.md`) : **diagnostic, pas pitch.** On ne présente pas la plateforme. On pose une question précise à quelqu'un qui a vécu le problème, et on écoute. Le premier message ne vend rien et ne demande rien d'autre qu'un échange court.

**Ce qui a changé le 01/09 et qui rend cet outreach urgent** : la découverte que HighByte livre déjà un serveur MCP au bord et de la génération de modèle par LLM (`questions_workshop_2026-09-01.md` §0.A) transforme la question centrale. Elle n'est plus *« est-ce que ce problème existe ? »* mais **« pourquoi les gens qui l'ont ne l'ont pas résolu avec les produits qui existent déjà ? »**. C'est la seule question qui décide si l'offre a une place.

---

## 1. Ce qu'on cherche à valider, par ordre de priorité

| # | Hypothèse à valider | Statut | Qui peut répondre |
|---|---|---|---|
| **H1** | Le goulot est la **sémantique** (se mettre d'accord sur ce que chaque tag signifie), pas la connectique. *Chiffre à confronter : 40-60 % de l'effort d'une implémentation UNS.* | Sourcé, jamais confronté au terrain | Segment A |
| **H2** | **Les produits existants (HighByte, Litmus, UMH) sont connus et pourtant peu déployés en ETI/PME.** Si oui, pourquoi — prix, complexité, méconnaissance ? | **Non testé. La plus importante depuis le 01/09** | Segments A et B |
| **H3** | Le **build interne** est le vrai réflexe, et il produit des outils orphelins quand la personne part | Sourcé (65-85 % du coût = maintenance), jamais confronté | Segments A et B |
| **H4** | Un budget de l'ordre de **18 k€/site/an** pour cette couche est pensable, ou hors de question | Non testé | Segment B en priorité |
| **H5** | **Robotique** : les flottes AMR reçoivent des ordres transactionnels du MES mais **pas l'état machine vif** | Hypothèse rétrécie le 31/08 | Segments C et D |
| **H6** | **Le vrai point dur est la JONCTION OT↔IT, pas le mapping des tags.** Savoir que `Machine1` (automate) = `machine1` (ERP) = `M-001` (MES) — c'est là que l'effort et la valeur se concentrent | **Nouvelle (02/09).** Recadrage majeur — remplace la lecture « tout est dans les tags » | Segments A1, A2 |
| **H7** | **En pratique, les noms OT et ERP ne se ressemblent pas**, donc une correspondance exacte est inopérante et quelqu'un fait le rapprochement à la main | **Nouvelle (02/09).** Si vrai, c'est le trou exact qu'on veut combler | Segments A1, A2 |
| **H8** | Un rapprochement **automatique** OT↔IT serait-il accepté, ou faut-il de toute façon qu'un humain valide chaque lien ? | **Nouvelle (02/09).** Détermine si la porte de validation est un atout ou une contrainte | Segments A2, B |

**Critère d'abandon rappelé** : si H2 revient « on connaît, on a évalué, on a acheté » — la place sur cette piste se réduit brutalement. Si H1 revient « non, le dur c'était le réseau/les protocoles » — le positionnement entier est à revoir. Ces réponses comptent, y compris quand elles nous déplaisent.

---

## 2. Contraintes LinkedIn — à respecter, sinon le message est tronqué

- **Invitation avec note : 300 caractères maximum.** Tous les contacts sont en 2e degré → c'est le premier point de contact. Chaque note ci-dessous est comptée et vérifiée.
- **Limite d'invitations** : ~100-200/semaine en compte gratuit. Non bloquant ici (13 contacts).
- **Le message long n'arrive qu'après acceptation.** D'où le format à deux temps ci-dessous : *note d'invitation* courte → *message de relance* une fois connecté.
- Pas de Sales Navigator supposé : pas d'InMail.

---

## 3. Segment A — Responsables Automatisme & Informatique Industrielle *(priorité 1)*

**Pourquoi eux d'abord** : ce sont littéralement les « techniciens d'usine » nommés dans la demande initiale. Leur intitulé — *Responsable Automatisme **et Informatique Industrielle*** — c'est la jonction OT/IT en une personne. Ils vivent H1, H2 et H3 au quotidien.

**Reclassement assumé** : ces profils étaient étiquetés « Track B — clients AMR » dans le fichier contacts, parce que la recherche d'origine visait la robotique. **Leur meilleur usage aujourd'hui est la validation OT/IT**, qui est la piste prioritaire de la matrice.

| Contact | Poste | Localisation |
|---|---|---|
| **Sami Aloui** | Responsable Automatisme & Informatique Industriel | Lens |
| **Emmanuel Lebreton** | Responsable Automatisme et Informatique Industrielle | Pays de la Loire |
| **Bastien Charrier** | Responsable projets automatisation — SBPROCESS | Lyon |
| **Khalil Mosrati** | Responsable Robotisation & Automatisation Industrielle | Paris — *couvre aussi H5* |

### Note d'invitation (≤300 car.)

> Bonjour {Prénom}, votre poste couvre l'automatisme et l'informatique industrielle — exactement le profil que je cherche. Je mène une étude sur ce qui coûte réellement le plus cher dans un projet de données usine. 15 min pour votre retour ? Rien à vendre, je cherche à comprendre.

### Message de relance (après acceptation)

> Merci d'avoir accepté {Prénom}.
>
> Je creuse une question précise, et je n'ai que des sources écrites — pas de terrain :
>
> **Quand vous avez dû connecter une nouvelle source de donnée en usine, qu'est-ce qui a réellement pris le plus de temps : brancher le protocole, ou se mettre d'accord sur ce que chaque tag voulait dire ?**
>
> Ce que je lis dit que la modélisation représente 40 à 60 % de l'effort. Je n'ai aucune idée si ça correspond à ce que vous vivez, et c'est précisément ce que je veux vérifier.
>
> Et la question qui m'intéresse le plus **(H6/H7)** :
>
> **Chez vous, une machine porte-t-elle le même nom dans l'automate et dans l'ERP ?** Si non — qui fait le rapprochement aujourd'hui, et comment ? C'est un fichier Excel, la tête de quelqu'un, ou c'est codé en dur quelque part ?
>
> Deux autres questions si vous avez le temps :
> - Avez-vous déjà regardé des outils type HighByte, Litmus ou United Manufacturing Hub ? Évalués, écartés — et pourquoi ?
> - Un outil interne important chez vous est-il devenu difficile à faire évoluer parce que la personne qui l'avait écrit est partie ?
>
> 15-20 min en visio quand ça vous arrange, ou par écrit si vous préférez. Je partagerai ce que j'apprends des autres échanges, ça me semble la moindre des choses.

---

## 4. Segment B — CTO / DSI *(priorité 2 — build vs buy et budget)*

**Ce qu'ils valident** : H2, H3, **H4** (le budget — question qu'un responsable automatisme ne tranchera pas).

| Contact | Poste | Angle |
|---|---|---|
| **Stéphane Jaud** | CTO — Directeur Technique & Innovation, VLAD (Tours) | Le meilleur match : titre CTO, profil industriel |
| **Frédéric Kieffer** | DSI en ETI/PME, orienté métier / stratégie IT / dette technique | Langage identique à notre ciblage ETI/PME |
| **Christophe Fournel** | Expert Data & Digital Transformation, conseil CDO/CIO PME/ETI | Profil conseil — utile pour challenger, pas pour un pilote |

### Note d'invitation (≤300 car.)

> Bonjour {Prénom}, je prépare une étude sur les projets de données industrielles en ETI/PME — en particulier pourquoi tant s'arrêtent au pilote. Votre parcours correspond exactement au terrain que j'essaie de comprendre. Auriez-vous 15 min pour un échange ? Démarche exploratoire, rien à vendre.

### Message de relance (après acceptation)

> Merci {Prénom}.
>
> Je cherche à comprendre une chose précise, et je préfère demander à quelqu'un qui l'a vécu plutôt que de me fier à ce que je lis :
>
> **Pourquoi n'avez-vous pas réussi à valoriser toutes vos bases de données jusqu'ici ?** Et concrètement, comment les équipes terrain enregistrent-elles la donnée aujourd'hui ?
>
> Trois questions plus ciblées si l'échange va plus loin :
> - Des outils comme HighByte ou Litmus (≈18 k€/site/an pour cette couche) — connus chez vous ? Évalués ? Ce budget, il est dans quelle catégorie : déjà provisionné, ou impensable ?
> - Avez-vous déjà construit ce genre d'outil en interne ? Il tourne encore, et qui le maintient ?
> - Sur les trois prochaines années, combien de cas d'usage différents espérez-vous servir avec cette donnée — un, ou plusieurs ?
>
> 15-20 min, au moment qui vous convient.

*Note pour Christophe Fournel* : profil conseil. Reformuler en « je cherche quelqu'un pour challenger une analyse » plutôt que « comprendre votre terrain » — il n'a pas d'usine, il a une vision transverse de plusieurs. C'est sa valeur ici.

---

## 5. Segment C — Robotique / vendeurs AMR *(priorité 3 — valide H5)*

**Rappel** : la question du 24/08 est périmée. Ne plus demander « avez-vous de la visibilité sur le contexte usine » — la réponse 2026 est « oui, on est intégrés au MES ».

| Contact | Poste | Rôle dans la validation |
|---|---|---|
| **Romain Desarzens** | Senior Robotics SW Engineer, Movu Robotics (Paris) | **Le plus utile** : technique, peut dire ce que le fleet manager reçoit réellement |
| **Lucas Heraud** | Manager & Référent technique Déploiement, Movu Robotics | Voit les frictions d'intégration client |
| **Lukasz Tomaszewski** | Robotics Support Engineer, Locus Robotics (Birmingham) | *En anglais* |
| Benoit Boyeau · Duc-Canh Nguyen | Senior Robotics SW Engineers (Paris) | Réserve |

### Note d'invitation (≤300 car.)

> Bonjour {Prénom}, je travaille sur l'intégration entre flottes AMR et systèmes de production, et je bute sur une question à laquelle seul quelqu'un qui construit ces systèmes peut répondre. Auriez-vous 15 min ? Démarche exploratoire, je ne vends rien.

### Version anglaise (Lukasz Tomaszewski)

> Hi Lukasz, I'm researching how AMR fleets integrate with production systems, and I'm stuck on a question only someone who works on these systems can answer. Would you have 15 minutes for a short chat? Purely exploratory — nothing to sell.

### Message de relance (après acceptation)

> Merci {Prénom}.
>
> Ma question, très concrètement :
>
> **Une flotte AMR reçoit des ordres du MES — des missions, des demandes de pièces. Est-ce qu'elle reçoit aussi l'état machine en temps réel ? Si une machine tombe en panne maintenant, combien de temps avant que la flotte le sache et se réordonne ?**
>
> Je fais la distinction entre une intégration *transactionnelle* (un ordre, une demande) et un *état vif* (la machine vient de s'arrêter il y a 4 secondes). De l'extérieur je n'arrive pas à savoir si cette distinction existe vraiment sur le terrain, ou si je me la suis inventée.
>
> Et si elle existe : est-ce que ça change quelque chose d'utile, ou est-ce que le dispatcher rattrape ça très bien à la main ?
>
> Une réponse même courte, par écrit, m'aiderait déjà beaucoup.

---

## 6. Segment D — Arnaud Lubespere *(double pertinence)*

Chef de projet PMP/SAFe, **Intégrateur Automatisation Logistique SAP chez Airbus** (Toulouse). Recoupe robotique **et** supply chain (l'exemple Airbus/RFQ de `tarik.md`). Un seul contact, deux fils testables.

### Note d'invitation (≤300 car.)

> Bonjour Arnaud, votre profil croise deux sujets sur lesquels je travaille : l'automatisation logistique et la donnée industrielle chez un donneur d'ordre aéronautique. J'aimerais beaucoup avoir votre lecture. 15 min quand ça vous arrange ? Rien à vendre, je cherche à comprendre.

**En relance** : commencer par le fil automatisation logistique (§5), et n'ouvrir le fil supply chain/RFQ que si l'échange se passe bien. Ne pas empiler deux sujets dans un premier message.

---

## 7. Séquencement

| Ordre | Qui | Pourquoi maintenant |
|---|---|---|
| **1** | Segment A — les 4 responsables automatisme | Valident H1, H2, H3 : les hypothèses de la piste prioritaire |
| **2** | Stéphane Jaud, puis Frédéric Kieffer | H4 (budget) — la seule question qu'eux peuvent trancher |
| **3** | Romain Desarzens, puis Khalil Mosrati | H5 — les 2 appels qui décident si la piste robotique vit ou meurt |
| **4** | Arnaud Lubespere | Double fil, une fois les questions rodées |
| Réserve | Fournel, Boyeau, Nguyen, Tomaszewski, Heraud | Si le taux de réponse est faible |

**Rythme réaliste** : 4-6 invitations, on attend, on ajuste le message selon ce qui répond. Ne pas envoyer les 13 d'un coup — le premier retour va probablement montrer que la question est mal posée, et on aura brûlé la liste.

---

## 8. Après chaque échange — quoi consigner

Court, mais systématiquement, sinon les enseignements se perdent :
- **La réponse à H1** telle qu'ils la formulent, dans leurs mots. Ne pas traduire dans notre vocabulaire — c'est précisément le langage client qu'on cherche.
- **H2 : connaissent-ils les produits existants ?** La réponse la plus décisive de tout cet outreach.
- Toute phrase qui **contredit** un de nos documents. Elles valent plus que les confirmations.
- Le nom de la personne qui, chez eux, tiendrait le budget.

À consigner dans `analysis_log.md`, une entrée par échange.

---

## 9. Garde-fous — à ne jamais écrire dans un message

1. ❌ **« Nous sommes les seuls à… »** — faux depuis le 01/09, et un technicien qui connaît HighByte le sait.
2. ❌ **Décrire la plateforme dans le premier message.** C'est un diagnostic, pas une démo. Geneviève a été explicite là-dessus.
3. ❌ **« Est-ce que vous avez ce problème ? »** — question fermée qui appelle « non ». Toujours demander un récit : *« racontez-moi la dernière fois que… »*.
4. ❌ **Annoncer la réponse attendue dans la question.** « Le vrai problème, c'est la sémantique, non ? » ne valide rien.
5. ❌ **Promettre une démo ou un pilote** à ce stade. Rien n'est vendu, rien n'est promis.
6. ❌ **Empiler plusieurs sujets** dans une première prise de contact (cas Arnaud Lubespere).

---

## 10. Statut

**Rien n'envoyé au 2026-09-01.** Ce document est prêt à exécuter — les messages sont rédigés et vérifiés sous la limite des 300 caractères. Il manque uniquement la décision de lancer.
