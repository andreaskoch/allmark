package repository

type Renderer interface {
	Execute()
}

type DocumentRenderer struct {
	Execute func()
}

func GetRenderer(repositoryItem *Item) Renderer {

	return DocumentRenderer{
		Execute: func() {

		},
	}

}
