# Wordcloud.art SVG for Go

A small, dependency-free Go library and command-line tool that turns text or weighted words into a simple SVG word cloud. It works locally; input text is not sent to a service.

This package makes a basic horizontal word layout. For shapes, animation, design controls, and more export formats, use [Wordcloud.art](https://www.wordcloud.art/).

## Install

```sh
go get github.com/YananLee/wordcloud-art-svg-go@v0.1.0
```

## Go API

```go
package main

import (
    "os"

    wordcloud "github.com/YananLee/wordcloud-art-svg-go"
)

func main() {
    svg, err := wordcloud.FromText("design design words colors", wordcloud.Options{})
    if err != nil {
        panic(err)
    }
    if err := os.WriteFile("wordcloud.svg", []byte(svg), 0644); err != nil {
        panic(err)
    }
}
```

For explicit weights, call `FromWords([]wordcloud.Word{{Text: "design", Weight: 5}}, wordcloud.Options{})`. Options allow `Width`, `Height`, and `MaxWords`; zero values use 800×480 pixels and 60 words.

## Command line

```sh
go run github.com/YananLee/wordcloud-art-svg-go/cmd/wordcloud-art-svg@v0.1.0 -text "design design words" -out wordcloud.svg
```

You can also omit `-text` and pipe text to stdin. Without `-out`, SVG is written to stdout.

## License

MIT.
