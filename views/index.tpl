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

<main class="container mt-5 mb-4" style="max-width: 640px;">
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

    <div class="card mb-4">
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

    <div class="card mb-4">
        <div class="card-body">
            <h2 class="h5 fw-bold mb-1">Tracking Weight</h2>
            <p class="text-muted small mb-3">Log your body weight daily to track trends over time.</p>
            <ul class="mb-0">
                <li class="mb-1">Go to <a href="/weight">Weight</a> and enter today's weight in lbs or kg.</li>
                <li class="mb-1">The dashboard shows a 3-day rolling average so short-term fluctuations don't distract you.</li>
                <li>Your preferred unit is set in <a href="/settings">Settings</a>.</li>
            </ul>
        </div>
    </div>

    <div class="card mb-4">
        <div class="card-body">
            <h2 class="h5 fw-bold mb-1">Tracking Macros</h2>
            <p class="text-muted small mb-3">Log food items with as much or as little detail as you want — no need to track everything.</p>
            <ul class="mb-0">
                <li class="mb-1">Go to <a href="/macros">Macros</a> and add a food item with its name and the macros you know.</li>
                <li class="mb-1">You can log just one macro — for example, only protein — if that's all you care about.</li>
                <li class="mb-1">Log all three (protein, carbs, and fat) and calories are calculated automatically (4 kcal/g protein &amp; carbs, 9 kcal/g fat).</li>
                <li class="mb-1">Set daily goals to see how your 3-day average stacks up — green means at or above goal, red means below.</li>
                <li>Each entry can include a serving size (g, oz, ml, or fl oz) for your reference.</li>
            </ul>
        </div>
    </div>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script src="/static/offline-sync.js"></script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
