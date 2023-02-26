package barcode

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/domgoodwin/bookscan/lookup"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

const csvFilePath = "./books.csv"

var disableBeep = false

func WaitForScan() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("bookscan v0.1")
	fmt.Println("---------------------")

	var streamer beep.StreamSeekCloser
	if os.Getenv("ENABLE_BEEP") == "true" {
		disableBeep = true
		setupBeep(streamer)
	}

	for true {
		fmt.Print("scan book-> ")
		playBeep(streamer)
		text, _ := reader.ReadString('\n')
		found, isbn := trimAndValidate(text)
		fmt.Println(isbn)
		if !found {
			fmt.Println("invalid isbn, ignoring.")
		}

		book, err := lookup.LookupISBN(isbn)
		if err != nil {
			fmt.Println("lookup err: %", err)
			continue
		}
		fmt.Println(book.Info())
		err = book.StoreInCSV(csvFilePath)
		if err != nil {
			fmt.Println("csv err: %", err)
			continue
		}
	}
}

func trimAndValidate(text string) (bool, string) {
	out := strings.TrimRight(text, "\n")
	out = strings.TrimRight(out, "\r")
	out = strings.Replace(out, "-", "", -1)
	out = strings.Replace(out, " ", "", -1)
	out = strings.Replace(out, "isbn", "", -1)

	re13 := regexp.MustCompile(`[0-9]{13}`)
	re10 := regexp.MustCompile(`[0-9]{10}`)

	if re13.MatchString(out) {
		return true, out
	}
	if re10.MatchString(out) {
		return true, out
	}
	return false, out
}

func setupBeep(streamer beep.StreamSeekCloser) {
	f, err := os.Open("./assets/beep.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}
func playBeep(streamer beep.StreamSeekCloser) {
	if disableBeep {
		return
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}
