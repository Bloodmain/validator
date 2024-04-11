package homework

import (
	"errors"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "invalid struct: interface",
			args: args{
				v: new(any),
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: map",
			args: args{
				v: map[string]string{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: string",
			args: args{
				v: "some string",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "valid struct with no fields",
			args: args{
				v: struct{}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with untagged fields",
			args: args{
				v: struct {
					f1 string
					f2 string
				}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with unexported fields",
			args: args{
				v: struct {
					foo string `validate:"len:10"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrValidateForUnexportedFields)
			},
		},
		{
			name: "invalid validator syntax",
			args: args{
				v: struct {
					Foo string `validate:"len:abcdef"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrInvalidValidatorSyntax)
			},
		},
		{
			name: "valid struct with tagged fields",
			args: args{
				v: struct {
					Len       string `validate:"len:20"`
					LenZ      string `validate:"len:0"`
					InInt     int    `validate:"in:20,25,30"`
					InNeg     int    `validate:"in:-20,-25,-30"`
					InStr     string `validate:"in:foo,bar"`
					MinInt    int    `validate:"min:10"`
					MinIntNeg int    `validate:"min:-10"`
					MinStr    string `validate:"min:10"`
					MinStrNeg string `validate:"min:-1"`
					MaxInt    int    `validate:"max:20"`
					MaxIntNeg int    `validate:"max:-2"`
					MaxStr    string `validate:"max:20"`
				}{
					Len:       "abcdefghjklmopqrstvu",
					LenZ:      "",
					InInt:     25,
					InNeg:     -25,
					InStr:     "bar",
					MinInt:    15,
					MinIntNeg: -9,
					MinStr:    "abcdefghjkl",
					MinStrNeg: "abc",
					MaxInt:    16,
					MaxIntNeg: -3,
					MaxStr:    "abcdefghjklmopqrst",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong length",
			args: args{
				v: struct {
					Lower    string `validate:"len:24"`
					Higher   string `validate:"len:5"`
					Zero     string `validate:"len:3"`
					BadSpec  string `validate:"len:%12"`
					Negative string `validate:"len:-6"`
				}{
					Lower:    "abcdef",
					Higher:   "abcdef",
					Zero:     "",
					BadSpec:  "abc",
					Negative: "abcd",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				expectedErrors := []struct {
					err   error
					field string
				}{
					{
						err:   ErrLenValidationFailed,
						field: "Lower",
					},
					{
						err:   ErrLenValidationFailed,
						field: "Higher",
					},
					{
						err:   ErrLenValidationFailed,
						field: "Zero",
					},
					{
						err:   ErrInvalidValidatorSyntax,
						field: "BadSpec",
					},
					{
						err:   ErrInvalidValidatorSyntax,
						field: "Negative",
					},
				}

				if _, ok := err.(interface{ Unwrap() []error }); !ok {
					assert.Fail(t, "err should be created with errors.Join(err...) function")
					return false
				}
				errs := err.(interface{ Unwrap() []error }).Unwrap()
				assert.Len(t, errs, 5)

				foundErrors := expectedErrors
				for i := range errs {
					actualErr := &ValidationError{}
					if errors.As(errs[i], &actualErr) {
						for ei := range expectedErrors {
							if errors.Is(actualErr, expectedErrors[ei].err) && actualErr.field == expectedErrors[ei].field {
								foundErrors = slices.Delete(foundErrors, ei, ei+1)
							}
						}
					}
				}

				assert.Empty(t, foundErrors, "unexpected errors found")
				return true
			},
		},
		{
			name: "wrong in",
			args: args{
				v: struct {
					InA     string `validate:"in:ab,cd"`
					InB     string `validate:"in:aa,bb,cd,ee"`
					InC     int    `validate:"in:-1,-3,5,7"`
					InD     int    `validate:"in:5-"`
					InEmpty string `validate:"in:"`
				}{
					InA:     "ef",
					InB:     "ab",
					InC:     2,
					InD:     12,
					InEmpty: "",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				expectedErrors := []struct {
					err   error
					field string
				}{
					{
						err:   ErrInValidationFailed,
						field: "InA",
					},
					{
						err:   ErrInValidationFailed,
						field: "InB",
					},
					{
						err:   ErrInValidationFailed,
						field: "InC",
					},
					{
						err:   ErrInvalidValidatorSyntax,
						field: "InD",
					},
					{
						err:   ErrInvalidValidatorSyntax,
						field: "InEmpty",
					},
				}

				if _, ok := err.(interface{ Unwrap() []error }); !ok {
					assert.Fail(t, "err should be created with errors.Join(err...) function")
					return false
				}
				errs := err.(interface{ Unwrap() []error }).Unwrap()
				assert.Len(t, errs, 5)

				foundErrors := expectedErrors
				for i := range errs {
					actualErr := &ValidationError{}
					if errors.As(errs[i], &actualErr) {
						for ei := range expectedErrors {
							if errors.Is(actualErr, expectedErrors[ei].err) && actualErr.field == expectedErrors[ei].field {
								foundErrors = slices.Delete(foundErrors, ei, ei+1)
							}
						}
					}
				}

				assert.Empty(t, foundErrors, "unexpected errors found")
				return true
			},
		},
		{
			name: "wrong min",
			args: args{
				v: struct {
					MinA string `validate:"min:12"`
					MinB int    `validate:"min:-12"`
					MinC int    `validate:"min:5-"`
					MinD int    `validate:"min:"`
					MinE string `validate:"min:"`
				}{
					MinA: "ef",
					MinB: -22,
					MinC: 12,
					MinD: 11,
					MinE: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				expectedErrors := []struct {
					err   error
					field string
				}{
					{
						err:   ErrMinValidationFailed,
						field: "MinA",
					},
					{
						err:   ErrMinValidationFailed,
						field: "MinB",
					},
					{
						err:   ErrInvalidValidatorSyntax,
						field: "MinC",
					},
					{
						err:   ErrInvalidValidatorSyntax,
						field: "MinD",
					},
					{
						err:   ErrInvalidValidatorSyntax,
						field: "MinE",
					},
				}

				if _, ok := err.(interface{ Unwrap() []error }); !ok {
					assert.Fail(t, "err should be created with errors.Join(err...) function")
					return false
				}
				errs := err.(interface{ Unwrap() []error }).Unwrap()
				assert.Len(t, errs, 5)

				foundErrors := expectedErrors
				for i := range errs {
					actualErr := &ValidationError{}
					if errors.As(errs[i], &actualErr) {
						for ei := range expectedErrors {
							if errors.Is(actualErr, expectedErrors[ei].err) && actualErr.field == expectedErrors[ei].field {
								foundErrors = slices.Delete(foundErrors, ei, ei+1)
							}
						}
					}
				}

				assert.Empty(t, foundErrors, "unexpected errors found")
				return true
			},
		},
		{
			name: "wrong max",
			args: args{
				v: struct {
					MaxA string `validate:"max:2"`
					MaxB string `validate:"max:-7"`
					MaxC int    `validate:"max:-12"`
					MaxD int    `validate:"max:5-"`
					MaxE int    `validate:"max:"`
					MaxF string `validate:"max:"`
				}{
					MaxA: "efgh",
					MaxB: "ab",
					MaxC: 22,
					MaxD: 12,
					MaxE: 11,
					MaxF: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				expectedErrors := []struct {
					err   error
					field string
				}{
					{
						err:   ErrMaxValidationFailed,
						field: "MaxA",
					},
					{
						err:   ErrMaxValidationFailed,
						field: "MaxB",
					},
					{
						err:   ErrMaxValidationFailed,
						field: "MaxC",
					},
					{
						err:   ErrInvalidValidatorSyntax,
						field: "MaxD",
					},
					{
						err:   ErrInvalidValidatorSyntax,
						field: "MaxE",
					},
					{
						err:   ErrInvalidValidatorSyntax,
						field: "MaxF",
					},
				}

				if _, ok := err.(interface{ Unwrap() []error }); !ok {
					assert.Fail(t, "err should be created with errors.Join(err...) function")
					return false
				}
				errs := err.(interface{ Unwrap() []error }).Unwrap()
				assert.Len(t, errs, 6)

				foundErrors := expectedErrors
				for i := range errs {
					actualErr := &ValidationError{}
					if errors.As(errs[i], &actualErr) {
						for ei := range expectedErrors {
							if errors.Is(actualErr, expectedErrors[ei].err) && actualErr.field == expectedErrors[ei].field {
								foundErrors = slices.Delete(foundErrors, ei, ei+1)
							}
						}
					}
				}

				assert.Empty(t, foundErrors, "unexpected errors found")
				return true
			},
		},
		{
			name: "bad types",
			args: args{
				v: struct {
					B  bool     `validate:"max:2"`
					S  struct{} `validate:"max:2"`
					It int      `validate:"len:3"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				expectedErrors := []struct {
					err   error
					field string
				}{
					{
						err:   ErrUnsupportedType,
						field: "B",
					},
					{
						err:   ErrUnsupportedType,
						field: "S",
					},
					{
						err:   ErrUnsupportedOperationForType,
						field: "It",
					},
				}

				if _, ok := err.(interface{ Unwrap() []error }); !ok {
					assert.Fail(t, "err should be created with errors.Join(err...) function")
					return false
				}
				errs := err.(interface{ Unwrap() []error }).Unwrap()
				assert.Len(t, errs, len(expectedErrors))

				foundErrors := expectedErrors
				for i := range errs {
					actualErr := &ValidationError{}
					if errors.As(errs[i], &actualErr) {
						for ei := range expectedErrors {
							if errors.Is(actualErr, expectedErrors[ei].err) && actualErr.field == expectedErrors[ei].field {
								foundErrors = slices.Delete(foundErrors, ei, ei+1)
							}
						}
					}
				}

				assert.Empty(t, foundErrors, "unexpected errors found")
				return true
			},
		},
		{
			name: "slices without Tag",
			args: args{
				v: struct {
					Strings []string
					Ints    []int
				}{},
			},
			wantErr: false,
		},
		{
			name: "valid slices",
			args: args{
				v: struct {
					Strings []string `validate:"len:3"`
					Ints    []int    `validate:"min:-2"`
					IntsIn  []int    `validate:"in:12,-1,3"`
				}{
					Strings: []string{"abc", "dda", "001"},
					Ints:    []int{0, -1, -2, 3, 2},
					IntsIn:  []int{12, -1, -1, -1, 12, -1, 3, 3},
				},
			},
			wantErr: false,
		},
		{
			name: "wrong slices",
			args: args{
				v: struct {
					Strings []string `validate:"len:3"`
					Ints    []int    `validate:"min:-2"`
					IntsIn  []int    `validate:"in:12,-1,3"`
				}{
					Strings: []string{"abc", "dda", "001", "ab"},
					Ints:    []int{0, -1, -2, 3, -3, 2},
					IntsIn:  []int{12, -1, -1, -1, 12, 1, 3, 3},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				expectedErrors := []struct {
					err   error
					field string
				}{
					{
						err:   ErrLenValidationFailed,
						field: "Strings",
					},
					{
						err:   ErrMinValidationFailed,
						field: "Ints",
					},
					{
						err:   ErrInValidationFailed,
						field: "IntsIn",
					},
				}

				if _, ok := err.(interface{ Unwrap() []error }); !ok {
					assert.Fail(t, "err should be created with errors.Join(err...) function")
					return false
				}
				errs := err.(interface{ Unwrap() []error }).Unwrap()
				assert.Len(t, errs, len(expectedErrors))

				foundErrors := expectedErrors
				for i := range errs {
					actualErr := &ValidationError{}
					if errors.As(errs[i], &actualErr) {
						for ei := range expectedErrors {
							if errors.Is(actualErr, expectedErrors[ei].err) && actualErr.field == expectedErrors[ei].field {
								foundErrors = slices.Delete(foundErrors, ei, ei+1)
							}
						}
					}
				}

				assert.Empty(t, foundErrors, "unexpected errors found")
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
