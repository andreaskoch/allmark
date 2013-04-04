package path

const (
	PatherTypeItem = "item"
	PatherTypeFile = "file"
)

type Pather interface {
	Path() string
	PathType() string
}
