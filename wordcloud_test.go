package wordcloud

import (
	"encoding/xml"
	"io"
	"strings"
	"testing"
)

func TestFromTextGeneratesParseableUnicodeSVG(t *testing.T) {
	svg, err := FromText("design design art 词云 词云", Options{})
	if err != nil {
		t.Fatal(err)
	}
	decoder := xml.NewDecoder(strings.NewReader(svg))
	var labels []string
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("invalid SVG XML: %v", err)
		}
		if text, ok := token.(xml.CharData); ok {
			labels = append(labels, string(text))
		}
	}
	joined := strings.Join(labels, " ")
	if !strings.Contains(joined, "design") || !strings.Contains(joined, "词云") {
		t.Fatalf("expected both Latin and Chinese labels, got %q", joined)
	}
}

func TestFromWordsEscapesLabels(t *testing.T) {
	svg, err := FromWords([]Word{{Text: "A&B", Weight: 3}}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "A&amp;B") {
		t.Fatalf("label was not XML escaped: %s", svg)
	}
}

func TestEmptyInputFails(t *testing.T) {
	if _, err := FromText("  -  ", Options{}); err == nil {
		t.Fatal("expected an error for input without words")
	}
}
