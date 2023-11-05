package git

import "testing"

func TestDisplayName(t *testing.T) {
	testCases := []struct {
		name     string
		file     File
		expected string
	}{
		{
			name: "Basic display name",
			file: File{
				name:           "file.md",
				secondName:     "",
				Parent:         nil,
				Dirpath:        "",
				IndexStatus:    '.',
				WorktreeStatus: '.',
			},
			expected: "file.md",
		},
		{
			name: "Second name display",
			file: File{
				name:           "file.md",
				secondName:     "oldfile.md",
				Parent:         nil,
				Dirpath:        "",
				IndexStatus:    '.',
				WorktreeStatus: '.',
			},
			expected: "oldfile.md -> file.md",
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		output := tc.file.getDisplayName()

		if output != tc.expected {
			t.Fatalf("Failed to get display name. Test %s expected '%s' got '%s'", tc.name, tc.expected, output)
		}
	}
}
