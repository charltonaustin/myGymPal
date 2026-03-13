<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Macros — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4" style="max-width: 600px;">
    <h1 class="h4 fw-bold mb-4">Macros</h1>

    <div class="card mb-4">
        <div class="card-body">
            <h2 class="h6 fw-semibold mb-3">Log Food</h2>
            <form method="POST" action="/macros">
                <div class="row g-2 mb-2">
                    <div class="col-4">
                        <input type="text" name="food_name" class="form-control form-control-sm" placeholder="Food name" required>
                    </div>
                    <div class="col-3">
                        <input type="date" name="date" class="form-control form-control-sm" value="{{.DefaultDate}}" required>
                    </div>
                    <div class="col-5">
                        <div class="input-group input-group-sm">
                            <input type="number" name="serving_weight" class="form-control" placeholder="Amount" min="0" step="0.1">
                            <select name="serving_unit" class="form-select" style="max-width: 90px;">
                                <option value="g">g</option>
                                <option value="oz">oz</option>
                                <option value="ml">ml</option>
                                <option value="fl oz">fl oz</option>
                            </select>
                        </div>
                    </div>
                </div>
                <div class="row g-2 align-items-end">
                    <div class="col-auto">
                        <label class="form-label small mb-1">Protein (g)</label>
                        <input type="number" name="protein" class="form-control form-control-sm" placeholder="0" min="0" step="0.1" style="width: 90px;">
                    </div>
                    <div class="col-auto">
                        <label class="form-label small mb-1">Carbs (g)</label>
                        <input type="number" name="carbs" class="form-control form-control-sm" placeholder="0" min="0" step="0.1" style="width: 90px;">
                    </div>
                    <div class="col-auto">
                        <label class="form-label small mb-1">Fat (g)</label>
                        <input type="number" name="fat" class="form-control form-control-sm" placeholder="0" min="0" step="0.1" style="width: 90px;">
                    </div>
                    <div class="col-auto">
                        <button type="submit" class="btn btn-dark btn-sm">Add</button>
                    </div>
                </div>
            </form>
        </div>
    </div>

    {{if .Days}}
    {{range .Days}}
    <div class="mb-4">
        <div class="d-flex align-items-baseline justify-content-between mb-2">
            <h2 class="h6 fw-semibold mb-0">{{.Date.Format "Mon, Jan 2, 2006"}}</h2>
            <span class="text-muted small">
                P {{printf "%.0f" .Protein}}g &middot;
                C {{printf "%.0f" .Carbs}}g &middot;
                F {{printf "%.0f" .Fat}}g &middot;
                {{printf "%.0f" .Calories}} kcal
            </span>
        </div>
        <div class="card">
            <ul class="list-group list-group-flush">
                {{range .Entries}}
                <li class="list-group-item py-2">
                    <div class="d-flex align-items-center justify-content-between view-row-{{.ID}}">
                        <div>
                            <span class="fw-semibold small">{{.FoodName}}</span>
                            {{if gt .ServingWeight 0.0}}<span class="text-muted small ms-1">{{printf "%g" .ServingWeight}}{{.ServingUnit}}</span>{{end}}
                            <span class="text-muted small ms-2">
                                Protein {{printf "%.0f" .Protein}}g &middot;
                                Carbs {{printf "%.0f" .Carbs}}g &middot;
                                Fat {{printf "%.0f" .Fat}}g
                            </span>
                        </div>
                        <div class="d-flex gap-2">
                            <button type="button" class="btn btn-link btn-sm p-0 text-secondary" onclick="showEdit({{.ID}})"><i class="bi bi-pencil"></i></button>
                            <form method="POST" action="/macros/{{.ID}}/delete" class="d-inline">
                                <button type="submit" class="btn btn-link btn-sm p-0 text-danger"><i class="bi bi-trash"></i></button>
                            </form>
                        </div>
                    </div>
                    <form method="POST" action="/macros/{{.ID}}" class="d-none edit-row-{{.ID}} mt-2">
                        <div class="row g-2 mb-2">
                            <div class="col-7">
                                <input type="text" name="food_name" class="form-control form-control-sm" value="{{.FoodName}}" required>
                            </div>
                            <div class="col-5">
                                <div class="input-group input-group-sm">
                                    <input type="number" name="serving_weight" class="form-control" value="{{if gt .ServingWeight 0.0}}{{.ServingWeight}}{{end}}" placeholder="Amount" min="0" step="0.1" style="width: 80px;">
                                    <select name="serving_unit" class="form-select" style="max-width: 90px;">
                                        <option value="g" {{if eq .ServingUnit "g"}}selected{{end}}>g</option>
                                        <option value="oz" {{if eq .ServingUnit "oz"}}selected{{end}}>oz</option>
                                        <option value="ml" {{if eq .ServingUnit "ml"}}selected{{end}}>ml</option>
                                        <option value="fl oz" {{if eq .ServingUnit "fl oz"}}selected{{end}}>fl oz</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="row g-2 align-items-center">
                            <div class="col-auto">
                                <label class="form-label small mb-1">Protein (g)</label>
                                <input type="number" name="protein" class="form-control form-control-sm" value="{{.Protein}}" min="0" step="0.1" style="width: 90px;">
                            </div>
                            <div class="col-auto">
                                <label class="form-label small mb-1">Carbs (g)</label>
                                <input type="number" name="carbs" class="form-control form-control-sm" value="{{.Carbs}}" min="0" step="0.1" style="width: 90px;">
                            </div>
                            <div class="col-auto">
                                <label class="form-label small mb-1">Fat (g)</label>
                                <input type="number" name="fat" class="form-control form-control-sm" value="{{.Fat}}" min="0" step="0.1" style="width: 90px;">
                            </div>
                            <div class="col-auto d-flex gap-2 align-self-end">
                                <button type="submit" class="btn btn-dark btn-sm">Save</button>
                                <button type="button" class="btn btn-outline-secondary btn-sm" onclick="hideEdit({{.ID}})">Cancel</button>
                            </div>
                        </div>
                    </form>
                </li>
                {{end}}
            </ul>
        </div>
    </div>
    {{end}}
    {{else}}
    <p class="text-muted">No food logged yet.</p>
    {{end}}
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>
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
