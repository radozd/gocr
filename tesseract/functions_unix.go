//go:build darwin || linux
// +build darwin linux

package tesseract

// #cgo CXXFLAGS: -std=c++0x
// #cgo CFLAGS: -I/opt/homebrew/include -I/usr/local/include
// #cgo CFLAGS: -Wno-unused-result
// #cgo LDFLAGS: -L/opt/homebrew/lib -L/usr/local/lib -ltesseract
// #include <stdlib.h>
// #include <leptonica/allheaders.h>
// #include <tesseract/capi.h>
import "C"
import (
	"unsafe"

	"github.com/radozd/gocr/leptonica"
)

type TessResultIterator struct {
	it *C.TessResultIterator
}

type TessPageIterator struct {
	it *C.TessPageIterator
}

type Api struct {
	handle *C.TessBaseAPI
}

var NullApi Api = Api{handle: nil}

func tessBaseAPICreate() Api {
	return Api{handle: C.TessBaseAPICreate()}
}

func tessBaseAPIEnd(api Api) {
	C.TessBaseAPIEnd(api.handle)
}

func tessBaseAPIDelete(api Api) {
	C.TessBaseAPIDelete(api.handle)
}

func tessBaseAPIInit2(api Api, datapath string, lang string, mode TessOcrEngineMode) {
	cDatapath := C.CString(datapath)
	defer C.free(unsafe.Pointer(cDatapath))

	cLang := C.CString(lang)
	defer C.free(unsafe.Pointer(cLang))

	C.TessBaseAPIInit2(api.handle, cDatapath, cLang, C.TessOcrEngineMode(mode))

}

func tessBaseAPISetPageSegMode(api Api, mode TessPageSegMode) {
	C.TessBaseAPISetPageSegMode(api.handle, C.TessPageSegMode(mode))
}

func tessBaseAPISetImage2(api Api, pix leptonica.Pix) {
	C.TessBaseAPISetImage2(api.handle, (*C.PIX)(unsafe.Pointer(leptonica.UnsafePix(pix))))
}

/* Utility */
func tessBaseAPISetVariable(api Api, name string, value string) bool {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	return C.TessBaseAPISetVariable(api.handle, cName, cValue) != 0
}

func tessDeleteText(text *C.char) {
	C.TessDeleteText(text)
}

/* Whole text */
func tessBaseAPIRecognize(api Api) int {
	return int(C.TessBaseAPIRecognize(api.handle, nil))
}

func tessBaseAPIGetUTF8Text(api Api) *C.char {
	return C.TessBaseAPIGetUTF8Text(api.handle)
}

func tessBaseAPIGetHOCRText(api Api, page_number int) *C.char {
	return C.TessBaseAPIGetHOCRText(api.handle, C.int32_t(page_number))
}

/* Result iterator */
func tessBaseAPIGetIterator(api Api) TessResultIterator {
	return TessResultIterator{it: C.TessBaseAPIGetIterator(api.handle)}
}

func tessResultIteratorGetPageIterator(handle TessResultIterator) TessPageIterator {
	return TessPageIterator{it: C.TessResultIteratorGetPageIterator(handle.it)}
}

func tessResultIteratorDelete(handle TessResultIterator) {
	C.TessResultIteratorDelete(handle.it)
}

func tessResultIteratorNext(handle TessResultIterator, level TessPageIteratorLevel) bool {
	return C.TessResultIteratorNext(handle.it, C.TessPageIteratorLevel(level)) != 0
}

func tessResultIteratorGetUTF8Text(handle TessResultIterator, level TessPageIteratorLevel) *C.char {
	return C.TessResultIteratorGetUTF8Text(handle.it, C.TessPageIteratorLevel(level))
}

func tessResultIteratorConfidence(handle TessResultIterator, level TessPageIteratorLevel) float32 {
	return float32(C.TessResultIteratorConfidence(handle.it, C.TessPageIteratorLevel(level)))
}

/* Page iterator */
func tessPageIteratorDelete(handle TessPageIterator) {
	C.TessPageIteratorDelete(handle.it)
}

func tessPageIteratorBegin(handle TessPageIterator) {
	C.TessPageIteratorBegin(handle.it)
}

func tessPageIteratorNext(handle TessPageIterator, level TessPageIteratorLevel) bool {
	return C.TessPageIteratorNext(handle.it, C.TessPageIteratorLevel(level)) != 0
}

func tessPageIteratorIsAtBeginningOf(handle TessPageIterator, level TessPageIteratorLevel) bool {
	return C.TessPageIteratorIsAtBeginningOf(handle.it, C.TessPageIteratorLevel(level)) != 0
}

func tessPageIteratorIsAtFinalElement(handle TessPageIterator, level TessPageIteratorLevel, element TessPageIteratorLevel) bool {
	return C.TessPageIteratorIsAtFinalElement(handle.it, C.TessPageIteratorLevel(level), C.TessPageIteratorLevel(element)) != 0
}

func tessPageIteratorBoundingBox(handle TessPageIterator, level TessPageIteratorLevel) (int, int, int, int) {
	l := C.int(0)
	t := C.int(0)
	r := C.int(0)
	b := C.int(0)

	res := C.TessPageIteratorBoundingBox(handle.it, C.TessPageIteratorLevel(level), &l, &t, &r, &b)
	if res != 0 {
		return int(l), int(t), int(r), int(b)
	}
	return 0, 0, 0, 0
}

func tessPageIteratorOrientation(handle TessPageIterator) TessPageOrientation {
	orientation := C.TessOrientation(0)
	writing_direction := C.TessWritingDirection(0)
	textline_order := C.TessTextlineOrder(0)
	deskew_angle := C.float(0)

	C.TessPageIteratorOrientation(handle.it, &orientation, &writing_direction, &textline_order, &deskew_angle)

	return TessPageOrientation(orientation)
}
