package helper

import (
	"fmt"
	"strings"

	"github.com/kholidss/xyz-skilltest/pkg/util"
	"github.com/spf13/cast"
)

func StructQueryWhereMysql(i interface{}, hideDeleted bool, tag string) (q string, vals []interface{}, limit, page uint64, err error) {

	var cols []string
	var startDate, endDate, sortOrder, groupBy string

	if i == nil {
		return q, vals, limit, page, nil
	}

	data, err := util.StructToMap(i, tag)
	if err != nil {
		return q, vals, limit, page, err
	}

	if len(data) == 0 {
		return q, vals, limit, page, err
	}

	for k, x := range data {
		if k == "page" {
			page = cast.ToUint64(x)
			continue
		}

		if k == "limit" {
			limit = cast.ToUint64(x)
			continue
		}

		if k == "start_date" {
			startDate = cast.ToString(x)
			continue
		}

		if k == "end_date" {
			endDate = cast.ToString(x)
			continue
		}

		if k == "sort_order" {
			sortOrder = cast.ToString(x)
			continue
		}

		if k == "group_by" {
			groupBy = cast.ToString(x)
			continue
		}

		vals = append(vals, x)
		cols = append(cols, k)
	}

	if len(cols) > 0 && !hideDeleted {
		q = fmt.Sprintf(`WHERE %s`, util.StringJoin(cols, "=? AND ", "=?"))
	}

	if len(cols) > 0 && hideDeleted {
		q = fmt.Sprintf(`WHERE %s AND is_deleted = false`, util.StringJoin(cols, "=? AND ", "=?"))
	}

	if len(cols) < 1 && hideDeleted {
		q = fmt.Sprint(`WHERE is_deleted = false`)
	}

	if startDate != "" && endDate != "" {
		q = fmt.Sprintf(`%s AND ( DATE(created_at) >= ?  AND DATE(created_at) <= ? )`, q)
		if hideDeleted {
			q = fmt.Sprintf(`%s AND ( DATE(created_at) >= ?  AND DATE(created_at) <= ? ) AND is_deleted = false`, q)
		}

		if len(cols) == 0 && !hideDeleted {
			q = fmt.Sprint(`WHERE ( DATE(created_at) >= ?  AND DATE(created_at) <= ? )`)
		}

		if len(cols) == 0 && hideDeleted {
			q = fmt.Sprint(`WHERE ( DATE(created_at) >= ?  AND DATE(created_at) <= ? ) AND is_deleted = false`)
		}

		vals = append(vals, startDate, endDate)
	}

	if groupBy != "" {
		q = fmt.Sprintf("%s GROUP BY %s", q, groupBy)
	}

	if sortOrder != "" {
		q = fmt.Sprintf("%s ORDER BY created_at %s", q, sortOrder)
	}

	return q, vals, limit, page, err
}

func StructQueryInsertMysql(param interface{}, tableName, tag string, returningID bool) (string, []interface{}, error) {
	var (
		keys   []string
		values []interface{}
		numArr []string
	)

	resMap, err := util.StructToMap(param, tag)
	if err != nil {
		return "", nil, err
	}

	for k, v := range resMap {
		keys = append(keys, k)
		values = append(values, v)
		numArr = append(numArr, "?")
	}

	var (
		q         string
		returnIDQ = func() string {
			if returningID {
				return "RETURNING id"
			}
			return ""
		}()
	)

	q = fmt.Sprintf(`
		INSERT INTO
			%s
		(
			%s
		)
		VALUES
		(
			%s
		)
		%s;
	`, tableName, strings.Join(keys, ","), strings.Join(numArr, ","), returnIDQ)

	return q, values, nil
}

func StructToQueryUpdateMysql(input interface{}, where interface{}, tableName, tag string) (string, []interface{}, error) {

	cols, vals, err := util.ToColumnsValues(input, tag)
	if err != nil {
		return "", vals, err
	}

	cu, vu, err := util.ToColumnsValues(where, tag)
	if err != nil {
		return "", vals, err
	}

	q := fmt.Sprintf(`UPDATE %s SET %s`, tableName, util.StringJoin(cols, "=?, ", "=?"))
	if len(cu) > 0 {
		q = fmt.Sprintf(`%s WHERE %s`, q, util.StringJoin(cu, "=? AND ", "=?"))
		vals = append(vals, vu...)
	}

	return q, vals, nil
}

func SelectCustom(selectColumn []string) string {
	if len(selectColumn) > 0 {
		return strings.Join(selectColumn, ",")
	}
	return "*"
}
