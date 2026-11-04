package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	wordcloud "github.com/YananLee/wordcloud-art-svg-go"
)

func main() {
	text := flag.String("text", "", "text to turn into a word cloud (reads stdin if empty)")
	output := flag.String("out", "", "output SVG file (writes stdout if empty)")
	flag.Parse()
	if *text == "" {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			fail(err)
		}
		*text = string(input)
	}
	svg, err := wordcloud.FromText(*text, wordcloud.Options{})
	if err != nil {
		fail(err)
	}
	if *output == "" {
		fmt.Println(svg)
		return
	}
	if err := os.WriteFile(*output, []byte(svg), 0644); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
