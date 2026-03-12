<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Exercise Library — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4 mb-4" style="max-width: 600px;">
    <div class="d-flex justify-content-between align-items-center mb-4">
        <h1 class="h4 fw-bold mb-0">Exercise Library</h1>
        <a href="/exercises/new" class="btn btn-dark btn-sm">+ New Exercise</a>
    </div>

    {{if .Success}}
    <div class="alert alert-success alert-dismissible fade show" id="success-alert">{{.Success}}</div>
    {{end}}

    {{if .Exercises}}
    <div class="list-group">
        {{range .Exercises}}
        <div class="list-group-item d-flex justify-content-between align-items-center">
            <div class="flex-grow-1">
                <div class="fw-semibold text-capitalize">
                    {{.Name}}
                    {{if .IsBodyweight}}<span class="badge bg-secondary ms-1">Bodyweight</span>{{end}}
                </div>
                {{if not .IsBodyweight}}
                <div class="text-muted small">Goal: {{.GoalWeight}} {{.WeightUnit}}</div>
                {{end}}
            </div>
            <div class="d-flex gap-2 flex-shrink-0 ms-3">
                <a href="/exercises/{{.ID}}/edit" class="btn btn-outline-secondary btn-sm">Edit</a>
                <button type="button" class="btn btn-outline-danger btn-sm"
                    data-bs-toggle="modal" data-bs-target="#deleteModal"
                    data-exercise-id="{{.ID}}" data-exercise-name="{{.Name}}">
                    Delete
                </button>
            </div>
        </div>
        {{end}}
    </div>
    {{else}}
    <p class="text-muted">No exercises yet. <a href="/exercises/new">Add one</a> to build your library.</p>
    {{end}}
</main>

<!-- Delete confirmation modal -->
<div class="modal fade" id="deleteModal" tabindex="-1">
    <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h5 class="modal-title">Delete exercise?</h5>
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
        document.getElementById('deleteModalName').textContent = btn.dataset.exerciseName;
        document.getElementById('deleteForm').action = '/exercises/' + btn.dataset.exerciseId + '/delete';
    });
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
