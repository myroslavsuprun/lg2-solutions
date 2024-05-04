package main

import (
	"slices"
	"testing"
)

func TestRanking(t *testing.T) {
	l := League{
		Teams: []Team{{
			Name:    "First",
			Players: []string{},
		}},
		Wins: map[string]int{},
	}

	l.MatchResult("First", 1, "Second", 2)

	expected := []string{"Second", "First"}
	got := l.Ranking()

	if !slices.Equal(expected, got) {
		t.Errorf("Ranking order is incorrect. \n got: %v \n expected: %v", got, expected)

	}
}
