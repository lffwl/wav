package wav

import (
	"fmt"
	"os"
	"testing"
)

func TestNewRead(t *testing.T) {

	file, err := os.Open("1.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	w, err := NewRead(file)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("format : %+v", w.Format())

}
