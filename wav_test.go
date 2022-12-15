package wav

import (
	"fmt"
	"os"
	"testing"
)

func TestWav(t *testing.T) {

	f, err := os.Create("2.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	ww, err := NewWrite(f, 10, 1, 16000, 16, 1, map[string][]byte{
		"aa": []byte("aaaaaaaaaaaa"),
		"bb": []byte("bbbbb"),
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("format : %+v , dataLen : %d \n", ww.Format(), ww.GetDataLen())

	file, err := os.Open("2.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	w, err := NewRead(file)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("format : %+v , dataLen : %d \n", w.Format(), w.GetDataLen())

}
