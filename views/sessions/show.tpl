<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Session #{{.Session.WorkoutNumber}} — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
    <style>.drag-handle { cursor: grab; touch-action: none; } .sortable-ghost { opacity: 0.4; }</style>
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4 mb-4" style="max-width: 600px;">
    <div class="mb-4">
        <a href="/programs/{{.Program.ID}}" class="text-muted small">&larr; {{.Program.Name}}</a>
        <div class="d-flex align-items-center gap-2 mt-1">
            <h1 class="h4 fw-bold mb-0">Session #{{.Session.WorkoutNumber}}</h1>
            {{if .Session.IsDeload}}
            <span class="badge bg-secondary">Deload</span>
            {{end}}
        </div>
        <p class="text-muted small mb-0 mt-1">
            Phase {{.Session.PhaseNumber}} &middot; Week {{.Session.WeekNumber}}
            &middot; {{.Session.Date.Format "Jan 2, 2006"}}
            {{if gt .PhaseRestSeconds 0}}&middot; {{.PhaseRestSeconds | restMinutes}}m {{.PhaseRestSeconds | restSecs}}s rest{{end}}
        </p>
    </div>

    {{range .ExerciseBlocks}}
    {{if ne .Block "main"}}
    <h2 class="h6 fw-semibold text-uppercase text-muted mt-4 mb-3">{{.Label}}</h2>
    {{end}}

    {{if eq .Block "cardio"}}

    {{range .Exercises}}
    {{if .CardioLogs}}
    {{$exID := .Exercise.ID}}
    {{$exName := .Exercise.Name}}
    {{range .CardioLogs}}

    <div class="card mb-2">
        <div class="card-body py-2">
            <div class="d-flex align-items-center justify-content-between">
                <div>
                    <span class="fw-semibold text-capitalize small">{{$exName}}</span>
                    {{if .CardioType}}<span class="text-muted small ms-2 text-capitalize">{{.CardioType}}</span>{{end}}
                </div>
                <div class="d-flex align-items-center gap-3">
                    <span class="text-muted small">
                       Goal: {{if gt .GoalDuration 0}}{{.GoalDuration}} | Actual: {{end}}{{.ActualDuration}} min
                    </span>
                    <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{$exID}}/cardio/{{.ID}}/delete" class="d-inline">
                        <button type="submit" class="btn btn-link btn-sm text-danger p-0" title="Delete"><i class="bi bi-trash"></i></button>
                    </form>
                </div>
            </div>
        </div>
    </div>
    {{end}}
    {{else}}
    {{$exID := .Exercise.ID}}
    <div class="card mb-3">
        <div class="card-body pb-2">
            <div class="d-flex align-items-baseline justify-content-between mb-2">
                <h2 class="h6 fw-semibold mb-0 text-capitalize">{{if .HitMax}}<button type="button" class="btn btn-link p-0 border-0 hit-max-btn" data-bs-toggle="modal" data-bs-target="#goalWeightModal" data-ex-name="{{.Exercise.Name}}" data-goal-weight="{{.Exercise.GoalWeight}}" data-weight-unit="{{.Exercise.WeightUnit}}" data-direction="up" title="Hit max reps last workout — tap to update goal weight" style="line-height:1;vertical-align:middle;"><i class="bi bi-arrow-up-circle-fill text-black" style="font-size:1.0em;"></i></button>&nbsp;{{else}}{{if .Exercise.IsTimeBased}}<button type="button" class="btn btn-link p-0 border-0" data-bs-toggle="modal" data-bs-target="#goalSecondsModal" data-ex-name="{{.Exercise.Name}}" data-goal-seconds="{{.Exercise.GoalSeconds}}" title="Set goal duration" style="line-height:1;vertical-align:middle;"><i class="bi bi-pencil text-black" style="font-size:1.0em;"></i></button>&nbsp;{{else}}{{if .Exercise.IsBodyweight}}<button type="button" class="btn btn-link p-0 border-0" data-bs-toggle="modal" data-bs-target="#goalRepsModal" data-ex-name="{{.Exercise.Name}}" data-goal-rep-min="{{.GoalRepMin}}" data-goal-rep-max="{{.GoalRepMax}}" title="Set goal reps" style="line-height:1;vertical-align:middle;"><i class="bi bi-pencil text-black" style="font-size:1.0em;"></i></button>{{else}}<button type="button" class="btn btn-link p-0 border-0 hit-max-btn" data-bs-toggle="modal" data-bs-target="#goalWeightModal" data-ex-name="{{.Exercise.Name}}" data-goal-weight="{{.Exercise.GoalWeight}}" data-weight-unit="{{.Exercise.WeightUnit}}" data-direction="down" title="{{if gt .Exercise.GoalWeight 0.0}}Missed max reps — tap to adjust goal weight{{else}}Set goal weight{{end}}" style="line-height:1;vertical-align:middle;"><i class="bi bi-{{if gt .Exercise.GoalWeight 0.0}}dash-circle-fill{{else}}pencil{{end}} text-black" style="font-size:1.0em;"></i></button>{{end}}&nbsp;{{end}}{{end}}{{.Exercise.Name}}</h2>
                <div class="d-flex align-items-center gap-2">
                    <span class="text-muted small">
                    {{if .Exercise.IsTimeBased}}
                    {{if gt .Exercise.GoalSeconds 0}}Goal: {{fmtDuration .Exercise.GoalSeconds}}{{end}}
                    {{else}}
                    {{if gt .Exercise.GoalWeight 0.0}}Goal: {{.Exercise.GoalWeight}} {{.Exercise.WeightUnit}}{{end}}
                    {{end}}
                    </span>
                    <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{$exID}}/delete" class="d-inline">
                        <button type="submit" class="btn btn-link btn-sm text-danger p-0" title="Remove exercise"><i class="bi bi-trash"></i></button>
                    </form>
                </div>
            </div>

            {{if .Sets}}
            {{if .Exercise.IsTimeBased}}
            <table class="table table-sm mb-2">
                <thead><tr>
                    <th class="text-muted fw-normal small ps-0">Set</th>
                    <th class="text-muted fw-normal small">Type</th>
                    <th class="text-muted fw-normal small">Duration</th>
                    <th></th>
                </tr></thead>
                <tbody>
                    {{range .Sets}}
                    <tr>
                        <td class="ps-0">{{.SetNumber}}</td>
                        <td class="text-capitalize">{{if .ActivityType}}{{.ActivityType}}{{else}}&mdash;{{end}}</td>
                        <td>{{fmtDuration .ActualSeconds}}</td>
                        <td class="text-end">
                            <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{$exID}}/sets/{{.ID}}/delete" class="d-inline">
                                <button type="submit" class="btn btn-link btn-sm text-danger p-0" title="Delete set"><i class="bi bi-trash"></i></button>
                            </form>
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{else}}
            <table class="table table-sm mb-2">
                <thead><tr>
                    <th class="text-muted fw-normal small ps-0">Set</th>
                    <th class="text-muted fw-normal small">Weight</th>
                    <th class="text-muted fw-normal small">Reps</th>
                    <th></th>
                </tr></thead>
                <tbody>
                    {{range .Sets}}
                    <tr>
                        <td class="ps-0">{{.SetNumber}}</td>
                        <td>{{.ActualWeight}} {{.WeightUnit}}</td>
                        <td>{{.ActualReps}}</td>
                        <td class="text-end">
                            <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{$exID}}/sets/{{.ID}}/delete" class="d-inline">
                                <button type="submit" class="btn btn-link btn-sm text-danger p-0" title="Delete set"><i class="bi bi-trash"></i></button>
                            </form>
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{end}}
            {{end}}

            {{if .Exercise.IsTimeBased}}
            <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{.Exercise.ID}}/sets" class="d-flex gap-2 align-items-end flex-wrap log-set-form" data-time-based="1" data-goal-seconds="{{.Exercise.GoalSeconds}}">
                <div>
                    <label class="form-label small mb-1">Type</label>
                    <select name="activity_type" class="form-select form-select-sm" style="width: 150px;">
                        <option value="">—</option>
                        <option value="steady state">Steady State</option>
                        <option value="fartlek">Fartlek</option>
                        <option value="intervals">Intervals</option>
                        <option value="hiit">HIIT</option>
                        <option value="easy">Easy / Recovery</option>
                    </select>
                </div>
                <div>
                    <label class="form-label small mb-1">Duration (hrs:mins:secs) </label>
                    <div class="d-flex gap-1 align-items-end">
                        <div class="text-center">
                            <input type="number" name="actual_h" class="form-control form-control-sm text-center" value="0" min="0" step="1" style="width: 56px;">
                        </div>
                        <div class="text-center">
                            <input type="number" name="actual_m" class="form-control form-control-sm text-center" value="0" min="0" max="59" step="1" style="width: 56px;">
                        </div>
                        <div class="text-center">
                            <input type="number" name="actual_s" class="form-control form-control-sm text-center" value="0" min="0" max="59" step="1" style="width: 56px;">
                        </div>
                    </div>
                </div>
                <button type="submit" class="btn btn-dark btn-sm mb-0" style="white-space:nowrap">+ Set</button>
            </form>
            {{else}}
            {{$last := .LastSet}}
            <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{.Exercise.ID}}/sets" class="d-flex gap-2 align-items-end w-100 log-set-form">
                <div class="flex-grow-1">
                    <label class="form-label small mb-1">Weight</label>
                    <div class="input-group input-group-sm">
                        <input type="number" name="actual_weight" class="form-control" placeholder="0" min="0" step="0.5"{{if $last}} value="{{$last.ActualWeight}}"{{else if gt .Exercise.GoalWeight 0.0}} value="{{.Exercise.GoalWeight}}"{{end}}>
                        <select name="weight_unit" class="form-select" style="max-width: 80px;">
                            <option value="lb" {{if $last}}{{if eq $last.WeightUnit "lb"}}selected{{end}}{{else if eq .Exercise.WeightUnit "lb"}}selected{{end}}>lb</option>
                            <option value="kg" {{if $last}}{{if eq $last.WeightUnit "kg"}}selected{{end}}{{else if eq .Exercise.WeightUnit "kg"}}selected{{end}}>kg</option>
                        </select>
                    </div>
                </div>
                <div>
                    <label class="form-label small mb-1">Reps</label>
                    <input type="number" name="actual_reps" class="form-control form-control-sm" placeholder="0" min="1" required style="width: 70px;"{{if $last}} value="{{$last.ActualReps}}"{{else if gt .Exercise.GoalReps 0}} value="{{.Exercise.GoalReps}}"{{end}}>
                </div>
                <button type="submit" class="btn btn-dark btn-sm mb-0" style="white-space:nowrap">+ Set</button>
            </form>
            {{end}}
        </div>
    </div>
    {{end}}
    {{end}}



    {{else}}

    <div class="sortable-block" data-session-id="{{$.Session.ID}}" data-block="{{.Block}}">
    {{range .Exercises}}
    {{$exID := .Exercise.ID}}
    <div class="card mb-3" data-ex-id="{{$exID}}">
        <div class="card-body pb-2">
            <div class="d-flex align-items-center justify-content-between mb-1">
                <div class="d-flex align-items-center gap-2 flex-grow-1 min-w-0">
                    <i class="bi bi-grip-vertical text-muted drag-handle flex-shrink-0" style="font-size:1.1rem;"></i>
                    <h2 class="h6 fw-semibold mb-0 text-capitalize text-truncate">{{if .HitMax}}
                      <button type="button" class="btn btn-link p-0 border-0 hit-max-btn" data-bs-toggle="modal" data-bs-target="#goalWeightModal" data-ex-name="{{.Exercise.Name}}" data-goal-weight="{{.Exercise.GoalWeight}}" data-weight-unit="{{.Exercise.WeightUnit}}" data-direction="up" title="Hit max reps last workout — tap to update goal weight" style="line-height:1;vertical-align:middle;">
                        <i class="bi bi-arrow-up-circle-fill text-black" style="font-size:1.0em;"></i>
                      </button>
                      &nbsp;{{else}}{{if (not .Exercise.IsTimeBased)}}{{if .Exercise.IsBodyweight}}
                      <button type="button" class="btn btn-link p-0 border-0" data-bs-toggle="modal" data-bs-target="#goalRepsModal" data-ex-name="{{.Exercise.Name}}" data-goal-rep-min="{{.GoalRepMin}}" data-goal-rep-max="{{.GoalRepMax}}" title="Set goal reps" style="line-height:1;vertical-align:middle;"><i class="bi bi-pencil text-black" style="font-size:1.0em;"></i></button>{{else}}<button type="button" class="btn btn-link p-0 border-0 hit-max-btn" data-bs-toggle="modal" data-bs-target="#goalWeightModal" data-ex-name="{{.Exercise.Name}}" data-goal-weight="{{.Exercise.GoalWeight}}" data-weight-unit="{{.Exercise.WeightUnit}}" data-direction="down" title="{{if gt .Exercise.GoalWeight 0.0}}Missed max reps — tap to adjust goal weight{{else}}Set goal weight{{end}}" style="line-height:1;vertical-align:middle;"><i class="bi bi-{{if gt .Exercise.GoalWeight 0.0}}dash-circle-fill{{else}}pencil{{end}} text-black" style="font-size:1.0em;"></i></button>{{end}}&nbsp;{{end}}{{end}}{{.Exercise.Name}}</h2>
                </div>
                <div class="flex-shrink-0">
                    <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{$exID}}/delete" class="d-inline">
                        <button type="submit" class="btn btn-link btn-sm text-danger p-0" title="Remove exercise"><i class="bi bi-trash"></i></button>
                    </form>
                </div>
            </div>
            <div class="text-muted small mb-2 ps-1">
            {{if .Exercise.IsTimeBased}}
            {{if gt .Exercise.GoalSeconds 0}}Goal: {{fmtDuration .Exercise.GoalSeconds}}{{end}}
            {{else}}
            {{if gt .Exercise.GoalWeight 0.0}}Goal: {{.Exercise.GoalWeight}} {{.Exercise.WeightUnit}}{{end}}
            {{if and .Exercise.IsBodyweight (gt .GoalRepMax 0)}}
            &nbsp;{{.GoalRepMin}}–{{.GoalRepMax}} reps
            {{else if and (gt $.PhaseRepMin 0) (gt $.PhaseRepMax 0)}}
            &nbsp;{{$.PhaseRepMin}}–{{$.PhaseRepMax}} reps
            {{end}}
            {{end}}
            </div>

            {{if .Sets}}
            {{if .Exercise.IsTimeBased}}
            <table class="table table-sm mb-2">
                <thead>
                    <tr>
                        <th class="text-muted fw-normal small ps-0">Set</th>
                        <th class="text-muted fw-normal small">Type</th>
                        <th class="text-muted fw-normal small">Duration</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Sets}}
                    <tr>
                        <td class="ps-0">{{.SetNumber}}</td>
                        <td class="text-capitalize">{{if .ActivityType}}{{.ActivityType}}{{else}}&mdash;{{end}}</td>
                        <td>{{fmtDuration .ActualSeconds}}</td>
                        <td class="text-end">
                            <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{$exID}}/sets/{{.ID}}/delete" class="d-inline">
                                <button type="submit" class="btn btn-link btn-sm text-danger p-0" title="Delete set"><i class="bi bi-trash"></i></button>
                            </form>
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{else}}
            <table class="table table-sm mb-2">
                <thead>
                    <tr>
                        <th class="text-muted fw-normal small ps-0">Set</th>
                        <th class="text-muted fw-normal small">Weight</th>
                        <th class="text-muted fw-normal small">Reps</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Sets}}
                    <tr>
                        <td class="ps-0">{{.SetNumber}}</td>
                        <td>{{.ActualWeight}} {{.WeightUnit}}</td>
                        <td>{{.ActualReps}}</td>
                        <td class="text-end">
                            <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{$exID}}/sets/{{.ID}}/delete" class="d-inline">
                                <button type="submit" class="btn btn-link btn-sm text-danger p-0" title="Delete set"><i class="bi bi-trash"></i></button>
                            </form>
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{end}}
            {{end}}

            {{$last := .LastSet}}
            {{if .Exercise.IsTimeBased}}
            <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{.Exercise.ID}}/sets" class="d-flex gap-2 align-items-end flex-wrap log-set-form" data-time-based="1" data-goal-seconds="{{.Exercise.GoalSeconds}}">
                <div>
                    <label class="form-label small mb-1">Type</label>
                    <select name="activity_type" class="form-select form-select-sm" style="width: 150px;">
                        <option value="">—</option>
                        <option value="steady state"{{if $last}}{{if eq $last.ActivityType "steady state"}} selected{{end}}{{end}}>Steady State</option>
                        <option value="fartlek"{{if $last}}{{if eq $last.ActivityType "fartlek"}} selected{{end}}{{end}}>Fartlek</option>
                        <option value="intervals"{{if $last}}{{if eq $last.ActivityType "intervals"}} selected{{end}}{{end}}>Intervals</option>
                        <option value="hiit"{{if $last}}{{if eq $last.ActivityType "hiit"}} selected{{end}}{{end}}>HIIT</option>
                        <option value="easy"{{if $last}}{{if eq $last.ActivityType "easy"}} selected{{end}}{{end}}>Easy / Recovery</option>
                    </select>
                </div>
                <div>
                    <label class="form-label small mb-1">Duration (hrs:mins:secs)</label>
                    <div class="d-flex gap-1 align-items-end">
                        <div class="text-center">
                            <input type="number" name="actual_h" class="form-control form-control-sm text-center" value="{{if $last}}{{$last.Hours}}{{else}}0{{end}}" min="0" step="1" style="width: 56px;">
                        </div>
                        <div class="text-center">
                            <input type="number" name="actual_m" class="form-control form-control-sm text-center" value="{{if $last}}{{$last.Minutes}}{{else}}0{{end}}" min="0" max="59" step="1" style="width: 56px;">
                        </div>
                        <div class="text-center">
                            <input type="number" name="actual_s" class="form-control form-control-sm text-center" value="{{if $last}}{{$last.Secs}}{{else}}0{{end}}" min="0" max="59" step="1" style="width: 56px;">
                        </div>
                    </div>
                </div>
                <button type="submit" class="btn btn-dark btn-sm mb-0" style="white-space:nowrap">+ Set</button>
            </form>
            {{else}}
            <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{.Exercise.ID}}/sets" class="d-flex gap-2 align-items-end w-100 log-set-form">
                <div class="flex-grow-1">
                    <label class="form-label small mb-1">Weight</label>
                    <div class="input-group input-group-sm">
                        <input type="number" name="actual_weight" class="form-control" placeholder="0" min="0" step="0.5"{{if $last}} value="{{$last.ActualWeight}}"{{else if gt .Exercise.GoalWeight 0.0}} value="{{.Exercise.GoalWeight}}"{{end}}>
                        <select name="weight_unit" class="form-select" style="max-width: 80px;">
                            <option value="lb" {{if $last}}{{if eq $last.WeightUnit "lb"}}selected{{end}}{{else if eq .Exercise.WeightUnit "lb"}}selected{{end}}>lb</option>
                            <option value="kg" {{if $last}}{{if eq $last.WeightUnit "kg"}}selected{{end}}{{else if eq .Exercise.WeightUnit "kg"}}selected{{end}}>kg</option>
                        </select>
                    </div>
                </div>
                <div>
                    <label class="form-label small mb-1">Reps</label>
                    <input type="number" name="actual_reps" class="form-control form-control-sm" placeholder="0" min="1" required style="width: 70px;"{{if $last}} value="{{$last.ActualReps}}"{{else if gt .Exercise.GoalReps 0}} value="{{.Exercise.GoalReps}}"{{end}}>
                </div>
                <button type="submit" class="btn btn-dark btn-sm mb-0" style="white-space:nowrap">+ Set</button>
            </form>
            {{end}}
        </div>
    </div>
    {{end}}
    </div>{{/* end sortable-block */}}

    {{end}}
    {{end}}

    <h2 class="h6 fw-semibold text-uppercase text-muted mt-4 mb-3">Add Exercise</h2>

    <form id="add-exercise-form" method="POST" action="/sessions/{{.Session.ID}}/exercises">
        <div class="card mb-3 p-3">
            {{template "partials/exercise_fields.tpl" .}}
            <div class="mt-2">
                <label class="form-label">Section</label>
                <select name="block" class="form-select form-select-sm">
                    <option value="main">Main</option>
                    <option value="abs">Abs</option>
                    <option value="cardio">Cardio</option>
                    <option value="stretch">Stretch</option>
                </select>
            </div>
        </div>
        <button type="submit" class="btn btn-dark btn-sm mb-4">Add Exercise</button>
    </form>
</main>

{{template "partials/modal_goal_weight.tpl" .}}
{{template "partials/modal_goal_reps.tpl" .}}
{{template "partials/modal_goal_seconds.tpl" .}}

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/sortablejs@1.15.3/Sortable.min.js"></script>
<script>

function fmtDuration(secs) {
    const h = Math.floor(secs / 3600);
    const m = Math.floor((secs % 3600) / 60);
    const s = secs % 60;
    if (h > 0) return `${h}:${String(m).padStart(2,'0')}:${String(s).padStart(2,'0')}`;
    return `${m}:${String(s).padStart(2,'0')}`;
}

// Pre-fill h/m/s log form inputs from data-goal-seconds attribute.
document.querySelectorAll('.log-set-form[data-time-based="1"]').forEach(form => {
    const goal = parseInt(form.dataset.goalSeconds || '0', 10);
    if (goal > 0) {
        form.querySelector('[name="actual_h"]').value = Math.floor(goal / 3600);
        form.querySelector('[name="actual_m"]').value = Math.floor((goal % 3600) / 60);
        form.querySelector('[name="actual_s"]').value = goal % 60;
    }
});

document.querySelectorAll('.log-set-form').forEach(form => {
    form.addEventListener('submit', async e => {
        e.preventDefault();

        const isTimeBased = form.dataset.timeBased === '1';
        const formData = new FormData(form);

        let data = null;
        try {
            const res = await fetch(form.action, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                    'X-Requested-With': 'XMLHttpRequest',
                },
                body: new URLSearchParams(formData),
            });
            if (!res.ok) { form.submit(); return; }
            data = await res.json();
        } catch {
            form.submit();
            return;
        }

        // Find or create the sets table inside this card.
        const cardBody = form.closest('.card-body');
        let table = cardBody.querySelector('table');
        if (!table) {
            table = document.createElement('table');
            table.className = 'table table-sm mb-2';
            if (isTimeBased) {
                table.innerHTML =
                    '<thead><tr>' +
                    '<th class="text-muted fw-normal small ps-0">Set</th>' +
                    '<th class="text-muted fw-normal small">Type</th>' +
                    '<th class="text-muted fw-normal small">Duration</th>' +
                    '<th></th>' +
                    '</tr></thead><tbody></tbody>';
            } else {
                table.innerHTML =
                    '<thead><tr>' +
                    '<th class="text-muted fw-normal small ps-0">Set</th>' +
                    '<th class="text-muted fw-normal small">Weight</th>' +
                    '<th class="text-muted fw-normal small">Reps</th>' +
                    '<th></th>' +
                    '</tr></thead><tbody></tbody>';
            }
            form.before(table);
        }

        // Extract exercise and session IDs from the form action URL.
        const parts   = form.action.split('/');
        const sessIdx = parts.indexOf('sessions');
        const sessID  = parts[sessIdx + 1];
        const exIdx   = parts.indexOf('exercises');
        const exID    = parts[exIdx + 1];

        const tbody  = table.querySelector('tbody');
        const setNum = data.set_number;
        const setID  = data.id;
        const deleteAction = `/sessions/${sessID}/exercises/${exID}/sets/${setID}/delete`;
        const row    = document.createElement('tr');

        if (isTimeBased) {
            const h = parseInt(formData.get('actual_h') || '0', 10);
            const m = parseInt(formData.get('actual_m') || '0', 10);
            const s = parseInt(formData.get('actual_s') || '0', 10);
            const totalSecs = h * 3600 + m * 60 + s;
            const actType = formData.get('activity_type') || '';
            const actTypeDisplay = actType ? `<span class="text-capitalize">${actType}</span>` : '&mdash;';
            row.innerHTML =
                `<td class="ps-0">${setNum}</td>` +
                `<td>${actTypeDisplay}</td>` +
                `<td>${fmtDuration(totalSecs)}</td>` +
                `<td class="text-end"><form method="POST" action="${deleteAction}" class="d-inline">` +
                `<button type="submit" class="btn btn-link btn-sm text-danger p-0" title="Delete set"><i class="bi bi-trash"></i></button>` +
                `</form></td>`;
        } else {
            const weight = formData.get('actual_weight') || '0';
            const unit   = formData.get('weight_unit')   || 'lb';
            const reps   = formData.get('actual_reps')   || '0';
            row.innerHTML =
                `<td class="ps-0">${setNum}</td>` +
                `<td>${weight} ${unit}</td>` +
                `<td>${reps}</td>` +
                `<td class="text-end"><form method="POST" action="${deleteAction}" class="d-inline">` +
                `<button type="submit" class="btn btn-link btn-sm text-danger p-0" title="Delete set"><i class="bi bi-trash"></i></button>` +
                `</form></td>`;
        }
        tbody.appendChild(row);

        if (typeof window.startRestTimer === 'function') window.startRestTimer();
    });
});
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
<script>
document.querySelectorAll('.sortable-block').forEach(function (container) {
    Sortable.create(container, {
        handle: '.drag-handle',
        animation: 150,
        ghostClass: 'sortable-ghost',
        onEnd: function () {
            const sessionID = container.dataset.sessionId;
            const ids = Array.from(container.querySelectorAll('.card[data-ex-id]'))
                .map(function (el) { return el.dataset.exId; })
                .join(',');
            fetch('/sessions/' + sessionID + '/exercises/reorder', {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body: 'ids=' + encodeURIComponent(ids),
            });
        },
    });
});
</script>

<!-- Rest Timer -->
<div id="rest-timer" class="d-none" style="position:fixed;bottom:0;left:0;right:0;z-index:1050;background:rgba(15,15,15,0.95);color:#fff;border-top:1px solid rgba(255,255,255,0.12);">
    <div class="container py-3 px-3" style="max-width:600px;">
        <div class="d-flex align-items-center justify-content-between gap-3">
            <div>
                <div class="text-secondary small mb-1">Rest</div>
                <div id="rest-countdown" class="fw-bold font-monospace" style="font-size:2rem;line-height:1;letter-spacing:0.04em;">0:00</div>
                <div class="text-secondary small mt-1">Rested: <span id="rest-elapsed">0:00</span></div>
            </div>
            <button id="rest-close" class="btn btn-outline-light btn-sm px-3">Done</button>
        </div>
    </div>
</div>

<script>
(function () {
    const PHASE_REST = {{.PhaseRestSeconds}};
    const SID        = '{{.Session.ID}}';
    const KEY_START  = 'restTimer_' + SID + '_start';
    const KEY_DUR    = 'restTimer_' + SID + '_dur';

    const timerEl      = document.getElementById('rest-timer');
    const countdownEl  = document.getElementById('rest-countdown');
    const elapsedEl    = document.getElementById('rest-elapsed');
    const closeBtn     = document.getElementById('rest-close');
    let   interval     = null;

    function fmt(secs) {
        const s = Math.abs(secs);
        return Math.floor(s / 60) + ':' + String(s % 60).padStart(2, '0');
    }

    let notified = false;

    function notify() {
        if (notified) return;
        notified = true;
        if (Notification.permission !== 'granted') return;
        const opts = {
            body: 'Time to get back to it.',
            icon: '/static/icons/icon-192.png',
            tag:  'rest-timer',
            renotify: false,
        };
        // Mobile browsers require showNotification via the service worker;
        // new Notification() is desktop-only.
        if (navigator.serviceWorker && navigator.serviceWorker.controller) {
            navigator.serviceWorker.ready.then(function (reg) {
                reg.showNotification('Rest complete!', opts);
            });
        } else {
            new Notification('Rest complete!', opts);
        }
    }

    function tick(startMs, durationSecs) {
        const elapsed   = Math.floor((Date.now() - startMs) / 1000);
        const remaining = durationSecs - elapsed;
        if (remaining <= 0 && remaining > -1) notify();
        countdownEl.textContent = remaining > 0 ? fmt(remaining) : 'Done!';
        elapsedEl.textContent   = fmt(elapsed);
    }

    function show(startMs, durationSecs) {
        if (interval) clearInterval(interval);
        notified = false;
        timerEl.classList.remove('d-none');
        tick(startMs, durationSecs);
        interval = setInterval(function () { tick(startMs, durationSecs); }, 500);
    }

    function stop() {
        if (interval) { clearInterval(interval); interval = null; }
        localStorage.removeItem(KEY_START);
        localStorage.removeItem(KEY_DUR);
        timerEl.classList.add('d-none');
    }

    closeBtn.addEventListener('click', stop);

    // Restore on reload
    const savedStart = localStorage.getItem(KEY_START);
    const savedDur   = localStorage.getItem(KEY_DUR);
    if (savedStart && savedDur) {
        show(parseInt(savedStart, 10), parseInt(savedDur, 10));
    }

    window.startRestTimer = function () {
        if (PHASE_REST <= 0) return;
        // Request permission on first use (must be called during a user gesture)
        if ('Notification' in window && Notification.permission === 'default') {
            Notification.requestPermission();
        }
        const now = Date.now();
        localStorage.setItem(KEY_START, now.toString());
        localStorage.setItem(KEY_DUR, PHASE_REST.toString());
        show(now, PHASE_REST);
    };
})();
</script>
<!-- goal modal scripts are in their respective partials -->
</body>
</html>
