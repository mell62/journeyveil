# Journeyveil — Product Specification

**Status:** Draft MVP specification  
**Platforms:** Android first; iOS is a later phase  
**Game type:** Location-based walking and exploration

## 1. Premise

Journeyveil turns real-world walking into the discovery of a hidden map. Players reveal previously unexplored hexagonal **Cells**, encounter **Anomalies**, complete photography or shape-walking challenges, and collect **Map Fragments**. Completed fragment puzzles reveal collectible pictures of landmarks around the world. A personal Journal preserves discoveries, photographs, walked shapes, and puzzles.

The primary experience is a tilted, third-person, avatar-centered 3D map. The character is shown on a map that follows the player's location and rotates with their heading.

## 2. Goals and boundaries

### Player goals

- Make an ordinary walk feel like progress and discovery.
- Offer a rewarding challenge when an Anomaly is found, without requiring a scenic landmark in every neighborhood.
- Give collected fragments a clear purpose: completing collectible pictures of landmarks around the world.
- Make past walks worth revisiting through an attractive, private Journal.

### MVP boundaries

- One default character; no character customization, combat, social features, or trading.
- Android foreground and background exploration. iOS behavior is out of scope for the first release.
- Photography challenges target stable subjects like landmarks (less XP) or variable subjects like sunrise/sets, birds, cloudy sky (more XP).

## 3. World and progression

### Cells

- The world is divided into stable hexagonal Cells. Each Cell has an ID, geographic boundary, and exploration state for each player.
- All Cells start **unexplored** for a new player. The map conceals unexplored terrain and reveals Cells as the player walks through them.
- Cells are about 500 m across. Discover one after two GPS readings inside it, at least 10 seconds apart and accurate to 25 m. Require a third reading near boundaries.
- The first valid discovery of a Cell grants **base XP once**.
- Field-test Cell size and GPS tolerance so shape walks fit and GPS drift does not reveal extra Cells.

### Anomalies

- Roughly one in eight Cells contains a discoverable Anomaly. Its presence is revealed when the player discovers the Cell, not while it remains unexplored.
- An Anomaly has one challenge and three player-facing states: **available**, **in progress**, and **completed**. A failed or interrupted attempt returns it to available; a completed Anomaly cannot award its reward again.
- An Anomaly near a curated landmark offers a photography challenge; elsewhere it offers a shape-walking challenge.
- Anomalies are assigned consistently so revisiting a Cell does not reroll its challenge or rewards.

### XP

- Start with 10 XP for discovering a Cell and 50 XP plus one Map Fragment for completing an Anomaly. Reaching level 2 takes 50 XP; each subsequent level requires 50 XP more than the previous one.
- The reward screen distinguishes discovery XP from challenge XP. The profile corner shows the character portrait, current level, and progress to the next level.
- Tune XP, levels, and Anomaly frequency using pilot walks and completion rates.
- Rewards are issued once per player and event, including after retries, restarts, or delayed background synchronization.

## 4. Challenges

### Photography

1. A photography Anomaly names a curated landmark to photograph (for example, “Photograph the old lighthouse”).
2. The player takes a photo in the game within 100 m of an approved public viewpoint.
3. The photo passes if its capture time and location are valid and the requested subject is verified by AI. Uncertain matches can be retried.
4. On success, the player receives bonus XP and a Map Fragment; the photo and challenge result appear in the Journal.
5. On failure, the player can retry.

### Shape walking

1. A shape-walking Anomaly asks the player to walk an approximate **triangle** or **square** in its Cell.
2. The game shows the target shape, starts recording when the player chooses **Start**, and renders their route on the map.
3. On **Finish**, accept a route of at least 200 m that ends within 40 m of its start, with three or four clear turns and sides no more than twice as long as each other. Allow for GPS noise and detours.
4. On success, the player receives bonus XP and a Map Fragment; a simplified route sketch and result appear in the Journal.
5. On failure, the player can retry.

## 5. Landmark catalog

Maintain a curated global set of landmarks for photography challenges and collectible pictures.

## 6. Map Fragments and landmark pictures

- Every completed Anomaly awards one fragment for an active landmark puzzle, regardless of where the player is walking.
- A puzzle represents a landmark anywhere in the world; its fragments form a collectible jigsaw-style picture.
- The collection view shows fragment slots and completion progress without revealing the landmark picture too early.
- Completing the puzzle reveals the landmark's name and location. The completed picture remains available in the Journal.
- Fragment awards must advance a puzzle rather than become duplicates. Once a puzzle is complete, subsequent awards advance another available puzzle.

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

- Request location permission with a plain-language explanation of why exploration needs it. Request background location if the player also opts into background exploration.
- On Android, background exploration uses an ongoing foreground location service with a persistent notification and an obvious stop control. The game must handle denied or revoked permissions, disabled location, process restarts, and battery restrictions.
- If tracking pauses, show that state clearly. Do not invent walked Cells, route segments, or challenge progress for a period without trustworthy readings.
- Batch and deduplicate discoveries when connectivity returns. A completed action must not grant XP or fragments twice.
- Minimize battery usage by adapting location sampling to movement and challenge state rather than continuously requesting maximum-accuracy GPS.

## 9. Tech stack

### Client

| Area | Technology |
|---|---|
| Game and app UI | Godot GDScript for the map experience, challenges, Journal, puzzles, achievements, and profile. |
| Avatar | Animated 3D character in Godot, with idle/walking transitions through AnimationTree. |
| Geographic map | Mapbox raster tiles rendered as textured ground tiles in Godot. Godot renders the tilted camera, avatar, concealed/revealed Cells, anomalies, and routes. |
| World grid | H3 |
| Android capabilities | A Kotlin Godot plugin for permissions, camera integration, authentication, and a foreground location service. |
| Location | Android Fused Location Provider, with sampling adapted to movement and challenge state. |
| Local persistence | SQLite for discoveries, challenge attempts, Journal metadata, and a sync queue. |

The MVP ground map is flat: buildings and labels are baked into tile images, so there are no extruded buildings and labels tilt with the ground. Validate this appearance on an Android device before committing beyond the prototype. Implement geographic positioning, tile loading, and permitted caching; verify Mapbox pricing, custom-renderer and offline terms, and required attribution before release.

### Backend and infrastructure

| Area | Technology |
|---|---|
| API | Go modular API. |
| Database | PostgreSQL + PostGIS for player progression, challenge assignments, landmark catalog, and geographic queries. |
| Authentication | Auth0, integrated through its Android SDK. |
| Media | Cloudflare R2 |
| Photo verification | An LLM API, called by the Go backend. |
| Background jobs | Cloud Tasks → Go worker on Cloud Run. |
