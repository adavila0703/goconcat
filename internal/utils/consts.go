package utils

type FileType string
type PrefixType string
type Directory string
type Destination string

const (
	// file types
	FileGo         FileType = ".go"
	FilePython     FileType = ".py"
	FileGit        FileType = ".git"
	FileJavaScript FileType = ".js"
	FileTypeScript FileType = ".ts"

	// prefix type
	PrefixMockery  PrefixType = "mock_"
	PrefixGoconcat PrefixType = "mocks_"

	// directories
	DirectoryGit Directory = ".git"

	// destination directories
	DestinationMockery Destination = "mockery"
)
