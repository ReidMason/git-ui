package filetree

// func TestGetSelectedLine(t *testing.T) {
// 	testCases := []struct {
// 		expectedResult FileTreeLine
// 		name           string
// 		filetree       FileTree
// 		expectError    bool
// 	}{
// 		{
// 			name:           "Single file tree line",
// 			expectedResult: FileTreeLine{Depth: 1},
// 			expectError:    false,
// 			filetree: FileTree{
// 				fileTreeLines: []FileTreeLine{{Depth: 1}},
// 				cursorIndex:   0,
// 				IsFocused:     false,
// 			},
// 		},
// 		{
// 			name:           "Index is minus",
// 			expectedResult: FileTreeLine{Depth: 1},
// 			expectError:    false,
// 			filetree: FileTree{
// 				fileTreeLines: []FileTreeLine{{Depth: 1}},
// 				cursorIndex:   -1,
// 				IsFocused:     false,
// 			},
// 		},
// 		{
// 			name:           "Index larger than list size",
// 			expectedResult: FileTreeLine{Depth: 2},
// 			expectError:    false,
// 			filetree: FileTree{
// 				fileTreeLines: []FileTreeLine{{Depth: 1}, {Depth: 2}},
// 				cursorIndex:   2,
// 				IsFocused:     false,
// 			},
// 		},
// 		{
// 			name:           "No file tree lines in list",
// 			expectedResult: FileTreeLine{},
// 			expectError:    true,
// 			filetree: FileTree{
// 				fileTreeLines: []FileTreeLine{},
// 				cursorIndex:   2,
// 				IsFocused:     false,
// 			},
// 		},
// 	}
//
// 	t.Parallel()
// 	for _, tc := range testCases {
// 		tc := tc
// 		selectedLine, err := tc.filetree.getSelectedLine()
// 		if tc.expectError && err == nil {
// 			t.Fatalf("Test: '%s' expected an error but didn't get one", tc.name)
// 		}
//
// 		if !tc.expectError && selectedLine.Depth != tc.expectedResult.Depth {
// 			t.Fatalf("Test: '%s' got the wrong file tree line. Expected: %v Got: %v", tc.name, tc.expectedResult.Depth, selectedLine.Depth)
// 		}
// 	}
// }
//
// func TestRender(t *testing.T) {
// 	rootDir := git.Directory{
// 		Name: "Root",
// 	}
// 	rootDir.ToggleExpanded()
//
// 	directory := git.Directory{
// 		Name:   "Directory",
// 		Parent: &rootDir,
// 	}
// 	directory.ToggleExpanded()
// 	rootDir.Directories = append(rootDir.Directories, &directory)
//
// 	file := git.File{
// 		Name:           "File",
// 		Parent:         &directory,
// 		IndexStatus:    'M',
// 		WorkTreeStatus: '.',
// 	}
// 	directory.Files = append(directory.Files, file)
//
// 	testCases := []struct {
// 		name, expectedResult string
// 		filetree             FileTree
// 	}{
// 		{
// 			name:           "No file tree lines in list",
// 			expectedResult: "No changes",
// 			filetree: FileTree{
// 				fileTreeLines: []FileTreeLine{},
// 				cursorIndex:   2,
// 				IsFocused:     false,
// 			},
// 		},
// 		{
// 			name:           "Single file tree line",
// 			expectedResult: "▶ Directory                                        ",
// 			filetree: FileTree{
// 				fileTreeLines: []FileTreeLine{
// 					{
// 						Depth: 0,
// 						Item: &git.Directory{
// 							Name: "Directory",
// 						},
// 					},
// 				},
// 				cursorIndex: 0,
// 				IsFocused:   false,
// 			},
// 		},
// 		{
// 			name: "Two deep file tree line",
// 			expectedResult: `▼ Root
//   ▼ Directory
//      M. File`,
// 			filetree: FileTree{
// 				fileTreeLines: []FileTreeLine{
// 					{
// 						Depth: 0,
// 						Item:  &rootDir,
// 					},
// 					{
// 						Depth: 1,
// 						Item:  &directory,
// 					},
// 					{
// 						Depth: 2,
// 						Item:  &file,
// 					},
// 				},
// 				cursorIndex: 0,
// 				IsFocused:   false,
// 			},
// 		},
// 	}
//
// 	t.Parallel()
// 	for _, tc := range testCases {
// 		tc := tc
// 		result := tc.filetree.Render()
//
// 		rArr := []rune(result)
// 		eArr := []rune(tc.expectedResult)
// 		for i, c := range eArr {
// 			if c != rArr[i] {
// 				t.Errorf("Non match at %d. Expected: '%s' Got '%s'", i, string(c), string(rArr[i]))
// 			} else {
// 				t.Logf("Match %s", string(c))
// 			}
// 		}
//
// 		if tc.expectedResult != result {
// 			t.Fatalf("Test: '%s' wrong output.\nExpected: '%v'\n     Got: '%v'", tc.name, formatDisplayString(tc.expectedResult), formatDisplayString(result))
// 		}
// 	}
// }
//
// func formatDisplayString(input string) string {
// 	input = strings.ReplaceAll(input, " ", "X")
// 	input = strings.ReplaceAll(input, "\n", "N")
// 	return input
// }
