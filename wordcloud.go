// Package wordcloud generates a simple SVG word cloud from text or weighted words.
package wordcloud

import (
	"errors"
	"fmt"
	"html"
	"math"
	"sort"
	"strings"
	"unicode"
)

// Word is a label and its positive relative weight.
type Word struct {
	Text   string
	Weight int
}

// Options controls the SVG canvas and maximum number of words.
type Options struct {
	Width    int
	Height   int
	MaxWords int
}

type box struct {
	x, y, width, height float64
}

var palette = [...]string{"#275dad", "#29a69a", "#ef7d57", "#6b5b95", "#c15c86"}

// FromText counts Unicode letters and numbers, then draws the most common words.
func FromText(text string, options Options) (string, error) {
	counts := make(map[string]int)
	var token strings.Builder
	flush := func() {
		if token.Len() > 0 {
			counts[token.String()]++
			token.Reset()
		}
	}
	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			token.WriteRune(r)
		} else {
			flush()
		}
	}
	flush()
	words := make([]Word, 0, len(counts))
	for label, weight := range counts {
		words = append(words, Word{Text: label, Weight: weight})
	}
	return FromWords(words, options)
}

// FromWords draws weighted labels as an SVG. Labels are escaped as XML text.
func FromWords(words []Word, options Options) (string, error) {
	width, height, maxWords := options.Width, options.Height, options.MaxWords
	if width == 0 {
		width = 800
	}
	if height == 0 {
		height = 480
	}
	if maxWords == 0 {
		maxWords = 60
	}
	if width < 100 || height < 100 || maxWords < 1 {
		return "", errors.New("width and height must be at least 100; max words must be positive")
	}

	entries := make([]Word, 0, len(words))
	for _, word := range words {
		if strings.TrimSpace(word.Text) == "" || word.Weight <= 0 {
			continue
		}
		entries = append(entries, word)
	}
	if len(entries) == 0 {
		return "", errors.New("provide at least one non-empty word with a positive weight")
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Weight != entries[j].Weight {
			return entries[i].Weight > entries[j].Weight
		}
		return entries[i].Text < entries[j].Text
	})
	if len(entries) > maxWords {
		entries = entries[:maxWords]
	}

	var svg strings.Builder
	fmt.Fprintf(&svg, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%d\" height=\"%d\" viewBox=\"0 0 %d %d\">", width, height, width, height)
	svg.WriteString("<rect width=\"100%\" height=\"100%\" fill=\"white\"/><g font-family=\"Arial, Helvetica, sans-serif\" font-weight=\"bold\">")
	placed := make([]box, 0, len(entries))
	for index, word := range entries {
		size := float64(16 + int(math.Round(54*math.Sqrt(float64(word.Weight)/float64(entries[0].Weight)))))
		labelWidth := 8.0
		for _, r := range word.Text {
			if r > 0x2e80 {
				labelWidth += size
			} else {
				labelWidth += size * 0.6
			}
		}
		labelHeight := size*1.25 + 8
		if labelWidth > float64(width) || labelHeight > float64(height) {
			continue
		}
		for step := 0; step < 1800; step++ {
			radius := 2.6 * math.Sqrt(float64(step))
			angle := float64(step)*2.399963229728653 + float64(index)
			x := float64(width)/2 + radius*math.Cos(angle) - labelWidth/2
			y := float64(height)/2 + radius*math.Sin(angle) - labelHeight/2
			candidate := box{x, y, labelWidth, labelHeight}
			if x < 0 || y < 0 || x+labelWidth > float64(width) || y+labelHeight > float64(height) {
				continue
			}
			collision := false
			for _, previous := range placed {
				if x < previous.x+previous.width && x+labelWidth > previous.x && y < previous.y+previous.height && y+labelHeight > previous.y {
					collision = true
					break
				}
			}
			if collision {
				continue
			}
			placed = append(placed, candidate)
			fmt.Fprintf(&svg, "<text x=\"%.1f\" y=\"%.1f\" font-size=\"%.0f\" fill=\"%s\">%s</text>", x+4, y+size+2, size, palette[index%len(palette)], html.EscapeString(word.Text))
			break
		}
	}
	if len(placed) == 0 {
		return "", errors.New("no words fit; try shorter words or a larger canvas")
	}
	svg.WriteString("</g></svg>")
	return svg.String(), nil
}
