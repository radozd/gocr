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

	_pixRead          = leptonicaDll.NewProc("pixRead")
	_pixReadMem       = leptonicaDll.NewProc("pixReadMem")
	_pixWrite         = leptonicaDll.NewProc("pixWrite")
	_pixWriteMem      = leptonicaDll.NewProc("pixWriteMem")
	_pixGetData       = leptonicaDll.NewProc("pixGetData")
	_pixGetDimensions = leptonicaDll.NewProc("pixGetDimensions")
	_pixGetWpl        = leptonicaDll.NewProc("pixGetWpl")
	_pixCountPixels   = leptonicaDll.NewProc("pixCountPixels")

	_pixRotate180         = leptonicaDll.NewProc("pixRotate180")
	_pixRotate            = leptonicaDll.NewProc("pixRotate")
	_pixAddBorder         = leptonicaDll.NewProc("pixAddBorder")
	_pixScaleToSize       = leptonicaDll.NewProc("pixScaleToSize")
	_pixDeskew            = leptonicaDll.NewProc("pixDeskew")
	_pixFindSkewAndDeskew = leptonicaDll.NewProc("pixFindSkewAndDeskew")

	_pixMaskOverColorPixels = leptonicaDll.NewProc("pixMaskOverColorPixels")
	_pixMaskOverGrayPixels  = leptonicaDll.NewProc("pixMaskOverGrayPixels")
	_pixSetMasked           = leptonicaDll.NewProc("pixSetMasked")
	_pixCombineMasked       = leptonicaDll.NewProc("pixCombineMasked")

	_pixContrastTRC      = leptonicaDll.NewProc("pixContrastTRC")
	_pixGammaTRC         = leptonicaDll.NewProc("pixGammaTRC")
	_pixBackgroundNorm   = leptonicaDll.NewProc("pixBackgroundNorm")
	_pixThresholdToValue = leptonicaDll.NewProc("pixThresholdToValue")

	_pixConvertRGBToGrayFast = leptonicaDll.NewProc("pixConvertRGBToGrayFast")
	_pixConvertTo8           = leptonicaDll.NewProc("pixConvertTo8")
	_pixConvertTo1           = leptonicaDll.NewProc("pixConvertTo1")

	_pixCopy        = leptonicaDll.NewProc("pixCopy")
	_pixInvert      = leptonicaDll.NewProc("pixInvert")
	_pixRasterop    = leptonicaDll.NewProc("pixRasterop")
	_pixOr          = leptonicaDll.NewProc("pixOr")
	_pixXor         = leptonicaDll.NewProc("pixXor")
	_pixAnd         = leptonicaDll.NewProc("pixAnd")
	_pixBlendInRect = leptonicaDll.NewProc("pixBlendInRect")

	_pixOpenBrick       = leptonicaDll.NewProc("pixOpenBrick")
	_pixCloseBrick      = leptonicaDll.NewProc("pixCloseBrick")
	_pixErodeBrick      = leptonicaDll.NewProc("pixErodeBrick")
	_pixDilateBrick     = leptonicaDll.NewProc("pixDilateBrick")
	_pixSobelEdgeFilter = leptonicaDll.NewProc("pixSobelEdgeFilter")

	_boxCreate      = leptonicaDll.NewProc("boxCreate")
	_boxDestroy     = leptonicaDll.NewProc("boxDestroy")
	_boxGetGeometry = leptonicaDll.NewProc("boxGetGeometry")
	_boxSetGeometry = leptonicaDll.NewProc("boxSetGeometry")

	_boxaDestroy  = leptonicaDll.NewProc("boxaDestroy")
	_boxaGetCount = leptonicaDll.NewProc("boxaGetCount")
	_boxaGetBox   = leptonicaDll.NewProc("boxaGetBox")

	_pixConnCompBB            = leptonicaDll.NewProc("pixConnCompBB")
	_pixRemoveBorderConnComps = leptonicaDll.NewProc("pixRemoveBorderConnComps")
)

type Pix struct {
	p unsafe.Pointer
}

var NullPix Pix = Pix{p: nil}

func UnsafePix(pix Pix) uintptr {
	return uintptr(pix.p)
}

func setMsgSeverity(level int) {
	_setMsgSeverity.Call(uintptr(C.int32_t(level)))
}

func pixDestroy(pix *Pix) {
	_pixDestroy.Call(uintptr(unsafe.Pointer(&pix.p)))
	pix.p = nil
}

func lept_free(cMem *C.uchar) {
	_lept_free.Call(uintptr(unsafe.Pointer(cMem)))
}

func pixRead(filename string) Pix {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	p, _, _ := _pixRead.Call(uintptr(unsafe.Pointer(cFilename)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixReadMem(image []byte) Pix {
	p, _, _ := _pixReadMem.Call(uintptr(unsafe.Pointer(unsafe.SliceData(image))), uintptr(C.size_t(len(image))))
	return Pix{p: unsafe.Pointer(p)}
}

func pixWrite(filename string, pix Pix, format int) bool {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	code, _, _ := _pixWrite.Call(uintptr(unsafe.Pointer(cFilename)), uintptr(pix.p), uintptr(C.int32_t(format)))
	return code == 0
}

func pixWriteMem(pix Pix, format int) []byte {
	var cMem *C.uchar = nil
	size := C.size_t(0)

	code, _, _ := _pixWriteMem.Call(uintptr(unsafe.Pointer(&cMem)), uintptr(unsafe.Pointer(&size)), uintptr(pix.p), uintptr(C.int32_t(format)))
	if code != 0 {
		return nil
	}
	defer lept_free(cMem)

	return C.GoBytes(unsafe.Pointer(cMem), C.int(size))
}

func pixGetData(pix Pix) *C.uint32_t {
	raw, _, _ := _pixGetData.Call(uintptr(pix.p))
	return (*C.uint32_t)(unsafe.Pointer(raw))
}

func pixGetDimensions(pix Pix) (int, int, int) {
	w := C.int(0)
	h := C.int(0)
	d := C.int(0)
	code, _, _ := _pixGetDimensions.Call(uintptr(pix.p), uintptr(unsafe.Pointer(&w)), uintptr(unsafe.Pointer(&h)), uintptr(unsafe.Pointer(&d)))
	if code != 0 {
		return 0, 0, 0
	}
	return int(w), int(h), int(d)
}

func pixGetWpl(pixs Pix) int {
	code, _, _ := _pixGetWpl.Call(uintptr(pixs.p))
	return int(code)
}

func pixCountPixels(pix Pix) int {
	count := C.int(0)
	code, _, _ := _pixCountPixels.Call(uintptr(pix.p), uintptr(unsafe.Pointer(&count)), uintptr(0))
	if code != 0 {
		return -1
	}
	return int(count)
}

////////////////////////////////////////////////////

func pixRotate180(pixd Pix, pixs Pix) Pix {
	p, _, _ := _pixRotate180.Call(uintptr(pixd.p), uintptr(pixs.p))
	return Pix{p: unsafe.Pointer(p)}
}

func pixRotate(pixs Pix, angle float32) Pix {
	deg2rad := float32(3.1415926535 / 180.0)
	p, _, _ := _pixRotate.Call(uintptr(pixs.p), uintptr(math.Float32bits(deg2rad*angle)), uintptr(1), uintptr(1), uintptr(0), uintptr(0))
	return Pix{p: unsafe.Pointer(p)}
}

func pixAddBorder(pixs Pix, npix int, val uint) Pix {
	p, _, _ := _pixAddBorder.Call(uintptr(pixs.p), uintptr(C.int32_t(npix)), uintptr(C.uint32_t(val)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixScaleToSize(pixs Pix, wd int, hd int) Pix {
	p, _, _ := _pixScaleToSize.Call(uintptr(pixs.p), uintptr(C.int32_t(wd)), uintptr(C.int32_t(hd)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixDeskew(pixs Pix, redsearch int) Pix {
	p, _, _ := _pixDeskew.Call(uintptr(pixs.p), uintptr(C.int32_t(redsearch)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixFindSkewAndDeskew(pixs Pix, redsearch int) (Pix, float32) {
	angle := C.float(0)
	conf := C.float(0)
	p, _, _ := _pixFindSkewAndDeskew.Call(uintptr(pixs.p), uintptr(C.int32_t(redsearch)), uintptr(unsafe.Pointer(&angle)), uintptr(unsafe.Pointer(&conf)))
	return Pix{p: unsafe.Pointer(p)}, float32(angle)
}

////////////////////////////////////////////////////

func pixMaskOverColorPixels(pixs Pix, threshdiff int, mindist int) Pix {
	p, _, _ := _pixMaskOverColorPixels.Call(uintptr(pixs.p), uintptr(C.int32_t(threshdiff)), uintptr(C.int32_t(mindist)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixMaskOverGrayPixels(pixs Pix, maxlimit int, satlimit int) Pix {
	p, _, _ := _pixMaskOverGrayPixels.Call(uintptr(pixs.p), uintptr(C.int32_t(maxlimit)), uintptr(C.int32_t(satlimit)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixSetMasked(pixd Pix, pixm Pix, val uint) int {
	code, _, _ := _pixSetMasked.Call(uintptr(pixd.p), uintptr(pixm.p), uintptr(C.uint32_t(val)))
	return int(code)
}

func pixCombineMasked(pixd Pix, pixs Pix, pixm Pix) int {
	code, _, _ := _pixCombineMasked.Call(uintptr(pixd.p), uintptr(pixs.p), uintptr(pixm.p))
	return int(code)
}

////////////////////////////////////////////////////

func pixContrastTRC(pixd Pix, pixs Pix, factor float32) Pix {
	p, _, _ := _pixContrastTRC.Call(uintptr(pixd.p), uintptr(pixs.p), uintptr(math.Float32bits(factor)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixGammaTRC(pixd Pix, pixs Pix, gamma float32, minval int, maxval int) Pix {
	p, _, _ := _pixGammaTRC.Call(uintptr(pixd.p), uintptr(pixs.p), uintptr(math.Float32bits(gamma)), uintptr(C.int32_t(minval)), uintptr(C.int32_t(maxval)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixBackgroundNorm(pixs Pix, pixim Pix, pixg Pix, sx int, sy int, thresh int, mincount int, bgval int, smoothx int, smoothy int) Pix {
	p, _, _ := _pixBackgroundNorm.Call(uintptr(pixs.p), uintptr(pixim.p), uintptr(pixg.p), uintptr(C.int32_t(sx)), uintptr(C.int32_t(sy)),
		uintptr(C.int32_t(thresh)), uintptr(C.int32_t(mincount)), uintptr(C.int32_t(bgval)), uintptr(C.int32_t(smoothx)), uintptr(C.int32_t(smoothy)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixThresholdToValue(pixd Pix, pixs Pix, threshval int, setval int) Pix {
	p, _, _ := _pixThresholdToValue.Call(uintptr(pixd.p), uintptr(pixs.p), uintptr(C.int32_t(threshval)), uintptr(C.int32_t(setval)))
	return Pix{p: unsafe.Pointer(p)}
}

////////////////////////////////////////////////////

func pixConvertRGBToGrayFast(pixs Pix) Pix {
	p, _, _ := _pixConvertRGBToGrayFast.Call(uintptr(pixs.p))
	return Pix{p: unsafe.Pointer(p)}
}

func pixConvertTo8(pixs Pix, cmapflag int) Pix {
	p, _, _ := _pixConvertTo8.Call(uintptr(pixs.p), uintptr(C.int32_t(cmapflag)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixConvertTo1(pixs Pix, threshold int) Pix {
	p, _, _ := _pixConvertTo1.Call(uintptr(pixs.p), uintptr(C.int32_t(threshold)))
	return Pix{p: unsafe.Pointer(p)}
}

////////////////////////////////////////////////////

func pixCopy(pixd Pix, pixs Pix) Pix {
	p, _, _ := _pixCopy.Call(uintptr(pixd.p), uintptr(pixs.p))
	return Pix{p: unsafe.Pointer(p)}
}

func pixInvert(pixd Pix, pixs Pix) Pix {
	p, _, _ := _pixInvert.Call(uintptr(pixd.p), uintptr(pixs.p))
	return Pix{p: unsafe.Pointer(p)}
}

func pixRasterop(pixd Pix, dx int, dy int, dw int, dh int, op int, pixs Pix, sx int, sy int) bool {
	res, _, _ := _pixRasterop.Call(uintptr(pixd.p), uintptr(C.int32_t(dx)), uintptr(C.int32_t(dy)), uintptr(C.int32_t(dw)), uintptr(C.int32_t(dh)),
		uintptr(C.int32_t(op)), uintptr(pixs.p), uintptr(C.int32_t(sx)), uintptr(C.int32_t(sy)))
	return res == 0
}

func pixOr(pixd Pix, pixs1 Pix, pixs2 Pix) {
	_pixOr.Call(uintptr(pixd.p), uintptr(pixs1.p), uintptr(pixs2.p))
}

func pixXor(pixd Pix, pixs1 Pix, pixs2 Pix) {
	_pixXor.Call(uintptr(pixd.p), uintptr(pixs1.p), uintptr(pixs2.p))
}

func pixAnd(pixd Pix, pixs1 Pix, pixs2 Pix) {
	_pixAnd.Call(uintptr(pixd.p), uintptr(pixs1.p), uintptr(pixs2.p))
}

func pixBlendInRect(pixs Pix, box Box, val int, fract float32) {
	_pixBlendInRect.Call(uintptr(pixs.p), uintptr(box.p), uintptr(val), uintptr(math.Float32bits(fract)))
}

////////////////////////////////////////////////////

func pixOpenBrick(pixd Pix, pixs Pix, hsize int, vsize int) Pix {
	p, _, _ := _pixOpenBrick.Call(uintptr(pixd.p), uintptr(pixs.p), uintptr(C.int32_t(hsize)), uintptr(C.int32_t(vsize)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixCloseBrick(pixd Pix, pixs Pix, hsize int, vsize int) Pix {
	p, _, _ := _pixCloseBrick.Call(uintptr(pixd.p), uintptr(pixs.p), uintptr(C.int32_t(hsize)), uintptr(C.int32_t(vsize)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixErodeBrick(pixd Pix, pixs Pix, hsize int, vsize int) Pix {
	p, _, _ := _pixErodeBrick.Call(uintptr(pixd.p), uintptr(pixs.p), uintptr(C.int32_t(hsize)), uintptr(C.int32_t(vsize)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixDilateBrick(pixd Pix, pixs Pix, hsize int, vsize int) Pix {
	p, _, _ := _pixDilateBrick.Call(uintptr(pixd.p), uintptr(pixs.p), uintptr(C.int32_t(hsize)), uintptr(C.int32_t(vsize)))
	return Pix{p: unsafe.Pointer(p)}
}

func pixSobelEdgeFilter(pixs Pix, orientflag int) Pix {
	p, _, _ := _pixSobelEdgeFilter.Call(uintptr(pixs.p), uintptr(C.int32_t(orientflag)))
	return Pix{p: unsafe.Pointer(p)}
}

////////////////////////////////////////////////////

type Box struct {
	p unsafe.Pointer
}

func boxCreate(x, y, w, h int) Box {
	p, _, _ := _boxCreate.Call(uintptr(x), uintptr(y), uintptr(w), uintptr(h))
	return Box{p: unsafe.Pointer(p)}
}

func boxDestroy(box *Box) {
	_boxDestroy.Call(uintptr(unsafe.Pointer(&box.p)))
	box.p = nil
}

func boxGetGeometry(box Box) (int, int, int, int) {
	x := C.int(0)
	y := C.int(0)
	w := C.int(0)
	h := C.int(0)
	code, _, _ := _boxGetGeometry.Call(uintptr(box.p), uintptr(unsafe.Pointer(&x)), uintptr(unsafe.Pointer(&y)), uintptr(unsafe.Pointer(&w)), uintptr(unsafe.Pointer(&h)))
	if code != 0 {
		return 0, 0, 0, 0
	}
	return int(x), int(y), int(w), int(h)
}

func boxSetGeometry(box Box, x int, y int, w int, h int) int {
	code, _, _ := _boxSetGeometry.Call(uintptr(box.p), uintptr(C.int32_t(x)), uintptr(C.int32_t(y)), uintptr(C.int32_t(w)), uintptr(C.int32_t(h)))
	return int(code)
}

////////////////////////////////////////////////////

type Boxa struct {
	p unsafe.Pointer
}

func boxaDestroy(boxa *Boxa) {
	_boxaDestroy.Call(uintptr(unsafe.Pointer(&boxa.p)))
	boxa.p = nil
}

func boxaGetCount(boxa Boxa) int {
	code, _, _ := _boxaGetCount.Call(uintptr(boxa.p))
	return int(code)
}

func boxaGetBox(boxa Boxa, index int, accessflag int) Box {
	p, _, _ := _boxaGetBox.Call(uintptr(boxa.p), uintptr(C.int32_t(index)), uintptr(C.int32_t(accessflag)))
	return Box{p: unsafe.Pointer(p)}
}

////////////////////////////////////////////////////

func pixConnCompBB(pixs Pix, connectivity int) Boxa {
	p, _, _ := _pixConnCompBB.Call(uintptr(pixs.p), uintptr(C.int32_t(connectivity)))
	return Boxa{p: unsafe.Pointer(p)}
}

func pixRemoveBorderConnComps(pixs Pix, connectivity int) Pix {
	p, _, _ := _pixRemoveBorderConnComps.Call(uintptr(pixs.p), uintptr(C.int32_t(connectivity)))
	return Pix{p: unsafe.Pointer(p)}
}
