<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Templates — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4" style="max-width: 600px;">
    <div class="d-flex justify-content-between align-items-center mb-4">
        <h1 class="h4 fw-bold mb-0">Workout Templates</h1>
        <a href="/templates/new" class="btn btn-dark btn-sm">+ New Template</a>
    </div>

    {{if .Success}}
    <div class="alert alert-success alert-dismissible fade show" id="success-alert">{{.Success}}</div>
    {{end}}

    {{if .Templates}}
    <div class="list-group">
        {{range .Templates}}
        <div class="list-group-item list-group-item-action d-flex justify-content-between align-items-center">
            <a href="/templates/{{.ID}}" class="text-decoration-none text-dark flex-grow-1">
                <div class="fw-semibold">{{.Name}}</div>
                {{if .Focus}}<div class="text-muted small">{{.Focus}}</div>{{end}}
            </a>
            <button type="button" class="btn btn-link btn-sm p-0 text-danger ms-3 flex-shrink-0"
                data-bs-toggle="modal" data-bs-target="#deleteModal"
                data-template-id="{{.ID}}" data-template-name="{{.Name}}">
                <i class="bi bi-trash"></i>
            </button>
        </div>
        {{end}}
    </div>
    {{else}}
    <p class="text-muted">No templates yet. <a href="/templates/new">Create one</a> to get started.</p>
    {{end}}
</main>

<!-- Delete confirmation modal -->
<div class="modal fade" id="deleteModal" tabindex="-1">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title">Delete template?</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
            </div>
            <div class="modal-body">
                <p class="mb-0">This will permanently delete <strong id="deleteModalName"></strong>. This cannot be undone.</p>
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
        document.getElementById('deleteModalName').textContent = btn.dataset.templateName;
        document.getElementById('deleteForm').action = '/templates/' + btn.dataset.templateId + '/delete';
    });
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
