package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

type Team struct {
	Name    string
	Players []string
}

type League struct {
	Teams []Team
	Wins  map[string]int
}

func (l *League) MatchResult(fTeam string, ftScore int, sTeam string, stScore int) {
	l.Wins[fTeam] += ftScore
	l.Wins[sTeam] += stScore
}

func (l League) Ranking() []string {
	teams := make([]string, 0, len(l.Wins))

	for k := range l.Wins {
		teams = append(teams, k)
	}

	sort.Slice(teams, func(i, j int) bool {
		return l.Wins[teams[i]] > l.Wins[teams[j]]
	})

	return teams
}

type Ranker interface {
	Ranking() []string
}

func main() {
	l := League{
		Teams: []Team{
			{
				Name:    "First",
				Players: []string{},
			},
			{
				Name:    "Second",
				Players: []string{},
			},
		},
		Wins: map[string]int{},
	}

	fmt.Printf("League: %v \n", l)

	l.MatchResult("First", 1, "Second", 2)

	RankPrinter(l, os.Stdout)
}

func RankPrinter(r Ranker, w io.Writer) {
	names := r.Ranking()
	for _, n := range names {
		io.WriteString(w, n+"\n")
	}
}
