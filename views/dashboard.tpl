<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Dashboard — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>

<nav class="navbar navbar-dark bg-dark">
    <div class="container">
        <a class="navbar-brand fw-bold" href="/">My Gym Pal</a>
        <div class="d-flex align-items-center gap-3">
            <span class="text-white-50 small">{{.Username}}</span>
            <a href="/logout" class="btn btn-outline-light btn-sm">Log out</a>
        </div>
    </div>
</nav>

<main class="container mt-5">
    <h1 class="h4 fw-bold">Welcome back, {{.Username}}!</h1>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>
