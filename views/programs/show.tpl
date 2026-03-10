<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Program.Name}} — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4" style="max-width: 600px;">
    <div class="mb-4">
        <a href="/programs" class="text-muted small">&larr; Programs</a>
        <div class="d-flex align-items-start justify-content-between mt-1">
            <div>
                <h1 class="h4 fw-bold mb-1">{{.Program.Name}}</h1>
                <p class="text-muted small mb-0">
                    Starts {{.Program.StartDate.Format "Jan 2, 2006"}}
                    &middot; {{.Program.NumPhases}} phase{{if gt .Program.NumPhases 1}}s{{end}}
                    &middot; {{.Program.WeeksPerPhase}} week{{if gt .Program.WeeksPerPhase 1}}s{{end}}/phase
                </p>
            </div>
            <form method="POST" action="/programs/{{.Program.ID}}/sessions" class="ms-3 flex-shrink-0 text-end">
                {{if .Templates}}
                <select name="template_id" class="form-select form-select-sm mb-1">
                    <option value="">No template</option>
                    {{range .Templates}}
                    <option value="{{.ID}}">{{.Name}}</option>
                    {{end}}
                </select>
                {{end}}
                <button type="submit" class="btn btn-dark btn-sm w-100">Start Workout</button>
            </form>
        </div>
    </div>

    {{if .Success}}
    <div class="alert alert-success alert-dismissible fade show" id="success-alert">{{.Success}}</div>
    {{end}}

    {{if .Error}}
    <div class="alert alert-danger">{{.Error}}</div>
    {{end}}

    <h2 class="h6 fw-semibold text-uppercase text-muted mb-3">Rep Ranges by Phase</h2>

    <form method="POST" action="/programs/{{.Program.ID}}" data-offline-sync>
        <div class="card">
            <ul class="list-group list-group-flush">
                {{range .Phases}}
                <li class="list-group-item">
                    <div class="d-flex align-items-center gap-3">
                        <span class="fw-semibold" style="min-width: 64px;">Phase {{.PhaseNumber}}</span>
                        <div class="d-flex align-items-center gap-2 flex-grow-1">
                            <input
                                type="number"
                                class="form-control form-control-sm phase-min"
                                name="rep_min_{{.PhaseNumber}}"
                                value="{{if gt .RepMin 0}}{{.RepMin}}{{end}}"
                                placeholder="Min"
                                min="1"
                                required
                                style="max-width: 80px;"
                            >
                            <span class="text-muted">–</span>
                            <input
                                type="number"
                                class="form-control form-control-sm phase-max"
                                name="rep_max_{{.PhaseNumber}}"
                                value="{{if gt .RepMax 0}}{{.RepMax}}{{end}}"
                                placeholder="Max"
                                min="1"
                                required
                                style="max-width: 80px;"
                            >
                            <span class="text-muted small">reps</span>
                        </div>
                        <button type="button" class="btn btn-outline-secondary btn-sm copy-to-all" title="Copy to all phases"><i class="bi bi-copy"></i></button>
                    </div>
                </li>
                {{end}}
            </ul>
        </div>

        <div class="mt-3">
            <button type="submit" class="btn btn-dark btn-sm">Save Rep Ranges</button>
        </div>
    </form>

    <h2 class="h6 fw-semibold text-uppercase text-muted mt-4 mb-3">Workout History</h2>

    {{if .Sessions}}
    <div class="card">
        <ul class="list-group list-group-flush">
            {{range .Sessions}}
            <li class="list-group-item">
                <a href="/sessions/{{.ID}}" class="text-decoration-none d-flex justify-content-between align-items-center">
                    <div>
                        <span class="fw-semibold">Workout #{{.WorkoutNumber}}</span>
                        <span class="text-muted ms-2 small">Phase {{.PhaseNumber}} · Week {{.WeekNumber}}{{if .IsDeload}} · <span class="text-warning">Deload</span>{{end}}</span>
                    </div>
                    <span class="text-muted small">{{.Date.Format "Jan 2, 2006"}}</span>
                </a>
            </li>
            {{end}}
        </ul>
    </div>
    {{else}}
    <p class="text-muted small">No workouts yet. Start one above.</p>
    {{end}}
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
    const alertEl = document.getElementById('success-alert');
    if (alertEl) {
        setTimeout(() => bootstrap.Alert.getOrCreateInstance(alertEl).close(), 3000);
    }

    document.querySelectorAll('.copy-to-all').forEach(btn => {
        btn.addEventListener('click', () => {
            const row = btn.closest('li');
            const min = row.querySelector('.phase-min').value;
            const max = row.querySelector('.phase-max').value;
            document.querySelectorAll('.phase-min').forEach(el => el.value = min);
            document.querySelectorAll('.phase-max').forEach(el => el.value = max);
        });
    });
</script>
<script src="/static/offline-sync.js"></script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
