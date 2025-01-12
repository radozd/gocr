package leptonica

import (
	"errors"
	"strconv"
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

func (pix Pix) Copy() Pix {
	return pixCopy(NullPix, pix)
}

func (pix Pix) GetScaledCopy(width int, height int) Pix {
	return pixScaleToSize(pix, width, height)
}

func (pix Pix) FillRect(x int, y int, w int, h int, black bool) {
	const PIX_SET int = 15
	const PIX_CLR int = 0

	_, _, d := pixGetDimensions(pix)

	var op int
	if black != (d != 1) {
		op = PIX_SET
	} else {
		op = PIX_CLR
	}
	pixRasterop(pix, x, y, w, h, op, NullPix, 0, 0)
}

func (pix Pix) pixSetMasked(mask Pix, color uint) {
	pixSetMasked(pix, mask, color)
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
