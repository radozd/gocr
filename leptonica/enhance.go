package leptonica

import "unsafe"

func (pix Pix) EnhancedCopy(opt EnhanceOptions) Pix {
	var enhanced Pix

	_, _, d := pix.GetDimensions()
	if d != 8 && d != 32 {
		enhanced = pixConvertTo8(pix, 0)
	} else {
		enhanced = pixCopy(NullPix, pix)
	}
	if enhanced == NullPix {
		return NullPix
	}

	if opt.TileX > 0 {
		tmp := pixBackgroundNorm(enhanced, NullPix, NullPix, opt.TileX, opt.TileY,
			opt.Thresh, opt.MinCount, opt.WhitePoint, opt.SmoothX, opt.SmoothY)
		enhanced.Destroy()

		if tmp == NullPix {
			return NullPix
		}
		enhanced = tmp
	}

	if opt.RemoveBorders > 0 {
		// https://github.com/DanBloomberg/leptonica/issues/590
		pix2 := pixConvertTo1(enhanced, opt.RemoveBorders)
		if pix2 == NullPix {
			return enhanced
		}
		pix3 := pixRemoveBorderConnComps(pix2, 8)
		if pix3 == NullPix {
			pix2.Destroy()
			return enhanced
		}
		pixXor(pix2, pix2, pix3)
		pix3.Destroy()

		enhanced.pixSetMasked(pix2, uint(opt.WhitePoint)+256*uint(opt.WhitePoint)+256*256*uint(opt.WhitePoint))
		pix2.Destroy()
	}

	if opt.Gamma > 0 {
		pixGammaTRC(enhanced, enhanced, float32(opt.Gamma)/100.0, opt.GammaMin, opt.GammaMax)
	}

	if opt.Factor > 0 {
		pixContrastTRC(enhanced, enhanced, float32(opt.Factor)/100.0)
	}

	return enhanced
}

func (pix Pix) GetGrayCopy(mode GrayCastMode, opt GrayOptions) Pix {
	var gray Pix

	_, _, d := pix.GetDimensions()
	if d == 32 {
		gray = pixConvertRGBToGrayFast(pix)
		if gray == NullPix {
			return NullPix
		}

		if mode != GRAY_SIMPLE {
			var mask Pix
			if mode == GRAY_CAST_KEEP_ONLY_COLORS {
				mask = pixMaskOverGrayPixels(pix, opt.WhitePoint, opt.Saturation)
			} else if mode == GRAY_CAST_REMOVE_COLORS {
				mask = pixMaskOverGrayPixels(pix, opt.WhitePoint, opt.Saturation)
				if mask == NullPix {
					return gray
				}
				pixInvert(mask, mask)
			} else if mode == GRAY_CAST_REMOVE_COLORS_2 {
				mask = pixMaskOverColorPixels(pix, opt.ThreshDiff, opt.MinDist)
			}

			if mask == NullPix {
				return gray
			}

			gray.pixSetMasked(mask, uint(opt.WhitePoint))
			mask.Destroy()
		}
	} else if d != 8 {
		gray = pixConvertTo8(pix, 0)
	} else {
		gray = pixCopy(NullPix, pix)
	}

	if gray == NullPix {
		return NullPix
	}

	if opt.MaximizeBrightness {
		pixThresholdToValue(gray, gray, opt.WhitePoint, 255)
	}

	return gray
}

func (pix Pix) GetRawGrayData() []byte {
	gray := pix.GetGrayCopy(GRAY_CAST_REMOVE_COLORS, DefaultGrayOptions)
	if gray == NullPix {
		return nil
	}
	defer gray.Destroy()

	w, h, _ := gray.GetDimensions()
	if w == 0 || h == 0 {
		return nil
	}

	wpl := pixGetWpl(gray)
	raw := pixGetData(gray)
	if raw == nil {
		return nil
	}

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
