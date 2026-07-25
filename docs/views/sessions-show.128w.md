---
level: 128w
parent: sessions.32w.md
deeper: sessions-show.256w.md
relates-to:
  - ../controllers/sessions.128w.md
source: views/sessions/show.tpl
---

Live workout session page. Includes navbar partial, exercise-fields partial (for the "Add Exercise" form), and all three
goal modals. Renders exercise blocks (main, abs, cardio, stretch) with sets tables. Set-logging forms submit via AJAX (
`fetch`) and append rows without reload; delete-set forms also use AJAX. Exercises in non-cardio blocks are
drag-to-reorder via SortableJS; reorder fires `POST /sessions/:id/exercises/reorder`. A fixed-bottom rest timer starts
on each set logged *unless the card flows straight into the next exercise* — `flowsIntoNext()` reads `data-linked`
(superset) or `data-in-circuit` (circuit member), both computed by the controller and never a raw column. It uses the
Web Audio API for its alarm and persists across page reloads via `localStorage`.

Supersets: every exercise card except the last in its block carries a chain button (`bi-link-45deg` off,
`bi-link` on) that toggles `POST /sessions/:id/exercises/:eid/link` and updates the card in place — the page never
reloads, so a running rest timer survives. Cards carry `data-link-raw` (the stored `linked_to_next`) and `data-linked`
(the computed effective link); members of a run show an `A1`/`A2` badge and share a left accent stripe. After a set is
logged, the timer starts only when `data-linked` is not `"true"`. `relabelBlock()` recomputes the effective links and
labels client-side after a toggle or a drag, mirroring `groupSessionExercises` in the controller, so a card dragged to
the bottom of its block loses its chain and rests again. A global lb/kg toggle converts all exercises at once; each non-time-based exercise card also has its own
per-exercise lb/kg toggle that persists the preference to the exercise library via
`POST /sessions/:id/exercises/:eid/unit`. Each exercise card header has a three-dots (⋮) dropdown button — icon shows ↑/↓ for HitMax/BelowGoal — whose menu
contains "Edit goal" (opens the appropriate goal modal) and "Change exercise" (modal with autocomplete that POSTs
`name` to `/sessions/:id/exercises/:eid/change` and reloads the page). Exercise autocomplete fills type radio from `.ExerciseLibraryJSON`.

Circuits render from `.Circuits` as one dark-headed card per circuit ("N rounds · Ns transition") holding its members'
cards, each with its work-seconds badge, its rounds table and a manual "+ Round" button. Members appear only here, never
also in their block, and carry no chain toggle.

The **Start Circuit** button opens `#circuit-runner`, a full-screen overlay on this same page — not a separate page,
because the session holds a running rest timer in JS that a navigation would lose. The runner flattens the circuit into
a plan (a lead-in gap, then each member's work interval per round, separated by the transition gap) and walks it:
work interval → **transition cue** → gap → **start cue** → next. Every countdown is a wall-clock delta from
`Date.now()`, never accumulated `setInterval` ticks, because a backgrounded mobile tab throttles timers badly. Pause
shifts the step's start time by the pause length; Skip abandons the current interval unlogged; Quit closes the overlay
and leaves every completed interval logged.

Each completed interval is POSTed to the ordinary set endpoint — round N is set N — and the row is appended to the
member's card underneath, so the page is already correct when the overlay closes and nothing reloads.

The two cues are deliberately opposite in shape so they can be told apart with the phone in a pocket: **start** is three
short rising beeps (660 → 880 → 1100 Hz), **transition** is two long falling tones (440 → 294 Hz). Both reuse the rest
timer's single `AudioContext` — iOS unlocks one per gesture and will not unlock a second — and both dispatch a
`circuit-cue` DOM event carrying the cue name and its tones, which the overlay turns into a colour flash for a muted
device.
