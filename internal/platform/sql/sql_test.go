package sql

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDynamicQuery_Successful(t *testing.T) {
	// Given
	type testStruct struct {
		Test         string  `db:"test"`
		TestValue    int     `db:"test_value"`
		TestForSQL   bool    `db:"test_for_sql"`
		TestForMySQL float32 `db:"test_for_mysql"`
	}

	testValue := testStruct{
		Test:         "Test String",
		TestValue:    1,
		TestForSQL:   true,
		TestForMySQL: 1.2,
	}

	columns := []string{"test", "test_value", "test_for_sql", "test_for_mysql"}

	// When
	dynamicQuery, values := DynamicQuery(columns, testValue)

	// Then
	require.NotEmpty(t, dynamicQuery)
	require.NotEmpty(t, values)
	require.Contains(t, dynamicQuery, "test = ?")
	require.Contains(t, dynamicQuery, "test_value = ?")
	require.Contains(t, dynamicQuery, "test_for_sql = ?")
	require.Contains(t, dynamicQuery, "test_for_mysql = ?")
}

func TestDynamicQuery_SuccessfulWithOnlyTwoColumns(t *testing.T) {
	// Given
	type testStruct struct {
		Test         string  `db:"test"`
		TestValue    int     `db:"test_value"`
		TestForSQL   bool    `db:"test_for_sql"`
		TestForMySQL float32 `db:"test_for_mysql"`
	}

	testValue := testStruct{
		Test:         "Test String",
		TestValue:    1,
		TestForSQL:   false,
		TestForMySQL: 1.2,
	}

	columns := []string{"test", "test_value"}

	// When
	dynamicQuery, values := DynamicQuery(columns, testValue)

	// Then
	require.NotEmpty(t, dynamicQuery)
	require.NotEmpty(t, values)
	require.Contains(t, dynamicQuery, "test = ?")
	require.Contains(t, dynamicQuery, "test_value = ?")
}
