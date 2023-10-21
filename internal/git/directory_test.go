package git

import "testing"

func TestGetStagedStatus(t *testing.T) {
	testCases := []struct {
		name           string
		directory      Directory
		expectedResult StagedStatus
	}{
		{
			name: "Child directory and files fully staged",
			directory: func() Directory {
				directory := Directory{}

				file := File{
					IndexStatus:    'M',
					WorktreeStatus: '.',
				}
				directory.Files = append(directory.Files, file)

				subDirectory := Directory{
					Files: []File{},
				}
				subDirectory.Files = append(subDirectory.Files, file)
				directory.Directories = append(directory.Directories, &subDirectory)

				return directory
			}(),
			expectedResult: FullyStaged,
		},
		{
			name: "Child directory partially staged",
			directory: func() Directory {
				directory := Directory{}

				file := File{
					IndexStatus:    'M',
					WorktreeStatus: '.',
				}
				directory.Files = append(directory.Files, file)

				subDirectory := Directory{
					Files: []File{},
				}
				directory.Directories = append(directory.Directories, &subDirectory)

				subDirectory.Files = append(subDirectory.Files, file)

				file.WorktreeStatus = 'M'
				subDirectory.Files = append(subDirectory.Files, file)

				return directory
			}(),
			expectedResult: PartiallyStaged,
		},
		{
			name: "One staged child directory another unstaged",
			directory: func() Directory {
				directory := Directory{}

				subDirectory := Directory{
					Files: []File{
						{
							IndexStatus:    'M',
							WorktreeStatus: 'M',
						},
					},
				}
				directory.Directories = append(directory.Directories, &subDirectory)

				subDirectoryStaged := Directory{
					Files: []File{
						{
							IndexStatus:    'M',
							WorktreeStatus: '.',
						},
					},
				}
				directory.Directories = append(directory.Directories, &subDirectoryStaged)

				return directory
			}(),
			expectedResult: PartiallyStaged,
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc

		result := tc.directory.GetStagedStatus()
		if result != tc.expectedResult {
			t.Fatalf("Wrong stage status found for directory. Test: '%s'. Expected: '%d' Got: '%d'", tc.name, tc.expectedResult, result)
		}
	}
}
