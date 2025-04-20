package orm

import (
	"database/sql/driver"
	"reflect"
	"slices"
	"strings"
)

// RewriteQueryAndArgs transform a slice argument to a list of arguments and rewrite the "?" in query to "(?,?,...)"
// so we can write sql like this:
//
//	query := "select * from userinfo where uid in ? and state = ?"
//	args := []int64{1,2,3}, 1
//
// after transform, actual sql will be:
//
//	query := "select * from userinfo where uid in (?,?,?) and state = ?"
//	args := 1, 2, 3, 1
//
// empty slice will rewrite "?" in query to  "(NULL)"
// take care of the behavior. especially when you use "not in" clause
func RewriteQueryAndArgs(query string, args ...any) (rewrittenQuery string, expandedArgs []any) {
	sliceIndexes := make([]int, 0)
	for i, arg := range args {
		if _, ok := arg.(driver.Valuer); ok {
			continue
		}
		rt := reflect.TypeOf(arg)
		if rt != nil && rt.Kind() == reflect.Slice {
			sliceIndexes = append(sliceIndexes, i)
		}
	}

	if len(sliceIndexes) == 0 {
		return query, args
	}

	queryParts := strings.Split(query, "?")
	sb := strings.Builder{}
	for idx, part := range queryParts {
		sb.WriteString(part)
		if idx != len(queryParts)-1 {
			_, found := slices.BinarySearch(sliceIndexes, idx)
			if found {
				sv := reflect.ValueOf(args[idx])
				if sv.Len() == 0 {
					sb.WriteString("(NULL)")
				} else {
					sb.WriteString("(")
					for i := 0; i < sv.Len(); i++ {
						expandedArgs = append(expandedArgs, sv.Index(i).Interface())
						sb.WriteString("?")
						if i != sv.Len()-1 {
							sb.WriteString(",")
						}
					}
					sb.WriteString(")")
				}
			} else {
				expandedArgs = append(expandedArgs, args[idx])
				sb.WriteString("?")
			}
		}
	}

	rewrittenQuery = sb.String()
	return
}
