package tesseract

type TextBlock struct {
	Goodness float32
	Value    string
	X        int
	Y        int
	Width    int
	Height   int
}

type TessPageIteratorLevel int

const (
	RIL_BLOCK TessPageIteratorLevel = iota
	RIL_PARA
	RIL_TEXTLINE
	RIL_WORD
	RIL_SYMBOL
)

type TessOcrEngineMode int

const (
	OEM_TESSERACT_ONLY TessOcrEngineMode = iota
	OEM_LSTM_ONLY
	OEM_TESSERACT_LSTM_COMBINED
	OEM_DEFAULT
)

type TessPageSegMode int

const (
	PSM_OSD_ONLY TessPageSegMode = iota
	PSM_AUTO_OSD
	PSM_AUTO_ONLY
	PSM_AUTO
	PSM_SINGLE_COLUMN
	PSM_SINGLE_BLOCK_VERT_TEXT
	PSM_SINGLE_BLOCK
	PSM_SINGLE_LINE
	PSM_SINGLE_WORD
	PSM_CIRCLE_WORD
	PSM_SINGLE_CHAR
	PSM_SPARSE_TEXT
	PSM_SPARSE_TEXT_OSD
	PSM_RAW_LINE
	PSM_COUNT
)

type TessPageOrientation int

const (
	ORIENTATION_PAGE_UP TessPageOrientation = iota
	ORIENTATION_PAGE_RIGHT
	ORIENTATION_PAGE_DOWN
	ORIENTATION_PAGE_LEFT
)
