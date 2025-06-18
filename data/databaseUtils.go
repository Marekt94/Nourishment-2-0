package database

import (
	"database/sql"
	"log"
	"strings"
)

func NullStringToString(n *sql.NullString) string{
	if n.Valid{
		return n.String
	}
	return ""
}	

func NullInt64ToInt(n *sql.NullInt64) int{
	if n.Valid{
		return int(n.Int64)
	}
	return 0
}

func NullFloat64ToFloat(n *sql.NullFloat64) float64{
	if n.Valid{
		return n.Float64
	}
	return 0
}

func QuestionMarks(n int) string {
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

func CreateColsToSelect(prefix string, c []string) string{
    res :=  prefix + `.` + strings.Join(c[:], `, ` + prefix + `.`)
    log.Println(res)
    return res
}