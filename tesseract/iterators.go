package tesseract

// #include <stdlib.h>
import "C"

import (
	"math"
	"unsafe"
)

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
