package queries

import (
	"fmt"
	"strings"
)

func (q PsQuery) Build(query string) (string, []interface{}) {

	equalsQueriesCols := []string{}
	equalsQueriesVals := []interface{}{}
	assignNum := 1
	for _, eq := range q.Equals {
		var dbcol string
		switch eq.Field {
		default:
			dbcol = fmt.Sprintf("%s=$%d", eq.Field, assignNum)
		}
		assignNum += 1
		equalsQueriesCols = append(equalsQueriesCols, dbcol)
		equalsQueriesVals = append(equalsQueriesVals, eq.Value)
	}
	if len(q.Date) == 2 {
		col := fmt.Sprintf("date_created between $%d and $%d;", assignNum, assignNum+1)
		equalsQueriesCols = append(equalsQueriesCols, col)
		vals := formatDate(q.Date)
		equalsQueriesVals = append(equalsQueriesVals, vals[0], vals[1])
	}
	query = fmt.Sprintf(query, strings.Join(equalsQueriesCols, " AND "))
	return query, equalsQueriesVals
}

func formatDate(dates []DateValue) [2]string {
	var stringDates [2]string
	for i, vals := range dates {
		stringDates[i] = fmt.Sprintf("%s-%s-%s", vals.Year, vals.Month, vals.Day)
	}
	return stringDates
}
