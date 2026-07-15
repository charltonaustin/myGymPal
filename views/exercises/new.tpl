<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Exercise — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4 mb-4" style="max-width: 560px;">
    <div class="mb-4">
        <a href="/exercises" class="text-muted small">&larr; Exercise Library</a>
        <h1 class="h4 fw-bold mt-1 mb-0">New Exercise</h1>
    </div>

    {{if .Error}}
    <div class="alert alert-danger" id="error-alert">{{.Error}}</div>
    {{end}}

    <form method="POST" action="/exercises/new" novalidate id="exercise-form">
        <div class="card p-3">
            {{template "partials/exercise_fields.tpl" .}}
        </div>

        <div class="d-flex gap-2 mt-3">
            <button type="submit" class="btn btn-dark">Add Exercise</button>
            <a href="/exercises" class="btn btn-outline-secondary">Cancel</a>
        </div>
    </form>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
    const alertEl = document.getElementById('error-alert');
    if (alertEl) setTimeout(() => alertEl.remove(), 4000);

    document.getElementById('exercise-form').addEventListener('submit', function (e) {
        if (!this.checkValidity()) {
            e.preventDefault();
            e.stopPropagation();
        }
        this.classList.add('was-validated');
    });
</script>
<script>
(function () {
    const availableNames = {{.AvailableNamesJSON}};
    const input = document.getElementById('ex_name');
    let dropdown = null;

    function showDropdown(matches) {
        if (!dropdown) {
            dropdown = document.createElement('div');
            dropdown.className = 'list-group shadow';
            dropdown.style.cssText = 'position:absolute;top:100%;left:0;right:0;z-index:1060;max-height:220px;overflow-y:auto;border-radius:0 0 .375rem .375rem;';
            input.parentElement.style.position = 'relative';
            input.parentElement.appendChild(dropdown);
        }
        dropdown.innerHTML = '';
        matches.forEach(function (name) {
            var btn = document.createElement('button');
            btn.type = 'button';
            btn.className = 'list-group-item list-group-item-action py-2 px-3';
            btn.style.fontSize = '0.95rem';
            btn.textContent = name;
            btn.addEventListener('mousedown', function (e) {
                e.preventDefault();
                input.value = name;
                hideDropdown();
            });
            dropdown.appendChild(btn);
        });
    }

    function hideDropdown() {
        if (dropdown) { dropdown.remove(); dropdown = null; }
    }

    input.addEventListener('input', function () {
        var q = this.value.trim().toLowerCase();
        if (!q) { hideDropdown(); return; }
        var matches = availableNames.filter(function (n) {
            return n.toLowerCase().includes(q);
        }).slice(0, 10);
        if (matches.length) showDropdown(matches); else hideDropdown();
    });

    input.addEventListener('blur', function () { setTimeout(hideDropdown, 150); });
    input.addEventListener('keydown', function (e) { if (e.key === 'Escape') hideDropdown(); });
})();
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
