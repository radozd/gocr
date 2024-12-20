//go:build darwin || linux
// +build darwin linux

package leptonica

// #cgo CXXFLAGS: -std=c++0x
// #cgo CFLAGS: -I/opt/homebrew/include -I/usr/local/include
// #cgo CFLAGS: -Wno-unused-result
// #cgo LDFLAGS: -L/opt/homebrew/lib -L/usr/local/lib -lleptonica
// #include <stdlib.h>
// #include <leptonica/allheaders.h>
import "C"
import (
	"unsafe"
)

type Pix struct {
	p *C.PIX
}

var NullPix Pix = Pix{p: nil}

func UnsafePix(pix Pix) uintptr {
	return uintptr(unsafe.Pointer(pix.p))
}

func setMsgSeverity(level int) {
	C.setMsgSeverity(C.l_int32(level))
}

func pixDestroy(pix *Pix) {
	C.pixDestroy((**C.PIX)(unsafe.Pointer(pix)))
	pix.p = nil
}

func lept_free(cMem *C.uchar) {
	C.lept_free(unsafe.Pointer(cMem))
}

func pixRead(filename string) Pix {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	return Pix{p: C.pixRead(cFilename)}
}

func pixReadMem(image []byte) Pix {
	return Pix{p: C.pixReadMem((*C.uchar)(unsafe.Pointer(unsafe.SliceData(image))), C.size_t(len(image)))}
}

func pixWrite(filename string, pix Pix, format int) bool {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	return int(C.pixWrite(cFilename, pix.p, C.l_int32(format))) == 0
}

func pixWriteMem(pix Pix, format int) []byte {
	var cMem *C.uchar = nil
	size := C.size_t(0)

	code := C.pixWriteMem((**C.uchar)(unsafe.Pointer(&cMem)), &size, pix.p, C.l_int32(format))
	if code != 0 {
		return nil
	}
	defer lept_free(cMem)

	return C.GoBytes(unsafe.Pointer(cMem), C.int(size))
}

func pixGetData(pixs Pix) *C.l_uint32 {
	return (*C.l_uint32)(C.pixGetData(pixs.p))
}

func pixGetDimensions(pixs Pix) (int, int, int) {
	w := C.int(0)
	h := C.int(0)
	d := C.int(0)
	code := C.pixGetDimensions(pixs.p, &w, &h, &d)
	if code != 0 {
		return 0, 0, 0
	}
	return int(w), int(h), int(d)
}

func pixGetWpl(pixs Pix) int {
	return int(C.pixGetWpl(pixs.p))
}

////////////////////////////////////////////////////

func pixRotate180(pixd Pix, pixs Pix) Pix {
	return Pix{p: C.pixRotate180(pixd.p, pixs.p)}
}

func pixRotate(pixs Pix, angle float32) Pix {
	deg2rad := float32(3.1415926535 / 180.0)
	return Pix{p: C.pixRotate(pixs.p, C.float(deg2rad*angle), 1, 1, 0, 0)}
}

func pixAddBorder(pixs Pix, npix int, val uint) Pix {
	return Pix{p: C.pixAddBorder(pixs.p, C.l_int32(npix), C.l_uint32(val))}
}

func pixScaleToSize(pixs Pix, wd int, hd int) Pix {
	return Pix{p: C.pixScaleToSize(pixs.p, C.l_int32(wd), C.l_int32(hd))}
}

func pixDeskew(pixs Pix, redsearch int) Pix {
	return Pix{p: C.pixDeskew(pixs.p, C.l_int32(redsearch))}
}

func pixFindSkewAndDeskew(pixs Pix, redsearch int) (Pix, float32) {
	angle := C.float(0)
	conf := C.float(0)
	pix := C.pixFindSkewAndDeskew(pixs.p, C.l_int32(redsearch), &angle, &conf)
	return Pix{p: pix}, float32(angle)
}

////////////////////////////////////////////////////

func pixMaskOverColorPixels(pixs Pix, threshdiff int, mindist int) Pix {
	return Pix{p: C.pixMaskOverColorPixels(pixs.p, C.l_int32(threshdiff), C.l_int32(mindist))}
}

func pixMaskOverGrayPixels(pixs Pix, maxlimit int, satlimit int) Pix {
	return Pix{p: C.pixMaskOverGrayPixels(pixs.p, C.l_int32(maxlimit), C.l_int32(satlimit))}
}

func pixPaintThroughMask(pixd Pix, pixm Pix, x int, y int, val uint) int {
	return int(C.pixPaintThroughMask(pixd.p, pixm.p, C.l_int32(x), C.l_int32(y), C.l_uint32(val)))
}

////////////////////////////////////////////////////

func pixContrastTRC(pixd Pix, pixs Pix, factor float32) Pix {
	return Pix{p: C.pixContrastTRC(pixd.p, pixs.p, C.l_float32(factor))}
}

func pixGammaTRC(pixd Pix, pixs Pix, gamma float32, minval int, maxval int) Pix {
	return Pix{p: C.pixGammaTRC(pixd.p, pixs.p, C.l_float32(gamma), C.l_int32(minval), C.l_int32(maxval))}
}

func pixBackgroundNorm(pixs Pix, pixim Pix, pixg Pix, sx int, sy int, thresh int, mincount int, bgval int, smoothx int, smoothy int) Pix {
	return Pix{p: C.pixBackgroundNorm(pixs.p, pixim.p, pixg.p, C.l_int32(sx), C.l_int32(sy),
		C.l_int32(thresh), C.l_int32(mincount), C.l_int32(bgval), C.l_int32(smoothx), C.l_int32(smoothy))}
}

////////////////////////////////////////////////////

func pixConvertRGBToGrayFast(pixs Pix) Pix {
	return Pix{p: C.pixConvertRGBToGrayFast(pixs.p)}
}

func pixConvertTo8(pixs Pix, cmapflag int) Pix {
	return Pix{p: C.pixConvertTo8(pixs.p, C.l_int32(cmapflag))}
}

func pixConvertTo1(pixs Pix, threshold int) Pix {
	return Pix{p: C.pixConvertTo1(pixs.p, C.l_int32(threshold))}
}

////////////////////////////////////////////////////

func pixCopy(pixd Pix, pixs Pix) Pix {
	return Pix{p: C.pixCopy(pixd.p, pixs.p)}
}

func pixInvert(pixd Pix, pixs Pix) Pix {
	return Pix{p: C.pixInvert(pixd.p, pixs.p)}
}

func pixRasterop(pixd Pix, dx int, dy int, dw int, dh int, op int, pixs Pix, sx int, sy int) bool {
	return C.pixRasterop(pixd.p, C.l_int32(dx), C.l_int32(dy), C.l_int32(dw), C.l_int32(dh), C.l_int32(op), pixs.p, C.l_int32(sx), C.l_int32(sy)) == 0
}

func pixOr(pixd Pix, pixs1 Pix, pixs2 Pix) {
	C.pixOr(pixd.p, pixs1.p, pixs2.p)
}

func pixXor(pixd Pix, pixs1 Pix, pixs2 Pix) {
	C.pixXor(pixd.p, pixs1.p, pixs2.p)
}

func pixAnd(pixd Pix, pixs1 Pix, pixs2 Pix) {
	C.pixAnd(pixd.p, pixs1.p, pixs2.p)
}

////////////////////////////////////////////////////

func pixRemoveBorderConnComps(pixs Pix, connectivity int) Pix {
	return Pix{p: C.pixRemoveBorderConnComps(pixs.p, C.l_int32(connectivity))}
}

func pixCloseGray(pixs Pix, hsize int, vsize int) Pix {
	return Pix{p: C.pixCloseGray(pixs.p, C.int32_t(hsize), C.int32_t(vsize))}
}

func pixThresholdToValue(pixd Pix, pixs Pix, threshval int, setval int) Pix {
	return Pix{p: C.pixThresholdToValue(pixd.p, pixs.p, C.int32_t(threshval), C.int32_t(setval))}
}

func pixErodeGray(pixs Pix, hsize int, vsize int) Pix {
	return Pix{p: C.pixErodeGray(pixs.p, C.int32_t(hsize), C.int32_t(vsize))}
}

func pixDilateGray(pixs Pix, hsize int, vsize int) Pix {
	return Pix{p: C.pixDilateGray(pixs.p, C.int32_t(hsize), C.int32_t(vsize))}
}
