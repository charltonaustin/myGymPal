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

<main class="container mt-4 mb-4" style="max-width: 600px;">
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
            <div class="d-flex gap-2 ms-3 flex-shrink-0">
                <a href="/programs/{{.Program.ID}}/sessions/new" class="btn btn-dark btn-sm">Start Workout</a>
                <a href="/programs/{{.Program.ID}}/sessions/new?sequential=1" class="btn btn-outline-secondary btn-sm">Log Workout</a>
            </div>
        </div>
    </div>

    {{if .Success}}
    <div class="alert alert-success alert-dismissible fade show" id="success-alert">{{.Success}}</div>
    {{end}}

    {{if .Error}}
    <div class="alert alert-danger">{{.Error}}</div>
    {{end}}

    <h2 class="h6 fw-semibold text-uppercase text-muted mt-4 mb-3">Workout History</h2>

    {{if .Sessions}}
    <div class="card">
        <ul class="list-group list-group-flush">
            {{range .Sessions}}
            <li class="list-group-item list-group-item-action d-flex align-items-center gap-2">
                <a href="/sessions/{{.ID}}" class="text-decoration-none text-dark flex-grow-1 d-flex justify-content-between align-items-center">
                    <div>
                        <span class="fw-semibold">Workout #{{.WorkoutNumber}}</span>
                        <span class="text-muted ms-2 small">Phase {{.PhaseNumber}} · Week {{.WeekNumber}}{{if .IsDeload}} · <span class="text-warning">Deload</span>{{end}}</span>
                    </div>
                    <span class="text-muted small d-flex align-items-center gap-2">{{.Date.Format "Jan 2, 2006"}} </span>
                </a>
                <form method="POST" action="/sessions/{{.ID}}/delete" class="flex-shrink-0">
                    <button type="submit" class="btn btn-link text-danger p-0" title="Delete workout" onclick="return confirm('Delete Workout #{{.WorkoutNumber}}?')">
                        <i class="bi bi-trash"></i>
                    </button>
                </form>
            </li>
            {{end}}
        </ul>
    </div>
    {{else}}
    <p class="text-muted small">No workouts yet. Start one above.</p>
    {{end}}

        <h2 class="h6 fw-semibold text-uppercase text-muted mt-5 mb-3">Phase Settings</h2>

        <form method="POST" action="/programs/{{.Program.ID}}" data-offline-sync>
            <div class="card">
                <ul class="list-group list-group-flush">
                    {{range .Phases}}
                    <li class="list-group-item py-2">
                        <div class="d-flex align-items-center justify-content-between mb-2">
                            <span class="fw-semibold small">Phase {{.PhaseNumber}}</span>
                            <button type="button" class="btn btn-outline-secondary btn-sm copy-to-all" title="Copy to all phases"><i class="bi bi-copy"></i></button>
                        </div>
                        <div class="d-flex align-items-center gap-3 flex-wrap">
                            <div class="d-flex align-items-center gap-1">
                                <span class="text-muted small">Reps</span>
                                <input type="number" class="form-control form-control-sm phase-min text-center" name="rep_min_{{.PhaseNumber}}" value="{{if gt .RepMin 0}}{{.RepMin}}{{end}}" placeholder="–" min="1" required style="width:56px;">
                                <span class="text-muted small">–</span>
                                <input type="number" class="form-control form-control-sm phase-max text-center" name="rep_max_{{.PhaseNumber}}" value="{{if gt .RepMax 0}}{{.RepMax}}{{end}}" placeholder="–" min="1" required style="width:56px;">
                            </div>
                            <div class="d-flex align-items-center gap-1">
                                <span class="text-muted small">Sets</span>
                                <input type="number" class="form-control form-control-sm phase-sets text-center" name="sets_{{.PhaseNumber}}" value="{{if gt .DefaultSets 0}}{{.DefaultSets}}{{else}}3{{end}}" min="1" required style="width:56px;">
                            </div>
                            <div class="d-flex align-items-center gap-1">
                                <span class="text-muted small">Rest</span>
                                <input type="number" class="form-control form-control-sm phase-rest-m text-center" name="rest_m_{{.PhaseNumber}}" value="{{.RestSeconds | restMinutes}}" placeholder="0" min="0" style="width:52px;">
                                <span class="text-muted small">m</span>
                                <input type="number" class="form-control form-control-sm phase-rest-s text-center" name="rest_s_{{.PhaseNumber}}" value="{{.RestSeconds | restSecs}}" placeholder="0" min="0" max="59" style="width:52px;">
                                <span class="text-muted small">s</span>
                            </div>
                        </div>
                    </li>
                    {{end}}
                </ul>
            </div>

            <div class="mt-3 mb-3">
                <button type="submit" class="btn btn-dark btn-sm">Save Phase Settings</button>
            </div>
        </form>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
    const alertEl = document.getElementById('success-alert');
    if (alertEl) {
        setTimeout(() => bootstrap.Alert.getOrCreateInstance(alertEl).close(), 3000);
    }

    document.querySelectorAll('.copy-to-all').forEach(btn => {
        btn.addEventListener('click', () => {
            const row   = btn.closest('li');
            const min   = row.querySelector('.phase-min').value;
            const max   = row.querySelector('.phase-max').value;
            const sets  = row.querySelector('.phase-sets').value;
            const restM = row.querySelector('.phase-rest-m').value;
            const restS = row.querySelector('.phase-rest-s').value;
            document.querySelectorAll('.phase-min').forEach(el => el.value = min);
            document.querySelectorAll('.phase-max').forEach(el => el.value = max);
            document.querySelectorAll('.phase-sets').forEach(el => el.value = sets);
            document.querySelectorAll('.phase-rest-m').forEach(el => el.value = restM);
            document.querySelectorAll('.phase-rest-s').forEach(el => el.value = restS);
        });
    });
</script>
<script src="/static/offline-sync.js"></script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
