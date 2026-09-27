# GURPS 4e skill gap review

Status: first pass. Tracks the follow-up in
[issue #1](https://github.com/mdhender/tnyc/issues/1).

## Purpose

Issue #1 maps each T'Nyc skill to its closest GURPS Fourth Edition skill. This
document goes the other way: it walks the GURPS 4e skill list and asks, for each
area of competence, whether T'Nyc already covers it. Areas with no reasonable
home become candidates for the unused skill numbers.

GURPS is a checklist, not a target. T'Nyc skills are deliberately broad, and a
candidate should only become a skill if it supports something a unit actually
does in the game: an order, an economic effect, a combat effect, or a detection
effect. Most GURPS skills fold into an existing T'Nyc skill or have no use in a
strategic PBM.

The GURPS skill names below come from the Basic Set skill list, filtered to
low-tech (TL0–4) and fantasy-relevant skills. They should be checked against
the book before this review is closed.

## How skill numbers are allocated

The original rules group skill numbers by kind, and within a group a higher
number costs more to learn and to maintain (rules 3.4). Keeping new skills in
the matching range is nice to have, not a design rule: the original designers
worked in the era of skill trees, and T'Nyc leans on GURPS, which does not map
naturally to one.

| Range | Group (inferred from current skills) | Free numbers |
|---|---|---|
| 1–25 | Personal, combat, and production skills | 26–31 |
| 32–41 | Magic | 42–49 |
| 50–52 | Provincial production (trade goods) | 53–64 |
| 65 | Movement | 66+ |

Magic costs are ordered separately (rules 3.4 lists them by increasing cost),
and that order does not match the skill numbers. New magic skills need a cost
position as well as a number.

## Existing gap in the rules

Rules 3.32.2 says that Cavalry "automatically bestows the skill Horse Riding",
but Horse Riding has no skill number. It belongs in the movement group (66+).
It should be resolved whether or not any other candidate is accepted.

## Candidates

Each candidate lists the GURPS skills behind it, what it would do in T'Nyc, and
a suggested range. "Strong" means an existing rule or order already needs it;
"Possible" means it fills a real gap but needs a design decision first.

### Personal, combat, and production (26–31)

| Candidate | GURPS source | Game use | Strength |
|---|---|---|---|
| Siegecraft | Engineer (Combat); Artillery (Catapult); Explosives (Demolition) | Attacking and defending cities, castles, towers, and fortresses. Nothing in the current list covers taking a walled place. | Strong |
| Construction | Architecture; Masonry; Carpentry; Engineer (Civil) | Building or improving castles, roads, bridges, and ports. Shipbuilding covers ships only. Needs a build order, which the rules do not have yet. | Possible |
| Logistics | Administration; Packing; Freight Handling; Soldier | Keeping armies supplied with food, pack animals, and freight. The game will model army supply. | Strong |
| Teaching | Teaching | Speeds STUDY for units under a teacher. Gives skilled leaders a reason to exist beyond using the skill themselves. | Possible |
| Intimidation | Intimidation; Interrogation | Improves TERRORIZE and the handling of captives (escape, EXECUTE, extracting information). | Possible |
| Law | Law; Administration (judicial) | Courts, justice, and loyalty in governed provinces. Could equally be folded into Administration. | Weak |

### Magic (42–49)

| Candidate | GURPS source | Game use | Strength |
|---|---|---|---|
| Enchantment | Enchantment college; Making and Breaking college; Alchemy | Making and identifying magical artifacts, which the rules already treat as THINGs. | Strong |
| Beast Mastery | Animal college; Animal Handling; Veterinary | Capturing and taming magical creatures. Rules 3.35 names taming as one way to fly, but no skill supports it. | Strong |
| Dispelling | Meta-Spells college (Dispel Magic, Counterspell); Exorcism | Ending spells, illusions, summonings, and undead in an area. Magic Resistance protects only its holder. | Strong |
| Mind Magic | Mind Control college; Communication and Empathy college | Magical PERSUADE, loyalty, and reading intentions. Overlaps Intrigue and Illusions, so needs a clear boundary. | Possible |
| Nature Magic | Plant college; Food college | Improving harvests and forests. The constructive counterpart to Fire Magic's "parching green valleys". | Possible |
| Warding | Protection and Warning college | Protecting a place or stack, as Magic Resistance protects a person. Could instead be a use of Dispelling. | Weak |

### Provincial production (53–64)

These skills exist to let a province produce trade goods (rules 3.36), so this
range should follow the trade-good list rather than the GURPS list. The trade
goods have not been designed yet; see
[issue #6](https://github.com/mdhender/tnyc/issues/6). GURPS suggests these crafts as a starting
checklist:

| Candidate | GURPS source |
|---|---|
| Smithing | Smith; Metallurgy; Armoury |
| Weaving | Sewing; Professional Skill (Weaver, Dyer) |
| Leatherworking | Leatherworking |
| Jewelry | Jeweler |
| Brewing | Professional Skill (Brewer, Vintner) |
| Carpentry | Carpentry |
| Horse Breeding | Animal Handling (Equines); Veterinary |

Smithing is the strongest of these: ARMOR changes a unit's armor rating, and
some supply of arms and armor presumably has to exist. Horse Breeding matters
because training Cavalry consumes a horse.

### Movement (66+)

| Candidate | GURPS source | Game use | Strength |
|---|---|---|---|
| Horse Riding | Riding (Horse) | Already referenced by the rules; see above. | Strong |
| Mountaineering | Climbing; Survival (Mountain) | Faster or safer movement through high hills, mountains, and plateaus. | Possible |
| Boating | Boating | River travel on the navigable rivers shown on the map. Only if river movement is distinct from Sailing. | Possible |
| Survival | Survival; Hiking; Navigation (Land) | Movement and attrition in deserts, wastes, and other hostile terrain. Overlaps Tracking. | Weak |

## Covered by existing skills

These GURPS areas matter to T'Nyc but already have a home. Listing them here
keeps them from being proposed again.

| GURPS skills | T'Nyc skill |
|---|---|
| Accounting; Economics; Finance | Administration, Stewardship |
| Market Analysis; Merchant; Appraisal-type uses of Connoisseur | Trading |
| Diplomacy; Politics; Streetwise; Fast-Talk; Detect Lies; Intelligence Analysis | Intrigue |
| Savoir-Faire; Heraldry; Carousing | Courtliness |
| Public Speaking; Propaganda; Leadership (civil) | Oratory |
| Shadowing; Observation; Search; Lockpicking; Filch; Escape; Disguise; Holdout; Poisons | Stealth |
| Tracking; Naturalist; Survival (Woodlands); Traps; Weather Sense | Tracking, Forestry |
| Navigation (Sea); Seamanship; Shiphandling; Meteorology (at sea) | Sailing |
| Strategy (Naval) | Naval Tactics |
| Leadership (military); Soldier | Military Leadership |
| Tactics; Strategy (Land) | Military Tactics |
| Brawling; Shield; Throwing; all melee weapon skills | Melee |
| Bow; Crossbow; Sling; Thrown Weapon | Archery |
| Riding; Lance | Cavalry |
| First Aid; Physician; Surgery; Diagnosis; Pharmacy | Medicine |
| Herb Lore; Gardening | Herbalist |
| Writing; Research; Literature | Scribe, Lore |
| History; Occultism; Hidden Lore; Area Knowledge; Cartography | Lore |
| Games; Gambling | Gaming |
| Performance; Singing; Musical Instrument; Dancing | Entertainment |
| Animal Handling (livestock) | Herding |
| Prospecting; Engineer (Mining) | Mining |

## Rejected

| GURPS area | Reason |
|---|---|
| Language; Linguistics | No language mechanics, and adding them would mostly create friction for players. |
| Theology; Religious Ritual | The setting has no religion mechanics. Revisit if temples or priests are added; Exorcism is covered by Dispelling. |
| Cooking; Housekeeping; Sewing (personal) | Too small-scale for a strategic PBM. Cooking may return as part of army supply. |
| Acrobatics; Jumping; Running; Climbing (personal) | Personal athletics with no strategic effect. Climbing is folded into Mountaineering. |
| Psionics; high-tech skills | Out of setting. |

## Decisions

1. Grouping skill numbers by kind is nice to have, not a design rule.
2. Siegecraft and Construction stay separate skills.
3. The game needs army supply, so Logistics is a candidate.
4. The trade-good list is tracked in
   [issue #6](https://github.com/mdhender/tnyc/issues/6).
5. Poisons fold into Stealth.
6. All Strong candidates are accepted. They are numbered and described in the
   [skill reference](../references/skills.md).

## Next steps

- Check the GURPS names and categories above against the Basic Set.
