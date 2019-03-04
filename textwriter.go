package combglasspattern

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func WriteTextfile(filenameprefix string, variant Variant) error {
	var (
		contents bytes.Buffer
		err      error
	)
	filename := fmt.Sprintf("%s_drillpoints.txt", filenameprefix)

	_, err = fmt.Fprintf(&contents, "// Board dimensions: %.1fmm x %.1fmm\n", variant.Board.Width, variant.Board.Height)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(&contents, "// Glass dimensions: inner Ø: %.1fmm ; outer Ø: %.1fmm\n", variant.Glass.InnerRadius*2, variant.Glass.OuterRadius*2)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(&contents, "")
	if err != nil {
		return err
	}
	for _, p := range variant.Points {
		_, _ = fmt.Fprintf(&contents, "%6.1f → %6.1f ↓\n", p.X, p.Y)
	}

	err = ioutil.WriteFile(filename, contents.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
