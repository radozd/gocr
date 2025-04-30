package tesseract

// #include <stdlib.h>
import "C"

import (
	"strings"
	"unsafe"

	"github.com/radozd/gocr/leptonica"
)

var DataPath string = "."

func NewApi(lang string) Api {
	api := tessBaseAPICreate()
	if api.handle == NullApi.handle {
		return NullApi
	}

	cDatapath := C.CString(DataPath)
	defer C.free(unsafe.Pointer(cDatapath))

	cLang := C.CString(lang)
	defer C.free(unsafe.Pointer(cLang))

	tessBaseAPIInit2(api, DataPath, lang, OEM_DEFAULT)

	return api
}

func (api Api) Close() {
	tessBaseAPIEnd(api)
	tessBaseAPIDelete(api)
	api.handle = NullApi.handle
}

func (api Api) SetVariable(name string, value string) {
	tessBaseAPISetVariable(api, name, value)
}

func (api Api) SetImagePix(pix leptonica.Pix) {
	tessBaseAPISetImage2(api, pix)
}

func (api Api) SetPageSegMode(mode TessPageSegMode) {
	tessBaseAPISetPageSegMode(api, mode)
}

func (api Api) GetIterator() TessResultIterator {
	return tessBaseAPIGetIterator(api)
}

// /////////////////////////////////////
func (api Api) Recognize() {
	tessBaseAPIRecognize(api)
}

func (api Api) Text() string {
	text := tessBaseAPIGetUTF8Text(api)
	if text == nil {
		return ""
	}

	res := C.GoString((*C.char)(unsafe.Pointer(text)))
	tessDeleteText(text)

	return res
}

// HOCRText returns the HOCR text for given pagenumber
func (api Api) HOCRText(pagenumber int) string {
	text := tessBaseAPIGetHOCRText(api, pagenumber)
	if text == nil {
		return ""
	}

	res := C.GoString((*C.char)(unsafe.Pointer(text)))
	tessDeleteText(text)

	return res
}

func (api Api) TextBlocks(level TessPageIteratorLevel) []TextBlock {
	resIt := api.GetIterator()
	defer resIt.Delete()

	pageIt := resIt.AsPageIterator()
	blocks := make([]TextBlock, 0)

	dynamic := level == RIL_DYNAMIC
	if dynamic {
		level = RIL_BLOCK
	}

	good := true
	for good {
		if text, goodness := resIt.GetUTF8Text(level); strings.TrimSpace(text) != "" {
			x1, y1, x2, y2 := pageIt.BoundingBox(level)
			if !dynamic {
				blocks = append(blocks, TextBlock{
					Goodness: goodness,
					Value:    text,
					X:        x1,
					Y:        y1,
					Width:    x2 - x1,
					Height:   y2 - y1,
				})
			} else {
				if sparsity := float32((x2-x1)*(y2-y1)) / float32(len(text)) / 1000; sparsity < 4 {
					blocks = append(blocks, TextBlock{
						Goodness: goodness,
						Value:    text,
						X:        x1,
						Y:        y1,
						Width:    x2 - x1,
						Height:   y2 - y1,
					})
				} else {
					it2 := resIt.Copy()
					pit2 := it2.AsPageIterator()
					good2 := true
					for good2 {
						if text, goodness := it2.GetUTF8Text(RIL_TEXTLINE); strings.TrimSpace(text) != "" {
							x1, y1, x2, y2 = pit2.BoundingBox(RIL_TEXTLINE)
							blocks = append(blocks, TextBlock{
								Goodness: goodness,
								Value:    text,
								X:        x1,
								Y:        y1,
								Width:    x2 - x1,
								Height:   y2 - y1,
							})
						}
						good2 = pit2.Next(RIL_TEXTLINE)
					}
					it2.Delete()
				}
			}
		}

		good = pageIt.Next(level)
	}
	return blocks
}

func (api Api) GetPageOrientation() TessPageOrientation {
	resIt := api.GetIterator()
	defer resIt.Delete()

	return tessPageIteratorOrientation(resIt.AsPageIterator())
}
