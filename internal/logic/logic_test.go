package logic

import (
	"testing"
)

func TestCurrentWord(t *testing.T) {
	logicWords := NewWords()
	logicWords.Start()
	word := logicWords.CurrentWord
	if word == "" {
		t.Fatalf("unexpected empty word")
	}
}

func TestEvaluateWord(t *testing.T) {
	logicWords := NewWords()
	logicWords.Start()
	word := logicWords.CurrentWord
	success, score := logicWords.EvaluateSuccess(word)
	if !success {
		t.Fatalf("unexpeted failure evaluation")
	}
	// check if out of range
	if score <= 0 || score > 1000 {
		t.Fatalf("unexpected score evaluated from logic, got " + string(score))
	}
	if logicWords.CurrentWord == word {
		t.Fatalf("same word returned after success matching")
	}
}

