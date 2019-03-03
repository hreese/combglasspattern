package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ajstarks/svgo/float"
	. "gitlab.com/hreese/combglasspattern"
)

type Variant struct {
	points []Point
	board  BoardConfiguration
	glass  GlassConfiguration
}

func main() {
	var (
		f        *os.File
		err      error
		canvas   *svg.SVG
		variants = make(map[string]Variant)
		//board    = PresetsBoard["DadantWeber"]
		//glass    = PresetGlas["BienenRuckWabengläserRund500"]
		board = PresetsBoard["ZanderSpec"]
		glass = PresetGlas["TestGlas"]
	)

	square, hexOne, hexTwo := GenerateHoles(board, glass)
	variants[`Square`] = Variant{square, board, glass}
	variants[`HexOne`] = Variant{hexOne, board, glass}
	variants[`HexTwo`] = Variant{hexTwo, board, glass}
	fmt.Printf("Sqare:   %d\nHexOne:  %d\nHexTwo:  %d\n", len(square), len(hexOne), len(hexTwo))
	if !board.IsSquare() {
		turnedBoard := board.Turn90()
		_, hexThree, hexFour := GenerateHoles(turnedBoard, glass)
		variants[`HexThree`] = Variant{hexThree, turnedBoard, glass}
		variants[`HexFour`] = Variant{hexFour, turnedBoard, glass}
		fmt.Printf("HexThree:  %d\nHexFour:  %d\n", len(hexThree), len(hexFour))
	}

	for filename, variant := range variants {
		var pointsAsText bytes.Buffer
		_, _ = fmt.Fprintf(&pointsAsText, "// board dimensions: %.1fmm x %.1fmm\n", board.Width, board.Height)
		_, _ = fmt.Fprintf(&pointsAsText, "// glass dimensions: inner Ø: %.1fmm ; outer Ø: %.1fmm\n", glass.InnerRadius*2, glass.OuterRadius*2)
		_, _ = fmt.Fprintln(&pointsAsText, "")
		// open svg file
		f, err = os.OpenFile(fmt.Sprintf("%s.svg", filename), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		canvas = svg.New(f)
		canvas.Startunit(variant.board.Width, variant.board.Height, "mm",
			fmt.Sprintf(`viewBox="0 0 %f %f"`, variant.board.Width, variant.board.Height))
		DrawBoard(canvas, variant.board, true)
		canvas.Group("Holes")
		for _, p := range variant.points {
			Throughhole(canvas, p, variant.glass.InnerRadius, variant.glass.OuterRadius, variant.glass.NumberOfSides)
			_, _ = fmt.Fprintf(&pointsAsText, "%6.1f → %6.1f ↓\n", p.X, p.Y)
		}
		canvas.Gend()
		canvas.End()

		err = ioutil.WriteFile(fmt.Sprintf("drillpoints_%s.txt", filename), pointsAsText.Bytes(), 0644)
		if err != nil {
			log.Println(err)
		}
	}

}
