package zbar

func newImage(width int, height int, data uintptr, size int) image {
	return zbar_image_create(width, height, data, size)
}

func (img image) destroy() {
	zbar_image_destroy(img)
}
