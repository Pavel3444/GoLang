package interfaces

type Storage interface {
	Save(content []byte, name string) error

	Read(name string) ([]byte, error)

	AddToIndex(name string, id string) error

	RemoveFromIndex(name string) error

	GetIndex() (map[string]string, error)
}
