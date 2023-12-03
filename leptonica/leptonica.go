package leptonica

// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

type Pix uintptr

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

func DestroyPix(pix Pix) {
	pixDestroy.Call(uintptr(unsafe.Pointer(&pix)))
}

func GetPixDimensions(pix Pix) (int, int, int) {
	cW := C.int(0)
	cH := C.int(0)
	cD := C.int(0)
	code, _, _ := pixGetDimensions.Call(uintptr(pix), uintptr(unsafe.Pointer(&cW)), uintptr(unsafe.Pointer(&cH)), uintptr(unsafe.Pointer(&cD)))
	if code != 0 {
		return 0, 0, 0
	}
	return int(cW), int(cH), int(cD)
}

func Rotate180(pixs Pix) Pix {
	pix, _, _ := pixRotate180.Call(uintptr(0), uintptr(pixs))
	return Pix(pix)
}

func ConvertRGBToGrayFast(pixs Pix) Pix {
	pix, _, _ := pixConvertRGBToGrayFast.Call(uintptr(pixs))
	return Pix(pix)
}

func GetRawGrayData(pix Pix) []byte {
	var gray uintptr

	_, _, d := GetPixDimensions(pix)
	if d == 32 {
		gray, _, _ = pixConvertRGBToGrayFast.Call(uintptr(pix))
	} else if d != 8 {
		gray, _, _ = pixConvertTo8.Call(uintptr(pix), 0)
	} else {
		gray = uintptr(pix)
	}
	if gray == 0 {
		return nil
	}
	defer func() {
		if gray != uintptr(pix) {
			DestroyPix(Pix(gray))
		}
	}()

	w, h, _ := GetPixDimensions(Pix(gray))
	wpl, _, _ := pixGetWpl.Call(gray)
	raw, _, _ := pixGetData.Call(gray)

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
