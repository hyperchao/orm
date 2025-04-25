package orm

import (
	"database/sql/driver"
	"github.com/hyperchao/orm/tag"
	"reflect"
	"slices"
	"strings"
)

type empty struct{}

func (empty) Scan(any) error {
	return nil
}

type columnAttr int

const (
	columnAttrPrimary columnAttr = 1 << iota
	columnAttrAutoincrement
	columnAttrOptimisticLock
)

func (c columnAttr) Has(attr columnAttr) bool {
	return c&attr != 0
}

const (
	separator   = ","
	placeholder = "?"
	quote       = "`"
)

var (
	tagParser = tag.NewParser(func(tagValue string) (field string, attributes columnAttr) {
		// tagValue example. `orm:"id,primary,autoincrement"`
		parts := strings.Split(tagValue, ",")
		field = parts[0]
		for _, attr := range parts[1:] {
			if strings.TrimSpace(attr) == TagPrimaryKey {
				attributes |= columnAttrPrimary
			}
			if strings.TrimSpace(attr) == TagAutoIncrement {
				attributes |= columnAttrAutoincrement
			}
			if strings.TrimSpace(attr) == TagVersion {
				attributes |= columnAttrOptimisticLock
			}
		}
		return
	})
)

// extract a slice of interfaces from struct for sql.Rows.Scan to use
func getColumnDest(conf *config, val any, columns []string) []any {
	if len(columns) == 0 {
		return nil
	}

	values := tagParser.Parse(conf.tagName, val)

	r := make([]any, len(columns))
	for i, col := range columns {
		if values.Contains(col) {
			r[i] = values.Get(col).Addr()
		} else {
			r[i] = empty{}
		}
	}
	return r
}

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

func parseInsertColumnsAndArgs(values tag.Values[columnAttr]) (columns []string, autoincrement string, args []any) {
	if values.Len() == 0 {
		return
	}

	columns = make([]string, 0, values.Len())
	args = make([]any, 0, values.Len())
	for field, value := range values.Iter() {
		if value.Meta().Attrs().Has(columnAttrAutoincrement) {
			autoincrement = field
		} else {
			columns = append(columns, field)
			args = append(args, value.Interface())
		}
	}

	return
}

func generateInsertSQL(tableName string, columns []string, count int) string {
	sb := strings.Builder{}
	sb.WriteString("INSERT INTO ")
	sb.WriteString(tableName)
	writeSQLColumns(&sb, columns)
	sb.WriteString(" VALUES ")
	writeSQLPlaceholders(&sb, len(columns))
	for count > 1 {
		sb.WriteString(",")
		writeSQLPlaceholders(&sb, len(columns))
		count--
	}
	return sb.String()
}

func writeSQLColumns(sb *strings.Builder, slice []string) {
	if len(slice) == 0 {
		return
	}
	sb.WriteString("(")
	sb.WriteString(slice[0])
	for i := 1; i < len(slice); i++ {
		sb.WriteString(separator)
		sb.WriteString(quote)
		sb.WriteString(slice[i])
		sb.WriteString(quote)
	}
	sb.WriteString(")")
}

func writeSQLPlaceholders(sb *strings.Builder, n int) {
	if n <= 0 {
		return
	}
	sb.WriteString("(")
	sb.WriteString(placeholder)
	for i := 1; i < n; i++ {
		sb.WriteString(separator)
		sb.WriteString(placeholder)
	}
	sb.WriteString(")")
}
