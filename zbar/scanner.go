package zbar

// #include <stdlib.h>
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/radozd/gocr/leptonica"
)

type Scanner struct {
	handle uintptr
}

const ZBAR_CFG_ENABLE int = 0

func NewScanner() *Scanner {
	s, _, _ := zbar_image_scanner_create.Call()
	return &Scanner{handle: s}
}

func (scn Scanner) setConfig(symbology int, config int, value int) {
	zbar_image_scanner_set_config.Call(scn.handle, uintptr(symbology), uintptr(config), uintptr(value))
}

func (scn Scanner) scan(img image) bool {
	res, _, _ := zbar_scan_image.Call(scn.handle, uintptr(img))
	return res > 0
}

func (scn *Scanner) Destroy() {
	zbar_image_scanner_destroy.Call(scn.handle)
	scn.handle = 0
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
