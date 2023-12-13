package leptonica

// #include <stdlib.h>
import "C"

import (
	"syscall"
)

var (
	leptonicaDll = syscall.NewLazyDLL("leptonica-1.82.0.dll")

	pixDestroy  = leptonicaDll.NewProc("pixDestroy")
	pixWrite    = leptonicaDll.NewProc("pixWrite")
	pixRead     = leptonicaDll.NewProc("pixRead")
	pixWriteMem = leptonicaDll.NewProc("pixWriteMem")
	pixReadMem  = leptonicaDll.NewProc("pixReadMem")

	pixRotate180 = leptonicaDll.NewProc("pixRotate180")
	pixAddBorder = leptonicaDll.NewProc("pixAddBorder") // PIX *pixAddBorder (PIX *pixs, l_int32 npix, l_uint32 val);

	pixMaskOverColorPixels = leptonicaDll.NewProc("pixMaskOverColorPixels") // PIX *pixMaskOverColorPixels(PIX *pixs, l_int32 threshdiff, l_int32 mindist)
	pixMaskOverGrayPixels  = leptonicaDll.NewProc("pixMaskOverGrayPixels")  // PIX *pixMaskOverGrayPixels(PIX *pixs, l_int32  maxlimit, l_int32  satlimit)
	pixInvert              = leptonicaDll.NewProc("pixInvert")              //PIX *pixInvert(PIX *pixd, PIX *pixs);
	pixPaintThroughMask    = leptonicaDll.NewProc("pixPaintThroughMask")    // l_int32 pixPaintThroughMask(PIX *pixd, PIX *pixm, l_int32 x, l_int32 y, l_uint32 val)

	pixConvertRGBToGrayFast = leptonicaDll.NewProc("pixConvertRGBToGrayFast")
	pixConvertTo8           = leptonicaDll.NewProc("pixConvertTo8")
	pixCopy                 = leptonicaDll.NewProc("pixCopy")

	pixGetData       = leptonicaDll.NewProc("pixGetData")
	pixGetDimensions = leptonicaDll.NewProc("pixGetDimensions")
	pixGetWpl        = leptonicaDll.NewProc("pixGetWpl")
)
