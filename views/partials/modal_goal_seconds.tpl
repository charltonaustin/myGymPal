<!-- Goal seconds update modal -->
<div class="modal fade" id="goalSecondsModal" tabindex="-1" aria-labelledby="goalSecondsModalLabel" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title" id="goalSecondsModalLabel">Update Goal Duration</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
                <p class="text-muted small mb-3">Set a new goal duration for <strong id="goalSecondsExName"></strong>.</p>
                <div class="d-flex gap-2 align-items-center">
                    <div class="flex-grow-1">
                        <label class="form-label small text-muted mb-1">Hours</label>
                        <input type="number" id="goalSecondsH" class="form-control text-center" min="0" step="1" placeholder="0">
                    </div>
                    <div class="pt-3 text-muted">:</div>
                    <div class="flex-grow-1">
                        <label class="form-label small text-muted mb-1">Minutes</label>
                        <input type="number" id="goalSecondsM" class="form-control text-center" min="0" max="59" step="1" placeholder="00">
                    </div>
                    <div class="pt-3 text-muted">:</div>
                    <div class="flex-grow-1">
                        <label class="form-label small text-muted mb-1">Seconds</label>
                        <input type="number" id="goalSecondsS" class="form-control text-center" min="0" max="59" step="1" placeholder="00">
                    </div>
                </div>
                <div id="goalSecondsError" class="text-danger small mt-2" style="display:none;"></div>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <button type="button" class="btn btn-dark btn-sm" id="goalSecondsSaveBtn">Save</button>
            </div>
        </div>
    </div>
</div>

<script>
(function () {
    let activeSecsBtn = null;

    function secsToHMS(total) {
        return {
            h: Math.floor(total / 3600),
            m: Math.floor((total % 3600) / 60),
            s: total % 60,
        };
    }

    function fmtHMS(h, m, s) {
        if (h > 0) return `${h}h ${String(m).padStart(2,'0')}m ${String(s).padStart(2,'0')}s`;
        if (m > 0) return `${m}m ${String(s).padStart(2,'0')}s`;
        return `${s}s`;
    }

    document.getElementById('goalSecondsModal').addEventListener('show.bs.modal', function (e) {
        activeSecsBtn = e.relatedTarget;
        const name  = activeSecsBtn.dataset.exName;
        const total = parseInt(activeSecsBtn.dataset.goalSeconds, 10) || 0;
        const {h, m, s} = secsToHMS(total);

        document.getElementById('goalSecondsExName').textContent = name;
        document.getElementById('goalSecondsH').value = h || '';
        document.getElementById('goalSecondsM').value = m || '';
        document.getElementById('goalSecondsS').value = s || '';
        document.getElementById('goalSecondsError').style.display = 'none';
    });

    document.getElementById('goalSecondsSaveBtn').addEventListener('click', async function () {
        if (!activeSecsBtn) return;
        const name  = activeSecsBtn.dataset.exName;
        const h     = document.getElementById('goalSecondsH').value || 0;
        const m     = document.getElementById('goalSecondsM').value || 0;
        const s     = document.getElementById('goalSecondsS').value || 0;
        const errEl = document.getElementById('goalSecondsError');
        errEl.style.display = 'none';

        try {
            const res = await fetch('/exercises/goal-seconds', {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body: new URLSearchParams({ name, goal_h: h, goal_m: m, goal_s: s }),
            });
            const data = await res.json();
            if (data.error) {
                errEl.textContent = data.error;
                errEl.style.display = '';
                return;
            }
            // Update the goal text displayed on the card.
            const card = activeSecsBtn.closest('.card');
            if (card) {
                const goalDiv = card.querySelector('.text-muted.small');
                if (goalDiv) {
                    const {h: rh, m: rm, s: rs} = secsToHMS(data.goal_seconds);
                    const newLabel = `Goal: ${fmtHMS(rh, rm, rs)}`;
                    if (/Goal:/.test(goalDiv.textContent)) {
                        goalDiv.textContent = goalDiv.textContent.replace(/Goal:\s*[\dhms ]+/, newLabel);
                    } else {
                        goalDiv.textContent = newLabel;
                    }
                }
                activeSecsBtn.dataset.goalSeconds = data.goal_seconds;
            }
            bootstrap.Modal.getInstance(document.getElementById('goalSecondsModal')).hide();
        } catch {
            errEl.textContent = 'Something went wrong. Please try again.';
            errEl.style.display = '';
        }
    });
})();
</script>
