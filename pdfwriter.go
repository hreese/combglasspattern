package combglasspattern

import (
	"fmt"
	"math"

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
	pdf.SetLineWidth(0)
	pdf.SetFillColor(0xdd, 0xdd, 0xdd)
	pdf.SetAlpha(0.2890625, "Normal")
	if numsides < 3 {
		pdf.Circle(x, y, outerRadius, "DF")
	} else {
		angleStep := math.Pi * 2 / float64(numsides)
		pdf.MoveTo(x+outerRadius, y)
		for step := 1; step < numsides; step++ {
			pdf.LineTo(x+outerRadius * math.Cos(float64(step)*angleStep), y+outerRadius * math.Sin(float64(step)*angleStep))
		}
		pdf.ClosePath()
		pdf.DrawPath("DF")
	}

	// hole outline
	pdf.SetAlpha(1.0, "Normal")
	//pdf.SetFillColor(0xff, 0xff, 0xff)
	pdf.Circle(x, y, innerRadius, "DF")

	// crossmark
	crossoffset := innerRadius / 12
	pdf.Line(x-crossoffset, y-crossoffset, x+crossoffset, y+crossoffset)
	pdf.Line(x+crossoffset, y-crossoffset, x-crossoffset, y+crossoffset)
	cellwidth := 200.0
	cellheight := 16.0
	celltext := fmt.Sprintf("x: %.1fmm y: %.1fmm", x, y)
	//celltext := fmt.Sprintf("→%.1fmm ↓%.1fmm", x, y)
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	// coordinates
	pdf.MoveTo(x-cellwidth, y-cellheight)
	pdf.CellFormat(2*cellwidth, cellheight, tr(celltext), "", 0, "C", false, 0 , "")
	// diameter
	pdf.MoveTo(x-cellwidth, y)
	pdf.CellFormat(2*cellwidth, cellheight, tr(fmt.Sprintf("diam: %.1fmm", innerRadius*2)), "", 0, "C", false, 0 , "")

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
