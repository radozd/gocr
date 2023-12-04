package zbar

import (
	"unsafe"

	"github.com/radozd/gocr/leptonica"
)

type ZBAR_CODETYPE int

const ZBAR_QRCODE ZBAR_CODETYPE = 64
const ZBAR_CODE128 ZBAR_CODETYPE = 128

type Code struct {
	CodeType ZBAR_CODETYPE
	Value    string
	X        int
	Y        int
	Width    int
	Height   int
}

func Process(pix leptonica.Pix) []Code {
	w, h, _ := pix.GetDimensions()
	gray := pix.GetRawGrayData()

	raw := unsafe.SliceData(gray)
	img := newImage(w, h, uintptr(unsafe.Pointer(raw)), len(gray))
	defer img.destroy()

	scanner := newScanner()
	defer scanner.destroy()

	codes := make([]Code, 0)
	if !scanner.scan(img) {
		return codes
	}

	img.first().each(func(code Code) {
		codes = append(codes, code)
	})
	return codes
}
