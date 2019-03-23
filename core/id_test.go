package core

import (
	"errors"
	"strconv"
	"testing"
	"time"
)

func TestIdGen_Generate(t *testing.T) {
	id := NewIdGen(1)

	m := make(map[string]struct{})
	for i := 0; i < 10000; i++ {
		g, err := id.Generate()
		if err != nil {
			t.Fatal(err)
		}
		_, ok := m[g]
		if ok {
			t.Fatal(errors.New("has duplicated key: " + g + " in round: " + strconv.Itoa(i)))
		}
		m[g] = struct{}{}
		time.Sleep(1 * time.Millisecond)
	}
}
