<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>New Program — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4" style="max-width: 480px;">
    <h1 class="h4 fw-bold mb-4">New Training Program</h1>

    {{if .Error}}
    <div class="alert alert-danger alert-dismissible fade show" id="error-alert">{{.Error}}</div>
    {{end}}

    <form method="POST" action="/programs" novalidate id="program-form" data-offline-sync>
        <div class="mb-3">
            <label for="name" class="form-label">Program Name</label>
            <input
                type="text"
                class="form-control"
                id="name"
                name="name"
                value="{{.Name}}"
                placeholder="e.g. Hypertrophy Block 1"
                required
            >
            <div class="invalid-feedback">Program name is required.</div>
        </div>

        <div class="mb-3">
            <label for="start_date" class="form-label">Start Date</label>
            <input
                type="date"
                class="form-control"
                id="start_date"
                name="start_date"
                value="{{.StartDate}}"
                required
            >
            <div class="invalid-feedback">Start date is required.</div>
        </div>

        <div class="mb-3">
            <label for="num_phases" class="form-label">Number of Phases</label>
            <input
                type="number"
                class="form-control"
                id="num_phases"
                name="num_phases"
                value="{{.NumPhases}}"
                min="1"
                required
            >
            <div class="invalid-feedback">Enter a number of phases (at least 1).</div>
        </div>

        <div class="d-flex gap-2">
            <button type="submit" class="btn btn-dark">Create Program</button>
            <a href="/programs" class="btn btn-outline-secondary">Cancel</a>
        </div>
    </form>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
    const alertEl = document.getElementById('error-alert');
    if (alertEl) {
        setTimeout(() => bootstrap.Alert.getOrCreateInstance(alertEl).close(), 3000);
    }

    document.getElementById('program-form').addEventListener('submit', function (e) {
        if (!this.checkValidity()) {
            e.preventDefault();
            e.stopPropagation();
        }
        this.classList.add('was-validated');
    });
</script>
<script src="/static/offline-sync.js"></script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
