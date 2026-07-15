<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Exercise History — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4 mb-4" style="max-width: 720px;">
    <div class="d-flex justify-content-between align-items-center mb-4">
        <div>
            <a href="/exercises" class="text-muted small">&larr; Exercise Library</a>
            <h1 class="h4 fw-bold mt-1 mb-0">Exercise History</h1>
        </div>
        <div class="btn-group btn-group-sm" role="group" aria-label="Weight unit">
            <input type="radio" class="btn-check" name="unit_toggle" id="unit_lb" value="lb" autocomplete="off" {{if eq .WeightUnit "lb"}}checked{{end}}>
            <label class="btn btn-outline-secondary" for="unit_lb">lb</label>
            <input type="radio" class="btn-check" name="unit_toggle" id="unit_kg" value="kg" autocomplete="off" {{if eq .WeightUnit "kg"}}checked{{end}}>
            <label class="btn btn-outline-secondary" for="unit_kg">kg</label>
        </div>
    </div>

    <div class="card mb-3">
        <div class="card-body">
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
    let selected = [];
    let currentUnit = document.querySelector('input[name="unit_toggle"]:checked').value;
    let chart = null;

    const searchInput = document.getElementById('exercise-search');
    makeAutocomplete(searchInput, exerciseNames, function (name) {
        addExercise(name);
        searchInput.value = '';
    });

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
        const params = new URLSearchParams({ names: selected.join(','), unit: currentUnit });
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
            s.points.forEach(function (p) { pointMap[p.date] = p.value; });
            var suffix = s.type === 'bodyweight' ? ' (reps)' : ' (' + data.weightUnit + ')';
            return {
                label: s.name + suffix,
                data: labels.map(function (d) { return pointMap.hasOwnProperty(d) ? pointMap[d] : null; }),
                borderColor: COLORS[i % COLORS.length],
                backgroundColor: COLORS[i % COLORS.length] + '33',
                tension: 0.3,
                spanGaps: true,
                yAxisID: s.type === 'bodyweight' ? 'y1' : 'y'
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
                                return ctx.dataset.label + ': ' + display;
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
})();
</script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
