package dataaccess

type ItemAccessor interface {
	GetRootItem() (*Item, error)
}
