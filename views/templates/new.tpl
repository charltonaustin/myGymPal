<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Template — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
    <style>.drag-handle { cursor: grab; touch-action: none; } .sortable-ghost { opacity: 0.4; }</style>
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4 mb-4" style="max-width: 560px;">
    <div class="mb-4">
        <a href="/templates" class="text-muted small">&larr; Templates</a>
        <h1 class="h4 fw-bold mt-1 mb-0">New Workout Template</h1>
    </div>

    {{if .Error}}
    <div class="alert alert-danger" id="error-alert">{{.Error}}</div>
    {{end}}

    <form method="POST" action="/templates/new" novalidate id="template-form">
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
                <div class="mb-2 d-flex align-items-center gap-2">
                    <i class="bi bi-grip-vertical text-muted drag-handle flex-shrink-0" style="font-size:1.1rem;"></i>
                    <input type="text" class="form-control" name="exercise_name_{{$i}}" value="{{$ex.Name}}" placeholder="Exercise name" required>
                </div>
                <div class="mb-2">
                    <div class="btn-group w-100" role="group">
                        <input type="radio" class="btn-check" name="ex_type_{{$i}}" id="ex_weighted_{{$i}}" value="weighted" autocomplete="off" {{if and (not $ex.IsBodyweight) (not $ex.IsTimeBased)}}checked{{end}}>
                        <label class="btn btn-outline-secondary btn-sm" for="ex_weighted_{{$i}}">Weighted</label>
                        <input type="radio" class="btn-check" name="ex_type_{{$i}}" id="ex_bw_{{$i}}" value="bodyweight" autocomplete="off" {{if $ex.IsBodyweight}}checked{{end}}>
                        <label class="btn btn-outline-secondary btn-sm" for="ex_bw_{{$i}}">Bodyweight</label>
                        <input type="radio" class="btn-check" name="ex_type_{{$i}}" id="ex_tb_{{$i}}" value="time_based" autocomplete="off" {{if $ex.IsTimeBased}}checked{{end}}>
                        <label class="btn btn-outline-secondary btn-sm" for="ex_tb_{{$i}}">Time-based</label>
                    </div>
                    <input type="hidden" name="is_bodyweight_{{$i}}" class="ex-bw-hidden" value="{{if $ex.IsBodyweight}}on{{end}}">
                    <input type="hidden" name="is_time_based_{{$i}}" class="ex-tb-hidden" value="{{if $ex.IsTimeBased}}on{{end}}">
                </div>
                <div class="d-flex align-items-center gap-2">
                    <select class="form-select form-select-sm flex-grow-1" name="block_{{$i}}">
                        <option value="main" {{if or (eq $ex.Block "") (eq $ex.Block "main")}}selected{{end}}>Main</option>
                        <option value="abs" {{if eq $ex.Block "abs"}}selected{{end}}>Abs</option>
                        <option value="cardio" {{if eq $ex.Block "cardio"}}selected{{end}}>Cardio</option>
                        <option value="stretch" {{if eq $ex.Block "stretch"}}selected{{end}}>Stretch</option>
                    </select>
                    <button type="button" class="btn btn-sm btn-outline-danger remove-exercise flex-shrink-0">Remove</button>
                </div>
            </div>
            {{end}}
        </div>

        <input type="hidden" name="exercise_count" id="exercise_count" value="{{len .Exercises}}">

        <button type="button" id="add-exercise" class="btn btn-outline-secondary btn-sm mb-4">+ Add Exercise</button>

        <div class="d-flex gap-2">
            <button type="submit" class="btn btn-dark">Create Template</button>
            <a href="/templates" class="btn btn-outline-secondary">Cancel</a>
        </div>
    </form>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/sortablejs@1.15.2/Sortable.min.js"></script>
<script>
    const alertEl = document.getElementById('error-alert');
    if (alertEl) setTimeout(() => alertEl.remove(), 4000);

    // Exercise library for autocomplete
    const exerciseLibraryArr = {{.ExerciseLibraryJSON}};
    const exerciseLibrary = {};
    exerciseLibraryArr.forEach(ex => { exerciseLibrary[ex.name] = ex; });

    function autofillFromLibrary(row, nameVal) {
        const entry = exerciseLibrary[nameVal];
        if (!entry) return;
        const type = entry.isTimeBased ? 'time_based' : entry.isBodyweight ? 'bodyweight' : 'weighted';
        const radio = row.querySelector(`input[value="${type}"]`);
        if (radio) { radio.checked = true; syncHiddens(row); }
    }

    function syncHiddens(row) {
        const checked = row.querySelector('input[name^="ex_type_"]:checked');
        const val = checked ? checked.value : 'weighted';
        const bwHidden = row.querySelector('.ex-bw-hidden');
        const tbHidden = row.querySelector('.ex-tb-hidden');
        if (bwHidden) bwHidden.value = val === 'bodyweight' ? 'on' : '';
        if (tbHidden) tbHidden.value = val === 'time_based' ? 'on' : '';
    }

    function attachAutocomplete(input, row) {
        const wrapper = input.parentElement;
        wrapper.style.position = 'relative';
        let dropdown = null;

        function showDropdown(matches) {
            if (!dropdown) {
                dropdown = document.createElement('div');
                dropdown.className = 'list-group shadow';
                dropdown.style.cssText = 'position:absolute;top:100%;left:0;right:0;z-index:1050;max-height:220px;overflow-y:auto;border-radius:0 0 .375rem .375rem;';
                wrapper.appendChild(dropdown);
            }
            dropdown.innerHTML = '';
            matches.forEach(ex => {
                const btn = document.createElement('button');
                btn.type = 'button';
                btn.className = 'list-group-item list-group-item-action py-2 px-3';
                btn.style.fontSize = '0.95rem';
                btn.textContent = ex.name;
                btn.addEventListener('mousedown', e => {
                    e.preventDefault();
                    input.value = ex.name;
                    autofillFromLibrary(row, ex.name);
                    hideDropdown();
                });
                dropdown.appendChild(btn);
            });
        }

        function hideDropdown() {
            if (dropdown) { dropdown.remove(); dropdown = null; }
        }

        input.addEventListener('input', () => {
            const val = input.value.toLowerCase().trim();
            if (!val) { hideDropdown(); return; }
            const matches = exerciseLibraryArr.filter(ex => ex.name.includes(val)).slice(0, 10);
            if (matches.length) showDropdown(matches);
            else hideDropdown();
        });

        input.addEventListener('blur', () => setTimeout(hideDropdown, 150));
    }

    let exerciseCount = parseInt(document.getElementById('exercise_count').value, 10);

    function reindexRows() {
        document.querySelectorAll('.exercise-row').forEach((row, i) => {
            row.dataset.index = i;
            const nameInput = row.querySelector('[name^="exercise_name_"]');
            if (nameInput) nameInput.name = `exercise_name_${i}`;
            row.querySelectorAll('input[type="radio"]').forEach(radio => {
                const oldId = radio.id;
                const newId = oldId.replace(/_\d+$/, `_${i}`);
                const label = row.querySelector(`label[for="${oldId}"]`);
                radio.name = `ex_type_${i}`;
                radio.id = newId;
                if (label) label.htmlFor = newId;
            });
            const bwHidden = row.querySelector('.ex-bw-hidden');
            if (bwHidden) bwHidden.name = `is_bodyweight_${i}`;
            const tbHidden = row.querySelector('.ex-tb-hidden');
            if (tbHidden) tbHidden.name = `is_time_based_${i}`;
            const blockSelect = row.querySelector('[name^="block_"]');
            if (blockSelect) blockSelect.name = `block_${i}`;
        });
        document.getElementById('exercise_count').value = document.querySelectorAll('.exercise-row').length;
    }

    function bindRow(row) {
        row.querySelector('.remove-exercise').addEventListener('click', () => { row.remove(); reindexRows(); });
        const nameInput = row.querySelector('[name^="exercise_name_"]');
        if (nameInput) attachAutocomplete(nameInput, row);
        row.querySelectorAll('input[name^="ex_type_"]').forEach(r => r.addEventListener('change', () => syncHiddens(row)));
    }

    document.querySelectorAll('.exercise-row').forEach(bindRow);

    document.getElementById('add-exercise').addEventListener('click', () => {
        const i = exerciseCount++;

        const row = document.createElement('div');
        row.className = 'exercise-row card mb-3 p-3';
        row.dataset.index = i;
        row.innerHTML = `
            <div class="mb-2 d-flex align-items-center gap-2">
                <i class="bi bi-grip-vertical text-muted drag-handle flex-shrink-0" style="font-size:1.1rem;"></i>
                <input type="text" class="form-control" name="exercise_name_${i}" placeholder="Exercise name" required>
            </div>
            <div class="mb-2">
                <div class="btn-group w-100" role="group">
                    <input type="radio" class="btn-check" name="ex_type_${i}" id="ex_weighted_${i}" value="weighted" autocomplete="off" checked>
                    <label class="btn btn-outline-secondary btn-sm" for="ex_weighted_${i}">Weighted</label>
                    <input type="radio" class="btn-check" name="ex_type_${i}" id="ex_bw_${i}" value="bodyweight" autocomplete="off">
                    <label class="btn btn-outline-secondary btn-sm" for="ex_bw_${i}">Bodyweight</label>
                    <input type="radio" class="btn-check" name="ex_type_${i}" id="ex_tb_${i}" value="time_based" autocomplete="off">
                    <label class="btn btn-outline-secondary btn-sm" for="ex_tb_${i}">Time-based</label>
                </div>
                <input type="hidden" name="is_bodyweight_${i}" class="ex-bw-hidden" value="">
                <input type="hidden" name="is_time_based_${i}" class="ex-tb-hidden" value="">
            </div>
            <div class="d-flex align-items-center gap-2">
                <select class="form-select form-select-sm flex-grow-1" name="block_${i}">
                    <option value="main" selected>Main</option>
                    <option value="abs">Abs</option>
                    <option value="cardio">Cardio</option>
                    <option value="stretch">Stretch</option>
                </select>
                <button type="button" class="btn btn-sm btn-outline-danger remove-exercise flex-shrink-0">Remove</button>
            </div>
        `;
        document.getElementById('exercises-container').appendChild(row);
        bindRow(row);
    });

    document.getElementById('template-form').addEventListener('submit', function (e) {
        reindexRows();
        if (!this.checkValidity()) {
            e.preventDefault();
            e.stopPropagation();
        }
        this.classList.add('was-validated');
    });

    Sortable.create(document.getElementById('exercises-container'), {
        handle: '.drag-handle',
        animation: 150,
        ghostClass: 'sortable-ghost',
        onEnd: reindexRows,
    });
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
