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

func pixCountPixels(pix Pix) int {
	count := C.int(0)
	code := int(C.pixCountPixels(pix.p, &count, nil))
	if code != 0 {
		return -1
	}
	return int(count)
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

func pixSetMasked(pixd Pix, pixm Pix, val uint) int {
	return int(C.pixSetMasked(pixd.p, pixm.p, C.l_uint32(val)))
}

func pixCombineMasked(pixd Pix, pixs Pix, pixm Pix) int {
	return int(C.pixCombineMasked(pixd.p, pixs.p, pixm.p))
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

func pixThresholdToValue(pixd Pix, pixs Pix, threshval int, setval int) Pix {
	return Pix{p: C.pixThresholdToValue(pixd.p, pixs.p, C.int32_t(threshval), C.int32_t(setval))}
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

func pixBlendInRect(pixs Pix, box Box, val int, fract float32) {
	C.pixBlendInRect(pixs.p, box.p, C.l_uint32(val), C.float(fract))
}

////////////////////////////////////////////////////

func pixOpenBrick(pixd Pix, pixs Pix, hsize int, vsize int) Pix {
	return Pix{p: C.pixOpenBrick(pixd.p, pixs.p, C.int32_t(hsize), C.int32_t(vsize))}
}

func pixCloseBrick(pixd Pix, pixs Pix, hsize int, vsize int) Pix {
	return Pix{p: C.pixCloseBrick(pixd.p, pixs.p, C.int32_t(hsize), C.int32_t(vsize))}
}

func pixErodeBrick(pixd Pix, pixs Pix, hsize int, vsize int) Pix {
	return Pix{p: C.pixErodeBrick(pixd.p, pixs.p, C.int32_t(hsize), C.int32_t(vsize))}
}

func pixDilateBrick(pixd Pix, pixs Pix, hsize int, vsize int) Pix {
	return Pix{p: C.pixDilateBrick(pixd.p, pixs.p, C.int32_t(hsize), C.int32_t(vsize))}
}

func pixSobelEdgeFilter(pixs Pix, orientflag int) Pix {
	return Pix{p: C.pixSobelEdgeFilter(pixs.p, C.int32_t(orientflag))}
}

////////////////////////////////////////////////////

type Box struct {
	p *C.BOX
}

func boxCreate(x, y, w, h int) Box {
	return Box{p: C.boxCreate(C.int32_t(x), C.int32_t(y), C.int32_t(w), C.int32_t(h))}
}

func boxDestroy(box *Box) {
	C.boxDestroy((**C.BOX)(unsafe.Pointer(box)))
	box.p = nil
}

func boxGetGeometry(box Box) (int, int, int, int) {
	x := C.int(0)
	y := C.int(0)
	w := C.int(0)
	h := C.int(0)
	code := C.boxGetGeometry(box.p, &x, &y, &w, &h)
	if code != 0 {
		return 0, 0, 0, 0
	}
	return int(x), int(y), int(w), int(h)
}

func boxSetGeometry(box Box, x int, y int, w int, h int) int {
	return int(C.boxSetGeometry(box.p, C.int32_t(x), C.int32_t(y), C.int32_t(w), C.int32_t(h)))
}

////////////////////////////////////////////////////

type Boxa struct {
	p *C.BOXA
}

func boxaDestroy(boxa *Boxa) {
	C.boxaDestroy((**C.BOXA)(unsafe.Pointer(boxa)))
	boxa.p = nil
}

func boxaGetCount(boxa Boxa) int {
	return int(C.boxaGetCount(boxa.p))
}

func boxaGetBox(boxa Boxa, index int, accessflag int) Box {
	return Box{p: C.boxaGetBox(boxa.p, C.int32_t(index), C.int32_t(accessflag))}
}

////////////////////////////////////////////////////

func pixConnCompBB(pixs Pix, connectivity int) Boxa {
	return Boxa{p: C.pixConnCompBB(pixs.p, C.int32_t(connectivity))}
}

func pixRemoveBorderConnComps(pixs Pix, connectivity int) Pix {
	return Pix{p: C.pixRemoveBorderConnComps(pixs.p, C.l_int32(connectivity))}
}
