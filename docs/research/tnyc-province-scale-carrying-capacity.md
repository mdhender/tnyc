# TNYC Province Scale and Carrying Capacity

## Summary

TNYC uses **wgvc Voronoi provinces** as its fundamental geographic unit. A typical generated TNYC map contains **10,000 land provinces**, distributed across a small number of islands/continents, with ocean occupying most of the generated world.

For economic and demographic modeling, use the traditional **six-mile wilderness hex** as a scale reference rather than as an actual subdivision that must exist in the game data.

The proposed scale is:

- **1 average land province = 7 six-mile wilderness hexes of area**
- A six-mile wilderness hex has a **3-mile apothem** (6 miles face-to-face)
- **1 wilderness hex ≈ 31.18 square miles**
- **1 average province ≈ 218.24 square miles**
- **1 average province ≈ 139,700 acres**
- **1 average province ≈ 565 km²**

The seven hexes are an **area equivalence**, not seven hexes laid end-to-end. A compact seven-hex cluster is roughly three hexes across, so a typical province is on the order of **20–25 miles across**, depending on shape.

## Current TNYC wgvc Map

A representative wgvc generation log reports:

```text
wrote tnyc (15965×6764): 45455 cells, 10000 land, 27→26 islands, merges=1, 78% ocean, 1 round(s)
```

This means:

- **45,455 total Voronoi cells** were used by the generator.
- Exactly **10,000 are land provinces**.
- Approximately **35,455 are ocean cells**.
- The generated map is approximately **78% ocean / 22% land**.
- The initial 27 landmasses became **26 islands/continents** after one merge.

The important distinction is that **10,000 is the number of playable land provinces**, not the total number of Voronoi cells.

The seven-wilderness-hex scale applies to the **land provinces**. Ocean cells are primarily part of wgvc's world geometry and need not be interpreted as seven wilderness hexes of physical area unless TNYC later gives them province-like gameplay semantics.

## Physical Scale

For a regular hex with apothem `a = 3 miles`:

```text
area = 2 × sqrt(3) × a²
     ≈ 31.177 square miles
```

Therefore:

```text
average province area
    = 7 × 31.177
    ≈ 218.24 square miles
```

Using 640 acres per square mile:

```text
218.24 × 640 ≈ 139,670 acres
```

For game purposes, **140,000 acres per average province** is a convenient round value.

## Hides and Agricultural Carrying Capacity

A medieval English **hide** was not consistently a fixed physical acreage. It evolved into an assessment of productive or fiscal capacity. The conventional modern approximation of **120 acres per hide** is nevertheless useful as a baseline for TNYC.

If every acre of an average province were ideally productive agricultural land:

```text
139,670 acres / 120 acres per hide ≈ 1,164 hides
```

Thus:

> **A perfectly arable average province has a maximum agricultural carrying capacity of approximately 1,164 equivalent hides.**

This does **not** mean the province is literally divided into 1,164 surveyed 120-acre parcels. A hide is an **equivalent productive-capacity unit**.

A useful first-order model is:

```text
Agricultural CC = 1,164 × agricultural productivity
```

where agricultural productivity ranges from 0.0 to 1.0 before other modifiers.

Examples:

| Effective productivity | Agricultural CC |
|---:|---:|
| 100% | 1,164 hides |
| 75% | 873 hides |
| 50% | 582 hides |
| 25% | 291 hides |
| 10% | 116 hides |
| 5% | 58 hides |
| 1% | 12 hides |
| 0% | 0 hides |

Population support can then be expressed separately:

```text
population capacity = agricultural CC × people supported per hide
```

The number of people supported per hide should be a game-rule parameter rather than embedded in the geography generator.

## Carrying Capacity as a General Resource Model

The same idea can be extended beyond farming.

Each province has a fixed physical area but several largely independent forms of **productive capacity**. Terrain, climate, water, geology, improvements, technology, and exploitation determine how much of each resource a province can sustainably or ultimately produce.

Agriculture, forestry, and fishing are principally **renewable flows**. Mining is fundamentally different: mineral resources are **stocks or reserves** that can be depleted.

A province might therefore have attributes conceptually like:

```text
Agricultural CC: 240 hides
Forestry CC:     780 forestry units
Fishery CC:       90 fishery units
Iron reserve:     35 ore units
Copper reserve:    8 ore units
```

The units do not all need to be called hides. The historical term **hide** is particularly appropriate for agricultural carrying capacity. Forestry and fishing can receive their own equivalent production units once their economic rules are defined.

### Land-use change

Capacities may change during play.

For example:

- clearing forest may increase agricultural CC while reducing forestry CC;
- irrigation or drainage may increase agricultural CC;
- erosion, war, or abandonment may reduce agricultural CC;
- improved forestry practices may change sustainable timber yield;
- overfishing may temporarily or permanently reduce fishery CC;
- mines consume finite reserves;
- technology may make previously marginal land or deposits economically useful.

This lets wgvc describe the **underlying potential of the land**, while TNYC models how people exploit and transform that potential.

## Recommended Raw Production Goods

The initial production system should use a deliberately small set of broad raw goods. The goal is to model strategically meaningful differences between provinces without turning the game into an inventory catalog.

### Agriculture

| Good | Scope |
|---|---|
| **Grain** | Wheat, barley, rye, oats, rice, maize, and similar staple field crops |
| **Produce** | Vegetables, fruit, legumes, roots, vineyards, orchards, and other food crops |
| **Fiber** | Flax, hemp, cotton, wool-producing agricultural capacity, and similar textile inputs |
| **Livestock** | Cattle, sheep, goats, pigs, horses, and other managed animals |

These goods draw primarily on **agricultural carrying capacity**, although livestock may also use pasture and marginal land unsuitable for crops.

### Forestry

| Good | Scope |
|---|---|
| **Timber** | Structural wood, ship timber, large beams, and quality lumber |
| **Wood** | Fuelwood, charcoal feedstock, poles, small wood, and ordinary woodland products |

Separating timber from generic wood allows shipbuilding and major construction to require high-quality forest resources without requiring a detailed tree-species simulation.

### Fishing

| Good | Scope |
|---|---|
| **Fish** | Marine, river, and lake fisheries, including shellfish and other ordinary aquatic food resources |

Fishing capacity should depend on coasts, rivers, lakes, wetlands, and local productivity rather than agricultural acreage.

### Mineral and Geological Resources

| Good | Scope |
|---|---|
| **Iron Ore** | Iron-bearing deposits |
| **Copper Ore** | Copper-bearing deposits |
| **Tin Ore** | Tin-bearing deposits; strategically important with copper for bronze |
| **Lead Ore** | Lead-bearing deposits; may also represent associated silver until silver deserves separate treatment |
| **Stone** | Building stone, including ordinary quarry products |
| **Salt** | Salt mines, brine, salt pans, and other bulk salt sources |

Unlike agriculture, forestry, and fishing, ores should normally be modeled as **finite reserves**, possibly with a production-rate limit as well as a remaining quantity.

### Initial Core List

The recommended initial raw-production vocabulary is therefore:

1. Grain
2. Produce
3. Fiber
4. Livestock
5. Timber
6. Wood
7. Fish
8. Iron Ore
9. Copper Ore
10. Tin Ore
11. Lead Ore
12. Stone
13. Salt

This list should remain intentionally broad. Specific crops, woods, metals, gemstones, spices, dyes, fine wines, furs, and similar special products can later enter the system as **luxury goods**, following the general approach of the GURPS realm-management rules.

## Comparison with Olympia, Marajanda, and Europe

TNYC deliberately occupies a different geographic scale from both its predecessor Olympia and the newer Marajanda world model.

| Property | Olympia | TNYC / wgvc | Marajanda | Europe |
|---|---:|---:|---:|---:|
| Fundamental map unit | Square province | Voronoi land province | 6-mile wilderness hex | — |
| Typical movement adjacency | 8 compass directions | Irregular graph, ~6 neighbors typical | 6 hex neighbors | — |
| Total map cells/provinces | 10,000 | 45,455 generator cells | Very large / world-scale | — |
| Land locations | ~6,000 | **10,000** | Earth-scale hex map | — |
| Water share | ~40% | **78%** | Generator-dependent | — |
| Landmasses | map-dependent | **26** in representative TNYC generation | World geography | — |
| Area represented by one land location | abstract | **~218 mi² / 140,000 acres** | **~31.18 mi² per hex** | — |
| 6-mile-hex equivalents per TNYC province | — | **7** | 1 | — |
| Approx. modeled land area | abstract | **~2.18 million mi²** | Earth-scale design | **~3.93 million mi²** |
| Relative to Europe | — | **~56%** | potentially many Europes / Earth-scale | 100% |
| Primary geographic question | Where can I move? | What can this province support? | What is in this hex? | — |

### Olympia

Olympia used a **100 × 100 square province map**, giving exactly **10,000 provinces**. Movement used the eight compass directions. With roughly 40% water, a typical map contained approximately **6,000 land provinces and 4,000 water provinces**.

TNYC therefore retains something of Olympia's familiar numerical scale while increasing the number of land locations substantially:

```text
10,000 TNYC land provinces / 6,000 Olympia land provinces ≈ 1.67
```

TNYC has roughly **67% more land provinces** than an Olympia map with 40% water.

The geometry also changes fundamentally. Instead of a square grid with eight fixed neighbors, wgvc produces irregular Voronoi provinces whose adjacencies form a planar graph.

### TNYC

At seven six-mile wilderness hexes per average province:

```text
10,000 × 7 = 70,000 wilderness-hex equivalents
```

and:

```text
70,000 × 31.177 ≈ 2.182 million square miles of land
```

That is approximately **5.65 million km²**.

TNYC is therefore best understood as a **large continental or multi-continental strategic theater**, not an Earth-sized simulation.

### Marajanda

Marajanda takes the opposite approach. Its fundamental geography is the **six-mile wilderness hex itself**, and its world-generation work has been aimed toward very large, potentially Earth-scale maps.

The conceptual distinction is useful:

> **Marajanda asks: “What is in this hex?”**
>
> **TNYC asks: “What can this province support?”**

TNYC's seven underlying wilderness hexes are therefore principally a **physical scale reference**. They do not need to exist as seven separately simulated game objects.

### Europe

Europe has an area of approximately **3.93 million square miles (10.18 million km²)**.

At approximately 2.18 million square miles of land, TNYC's 10,000 land provinces collectively represent roughly:

```text
2.18 / 3.93 ≈ 56%
```

of Europe's area.

At the TNYC province scale, all of Europe would contain approximately:

```text
3.93 million / 218.24 ≈ 18,000 provinces
```

Thus the current TNYC map is **a little over half a Europe by land area**, while still containing 10,000 strategically distinct land provinces.

## Design Principles

1. **Province geometry and productive capacity are separate concepts.** A province's physical size does not directly dictate its population.
2. **Seven wilderness hexes is an area equivalence, not an internal map.** TNYC need not simulate those hexes individually.
3. **Agricultural carrying capacity is measured in equivalent hides.** A perfect average province tops out near 1,164 hides.
4. **Carrying capacity measures productive potential, not surveyed acreage.** Poor terrain may require many physical acres to equal one productive hide.
5. **Resource capacities are multidimensional.** A province can simultaneously have agricultural, forestry, fishing, and mineral potential.
6. **Renewable resources and mineral reserves should behave differently.** Agriculture, forestry, and fishing produce sustainable flows; ores are finite stocks.
7. **wgvc generates potential; TNYC simulates exploitation.** Geography establishes the opportunity set, while population, technology, investment, war, and player decisions determine actual production.
8. **Luxury goods are a later layer.** Keep the initial raw-resource vocabulary small and add unusual or high-value goods only after the basic production economy works.
