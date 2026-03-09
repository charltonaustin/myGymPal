<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Programs — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4" style="max-width: 640px;">
    <div class="d-flex justify-content-between align-items-center mb-4">
        <h1 class="h4 fw-bold mb-0">Training Programs</h1>
        <a href="/programs/new" class="btn btn-dark btn-sm">+ New Program</a>
    </div>

    {{if .Success}}
    <div class="alert alert-success alert-dismissible fade show" id="success-alert">{{.Success}}</div>
    {{end}}

    {{if .Programs}}
    <ul class="list-group">
        {{range .Programs}}
        <li class="list-group-item d-flex justify-content-between align-items-center">
            <div>
                <div class="fw-semibold">{{.Name}}</div>
                <div class="text-muted small">
                    Starts {{.StartDate.Format "Jan 2, 2006"}} &middot; {{.NumPhases}} phase{{if gt .NumPhases 1}}s{{end}} &middot; {{.WeeksPerPhase}} week{{if gt .WeeksPerPhase 1}}s{{end}}/phase
                </div>
            </div>
        </li>
        {{end}}
    </ul>
    {{else}}
    <p class="text-muted">No programs yet. <a href="/programs/new">Create your first one.</a></p>
    {{end}}
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
    const alertEl = document.getElementById('success-alert');
    if (alertEl) {
        setTimeout(() => bootstrap.Alert.getOrCreateInstance(alertEl).close(), 3000);
    }
</script>
<script src="/static/offline-sync.js"></script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
