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

<main class="container mt-5" style="max-width: 480px;">
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

        <button type="submit" class="btn btn-primary">Save Settings</button>
    </form>
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
