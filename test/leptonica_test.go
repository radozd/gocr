package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/radozd/gocr/leptonica"
)

func TestConvertPix(t *testing.T) {
	buf, err := os.ReadFile("test.png")
	if err != nil {
		t.Error(err)
		return
	}
	pix := leptonica.NewPixFromMem(buf)
	if pix == leptonica.NullPix {
		t.Error("error loading pix from file")
		return
	}
	defer pix.Destroy()

	pix.WriteToFile("test_write1.jpg", leptonica.JFIF_JPEG)

	tmpbuf, err := pix.WriteToMem(leptonica.PNG)
	if err != nil {
		t.Error(err)
		return
	}
	os.WriteFile("test_write2.png", tmpbuf, os.ModePerm)
}

func TestConvertGrey(t *testing.T) {
	pix := leptonica.NewPixFromFile("colors.png")
	if pix == leptonica.NullPix {
		t.Error("error loading pix from file")
		return
	}
	defer pix.Destroy()

	w, h, d := pix.GetDimensions()
	t.Log("img:", w, h, d)

	gray := pix.GetRawGrayData()
	if gray == nil {
		t.Error("error loading pix")
		return
	}

	t.Log("gry size=", len(gray))

	f, err := os.Create("gray.pgm")
	if err != nil {
		t.Error(err)
		return
	}

	f.WriteString(fmt.Sprintf("P5\n%d %d\n255\n", w, h))
	f.Write(gray)
}

func TestDeskew(t *testing.T) {
	pix := leptonica.NewPixFromFile("rotated.jpg")
	if pix == leptonica.NullPix {
		t.Error("error loading pix from file")
		return
	}
	defer pix.Destroy()

	w, h, d := pix.GetDimensions()
	t.Log("img:", w, h, d)

	dpix := pix.GetDeskewedCopy(0)
	if dpix == leptonica.NullPix {
		t.Error("error loading pix")
		return
	}
	defer dpix.Destroy()

	dpix.WriteToFile("deskewed.jpg", leptonica.JFIF_JPEG)
}

func TestRotate(t *testing.T) {
	pix := leptonica.NewPixFromFile("test.png")
	if pix == leptonica.NullPix {
		t.Error("error loading pix from file")
		return
	}
	defer pix.Destroy()

	w, h, d := pix.GetDimensions()
	t.Log("img:", w, h, d)

	dpix := pix.GetRotatedCopy(-3)
	if dpix == leptonica.NullPix {
		t.Error("error loading pix")
		return
	}
	defer dpix.Destroy()

	pix2, angle := dpix.GetDeskewedCopyAndAngle(0)
	defer pix2.Destroy()

	t.Log(angle)

	pix2.WriteToFile("rotated.jpg", leptonica.JFIF_JPEG)
}

func TestRemoveBlackBorders(t *testing.T) {
	pix := leptonica.NewPixFromFile("image_black_borders.jpg")
	if pix == leptonica.NullPix {
		t.Error("error loading pix from file")
		return
	}
	defer pix.Destroy()

	tmp := pix.GetGrayCopy(leptonica.GRAY_CAST_REMOVE_COLORS, leptonica.DefaultGrayOptions)
	defer tmp.Destroy()

	pix1 := tmp.EnhancedCopy(leptonica.DefaultEnhanceOptions)
	defer pix1.Destroy()

	pix1.WriteToFile("image_black_borders_removed.jpg", leptonica.JFIF_JPEG)
}

func TestBlentRect(t *testing.T) {
	pix := leptonica.NewPixFromFile("blank32.png")
	if pix == leptonica.NullPix {
		t.Error("error loading pix from file")
		return
	}
	defer pix.Destroy()

	pix.BlendRect(100, 100, 200, 300, 0, 0.2)
	pix.BlendRect(150, 150, 200, 300, 0, 0.2)
	pix.WriteToFile("blank_blend2.png", leptonica.PNG)
}
