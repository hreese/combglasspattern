package combglasspattern

type BoardConfiguration struct {
	Width, Height   float64
	WallOffset      float64
	MinHoleDistance float64
}

type GlassConfiguration struct {
	InnerRadius   float64
	OuterRadius   float64
	NumberOfSides int
}

type Variant struct {
	Points []Point
	Board  BoardConfiguration
	Glass  GlassConfiguration
}

var (
	PresetsBoard = map[string]BoardConfiguration{
		"DadantWeber": {
			Width:           464,
			Height:          464,
			WallOffset:      10,
			MinHoleDistance: 10,
		},
		"ZanderSpec": {
			Width:           435,
			Height:          380,
			WallOffset:      10,
			MinHoleDistance: 10,
		},
		"TestBrett": {
			Width:           500,
			Height:          600,
			WallOffset:      10,
			MinHoleDistance: 10,
		},
		"DemoBrettA4": {
			Width:           210 - 20,
			Height:          297 - 20,
			WallOffset:      5,
			MinHoleDistance: 5,
		},
	}
	PresetGlas = map[string]GlassConfiguration{
		// https://www.holtermann-glasshop.de/Sechseckglaeser/Sechseckglas-580-ml/
		"HolterMannTwistOffSechseckglas580": {
			InnerRadius:   82 / 2,
			OuterRadius:   95 / 2,
			NumberOfSides: 6,
		},
		// https://www.holtermann-glasshop.de/Designglaeser/Viereckglas-312-ml/Viereckglas-312-ml-Biene.html
		"HolterMannTwistOffViereckglas312": {
			InnerRadius:   60 / 2,
			OuterRadius:   75 / 2,
			NumberOfSides: 4,
		},
		// https://www.flaschenbauer.de/einmachglaeser/sechskantglaeser/sechskantglas-580-ml-to-82
		"FlaschenBauerSechskantglas580mlTO82": {
			InnerRadius:   82 / 2,
			OuterRadius:   95 / 2,
			NumberOfSides: 6,
		},
		// https://www.bienen-ruck.de/imkershop/honigverkauf-werbemittel/twist-off-glaeser/1902/wabenglaeser-rund
		"BienenRuckWabengläserRund500": {
			InnerRadius:   82 / 2,
			OuterRadius:   90 / 2,
			NumberOfSides: 0,
		},
		"TestGlas": {
			InnerRadius:   60 / 2,
			OuterRadius:   88 / 2,
			NumberOfSides: 0,
		},
		"DemoGlasEckig": {
			InnerRadius:   45 / 2,
			OuterRadius:   54 / 2,
			NumberOfSides: 8,
		},
	}
)

func (board BoardConfiguration) CenterPoint() Point {
	return Point{board.Width / 2, board.Height / 2}
}

func (board BoardConfiguration) IsSquare() bool {
	return board.Width == board.Height
}

func (board BoardConfiguration) Turn90() BoardConfiguration {
	board.Width, board.Height = board.Height, board.Width
	return board
}
