<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Weight — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4 mb-4" style="max-width: 560px;">
    <h1 class="h4 fw-bold mb-4">Weight</h1>

    {{if .Success}}
    <div class="alert alert-success alert-dismissible fade show" id="success-alert">{{.Success}}</div>
    {{end}}

    {{if .WeightAvg}}
    <div class="card mb-4">
        <div class="card-body">
            <h2 class="h6 fw-semibold mb-1">{{.WeightAvg.Days}}-Day Average</h2>
            <p class="display-6 fw-bold mb-0 mt-2">{{printf "%.1f" .WeightAvg.Weight}} <span class="fs-5 fw-normal text-muted">{{.WeightAvg.Unit}}</span></p>
        </div>
    </div>
    {{end}}

    <div class="card mb-4">
        <div class="card-body">
            <h2 class="h6 fw-semibold mb-3">Log Weight</h2>
            <form method="POST" action="/weight" class="d-flex gap-2 align-items-end flex-wrap">
                <div>
                    <label class="form-label small mb-1">Date</label>
                    <input type="date" name="date" class="form-control form-control-sm" value="{{.DefaultDate}}" required style="width: 150px;">
                </div>
                <div>
                    <label class="form-label small mb-1">Weight</label>
                    <div class="input-group input-group-sm" style="width: 160px;">
                        <input type="number" name="weight" class="form-control" placeholder="0" min="0" step="0.1" required>
                        <select name="weight_unit" class="form-select" style="max-width: 72px;">
                            <option value="lb" {{if eq .WeightUnit "lb"}}selected{{end}}>lb</option>
                            <option value="kg" {{if eq .WeightUnit "kg"}}selected{{end}}>kg</option>
                        </select>
                    </div>
                </div>
                <button type="submit" class="btn btn-dark btn-sm mb-0">Add</button>
            </form>
        </div>
    </div>

    {{if .Entries}}
    <div class="card">
        <ul class="list-group list-group-flush">
            {{range .Entries}}
            <li class="list-group-item" id="entry-{{.ID}}">
                <div class="d-flex align-items-center justify-content-between view-row-{{.ID}}">
                    <div>
                        <span class="fw-semibold">{{printf "%.1f" .Weight}} {{.WeightUnit}}</span>
                        <span class="text-muted small ms-2">{{.Date.Format "Jan 2, 2006"}}</span>
                    </div>
                    <div class="d-flex gap-2">
                        <button type="button" class="btn btn-link btn-sm p-0 text-secondary" onclick="showEdit({{.ID}})"><i class="bi bi-pencil"></i></button>
                        <form method="POST" action="/weight/{{.ID}}/delete" class="d-inline">
                            <button type="submit" class="btn btn-link btn-sm p-0 text-danger"><i class="bi bi-trash"></i></button>
                        </form>
                    </div>
                </div>
                <form method="POST" action="/weight/{{.ID}}" class="d-none edit-row-{{.ID}} d-flex gap-2 align-items-end flex-wrap mt-2">
                    <div>
                        <div class="input-group input-group-sm" style="width: 160px;">
                            <input type="number" name="weight" class="form-control" value="{{printf "%.1f" .Weight}}" min="0" step="0.1" required>
                            <select name="weight_unit" class="form-select" style="max-width: 72px;">
                                <option value="lb" {{if eq .WeightUnit "lb"}}selected{{end}}>lb</option>
                                <option value="kg" {{if eq .WeightUnit "kg"}}selected{{end}}>kg</option>
                            </select>
                        </div>
                    </div>
                    <button type="submit" class="btn btn-dark btn-sm">Save</button>
                    <button type="button" class="btn btn-outline-secondary btn-sm" onclick="hideEdit({{.ID}})">Cancel</button>
                </form>
            </li>
            {{end}}
        </ul>
    </div>
    {{else}}
    <p class="text-muted">No weight entries yet.</p>
    {{end}}
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
    const alertEl = document.getElementById('success-alert');
    if (alertEl) setTimeout(() => bootstrap.Alert.getOrCreateInstance(alertEl).close(), 3000);

    function showEdit(id) {
        document.querySelector(`.view-row-${id}`).classList.add('d-none');
        document.querySelector(`.edit-row-${id}`).classList.remove('d-none');
    }
    function hideEdit(id) {
        document.querySelector(`.edit-row-${id}`).classList.add('d-none');
        document.querySelector(`.view-row-${id}`).classList.remove('d-none');
    }
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
