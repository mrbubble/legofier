// Copyright 2014 Leonardo "Bubble" Mesquita
package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/mrbubble/lego"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"sort"
)

var panelWidth = flag.Uint("panel_width", 64, "Panel width to use")
var outputScale = flag.Int("output_scale", 12, "Size of brick on output")
var dither = flag.Bool("dither", true, "Dither?")
var outline = flag.Bool("outline", true, "Draw outline?")

type elem struct {
	brick lego.Brick
	count int
}
type table []elem

func (t table) Len() int           { return len(t) }
func (t table) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t table) Less(i, j int) bool { return t[i].count > t[j].count }

func main() {
	flag.Parse()
	reader, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	src, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	panel := lego.NewPanel(src,
		&lego.Options{*panelWidth, lego.ALL_BRICKS, *dither})
	out := panel.Draw(*outputScale, *outline)
	var buf bytes.Buffer
	if err := png.Encode(&buf, out); err != nil {
		log.Fatal(err)
	}
	enc := base64.StdEncoding.EncodeToString(buf.Bytes())
	fmt.Printf("<img src=\"data:image/png;base64,%s\">\n", enc)
	fmt.Printf("<br>Dimensions: %v\n", panel.Size())
	bricks := panel.CountBricks()
	total := 0
	var t table
	for brick, count := range bricks {
		total += count
		t = append(t, elem{brick, count})
	}
	sort.Sort(t)
	fmt.Printf("<br>Number of bricks: %v\n", total)
	fmt.Printf("<br>Bricks:\n<table>\n")
	for _, e := range t {
		fmt.Printf("<tr><td>%s</td><td>%d</td></tr>\n", e.brick, e.count)
	}
	fmt.Printf("</table>\n")
}
