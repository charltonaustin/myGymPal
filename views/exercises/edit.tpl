<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Edit Exercise — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4 mb-4" style="max-width: 560px;">
    <div class="mb-4">
        <a href="/exercises" class="text-muted small">&larr; Exercise Library</a>
        <h1 class="h4 fw-bold mt-1 mb-0">Edit Exercise</h1>
    </div>

    {{if .Error}}
    <div class="alert alert-danger" id="error-alert">{{.Error}}</div>
    {{end}}

    <form method="POST" action="/exercises/{{.Exercise.ID}}/edit" novalidate id="exercise-form">
        <div class="card p-3">
            <div class="mb-3">
                <label for="name" class="form-label">Exercise Name</label>
                <input
                    type="text"
                    class="form-control"
                    id="name"
                    name="name"
                    value="{{.Name}}"
                    placeholder="e.g. Bench Press"
                    required
                >
                <div class="invalid-feedback">Exercise name is required.</div>
            </div>

            <div class="form-check mb-2">
                <input
                    type="checkbox"
                    class="form-check-input"
                    id="is_bodyweight"
                    name="is_bodyweight"
                    {{if .IsBodyweight}}checked{{end}}
                >
                <label class="form-check-label" for="is_bodyweight">Bodyweight exercise</label>
            </div>

            <div class="form-check mb-3">
                <input
                    type="checkbox"
                    class="form-check-input"
                    id="is_time_based"
                    name="is_time_based"
                    {{if .IsTimeBased}}checked{{end}}
                >
                <label class="form-check-label" for="is_time_based">Time-based exercise</label>
            </div>

            <div class="weight-row mb-3 {{if or .IsBodyweight .IsTimeBased}}d-none{{end}}">
                <label class="form-label">Goal Weight</label>
                <div class="input-group input-group-sm">
                    <input
                        type="number"
                        class="form-control"
                        name="goal_weight"
                        id="goal_weight"
                        value="{{.GoalWeight}}"
                        placeholder="0"
                        min="0"
                        step="0.5"
                    >
                    <select name="weight_unit" id="weight_unit" class="form-select" style="max-width: 72px;">
                        <option value="lb" {{if eq .ExWeightUnit "lb"}}selected{{end}}>lb</option>
                        <option value="kg" {{if eq .ExWeightUnit "kg"}}selected{{end}}>kg</option>
                    </select>
                </div>
            </div>

            <div class="time-row mb-3 {{if not .IsTimeBased}}d-none{{end}}">
                <label class="form-label">Goal Duration</label>
                <div class="d-flex gap-2 align-items-center">
                    <div class="text-center">
                        <input type="number" name="goal_h" class="form-control form-control-sm text-center" value="{{.GoalHours}}" min="0" step="1" style="width: 64px;">
                        <div class="text-muted" style="font-size: 0.75rem;">hrs</div>
                    </div>
                    <div class="text-center">
                        <input type="number" name="goal_m" class="form-control form-control-sm text-center" value="{{.GoalMinutes}}" min="0" max="59" step="1" style="width: 64px;">
                        <div class="text-muted" style="font-size: 0.75rem;">min</div>
                    </div>
                    <div class="text-center">
                        <input type="number" name="goal_s" class="form-control form-control-sm text-center" value="{{.GoalSecsRemainder}}" min="0" max="59" step="1" style="width: 64px;">
                        <div class="text-muted" style="font-size: 0.75rem;">sec</div>
                    </div>
                </div>
            </div>
        </div>

        <div class="d-flex gap-2 mt-3">
            <button type="submit" class="btn btn-dark">Save Changes</button>
            <a href="/exercises" class="btn btn-outline-secondary">Cancel</a>
        </div>
    </form>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
    const alertEl = document.getElementById('error-alert');
    if (alertEl) setTimeout(() => alertEl.remove(), 4000);

    const bwCheck = document.getElementById('is_bodyweight');
    const tbCheck = document.getElementById('is_time_based');
    const weightRow = document.querySelector('.weight-row');
    const timeRow = document.querySelector('.time-row');

    function updateGoalRows() {
        weightRow.classList.toggle('d-none', bwCheck.checked || tbCheck.checked);
        timeRow.classList.toggle('d-none', !tbCheck.checked);
    }

    bwCheck.addEventListener('change', updateGoalRows);
    tbCheck.addEventListener('change', updateGoalRows);

    document.getElementById('exercise-form').addEventListener('submit', function (e) {
        if (!this.checkValidity()) {
            e.preventDefault();
            e.stopPropagation();
        }
        this.classList.add('was-validated');
    });
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
