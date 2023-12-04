package tesseract

// #include <stdlib.h>
import "C"

import (
	"math"
	"unsafe"

	"github.com/radozd/gocr/leptonica"
)

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

func NewApi(lang string) *Api {
	api, _, _ := tessCreate.Call()
	if api == 0 {
		return nil
	}

	datapath := "./"
	cDatapath := C.CString(datapath)
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

// /////////////////////////////////////
func (resIt ResultIterator) GetPageIterator() PageIterator {
	pageIt, _, _ := tessResultIteratorGetPageIterator.Call(uintptr(resIt))
	return PageIterator(pageIt)
}

func (resIt ResultIterator) Delete() {
	tessResultIteratorDelete.Call(uintptr(resIt))
}

func (resIt ResultIterator) GetUTF8Text(level PageIteratorLevel) (string, float32) {
	text, _, _ := tessResultIteratorGetUTF8Text.Call(uintptr(resIt), uintptr(level))
	if text == 0 {
		return "", 0.0
	}

	res := C.GoString((*C.char)(unsafe.Pointer(text)))
	tessFreeUTF8Text.Call(text)

	_, confidence, _ := tessResultIteratorConfidence.Call(uintptr(resIt), uintptr(level))
	c := math.Float32frombits(uint32(confidence))

	return res, c
}

// /////////////////////////////////////
func (pageIt PageIterator) Begin() {
	tessPageIteratorBegin.Call(uintptr(pageIt))
}

func (pageIt PageIterator) BoundingBox(level PageIteratorLevel) (int, int, int, int) {
	cL := C.int(0)
	cT := C.int(0)
	cR := C.int(0)
	cB := C.int(0)

	res, _, _ := tessPageIteratorBoundingBox.Call(uintptr(pageIt), uintptr(level),
		uintptr(unsafe.Pointer(&cL)), uintptr(unsafe.Pointer(&cT)), uintptr(unsafe.Pointer(&cR)), uintptr(unsafe.Pointer(&cB)))
	if res != 0 {
		return int(cL), int(cT), int(cR), int(cB)
	}
	return 0, 0, 0, 0
}

func (pageIt PageIterator) IsAtBeginningOf(level PageIteratorLevel) bool {
	res, _, _ := tessPageIteratorIsAtBeginningOf.Call(uintptr(pageIt), uintptr(level))
	return res != 0
}

func (pageIt PageIterator) IsAtFinalElement(level PageIteratorLevel, element PageIteratorLevel) bool {
	res, _, _ := tessPageIteratorIsAtFinalElement.Call(uintptr(pageIt), uintptr(level), uintptr(element))
	return res != 0
}

func (pageIt PageIterator) Next(level PageIteratorLevel) bool {
	res, _, _ := tessPageIteratorNext.Call(uintptr(pageIt), uintptr(level))
	return res != 0
}
