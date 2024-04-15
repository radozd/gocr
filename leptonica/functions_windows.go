package leptonica

// #include <stdlib.h>
// #include <stdint.h>
import "C"

import (
	"math"
	"syscall"
	"unsafe"
)

var (
	leptonicaDll = syscall.NewLazyDLL("leptonica-1.84.0.dll")

	_setMsgSeverity = leptonicaDll.NewProc("setMsgSeverity")
	_pixDestroy     = leptonicaDll.NewProc("pixDestroy")
	_lept_free      = leptonicaDll.NewProc("lept_free")

	_pixRead     = leptonicaDll.NewProc("pixRead")
	_pixReadMem  = leptonicaDll.NewProc("pixReadMem")
	_pixWrite    = leptonicaDll.NewProc("pixWrite")
	_pixWriteMem = leptonicaDll.NewProc("pixWriteMem")

	_pixRotate180         = leptonicaDll.NewProc("pixRotate180")
	_pixRotate            = leptonicaDll.NewProc("pixRotate")
	_pixAddBorder         = leptonicaDll.NewProc("pixAddBorder")
	_pixRasterop          = leptonicaDll.NewProc("pixRasterop")
	_pixScaleToSize       = leptonicaDll.NewProc("pixScaleToSize")
	_pixDeskew            = leptonicaDll.NewProc("pixDeskew")
	_pixFindSkewAndDeskew = leptonicaDll.NewProc("pixFindSkewAndDeskew")

	_pixMaskOverColorPixels = leptonicaDll.NewProc("pixMaskOverColorPixels")
	_pixMaskOverGrayPixels  = leptonicaDll.NewProc("pixMaskOverGrayPixels")
	_pixInvert              = leptonicaDll.NewProc("pixInvert")
	_pixPaintThroughMask    = leptonicaDll.NewProc("pixPaintThroughMask")

	_pixContrastTRC    = leptonicaDll.NewProc("pixContrastTRC")
	_pixGammaTRC       = leptonicaDll.NewProc("pixGammaTRC")
	_pixBackgroundNorm = leptonicaDll.NewProc("pixBackgroundNorm")

	_pixConvertRGBToGrayFast = leptonicaDll.NewProc("pixConvertRGBToGrayFast")
	_pixConvertTo8           = leptonicaDll.NewProc("pixConvertTo8")
	_pixCopy                 = leptonicaDll.NewProc("pixCopy")

	_pixGetData       = leptonicaDll.NewProc("pixGetData")
	_pixGetDimensions = leptonicaDll.NewProc("pixGetDimensions")
	_pixGetWpl        = leptonicaDll.NewProc("pixGetWpl")
)

type Pix uintptr

const NullPix Pix = 0

func UnsafePix(pix Pix) uintptr {
	return uintptr(pix)
}

func setMsgSeverity(level int) {
	_setMsgSeverity.Call(uintptr(C.int32_t(level)))
}

func pixDestroy(pix *Pix) {
	_pixDestroy.Call(uintptr(unsafe.Pointer(pix)))
	*pix = NullPix
}

func lept_free(cMem *C.uchar) {
	_lept_free.Call(uintptr(unsafe.Pointer(cMem)))
}

func pixRead(filename string) Pix {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	pix, _, _ := _pixRead.Call(uintptr(unsafe.Pointer(cFilename)))
	return Pix(pix)
}

func pixReadMem(image []byte) Pix {
	pix, _, _ := _pixReadMem.Call(uintptr(unsafe.Pointer(unsafe.SliceData(image))), uintptr(C.size_t(len(image))))
	return Pix(pix)
}

func pixWrite(filename string, pix Pix, format int) bool {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	code, _, _ := _pixWrite.Call(uintptr(unsafe.Pointer(cFilename)), uintptr(pix), uintptr(C.int32_t(format)))
	return code == 0
}

func pixWriteMem(pix Pix, format int) []byte {
	var cMem *C.uchar = nil
	size := C.size_t(0)

	code, _, _ := _pixWriteMem.Call(uintptr(unsafe.Pointer(&cMem)), uintptr(unsafe.Pointer(&size)), uintptr(pix), uintptr(C.int32_t(format)))
	if code != 0 {
		return nil
	}
	defer lept_free(cMem)

	return C.GoBytes(unsafe.Pointer(cMem), C.int(size))
}

func pixRotate180(pixd Pix, pixs Pix) Pix {
	pix, _, _ := _pixRotate180.Call(uintptr(pixd), uintptr(pixs))
	return Pix(pix)
}

func pixRotate(pixs Pix, angle float32) Pix {
	deg2rad := float32(3.1415926535 / 180.0)
	pix, _, _ := _pixRotate.Call(uintptr(pixs), uintptr(math.Float32bits(deg2rad*angle)), uintptr(1), uintptr(1), uintptr(0), uintptr(0))
	return Pix(pix)
}

func pixAddBorder(pixs Pix, npix int, val uint) Pix {
	pix, _, _ := _pixAddBorder.Call(uintptr(pixs), uintptr(C.int32_t(npix)), uintptr(C.uint32_t(val)))
	return Pix(pix)
}

func pixRasterop(pixd Pix, dx int, dy int, dw int, dh int, op int, pixs Pix, sx int, sy int) bool {
	res, _, _ := _pixRasterop.Call(uintptr(pixd), uintptr(C.int32_t(dx)), uintptr(C.int32_t(dy)), uintptr(C.int32_t(dw)), uintptr(C.int32_t(dh)),
		uintptr(C.int32_t(op)), uintptr(pixs), uintptr(C.int32_t(sx)), uintptr(C.int32_t(sy)))
	return res == 0
}

func pixScaleToSize(pixs Pix, wd int, hd int) Pix {
	pix, _, _ := _pixScaleToSize.Call(uintptr(pixs), uintptr(C.int32_t(wd)), uintptr(C.int32_t(hd)))
	return Pix(pix)
}

func pixDeskew(pixs Pix, redsearch int) Pix {
	pix, _, _ := _pixDeskew.Call(uintptr(pixs), uintptr(C.int32_t(redsearch)))
	return Pix(pix)
}

func pixFindSkewAndDeskew(pixs Pix, redsearch int) (Pix, float32) {
	angle := C.float(0)
	conf := C.float(0)
	pix, _, _ := _pixFindSkewAndDeskew.Call(uintptr(pixs), uintptr(C.int32_t(redsearch)), uintptr(unsafe.Pointer(&angle)), uintptr(unsafe.Pointer(&conf)))
	return Pix(pix), float32(angle)
}

func pixMaskOverColorPixels(pixs Pix, threshdiff int, mindist int) Pix {
	pix, _, _ := _pixMaskOverColorPixels.Call(uintptr(pixs), uintptr(C.int32_t(threshdiff)), uintptr(C.int32_t(mindist)))
	return Pix(pix)
}

func pixMaskOverGrayPixels(pixs Pix, maxlimit int, satlimit int) Pix {
	pix, _, _ := _pixMaskOverGrayPixels.Call(uintptr(pixs), uintptr(C.int32_t(maxlimit)), uintptr(C.int32_t(satlimit)))
	return Pix(pix)
}

func pixInvert(pixd Pix, pixs Pix) Pix {
	pix, _, _ := _pixInvert.Call(uintptr(pixd), uintptr(pixs))
	return Pix(pix)
}

func pixPaintThroughMask(pixd Pix, pixm Pix, x int, y int, val uint) int {
	code, _, _ := _pixPaintThroughMask.Call(uintptr(pixd), uintptr(pixm), uintptr(C.int32_t(x)), uintptr(C.int32_t(y)), uintptr(C.uint32_t(val)))
	return int(code)
}

func pixContrastTRC(pixd Pix, pixs Pix, factor float32) Pix {
	pix, _, _ := _pixContrastTRC.Call(uintptr(pixd), uintptr(pixs), uintptr(math.Float32bits(factor)))
	return Pix(pix)
}

func pixGammaTRC(pixd Pix, pixs Pix, gamma float32, minval int, maxval int) Pix {
	pix, _, _ := _pixGammaTRC.Call(uintptr(pixd), uintptr(pixs), uintptr(math.Float32bits(gamma)), uintptr(C.int32_t(minval)), uintptr(C.int32_t(maxval)))
	return Pix(pix)
}

func pixBackgroundNorm(pixs Pix, pixim Pix, pixg Pix, sx int, sy int, thresh int, mincount int, bgval int, smoothx int, smoothy int) Pix {
	pix, _, _ := _pixBackgroundNorm.Call(uintptr(pixs), uintptr(pixim), uintptr(pixg), uintptr(C.int32_t(sx)), uintptr(C.int32_t(sy)),
		uintptr(C.int32_t(thresh)), uintptr(C.int32_t(mincount)), uintptr(C.int32_t(bgval)), uintptr(C.int32_t(smoothx)), uintptr(C.int32_t(smoothy)))
	return Pix(pix)
}

func pixConvertRGBToGrayFast(pixs Pix) Pix {
	pix, _, _ := _pixConvertRGBToGrayFast.Call(uintptr(pixs))
	return Pix(pix)
}

func pixConvertTo8(pixs Pix, cmapflag int) Pix {
	pix, _, _ := _pixConvertTo8.Call(uintptr(pixs), uintptr(C.int32_t(cmapflag)))
	return Pix(pix)
}

func pixCopy(pixd Pix, pixs Pix) Pix {
	pix, _, _ := _pixCopy.Call(uintptr(pixd), uintptr(pixs))
	return Pix(pix)
}

func pixGetData(pix Pix) *C.uint32_t {
	raw, _, _ := _pixGetData.Call(uintptr(pix))
	return (*C.uint32_t)(unsafe.Pointer(raw))
}

func pixGetDimensions(pix Pix) (int, int, int) {
	w := C.int(0)
	h := C.int(0)
	d := C.int(0)
	code, _, _ := _pixGetDimensions.Call(uintptr(pix), uintptr(unsafe.Pointer(&w)), uintptr(unsafe.Pointer(&h)), uintptr(unsafe.Pointer(&d)))
	if code != 0 {
		return 0, 0, 0
	}
	return int(w), int(h), int(d)
}

func pixGetWpl(pixs Pix) int {
	code, _, _ := _pixGetWpl.Call(uintptr(pixs))
	return int(code)
}
