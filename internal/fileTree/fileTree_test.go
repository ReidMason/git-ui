package filetree

import "testing"

func TestGetSelectedLine(t *testing.T) {
	testCases := []struct {
		expectedResult FileTreeLine
		name           string
		filetree       FileTree
		expectError    bool
	}{
		{
			name:           "Single file tree line",
			expectedResult: FileTreeLine{Depth: 1},
			expectError:    false,
			filetree: FileTree{
				fileTreeLines: []FileTreeLine{{Depth: 1}},
				cursorIndex:   0,
				IsFocused:     false,
			},
		},
		{
			name:           "Index is minus",
			expectedResult: FileTreeLine{Depth: 1},
			expectError:    false,
			filetree: FileTree{
				fileTreeLines: []FileTreeLine{{Depth: 1}},
				cursorIndex:   -1,
				IsFocused:     false,
			},
		},
		{
			name:           "Index larger than list size",
			expectedResult: FileTreeLine{Depth: 2},
			expectError:    false,
			filetree: FileTree{
				fileTreeLines: []FileTreeLine{{Depth: 1}, {Depth: 2}},
				cursorIndex:   2,
				IsFocused:     false,
			},
		},
		{
			name:           "No file tree lines in list",
			expectedResult: FileTreeLine{},
			expectError:    true,
			filetree: FileTree{
				fileTreeLines: []FileTreeLine{},
				cursorIndex:   2,
				IsFocused:     false,
			},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc
		selectedLine, err := tc.filetree.getSelectedLine()
		if tc.expectError && err == nil {
			t.Fatalf("Test: '%s' expected an error but didn't get one", tc.name)
		}

		if !tc.expectError && selectedLine.Depth != tc.expectedResult.Depth {
			t.Fatalf("Test: '%s' got the wrong file tree line. Expected: %v Got: %v", tc.name, tc.expectedResult.Depth, selectedLine.Depth)
		}
	}
}
