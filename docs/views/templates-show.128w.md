---
level: 128w
parent: templates.32w.md
deeper: templates-show.256w.md
relates-to:
  - ../controllers/templates.128w.md
source: views/templates/show.tpl
---

Read-only template detail page. Includes `partials/navbar.tpl`. Shows template name and optional focus text. An "Edit"
button links to `/templates/:id/edit`. Below the header, exercises are grouped into blocks (`.ExerciseBlocks`); each
block renders a section heading and a list of exercise names (with "Bodyweight" note where applicable). Empty state
shown when no exercises. `.Success` flash auto-dismissed after 3 s. No AJAX.
