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

func TestSort(t *testing.T) {
	testCases := []struct {
		name          string
		directory     Directory
		expectedDirs  []string
		expectedFiles []string
	}{
		{
			name: "Sort directories",
			directory: func() Directory {
				directory := Directory{}

				bSubDirectory := Directory{
					Name:  "bSubDirectory",
					Files: []File{},
				}
				directory.Directories = append(directory.Directories, &bSubDirectory)

				aSubDirectory := Directory{
					Name:  "aSubDirectory",
					Files: []File{},
				}
				directory.Directories = append(directory.Directories, &aSubDirectory)

				return directory
			}(),
			expectedDirs:  []string{"aSubDirectory", "bSubDirectory"},
			expectedFiles: []string{},
		},
		{
			name: "Sort files",
			directory: func() Directory {
				directory := Directory{
					Files: []File{
						{name: "bFile"},
						{name: "aFile"},
					},
				}

				return directory
			}(),
			expectedDirs:  []string{},
			expectedFiles: []string{"aFile", "bFile"},
		},
	}

	t.Parallel()
	for _, tc := range testCases {
		tc := tc

		tc.directory.Sort()

		for i, subDir := range tc.expectedDirs {
			if subDir != tc.directory.Directories[i].Name {
				t.Fatalf("Got wrong directory sort index %d. Test %s expected '%s' got '%s'", i, tc.name, subDir, tc.directory.Directories[i].Name)
			}
		}

		for i, file := range tc.expectedFiles {
			if file != tc.directory.Files[i].name {
				t.Fatalf("Got wrong file sort index %d. Test %s expected '%s' got '%s'", i, tc.name, file, tc.directory.Files[i].name)
			}
		}
	}
}
