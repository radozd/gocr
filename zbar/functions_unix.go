//go:build darwin || linux
// +build darwin linux

package zbar

// #include <stdlib.h>
// #cgo LDFLAGS: libzbar64.a /Library/Developer/CommandLineTools/SDKs/MacOSX15.sdk/usr/lib/libiconv.tbd
// #include <zbar.h>
import "C"
import "unsafe"

type Scanner struct {
	handle *C.zbar_image_scanner_t
}

type image struct {
	handle *C.zbar_image_t
}

type symbol struct {
	handle *C.zbar_symbol_t
}

var nullSymbol symbol = symbol{}

func zbar_image_create(width int, height int, data uintptr, size int) image {
	img := C.zbar_image_create()
	C.zbar_image_set_size(img, C.uint(width), C.uint(height))
	C.zbar_image_set_format(img, C.ulong(0x30303859)) // unsafe.Pointer(cFormat) // Y800 (grayscale)
	C.zbar_image_set_data(img, unsafe.Pointer(data), C.ulong(size), nil)

	return image{handle: img}
}

func zbar_image_destroy(img image) {
	C.zbar_image_destroy(img.handle)
}

// scanner
func zbar_image_scanner_create() Scanner {
	s := C.zbar_image_scanner_create()
	return Scanner{handle: s}
}

func zbar_image_scanner_destroy(scn *Scanner) {
	C.zbar_image_scanner_destroy(scn.handle)
	scn.handle = nil
}

func zbar_image_scanner_set_config(scn Scanner, symbology int, config int, value int) {
	C.zbar_image_scanner_set_config(scn.handle, C.zbar_symbol_type_t(symbology), C.zbar_config_t(config), C.int(value))
}

func zbar_scan_image(scn Scanner, img image) bool {
	res := C.zbar_scan_image(scn.handle, img.handle)
	return res > 0
}

// symbol
func zbar_image_first_symbol(img image) symbol {
	sym := C.zbar_image_first_symbol(img.handle)
	return symbol{handle: sym}
}

func zbar_symbol_next(sym symbol) symbol {
	n := C.zbar_symbol_next(sym.handle)
	return symbol{handle: n}
}

func zbar_symbol_get_data(sym symbol) string {
	text := C.zbar_symbol_get_data(sym.handle)
	if text == nil {
		return ""
	}
	return C.GoString((*C.char)(unsafe.Pointer(text)))
}

func zbar_symbol_get_loc_size(sym symbol) int {
	res := C.zbar_symbol_get_loc_size(sym.handle)
	return int(res)
}

func zbar_symbol_get_loc_x(sym symbol, idx int) int {
	res := C.zbar_symbol_get_loc_x(sym.handle, C.uint(idx))
	return int(res)
}

func zbar_symbol_get_loc_y(sym symbol, idx int) int {
	res := C.zbar_symbol_get_loc_y(sym.handle, C.uint(idx))
	return int(res)
}

func zbar_symbol_get_type(sym symbol) ZBAR_CODETYPE {
	t := C.zbar_symbol_get_type(sym.handle)
	return ZBAR_CODETYPE(t)
}
