package repository

type Addresser interface {
	GetAbsolutePath() string
	GetRelativePath(basePath string) string
}
