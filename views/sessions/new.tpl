<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Start Workout — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4 mb-4" style="max-width: 480px;">
    <div class="mb-4">
        <a href="/programs/{{.Program.ID}}" class="text-muted small">&larr; {{.Program.Name}}</a>
        <h1 class="h4 fw-bold mt-1 mb-0">{{if .LogMode}}Log Workout{{else}}Start Workout{{end}}</h1>
    </div>

    <form method="POST" action="/programs/{{.Program.ID}}/sessions">
        <div class="card mb-3">
            <div class="card-body">
                <div class="row g-3">
                    <div class="col-4">
                        <label class="form-label fw-semibold">Phase</label>
                        <input
                            type="number"
                            class="form-control"
                            name="phase_number"
                            value="{{.PhaseNumber}}"
                            min="1"
                            required
                        >
                    </div>
                    <div class="col-4">
                        <label class="form-label fw-semibold">Week</label>
                        <input
                            type="number"
                            class="form-control"
                            name="week_number"
                            id="week_number"
                            value="{{.WeekNumber}}"
                            min="1"
                            max="{{.Program.WeeksPerPhase}}"
                            required
                        >
                    </div>
                    <div class="col-4">
                        <label class="form-label fw-semibold">Workout #</label>
                        <input
                            type="number"
                            class="form-control"
                            name="workout_number"
                            value="{{.WorkoutNumber}}"
                            min="1"
                            required
                        >
                    </div>
                </div>
                <div id="deload-notice" class="mt-2 text-warning small fw-semibold" style="display:none;">
                    Deload week — reps will be set 2 below phase minimum.
                </div>
            </div>
        </div>

        <div class="mb-3">
            <label class="form-label fw-semibold">Date</label>
            <input
                type="date"
                class="form-control"
                name="date"
                value="{{.DefaultDate}}"
                required
            >
        </div>

        {{if .Templates}}
        <div class="mb-3">
            <label class="form-label fw-semibold">Template</label>
            <select name="template_id" class="form-select">
                <option value="">No template</option>
                {{range .Templates}}
                <option value="{{.ID}}">{{.Name}}</option>
                {{end}}
            </select>
        </div>
        {{end}}

        <button type="submit" class="btn btn-dark w-100">{{if .LogMode}}Log Workout{{else}}Start Workout{{end}}</button>
    </form>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
    const weeksPerPhase = {{.Program.WeeksPerPhase}};
    const weekInput = document.getElementById('week_number');
    const notice = document.getElementById('deload-notice');

    function updateDeloadNotice() {
        notice.style.display = parseInt(weekInput.value) === weeksPerPhase ? '' : 'none';
    }

    weekInput.addEventListener('input', updateDeloadNotice);
    updateDeloadNotice();
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
