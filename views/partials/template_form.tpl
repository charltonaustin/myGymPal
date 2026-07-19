{{template "partials/navbar.tpl" .}}

<style>
    /* Whether a row is in a circuit is decided by the container it sits in, not
       by its markup, so the same row partial serves both cases and a row dragged
       or re-rendered into a circuit picks up the right fields with no JS. */
    .row-work { display: none; }
    .circuit-exercises .row-work { display: flex; }
    .circuit-exercises .row-block { display: none; }
    .circuit-card { border-left: 3px solid var(--bs-dark); }
    .circuit-card .exercise-row { background-color: var(--bs-tertiary-bg); }
</style>

<main class="container mt-4 mb-4" style="max-width: 560px;">
    <div class="mb-4">
        <a href="{{.BackURL}}" class="text-muted small">&larr; {{.BackLabel}}</a>
        <h1 class="h4 fw-bold mt-1 mb-0">{{.Heading}}</h1>
    </div>

    {{if .Error}}
    <div class="alert alert-danger" id="error-alert">{{.Error}}</div>
    {{end}}

    <form method="POST" action="{{.FormAction}}" novalidate id="template-form">
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
            {{range .Exercises}}
            {{template "partials/template_exercise_row.tpl" .}}
            {{end}}
        </div>

        <button type="button" id="add-exercise" class="btn btn-outline-secondary btn-sm mb-4">+ Add Exercise</button>

        <h2 class="h6 fw-semibold text-uppercase text-muted mt-4 mb-3">Circuits</h2>

        <div id="circuits-container">
            {{range .Circuits}}
            <div class="circuit-card card mb-3 p-3" data-circuit-index="{{.Index}}">
                <div class="mb-2 d-flex align-items-center gap-2">
                    <input type="text" class="form-control circuit-name" name="circuit_name_{{.Index}}" value="{{.Name}}" placeholder="Circuit name" required>
                    <button type="button" class="btn btn-sm btn-outline-danger remove-circuit flex-shrink-0">Remove</button>
                </div>
                <div class="mb-3 d-flex align-items-center gap-2">
                    <div class="input-group input-group-sm">
                        <span class="input-group-text">Rounds</span>
                        <input type="number" class="form-control circuit-rounds" name="circuit_rounds_{{.Index}}" value="{{.Rounds}}" min="1" step="1">
                    </div>
                    <div class="input-group input-group-sm">
                        <span class="input-group-text">Transition</span>
                        <input type="number" class="form-control circuit-transition" name="circuit_transition_{{.Index}}" value="{{.TransitionSeconds}}" min="0" step="5">
                        <span class="input-group-text">sec</span>
                    </div>
                </div>
                <div class="circuit-exercises">
                    {{range .Exercises}}
                    {{template "partials/template_exercise_row.tpl" .}}
                    {{end}}
                </div>
                <button type="button" class="btn btn-outline-secondary btn-sm add-circuit-exercise align-self-start">+ Add Exercise to Circuit</button>
            </div>
            {{end}}
        </div>

        <button type="button" id="add-circuit" class="btn btn-outline-secondary btn-sm mb-4">+ Add Circuit</button>

        <input type="hidden" name="exercise_count" id="exercise_count" value="{{.ExerciseCount}}">
        <input type="hidden" name="circuit_count" id="circuit_count" value="{{.CircuitCount}}">

        <div class="d-flex gap-2">
            <button type="submit" class="btn btn-dark">{{.SubmitLabel}}</button>
            <a href="{{.BackURL}}" class="btn btn-outline-secondary">Cancel</a>
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

    // Circuits are renumbered before the exercises, because each exercise reads
    // the circuit index off the card it sits in.
    function reindexCircuits() {
        document.querySelectorAll('.circuit-card').forEach((card, i) => {
            card.dataset.circuitIndex = i;
            const name = card.querySelector('.circuit-name');
            if (name) name.name = `circuit_name_${i}`;
            const rounds = card.querySelector('.circuit-rounds');
            if (rounds) rounds.name = `circuit_rounds_${i}`;
            const transition = card.querySelector('.circuit-transition');
            if (transition) transition.name = `circuit_transition_${i}`;
        });
        document.getElementById('circuit_count').value = document.querySelectorAll('.circuit-card').length;
    }

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
            const blockSelect = row.querySelector('.row-block');
            if (blockSelect) blockSelect.name = `block_${i}`;
            // Membership follows the DOM: a row inside a circuit card belongs to
            // that circuit, a row anywhere else is loose. Nothing else records it,
            // so this must run before every submit — and the field has to be
            // renamed as well as re-valued, or row i submits row j's membership.
            const card = row.closest('.circuit-card');
            const circuitIndex = row.querySelector('.ex-circuit-index');
            if (circuitIndex) {
                circuitIndex.name = `circuit_index_${i}`;
                circuitIndex.value = card ? card.dataset.circuitIndex : '-1';
            }

            // Work seconds only mean anything inside a circuit. Zeroing them on the
            // way out keeps a row that was dragged out of one — or whose circuit was
            // removed — from carrying a duration nothing will ever read.
            const work = row.querySelector('.ex-work-seconds');
            if (work) {
                work.name = `work_seconds_${i}`;
                if (!card) work.value = '0';
            }
        });
        document.getElementById('exercise_count').value = document.querySelectorAll('.exercise-row').length;
    }

    function reindexAll() {
        reindexCircuits();
        reindexRows();
    }

    function bindRow(row) {
        row.querySelector('.remove-exercise').addEventListener('click', () => { row.remove(); reindexAll(); });
        const nameInput = row.querySelector('[name^="exercise_name_"]');
        if (nameInput) attachAutocomplete(nameInput, row);
        row.querySelectorAll('input[name^="ex_type_"]').forEach(r => r.addEventListener('change', () => syncHiddens(row)));
    }

    // Mirrors partials/template_exercise_row.tpl. The indices here are
    // provisional: reindexAll rewrites every name before the form is submitted.
    function exerciseRowHTML(i) {
        return `
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
                <select class="form-select form-select-sm flex-grow-1 row-block" name="block_${i}">
                    <option value="main" selected>Main</option>
                    <option value="abs">Abs</option>
                    <option value="cardio">Cardio</option>
                    <option value="stretch">Stretch</option>
                </select>
                <div class="input-group input-group-sm row-work flex-grow-1">
                    <input type="number" class="form-control ex-work-seconds" name="work_seconds_${i}" value="0" min="0" step="5" placeholder="Work" aria-label="Work seconds">
                    <span class="input-group-text">sec</span>
                </div>
                <button type="button" class="btn btn-sm btn-outline-danger remove-exercise flex-shrink-0">Remove</button>
            </div>
            <input type="hidden" name="circuit_index_${i}" class="ex-circuit-index" value="-1">
        `;
    }

    function newExerciseRow() {
        const i = document.querySelectorAll('.exercise-row').length;
        const row = document.createElement('div');
        row.className = 'exercise-row card mb-3 p-3';
        row.dataset.index = i;
        row.innerHTML = exerciseRowHTML(i);
        bindRow(row);
        return row;
    }

    function bindCircuit(card) {
        card.querySelector('.remove-circuit').addEventListener('click', () => { card.remove(); reindexAll(); });
        card.querySelector('.add-circuit-exercise').addEventListener('click', () => {
            const row = newExerciseRow();
            // A circuit exercise is worked for a length of time by definition, so it
            // opens on a usable duration rather than on zero.
            row.querySelector('.ex-work-seconds').value = '30';
            card.querySelector('.circuit-exercises').appendChild(row);
            reindexAll();
        });
        makeSortable(card.querySelector('.circuit-exercises'));
    }

    function makeSortable(container) {
        Sortable.create(container, {
            handle: '.drag-handle',
            animation: 150,
            ghostClass: 'sortable-ghost',
            onEnd: reindexAll,
        });
    }

    document.querySelectorAll('.exercise-row').forEach(bindRow);
    document.querySelectorAll('.circuit-card').forEach(bindCircuit);

    document.getElementById('add-exercise').addEventListener('click', () => {
        document.getElementById('exercises-container').appendChild(newExerciseRow());
        reindexAll();
    });

    document.getElementById('add-circuit').addEventListener('click', () => {
        const i = document.querySelectorAll('.circuit-card').length;
        const card = document.createElement('div');
        card.className = 'circuit-card card mb-3 p-3';
        card.dataset.circuitIndex = i;
        card.innerHTML = `
            <div class="mb-2 d-flex align-items-center gap-2">
                <input type="text" class="form-control circuit-name" name="circuit_name_${i}" placeholder="Circuit name" required>
                <button type="button" class="btn btn-sm btn-outline-danger remove-circuit flex-shrink-0">Remove</button>
            </div>
            <div class="mb-3 d-flex align-items-center gap-2">
                <div class="input-group input-group-sm">
                    <span class="input-group-text">Rounds</span>
                    <input type="number" class="form-control circuit-rounds" name="circuit_rounds_${i}" value="1" min="1" step="1">
                </div>
                <div class="input-group input-group-sm">
                    <span class="input-group-text">Transition</span>
                    <input type="number" class="form-control circuit-transition" name="circuit_transition_${i}" value="15" min="0" step="5">
                    <span class="input-group-text">sec</span>
                </div>
            </div>
            <div class="circuit-exercises"></div>
            <button type="button" class="btn btn-outline-secondary btn-sm add-circuit-exercise align-self-start">+ Add Exercise to Circuit</button>
        `;
        document.getElementById('circuits-container').appendChild(card);
        bindCircuit(card);
        reindexAll();
    });

    document.getElementById('template-form').addEventListener('submit', function (e) {
        reindexAll();
        if (!this.checkValidity()) {
            e.preventDefault();
            e.stopPropagation();
        }
        this.classList.add('was-validated');
    });

    makeSortable(document.getElementById('exercises-container'));

    // The edit page arrives with rows already server-rendered; the new page starts
    // empty. Running this once on load settles both into the same state, and in
    // particular writes each row's circuit membership into its hidden field.
    reindexAll();
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
