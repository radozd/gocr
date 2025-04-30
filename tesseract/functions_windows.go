package tesseract

// #include <stdlib.h>
// #include <stdint.h>
import "C"
import (
	"math"
	"syscall"
	"unsafe"

	"github.com/radozd/gocr/leptonica"
)

var (
	tessDll = syscall.NewLazyDLL("tesseract53.dll")

	_tessBaseAPICreate         = tessDll.NewProc("TessBaseAPICreate")
	_tessBaseAPIEnd            = tessDll.NewProc("TessBaseAPIEnd")
	_tessBaseAPIDelete         = tessDll.NewProc("TessBaseAPIDelete")
	_tessBaseAPIInit2          = tessDll.NewProc("TessBaseAPIInit2")
	_tessBaseAPISetPageSegMode = tessDll.NewProc("TessBaseAPISetPageSegMode")
	_tessBaseAPISetImage2      = tessDll.NewProc("TessBaseAPISetImage2")

	/* Utility */
	_tessBaseAPISetVariable = tessDll.NewProc("TessBaseAPISetVariable")
	_tessDeleteText         = tessDll.NewProc("TessDeleteText")

	/* Whole text */
	_tessBaseAPIRecognize   = tessDll.NewProc("TessBaseAPIRecognize")
	_tessBaseAPIGetUTF8Text = tessDll.NewProc("TessBaseAPIGetUTF8Text")
	_tessBaseAPIGetHOCRText = tessDll.NewProc("TessBaseAPIGetHOCRText")

	/* Result iterator */
	_tessBaseAPIGetIterator            = tessDll.NewProc("TessBaseAPIGetIterator")
	_tessResultIteratorGetPageIterator = tessDll.NewProc("TessResultIteratorGetPageIterator")
	_tessResultIteratorCopy            = tessDll.NewProc("TessResultIteratorCopy")
	_tessResultIteratorDelete          = tessDll.NewProc("TessResultIteratorDelete")
	_tessResultIteratorNext            = tessDll.NewProc("TessResultIteratorNext")
	_tessResultIteratorGetUTF8Text     = tessDll.NewProc("TessResultIteratorGetUTF8Text")
	_tessResultIteratorConfidence      = tessDll.NewProc("TessResultIteratorConfidence")

	/* Page iterator */
	_tessPageIteratorBegin            = tessDll.NewProc("TessPageIteratorBegin")
	_tessPageIteratorNext             = tessDll.NewProc("TessPageIteratorNext")
	_tessPageIteratorIsAtBeginningOf  = tessDll.NewProc("TessPageIteratorIsAtBeginningOf")
	_tessPageIteratorIsAtFinalElement = tessDll.NewProc("TessPageIteratorIsAtFinalElement")
	_tessPageIteratorBoundingBox      = tessDll.NewProc("TessPageIteratorBoundingBox")
	_tessPageIteratorOrientation      = tessDll.NewProc("TessPageIteratorOrientation")
)

type TessResultIterator struct {
	it uintptr
}

type TessPageIterator struct {
	it uintptr
}

type Api struct {
	handle uintptr
}

var NullApi Api = Api{handle: 0}

func tessBaseAPICreate() Api {
	h, _, _ := _tessBaseAPICreate.Call()
	if h == 0 {
		return Api{}
	}

	return Api{handle: h}
}

func tessBaseAPIEnd(api Api) {
	_tessBaseAPIEnd.Call(api.handle)
}

func tessBaseAPIDelete(api Api) {
	_tessBaseAPIDelete.Call(api.handle)
}

func tessBaseAPIInit2(api Api, datapath string, lang string, mode TessOcrEngineMode) {
	cDatapath := C.CString(datapath)
	defer C.free(unsafe.Pointer(cDatapath))

	cLang := C.CString(lang)
	defer C.free(unsafe.Pointer(cLang))

	_tessBaseAPIInit2.Call(api.handle, uintptr(unsafe.Pointer(cDatapath)), uintptr(unsafe.Pointer(cLang)), uintptr(mode))
}

func tessBaseAPISetPageSegMode(api Api, mode TessPageSegMode) {
	_tessBaseAPISetPageSegMode.Call(api.handle, uintptr(mode))
}

func tessBaseAPISetImage2(api Api, pix leptonica.Pix) {
	_tessBaseAPISetImage2.Call(api.handle, leptonica.UnsafePix(pix))
}

/* Utility */
func tessBaseAPISetVariable(api Api, name string, value string) bool {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	code, _, _ := _tessBaseAPISetVariable.Call(api.handle, uintptr(unsafe.Pointer(cName)), uintptr(unsafe.Pointer(cValue)))
	return code != 0
}

func tessDeleteText(text *C.char) {
	_tessDeleteText.Call(uintptr(unsafe.Pointer(text)))
}

/* Whole text */
func tessBaseAPIRecognize(api Api) int {
	code, _, _ := _tessBaseAPIRecognize.Call(api.handle, 0)
	return int(code)
}

func tessBaseAPIGetUTF8Text(api Api) *C.char {
	text, _, _ := _tessBaseAPIGetUTF8Text.Call(api.handle)
	return (*C.char)(unsafe.Pointer(text))
}

func tessBaseAPIGetHOCRText(api Api, page_number int) *C.char {
	text, _, _ := _tessBaseAPIGetHOCRText.Call(api.handle, uintptr(page_number))
	return (*C.char)(unsafe.Pointer(text))
}

/* Result iterator */
func tessBaseAPIGetIterator(api Api) TessResultIterator {
	h, _, _ := _tessBaseAPIGetIterator.Call(api.handle)
	return TessResultIterator{it: h}
}

func tessResultIteratorGetPageIterator(handle TessResultIterator) TessPageIterator {
	h, _, _ := _tessResultIteratorGetPageIterator.Call(handle.it)
	return TessPageIterator{it: h}
}

func tessResultIteratorCopy(handle TessResultIterator) TessResultIterator {
	h, _, _ := _tessResultIteratorCopy.Call(handle.it)
	return TessResultIterator{it: h}
}

func tessResultIteratorDelete(handle TessResultIterator) {
	_tessResultIteratorDelete.Call(handle.it)
}

func tessResultIteratorNext(handle TessResultIterator, level TessPageIteratorLevel) bool {
	code, _, _ := _tessResultIteratorNext.Call(handle.it, uintptr(level))
	return code != 0
}

func tessResultIteratorGetUTF8Text(handle TessResultIterator, level TessPageIteratorLevel) *C.char {
	text, _, _ := _tessResultIteratorGetUTF8Text.Call(handle.it, uintptr(level))
	return (*C.char)(unsafe.Pointer(text))
}

func tessResultIteratorConfidence(handle TessResultIterator, level TessPageIteratorLevel) float32 {
	_, confidence, _ := _tessResultIteratorConfidence.Call(handle.it, uintptr(level))
	c := math.Float32frombits(uint32(confidence))

	return c
}

/* Page iterator */
func tessPageIteratorBegin(handle TessPageIterator) {
	_tessPageIteratorBegin.Call(handle.it)
}

func tessPageIteratorNext(handle TessPageIterator, level TessPageIteratorLevel) bool {
	code, _, _ := _tessPageIteratorNext.Call(handle.it, uintptr(level))
	return code != 0
}

func tessPageIteratorIsAtBeginningOf(handle TessPageIterator, level TessPageIteratorLevel) bool {
	code, _, _ := _tessPageIteratorIsAtBeginningOf.Call(handle.it, uintptr(level))
	return code != 0
}

func tessPageIteratorIsAtFinalElement(handle TessPageIterator, level TessPageIteratorLevel, element TessPageIteratorLevel) bool {
	code, _, _ := _tessPageIteratorIsAtFinalElement.Call(handle.it, uintptr(level), uintptr(element))
	return code != 0
}

func tessPageIteratorBoundingBox(handle TessPageIterator, level TessPageIteratorLevel) (int, int, int, int) {
	l := C.int(0)
	t := C.int(0)
	r := C.int(0)
	b := C.int(0)

	res, _, _ := _tessPageIteratorBoundingBox.Call(handle.it, uintptr(level),
		uintptr(unsafe.Pointer(&l)), uintptr(unsafe.Pointer(&t)), uintptr(unsafe.Pointer(&r)), uintptr(unsafe.Pointer(&b)))
	if res != 0 {
		return int(l), int(t), int(r), int(b)
	}
	return 0, 0, 0, 0
}

func tessPageIteratorOrientation(handle TessPageIterator) TessPageOrientation {
	orientation := C.int(0)
	writing_direction := C.int(0)
	textline_order := C.int(0)
	deskew_angle := C.float(0)

	_tessPageIteratorOrientation.Call(handle.it,
		uintptr(unsafe.Pointer(&orientation)),
		uintptr(unsafe.Pointer(&writing_direction)),
		uintptr(unsafe.Pointer(&textline_order)),
		uintptr(unsafe.Pointer(&deskew_angle)))

	return TessPageOrientation(orientation)
}
