<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Create Account — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body class="bg-light">

<nav class="navbar navbar-dark bg-dark">
    <div class="container">
        <a class="navbar-brand fw-bold" href="/">My Gym Pal</a>
    </div>
</nav>

<main class="container" style="max-width: 420px; margin-top: 3rem;">
    <div class="card shadow-sm">
        <div class="card-body p-4">
            <h1 class="h4 fw-bold mb-4 text-center">Create Account</h1>

            {{if .Error}}
            <div class="alert alert-danger alert-dismissible fade show" id="error-alert">{{.Error}}</div>
            {{end}}

            <form method="POST" action="/register" novalidate id="register-form">
                <div class="mb-3">
                    <label for="username" class="form-label">Username</label>
                    <input
                        type="text"
                        class="form-control"
                        id="username"
                        name="username"
                        value="{{.Username}}"
                        autocomplete="username"
                        required
                    >
                    <div class="invalid-feedback">Username is required.</div>
                </div>

                <div class="mb-3">
                    <label for="password" class="form-label">Password</label>
                    <input
                        type="password"
                        class="form-control"
                        id="password"
                        name="password"
                        autocomplete="new-password"
                        minlength="8"
                        required
                    >
                    <div class="form-text">At least 8 characters.</div>
                    <div class="invalid-feedback">Password must be at least 8 characters.</div>
                </div>

                <div class="mb-3">
                    <label for="confirm_password" class="form-label">Confirm Password</label>
                    <input
                        type="password"
                        class="form-control"
                        id="confirm_password"
                        name="confirm_password"
                        autocomplete="new-password"
                        required
                    >
                    <div class="invalid-feedback">Passwords do not match.</div>
                </div>

                <button type="submit" class="btn btn-dark w-100 mt-2">Create Account</button>
            </form>

            <p class="text-center text-muted mt-3 mb-0 small">
                Already have an account? <a href="/login">Log in</a>
            </p>
        </div>
    </div>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
    const alertEl = document.getElementById('error-alert');
    if (alertEl) {
        setTimeout(() => bootstrap.Alert.getOrCreateInstance(alertEl).close(), 3000);
    }

    const password = document.getElementById('password');
    const confirm = document.getElementById('confirm_password');

    function checkPasswordMatch() {
        confirm.setCustomValidity(confirm.value && confirm.value !== password.value ? 'Passwords do not match.' : '');
    }

    password.addEventListener('input', checkPasswordMatch);
    confirm.addEventListener('input', checkPasswordMatch);

    document.getElementById('register-form').addEventListener('submit', function (e) {
        checkPasswordMatch();
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
