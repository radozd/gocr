package zbar

// #include <stdlib.h>
import "C"

import (
	"syscall"
)

var (
	zbarDll = syscall.NewLazyDLL("libzbar64.dll")

	// init
	zbar_image_create     = zbarDll.NewProc("zbar_image_create")
	zbar_image_destroy    = zbarDll.NewProc("zbar_image_destroy")
	zbar_image_set_size   = zbarDll.NewProc("zbar_image_set_size")   // (width, height)
	zbar_image_set_format = zbarDll.NewProc("zbar_image_set_format") // always "Y800"
	zbar_image_set_data   = zbarDll.NewProc("zbar_image_set_data")   // (data, length)

	// scanner
	zbar_image_scanner_create     = zbarDll.NewProc("zbar_image_scanner_create")
	zbar_image_scanner_destroy    = zbarDll.NewProc("zbar_image_scanner_destroy")
	zbar_image_scanner_set_config = zbarDll.NewProc("zbar_image_scanner_set_config")

	zbar_scan_image = zbarDll.NewProc("zbar_scan_image")

	// symbol
	zbar_image_first_symbol  = zbarDll.NewProc("zbar_image_first_symbol")
	zbar_symbol_next         = zbarDll.NewProc("zbar_symbol_next")
	zbar_symbol_get_type     = zbarDll.NewProc("zbar_symbol_get_type")
	zbar_symbol_get_data     = zbarDll.NewProc("zbar_symbol_get_data")
	zbar_symbol_get_loc_size = zbarDll.NewProc("zbar_symbol_get_loc_size")
	zbar_symbol_get_loc_x    = zbarDll.NewProc("zbar_symbol_get_loc_x")
	zbar_symbol_get_loc_y    = zbarDll.NewProc("zbar_symbol_get_loc_y")
)
