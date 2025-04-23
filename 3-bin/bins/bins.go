package bins

import (
	"github.com/google/uuid"
	"time"
)

type Bin struct {
	Id         string
	Private    bool
	CreateTime time.Time
	Name       string
}

type BinList struct {
	Bins []Bin
}

func (list *BinList) NewBin(private bool, name string) *Bin {
	bin := Bin{
		Private:    private,
		Name:       name,
		Id:         uuid.New().String(),
		CreateTime: time.Now(),
	}
	list.Bins = append(list.Bins, bin)
	return &list.Bins[len(list.Bins)-1]
}

func NewBinList() *BinList {
	return &BinList{}
}
