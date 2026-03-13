<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Programs — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4" style="max-width: 640px;">
    <div class="d-flex justify-content-between align-items-center mb-4">
        <h1 class="h4 fw-bold mb-0">Training Programs</h1>
        <a href="/programs/new" class="btn btn-dark btn-sm">+ New Program</a>
    </div>

    {{if .Success}}
    <div class="alert alert-success alert-dismissible fade show" id="success-alert">{{.Success}}</div>
    {{end}}

    {{if .Programs}}
    <ul class="list-group">
        {{range .Programs}}
        <li class="list-group-item list-group-item-action d-flex justify-content-between align-items-center">
            <a href="/programs/{{.ID}}" class="text-decoration-none text-dark flex-grow-1">
                <div class="fw-semibold text-capitalize">{{.Name}} </div>
                <div class="text-muted small">
                    Starts {{.StartDate.Format "Jan 2, 2006"}} &middot; {{.NumPhases}} phase{{if gt .NumPhases 1}}s{{end}} &middot; {{.WeeksPerPhase}} week{{if gt .WeeksPerPhase 1}}s{{end}}/phase
                </div>
            </a>
            <button type="button" class="btn btn-link btn-sm p-0 text-danger ms-3 flex-shrink-0"
                data-bs-toggle="modal" data-bs-target="#deleteModal"
                data-program-id="{{.ID}}" data-program-name="{{.Name}}">
                <i class="bi bi-trash"></i>
            </button>
        </li>
        {{end}}
    </ul>
    {{else}}
    <p class="text-muted">No programs yet. <a href="/programs/new">Create your first one.</a></p>
    {{end}}
</main>

<!-- Delete confirmation modal -->
<div class="modal fade" id="deleteModal" tabindex="-1">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title">Delete program?</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
                <p class="mb-0">This will permanently delete <strong id="deleteModalName"></strong> and all its sessions. This cannot be undone.</p>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Cancel</button>
                <form id="deleteForm" method="POST">
                    <button type="submit" class="btn btn-danger btn-sm">Delete</button>
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

    document.getElementById('deleteModal').addEventListener('show.bs.modal', e => {
        const btn = e.relatedTarget;
        document.getElementById('deleteModalName').textContent = btn.dataset.programName;
        document.getElementById('deleteForm').action = '/programs/' + btn.dataset.programId + '/delete';
    });
</script>
<script src="/static/offline-sync.js"></script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
