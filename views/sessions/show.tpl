<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Session #{{.Session.WorkoutNumber}} — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4" style="max-width: 600px;">
    <div class="mb-4">
        <a href="/programs/{{.Program.ID}}" class="text-muted small">&larr; {{.Program.Name}}</a>
        <div class="d-flex align-items-center gap-2 mt-1">
            <h1 class="h4 fw-bold mb-0">Session #{{.Session.WorkoutNumber}}</h1>
            {{if .Session.IsDeload}}
            <span class="badge bg-secondary">Deload</span>
            {{end}}
        </div>
        <p class="text-muted small mb-0 mt-1">
            Phase {{.Session.PhaseNumber}} &middot; Week {{.Session.WeekNumber}}
            &middot; {{.Session.Date.Format "Jan 2, 2006"}}
        </p>
    </div>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
