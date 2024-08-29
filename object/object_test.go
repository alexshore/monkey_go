package object

import "testing"

func TestStringHashKey(t *testing.T) {
	h1 := &String{Value: "Hello World"}
	h2 := &String{Value: "Hello World"}

	d1 := &String{Value: "AHHHHHHHHHH"}
	d2 := &String{Value: "AHHHHHHHHHH"}

	if h1.HashKey() != h2.HashKey() {
		t.Errorf("strings with same content have distinct hashes")
	}

	if d1.HashKey() != d2.HashKey() {
		t.Errorf("strings with same content have distinct hashes")
	}

	if h1.HashKey() == d1.HashKey() {
		t.Errorf("strings with distinct content have same hashes")
	}
}
