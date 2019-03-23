package core

import (
	"errors"
	"sync"
	"time"
)

type IdGen struct {
	nodeId   uint32
	nodeMask uint64

	prevTimeOffset int64
	subSeq         uint32

	timeOffset int64
	startTime  int64

	ticker *time.Ticker
	lock   sync.Mutex
}

func NewIdGen(nodeId uint32) *IdGen {
	idgen := &IdGen{
		nodeId:   nodeId,
		nodeMask: 0x00000000ffffffff & uint64(nodeId),
		ticker:   time.NewTicker(1 * time.Second),
	}
	t, _ := time.Parse(time.RFC3339, "2019-03-23T12:00:00+08:00")
	idgen.startTime = t.Unix()
	idgen.prevTimeOffset = 0
	idgen.timeOffset = 0

	go idgen.task()

	return idgen
}

func (this *IdGen) task() {
	ticker := this.ticker
	c := ticker.C
	startTime := this.startTime
	for {
		t := <-c
		this.timeOffset = t.Unix() - startTime
	}
}

func (this *IdGen) Generate() (string, error) {
	this.lock.Lock()
	prevOffset := this.prevTimeOffset
	curOffset := this.timeOffset
	subSeq := this.subSeq
	if prevOffset > curOffset {
		this.lock.Unlock()
		return "", errors.New("time back track")
	} else if prevOffset < curOffset {
		prevOffset = curOffset
		this.prevTimeOffset = curOffset
		subSeq = 0
		this.subSeq = 0
	}
	this.subSeq++
	this.lock.Unlock()

	timeAndNode := uint64(prevOffset&0x00000000ffffffff)<<32 | this.nodeMask
	subSeqForm := 0x00000000ffffffff & uint64(subSeq)

	return convertUint64ToString(timeAndNode) + convertUint64ToString(subSeqForm), nil
}

func convertUint64ToString(data uint64) string {
	b := make([]byte, 16)
	for i := 15; i >= 0; i-- {
		b[i] = mapping[int(data&0x000000000000000f)]
		data = data >> 4
	}
	return string(b)
}

var (
	mapping = []byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F',
	}
)
