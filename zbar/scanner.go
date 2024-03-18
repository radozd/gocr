package zbar

import (
	"runtime"
	"unsafe"

	"github.com/radozd/gocr/leptonica"
)

const ZBAR_CFG_ENABLE int = 0

func NewScanner() Scanner {
	return zbar_image_scanner_create()
}

func (scn Scanner) setConfig(symbology int, config int, value int) {
	zbar_image_scanner_set_config(scn, symbology, config, value)
}

func (scn Scanner) scan(img image) bool {
	return zbar_scan_image(scn, img)
}

func (scn *Scanner) Destroy() {
	zbar_image_scanner_destroy(scn)
}

func (scn Scanner) Process(pix leptonica.Pix) []Code {
	codes := make([]Code, 0)

	w, h, _ := pix.GetDimensions()
	if w == 0 || h == 0 {
		return codes
	}

	gray := pix.GetRawGrayData()
	if len(gray) == 0 {
		return codes
	}

	raw := unsafe.Pointer(unsafe.SliceData(gray))
	img := newImage(w, h, uintptr(raw), len(gray))

	if scn.scan(img) {
		img.first().each(func(code Code) {
			codes = append(codes, code)
		})
	}

	img.destroy()
	runtime.KeepAlive(raw)

	return codes
}
