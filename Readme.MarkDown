# wav
reader and write wav

# Simple Demo
    
    package main

    import (
        "fmt"
        "github.com/lffwl/wav"
        "log"
        "os"
    )

    func main() {
    // write wav
    f, err := os.Create("1.wav")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    ww, err := wav.NewWrite(f, 10, 1, 16000, 16, 1, map[string][]byte{
        "aa": []byte("aaaaaaaaaaaa"),
        "bb": []byte("bbbbb"),
    })
    if err != nil {
        log.Fatal(err)
    }

	fmt.Printf("format : %+v , dataLen : %d \n", ww.Format(), ww.GetDataLen())
	fmt.Printf("other : %+v ,\n", ww.GetOther())

	// read wav
	file, err := os.Open("1.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w, err := wav.NewRead(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("format : %+v , dataLen : %d \n", w.Format(), w.GetDataLen())
	fmt.Printf("other : %+v ,\n", ww.GetOther())
}
