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

    {{range .Exercises}}
    {{$exID := .Exercise.ID}}
    <div class="card mb-3">
        <div class="card-body pb-2">
            <div class="d-flex align-items-baseline justify-content-between mb-2">
                <h2 class="h6 fw-semibold mb-0">{{.Exercise.Name}}</h2>
                <span class="text-muted small">
                {{if gt .Exercise.GoalWeight 0.0}}
                Goal: {{.Exercise.GoalWeight}} {{.Exercise.WeightUnit}}
                {{end}}
                {{if and (gt $.PhaseRepMin 0) (gt $.PhaseRepMax 0)}}
                {{$.PhaseRepMin}}–{{$.PhaseRepMax}} reps
                {{end}}
                </span>
            </div>

            {{if .Sets}}
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
        </div>
    </div>
    {{end}}

    <div class="card mt-4">
        <div class="card-body">
            <h2 class="h6 fw-semibold mb-3">Add Exercise</h2>
            <form method="POST" action="/sessions/{{.Session.ID}}/exercises">
                <div class="mb-2">
                    <input type="text" name="name" class="form-control form-control-sm" placeholder="Exercise name" required>
                </div>
                <div class="d-flex gap-2 align-items-end">
                    <div>
                        <label class="form-label small mb-1">Goal Weight</label>
                        <div class="input-group input-group-sm" style="width: 150px;">
                            <input type="number" name="goal_weight" class="form-control" placeholder="0" min="0" step="0.5">
                            <select name="weight_unit" class="form-select" style="max-width: 60px;">
                                <option value="lb" {{if eq .WeightUnit "lb"}}selected{{end}}>lb</option>
                                <option value="kg" {{if eq .WeightUnit "kg"}}selected{{end}}>kg</option>
                            </select>
                        </div>
                    </div>
                    <button type="submit" class="btn btn-dark btn-sm">Add</button>
                </div>
            </form>
        </div>
    </div>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
document.querySelectorAll('.log-set-form').forEach(form => {
    form.addEventListener('submit', async e => {
        e.preventDefault();

        const formData = new FormData(form);
        const weight = formData.get('actual_weight') || '0';
        const unit   = formData.get('weight_unit')   || 'lb';
        const reps   = formData.get('actual_reps')   || '0';

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
            table.innerHTML =
                '<thead><tr>' +
                '<th class="text-muted fw-normal small ps-0">Set</th>' +
                '<th class="text-muted fw-normal small">Weight</th>' +
                '<th class="text-muted fw-normal small">Reps</th>' +
                '<th></th>' +
                '</tr></thead><tbody></tbody>';
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
        row.innerHTML =
            `<td class="ps-0">${setNum}</td>` +
            `<td>${weight} ${unit}</td>` +
            `<td>${reps}</td>` +
            `<td class="text-end"><form method="POST" action="${deleteAction}" class="d-inline">` +
            `<button type="submit" class="btn btn-link btn-sm text-danger p-0" title="Delete set"><i class="bi bi-trash"></i></button>` +
            `</form></td>`;
        tbody.appendChild(row);
    });
});
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
