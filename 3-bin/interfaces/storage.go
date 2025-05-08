package interfaces

type Storage interface {
	Save(content []byte, name string) error

	Read(name string) ([]byte, error)
}
