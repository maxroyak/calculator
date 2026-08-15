package calc

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 5, 3, 8},
		{"negative numbers", -5, -3, -8},
		{"mixed signs", -5, 3, -2},
		{"zero", 0, 0, 0},
		{"fractional", 1.5, 2.5, 4},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Add(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("Add(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 10, 4, 6},
		{"negative result", 4, 10, -6},
		{"both negative", -10, -4, -6},
		{"zero", 0, 0, 0},
		{"fractional", 5.5, 2.5, 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Subtract(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want float64
	}{
		{"positive numbers", 6, 7, 42},
		{"negative numbers", -6, -7, 42},
		{"mixed signs", -6, 7, -42},
		{"by zero", 5, 0, 0},
		{"fractional", 2.5, 4, 10},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Multiply(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr bool
		errMsg  string
	}{
		{"positive numbers", 10, 2, 5, false, ""},
		{"negative numbers", -10, -2, 5, false, ""},
		{"mixed signs", -10, 2, -5, false, ""},
		{"fractional", 7.5, 2.5, 3, false, ""},
		{"dividend zero", 0, 5, 0, false, ""},
		{"division by zero", 10, 0, 0, true, "division by zero is not allowed"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Divide(tc.a, tc.b)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Divide(%v, %v) expected error, got nil", tc.a, tc.b)
				}
				if err.Error() != tc.errMsg {
					t.Errorf("Divide(%v, %v) error = %q, want %q", tc.a, tc.b, err.Error(), tc.errMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("Divide(%v, %v) unexpected error: %v", tc.a, tc.b, err)
			}
			if got != tc.want {
				t.Errorf("Divide(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}
