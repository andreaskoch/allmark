package repository

type Pather interface {
	PathAbsolute() string
	PathRelative() string
}
