<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Program — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4 mb-4" style="max-width: 480px;">
    <h1 class="h4 fw-bold mb-4">New Training Program</h1>

    {{if .Error}}
    <div class="alert alert-danger alert-dismissible fade show" id="error-alert">{{.Error}}</div>
    {{end}}

    <form method="POST" action="/programs" novalidate id="program-form" data-offline-sync>
        <div class="mb-3">
            <label for="name" class="form-label">Program Name</label>
            <input
                type="text"
                class="form-control"
                id="name"
                name="name"
                value="{{.Name}}"
                placeholder="e.g. Hypertrophy Block 1"
                required
            >
            <div class="invalid-feedback">Program name is required.</div>
        </div>

        <div class="mb-3">
            <label for="start_date" class="form-label">Start Date</label>
            <input
                type="date"
                class="form-control"
                id="start_date"
                name="start_date"
                value="{{.StartDate}}"
                required
            >
            <div class="invalid-feedback">Start date is required.</div>
        </div>

        <div class="mb-3">
            <label for="num_phases" class="form-label">Number of Phases</label>
            <input
                type="number"
                class="form-control"
                id="num_phases"
                name="num_phases"
                value="{{.NumPhases}}"
                min="1"
                required
            >
            <div class="invalid-feedback">Enter a number of phases (at least 1).</div>
        </div>

        <div class="mb-3">
            <label for="weeks_per_phase" class="form-label">Weeks per Phase</label>
            <input
                type="number"
                class="form-control"
                id="weeks_per_phase"
                name="weeks_per_phase"
                value="{{.WeeksPerPhase}}"
                min="1"
                required
            >
            <div class="invalid-feedback">Enter a number of weeks (at least 1).</div>
        </div>

        <div class="mb-3">
            <label for="workouts_per_week" class="form-label">Workouts per Week</label>
            <input
                type="number"
                class="form-control"
                id="workouts_per_week"
                name="workouts_per_week"
                value="{{.WorkoutsPerWeek}}"
                min="1"
                required
            >
            <div class="invalid-feedback">Enter a number of workouts (at least 1).</div>
        </div>

        <div class="mb-3">
            <label class="form-label">Default Rep Range</label>
            <div class="d-flex align-items-center gap-2">
                <input
                    type="number"
                    class="form-control"
                    id="default_rep_min"
                    name="default_rep_min"
                    value="{{.DefaultRepMin}}"
                    placeholder="Min"
                    min="1"
                    required
                    style="max-width: 100px;"
                >
                <span class="text-muted">–</span>
                <input
                    type="number"
                    class="form-control"
                    id="default_rep_max"
                    name="default_rep_max"
                    value="{{.DefaultRepMax}}"
                    placeholder="Max"
                    min="1"
                    required
                    style="max-width: 100px;"
                >
                <span class="text-muted small">reps</span>
            </div>
            <div class="form-text">Applied to all phases — adjust per phase after creating.</div>
            <div class="invalid-feedback">Enter a valid rep range (min ≥ 1, max ≥ min).</div>
        </div>

        <div class="mb-3">
            <label for="default_sets" class="form-label">Default Sets per Exercise</label>
            <input
                type="number"
                class="form-control"
                id="default_sets"
                name="default_sets"
                value="{{.DefaultSets}}"
                min="1"
                required
                style="max-width: 100px;"
            >
            <div class="form-text">Applied to all phases — adjust per phase after creating.</div>
            <div class="invalid-feedback">Enter a number of sets (at least 1).</div>
        </div>

        <div class="d-flex gap-2">
            <button type="submit" class="btn btn-dark">Create Program</button>
            <a href="/programs" class="btn btn-outline-secondary">Cancel</a>
        </div>
    </form>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
    const alertEl = document.getElementById('error-alert');
    if (alertEl) {
        setTimeout(() => bootstrap.Alert.getOrCreateInstance(alertEl).close(), 3000);
    }

    document.getElementById('program-form').addEventListener('submit', function (e) {
        if (!this.checkValidity()) {
            e.preventDefault();
            e.stopPropagation();
        }
        this.classList.add('was-validated');
    });
</script>
<script src="/static/offline-sync.js"></script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
