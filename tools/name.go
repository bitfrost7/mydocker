package tools

import (
	"fmt"
	"math/rand"
	"time"
)

var adjectives = []string{"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry", "dark", "summer", "icy"}
var nouns = []string{"waterfall", "river", "breeze", "moon", "rain", "wind", "sea", "morning", "snow", "lake"}

func GenerateDefaultName() string {
	rand.Seed(time.Now().UnixNano())
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	return fmt.Sprintf("%s_%s", adj, noun)
}
