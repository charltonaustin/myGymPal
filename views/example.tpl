<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>See It In Action — My Gym Pal</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css" rel="stylesheet">
    <link rel="manifest" href="/manifest.json">
</head>
<body>

{{template "partials/navbar.tpl" .}}

<main class="container mt-4 mb-5" style="max-width: 640px;">
    <div class="mb-4">
        <a href="/" class="text-muted small">&larr; Home</a>
        <h1 class="h4 fw-bold mt-1 mb-1">See it in action</h1>
        <p class="text-muted small mb-0">A peek at what My Gym Pal looks like for a real user. Jack has been training for 8 weeks.</p>
    </div>

    <!-- Dashboard -->
    <h2 class="h6 fw-semibold text-uppercase text-muted mb-2">Dashboard</h2>
    <div class="card mb-4">
        <div class="card-body">
            <h3 class="h5 fw-bold mb-3">Welcome back, Jack!</h3>
            <div class="row g-3 mb-3">
                <div class="col-6">
                    <div class="card h-100 border">
                        <div class="card-body py-2 px-3">
                            <div class="d-flex justify-content-between align-items-start mb-1">
                                <span class="h6 fw-semibold mb-0">Weight</span>
                                <span class="text-muted small">Log Weight</span>
                            </div>
                            <p class="fs-4 fw-bold mb-0">203.4 <span class="fs-6 fw-normal text-muted">lb</span></p>
                            <p class="text-muted small mb-0">3-day avg</p>
                        </div>
                    </div>
                </div>
                <div class="col-6">
                    <div class="card h-100 border">
                        <div class="card-body py-2 px-3">
                            <div class="d-flex justify-content-between align-items-start mb-1">
                                <span class="h6 fw-semibold mb-0">Macros</span>
                                <span class="text-muted small">Log Macros</span>
                            </div>
                            <p class="text-muted small mb-1">3-day avg</p>
                            <table class="table table-sm mb-0">
                                <tbody class="small">
                                    <tr>
                                        <td class="ps-0 text-muted">Protein</td>
                                        <td class="text-end fw-semibold">196g</td>
                                        <td class="text-end fw-semibold text-success">98%</td>
                                    </tr>
                                    <tr>
                                        <td class="ps-0 text-muted">Carbs</td>
                                        <td class="text-end fw-semibold">238g</td>
                                        <td class="text-end fw-semibold text-danger">79%</td>
                                    </tr>
                                    <tr>
                                        <td class="ps-0 text-muted">Fat</td>
                                        <td class="text-end fw-semibold">68g</td>
                                        <td class="text-end fw-semibold text-success">102%</td>
                                    </tr>
                                    <tr class="border-top">
                                        <td class="ps-0 text-muted">Calories</td>
                                        <td class="text-end fw-semibold">2,316</td>
                                        <td class="text-end fw-semibold text-success">96%</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
            <div class="d-flex justify-content-between align-items-baseline mb-2">
                <span class="h6 fw-semibold text-uppercase text-muted mb-0" style="font-size: 0.75rem;">Recent Workouts</span>
                <span class="text-muted small">Log next workout</span>
            </div>
            <div class="list-group">
                <div class="list-group-item d-flex justify-content-between align-items-center">
                    <div>
                        <div class="fw-semibold">Strength Block A</div>
                        <div class="text-muted small">Phase 2 &middot; Week 3 &middot; Workout 2</div>
                    </div>
                    <span class="text-muted small">Mar 14, 2026</span>
                </div>
                <div class="list-group-item d-flex justify-content-between align-items-center">
                    <div>
                        <div class="fw-semibold">Strength Block A</div>
                        <div class="text-muted small">Phase 2 &middot; Week 3 &middot; Workout 1</div>
                    </div>
                    <span class="text-muted small">Mar 12, 2026</span>
                </div>
                <div class="list-group-item d-flex justify-content-between align-items-center">
                    <div>
                        <div class="fw-semibold">Strength Block A</div>
                        <div class="text-muted small">Phase 2 &middot; Week 2 &middot; Workout 3</div>
                    </div>
                    <span class="text-muted small">Mar 10, 2026</span>
                </div>
            </div>
        </div>
    </div>

    <!-- Active session -->
    <h2 class="h6 fw-semibold text-uppercase text-muted mb-2">Logging a Workout</h2>
    <div class="card mb-4">
        <div class="card-body">
            <a class="text-muted small d-block mb-1">&larr; Strength Block A</a>
            <div class="d-flex align-items-center gap-2 mb-0">
                <h3 class="h5 fw-bold mb-0">Session #2</h3>
            </div>
            <p class="text-muted small mb-3">Phase 2 &middot; Week 3 &middot; Mar 14, 2026</p>

            <!-- Exercise 1 -->
            <div class="card mb-3 border">
                <div class="card-body pb-2">
                    <div class="d-flex align-items-baseline justify-content-between mb-2">
                        <span class="h6 fw-semibold mb-0">Bench Press</span>
                        <span class="text-muted small">Goal: 185 lb &nbsp; 6–8 reps</span>
                    </div>
                    <table class="table table-sm mb-2">
                        <thead>
                            <tr>
                                <th class="text-muted fw-normal small ps-0">Set</th>
                                <th class="text-muted fw-normal small">Weight</th>
                                <th class="text-muted fw-normal small">Reps</th>
                                <th></th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr><td class="ps-0">1</td><td>185 lb</td><td>8</td><td class="text-end"><i class="bi bi-trash text-danger" style="font-size: 0.8rem;"></i></td></tr>
                            <tr><td class="ps-0">2</td><td>185 lb</td><td>8</td><td class="text-end"><i class="bi bi-trash text-danger" style="font-size: 0.8rem;"></i></td></tr>
                            <tr><td class="ps-0">3</td><td>185 lb</td><td>7</td><td class="text-end"><i class="bi bi-trash text-danger" style="font-size: 0.8rem;"></i></td></tr>
                        </tbody>
                    </table>
                    <div class="d-flex gap-2 align-items-end">
                        <div>
                            <label class="form-label small mb-1">Weight</label>
                            <div class="input-group input-group-sm" style="width: 160px;">
                                <input type="number" class="form-control" value="185" disabled>
                                <span class="input-group-text">lb</span>
                            </div>
                        </div>
                        <div>
                            <label class="form-label small mb-1">Reps</label>
                            <input type="number" class="form-control form-control-sm" value="8" disabled style="width: 70px;">
                        </div>
                        <button class="btn btn-dark btn-sm mb-0" disabled>+ Set</button>
                    </div>
                </div>
            </div>

            <!-- Exercise 2 — bodyweight -->
            <div class="card mb-3 border">
                <div class="card-body pb-2">
                    <div class="d-flex align-items-baseline justify-content-between mb-2">
                        <span class="h6 fw-semibold mb-0">Pull-up</span>
                        <span class="text-muted small">6–8 reps</span>
                    </div>
                    <table class="table table-sm mb-2">
                        <thead>
                            <tr>
                                <th class="text-muted fw-normal small ps-0">Set</th>
                                <th class="text-muted fw-normal small">Weight</th>
                                <th class="text-muted fw-normal small">Reps</th>
                                <th></th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr><td class="ps-0">1</td><td>0 lb</td><td>8</td><td class="text-end"><i class="bi bi-trash text-danger" style="font-size: 0.8rem;"></i></td></tr>
                            <tr><td class="ps-0">2</td><td>0 lb</td><td>7</td><td class="text-end"><i class="bi bi-trash text-danger" style="font-size: 0.8rem;"></i></td></tr>
                        </tbody>
                    </table>
                    <div class="d-flex gap-2 align-items-end">
                        <div>
                            <label class="form-label small mb-1">Weight</label>
                            <div class="input-group input-group-sm" style="width: 160px;">
                                <input type="number" class="form-control" value="0" disabled>
                                <span class="input-group-text">lb</span>
                            </div>
                        </div>
                        <div>
                            <label class="form-label small mb-1">Reps</label>
                            <input type="number" class="form-control form-control-sm" value="8" disabled style="width: 70px;">
                        </div>
                        <button class="btn btn-dark btn-sm mb-0" disabled>+ Set</button>
                    </div>
                </div>
            </div>

            <!-- Cardio section -->
            <h4 class="h6 fw-semibold text-uppercase text-muted mt-4 mb-2" style="font-size: 0.75rem;">Cardio</h4>
            <div class="card mb-3 border">
                <div class="card-body py-2">
                    <div class="d-flex align-items-center justify-content-between">
                        <div>
                            <span class="fw-semibold text-capitalize small">run</span>
                            <span class="text-muted small ms-2 text-capitalize">Steady State</span>
                        </div>
                        <span class="text-muted small">Goal: 30 | Actual: 32 min</span>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Exercise library -->
    <h2 class="h6 fw-semibold text-uppercase text-muted mb-2">Exercise Library</h2>
    <div class="card mb-4">
        <div class="card-body">
            <p class="text-muted small mb-3">Jack keeps a library of exercises with goal weights. After 3 sets at or above goal reps, the goal updates automatically.</p>
            <div class="list-group">
                <div class="list-group-item d-flex justify-content-between align-items-center">
                    <div>
                        <div class="fw-semibold text-capitalize">Bench Press</div>
                        <div class="text-muted small">Goal: 185 lb</div>
                    </div>
                    <i class="bi bi-pencil text-secondary" style="font-size: 0.85rem;"></i>
                </div>
                <div class="list-group-item d-flex justify-content-between align-items-center">
                    <div>
                        <div class="fw-semibold text-capitalize">Squat <span class="badge bg-secondary ms-1" style="font-size: 0.65rem;">Bodyweight</span></div>
                    </div>
                    <i class="bi bi-pencil text-secondary" style="font-size: 0.85rem;"></i>
                </div>
                <div class="list-group-item d-flex justify-content-between align-items-center">
                    <div>
                        <div class="fw-semibold text-capitalize">Plank <span class="badge bg-info ms-1" style="font-size: 0.65rem;">Time-based</span></div>
                        <div class="text-muted small">Goal: 1:30</div>
                    </div>
                    <i class="bi bi-pencil text-secondary" style="font-size: 0.85rem;"></i>
                </div>
                <div class="list-group-item d-flex justify-content-between align-items-center">
                    <div>
                        <div class="fw-semibold text-capitalize">Overhead Press</div>
                        <div class="text-muted small">Goal: 115 lb</div>
                    </div>
                    <i class="bi bi-pencil text-secondary" style="font-size: 0.85rem;"></i>
                </div>
            </div>
        </div>
    </div>

    <div class="text-center pt-2">
        <p class="text-muted mb-3">Ready to start tracking your own progress?</p>
        <a href="/register" class="btn btn-dark">Create a Free Account</a>
    </div>
</main>

<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"></script>
<script>if ('serviceWorker' in navigator) { navigator.serviceWorker.register('/sw.js').catch(console.error); }</script>
</body>
</html>
