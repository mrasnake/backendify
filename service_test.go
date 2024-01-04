package main

import (
	"os"
	"testing"
)

// errString is used for comparison of error strings
// by returning the string "nil" when given a nil value
func errString(err error) string {
	if err == nil {
		return "nil"
	}
	return err.Error()
}

// Test_isActive builds a list of custom structs and loops through each of them
// performing the associated unit test on isActive() with the specified parameters.
func Test_isActive(t *testing.T) {

	t.Parallel()

	// tests contains all the parameters, checks and expected results of each test.
	tests := []struct {
		testName string
		input    string
		want     bool
		wantErr  string
	}{
		// isActive Test #1 checks when passing blank date
		// to isActive, the appropriate bool is returned.
		{
			testName: "IsActive Test #1 - Blank",
			input:    "",
			want:     true,
			wantErr:  "nil",
		},
		// isActive Test #2 checks when passing invalid date
		// to isActive, the appropriate error is returned.
		{
			testName: "IsActive Test #2 - Invalid date string",
			input:    "not a date",
			want:     false,
			wantErr:  `Error while parsing the date time: parsing time "not a date" as "2006-01-02T15:04:05Z07:00": cannot parse "not a date" as "2006"`,
		},
		// isActive Test #3 checks when passing a valid date after
		// time.Now() to isActive, the appropriate values are returned.
		{
			testName: "IsActive Test #3 - Closed date has not passed",
			input:    "2024-01-02T15:04:05Z",
			want:     true,
			wantErr:  "nil",
		},
		// isActive Test #4 checks when passing a valid date before
		// time.Now() to isActive, the appropriate values are returned.
		{
			testName: "IsActive Test #4 - Closed date has passed",
			input:    "2021-01-02T15:04:05Z",
			want:     false,
			wantErr:  "nil",
		},
	}

	// This loops through each item in the tests list, uses the individual parameters
	// to prepare and perform the unit test and compares the results.
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			got, err := isActive(tt.input)

			if got != tt.want {
				t.Errorf("Unexpected results, got= %+#v, want= %+#v", got, tt.want)
			}
			if errString(err) != tt.wantErr {
				t.Errorf("Unexpected error returned, got= %+#v, want= %+#v", errString(err), tt.wantErr)
			}
		})
	}
}

// Test_readArgs builds a list of custom structs and loops through each of them
// performing the associated unit test on readArgs() with the specified parameters.
func Test_readArgs(t *testing.T) {

	t.Parallel()

	// tests contains all the parameters, checks and expected results of each test.
	tests := []struct {
		testName string
		args     []string
		want     map[string]string
		wantErr  string
	}{
		// readArgs Test #1 checks when passing no additional
		// arguments to readArgs, the appropriate error is returned.
		{
			testName: "readArgs Test #1 - No Arguments",
			args:     []string{"./program"},
			want:     map[string]string{},
			wantErr:  "no valid parameters",
		},
		// readArgs Test #2 checks when passing an invalid parameter
		// format to readArgs, the appropriate error is returned.
		{
			testName: "readArgs Test #2 - Invalid Argument Format",
			args:     []string{"./program", "invalid-param"},
			want:     map[string]string{},
			wantErr:  "invalid parameter format",
		},
		// readArgs Test #3 checks when passing an invalid
		// URL to readArgs, the appropriate error is returned.
		{
			testName: "readArgs Test #3 - Invalid URL",
			args:     []string{"./program", "ru=invalidURL"},
			want:     map[string]string{},
			wantErr:  "invalid parameter, must contain valid URL",
		},
		// readArgs Test #4 checks when passing a single valid
		// argument to readArgs, the appropriate value is returned.
		{
			testName: "readArgs Test #4 - Single Valid Parameter",
			args:     []string{"./program", "us=http://localhost:9002"},
			want: map[string]string{
				"us": "http://localhost:9002",
			},
			wantErr: "nil",
		},
		// readArgs Test #5 checks when passing a single valid
		// argument to readArgs, the appropriate value is returned.
		{
			testName: "readArgs Test #5 - Multi Valid Parameters",
			args:     []string{"./program", "us=http://localhost:9002", "uk=http://localhost:9001"},
			want: map[string]string{
				"us": "http://localhost:9002",
				"uk": "http://localhost:9001",
			},
			wantErr: "nil",
		},
	}

	// This loops through each item in the tests list, uses the individual parameters
	// to prepare and perform the unit test and compares the results.
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			os.Args = tt.args

			got, err := readArgs()
			if errString(err) != tt.wantErr {
				t.Errorf("Unexpected error returned, got= %+#v, want= %+#v", errString(err), tt.wantErr)
			}

			if len(got) != len(tt.want) {
				t.Errorf("Unexpected length of results, got= %+#v, want= %+#v", len(got), len(tt.want))
			} else {
				for k, v := range got {
					if v != tt.want[k] {
						t.Errorf("Unexpected map value for key: %+#v, got=%+#v, want=%+#v", k, v, tt.want[k])
					}
				}
			}
		})
	}
}

// Test_formRequest builds a list of custom structs and loops through each of them
// performing the associated unit test on formRequest() with the specified parameters.
func Test_formRequest(t *testing.T) {

	t.Parallel()

	// tests contains all the parameters, checks and expected results of each test.
	tests := []struct {
		testName string
		input    *GetCompanyRequest
		args     []string
		want     string
	}{
		// formRequest Test #1 checks when passing valid argument
		// to formRequest, the appropriate string is returned.
		{
			testName: "formRequest Test #1 - Success",
			input: &GetCompanyRequest{
				ID:   "1234",
				Code: "us",
			},
			args: []string{"./program", "us=http://localhost:9002"},
			want: "http://localhost:9002/companies/1234",
		},
	}

	// This loops through each item in the tests list, uses the individual parameters
	// to prepare and perform the unit test and compares the results.
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			os.Args = tt.args
			s, err := NewService()
			if err != nil {
				t.Error("error creating service")
			}

			got := s.formRequest(tt.input)

			if got != tt.want {
				t.Errorf("Unexpected results, got= %+#v, want= %+#v", got, tt.want)
			}
		})
	}
}