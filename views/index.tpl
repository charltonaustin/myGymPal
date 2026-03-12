<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-5" style="max-width: 640px;">
    <h1 class="display-4 fw-bold mb-3 text-center">Welcome to My Gym Pal</h1>
    <p class="lead text-muted text-center mb-5">Your personal fitness companion.</p>

    <div class="card mb-4">
        <div class="card-body">
            <h2 class="h5 fw-bold mb-1">Quick Start</h2>
            <p class="text-muted small mb-3">Get up and running in two steps.</p>
            <ol class="mb-0">
                <li class="mb-1"><a href="/programs/new">Create a Program</a> — set your schedule and rep ranges by phase.</li>
                <li>Start a workout — open your program and log a session.</li>
            </ol>
        </div>
    </div>

    <div class="card">
        <div class="card-body">
            <h2 class="h5 fw-bold mb-1">Full Setup</h2>
            <p class="text-muted small mb-3">Build a reusable exercise library and templates for a more structured experience.</p>
            <ol class="mb-0">
                <li class="mb-1"><a href="/exercises/new">Create exercises</a> — add names and goal weights to your library.</li>
                <li class="mb-1"><a href="/templates/new">Create a template</a> — build a workout layout from your exercises.</li>
                <li class="mb-1"><a href="/programs/new">Create a program</a> — attach templates to your training schedule.</li>
                <li>Start a workout from a template and track your sets.</li>
            </ol>
        </div>
    </div>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script src="/static/offline-sync.js"></script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
