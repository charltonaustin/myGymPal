{{/*
    One exercise row on the template form. The dot is an exerciseForm.

    The same row is used for a loose exercise and for one inside a circuit; which
    it is depends on the container it sits in, not on the markup. The work-seconds
    field and the block select are shown or hidden by CSS on that container, so a
    row keeps working when it is added, removed, or re-rendered on either page.
*/}}
<div class="exercise-row card mb-3 p-3" data-index="{{.Index}}">
    <div class="mb-2 d-flex align-items-center gap-2">
        <i class="bi bi-grip-vertical text-muted drag-handle flex-shrink-0" style="font-size:1.1rem;"></i>
        <input type="text" class="form-control" name="exercise_name_{{.Index}}" value="{{.Name}}" placeholder="Exercise name" required>
    </div>
    <div class="mb-2">
        <div class="btn-group w-100" role="group">
            <input type="radio" class="btn-check" name="ex_type_{{.Index}}" id="ex_weighted_{{.Index}}" value="weighted" autocomplete="off" {{if and (not .IsBodyweight) (not .IsTimeBased)}}checked{{end}}>
            <label class="btn btn-outline-secondary btn-sm" for="ex_weighted_{{.Index}}">Weighted</label>
            <input type="radio" class="btn-check" name="ex_type_{{.Index}}" id="ex_bw_{{.Index}}" value="bodyweight" autocomplete="off" {{if .IsBodyweight}}checked{{end}}>
            <label class="btn btn-outline-secondary btn-sm" for="ex_bw_{{.Index}}">Bodyweight</label>
            <input type="radio" class="btn-check" name="ex_type_{{.Index}}" id="ex_tb_{{.Index}}" value="time_based" autocomplete="off" {{if .IsTimeBased}}checked{{end}}>
            <label class="btn btn-outline-secondary btn-sm" for="ex_tb_{{.Index}}">Time-based</label>
        </div>
        <input type="hidden" name="is_bodyweight_{{.Index}}" class="ex-bw-hidden" value="{{if .IsBodyweight}}on{{end}}">
        <input type="hidden" name="is_time_based_{{.Index}}" class="ex-tb-hidden" value="{{if .IsTimeBased}}on{{end}}">
    </div>
    <div class="d-flex align-items-center gap-2">
        <select class="form-select form-select-sm flex-grow-1 row-block" name="block_{{.Index}}">
            <option value="main" {{if or (eq .Block "") (eq .Block "main")}}selected{{end}}>Main</option>
            <option value="abs" {{if eq .Block "abs"}}selected{{end}}>Abs</option>
            <option value="cardio" {{if eq .Block "cardio"}}selected{{end}}>Cardio</option>
            <option value="stretch" {{if eq .Block "stretch"}}selected{{end}}>Stretch</option>
        </select>
        <div class="input-group input-group-sm row-work flex-grow-1">
            <input type="number" class="form-control ex-work-seconds" name="work_seconds_{{.Index}}" value="{{.WorkSeconds}}" min="0" step="5" placeholder="Work" aria-label="Work seconds">
            <span class="input-group-text">sec</span>
        </div>
        <button type="button" class="btn btn-sm btn-outline-danger remove-exercise flex-shrink-0">Remove</button>
    </div>
    <input type="hidden" name="circuit_index_{{.Index}}" class="ex-circuit-index" value="-1">
</div>
