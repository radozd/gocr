package zbar

// #include <stdlib.h>
import "C"

type scanner uintptr

const ZBAR_CFG_ENABLE int = 0

func newScanner() scanner {
	s, _, _ := zbar_image_scanner_create.Call()
	scn := scanner(s)

	scn.setConfig(0, ZBAR_CFG_ENABLE, 1)
	return scn
}

func (scn scanner) setConfig(symbology int, config int, value int) {
	zbar_image_scanner_set_config.Call(uintptr(scn), uintptr(symbology), uintptr(config), uintptr(value))
}

func (scn scanner) scan(img image) bool {
	res, _, _ := zbar_scan_image.Call(uintptr(scn), uintptr(img))
	return res > 0
}

func (scn scanner) destroy() {
	zbar_image_scanner_destroy.Call(uintptr(scn))
}
