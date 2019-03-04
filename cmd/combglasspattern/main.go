package main

import (
	"fmt"
	. "gitlab.com/hreese/combglasspattern"
	"log"
)

func main() {
	var (
		err      error
		variants = make(map[string]Variant)
		board    = PresetsBoard["DadantWeber"]
		glass    = PresetGlas["BienenRuckWabengläserRund500"]
		//board = PresetsBoard["ZanderSpec"]
		//glass = PresetGlas["TestGlas"]
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
		err = WriteTextfile(filename, variant)
		if err != nil {
			log.Fatal(err)
		}
		err = WriteSVG(filename, variant)
		if err != nil {
			log.Fatal(err)
		}
		err = WritePDF(filename, variant)
		if err != nil {
			log.Fatal(err)
		}
	}
}
