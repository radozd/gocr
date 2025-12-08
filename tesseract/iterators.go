package tesseract

// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

func (resIt TessResultIterator) AsPageIterator() TessPageIterator {
	return tessResultIteratorGetPageIterator(resIt)
}

func (resIt TessResultIterator) Copy() TessResultIterator {
	return tessResultIteratorCopy(resIt)
}

func (resIt TessResultIterator) Delete() {
	tessResultIteratorDelete(resIt)
}

func (resIt TessResultIterator) Next(level TessPageIteratorLevel) bool {
	return tessResultIteratorNext(resIt, level)
}

func (resIt TessResultIterator) GetUTF8Text(level TessPageIteratorLevel) (string, float32) {
	text := tessResultIteratorGetUTF8Text(resIt, level)
	if text == nil {
		return "", 0.0
	}

	res := C.GoString((*C.char)(unsafe.Pointer(text)))
	tessDeleteText(text)

	confidence := tessResultIteratorConfidence(resIt, level)

	return res, confidence
}

// /////////////////////////////////////
func (pageIt TessPageIterator) Begin() {
	tessPageIteratorBegin(pageIt)
}

func (pageIt TessPageIterator) IsAtBeginningOf(level TessPageIteratorLevel) bool {
	return tessPageIteratorIsAtBeginningOf(pageIt, level)
}

func (pageIt TessPageIterator) IsAtFinalElement(level TessPageIteratorLevel, element TessPageIteratorLevel) bool {
	return tessPageIteratorIsAtFinalElement(pageIt, level, element)
}

func (pageIt TessPageIterator) Next(level TessPageIteratorLevel) bool {
	return tessPageIteratorNext(pageIt, level)
}

func (pageIt TessPageIterator) BoundingBox(level TessPageIteratorLevel) (int, int, int, int) {
	return tessPageIteratorBoundingBox(pageIt, level)
}
