package goconcat

type FileType string
type PrefixType string
type Directory string
type Destination string

const (
	// file types
	FileGo         FileType = ".go"
	fileGit        FileType = ".git"
	fileJavaScript FileType = ".js"
	fileTypeScript FileType = ".ts"

	// prefix type
	prefixMockery  PrefixType = "mock_"
	prefixGoconcat PrefixType = "mocks_"

	// directories
	girectoryGit Directory = ".git"

	// destination directories
	destinationMockery Destination = "mockery"

	mainFile      string = "main.go"
	rootDirectory string = "."
	goconcat      string = "goconcat"
)
