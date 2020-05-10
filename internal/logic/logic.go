package logic

import (
	"github.com/tjarratt/babble"
	"math"
	"time"
)

type Words struct {
	CurrentWord     string
	lastRenewedTime time.Time
	babbler         babble.Babbler
}

func NewWords() *Words {
	babbler := babble.NewBabbler()
	babbler.Count = 1
	words := &Words{
		babbler: babbler,
	}
	return words
}

// Start generate random word and start timer
func (w *Words) Start() {
	w.generateRandomWord()
}

func (w *Words) generateRandomWord() {
	w.CurrentWord = w.babbler.Babble()
	w.lastRenewedTime = time.Now()
}

func (w *Words) EvaluateSuccess(word string) (bool, int64) {
	if word != w.CurrentWord {
		return false, 0
	}

	elapsedTime := time.Now().Sub(w.lastRenewedTime).Milliseconds()
	// renew current word
	w.generateRandomWord()

	// if time is >= 10000ms the score is 0, then as lower time higher will be the score
	// score range will be between 0 and 1000
	// e.g. 500ms gives a score of 150 (10000 - 500 / 10 = 950)
	//		3500ms gives a score of 90 (10000 - 3500 / 10 = 650)
	return true, int64(math.Max(0, float64((10000 - elapsedTime)/10)))
}
