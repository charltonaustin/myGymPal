<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Exercise History — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css">
    <link rel="manifest" href="/manifest.json">
    <style>
        .hm-scroll { overflow-x: auto; padding-bottom: 4px; }
        .hm-inner { display: inline-block; }
        .hm-months { display: grid; grid-auto-flow: column; grid-auto-columns: 13px;
                     margin-left: 30px; height: 14px; font-size: 10px; color: #898781; }
        .hm-months span { overflow: visible; white-space: nowrap; }
        .hm-body { display: flex; }
        .hm-weekdays { display: grid; grid-template-rows: repeat(7, 11px); gap: 2px;
                       width: 30px; font-size: 9px; color: #898781; }
        .hm-weekdays span { line-height: 11px; }
        .hm-grid { display: grid; grid-template-rows: repeat(7, 11px); grid-auto-flow: column;
                   grid-auto-columns: 11px; gap: 2px; }
        .hm-cell { width: 11px; height: 11px; border-radius: 2px; background: #ebedf0; }
        .hm-l1 { background: #b7d3f6; }
        .hm-l2 { background: #6da7ec; }
        .hm-l3 { background: #2a78d6; }
        .hm-l4 { background: #184f95; }
        .hm-legend { display: flex; align-items: center; gap: 3px; justify-content: flex-end;
                     margin-top: 6px; font-size: 10px; color: #898781; }
        .hm-legend .hm-cell { display: inline-block; }
    </style>
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4 mb-4" style="max-width: 720px;">
    <div class="d-flex justify-content-between align-items-center mb-4">
        <div>
            <a href="/exercises" class="text-muted small">&larr; Exercise Library</a>
            <h1 class="h4 fw-bold mt-1 mb-0">Exercise History</h1>
        </div>
        <div class="d-flex align-items-center gap-2">
            <select id="range-select" class="form-select form-select-sm" aria-label="Time range" style="width:auto;">
                <option value="14" {{if eq .DefaultDays 14}}selected{{end}}>2 weeks</option>
                <option value="30">1 month</option>
                <option value="90">3 months</option>
                <option value="180">6 months</option>
                <option value="365">1 year</option>
                <option value="0">All time</option>
            </select>
            <div class="btn-group btn-group-sm" role="group" aria-label="Weight unit">
                <input type="radio" class="btn-check" name="unit_toggle" id="unit_lb" value="lb" autocomplete="off" {{if eq .WeightUnit "lb"}}checked{{end}}>
                <label class="btn btn-outline-secondary" for="unit_lb">lb</label>
                <input type="radio" class="btn-check" name="unit_toggle" id="unit_kg" value="kg" autocomplete="off" {{if eq .WeightUnit "kg"}}checked{{end}}>
                <label class="btn btn-outline-secondary" for="unit_kg">kg</label>
            </div>
        </div>
    </div>

    <div class="card mb-3">
        <div class="card-body">
            <div class="d-flex justify-content-between align-items-center mb-2">
                <h2 class="h6 fw-bold mb-0">Workout Activity</h2>
                <span class="text-muted small" id="heatmap-summary"></span>
            </div>
            <div class="hm-scroll">
                <div class="hm-inner">
                    <div class="hm-months" id="hm-months"></div>
                    <div class="hm-body">
                        <div class="hm-weekdays">
                            <span></span><span>Mon</span><span></span><span>Wed</span>
                            <span></span><span>Fri</span><span></span>
                        </div>
                        <div class="hm-grid" id="hm-grid"></div>
                    </div>
                </div>
            </div>
            <div class="hm-legend">
                <span>Less</span>
                <span class="hm-cell"></span>
                <span class="hm-cell hm-l1"></span>
                <span class="hm-cell hm-l2"></span>
                <span class="hm-cell hm-l3"></span>
                <span class="hm-cell hm-l4"></span>
                <span>More</span>
            </div>
        </div>
    </div>

    <div class="card mb-3">
        <div class="card-body">
            <div class="d-flex flex-wrap gap-2 align-items-center mb-2">
                <div style="position:relative;">
                    <input type="text" id="template-search" class="form-control form-control-sm"
                           placeholder="Load a template…" autocomplete="off" style="width:200px;">
                </div>
                <button type="button" id="clear-all" class="btn btn-sm btn-outline-secondary">Clear all</button>
            </div>
            <div class="d-flex flex-wrap gap-2 align-items-center" id="chips-area">
                <div style="position:relative;">
                    <input type="text" id="exercise-search" class="form-control form-control-sm"
                           placeholder="Add exercise…" autocomplete="off" style="width:200px;">
                </div>
            </div>
        </div>
    </div>

    <div class="card">
        <div class="card-body">
            <div id="chart-empty" class="text-muted text-center py-4">
                Add an exercise above to see its history.
            </div>
            <canvas id="history-chart" style="display:none;"></canvas>
        </div>
    </div>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.3/dist/chart.umd.min.js"></script>
<script src="/static/autocomplete.js"></script>
<script>
(function () {
    const exerciseNames = {{.UserExerciseNamesJSON}};
    const defaultNames = {{.DefaultExerciseNamesJSON}};
    const templates = {{.TemplatesJSON}};
    let selected = [];
    let currentUnit = document.querySelector('input[name="unit_toggle"]:checked').value;
    const rangeSelect = document.getElementById('range-select');
    let currentDays = rangeSelect.value;
    let chart = null;

    rangeSelect.addEventListener('change', function () {
        currentDays = this.value;
        if (selected.length > 0) fetchAndRender();
    });

    const searchInput = document.getElementById('exercise-search');
    makeAutocomplete(searchInput, exerciseNames, function (name) {
        addExercise(name);
        searchInput.value = '';
    });

    const templateInput = document.getElementById('template-search');
    makeAutocomplete(templateInput, templates.map(function (t) { return t.name; }), function (name) {
        loadTemplate(name);
        templateInput.value = '';
    });

    document.getElementById('clear-all').addEventListener('click', clearAll);

    function loadTemplate(name) {
        const tpl = templates.find(function (t) { return t.name === name; });
        if (!tpl) return;
        // Replace the current selection with this template's exercises (de-duplicated).
        selected = tpl.exercises.filter(function (n, i) { return tpl.exercises.indexOf(n) === i; });
        renderChips();
        if (selected.length > 0) {
            fetchAndRender();
        } else {
            showEmpty();
        }
    }

    function clearAll() {
        selected = [];
        renderChips();
        showEmpty();
    }

    document.querySelectorAll('input[name="unit_toggle"]').forEach(function (radio) {
        radio.addEventListener('change', function () {
            currentUnit = this.value;
            if (selected.length > 0) fetchAndRender();
        });
    });

    function addExercise(name) {
        if (selected.includes(name)) return;
        selected.push(name);
        renderChips();
        fetchAndRender();
    }

    function removeExercise(name) {
        selected = selected.filter(function (n) { return n !== name; });
        renderChips();
        if (selected.length > 0) {
            fetchAndRender();
        } else {
            showEmpty();
        }
    }

    function renderChips() {
        const area = document.getElementById('chips-area');
        area.querySelectorAll('.exercise-chip').forEach(function (c) { c.remove(); });
        const searchWrapper = area.querySelector('div');
        selected.forEach(function (name) {
            const chip = document.createElement('span');
            chip.className = 'badge bg-dark d-inline-flex align-items-center gap-1 exercise-chip';
            chip.style.fontSize = '0.85rem';
            chip.style.padding = '0.4em 0.6em';
            const label = document.createTextNode(name + ' ');
            chip.appendChild(label);
            const btn = document.createElement('button');
            btn.type = 'button';
            btn.className = 'btn-close btn-close-white';
            btn.style.fontSize = '0.6rem';
            btn.setAttribute('aria-label', 'Remove ' + name);
            btn.addEventListener('click', function () { removeExercise(name); });
            chip.appendChild(btn);
            area.insertBefore(chip, searchWrapper);
        });
    }

    function fetchAndRender() {
        const params = new URLSearchParams({ names: selected.join(','), unit: currentUnit, days: currentDays });
        fetch('/exercises/history/data?' + params)
            .then(function (r) { return r.json(); })
            .then(function (data) { renderChart(data); })
            .catch(function (err) { console.error('history fetch error', err); });
    }

    var COLORS = [
        '#0d6efd', '#dc3545', '#198754', '#fd7e14',
        '#6f42c1', '#20c997', '#ffc107', '#0dcaf0'
    ];

    function renderChart(data) {
        var canvas = document.getElementById('history-chart');
        var empty = document.getElementById('chart-empty');

        var weightSeries = (data.series || []).filter(function (s) {
            return s.type === 'weight' && s.points && s.points.length > 0;
        });
        var bwSeries = (data.series || []).filter(function (s) {
            return s.type === 'bodyweight' && s.points && s.points.length > 0;
        });
        var plottable = weightSeries.concat(bwSeries);

        if (plottable.length === 0) {
            canvas.style.display = 'none';
            empty.style.display = 'block';
            if (chart) { chart.destroy(); chart = null; }
            return;
        }

        var dateSet = {};
        plottable.forEach(function (s) {
            s.points.forEach(function (p) { dateSet[p.date] = true; });
        });
        var labels = Object.keys(dateSet).sort();

        var datasets = plottable.map(function (s, i) {
            var pointMap = {};
            var repsMap = {};
            s.points.forEach(function (p) {
                pointMap[p.date] = p.value;
                repsMap[p.date] = p.reps;
            });
            var suffix = s.type === 'bodyweight' ? ' (reps)' : ' (' + data.weightUnit + ')';
            var repsData = labels.map(function (d) { return repsMap.hasOwnProperty(d) ? repsMap[d] : 0; });
            return {
                label: s.name + suffix,
                data: labels.map(function (d) { return pointMap.hasOwnProperty(d) ? pointMap[d] : null; }),
                borderColor: COLORS[i % COLORS.length],
                backgroundColor: COLORS[i % COLORS.length] + '33',
                tension: 0.3,
                spanGaps: true,
                yAxisID: s.type === 'bodyweight' ? 'y1' : 'y',
                exerciseType: s.type,
                repsData: repsData,
                pointRadius: repsData.map(function (r) { return r > 0 ? Math.min(4 + r, 20) : 4; }),
                pointHoverRadius: repsData.map(function (r) { return r > 0 ? Math.min(6 + r, 22) : 6; })
            };
        });

        var hasWeight = weightSeries.length > 0;
        var hasBW = bwSeries.length > 0;
        var scales = {};
        if (hasWeight) {
            scales.y = {
                type: 'linear',
                position: 'left',
                title: { display: true, text: data.weightUnit }
            };
        }
        if (hasBW) {
            scales.y1 = {
                type: 'linear',
                position: 'right',
                title: { display: true, text: 'reps' },
                grid: { drawOnChartArea: !hasWeight }
            };
        }

        if (chart) { chart.destroy(); }
        canvas.style.display = 'block';
        empty.style.display = 'none';

        chart = new Chart(canvas, {
            type: 'line',
            data: { labels: labels, datasets: datasets },
            options: {
                responsive: true,
                interaction: { mode: 'index', intersect: false },
                plugins: {
                    legend: { position: 'bottom' },
                    tooltip: {
                        callbacks: {
                            label: function (ctx) {
                                var val = ctx.parsed.y;
                                if (val === null || val === undefined) return null;
                                var display = (val % 1 === 0) ? val.toString() : val.toFixed(1);
                                var label = ctx.dataset.label + ': ' + display;
                                if (ctx.dataset.exerciseType !== 'bodyweight') {
                                    var reps = ctx.dataset.repsData && ctx.dataset.repsData[ctx.dataIndex];
                                    if (reps) { label += ' × ' + reps + ' reps'; }
                                }
                                return label;
                            }
                        }
                    }
                },
                scales: scales
            }
        });
    }

    function showEmpty() {
        var canvas = document.getElementById('history-chart');
        var empty = document.getElementById('chart-empty');
        canvas.style.display = 'none';
        empty.style.display = 'block';
        if (chart) { chart.destroy(); chart = null; }
    }

    // Pre-populate with exercises performed in the last two weeks.
    if (Array.isArray(defaultNames) && defaultNames.length > 0) {
        selected = defaultNames.slice();
        renderChips();
        fetchAndRender();
    }
})();
</script>
<script>
(function () {
    var activity = {{.HeatmapDataJSON}};
    var counts = {};
    (activity || []).forEach(function (d) { counts[d.date] = d.count; });

    var MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
                  'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

    function key(d) {
        var m = String(d.getMonth() + 1).padStart(2, '0');
        var day = String(d.getDate()).padStart(2, '0');
        return d.getFullYear() + '-' + m + '-' + day;
    }
    function human(d) {
        return MONTHS[d.getMonth()] + ' ' + d.getDate() + ', ' + d.getFullYear();
    }
    function level(c) { return c <= 9 ? 1 : c <= 19 ? 2 : c <= 29 ? 3 : 4; }

    var today = new Date();
    today.setHours(0, 0, 0, 0);
    var start = new Date(today);
    start.setDate(start.getDate() - 364);
    start.setDate(start.getDate() - start.getDay()); // back up to Sunday
    var weeks = Math.ceil(((today - start) / 86400000 + 1) / 7);

    var grid = document.getElementById('hm-grid');
    for (var i = 0; i < weeks * 7; i++) {
        var d = new Date(start);
        d.setDate(start.getDate() + i);
        var cell = document.createElement('div');
        cell.className = 'hm-cell';
        if (d > today) {
            cell.style.visibility = 'hidden';
        } else if (counts.hasOwnProperty(key(d))) {
            var c = counts[key(d)];
            cell.classList.add('hm-l' + level(c));
            cell.title = (c === 1 ? '1 set' : c + ' sets') + ' — ' + human(d);
        } else {
            cell.title = 'No workout — ' + human(d);
        }
        grid.appendChild(cell);
    }

    var months = document.getElementById('hm-months');
    var prevMonth = -1;
    for (var w = 0; w < weeks; w++) {
        var wd = new Date(start);
        wd.setDate(start.getDate() + w * 7);
        var span = document.createElement('span');
        if (wd.getMonth() !== prevMonth) {
            span.textContent = MONTHS[wd.getMonth()];
            prevMonth = wd.getMonth();
        }
        months.appendChild(span);
    }

    var n = (activity || []).length;
    document.getElementById('heatmap-summary').textContent =
        n + (n === 1 ? ' day' : ' days') + ' in the last year';
})();
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
