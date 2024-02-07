package leptonica

// #include <stdlib.h>
import "C"

import (
	"errors"
	"math"
	"strconv"
	"unsafe"
)

type Pix uintptr

type ImageType int32

const (
	UNKNOWN ImageType = iota
	BMP
	JFIF_JPEG
	PNG
	TIFF
	TIFF_PACKBITS
	TIFF_RLE
	TIFF_G3
	TIFF_G4
	TIFF_LZW
	TIFF_ZIP
	PNM
	PS
	GIF
	JP2
	WEBP
	LPDF
	DEFAULT
	SPIX
)

type GrayCastMode int

const (
	GRAY_SIMPLE GrayCastMode = iota
	GRAY_CAST_REMOVE_COLORS
	GRAY_CAST_KEEP_ONLY_COLORS
)

func NewPixFromFile(filename string) Pix {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	pix, _, _ := pixRead.Call(uintptr(unsafe.Pointer(cFilename)))
	return Pix(pix)
}

func NewPixFromMem(image []byte) Pix {
	size := C.size_t(len(image))
	pix, _, _ := pixReadMem.Call(uintptr(unsafe.Pointer(unsafe.SliceData(image))), uintptr(size))
	return Pix(pix)
}

func (pix Pix) Destroy() {
	pixDestroy.Call(uintptr(unsafe.Pointer(&pix)))
}

func (pix Pix) WriteToFile(filename string, format ImageType) error {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	cFormat := C.int(format)

	code, _, _ := pixWrite.Call(uintptr(unsafe.Pointer(cFilename)), uintptr(pix), uintptr(cFormat))
	if code != 0 {
		return errors.New("error saving pix: " + filename + " (format: " + strconv.Itoa(int(format)) + ")")
	}
	return nil
}

func (pix Pix) WriteToMem(format ImageType) ([]byte, error) {
	var cMem *C.char = nil
	cSize := C.size_t(0)
	cFormat := C.int(format)

	code, _, _ := pixWriteMem.Call(uintptr(unsafe.Pointer(&cMem)), uintptr(unsafe.Pointer(&cSize)), uintptr(pix), uintptr(cFormat))
	if code != 0 {
		return nil, errors.New("error " + strconv.Itoa(int(code)))
	}
	defer lept_free.Call(uintptr(unsafe.Pointer(cMem)))

	data := C.GoBytes(unsafe.Pointer(cMem), C.int(cSize))

	return data, nil
}

func (pix Pix) GetDimensions() (int, int, int) {
	cW := C.int(0)
	cH := C.int(0)
	cD := C.int(0)
	code, _, _ := pixGetDimensions.Call(uintptr(pix), uintptr(unsafe.Pointer(&cW)), uintptr(unsafe.Pointer(&cH)), uintptr(unsafe.Pointer(&cD)))
	if code != 0 {
		return 0, 0, 0
	}
	return int(cW), int(cH), int(cD)
}

func (pix Pix) GetScaledCopy(width int, height int) Pix {
	scaled, _, _ := pixScaleToSize.Call(uintptr(pix), uintptr(width), uintptr(height))
	return Pix(scaled)
}

func (pix Pix) GetRotated180Copy() Pix {
	pix180, _, _ := pixRotate180.Call(uintptr(0), uintptr(pix))
	return Pix(pix180)
}

func (pix Pix) GetDeskewedCopy(redsearch int) Pix {
	dpix, _, _ := pixDeskew.Call(uintptr(pix), uintptr(redsearch))
	return Pix(dpix)
}

type GrayOptions struct {
	Saturation int
	WhitePoint int
}

var DefaultGrayOptions = GrayOptions{
	Saturation: 40,
	WhitePoint: 250,
}

func (pix Pix) GetGrayCopy(mode GrayCastMode, opt GrayOptions) Pix {
	var gray uintptr

	_, _, d := pix.GetDimensions()
	if d == 32 {
		gray, _, _ = pixConvertRGBToGrayFast.Call(uintptr(pix))
		if mode != GRAY_SIMPLE {
			mask, _, _ := pixMaskOverGrayPixels.Call(uintptr(pix), uintptr(255), uintptr(opt.Saturation))
			defer Pix(mask).Destroy()

			if mode == GRAY_CAST_REMOVE_COLORS {
				//mask, _, _ = pixMaskOverColorPixels.Call(uintptr(pix), uintptr(50), uintptr(1))
				pixInvert.Call(mask, mask)
			}
			_, _, _ = pixPaintThroughMask.Call(gray, mask, uintptr(0), uintptr(0), uintptr(opt.WhitePoint))
		}
	} else if d != 8 {
		gray, _, _ = pixConvertTo8.Call(uintptr(pix), 0)
	} else {
		gray, _, _ = pixCopy.Call(uintptr(0), uintptr(pix))
	}
	return Pix(gray)
}

type EnhanceOptions struct {
	TileX    int
	TileY    int
	Thresh   int
	MinCount int
	BgVal    int
	SmoothX  int
	SmoothY  int

	Gamma    float32
	GammaMin int
	GammaMax int

	Factor float32
}

var DefaultEnhanceOptions = EnhanceOptions{
	TileX:    10,
	TileY:    10,
	Thresh:   40,
	MinCount: 50,
	BgVal:    250,
	SmoothX:  1,
	SmoothY:  1,
	Gamma:    0,
	GammaMin: 20,
	GammaMax: 240,
	Factor:   0.5,
}

func (pix Pix) EnhancedCopy(opt EnhanceOptions) Pix {
	p := pix
	_, _, d := pix.GetDimensions()
	if d != 8 && d != 32 {
		pp, _, _ := pixConvertTo8.Call(uintptr(pix), 0)
		p = Pix(pp)
		defer p.Destroy()
	}

	var enhanced uintptr
	if opt.TileX > 0 {
		enhanced, _, _ = pixBackgroundNorm.Call(uintptr(p), 0, 0, uintptr(opt.TileX), uintptr(opt.TileY),
			uintptr(opt.Thresh), uintptr(opt.MinCount), uintptr(opt.BgVal), uintptr(opt.SmoothX), uintptr(opt.SmoothY))
	} else {
		enhanced, _, _ = pixCopy.Call(uintptr(0), uintptr(p))
	}

	if opt.Gamma > 0 {
		pixGammaTRC.Call(enhanced, enhanced, uintptr(math.Float32bits(opt.Gamma)), uintptr(opt.GammaMin), uintptr(opt.GammaMax))
	}

	if opt.Factor > 0 {
		pixContrastTRC.Call(enhanced, enhanced, uintptr(math.Float32bits(opt.Factor)))
	}

	return Pix(enhanced)
}

func (pix Pix) GetRawGrayData() []byte {
	gray := pix.GetGrayCopy(GRAY_CAST_REMOVE_COLORS, DefaultGrayOptions)
	if gray == 0 {
		return nil
	}
	defer gray.Destroy()

	w, h, _ := gray.GetDimensions()
	wpl, _, _ := pixGetWpl.Call(uintptr(gray))
	raw, _, _ := pixGetData.Call(uintptr(gray))

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
