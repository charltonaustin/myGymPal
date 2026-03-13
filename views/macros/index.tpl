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

    {{if .Summary}}
    <div class="card mb-4">
        <div class="card-body">
            <h2 class="h6 fw-semibold mb-1">{{.Summary.Days}}-Day Average</h2>
            <table class="table table-sm mb-0 mt-2">
                <thead>
                    <tr class="text-muted small">
                        <th class="fw-normal ps-0">Macro</th>
                        <th class="fw-normal text-end">Avg</th>
                        {{if .Summary.HasGoal}}<th class="fw-normal text-end">Goal</th><th class="fw-normal text-end">%</th>{{end}}
                    </tr>
                </thead>
                <tbody class="small">
                    <tr>
                        <td class="ps-0">Protein</td>
                        <td class="text-end">{{printf "%.0f" .Summary.Protein.Actual}}g</td>
                        {{if .Summary.HasGoal}}
                        <td class="text-end text-muted">{{printf "%.0f" .Summary.Protein.Goal}}g</td>
                        <td class="text-end fw-semibold {{if .Summary.Protein.AtGoal}}text-success{{else}}text-danger{{end}}">{{.Summary.Protein.Pct}}%</td>
                        {{end}}
                    </tr>
                    <tr>
                        <td class="ps-0">Carbs</td>
                        <td class="text-end">{{printf "%.0f" .Summary.Carbs.Actual}}g</td>
                        {{if .Summary.HasGoal}}
                        <td class="text-end text-muted">{{printf "%.0f" .Summary.Carbs.Goal}}g</td>
                        <td class="text-end fw-semibold {{if .Summary.Carbs.AtGoal}}text-success{{else}}text-danger{{end}}">{{.Summary.Carbs.Pct}}%</td>
                        {{end}}
                    </tr>
                    <tr>
                        <td class="ps-0">Fat</td>
                        <td class="text-end">{{printf "%.0f" .Summary.Fat.Actual}}g</td>
                        {{if .Summary.HasGoal}}
                        <td class="text-end text-muted">{{printf "%.0f" .Summary.Fat.Goal}}g</td>
                        <td class="text-end fw-semibold {{if .Summary.Fat.AtGoal}}text-success{{else}}text-danger{{end}}">{{.Summary.Fat.Pct}}%</td>
                        {{end}}
                    </tr>
                    <tr class="border-top">
                        <td class="ps-0">Calories</td>
                        <td class="text-end">{{printf "%.0f" .Summary.Calories.Actual}} kcal</td>
                        {{if .Summary.HasGoal}}
                        <td class="text-end text-muted">{{printf "%.0f" .Summary.Calories.Goal}} kcal</td>
                        <td class="text-end fw-semibold {{if .Summary.Calories.AtGoal}}text-success{{else}}text-danger{{end}}">{{.Summary.Calories.Pct}}%</td>
                        {{end}}
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
    {{end}}

    <div class="card mb-4">
        <div class="card-body">
            <h2 class="h6 fw-semibold mb-3">Daily Goals</h2>
            <form method="POST" action="/macros/goals">
                <div class="row g-2 align-items-end">
                    <div class="col-auto">
                        <label class="form-label small mb-1">Protein (g)</label>
                        <input type="number" name="protein_goal" class="form-control form-control-sm" value="{{if .Goal}}{{printf "%.0f" .Goal.Protein}}{{end}}" placeholder="0" min="0" step="1" style="width: 90px;">
                    </div>
                    <div class="col-auto">
                        <label class="form-label small mb-1">Carbs (g)</label>
                        <input type="number" name="carbs_goal" class="form-control form-control-sm" value="{{if .Goal}}{{printf "%.0f" .Goal.Carbs}}{{end}}" placeholder="0" min="0" step="1" style="width: 90px;">
                    </div>
                    <div class="col-auto">
                        <label class="form-label small mb-1">Fat (g)</label>
                        <input type="number" name="fat_goal" class="form-control form-control-sm" value="{{if .Goal}}{{printf "%.0f" .Goal.Fat}}{{end}}" placeholder="0" min="0" step="1" style="width: 90px;">
                    </div>
                    <div class="col-auto">
                        <button type="submit" class="btn btn-outline-secondary btn-sm">Save</button>
                    </div>
                </div>
            </form>
        </div>
    </div>

    <div class="card mb-4">
        <div class="card-body">
            <h2 class="h6 fw-semibold mb-3">Log Food</h2>
            <form method="POST" action="/macros">
                <div class="row g-2 mb-2">
                    <div class="col-4" id="food-name-wrapper">
                        <input type="text" id="food-name-input" name="food_name" class="form-control form-control-sm" placeholder="Food name" required autocomplete="off">
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
        {{$day := .}}
        <div class="d-flex align-items-baseline justify-content-between mb-2">
            <h2 class="h6 fw-semibold mb-0">{{.Date.Format "Mon, Jan 2, 2006"}}</h2>
            <span class="text-muted small">
                {{with $.Goal}}
                P {{printf "%.0f" $day.Protein}}/{{printf "%.0f" .Protein}}g &middot;
                C {{printf "%.0f" $day.Carbs}}/{{printf "%.0f" .Carbs}}g &middot;
                F {{printf "%.0f" $day.Fat}}/{{printf "%.0f" .Fat}}g &middot;
                {{printf "%.0f" $day.Calories}} kcal
                {{else}}
                P {{printf "%.0f" .Protein}}g &middot;
                C {{printf "%.0f" .Carbs}}g &middot;
                F {{printf "%.0f" .Fat}}g &middot;
                {{printf "%.0f" .Calories}} kcal
                {{end}}
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

    // Food history autocomplete
    const foodHistory = {{.FoodHistoryJSON}};
    const foodMap = {};
    foodHistory.forEach(f => { foodMap[f.name] = f; });

    const foodInput = document.getElementById('food-name-input');
    const foodWrapper = document.getElementById('food-name-wrapper');
    foodWrapper.style.position = 'relative';
    let dropdown = null;

    function fillForm(food) {
        foodInput.value = food.name;
        const form = foodInput.closest('form');
        if (food.servingWeight > 0) form.querySelector('[name="serving_weight"]').value = food.servingWeight;
        const unitSel = form.querySelector('[name="serving_unit"]');
        if (food.servingUnit) unitSel.value = food.servingUnit;
        if (food.protein) form.querySelector('[name="protein"]').value = food.protein;
        if (food.carbs)   form.querySelector('[name="carbs"]').value = food.carbs;
        if (food.fat)     form.querySelector('[name="fat"]').value = food.fat;
    }

    function showDropdown(matches) {
        if (!dropdown) {
            dropdown = document.createElement('div');
            dropdown.className = 'list-group shadow';
            dropdown.style.cssText = 'position:absolute;top:100%;left:0;right:0;z-index:1050;max-height:220px;overflow-y:auto;border-radius:0 0 .375rem .375rem;';
            foodWrapper.appendChild(dropdown);
        }
        dropdown.innerHTML = '';
        matches.forEach(food => {
            const btn = document.createElement('button');
            btn.type = 'button';
            btn.className = 'list-group-item list-group-item-action py-2 px-3';
            btn.style.fontSize = '0.9rem';
            btn.textContent = food.name;
            btn.addEventListener('mousedown', e => {
                e.preventDefault();
                fillForm(food);
                hideDropdown();
            });
            dropdown.appendChild(btn);
        });
    }

    function hideDropdown() {
        if (dropdown) { dropdown.remove(); dropdown = null; }
    }

    foodInput.addEventListener('input', () => {
        const val = foodInput.value.toLowerCase().trim();
        if (!val) { hideDropdown(); return; }
        const matches = foodHistory.filter(f => f.name.toLowerCase().includes(val)).slice(0, 10);
        if (matches.length) showDropdown(matches);
        else hideDropdown();
    });

    foodInput.addEventListener('blur', () => setTimeout(hideDropdown, 150));
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
