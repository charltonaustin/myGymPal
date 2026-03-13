<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Session #{{.Session.WorkoutNumber}} — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4" style="max-width: 600px;">
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
                <h2 class="h6 fw-semibold mb-0 text-capitalize">{{.Exercise.Name}}</h2>
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
                <button type="submit" class="btn btn-dark btn-sm mb-0">+ Set</button>
            </form>
            {{else}}
            <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{.Exercise.ID}}/sets" class="d-flex gap-2 align-items-end log-set-form">
                <div>
                    <label class="form-label small mb-1">Weight</label>
                    <div class="input-group input-group-sm" style="width: 160px;">
                        <input type="number" name="actual_weight" class="form-control" placeholder="0" min="0" step="0.5"{{if gt .Exercise.GoalWeight 0.0}} value="{{.Exercise.GoalWeight}}"{{end}}>
                        <select name="weight_unit" class="form-select" style="max-width: 90px;">
                            <option value="lb" {{if eq $.WeightUnit "lb"}}selected{{end}}>lb</option>
                            <option value="kg" {{if eq $.WeightUnit "kg"}}selected{{end}}>kg</option>
                        </select>
                    </div>
                </div>
                <div>
                    <label class="form-label small mb-1">Reps</label>
                    <input type="number" name="actual_reps" class="form-control form-control-sm" placeholder="0" min="1" required style="width: 70px;"{{if gt .Exercise.GoalReps 0}} value="{{.Exercise.GoalReps}}"{{end}}>
                </div>
                <button type="submit" class="btn btn-dark btn-sm mb-0">+ Set</button>
            </form>
            {{end}}
        </div>
    </div>
    {{end}}
    {{end}}

    <div class="card mb-3 p-3">
        <form method="POST" action="/sessions/{{$.Session.ID}}/cardio" class="d-flex flex-column gap-2">
            <input type="text" name="name" class="form-control form-control-sm" placeholder="Activity (e.g. run, bike) — optional">
            <div class="d-flex gap-2 align-items-end flex-wrap">
                <div>
                    <label class="form-label small mb-1">Type</label>
                    <select name="cardio_type" class="form-select form-select-sm" style="width: 150px;">
                        <option value="">—</option>
                        <option value="steady state">Steady State</option>
                        <option value="fartlek">Fartlek</option>
                        <option value="intervals">Intervals</option>
                        <option value="hiit">HIIT</option>
                        <option value="easy">Easy / Recovery</option>
                    </select>
                </div>
                <div>
                    <label class="form-label small mb-1">Goal (min)</label>
                    <input type="number" name="goal_duration" class="form-control form-control-sm" placeholder="0" min="0" style="width: 80px;">
                </div>
                <div>
                    <label class="form-label small mb-1">Actual (min)</label>
                    <input type="number" name="actual_duration" class="form-control form-control-sm" placeholder="0" min="0" required style="width: 80px;">
                </div>
                <button type="submit" class="btn btn-dark btn-sm mb-0">Log</button>
            </div>
        </form>
    </div>

    {{else}}

    {{range .Exercises}}
    {{$exID := .Exercise.ID}}
    <div class="card mb-3">
        <div class="card-body pb-2">
            <div class="d-flex align-items-baseline justify-content-between mb-2">
                <h2 class="h6 fw-semibold mb-0 text-capitalize">{{.Exercise.Name}}</h2>
                <div class="d-flex align-items-center gap-2">
                    <span class="text-muted small">
                    {{if .Exercise.IsTimeBased}}
                    {{if gt .Exercise.GoalSeconds 0}}Goal: {{fmtDuration .Exercise.GoalSeconds}}{{end}}
                    {{else}}
                    {{if gt .Exercise.GoalWeight 0.0}}
                    Goal: {{.Exercise.GoalWeight}} {{.Exercise.WeightUnit}}
                    {{end}}
                    {{if and (gt $.PhaseRepMin 0) (gt $.PhaseRepMax 0)}}
                    {{$.PhaseRepMin}}–{{$.PhaseRepMax}} reps
                    {{end}}
                    {{end}}
                    </span>
                    <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{$exID}}/delete" class="d-inline">
                        <button type="submit" class="btn btn-link btn-sm text-danger p-0" title="Remove exercise"><i class="bi bi-trash"></i></i></button>
                    </form>
                </div>
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
                    <label class="form-label small mb-1">Duration (hrs:mins:secs)</label>
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
                <button type="submit" class="btn btn-dark btn-sm mb-0">+ Set</button>
            </form>
            {{else}}
            <form method="POST" action="/sessions/{{$.Session.ID}}/exercises/{{.Exercise.ID}}/sets" class="d-flex gap-2 align-items-end log-set-form">
                <div>
                    <label class="form-label small mb-1">Weight</label>
                    <div class="input-group input-group-sm" style="width: 160px;">
                        <input type="number" name="actual_weight" class="form-control" placeholder="0" min="0" step="0.5"{{if gt .Exercise.GoalWeight 0.0}} value="{{.Exercise.GoalWeight}}"{{end}}>
                        <select name="weight_unit" class="form-select" style="max-width: 90px;">
                            <option value="lb" {{if eq $.WeightUnit "lb"}}selected{{end}}>lb</option>
                            <option value="kg" {{if eq $.WeightUnit "kg"}}selected{{end}}>kg</option>
                        </select>
                    </div>
                </div>
                <div>
                    <label class="form-label small mb-1">Reps</label>
                    <input type="number" name="actual_reps" class="form-control form-control-sm" placeholder="0" min="1" required style="width: 70px;"{{if gt .Exercise.GoalReps 0}} value="{{.Exercise.GoalReps}}"{{end}}>
                </div>
                <button type="submit" class="btn btn-dark btn-sm mb-0">+ Set</button>
            </form>
            {{end}}
        </div>
    </div>
    {{end}}

    {{end}}
    {{end}}

    <h2 class="h6 fw-semibold text-uppercase text-muted mt-4 mb-3">Add Exercise</h2>
    <form method="POST" action="/sessions/{{.Session.ID}}/exercises">
        <div class="card mb-3 p-3">
            <div class="mb-2">
                <input type="text" name="name" class="form-control" placeholder="Exercise name" required>
            </div>
            <div class="d-flex align-items-center gap-3 mb-2">
                <div class="form-check mb-0">
                    <input type="checkbox" class="form-check-input add-ex-bw-check" name="is_bodyweight" id="add_ex_bw">
                    <label class="form-check-label" for="add_ex_bw">Bodyweight</label>
                </div>
                <div class="form-check mb-0">
                    <input type="checkbox" class="form-check-input add-ex-tb-check" name="is_time_based" id="add_ex_tb">
                    <label class="form-check-label" for="add_ex_tb">Time-based</label>
                </div>
                <select name="block" class="form-select form-select-sm" style="width: auto;">
                    <option value="main">Main</option>
                    <option value="abs">Abs</option>
                    <option value="cardio">Cardio</option>
                    <option value="stretch">Stretch</option>
                </select>
            </div>
            <div class="add-ex-weight-row">
                <div class="input-group input-group-sm">
                    <input type="number" name="goal_weight" class="form-control" placeholder="Goal weight" min="0" step="0.5">
                    <select name="weight_unit" class="form-select" style="max-width: 72px;">
                        <option value="lb" {{if eq .WeightUnit "lb"}}selected{{end}}>lb</option>
                        <option value="kg" {{if eq .WeightUnit "kg"}}selected{{end}}>kg</option>
                    </select>
                </div>
            </div>
            <div class="add-ex-time-row d-none mt-2">
                <input type="number" name="goal_seconds" class="form-control form-control-sm" placeholder="Goal duration (sec)" min="0" step="1" style="width: 180px;">
            </div>
        </div>
        <button type="submit" class="btn btn-dark btn-sm mb-4">Add Exercise</button>
    </form>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
const addExBwCheck = document.querySelector('.add-ex-bw-check');
const addExTbCheck = document.querySelector('.add-ex-tb-check');
const addExWeightRow = document.querySelector('.add-ex-weight-row');
const addExTimeRow = document.querySelector('.add-ex-time-row');

function updateAddExRows() {
    addExWeightRow.classList.toggle('d-none', addExBwCheck.checked || addExTbCheck.checked);
    addExTimeRow.classList.toggle('d-none', !addExTbCheck.checked);
}

if (addExBwCheck) {
    addExBwCheck.addEventListener('change', updateAddExRows);
    addExTbCheck.addEventListener('change', updateAddExRows);
}

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
    });
});
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
