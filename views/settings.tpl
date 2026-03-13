<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Settings — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-5 mb-4" style="max-width: 480px;">
    <h1 class="h4 fw-bold mb-4">Account Settings</h1>

    {{if .Success}}
    <div class="alert alert-success alert-dismissible fade show" id="success-alert">{{.Success}}</div>
    {{end}}

    {{if .Error}}
    <div class="alert alert-danger">{{.Error}}</div>
    {{end}}

    <form method="post" action="/settings" data-offline-sync>
        <fieldset class="mb-4">
            <legend class="form-label fw-semibold">Weight Unit</legend>
            <div class="form-check">
                <input class="form-check-input" type="radio" name="weight_unit" id="unit_lb" value="lb"
                    {{if eq .WeightUnit "lb"}}checked{{end}}>
                <label class="form-check-label" for="unit_lb">lb (pounds)</label>
            </div>
            <div class="form-check">
                <input class="form-check-input" type="radio" name="weight_unit" id="unit_kg" value="kg"
                    {{if eq .WeightUnit "kg"}}checked{{end}}>
                <label class="form-check-label" for="unit_kg">kg (kilograms)</label>
            </div>
        </fieldset>

        <button type="submit" class="btn btn-dark">Save Settings</button>
    </form>

    <hr class="my-5">

    <div class="mb-4">
        <h2 class="h6 fw-semibold text-uppercase text-muted mb-1">Danger Zone</h2>
        <p class="text-muted small mb-3">This will permanently delete your account and all associated data including programs, sessions, and workout history. This cannot be undone.</p>
        <button type="button" class="btn btn-outline-danger btn-sm" data-bs-toggle="modal" data-bs-target="#deleteAccountModal">
            Delete Account
        </button>
    </div>
</main>

<div class="modal fade" id="deleteAccountModal" tabindex="-1">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title">Delete Account</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
                <p>Are you sure you want to delete your account? This will permanently remove all your programs, sessions, and workout history.</p>
                <p class="fw-semibold mb-0">This action cannot be undone.</p>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>
                <form method="post" action="/account/delete">
                    <button type="submit" class="btn btn-danger">Yes, Delete Everything</button>
                </form>
            </div>
        </div>
    </div>
</div>

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
