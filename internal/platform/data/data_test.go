package data

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOverwriteStruct_Successful(t *testing.T) {
	// Given
	type testA struct {
		Value1 int
		Value2 float32
		Value3 bool
		Value4 string
		Value5 int
		Value6 bool
		Value7 *string
	}

	type testB struct {
		Value1 int
		Value2 float32
		Value3 bool
		Value4 string
		Value5 string
		Value6 bool
		Value7 *string
	}

	var a testA

	b := testB{
		Value1: 5,
		Value2: 1.2,
		Value3: true,
		Value4: "string",
		Value5: "not_possible",
	}

	columns := []string{"Value1", "Value2", "Value3", "Value4", "Value5", "Value6", "Value7"}

	// When
	OverwriteStruct(&a, b, columns)

	// Then
	require.Equal(t, b.Value1, a.Value1)
}
