<nav class="navbar navbar-dark bg-dark">
    <div class="container">
        <a class="navbar-brand fw-bold" href="/">My Gym Pal</a>
        <div class="navbar-nav flex-row gap-3 me-auto ms-4">
            {{if .LoggedIn}}
            <a href="/" class="nav-link {{if eq .ActivePage "home"}}active fw-semibold{{end}}">Home</a>
            <a href="/dashboard" class="nav-link {{if eq .ActivePage "dashboard"}}active fw-semibold{{end}}">Dashboard</a>
            <a href="/settings" class="nav-link {{if eq .ActivePage "settings"}}active fw-semibold{{end}}">Settings</a>
            {{else}}
            <a href="/" class="nav-link active fw-semibold">Home</a>
            {{end}}
        </div>
        <div class="d-flex align-items-center gap-2">
            {{if .LoggedIn}}
            <a href="/logout" class="btn btn-outline-light btn-sm">Log Out</a>
            {{else}}
            <a href="/login" class="btn btn-outline-light btn-sm">Log In</a>
            <a href="/register" class="btn btn-light btn-sm">Create Account</a>
            {{end}}
        </div>
    </div>
</nav>
