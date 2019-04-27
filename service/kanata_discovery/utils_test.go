package kanata_discovery

import "testing"

func TestCheckPathNode(t *testing.T) {
	if err := checkPathNode("/"); err == nil {
		t.Fatal("char '/' should be invalid")
	}

	if err := checkPathNode("\\"); err == nil {
		t.Fatal("char '\\' should be invalid")
	}

	if err := checkPathNode(""); err == nil {
		t.Fatal("empty should be invalid")
	}

	if err := checkPathNode("0123456789"); err != nil {
		t.Fatal("char '1' should be accepted")
	}

}
