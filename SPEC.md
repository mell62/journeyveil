# Journeyveil — Product Specification

**Status:** Draft MVP specification  
**Platforms:** Android first; iOS is a later phase  
**Game type:** Location-based walking and exploration

## 1. Premise

Journeyveil turns real-world walking into the discovery of a hidden map. Players reveal previously unexplored hexagonal **Cells**, encounter **Anomalies**, complete photography or shape-walking challenges, and collect **Map Fragments**. Completed fragment puzzles reveal collectible pictures of landmarks around the world. A personal Journal preserves discoveries, photographs, walked shapes, and puzzles.

The primary experience is a tilted, third-person, avatar-centered 3D map. This is **not camera-based augmented reality**: the character is shown on a map that follows the player's location and rotates with their heading.

## 2. Goals and boundaries

### Player goals

- Make an ordinary walk feel like progress and discovery.
- Offer a rewarding challenge when an Anomaly is found, without requiring a scenic landmark in every neighborhood.
- Give collected fragments a clear purpose: completing collectible pictures of landmarks around the world.
- Make past walks worth revisiting through an attractive, private Journal.

### MVP boundaries

- One default character; no character customization, combat, social features, or trading.
- Android foreground and background exploration. iOS behavior is out of scope for the first release.
- Exploration and challenges should work wherever location and map data support them. The landmark-picture catalog is global and independent of the player's location.
- Photography prompts target stable, identifiable subjects (such as monuments or viewpoints). Sunset and weather-dependent prompts are deferred or optional, never the sole way to finish an Anomaly.
- No AR camera overlay. The camera is used only to take challenge photographs.

## 3. World and progression

### Cells

- The world is divided into stable hexagonal Cells. Each Cell has an ID, geographic boundary, and exploration state for each player.
- All Cells start **unexplored** for a new player. The map conceals unexplored terrain and reveals Cells as the player walks through them.
- A Cell becomes **discovered** only after sufficiently reliable location readings place the player inside it. No manual tap-to-reveal or retroactive import of location history.
- The first valid discovery of a Cell grants **base XP once**. Revisits do not grant discovery XP again.
- Hex size and GPS tolerance must be tuned together in outdoor testing: the Cell must be large enough to support useful walks and shape challenges without location noise repeatedly moving a player between Cells.

### Anomalies

- Some Cells contain a discoverable Anomaly. Its presence is revealed when the player discovers the Cell, not while it remains unexplored.
- An Anomaly has one challenge and three player-facing states: **available**, **in progress**, and **completed**. A failed or interrupted attempt returns it to available; a completed Anomaly cannot award its reward again.
- In an eligible scenic Cell, the Anomaly offers a photography challenge. Otherwise it offers a shape-walking challenge. A photography challenge that is inaccessible, unsafe, or cannot be verified must have a non-photography fallback; it must not permanently strand the player's reward.
- Anomalies are assigned consistently so revisiting a Cell does not reroll its challenge or rewards.

### XP

- Discovering a Cell grants a small amount of XP. Completing an Anomaly grants a substantially larger bonus and one Map Fragment.
- The reward screen distinguishes discovery XP from challenge XP. The profile corner shows the character portrait, current level, and progress to the next level.
- Exact XP amounts, level curve, Anomaly frequency, and fragment counts are balancing parameters to settle during playtesting, not hard-coded product promises.
- Rewards are issued once per player and event, including after retries, restarts, or delayed background synchronization.

## 4. Challenges

### Photography

1. A scenic Anomaly presents a specific, observable subject and clear instructions (for example, “Photograph the old lighthouse”).
2. The player takes a photo in the game while near the Anomaly's eligible public location.
3. The game checks capture time and location eligibility, then analyzes whether the requested subject appears in the image. **Beauty or photographic quality is not a pass condition.**
4. On success, the player receives bonus XP and a Map Fragment; the photo and challenge result appear in the Journal.
5. On failure, the player receives a useful reason and can retry or choose the fallback where available. An uncertain AI result must not be presented as a confident rejection.

Image analysis is only one signal: it cannot prove where an image was taken, and GPS alone cannot prove what a camera could see. Define review thresholds, retry behavior, and abuse controls before release. Do not require players to photograph strangers, enter private property, or cross unsafe terrain.

### Shape walking

1. A non-scenic Anomaly asks the player to walk an approximate **triangle** or **square** in its Cell.
2. The game shows the target shape, starts recording when the player chooses **Start**, and renders their route on the map.
3. On **Finish**, the route is checked for a minimum meaningful distance, a closed loop, the requested number of turns, and approximate similarity to the target. GPS drift, reasonable detours, and imperfect angles are tolerated.
4. On success, the player receives bonus XP and a Map Fragment; a simplified route sketch and result appear in the Journal.
5. An interruption can pause or end the attempt without losing previously earned rewards. Unreliable tracking produces an inconclusive attempt rather than a false failure.

Shape walking must not encourage road crossing, trespassing, or staring at the screen while moving. Test whether a shape can realistically be walked on public routes within the chosen Cell size; provide a way to report unsafe or impractical challenges. Do not assume a mathematically perfect shape is possible on local streets.

## 5. Scenic-place selection

Journeyveil does **not** manually inspect every Cell or classify every location as scenic.

- Gather candidate subjects from geographic sources such as OpenStreetMap and Wikidata: monuments, landmarks, viewpoints, notable natural features, and similar places.
- Rank candidates by subject type, documented significance, available evidence (including geotagged imagery where licensing allows), public accessibility, and confidence in location and identity.
- Generate a specific photography prompt only above a conservative confidence threshold. A missing or low-confidence candidate produces a shape-walking Anomaly instead.
- Check that the player can reach an eligible public position; a POI coordinate alone does not guarantee visibility, access, or safety.
- Provide player reports for missing subjects, inaccessible locations, inaccurate prompts, and unsafe routes. Use reports and challenge outcomes to retire bad prompts and prioritize human review of popular candidates.
- Curate or review the much smaller global set of **featured puzzle landmarks**. Automated scenic scoring is a candidate generator, not a guarantee that a landmark is worth featuring.

Location data coverage and access rules vary by region. Use shape-walking challenges where reliable photography prompts are unavailable.

## 6. Map Fragments and landmark pictures

- Every completed Anomaly awards one fragment for an active landmark puzzle, regardless of where the player is walking.
- A puzzle represents a landmark anywhere in the world; its fragments form a collectible jigsaw-style picture.
- The collection view shows fragment slots and completion progress without revealing the landmark picture too early.
- Completing the puzzle reveals the landmark's name and location. The completed picture remains available in the Journal.
- Fragment awards must advance a puzzle rather than become unusable duplicates. Once a puzzle is complete, subsequent awards advance another available puzzle; if the catalog is exhausted, the game must offer a defined alternative reward before enabling more Anomalies.
- Featured landmarks must be verified as real and suitable to show publicly. A reveal is not a navigation instruction or assurance of current access.

**Catalog decision before launch:** provide enough reviewed landmark pictures to sustain expected Anomaly completions, or define the alternative reward when the catalog is exhausted.

## 7. Journal and interface

### Main map

- Tilted 3D perspective, a single default avatar at the player's position, and a map that follows movement.
- The map rotates with device heading when heading quality is usable; provide a north-up control and a stable fallback when compass data is noisy.
- Day and night map palettes follow local time. Revealed and unexplored Cells must remain distinguishable in both palettes.
- A corner profile displays character portrait, level, and XP bar. Primary navigation provides **Journal**, **Achievements**, and **Statistics**.
- Anomalies and challenge states have distinct, readable map markers.

### Journal

- A chronological discovery feed alongside collections for **Photos**, **Walked Shapes**, **Map Fragments**, and **Completed Puzzles**.
- Entries show date, approximate location, challenge, and reward. Photo entries show the player's photograph; shape entries show a simplified route rather than requiring raw location history.
- Players can view and delete their own photos and detailed route data. Journal entries remain legible if associated media is removed.

### Achievements and statistics

- MVP achievements cover basic milestones (first Cell, first Anomaly, first photo, first shape, first completed puzzle).
- Statistics include Cells discovered, Anomalies completed, XP/level, fragments collected, and puzzles completed. Do not equate GPS recording time with distance walked.

## 8. Location, privacy, and reliability

- Request location permission with a plain-language explanation of why exploration needs it. Request background location only when the player opts into background exploration.
- On Android, background exploration uses an ongoing foreground location service with a persistent notification and an obvious stop control. The game must handle denied or revoked permissions, disabled location, process restarts, and battery restrictions.
- If tracking pauses, show that state clearly. Do not invent walked Cells, route segments, or challenge progress for a period without trustworthy readings.
- Batch and deduplicate discoveries when connectivity returns. A completed action must not grant XP or fragments twice.
- Minimize battery usage by adapting location sampling to movement and challenge state rather than continuously requesting maximum-accuracy GPS.
- Location history and photos are sensitive. Specify what stays on-device, what is uploaded for validation/sync, retention periods, deletion, and consent **before implementation**. Do not expose precise personal routes in public features.
- Apply basic spoofing and abuse controls without treating legitimate GPS uncertainty as cheating.

## 9. MVP acceptance criteria

1. A new account sees unexplored Cells; a valid outdoor walk reveals them and awards discovery XP once each.
2. With background exploration enabled, Android continues recording eligible discoveries while the game is not foregrounded, displays a persistent notification, and stops when the player chooses Stop.
3. Discovering an Anomaly reveals the appropriate challenge. A successful photo or shape awards bonus XP and one fragment exactly once.
4. A low-confidence scenic location does not receive an invented photography prompt. A bad or unsafe prompt can be reported and does not permanently block progress.
5. Completing a puzzle reveals its featured landmark and completed picture; awarded fragments never become unusable duplicates.
6. The Journal shows completed photos, simplified shapes, fragments, and puzzles. Deleting sensitive media or route detail works as described in the privacy policy.
7. The avatar follows the player; map rotation, north-up fallback, day/night palettes, XP bar, and primary navigation work without obscuring the walking route.
8. Permission denial, noisy GPS, interruption, and offline synchronization do not fabricate progress or duplicate rewards.

## 10. Decisions to validate in the pilot

- Initial worldwide landmark-picture catalog size and how to replenish it as players complete puzzles.
- Hex size, discovery accuracy/dwell rules, Anomaly density, and XP/level progression.
- Shape-scoring tolerance and whether local public routes support triangle and square challenges inside a Cell.
- Photography eligibility radius, subject-verification threshold, fallback policy, and validation costs.
- Photo/route storage, synchronization, retention, and deletion policy.
- Accessibility alternatives for players unable to complete a particular walking or photography task.
