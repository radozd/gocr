package tesseract

// #include <stdlib.h>
import "C"

import (
	"unsafe"

	"github.com/radozd/gocr/leptonica"
)

type TextBlock struct {
	Goodness float32
	Value    string
	X        int
	Y        int
	Width    int
	Height   int
}

type Api struct {
	handle uintptr
}

type ResultIterator uintptr
type PageIterator uintptr

type PageIteratorLevel int

const (
	RIL_BLOCK PageIteratorLevel = iota
	RIL_PARA
	RIL_TEXTLINE
	RIL_WORD
	RIL_SYMBOL
)

type OcrEngineMode int

const (
	OEM_TESSERACT_ONLY OcrEngineMode = iota
	OEM_LSTM_ONLY
	OEM_TESSERACT_LSTM_COMBINED
	OEM_DEFAULT
)

type PageSegMode int

const (
	PSM_OSD_ONLY PageSegMode = iota
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

var DataPath string = "."

func NewApi(lang string) *Api {
	api, _, _ := tessCreate.Call()
	if api == 0 {
		return nil
	}

	cDatapath := C.CString(DataPath)
	defer C.free(unsafe.Pointer(cDatapath))

	// Set the language
	cLang := C.CString(lang)
	defer C.free(unsafe.Pointer(cLang))

	mode := OEM_DEFAULT
	_, _, _ = tessInit2.Call(api, uintptr(unsafe.Pointer(cDatapath)), uintptr(unsafe.Pointer(cLang)), uintptr(mode))

	return &Api{handle: api}
}

func (api *Api) Close() {
	tessBaseAPIEnd.Call(api.handle)
	tessDelete.Call(api.handle)
	api.handle = 0
}

func (api *Api) SetVariable(name string, val string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))

	tessSetVariable.Call(api.handle, uintptr(unsafe.Pointer(cName)), uintptr(unsafe.Pointer(cVal)))
}

func (api *Api) SetImagePix(pix leptonica.Pix) {
	tessSetImage2.Call(api.handle, uintptr(pix))
}

func (api *Api) SetPageSegMode(mode PageSegMode) {
	tessSetPageSegMode.Call(api.handle, uintptr(mode))
}

func (api *Api) GetIterator() ResultIterator {
	resIt, _, _ := tessGetIterator.Call(api.handle)
	return ResultIterator(resIt)
}

// /////////////////////////////////////
func (api *Api) Recognize() {
	par := 0
	tessRecognize.Call(api.handle, uintptr(par))
}

func (api *Api) Text() string {
	text, _, _ := tessGetUTF8Text.Call(api.handle)
	if text == 0 {
		return ""
	}

	res := C.GoString((*C.char)(unsafe.Pointer(text)))
	tessFreeUTF8Text.Call(text)

	return res
}

// HOCRText returns the HOCR text for given pagenumber
func (api *Api) HOCRText(pagenumber int) string {
	text, _, _ := tessGetHOCRText.Call(api.handle, uintptr(pagenumber))
	if text == 0 {
		return ""
	}

	res := C.GoString((*C.char)(unsafe.Pointer(text)))
	tessFreeUTF8Text.Call(text)

	return res
}

func (api *Api) TextBlocks(level PageIteratorLevel) []TextBlock {
	resIt := api.GetIterator()
	defer resIt.Delete()

	pageIt := resIt.GetPageIterator()

	blocks := make([]TextBlock, 0)

	good := true
	for good {
		if pageIt.IsAtBeginningOf(level) {
			x1, y1, x2, y2 := pageIt.BoundingBox(level)
			text, goodness := resIt.GetUTF8Text(level)
			blocks = append(blocks, TextBlock{
				Goodness: goodness,
				Value:    text,
				X:        x1,
				Y:        y1,
				Width:    x2 - x1,
				Height:   y2 - y1,
			})
		}

		good = pageIt.Next(level)
	}
	return blocks
}
