package tiktok_go

import (
	"testing"
	"time"
)

func TestGetter_Text(t *testing.T) {
	getter, err := NewTorGetter("127.0.0.1", 9050)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	//txt, err := getter.Text("https://www.tiktok.com/@der_verrueckte_mutmacher", "[data-e2e=\"user-bio\"]")
	for i := 0; i < 3; i++ {
		txt, _ := getter.Text("https://check.torproject.org/?lang=en_US", "strong")
		getter.NewCircuit()
		t.Log(txt)
	}
	time.Sleep(6 * time.Second)
	for i := 0; i < 3; i++ {
		txt, _ := getter.Text("https://check.torproject.org/?lang=en_US", "strong")
		getter.NewCircuit()
		t.Log(txt)
	}
}
