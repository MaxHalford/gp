package op

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestHeight(t *testing.T) {
	var testCases = []struct {
		op Operator
		d  uint
	}{
		{
			op: Const{42},
			d:  0,
		},
		{
			op: Cos{Const{42}},
			d:  1,
		},
		{
			op: Mul{Cos{Const{42}}, Cos{Const{42}}},
			d:  2,
		},
		{
			op: Mul{Cos{Const{42}}, Cos{Cos{Const{42}}}},
			d:  3,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			d := CalcHeight(tc.op)
			if d != tc.d {
				t.Errorf("Expected %d, got %d", tc.d, d)
			}
		})
	}
}

func TestCountOps(t *testing.T) {
	var testCases = []struct {
		op Operator
		n  uint
	}{
		{
			op: Const{42},
			n:  1,
		},
		{
			op: Cos{Const{42}},
			n:  2,
		},
		{
			op: Mul{Cos{Const{42}}, Cos{Const{42}}},
			n:  5,
		},
		{
			op: Mul{Cos{Const{42}}, Cos{Cos{Const{42}}}},
			n:  6,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			n := CountOps(tc.op)
			if n != tc.n {
				t.Errorf("Expected %d, got %d", tc.n, n)
			}
		})
	}
}

func TestSelect(t *testing.T) {
	var testCases = []struct {
		in  Operator
		pos uint
		out Operator
	}{
		{
			in:  Const{42},
			pos: 0,
			out: Const{42},
		},
		{
			in:  Cos{Const{42}},
			pos: 0,
			out: Cos{Const{42}},
		},
		{
			in:  Cos{Const{42}},
			pos: 1,
			out: Const{42},
		},
		{
			in:  Add{Const{42}, Var{1}},
			pos: 0,
			out: Add{Const{42}, Var{1}},
		},
		{
			in:  Add{Const{42}, Var{1}},
			pos: 1,
			out: Const{42},
		},
		{
			in:  Add{Const{42}, Var{1}},
			pos: 2,
			out: Var{1},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := Select(tc.in, tc.pos)
			if out != tc.out {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}
}

func TestSample(t *testing.T) {
	var (
		rng       = rand.New(rand.NewSource(time.Now().UnixNano()))
		testCases = []struct {
			in     Operator
			weight func(op Operator, depth uint, rng *rand.Rand) float64
			out    Operator
		}{
			{
				in: Const{42},
				weight: func(op Operator, depth uint, rng *rand.Rand) float64 {
					return 1
				},
				out: Const{42},
			},
			{
				in: Cos{Const{42}},
				weight: func(op Operator, depth uint, rng *rand.Rand) float64 {
					switch op.(type) {
					case Const:
						return 1
					default:
						return 0
					}
				},
				out: Const{42},
			},
			{
				in: Add{Var{1}, Cos{Const{42}}},
				weight: func(op Operator, depth uint, rng *rand.Rand) float64 {
					switch depth {
					case 2:
						return 1
					default:
						return 0
					}
				},
				out: Const{42},
			},
		}
	)
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out, _ := Sample(tc.in, tc.weight, rng)
			if out != tc.out {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}
}

func TestReplaceAt(t *testing.T) {
	var testCases = []struct {
		in   Operator
		pos  uint
		with Operator
		out  Operator
	}{
		{
			in:   Add{Const{0}, Var{0}},
			pos:  0,
			with: Const{42},
			out:  Const{42},
		},
		{
			in:   Add{Const{0}, Var{0}},
			pos:  1,
			with: Const{42},
			out:  Add{Const{42}, Var{0}},
		},
		{
			in:   Add{Const{0}, Var{0}},
			pos:  2,
			with: Const{42},
			out:  Add{Const{0}, Const{42}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := ReplaceAt(tc.in, tc.pos, tc.with)
			if out != tc.out {
				t.Errorf("Expected %s, got %s", tc.out, out)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	var (
		op         = Add{Var{0}, Const{42.24}}
		bytes, err = MarshalJSON(op)
	)
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
		return
	}
	newOp, err := UnmarshalJSON(bytes)
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
		return
	}
	if newOp != op {
		t.Errorf("Expected %s, got %s", op, newOp)
		return
	}
}

func TestGetConsts(t *testing.T) {
	var testCases = []struct {
		op     Operator
		values []float64
	}{
		{
			op:     Const{42},
			values: []float64{42},
		},
		{
			op:     Add{Const{1}, Const{0}},
			values: []float64{1, 0},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			values := GetConsts(tc.op)
			if !reflect.DeepEqual(values, tc.values) {
				t.Errorf("Expected %v, got %v", tc.values, values)
			}
		})
	}
}

func TestSetConsts(t *testing.T) {
	var testCases = []struct {
		in     Operator
		values []float64
		out    Operator
	}{
		{
			in:     Const{42},
			values: []float64{43},
			out:    Const{43},
		},
		{
			in:     Add{Const{1}, Const{0}},
			values: []float64{0, 1},
			out:    Add{Const{0}, Const{1}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) {
			out := SetConsts(tc.in, tc.values)
			if !reflect.DeepEqual(out, tc.out) {
				t.Errorf("Expected %v, got %v", tc.out, out)
			}
		})
	}
}
