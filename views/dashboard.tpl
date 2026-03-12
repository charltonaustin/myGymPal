<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dashboard — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-5" style="max-width: 640px;">
    <h1 class="h4 fw-bold mb-4">Welcome back, {{.Username}}!</h1>

    <h2 class="h6 fw-semibold text-uppercase text-muted mb-3">Recent Workouts</h2>

    {{if .RecentSessions}}
    <div class="list-group">
        {{range .RecentSessions}}
        <a href="/sessions/{{.ID}}" class="list-group-item list-group-item-action d-flex justify-content-between align-items-center">
            <div>
                <div class="fw-semibold">{{.ProgramName}}</div>
                <div class="text-muted small">Phase {{.PhaseNumber}} &middot; Week {{.WeekNumber}} &middot; Workout {{.WorkoutNumber}}{{if .IsDeload}} &middot; Deload{{end}}</div>
            </div>
            <span class="text-muted small">{{.Date.Format "Jan 2, 2006"}}</span>
        </a>
        {{end}}
    </div>
    {{else}}
    <p class="text-muted">No workouts logged yet. <a href="/programs">Start a session</a> from a program.</p>
    {{end}}
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script src="/static/offline-sync.js"></script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
