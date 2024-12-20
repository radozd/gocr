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

const (
	L_SELECT_WIDTH     = 1 /*!< width must satisfy constraint         */
	L_SELECT_HEIGHT    = 2 /*!< height must satisfy constraint        */
	L_SELECT_XVAL      = 3 /*!< x value must satisfy constraint       */
	L_SELECT_YVAL      = 4 /*!< y value must satisfy constraint       */
	L_SELECT_IF_EITHER = 5 /*!< either width or height (or xval       */
	/*!< or yval) can satisfy constraint       */
	L_SELECT_IF_BOTH = 6 /*!< both width and height (or xval        */
	/*!< and yval must satisfy constraint      */
)

const (
	L_SELECT_IF_LT  = 1 /* save if value is less than threshold  */
	L_SELECT_IF_GT  = 2 /* save if value is more than threshold  */
	L_SELECT_IF_LTE = 3 /* save if value is <= to the threshold  */
	L_SELECT_IF_GTE = 4 /* save if value is >= to the threshold  */
)

const (
	L_INSERT     = 0 /* stuff it in; no copy, clone or copy-clone    */
	L_COPY       = 1 /* make/use a copy of the object                */
	L_CLONE      = 2 /* make/use clone (ref count) of the object     */
	L_COPY_CLONE = 3 /* make a new object and fill with with clones  */
	/* of each object in the array(s)               */
)

type EnhanceOptions struct {
	TileX      int
	TileY      int
	Thresh     int
	MinCount   int
	WhitePoint int
	SmoothX    int
	SmoothY    int

	Gamma    int
	GammaMin int
	GammaMax int

	Factor int

	RemoveBorders int
}

var DefaultEnhanceOptions = EnhanceOptions{
	TileX:         10,
	TileY:         10,
	Thresh:        40,
	MinCount:      40,
	WhitePoint:    250,
	SmoothX:       4,
	SmoothY:       4,
	Gamma:         50,
	GammaMin:      20,
	GammaMax:      240,
	Factor:        80,
	RemoveBorders: 180,
}

type GrayCastMode int

const (
	GRAY_SIMPLE GrayCastMode = iota
	GRAY_CAST_REMOVE_COLORS
	GRAY_CAST_KEEP_ONLY_COLORS
	GRAY_CAST_REMOVE_COLORS_2
)

type GrayOptions struct {
	Saturation int
	WhitePoint int

	ThreshDiff int
	MinDist    int

	MaximizeBrightness bool
}

var DefaultGrayOptions = GrayOptions{
	Saturation: 150,
	WhitePoint: 250,

	ThreshDiff: 90,
	MinDist:    2,

	MaximizeBrightness: true,
}

type MaskOptions struct {
	Thresh int

	SqrBlock int
	SqrMin   int
	SqrMax   int

	BarMin    int
	BarMax    int
	BarWidth  int
	BarHeight int

	LinMin int

	SpMax    int
	SpWeight int
}

var DefaultMaskOptions = MaskOptions{
	Thresh: 90,

	SqrBlock: 15,
	SqrMin:   75,
	SqrMax:   450,

	BarMin:    2,
	BarMax:    20,
	BarWidth:  40,
	BarHeight: 80,

	LinMin: 59,

	SpMax:    10,
	SpWeight: 180,
}
