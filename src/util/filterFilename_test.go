package util

import (
	"strings"
	"testing"
)

func TestFilterFilename(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"file.txt", "file.txt", false},
		{"file/with/slashes.txt", "filewithslashes.txt", false},
		{"file\\with\\backslashes.txt", "filewithbackslashes.txt", false},
		{"file/with/more/than/one.period.txt", "file/with/more/than/one.period.txt", true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := FilterFilename(tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("FilterFilename(%s) error = %v, wantErr %v", tc.input, err, tc.wantErr)
				return
			}

			if result != tc.expected {
				t.Errorf("FilterFilename(%s) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"file.txt", "file.txt", false},
		{"file with spaces.txt", "file_with_spaces.txt", false},
		{"file@with#special$chars.txt", "filewithspecialchars.txt", false},
		{"file/with/slashes.txt", "filewithslashes.txt", false},
		{"file\\with\\backslashes.txt", "filewithbackslashes.txt", false},
		{"CON.txt", "", true},                         // Reserved name
		{"AUX.txt", "", true},                         // Reserved name
		{"COM1.txt", "", true},                        // Reserved name
		{"LPT1.txt", "", true},                        // Reserved name
		{"", "", true},                                // Empty filename
		{"a", "a", false},                             // Single character
		{strings.Repeat("a", 256) + ".txt", "", true}, // Too long
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := SanitizeFilename(tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("SanitizeFilename(%s) error = %v, wantErr %v", tc.input, err, tc.wantErr)
				return
			}

			if result != tc.expected {
				t.Errorf("SanitizeFilename(%s) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestValidateFilename(t *testing.T) {
	testCases := []struct {
		input   string
		wantErr bool
	}{
		{"file.txt", false},
		{"file_with_spaces.txt", false},
		{"file-with-dashes.txt", false},
		{"file_with_underscores.txt", false},
		{"file.with.dots.txt", false},
		{"file<>with<>invalid.txt", true},         // Invalid characters
		{"file\"with\"quotes.txt", true},          // Invalid characters
		{"file|with|pipes.txt", true},             // Invalid characters
		{"file?with?question.txt", true},          // Invalid characters
		{"file*with*asterisks.txt", true},         // Invalid characters
		{"CON.txt", true},                         // Reserved name
		{"AUX.txt", true},                         // Reserved name
		{"COM1.txt", true},                        // Reserved name
		{"LPT1.txt", true},                        // Reserved name
		{"", true},                                // Empty filename
		{strings.Repeat("a", 256) + ".txt", true}, // Too long
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			err := ValidateFilename(tc.input)

			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateFilename(%s) error = %v, wantErr %v", tc.input, err, tc.wantErr)
				return
			}
		})
	}
}

func TestGenerateUniqueFilename(t *testing.T) {
	testCases := []struct {
		baseFilename string
		existing     []string
		expected     string
	}{
		{"file.txt", []string{}, "file.txt"},
		{"file.txt", []string{"file.txt"}, "file_1.txt"},
		{"file.txt", []string{"file.txt", "file_1.txt"}, "file_2.txt"},
		{"file.txt", []string{"file.txt", "file_1.txt", "file_2.txt"}, "file_3.txt"},
		{"file", []string{"file"}, "file_1"},
		{"file", []string{"file", "file_1"}, "file_2"},
	}

	for _, tc := range testCases {
		t.Run(tc.baseFilename, func(t *testing.T) {
			checkExists := func(filename string) bool {
				for _, existing := range tc.existing {
					if filename == existing {
						return true
					}
				}
				return false
			}

			result := GenerateUniqueFilename(tc.baseFilename, checkExists)

			if result != tc.expected {
				t.Errorf("GenerateUniqueFilename(%s) = %v, want %v", tc.baseFilename, result, tc.expected)
			}
		})
	}
}
