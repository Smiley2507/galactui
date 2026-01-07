package main

import (
	"os"
	"sort"
	"strconv"
	"strings"
)

func getHighScoreFilename(mode int) string {
	if mode == 1 {
		return "highscore_boss.txt"
	}
	return "highscore.txt"
}

func loadHighScores(mode int) []int {
	filename := getHighScoreFilename(mode)
	data, err := os.ReadFile(filename)
	var scores []int
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			s, _ := strconv.Atoi(strings.TrimSpace(line))
			scores = append(scores, s)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(scores)))
	if len(scores) > 5 {
		scores = scores[:5]
	}
	for len(scores) < 5 {
		scores = append(scores, 0)
	}
	return scores
}

func saveScore(score, mode int) {
	scores := loadHighScores(mode)
	scores = append(scores, score)
	sort.Sort(sort.Reverse(sort.IntSlice(scores)))
	if len(scores) > 5 {
		scores = scores[:5]
	}
	var lines []string
	for _, s := range scores {
		lines = append(lines, strconv.Itoa(s))
	}
	filename := getHighScoreFilename(mode)
	os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
}