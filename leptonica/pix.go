package leptonica

// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

type Pix uintptr

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

func NewPixFromMem(image *[]byte) Pix {
	size := C.size_t(len(*image))
	pix, _, _ := pixReadMem.Call(uintptr(unsafe.Pointer(&(*image)[0])), uintptr(size))
	return Pix(pix)
}

func (pix Pix) Destroy() {
	pixDestroy.Call(uintptr(unsafe.Pointer(&pix)))
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

func (pix Pix) GetRotated180Copy() Pix {
	pix180, _, _ := pixRotate180.Call(uintptr(0), uintptr(pix))
	return Pix(pix180)
}

func (pix Pix) GetGrayCopy(mode GrayCastMode) Pix {
	var gray uintptr

	_, _, d := pix.GetDimensions()
	if d == 32 {
		gray, _, _ = pixConvertRGBToGrayFast.Call(uintptr(pix))
		if mode != GRAY_SIMPLE {
			mask, _, _ := pixMaskOverGrayPixels.Call(uintptr(pix), uintptr(255), 60)
			if mode == GRAY_CAST_REMOVE_COLORS {
				//mask, _, _ = pixMaskOverColorPixels.Call(uintptr(pix), uintptr(50), uintptr(1))
				pixInvert.Call(mask, mask)
			}
			defer Pix(mask).Destroy()

			_, _, _ = pixPaintThroughMask.Call(gray, mask, uintptr(0), uintptr(0), uintptr(255))
		}
	} else if d != 8 {
		gray, _, _ = pixConvertTo8.Call(uintptr(pix), 0)
	} else {
		gray, _, _ = pixCopy.Call(uintptr(0), uintptr(pix))
	}
	return Pix(gray)
}

func (pix Pix) GetRawGrayData() []byte {
	gray := pix.GetGrayCopy(GRAY_CAST_REMOVE_COLORS)
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
