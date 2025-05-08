package bins

type BinRepository interface {
	Add(private bool, name string) (*Bin, error)

	List() ([]Bin, error)
}
