package zbar

// #include <stdlib.h>
import "C"

func newImage(width int, height int, data uintptr, size int) image {
	//cFormat := C.CString(format)
	//defer C.free(unsafe.Pointer(cFormat))

	return zbar_image_create(width, height, data, size)
}

func (img image) destroy() {
	zbar_image_destroy(img)
}
