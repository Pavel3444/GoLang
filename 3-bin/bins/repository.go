package bins

import (
	"bin/interfaces"
	"fmt"
	"time"
)

type binRepo struct {
	store    interfaces.Storage
	fileName string
	list     *BinList
}

func NewBinRepository(store interfaces.Storage, fileName string) *binRepo {
	bl := NewBinList()
	if raw, err := store.Read(fileName); err == nil {
		if err := bl.FromBytes(raw); err != nil {
			fmt.Printf("warning: cannot parse %s: %v\n", fileName, err)
		}
	}
	return &binRepo{store, fileName, bl}
}

func (r *binRepo) Add(private bool, name string) (*Bin, error) {
	b := r.list.AddBin(private, name)
	r.list.DateOfUpdate = time.Now()

	data, err := r.list.ToBytes()
	if err != nil {
		return nil, fmt.Errorf("marshal bins: %w", err)
	}
	if err := r.store.Save(data, r.fileName); err != nil {
		return nil, fmt.Errorf("save bins: %w", err)
	}
	return b, nil
}

func (r *binRepo) List() ([]Bin, error) {
	return r.list.Bins, nil
}
