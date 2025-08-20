package database

import (
	"database/sql"
	"nourishment_20/internal/logging"
	"strings"
)

// TODO te metody powinny byc metodami silnika, a nie luzem, poniewaz np: w innej bazie danych  QuestionMarks moze wygladac inaczej
func NullStringToString(n *sql.NullString) string { // [AI REFACTOR]
	if n.Valid {
		return n.String
	}
	return ""
}

func NullInt64ToInt(n *sql.NullInt64) int { // [AI REFACTOR]
	if n.Valid {
		return int(n.Int64)
	}
	return 0
}

func NullFloat64ToFloat(n *sql.NullFloat64) float64 { // [AI REFACTOR]
	if n.Valid {
		return n.Float64
	}
	return 0
}

func QuestionMarks(n int) string { // [AI REFACTOR]
	if n <= 0 {
		return ""
	}
	marks := make([]string, n)
	for i := range marks {
		marks[i] = "?"
	}
	return strings.Join(marks, ", ")
}

func UpdateValues(v []string) string {
	if len(v) == 0 {
		return ""
	}
	updateValues := make([]string, len(v))
	for i := range updateValues {
		updateValues[i] = v[i] + `=?`
	}
	return strings.Join(updateValues, ", ")
}

func CreateColsToSelect(prefix string, c []string) string {
	res := prefix + `.` + strings.Join(c[:], `, `+prefix+`.`)
	logging.Global.Debugf(res)
	return res
}
