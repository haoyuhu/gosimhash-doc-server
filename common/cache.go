package common

import (
	"github.com/HaoyuHu/gosimhash"
	"github.com/HaoyuHu/gosimhash-doc-server/model"
)

type SimhashLimit int

const (
	LIMIT_1  = SimhashLimit(1)
	LIMIT_3  = SimhashLimit(3)
	LIMIT_7  = SimhashLimit(7)
	LIMIT_15 = SimhashLimit(15)
)

type SimhashCache interface {
	Init(docIds []string, simhashList []uint64, timeouts []int64) int

	InsertIfNotDuplicated(docId string, simhash uint64, age int64) (bool, *model.Document, error)
}

const (
	MASK_2  uint64 = 0xffffffff
	MASK_4  uint64 = 0xffff
	MASK_8  uint64 = 0xff
	MASK_16 uint64 = 0xf
)

var MASKS = map[int]uint64{2: MASK_2, 4: MASK_4, 8: MASK_8, 16: MASK_16}

type SimhashOperator struct {
	partNumber int
	mask       uint64
}

func NewSimhashOperator(number int) *SimhashOperator {
	op := &SimhashOperator{}
	op.partNumber = number
	op.mask = MASKS[number]
	return op
}

func (op *SimhashOperator) Check() bool {
	return Check2Power(op.partNumber)
}

func (op *SimhashOperator) Cut(simhash uint64) []uint64 {
	var ret []uint64 = make([]uint64, op.partNumber)
	var move = uint(gosimhash.BITS_LENGTH / op.partNumber)
	for i := 0; i < op.partNumber; i++ {
		ret[i] = op.mask & simhash
		simhash >>= move
	}
	return ret
}
