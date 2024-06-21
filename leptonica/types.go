package leptonica

type ImageType int32

const (
	UNKNOWN ImageType = iota
	BMP
	JFIF_JPEG
	PNG
	TIFF
	TIFF_PACKBITS
	TIFF_RLE
	TIFF_G3
	TIFF_G4
	TIFF_LZW
	TIFF_ZIP
	PNM
	PS
	GIF
	JP2
	WEBP
	LPDF
	DEFAULT
	SPIX
)

const (
	L_SEVERITY_EXTERNAL = 0 /* Get the severity from the environment   */
	L_SEVERITY_ALL      = 1 /* Lowest severity: print all messages     */
	L_SEVERITY_DEBUG    = 2 /* Print debugging and higher messages     */
	L_SEVERITY_INFO     = 3 /* Print informational and higher messages */
	L_SEVERITY_WARNING  = 4 /* Print warning and higher messages       */
	L_SEVERITY_ERROR    = 5 /* Print error and higher messages         */
	L_SEVERITY_NONE     = 6 /* Highest severity: print no messages     */
)

type GrayCastMode int

const (
	GRAY_SIMPLE GrayCastMode = iota
	GRAY_CAST_REMOVE_COLORS
	GRAY_CAST_KEEP_ONLY_COLORS
	GRAY_CAST_REMOVE_COLORS_2
)

type EnhanceOptions struct {
	TileX    int
	TileY    int
	Thresh   int
	MinCount int
	BgVal    int
	SmoothX  int
	SmoothY  int

	Gamma    float32
	GammaMin int
	GammaMax int

	Factor float32

	RemoveBorders int
}

var DefaultEnhanceOptions = EnhanceOptions{
	TileX:         10,
	TileY:         10,
	Thresh:        40,
	MinCount:      40,
	BgVal:         250,
	SmoothX:       4,
	SmoothY:       4,
	Gamma:         0.5,
	GammaMin:      20,
	GammaMax:      240,
	Factor:        0.8,
	RemoveBorders: 128,
}

type GrayOptions struct {
	Saturation int
	WhitePoint int

	ThreshDiff int
	MinDist    int
}

var DefaultGrayOptions = GrayOptions{
	Saturation: 150,
	WhitePoint: 250,

	ThreshDiff: 90,
	MinDist:    2,
}
