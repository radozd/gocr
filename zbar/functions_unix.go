//go:build darwin || linux
// +build darwin linux

package zbar

// #include <stdlib.h>
// #cgo LDFLAGS: libzbar64.a /Library/Developer/CommandLineTools/SDKs/MacOSX15.sdk/usr/lib/libiconv.tbd
// #include <zbar.h>
import "C"
import "unsafe"

type Scanner struct {
	p *C.zbar_image_scanner_t
}

type image struct {
	p *C.zbar_image_t
}

type symbol struct {
	p *C.zbar_symbol_t
}

var nullSymbol symbol = symbol{}

func zbar_image_create(width int, height int, data uintptr, size int) image {
	img := C.zbar_image_create()
	C.zbar_image_set_size(img, C.uint(width), C.uint(height))
	C.zbar_image_set_format(img, C.ulong(0x30303859)) // unsafe.Pointer(cFormat) // Y800 (grayscale)
	C.zbar_image_set_data(img, unsafe.Pointer(data), C.ulong(size), nil)

	return image{p: img}
}

func zbar_image_destroy(img image) {
	C.zbar_image_destroy(img.p)
}

// scanner
func zbar_image_scanner_create() Scanner {
	s := C.zbar_image_scanner_create()
	return Scanner{p: s}
}

func zbar_image_scanner_destroy(scn *Scanner) {
	C.zbar_image_scanner_destroy(scn.p)
	scn.p = nil
}

func zbar_image_scanner_set_config(scn Scanner, symbology int, config int, value int) {
	C.zbar_image_scanner_set_config(scn.p, C.zbar_symbol_type_t(symbology), C.zbar_config_t(config), C.int(value))
}

func zbar_scan_image(scn Scanner, img image) bool {
	res := C.zbar_scan_image(scn.p, img.p)
	return res > 0
}

// symbol
func zbar_image_first_symbol(img image) symbol {
	sym := C.zbar_image_first_symbol(img.p)
	return symbol{p: sym}
}

func zbar_symbol_next(sym symbol) symbol {
	n := C.zbar_symbol_next(sym.p)
	return symbol{p: n}
}

func zbar_symbol_get_data(sym symbol) string {
	text := C.zbar_symbol_get_data(sym.p)
	if text == nil {
		return ""
	}
	return C.GoString((*C.char)(unsafe.Pointer(text)))
}

func zbar_symbol_get_loc_size(sym symbol) int {
	res := C.zbar_symbol_get_loc_size(sym.p)
	return int(res)
}

func zbar_symbol_get_loc_x(sym symbol, idx int) int {
	res := C.zbar_symbol_get_loc_x(sym.p, C.uint(idx))
	return int(res)
}

func zbar_symbol_get_loc_y(sym symbol, idx int) int {
	res := C.zbar_symbol_get_loc_y(sym.p, C.uint(idx))
	return int(res)
}

func zbar_symbol_get_type(sym symbol) ZBAR_CODETYPE {
	t := C.zbar_symbol_get_type(sym.p)
	return ZBAR_CODETYPE(t)
}
