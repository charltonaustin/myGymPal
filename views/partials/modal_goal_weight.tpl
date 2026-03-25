<!-- Goal weight update modal -->
<div class="modal fade" id="goalWeightModal" tabindex="-1" aria-labelledby="goalWeightModalLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="goalWeightModalLabel">Update Goal Weight</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
                <p id="goalWeightDesc" class="text-muted small mb-3">You hit the max reps last session. Set a new goal weight for <strong id="goalWeightExName"></strong>.</p>
                <div class="input-group">
                    <input type="number" id="goalWeightInput" class="form-control" min="0" step="0.5" placeholder="New goal weight">
                    <select id="goalWeightUnit" class="form-select" style="max-width: 90px;">
                        <option value="lb">lb</option>
                        <option value="kg">kg</option>
                    </select>
                </div>
                <div id="goalWeightError" class="text-danger small mt-2" style="display:none;"></div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="button" class="btn btn-dark btn-sm" id="goalWeightSaveBtn">Save</button>
            </div>
        </div>
    </div>
</div>

<script>
(function () {
    let activeBtn = null;

    document.getElementById('goalWeightModal').addEventListener('show.bs.modal', function (e) {
        activeBtn = e.relatedTarget;
        const name       = activeBtn.dataset.exName;
        const goalWeight = parseFloat(activeBtn.dataset.goalWeight) || 0;
        const unit       = activeBtn.dataset.weightUnit || 'lb';
        const direction  = activeBtn.dataset.direction || 'up';

        document.getElementById('goalWeightExName').textContent = name;
        document.getElementById('goalWeightInput').value = goalWeight > 0 ? goalWeight : '';
        document.getElementById('goalWeightUnit').value  = unit === 'kg' ? 'kg' : 'lb';
        document.getElementById('goalWeightError').style.display = 'none';

        const desc = document.getElementById('goalWeightDesc');
        if (direction === 'down') {
            desc.firstChild.textContent = 'Set a new goal weight for ';
        } else {
            desc.firstChild.textContent = 'You hit the max reps last session. Set a new goal weight for ';
        }
    });

    document.getElementById('goalWeightSaveBtn').addEventListener('click', async function () {
        if (!activeBtn) return;
        const name       = activeBtn.dataset.exName;
        const goalWeight = document.getElementById('goalWeightInput').value;
        const unit       = document.getElementById('goalWeightUnit').value;
        const errEl      = document.getElementById('goalWeightError');
        errEl.style.display = 'none';

        try {
            const res = await fetch('/exercises/goal-weight', {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body: new URLSearchParams({ name, goal_weight: goalWeight, weight_unit: unit }),
            });
            const data = await res.json();
            if (data.error) {
                errEl.textContent = data.error;
                errEl.style.display = '';
                return;
            }
            // Update the goal text displayed on the card.
            const card = activeBtn.closest('.card');
            if (card) {
                const goalSpan = card.querySelector('.text-muted.small');
                if (goalSpan) {
                    const goalLabel = `Goal: ${data.goal_weight} ${data.weight_unit}`;
                    if (/Goal:\s*[\d.]+ (?:lb|kg)/.test(goalSpan.textContent)) {
                        goalSpan.textContent = goalSpan.textContent.replace(/Goal:\s*[\d.]+ (?:lb|kg)/, goalLabel);
                    } else {
                        // No goal text yet — prepend it.
                        const existing = goalSpan.textContent.trim();
                        goalSpan.textContent = existing ? `${goalLabel}  ${existing}` : goalLabel;
                    }
                }
                // Update the button's data attributes for the next modal open.
                activeBtn.dataset.goalWeight = data.goal_weight;
                activeBtn.dataset.weightUnit = data.weight_unit;
                // Swap pencil icon to dash-circle-fill now that a goal weight is set.
                const icon = activeBtn.querySelector('i');
                if (icon) {
                    icon.classList.replace('bi-pencil', 'bi-dash-circle-fill');
                }
                // Update the weight input if it's empty or showing the old default of 0.
                const weightInput = card.querySelector('input[name="actual_weight"]');
                if (weightInput && (!weightInput.value || weightInput.value === '0')) {
                    weightInput.value = data.goal_weight;
                }
                const unitSelect = card.querySelector('select[name="weight_unit"]');
                if (unitSelect) {
                    unitSelect.value = data.weight_unit;
                }
            }
            bootstrap.Modal.getInstance(document.getElementById('goalWeightModal')).hide();
        } catch {
            errEl.textContent = 'Something went wrong. Please try again.';
            errEl.style.display = '';
        }
    });
})();
</script>
