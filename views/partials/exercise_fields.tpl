<div class="mb-3">
    <label for="ex_name" class="form-label">Exercise Name</label>
    <input type="text" class="form-control" id="ex_name" name="name" value="{{.Name}}" placeholder="e.g. Bench Press" required>
</div>

<div class="mb-3">
    <label class="form-label">Type</label>
    <div class="btn-group w-100" role="group">
        <input type="radio" class="btn-check" name="ex_type_radio" id="ex_type_weighted" value="weighted" autocomplete="off" {{if and (not .IsBodyweight) (not .IsTimeBased)}}checked{{end}}>
        <label class="btn btn-outline-secondary btn-sm" for="ex_type_weighted">Weighted</label>
        <input type="radio" class="btn-check" name="ex_type_radio" id="ex_type_bodyweight" value="bodyweight" autocomplete="off" {{if .IsBodyweight}}checked{{end}}>
        <label class="btn btn-outline-secondary btn-sm" for="ex_type_bodyweight">Bodyweight</label>
        <input type="radio" class="btn-check" name="ex_type_radio" id="ex_type_timebased" value="time_based" autocomplete="off" {{if .IsTimeBased}}checked{{end}}>
        <label class="btn btn-outline-secondary btn-sm" for="ex_type_timebased">Time-based</label>
    </div>
    <input type="hidden" name="is_bodyweight" id="ex_is_bw" value="{{if .IsBodyweight}}on{{end}}">
    <input type="hidden" name="is_time_based" id="ex_is_tb" value="{{if .IsTimeBased}}on{{end}}">
</div>

<div class="ex-weight-row {{if or .IsBodyweight .IsTimeBased}}d-none{{end}}">
    <label class="form-label">Goal Weight</label>
    <div class="input-group input-group-sm">
        <input type="number" class="form-control" name="goal_weight" value="{{.GoalWeight}}" placeholder="0" min="0" step="0.5">
        <select name="weight_unit" class="form-select" style="max-width: 72px;">
            <option value="lb" {{if eq .WeightUnit "lb"}}selected{{end}}>lb</option>
            <option value="kg" {{if eq .WeightUnit "kg"}}selected{{end}}>kg</option>
        </select>
    </div>
</div>

<div class="ex-time-row {{if not .IsTimeBased}}d-none{{end}}">
    <label class="form-label">Goal Duration</label>
    <div class="d-flex align-items-center gap-1">
        <input type="number" name="goal_h" class="form-control form-control-sm text-center flex-grow-1" value="{{.GoalHours}}" min="0" step="1" placeholder="0" style="min-width: 0;">
        <span class="text-muted small">h</span>
        <span class="text-muted px-1">:</span>
        <input type="number" name="goal_m" class="form-control form-control-sm text-center flex-grow-1" value="{{.GoalMinutes}}" min="0" max="59" step="1" placeholder="00" style="min-width: 0;">
        <span class="text-muted small">m</span>
        <span class="text-muted px-1">:</span>
        <input type="number" name="goal_s" class="form-control form-control-sm text-center flex-grow-1" value="{{.GoalSecsRemainder}}" min="0" max="59" step="1" placeholder="00" style="min-width: 0;">
        <span class="text-muted small">s</span>
    </div>
</div>

<script>
(function () {
    const radios    = document.querySelectorAll('input[name="ex_type_radio"]');
    const isBwInput = document.getElementById('ex_is_bw');
    const isTbInput = document.getElementById('ex_is_tb');
    const weightRow = document.querySelector('.ex-weight-row');
    const timeRow   = document.querySelector('.ex-time-row');
    if (!radios.length) return;
    function updateRows() {
        const val = document.querySelector('input[name="ex_type_radio"]:checked').value;
        isBwInput.value = val === 'bodyweight' ? 'on' : '';
        isTbInput.value = val === 'time_based' ? 'on' : '';
        weightRow.classList.toggle('d-none', val !== 'weighted');
        timeRow.classList.toggle('d-none', val !== 'time_based');
    }
    radios.forEach(r => r.addEventListener('change', updateRows));
})();
</script>
