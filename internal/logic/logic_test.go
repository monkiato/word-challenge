package logic

import (
	"testing"
)

func TestCurrentWord(t *testing.T) {
	logicWords := NewWords()
	word := logicWords.CurrentWord
	if word == "" {
		t.Fatalf("unexpected empty word")
	}
}

func TestEvaluateWord(t *testing.T) {
	logicWords := NewWords()
	word := logicWords.CurrentWord
	success, score := logicWords.EvaluateSuccess(word)
	if !success {
		t.Fatalf("unexpeted failure evaluation")
	}
	// check if out of range
	if score <= 0 || score > 500 {
		t.Fatalf("unexpeted score evaluated from logic")
	}
	if logicWords.CurrentWord == word {
		t.Fatalf("same word returned after success matching")
	}
}

