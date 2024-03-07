package leptonica

// #include <stdlib.h>
import "C"

import (
	"errors"
	"strconv"
	"unsafe"
)

func NewPixFromFile(filename string) Pix {
	return pixRead(filename)
}

func NewPixFromMem(image []byte) Pix {
	return pixReadMem(image)
}

func (pix *Pix) Destroy() {
	pixDestroy(pix)
}

func (pix Pix) WriteToFile(filename string, format ImageType) error {
	if !pixWrite(filename, pix, int(format)) {
		return errors.New("error saving pix: " + filename + " (format: " + strconv.Itoa(int(format)) + ")")
	}
	return nil
}

func (pix Pix) WriteToMem(format ImageType) ([]byte, error) {
	return pixWriteMem(pix, int(format)), nil
}

func (pix Pix) GetDimensions() (int, int, int) {
	return pixGetDimensions(pix)
}

func (pix Pix) GetScaledCopy(width int, height int) Pix {
	return pixScaleToSize(pix, width, height)
}

func (pix Pix) FillRect(x int, y int, w int, h int, white bool) {
	PIX_SET := 15
	PIX_CLR := 0
	var op int
	if white {
		op = PIX_SET
	} else {
		op = PIX_CLR
	}
	pixRasterop(pix, x, y, w, h, op, pix, x, y)
}

func (pix Pix) GetRotated180Copy() Pix {
	return pixRotate180(NullPix, pix)
}

func (pix Pix) GetRotatedCopy(angle float32) Pix {
	return pixRotate(pix, angle)
}

func (pix Pix) GetDeskewedCopy(redsearch int) Pix {
	return pixDeskew(pix, redsearch)
}

func (pix Pix) GetDeskewedCopyAndAngle(redsearch int) (Pix, float32) {
	return pixFindSkewAndDeskew(pix, redsearch)
}

func (pix Pix) GetGrayCopy(mode GrayCastMode, opt GrayOptions) Pix {
	var gray Pix

	_, _, d := pix.GetDimensions()
	if d == 32 {
		gray = pixConvertRGBToGrayFast(pix)
		if mode != GRAY_SIMPLE {
			var mask Pix
			if mode == GRAY_CAST_KEEP_ONLY_COLORS {
				mask = pixMaskOverGrayPixels(pix, opt.WhitePoint, opt.Saturation)
			} else if mode == GRAY_CAST_REMOVE_COLORS {
				mask = pixMaskOverGrayPixels(pix, opt.WhitePoint, opt.Saturation)
				pixInvert(mask, mask)
			} else if mode == GRAY_CAST_REMOVE_COLORS_2 {
				mask = pixMaskOverColorPixels(pix, opt.ThreshDiff, opt.MinDist)
			}
			pixPaintThroughMask(gray, mask, 0, 0, uint(opt.WhitePoint))
			mask.Destroy()
		}
	} else if d != 8 {
		gray = pixConvertTo8(pix, 0)
	} else {
		gray = pixCopy(NullPix, pix)
	}
	return Pix(gray)
}

func (pix Pix) EnhancedCopy(opt EnhanceOptions) Pix {
	var enhanced Pix

	_, _, d := pix.GetDimensions()
	if d != 8 && d != 32 {
		enhanced = pixConvertTo8(pix, 0)
	} else {
		enhanced = pixCopy(NullPix, pix)
	}

	if opt.TileX > 0 {
		tmp := pixBackgroundNorm(enhanced, NullPix, NullPix, opt.TileX, opt.TileY,
			opt.Thresh, opt.MinCount, opt.BgVal, opt.SmoothX, opt.SmoothY)
		enhanced.Destroy()
		enhanced = tmp
	}

	if opt.Gamma > 0 {
		pixGammaTRC(enhanced, enhanced, opt.Gamma, opt.GammaMin, opt.GammaMax)
	}

	if opt.Factor > 0 {
		pixContrastTRC(enhanced, enhanced, opt.Factor)
	}

	return enhanced
}

func (pix Pix) GetRawGrayData() []byte {
	gray := pix.GetGrayCopy(GRAY_CAST_REMOVE_COLORS, DefaultGrayOptions)
	if gray == NullPix {
		return nil
	}
	defer gray.Destroy()

	w, h, _ := gray.GetDimensions()
	wpl := pixGetWpl(gray)
	raw := pixGetData(gray)

	rowlen := 4 * int(wpl)
	pixels := unsafe.Slice((*byte)(unsafe.Pointer(raw)), rowlen*h)

	bytes := make([]byte, w*h)
	for i := 0; i < h; i++ {
		ofs := i * rowlen
		line := pixels[ofs : ofs+rowlen]
		for j := 0; j < w; j++ {
			val := line[j^3]
			bytes[w*i+j] = val
		}
	}

	return bytes
}
