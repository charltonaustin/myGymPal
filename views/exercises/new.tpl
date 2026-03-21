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
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
