package repository

type Pather interface {
	AbsolutePath() string
	RelativePath(basePath string) string
}
