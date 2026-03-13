<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dashboard — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-5" style="max-width: 640px;">
    <h1 class="h4 fw-bold mb-4">Welcome back, {{.Username}}!</h1>

    <div class="row g-3 mb-4">
        <div class="col-6">
            <div class="card h-100">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-2">
                        <h2 class="h6 fw-semibold mb-0">Weight</h2>
                        <a href="/weight" class="text-muted small">Log Weight</a>
                    </div>
                    {{if .WeightAvg}}
                    <p class="fs-4 fw-bold mb-0">{{printf "%.1f" .WeightAvg.Weight}} <span class="fs-6 fw-normal text-muted">{{.WeightAvg.Unit}}</span></p>
                    <p class="text-muted small mb-0">{{.WeightAvg.Days}}-day avg</p>
                    {{else}}
                    <p class="text-muted small mb-0">No entries yet.</p>
                    {{end}}
                </div>
            </div>
        </div>
        <div class="col-6">
            <div class="card h-100">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-2">
                        <h2 class="h6 fw-semibold mb-0">Macros</h2>
                        <a href="/macros" class="text-muted small">Log Macros</a>
                    </div>
                    {{if .MacroSummary}}
                    <p class="text-muted small mb-1">{{.MacroSummary.Days}}-day avg</p>
                    <div class="small">
                        <div class="d-flex justify-content-between">
                            <span class="text-muted">Protein</span>
                            <span class="fw-semibold">{{printf "%.0f" .MacroSummary.Protein.Actual}}g
                            {{if .MacroSummary.HasGoal}}<span class="{{if .MacroSummary.Protein.AtGoal}}text-success{{else}}text-danger{{end}}">{{.MacroSummary.Protein.Pct}}%</span>{{end}}</span>
                        </div>
                        <div class="d-flex justify-content-between">
                            <span class="text-muted">Carbs</span>
                            <span class="fw-semibold">{{printf "%.0f" .MacroSummary.Carbs.Actual}}g
                            {{if .MacroSummary.HasGoal}}<span class="{{if .MacroSummary.Carbs.AtGoal}}text-success{{else}}text-danger{{end}}">{{.MacroSummary.Carbs.Pct}}%</span>{{end}}</span>
                        </div>
                        <div class="d-flex justify-content-between">
                            <span class="text-muted">Fat</span>
                            <span class="fw-semibold">{{printf "%.0f" .MacroSummary.Fat.Actual}}g
                            {{if .MacroSummary.HasGoal}}<span class="{{if .MacroSummary.Fat.AtGoal}}text-success{{else}}text-danger{{end}}">{{.MacroSummary.Fat.Pct}}%</span>{{end}}</span>
                        </div>
                        <div class="d-flex justify-content-between border-top mt-1 pt-1">
                            <span class="text-muted">Calories</span>
                            <span class="fw-semibold">{{printf "%.0f" .MacroSummary.Calories.Actual}}
                            {{if .MacroSummary.HasGoal}}<span class="{{if .MacroSummary.Calories.AtGoal}}text-success{{else}}text-danger{{end}}">{{.MacroSummary.Calories.Pct}}%</span>{{end}}</span>
                        </div>
                    </div>
                    {{else}}
                    <p class="text-muted small mb-0">No entries yet.</p>
                    {{end}}
                </div>
            </div>
        </div>
    </div>
    
    <div class="d-flex justify-content-between align-items-baseline mb-3">
        <h2 class="h6 fw-semibold text-uppercase text-muted mb-0">Recent Workouts</h2>
        {{if .RecentSessions}}<a href="/programs/{{(index .RecentSessions 0).ProgramID}}/sessions/new" class="text-muted small">Log next workout</a>{{end}}
    </div>

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
