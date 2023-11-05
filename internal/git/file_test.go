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

func TestGetFilePaths(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		file     File
	}{
		{
			name: "Single file name",
			file: File{
				name:           "file.md",
				secondName:     "oldfile.md",
				Parent:         nil,
				Dirpath:        "",
				IndexStatus:    '.',
				WorktreeStatus: '.',
			},
			expected: "file.md",
		},
		{
			name: "Dirpath name display",
			file: File{
				name:           "file.md",
				secondName:     "oldfile.md",
				Parent:         nil,
				Dirpath:        "directory/files",
				IndexStatus:    '.',
				WorktreeStatus: '.',
			},
			expected: "directory/files/file.md",
		},
		{
			name: "Dirpath name display trailing slash",
			file: File{
				name:           "file.md",
				secondName:     "oldfile.md",
				Parent:         nil,
				Dirpath:        "directory/files/",
				IndexStatus:    '.',
				WorktreeStatus: '.',
			},
			expected: "directory/files/file.md",
		},
		{
			name: "Dirpath name display leading slash",
			file: File{
				name:           "file.md",
				secondName:     "oldfile.md",
				Parent:         nil,
				Dirpath:        "/directory/files/",
				IndexStatus:    '.',
				WorktreeStatus: '.',
			},
			expected: "/directory/files/file.md",
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		output := tc.file.GetFilePath()

		if output != tc.expected {
			t.Fatalf("Failed to get filepath. Test %s expected '%s' got '%s'", tc.name, tc.expected, output)
		}
	}
}
