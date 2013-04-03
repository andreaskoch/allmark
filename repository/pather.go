package repository

type Pather interface {
	DirectoryAbsolute() string
	PathAbsolute() string
}

type RenderPather interface {
	RenderPathAbsolute() string
	RenderPathRelative() string
}
