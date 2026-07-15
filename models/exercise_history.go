package models

import "github.com/beego/beego/v2/client/orm"

// ExerciseHistoryPoint is one session's aggregated best value for one exercise.
type ExerciseHistoryPoint struct {
	Date  string  `orm:"column(session_date)" json:"date"`
	Value float64 `orm:"column(max_value)"    json:"value"`
}

// ExerciseHistorySeries holds the full time series for one exercise.
type ExerciseHistorySeries struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"` // "weight", "bodyweight", or "time_based"
	Unit   string                 `json:"unit"` // "lb", "kg", or "reps"
	Points []ExerciseHistoryPoint `json:"points"`
}

const weightHistorySQL = `
SELECT
    TO_CHAR(s.date, 'YYYY-MM-DD') AS session_date,
    MAX(CASE ss.weight_unit WHEN 'kg' THEN ss.actual_weight * 2.20462 ELSE ss.actual_weight END) AS max_value
FROM sessions s
JOIN session_exercises se ON se.session_id = s.id
JOIN session_sets ss ON ss.session_exercise_id = se.id
WHERE s.user_id = ?
  AND LOWER(TRIM(se.name)) = LOWER(TRIM(?))
  AND se.is_bodyweight = false
  AND se.is_time_based = false
  AND ss.actual_weight > 0
GROUP BY s.date
ORDER BY s.date ASC
`

const bodyweightHistorySQL = `
SELECT
    TO_CHAR(s.date, 'YYYY-MM-DD') AS session_date,
    MAX(ss.actual_reps) AS max_value
FROM sessions s
JOIN session_exercises se ON se.session_id = s.id
JOIN session_sets ss ON ss.session_exercise_id = se.id
WHERE s.user_id = ?
  AND LOWER(TRIM(se.name)) = LOWER(TRIM(?))
  AND se.is_bodyweight = true
  AND ss.actual_reps > 0
GROUP BY s.date
ORDER BY s.date ASC
`

// GetExerciseHistory returns per-session max weight (or max reps for bodyweight) for each named exercise.
// Weights are normalized to targetUnit; bodyweight exercises always report reps.
// Time-based exercises are returned with Type "time_based" and no data points.
func GetExerciseHistory(userID int64, names []string, targetUnit string) ([]ExerciseHistorySeries, error) {
	o := orm.NewOrm()
	result := make([]ExerciseHistorySeries, 0, len(names))

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
			_, err = o.Raw(bodyweightHistorySQL, userID, name).QueryRows(&points)
		} else {
			series.Type = "weight"
			series.Unit = targetUnit
			_, err = o.Raw(weightHistorySQL, userID, name).QueryRows(&points)
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
