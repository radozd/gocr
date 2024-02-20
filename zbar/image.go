//go:build windows

package zbar

// #include <stdlib.h>
import "C"

type image uintptr

func newImage(width int, height int, data uintptr, size int) image {
	//cFormat := C.CString(format)
	//defer C.free(unsafe.Pointer(cFormat))

	img, _, _ := zbar_image_create.Call()
	zbar_image_set_size.Call(img, uintptr(width), uintptr(height))
	zbar_image_set_format.Call(img, uintptr(C.ulong(0x30303859))) // unsafe.Pointer(cFormat) // Y800 (grayscale)
	zbar_image_set_data.Call(img, data, uintptr(size), uintptr(0))

	return image(img)
}

func (img image) destroy() {
	zbar_image_destroy.Call(uintptr(img))
}
