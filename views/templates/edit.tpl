<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Edit {{.Template.Name}} — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4" style="max-width: 560px;">
    <div class="mb-4">
        <a href="/templates/{{.Template.ID}}" class="text-muted small">&larr; {{.Template.Name}}</a>
        <h1 class="h4 fw-bold mt-1 mb-0">Edit Template</h1>
    </div>

    {{if .Error}}
    <div class="alert alert-danger" id="error-alert">{{.Error}}</div>
    {{end}}

    <form method="POST" action="/templates/{{.Template.ID}}" novalidate id="template-form">
        <div class="mb-3">
            <label for="name" class="form-label">Template Name</label>
            <input
                type="text"
                class="form-control"
                id="name"
                name="name"
                value="{{.Name}}"
                placeholder="e.g. Upper Body A"
                required
            >
            <div class="invalid-feedback">Template name is required.</div>
        </div>

        <div class="mb-3">
            <label for="focus" class="form-label">Focus <span class="text-muted fw-normal">(optional)</span></label>
            <input
                type="text"
                class="form-control"
                id="focus"
                name="focus"
                value="{{.Focus}}"
                placeholder="e.g. Chest &amp; Shoulders"
            >
        </div>

        <h2 class="h6 fw-semibold text-uppercase text-muted mt-4 mb-3">Exercises</h2>

        <div id="exercises-container">
            {{range $i, $ex := .Exercises}}
            <div class="exercise-row card mb-3 p-3" data-index="{{$i}}">
                <div class="mb-2">
                    <input
                        type="text"
                        class="form-control"
                        name="exercise_name_{{$i}}"
                        value="{{$ex.Name}}"
                        placeholder="Exercise name"
                        required
                    >
                </div>
                <div class="form-check mb-2">
                    <input
                        type="checkbox"
                        class="form-check-input bw-check"
                        name="is_bodyweight_{{$i}}"
                        id="bw_{{$i}}"
                        {{if $ex.IsBodyweight}}checked{{end}}
                    >
                    <label class="form-check-label" for="bw_{{$i}}">Bodyweight exercise</label>
                </div>
                <div class="weight-row mb-2 {{if $ex.IsBodyweight}}d-none{{end}}">
                    <div class="input-group input-group-sm">
                        <input
                            type="number"
                            class="form-control"
                            name="goal_weight_{{$i}}"
                            value="{{$ex.GoalWeight}}"
                            placeholder="Goal weight"
                            min="0"
                            step="0.5"
                        >
                        <select name="weight_unit_{{$i}}" class="form-select" style="max-width: 72px;">
                            <option value="lb" {{if eq $ex.WeightUnit "lb"}}selected{{end}}>lb</option>
                            <option value="kg" {{if eq $ex.WeightUnit "kg"}}selected{{end}}>kg</option>
                        </select>
                    </div>
                </div>
                <div class="d-flex align-items-center gap-2">
                    <input
                        type="number"
                        class="form-control form-control-sm"
                        name="rep_min_{{$i}}"
                        value="{{$ex.RepMin}}"
                        placeholder="Min"
                        min="1"
                        required
                        style="max-width: 90px;"
                    >
                    <span class="text-muted">–</span>
                    <input
                        type="number"
                        class="form-control form-control-sm"
                        name="rep_max_{{$i}}"
                        value="{{$ex.RepMax}}"
                        placeholder="Max"
                        min="1"
                        required
                        style="max-width: 90px;"
                    >
                    <span class="text-muted small">reps</span>
                    <button type="button" class="btn btn-sm btn-outline-danger ms-auto remove-exercise">Remove</button>
                </div>
            </div>
            {{end}}
        </div>

        <input type="hidden" name="exercise_count" id="exercise_count" value="{{len .Exercises}}">

        <button type="button" id="add-exercise" class="btn btn-outline-secondary btn-sm mb-4">+ Add Exercise</button>

        <div class="d-flex gap-2">
            <button type="submit" class="btn btn-dark">Save Changes</button>
            <a href="/templates/{{.Template.ID}}" class="btn btn-outline-secondary">Cancel</a>
        </div>
    </form>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
    const alertEl = document.getElementById('error-alert');
    if (alertEl) setTimeout(() => alertEl.remove(), 4000);

    const weightUnit = "{{.WeightUnit}}";
    let exerciseCount = parseInt(document.getElementById('exercise_count').value, 10);

    function bindRow(row) {
        const bwCheck = row.querySelector('.bw-check');
        const weightRow = row.querySelector('.weight-row');
        bwCheck.addEventListener('change', () => {
            weightRow.classList.toggle('d-none', bwCheck.checked);
        });
        row.querySelector('.remove-exercise').addEventListener('click', () => {
            row.remove();
        });
    }

    document.querySelectorAll('.exercise-row').forEach(bindRow);

    document.getElementById('add-exercise').addEventListener('click', () => {
        const i = exerciseCount++;
        document.getElementById('exercise_count').value = exerciseCount;

        const row = document.createElement('div');
        row.className = 'exercise-row card mb-3 p-3';
        row.dataset.index = i;
        row.innerHTML = `
            <div class="mb-2">
                <input type="text" class="form-control" name="exercise_name_${i}" placeholder="Exercise name" required>
            </div>
            <div class="form-check mb-2">
                <input type="checkbox" class="form-check-input bw-check" name="is_bodyweight_${i}" id="bw_${i}">
                <label class="form-check-label" for="bw_${i}">Bodyweight exercise</label>
            </div>
            <div class="weight-row mb-2">
                <div class="input-group input-group-sm">
                    <input type="number" class="form-control" name="goal_weight_${i}" placeholder="Goal weight" min="0" step="0.5">
                    <select name="weight_unit_${i}" class="form-select" style="max-width: 72px;">
                        <option value="lb" ${weightUnit === 'lb' ? 'selected' : ''}>lb</option>
                        <option value="kg" ${weightUnit === 'kg' ? 'selected' : ''}>kg</option>
                    </select>
                </div>
            </div>
            <div class="d-flex align-items-center gap-2">
                <input type="number" class="form-control form-control-sm" name="rep_min_${i}" placeholder="Min" min="1" required style="max-width: 90px;">
                <span class="text-muted">–</span>
                <input type="number" class="form-control form-control-sm" name="rep_max_${i}" placeholder="Max" min="1" required style="max-width: 90px;">
                <span class="text-muted small">reps</span>
                <button type="button" class="btn btn-sm btn-outline-danger ms-auto remove-exercise">Remove</button>
            </div>
        `;
        document.getElementById('exercises-container').appendChild(row);
        bindRow(row);
    });

    document.getElementById('template-form').addEventListener('submit', function (e) {
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
