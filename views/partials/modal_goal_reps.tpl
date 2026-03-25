<!-- Goal reps update modal -->
<div class="modal fade" id="goalRepsModal" tabindex="-1" aria-labelledby="goalRepsModalLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="goalRepsModalLabel">Update Goal Reps</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
                <p class="text-muted small mb-3">Set a new goal rep range for <strong id="goalRepsExName"></strong>.</p>
                <div class="d-flex gap-2 align-items-center">
                    <div class="flex-grow-1">
                        <label class="form-label small text-muted mb-1">Min reps</label>
                        <input type="number" id="goalRepMinInput" class="form-control" min="0" step="1" placeholder="Min">
                    </div>
                    <div class="pt-3 text-muted">–</div>
                    <div class="flex-grow-1">
                        <label class="form-label small text-muted mb-1">Max reps</label>
                        <input type="number" id="goalRepMaxInput" class="form-control" min="0" step="1" placeholder="Max">
                    </div>
                </div>
                <div id="goalRepsError" class="text-danger small mt-2" style="display:none;"></div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="button" class="btn btn-dark btn-sm" id="goalRepsSaveBtn">Save</button>
            </div>
        </div>
    </div>
</div>

<script>
(function () {
    let activeRepsBtn = null;

    document.getElementById('goalRepsModal').addEventListener('show.bs.modal', function (e) {
        activeRepsBtn = e.relatedTarget;
        const name      = activeRepsBtn.dataset.exName;
        const repMin    = parseInt(activeRepsBtn.dataset.goalRepMin, 10) || 0;
        const repMax    = parseInt(activeRepsBtn.dataset.goalRepMax, 10) || 0;

        document.getElementById('goalRepsExName').textContent = name;
        document.getElementById('goalRepMinInput').value = repMin > 0 ? repMin : '';
        document.getElementById('goalRepMaxInput').value = repMax > 0 ? repMax : '';
        document.getElementById('goalRepsError').style.display = 'none';
    });

    document.getElementById('goalRepsSaveBtn').addEventListener('click', async function () {
        if (!activeRepsBtn) return;
        const name    = activeRepsBtn.dataset.exName;
        const repMin  = document.getElementById('goalRepMinInput').value;
        const repMax  = document.getElementById('goalRepMaxInput').value;
        const errEl   = document.getElementById('goalRepsError');
        errEl.style.display = 'none';

        try {
            const res = await fetch('/exercises/goal-reps', {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body: new URLSearchParams({ name, goal_rep_min: repMin, goal_rep_max: repMax }),
            });
            const data = await res.json();
            if (data.error) {
                errEl.textContent = data.error;
                errEl.style.display = '';
                return;
            }
            // Update the rep range text displayed on the card.
            const card = activeRepsBtn.closest('.card');
            if (card) {
                const goalDiv = card.querySelector('.text-muted.small');
                if (goalDiv) {
                    goalDiv.textContent = goalDiv.textContent.replace(/\d+–\d+ reps/, `${data.goal_rep_min}–${data.goal_rep_max} reps`);
                }
                // Update the button's data attributes for next open.
                activeRepsBtn.dataset.goalRepMin = data.goal_rep_min;
                activeRepsBtn.dataset.goalRepMax = data.goal_rep_max;
            }
            bootstrap.Modal.getInstance(document.getElementById('goalRepsModal')).hide();
        } catch {
            errEl.textContent = 'Something went wrong. Please try again.';
            errEl.style.display = '';
        }
    });
})();
</script>
