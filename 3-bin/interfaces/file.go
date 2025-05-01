package interfaces

type File interface {
	Read(source string) ([]byte, string, error)

	Write(content []byte, name string) error
}
