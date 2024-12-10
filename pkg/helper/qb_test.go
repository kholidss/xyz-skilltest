package helper

import (
	"reflect"
	"testing"
)

func TestStructToQueryUpdate(t *testing.T) {
	type UpdateData struct {
		Name  string `db:"name"`
		Email string `db:"email"`
	}

	type WhereData struct {
		ID int `db:"id"`
	}

	updateInput := UpdateData{Name: "John", Email: "john@example.com"}
	whereInput := WhereData{ID: 1}
	tableName := "users"
	tag := "db"

	expectedQuery := "UPDATE users SET name=?, email=? WHERE id=?"
	expectedValues := []interface{}{"John", "john@example.com", 1}

	query, values, err := StructToQueryUpdateMysql(updateInput, whereInput, tableName, tag)

	if err != nil {
		t.Errorf("Expected no error, but got error: %v", err)
	}

	if query != expectedQuery {
		t.Errorf("Query does not match. Expected: %s, Actual: %s", expectedQuery, query)
	}

	if !reflect.DeepEqual(values, expectedValues) {
		t.Errorf("Values do not match. Expected: %v, Actual: %v", expectedValues, values)
	}
}
