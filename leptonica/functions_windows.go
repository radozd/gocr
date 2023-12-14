package leptonica

// #include <stdlib.h>
import "C"

import (
	"syscall"
)

var (
	leptonicaDll = syscall.NewLazyDLL("leptonica-1.82.0.dll")

	pixDestroy = leptonicaDll.NewProc("pixDestroy")
	lept_free  = leptonicaDll.NewProc("lept_free")

	pixRead     = leptonicaDll.NewProc("pixRead")
	pixReadMem  = leptonicaDll.NewProc("pixReadMem")
	pixWrite    = leptonicaDll.NewProc("pixWrite")
	pixWriteMem = leptonicaDll.NewProc("pixWriteMem")

	pixRotate180   = leptonicaDll.NewProc("pixRotate180")
	pixAddBorder   = leptonicaDll.NewProc("pixAddBorder")   //PIX *pixAddBorder (PIX *pixs, l_int32 npix, l_uint32 val);
	pixScaleToSize = leptonicaDll.NewProc("pixScaleToSize") //PIX *pixScaleToSize(PIX *pixs, l_int32 wd, l_int32 hd)

	pixMaskOverColorPixels = leptonicaDll.NewProc("pixMaskOverColorPixels") //PIX *pixMaskOverColorPixels(PIX *pixs, l_int32 threshdiff, l_int32 mindist)
	pixMaskOverGrayPixels  = leptonicaDll.NewProc("pixMaskOverGrayPixels")  //PIX *pixMaskOverGrayPixels(PIX *pixs, l_int32  maxlimit, l_int32  satlimit)
	pixInvert              = leptonicaDll.NewProc("pixInvert")              //PIX *pixInvert(PIX *pixd, PIX *pixs);
	pixPaintThroughMask    = leptonicaDll.NewProc("pixPaintThroughMask")    //l_int32 pixPaintThroughMask(PIX *pixd, PIX *pixm, l_int32 x, l_int32 y, l_uint32 val)

	pixContrastTRC    = leptonicaDll.NewProc("pixContrastTRC")    //PIX *pixContrastTRC(PIX *pixd, PIX *pixs, l_float32 factor)
	pixGammaTRC       = leptonicaDll.NewProc("pixGammaTRC")       //PIX* pixGammaTRC(PIX *pixd, PIX *pixs, l_float32 gamma, l_int32 minval, l_int32 maxval)
	pixBackgroundNorm = leptonicaDll.NewProc("pixBackgroundNorm") // PIX *pixBackgroundNormSimple(PIX *pixs, PIX *pixim, PIX *pixg,
	// l_int32 sx, l_int32 sy, l_int32 thresh, l_int32 mincount, l_int32 bgval, l_int32 smoothx, l_int32 smoothy)

	pixConvertRGBToGrayFast = leptonicaDll.NewProc("pixConvertRGBToGrayFast")
	pixConvertTo8           = leptonicaDll.NewProc("pixConvertTo8")
	pixCopy                 = leptonicaDll.NewProc("pixCopy")

	pixGetData       = leptonicaDll.NewProc("pixGetData")
	pixGetDimensions = leptonicaDll.NewProc("pixGetDimensions")
	pixGetWpl        = leptonicaDll.NewProc("pixGetWpl")
)
