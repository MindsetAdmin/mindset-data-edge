# L2S — outreach draft (2026-08-11)

Purpose: technical co-design partnership, not investment. Target: **Didier Dumur** (Automatique et systèmes) and **Salah El Ayoubi** (Télécoms et réseaux), co-leads of L2S's named "Industry 4.0" research axis — chosen first because they're the only contacts in `docs/paris_saclay_ip_paris_outreach_2026-07-30.md` with a direct, named email (the other labs need a named contact found before a first message makes sense).

**Nothing has been sent.** Draft for review, same posture as the lemlist batch.

---

## Email (French — academic/lab audience)

**To:** didier.dumur@l2s.centralesupelec.fr, salah.elayoubi@l2s.centralesupelec.fr
**Objet :** Collaboration technique — plateforme OT/IT industrielle (Industrie du futur)

**TL;DR**
**Problème :** L'IA industrielle, c'est 80 % de préparation de données. La fragmentation OT/IT bloque l'adoption de l'IA et paralyse la prise de décision — et se prolonge en amont, jusqu'à la supply chain, dès qu'on sort du seul site industriel.
**Solution :** Une plateforme native UNS/ISA-95 ouverte qui unifie les flux OT/IT en quelques jours, avec un score de confiance et un auto-diagnostic temps réel pour rendre la donnée AI-Ready — même logique appliquée à la donnée fournisseur pour la visibilité supply chain.
**Proposition :** Pas un pilote commercial — un partenariat de co-design académique avec le L2S pour challenger l'architecture et le modèle de score sur un cas réel.

Bonjour,

Je suis Mohamed, co-fondateur et CTO de Mindset Data. Vu votre rôle de co-responsables de l'axe "Industrie du futur" du L2S, je me permets de vous contacter directement — la description de cet axe recoupe très précisément ce qu'on construit : vous parlez de produits et moyens de production qui doivent devenir "connected, interoperable, and secure", et d'une nouvelle organisation qui touche plusieurs niveaux, dont explicitement la supply chain, pas seulement l'usine isolée.

Pendant 20 ans, interconnecter l'OT et l'IT a imposé soit une pyramide rigide à 5 niveaux, soit du code point-à-point sur-mesure. Résultat : chaque nouveau besoin d'analyse devient un projet IT lourd, long, et risqué pour la stabilité du réseau usine.

Ce qui change : Mindset Data apporte une couche d'interopérabilité au-dessus de l'existant, en lecture seule (OPC-UA, MQTT, bases MySQL aujourd'hui — Modbus, S7 et d'autres protocoles industriels en cours d'ajout à la roadmap) — un bus central qui contextualise la donnée à la source selon un modèle ISA-95 et calcule un score de confiance en temps réel, plutôt que des tuyaux rigides à reconstruire à chaque nouveau besoin. Rien qui modifie le réseau, les automates ou les équipements du site. On étend aujourd'hui la même logique à la donnée fournisseur (visibilité supply chain multi-niveaux) — c'est là, concrètement, que le lien avec votre axe dépasse le seul site industriel.

Notre démarche est purement exploratoire : la première brique de la plateforme est posée, et on cherche un partenaire pour challenger la suite — en particulier le modèle de score de confiance (règles pondérées, pas du ML) sur un cas réel, ce qui est un terrain où votre expertise en systèmes cyberphysiques serait précieuse. Le co-design est 100 % gratuit et peut démarrer rapidement — je peux aussi partager une démonstration technique en amont si utile.

Est-ce un sujet que vous explorez en ce moment au L2S ? Si oui, avec plaisir pour en discuter très prochainement.

Bien cordialement,
Mohamed Khenafif
CTO, Mindset Data

---

## Notes

- **Honesty check applied**: no capability claimed beyond what's actually built (`CLAUDE.md`) — "score de confiance... approche par règles pondérées, pas du ML" matches the real `kg.AutoAcceptThreshold` mechanism, not an overclaim. "Phase de découverte, pas encore de client signé" matches `docs/decisions.md`/`docs/mindset.md`'s actual current-stage framing.
- **Not drafted yet**: IRT SystemX, LGI, LATIN, FAPS — each needs a named individual contact found first (the doc flags this gap explicitly for all four). Say the word if you want that research done next.

---

## Individual alumni — beyond labs, for technical help (2026-08-11)

Broader than the lab targets above: individual École Polytechnique alumni whose actual background could help validate or strengthen the platform/technical model — not investors. Found via 3 keyword passes on the LinkedIn alumni tool (`supply chain`, `aerospace`, `industrial IoT`) — same tool that worked for the 2026-08-05 VC search. First-page results only per pass (each returned 500-900+ total matches), picked for direct relevance, not exhaustive.

**Top pick — Christophe Marchive.** Group CISO / CISO, "Sovereign OT-IT Cyber," 20 years in critical industries — IEC 62443, EBIOS RM, ISO 27001, NIS2, CRA 2027, sectors listed: Aerospace · Naval · Sovereign Space · Rail. 2nd-degree connection. This is the single most relevant name across all three searches: he's a working OT-IT security expert in exactly the industries (aerospace, critical infrastructure) the platform targets, and exactly the credential (ISO 27001) that `docs/tarik.md` §2 flags as the real gate before a use-case-2-style deployment. A 20-30 min call with him would be a genuine sanity-check on the security roadmap, not a courtesy chat.

**Other strong candidates, by what they'd actually help with:**

| Name | Role | Degree | Why relevant |
|---|---|---|---|
| Steeve Boniface | Chief Supply Chain Data Officer | 2nd (mutual: Aurélie Pillet) | Surfaced in both the "supply chain" and "industrial IoT" searches — sits exactly at the data/supply-chain intersection `docs/tarik.md` is about |
| Xavier Dupont | Industrial & Energy Technology Strategist — "Bridging IoT, Manufacturing & Sustainability," 18K followers | 2nd | Visible industry voice in the exact OT/IoT space, could give real critique or open doors |
| Arnold Brice Michael Tayou | Industrial Data & Digitalization Engineer — Palantir Foundry, energy/environmental data platforms | 3rd | Hands-on peer: has built the same category of thing (industrial data platform) at a larger company — useful technical sounding board |
| Romain Lefrançois | Aéronautique & Espace — Supply Chain & Operations Excellence | 2nd (mutual: Nazim Nachi) | Aerospace supply chain specifically — directly maps to the Airbus RFQ example in `docs/tarik.md` |
| Laurent LeNezet | Supply chain manager at Safran | 3rd | Real practitioner inside a named aerospace Tier 1 — could sanity-check the RFQ use case from the buyer's side |
| Nicolas Sebe / Tristan Oualid | Supply Chain Scientists at Lokad | 3rd | Lokad is itself a supply-chain-optimization software company — closest thing to a technical peer/competitor conversation on this list |

**Drafted (2026-08-11, corrected same day)** — LinkedIn connection notes, French, **200-char limit** (corrected — LinkedIn's actual current cap; the earlier UK/US campaign's notes in `docs/personalized_linkedin_messages_2026-07-28.md` were built against an assumed 300 and haven't been re-checked against this — flagging, not fixing that batch here since it's 38 messages and out of scope for this ask). All under 200 chars including punctuation. **Nothing sent** — review before anything goes out.

### Christophe Marchive — priority 1 (OT-IT security) — 171 chars
> Bonjour Christophe, je construis une plateforme data OT/IT industrielle (Mindset Data). Votre expertise cyber OT-IT/IEC 62443 m'intéresse beaucoup — 20 min pour échanger ?

### Steeve Boniface — Chief Supply Chain Data Officer — 165 chars
> Bonjour Steeve, je construis Mindset Data, une plateforme data pour la supply chain industrielle. Votre profil data/supply chain m'intéresse — 20 min pour échanger ?

### Xavier Dupont — Industrial & Energy Technology Strategist — 154 chars
> Bonjour Xavier, je construis Mindset Data, une plateforme data OT/IT industrielle. Votre travail sur l'IoT industriel recoupe ce qu'on fait — échangeons ?

### Arnold Brice Michael Tayou — Industrial Data Engineer (Palantir Foundry) — 145 chars
> Bonjour Arnold, je construis Mindset Data — plateforme data industrielle, dans l'esprit Foundry mais pour l'OT/IT terrain. 20 min pour échanger ?

### Romain Lefrançois — Aéronautique & Espace, Supply Chain & Operations Excellence — 156 chars
> Bonjour Romain, je travaille sur la sélection fournisseur par la donnée en aéronautique (Mindset Data). Votre expérience terrain serait précieuse — 20 min ?

### Laurent LeNezet — Supply chain manager, Safran — 160 chars
> Bonjour Laurent, je travaille sur un outil de sélection fournisseur pour l'aéronautique (Mindset Data). Votre expérience chez Safran serait utile — échangeons ?

### Nicolas Sebe — Senior Supply Chain Scientist, Lokad — 167 chars
> Bonjour Nicolas, je construis Mindset Data — data appliquée à la supply chain, scoring traçable (pas de boîte noire). Votre expérience Lokad m'intéresse — échangeons ?

*(Tristan Oualid, also at Lokad in the same Supply Chain Scientist role, isn't drafted separately — sending both would likely read as redundant to Lokad's team; pick whichever of the two responds better to a first message, or Nicolas Sebe by default since he came up first.)*

*(Character counts above are manual estimates, not a programmatic count — recount in LinkedIn's own composer before sending, especially any note close to 190+.)*

**Sequencing suggestion**: Christophe Marchive first (highest-value, most specific fit), then Romain Lefrançois and Laurent LeNezet together (both aerospace-supply-chain, reinforce the Airbus use case from two angles), then the rest as capacity allows.
