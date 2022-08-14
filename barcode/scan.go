package barcode

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/domgoodwin/bookscan/lookup"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

const csvFilePath = "./books.csv"

func WaitForScan() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("bookscan v0.1")
	fmt.Println("---------------------")

	for true {
		fmt.Print("scan book-> ")
		beep()
		text, _ := reader.ReadString('\n')
		book, err := lookup.LookupISBN(text)
		if err != nil {
			fmt.Println("lookup err: %", err)
		}
		fmt.Println(book.Info())
		err = book.StoreInCSV(csvFilePath)
		if err != nil {
			fmt.Println("csv err: %", err)
		}
	}
}

func beep() {
	f, err := os.Open("Rockafeller Skank.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	speaker.Play(streamer)
}
