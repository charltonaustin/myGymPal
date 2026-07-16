package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// DayActivity is the aggregated workout volume for a single calendar day.
type DayActivity struct {
	Date  string `orm:"column(day)"       json:"date"`
	Count int    `orm:"column(set_count)" json:"count"`
}

// dailyActivitySQL counts logged sets per day for one user within a date window.
// A day appears in the result if it has a session, even with zero sets logged
// (e.g. a cardio-only day), so callers can treat any returned day as "exercised".
const dailyActivitySQL = `
SELECT TO_CHAR(s.date, 'YYYY-MM-DD') AS day, COUNT(ss.id) AS set_count
FROM sessions s
LEFT JOIN session_exercises se ON se.session_id = s.id
LEFT JOIN session_sets ss ON ss.session_exercise_id = se.id
WHERE s.user_id = ?
  AND s.date >= ?::date
GROUP BY s.date
ORDER BY s.date ASC
`

// GetDailyActivity returns the days the user exercised within the last `days`
// days, each with the number of sets logged that day (its heatmap intensity).
func GetDailyActivity(userID int64, days int) ([]DayActivity, error) {
	o := orm.NewOrm()
	cutoff := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	var rows []*DayActivity
	if _, err := o.Raw(dailyActivitySQL, userID, cutoff).QueryRows(&rows); err != nil {
		return nil, err
	}

	result := make([]DayActivity, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}
