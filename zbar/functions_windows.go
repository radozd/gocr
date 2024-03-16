package zbar

// #include <stdlib.h>
import "C"

import (
	"syscall"
	"unsafe"
)

var (
	zbarDll = syscall.NewLazyDLL("libzbar64.dll")

	// init
	_zbar_image_create     = zbarDll.NewProc("zbar_image_create")
	_zbar_image_destroy    = zbarDll.NewProc("zbar_image_destroy")
	_zbar_image_set_size   = zbarDll.NewProc("zbar_image_set_size")   // (width, height)
	_zbar_image_set_format = zbarDll.NewProc("zbar_image_set_format") // always "Y800"
	_zbar_image_set_data   = zbarDll.NewProc("zbar_image_set_data")   // (data, length)

	// scanner
	_zbar_image_scanner_create     = zbarDll.NewProc("zbar_image_scanner_create")
	_zbar_image_scanner_destroy    = zbarDll.NewProc("zbar_image_scanner_destroy")
	_zbar_image_scanner_set_config = zbarDll.NewProc("zbar_image_scanner_set_config")

	_zbar_scan_image = zbarDll.NewProc("zbar_scan_image")

	// symbol
	_zbar_image_first_symbol  = zbarDll.NewProc("zbar_image_first_symbol")
	_zbar_symbol_next         = zbarDll.NewProc("zbar_symbol_next")
	_zbar_symbol_get_type     = zbarDll.NewProc("zbar_symbol_get_type")
	_zbar_symbol_get_data     = zbarDll.NewProc("zbar_symbol_get_data")
	_zbar_symbol_get_loc_size = zbarDll.NewProc("zbar_symbol_get_loc_size")
	_zbar_symbol_get_loc_x    = zbarDll.NewProc("zbar_symbol_get_loc_x")
	_zbar_symbol_get_loc_y    = zbarDll.NewProc("zbar_symbol_get_loc_y")
)

type Scanner struct {
	handle uintptr
}

type image uintptr

type symbol uintptr

var nullSymbol symbol = symbol(0)

func zbar_image_create(width int, height int, data uintptr, size int) image {
	img, _, _ := _zbar_image_create.Call()
	_zbar_image_set_size.Call(img, uintptr(width), uintptr(height))
	_zbar_image_set_format.Call(img, uintptr(C.ulong(0x30303859))) // unsafe.Pointer(cFormat) // Y800 (grayscale)
	_zbar_image_set_data.Call(img, data, uintptr(size), uintptr(0))

	return image(img)
}

func zbar_image_destroy(img image) {
	_zbar_image_destroy.Call(uintptr(img))
}

// scanner
func zbar_image_scanner_create() Scanner {
	s, _, _ := _zbar_image_scanner_create.Call()
	return Scanner{handle: s}
}

func zbar_image_scanner_destroy(scn *Scanner) {
	_zbar_image_scanner_destroy.Call(scn.handle)
	scn.handle = 0
}

func zbar_image_scanner_set_config(scn Scanner, symbology int, config int, value int) {
	_zbar_image_scanner_set_config.Call(scn.handle, uintptr(symbology), uintptr(config), uintptr(value))
}

func zbar_scan_image(scn Scanner, img image) bool {
	res, _, _ := _zbar_scan_image.Call(scn.handle, uintptr(img))
	return res > 0
}

// symbol
func zbar_image_first_symbol(img image) symbol {
	sym, _, _ := _zbar_image_first_symbol.Call(uintptr(img))
	return symbol(sym)
}

func zbar_symbol_next(sym symbol) symbol {
	n, _, _ := _zbar_symbol_next.Call(uintptr(sym))
	return symbol(n)
}

func zbar_symbol_get_data(sym symbol) string {
	text, _, _ := _zbar_symbol_get_data.Call(uintptr(sym))
	if text == 0 {
		return ""
	}
	return C.GoString((*C.char)(unsafe.Pointer(text)))
}

func zbar_symbol_get_loc_size(sym symbol) int {
	res, _, _ := _zbar_symbol_get_loc_size.Call(uintptr(sym))
	return int(res)
}

func zbar_symbol_get_loc_x(sym symbol, idx int) int {
	res, _, _ := _zbar_symbol_get_loc_x.Call(uintptr(sym), uintptr(idx))
	return int(res)
}

func zbar_symbol_get_loc_y(sym symbol, idx int) int {
	res, _, _ := _zbar_symbol_get_loc_y.Call(uintptr(sym), uintptr(idx))
	return int(res)
}

func zbar_symbol_get_type(sym symbol) ZBAR_CODETYPE {
	t, _, _ := _zbar_symbol_get_type.Call(uintptr(sym))
	return ZBAR_CODETYPE(t)
}
