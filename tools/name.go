package tools

import (
	"fmt"
	"math/rand"
)

var adjectives = []string{"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry", "dark", "summer", "icy"}
var nouns = []string{"waterfall", "river", "breeze", "moon", "rain", "wind", "sea", "morning", "snow", "lake"}

func GenerateDefaultName() string {
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	return fmt.Sprintf("%s_%s", adj, noun)
}
