package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// ExerciseHistoryPoint is one session's aggregated best value for one exercise.
type ExerciseHistoryPoint struct {
	Date  string  `orm:"column(session_date)" json:"date"`
	Value float64 `orm:"column(max_value)"    json:"value"`
	Reps  int     `orm:"column(reps)"         json:"reps"`
}

// ExerciseHistorySeries holds the full time series for one exercise.
type ExerciseHistorySeries struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"` // "weight", "bodyweight", or "time_based"
	Unit   string                 `json:"unit"` // "lb", "kg", or "reps"
	Points []ExerciseHistoryPoint `json:"points"`
}

const weightHistorySQL = `
WITH ranked AS (
    SELECT
        TO_CHAR(s.date, 'YYYY-MM-DD') AS session_date,
        CASE ss.weight_unit WHEN 'kg' THEN ss.actual_weight * 2.20462 ELSE ss.actual_weight END AS normalized_weight,
        ss.actual_reps,
        ROW_NUMBER() OVER (
            PARTITION BY s.date
            ORDER BY CASE ss.weight_unit WHEN 'kg' THEN ss.actual_weight * 2.20462 ELSE ss.actual_weight END DESC
        ) AS rn
    FROM sessions s
    JOIN session_exercises se ON se.session_id = s.id
    JOIN session_sets ss ON ss.session_exercise_id = se.id
    WHERE s.user_id = ?
      AND LOWER(TRIM(se.name)) = LOWER(TRIM(?))
      AND se.is_bodyweight = false
      AND se.is_time_based = false
      AND ss.actual_weight > 0
      AND s.date >= ?::date
)
SELECT session_date, normalized_weight AS max_value, actual_reps AS reps
FROM ranked
WHERE rn = 1
ORDER BY session_date ASC
`

const bodyweightHistorySQL = `
WITH ranked AS (
    SELECT
        TO_CHAR(s.date, 'YYYY-MM-DD') AS session_date,
        ss.actual_reps,
        ROW_NUMBER() OVER (PARTITION BY s.date ORDER BY ss.actual_reps DESC) AS rn
    FROM sessions s
    JOIN session_exercises se ON se.session_id = s.id
    JOIN session_sets ss ON ss.session_exercise_id = se.id
    WHERE s.user_id = ?
      AND LOWER(TRIM(se.name)) = LOWER(TRIM(?))
      AND se.is_bodyweight = true
      AND ss.actual_reps > 0
      AND s.date >= ?::date
)
SELECT session_date, actual_reps AS max_value, actual_reps AS reps
FROM ranked
WHERE rn = 1
ORDER BY session_date ASC
`

// recentExerciseNamesSQL returns the distinct names of exercises the user has
// actually logged sets for within the given date window, most-recently-performed first.
const recentExerciseNamesSQL = `
SELECT se.name, MAX(s.date) AS last_done
FROM sessions s
JOIN session_exercises se ON se.session_id = s.id
WHERE s.user_id = ?
  AND s.date >= ?
  AND EXISTS (SELECT 1 FROM session_sets ss WHERE ss.session_exercise_id = se.id)
GROUP BY se.name
ORDER BY last_done DESC, se.name ASC
`

// GetRecentExerciseNames returns the names of exercises the user has performed
// (logged at least one set for) within the last `days` days, ordered by how
// recently they were last done. Used to pre-populate the history graph.
func GetRecentExerciseNames(userID int64, days int) ([]string, error) {
	o := orm.NewOrm()
	cutoff := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	var rows []struct {
		Name     string    `orm:"column(name)"`
		LastDone time.Time `orm:"column(last_done)"`
	}
	if _, err := o.Raw(recentExerciseNamesSQL, userID, cutoff).QueryRows(&rows); err != nil {
		return nil, err
	}

	names := make([]string, len(rows))
	for i, r := range rows {
		names[i] = r.Name
	}
	return names, nil
}

// GetExerciseHistory returns per-session max weight (or max reps for bodyweight) for each named exercise.
// Weights are normalized to targetUnit; bodyweight exercises always report reps.
// Time-based exercises are returned with Type "time_based" and no data points.
// Only sessions within the last `days` days are included; days <= 0 means no date limit.
func GetExerciseHistory(userID int64, names []string, targetUnit string, days int) ([]ExerciseHistorySeries, error) {
	o := orm.NewOrm()
	result := make([]ExerciseHistorySeries, 0, len(names))

	// Sentinel cutoff older than any data acts as a no-op when no window is requested.
	cutoff := "0001-01-01"
	if days > 0 {
		cutoff = time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	}

	for _, name := range names {
		ex, err := GetExerciseByName(userID, name)
		if err != nil || ex.IsTimeBased {
			result = append(result, ExerciseHistorySeries{Name: name, Type: "time_based", Unit: "s", Points: nil})
			continue
		}

		series := ExerciseHistorySeries{Name: name}
		var points []*ExerciseHistoryPoint

		if ex.IsBodyweight {
			series.Type = "bodyweight"
			series.Unit = "reps"
			_, err = o.Raw(bodyweightHistorySQL, userID, name, cutoff).QueryRows(&points)
		} else {
			series.Type = "weight"
			series.Unit = targetUnit
			_, err = o.Raw(weightHistorySQL, userID, name, cutoff).QueryRows(&points)
			if err == nil && targetUnit != "lb" {
				for _, p := range points {
					p.Value = ConvertWeight(p.Value, "lb", targetUnit)
				}
			}
		}

		if err != nil {
			return nil, err
		}

		series.Points = make([]ExerciseHistoryPoint, len(points))
		for i, p := range points {
			series.Points[i] = *p
		}
		result = append(result, series)
	}

	return result, nil
}
