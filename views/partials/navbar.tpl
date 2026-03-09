<nav class="navbar navbar-dark bg-dark navbar-expand-md">
    <div class="container">
        <a class="navbar-brand fw-bold" href="/">My Gym Pal</a>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#mainNav"
            aria-controls="mainNav" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="mainNav">
            <div class="navbar-nav me-auto">
                {{if .LoggedIn}}
                <a href="/" class="nav-link {{if eq .ActivePage "home"}}active fw-semibold{{end}}">Home</a>
                <a href="/dashboard" class="nav-link {{if eq .ActivePage "dashboard"}}active fw-semibold{{end}}">Dashboard</a>
                <a href="/settings" class="nav-link {{if eq .ActivePage "settings"}}active fw-semibold{{end}}">Settings</a>
                {{else}}
                <a href="/" class="nav-link active fw-semibold">Home</a>
                {{end}}
            </div>
            <div class="d-flex align-items-center gap-2 mt-2 mt-md-0">
                {{if .LoggedIn}}
                <a href="/logout" class="btn btn-outline-light btn-sm">Log Out</a>
                {{else}}
                <a href="/login" class="btn btn-outline-light btn-sm">Log In</a>
                <a href="/register" class="btn btn-light btn-sm">Create Account</a>
                {{end}}
            </div>
        </div>
    </div>
</nav>
