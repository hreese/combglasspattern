package combglasspattern

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func PDFDrawBoard(pdf *gofpdf.Fpdf, board BoardConfiguration, showOrigin bool) {
	// save old FillColor
	var (
		oldred, oldgreen, oldblue = pdf.GetFillColor()
		oldlinewidth              = pdf.GetLineWidth()
	)

	// draw outer rect with border
	pdf.SetFillColor(0xdd, 0xdd, 0xdd)

	pdf.MoveTo(0, 0)
	pdf.LineTo(board.Width, 0)
	pdf.LineTo(board.Width, board.Height)
	pdf.LineTo(0, board.Height)
	pdf.ClosePath()

	pdf.MoveTo(board.WallOffset, board.WallOffset)
	pdf.LineTo(board.WallOffset, board.Height-board.WallOffset)
	pdf.LineTo(board.Width-board.WallOffset, board.Height-board.WallOffset)
	pdf.LineTo(board.Width-board.WallOffset, board.WallOffset)
	pdf.ClosePath()

	pdf.DrawPath("f*")

	if showOrigin {
		pdf.SetFillColor(0, 0, 0)
		pdf.Polygon([]gofpdf.PointType{
			gofpdf.PointType{0, 0},
			gofpdf.PointType{board.WallOffset, 0},
			gofpdf.PointType{0, board.WallOffset},
			//gofpdf.PointType{0, 0},
		}, "F")
	}

	// restore old FillColor and Linewidth
	pdf.SetFillColor(oldred, oldgreen, oldblue)
	pdf.SetLineWidth(oldlinewidth)
}

func PDFThroughhole(pdf *gofpdf.Fpdf, center Point, innerRadius, outerRadius float64, numsides int) {
	// save old FillColor
	var (
		oldred, oldgreen, oldblue = pdf.GetFillColor()
		oldlinewidth              = pdf.GetLineWidth()
		x, y                      = center.X, center.Y
	)

	// glass outline
	if numsides < 3 {
		pdf.SetLineWidth(0)
		pdf.SetFillColor(0xbd, 0xbd, 0xbd)
		pdf.SetAlpha(0.2890625, "Screen")
		pdf.Circle(x, y, outerRadius, "DF")
	} else {
		// TODO
	}

	// hole outline
	pdf.SetAlpha(1.0, "Normal")
	pdf.SetFillColor(0xff, 0xff, 0xff)
	pdf.Circle(x, y, innerRadius, "DF")

	// crossmark
	crossoffset := innerRadius / 12
	pdf.Line(x-crossoffset, y-crossoffset, x+crossoffset, y+crossoffset)
	pdf.Line(x+crossoffset, y-crossoffset, x-crossoffset, y+crossoffset)
	cellwidth := 200.0
	cellheight := 16.0
	celltext := fmt.Sprintf("=> %.1fmm v%.1fmm", x, y)
	pdf.MoveTo(x-cellwidth, y-cellheight)
	pdf.CellFormat(2*cellwidth, cellheight, celltext, "", 0, "C", false, 0 , "")

	// restore old FillColor and Linewidth
	pdf.SetFillColor(oldred, oldgreen, oldblue)
	pdf.SetLineWidth(oldlinewidth)
}

func WritePDF(filenameprefix string, variant Variant) error {
	var (
		//err      error
		filename string
	)
	filename = fmt.Sprintf("%s.pdf", filenameprefix)

	// new pdf
	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "mm",
		SizeStr:        "A4",
		Size: gofpdf.SizeType{
			Ht: variant.Board.Height,
			Wd: variant.Board.Width,
		},
		FontDirStr: "",
	})

	// metadata
	pdf.SetAuthor("https://gitlab.com/hreese/combglasspattern", false)
	pdf.SetDisplayMode("fullpage", "single")
	pdf.SetMargins(0, 0, 0)
	pdf.SetFont("Helvetica", "", 14)

	pdf.AddPage()

	PDFDrawBoard(pdf, variant.Board, true)
	for _, p := range variant.Points {
		PDFThroughhole(pdf, p, variant.Glass.InnerRadius, variant.Glass.OuterRadius, variant.Glass.NumberOfSides)
	}

	return pdf.OutputFileAndClose(filename)
}
