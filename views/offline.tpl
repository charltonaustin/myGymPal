<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Offline — My Gym Pal</title>
    <link rel="manifest" href="/manifest.json">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body class="bg-light">

<nav class="navbar navbar-dark bg-dark">
    <div class="container">
        <a class="navbar-brand fw-bold" href="/">My Gym Pal</a>
    </div>
</nav>

<main class="container text-center mt-5">
    <div class="py-5">
        <div class="display-1 mb-3">📶</div>
        <h1 class="h3 fw-bold mb-3">You're Offline</h1>
        <p class="text-muted mb-4">This page isn't available offline. Check your connection and try again.</p>
        <button class="btn btn-dark" onclick="window.location.reload()">Retry</button>
    </div>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>
