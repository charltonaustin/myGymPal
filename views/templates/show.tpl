<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Template.Name}} — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4" style="max-width: 600px;">
    <div class="mb-4">
        <a href="/templates" class="text-muted small">&larr; Templates</a>
        <div class="d-flex align-items-center justify-content-between mt-1">
            <div>
                <h1 class="h4 fw-bold mb-1">{{.Template.Name}}</h1>
                {{if .Template.Focus}}<p class="text-muted small mb-0">{{.Template.Focus}}</p>{{end}}
            </div>
            <a href="/templates/{{.Template.ID}}/edit" class="btn btn-outline-secondary btn-sm">Edit</a>
        </div>
    </div>

    {{if .Success}}
    <div class="alert alert-success alert-dismissible fade show" id="success-alert">{{.Success}}</div>
    {{end}}

    <h2 class="h6 fw-semibold text-uppercase text-muted mb-3">Exercises</h2>

    {{if .Exercises}}
    <div class="card">
        <ul class="list-group list-group-flush">
            {{range .Exercises}}
            <li class="list-group-item">
                <div class="fw-semibold">{{.Name}}</div>
                {{if .IsBodyweight}}<div class="text-muted small">Bodyweight</div>{{end}}
            </li>
            {{end}}
        </ul>
    </div>
    {{else}}
    <p class="text-muted">No exercises in this template.</p>
    {{end}}
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
    const alertEl = document.getElementById('success-alert');
    if (alertEl) {
        setTimeout(() => bootstrap.Alert.getOrCreateInstance(alertEl).close(), 3000);
    }
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
