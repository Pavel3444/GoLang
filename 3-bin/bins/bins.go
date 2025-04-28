package bins

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Bin struct {
	Id         string    `json:"id"`
	Private    bool      `json:"private"`
	CreateTime time.Time `json:"create_time"`
	Name       string    `json:"name"`
}

type BinList struct {
	Bins         []Bin     `json:"bins"`
	DateOfUpdate time.Time `json:"date_of_update"`
}

func NewBinList() *BinList {
	return &BinList{Bins: make([]Bin, 0)}
}

func (bins *BinList) AddBin(private bool, name string) *Bin {
	bin := Bin{
		Private:    private,
		Name:       name,
		Id:         uuid.New().String(),
		CreateTime: time.Now(),
	}
	bins.Bins = append(bins.Bins, bin)
	return &bin
}

func (bl BinList) ToBytes() ([]byte, error) {
	return json.Marshal(bl)
}

func (bl *BinList) FromBytes(data []byte) error {
	return json.Unmarshal(data, bl)
}
