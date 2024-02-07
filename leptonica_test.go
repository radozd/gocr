package gocr

import (
	"fmt"
	"os"
	"testing"

	"github.com/radozd/gocr/leptonica"
)

func TestConvertPix(t *testing.T) {
	buf, err := os.ReadFile("tst.jpg")
	if err != nil {
		t.Error(err)
	}
	pix := leptonica.NewPixFromMem(buf)
	if pix == 0 {
		t.Error("error loading pix from file")
	}
	defer pix.Destroy()

	pix.WriteToFile("test_write1.jpg", leptonica.JFIF_JPEG)

	tmpbuf, err := pix.WriteToMem(leptonica.PNG)
	if err != nil {
		t.Error(err)
	}
	os.WriteFile("test_write2.png", tmpbuf, os.ModePerm)
}

func TestConvertGrey(t *testing.T) {
	pix := leptonica.NewPixFromFile("colors.png")
	if pix == 0 {
		t.Error("error loading pix from file")
	}
	defer pix.Destroy()

	w, h, d := pix.GetDimensions()
	t.Log("img:", w, h, d)

	gray := pix.GetRawGrayData()
	if gray == nil {
		t.Error("error loading pix")
	}

	t.Log("gry size=", len(gray))

	f, err := os.Create("gray.pgm")
	if err != nil {
		t.Error(err)
	}

	f.WriteString(fmt.Sprintf("P5\n%d %d\n255\n", w, h))
	f.Write(gray)
}

func TestDeskew(t *testing.T) {
	pix := leptonica.NewPixFromFile("skewed.jpg")
	if pix == 0 {
		t.Error("error loading pix from file")
	}
	defer pix.Destroy()

	w, h, d := pix.GetDimensions()
	t.Log("img:", w, h, d)

	dpix := pix.GetDeskewedCopy(0)
	if dpix == 0 {
		t.Error("error loading pix")
	}

	dpix.WriteToFile("deskewed.jpg", leptonica.JFIF_JPEG)
}

func TestEnhancePix(t *testing.T) {
	pix := leptonica.NewPixFromFile("bad.jpg")
	if pix == 0 {
		t.Error("error loading pix from file")
	}
	defer pix.Destroy()

	tmp := pix.GetGrayCopy(leptonica.GRAY_CAST_REMOVE_COLORS, leptonica.DefaultGrayOptions)
	defer tmp.Destroy()

	opt := leptonica.EnhanceOptions{
		TileX:    10,
		TileY:    10,
		Thresh:   40,
		MinCount: 50,
		BgVal:    250,
		SmoothX:  1,
		SmoothY:  1,
		Gamma:    0.6,
		GammaMin: 20,
		GammaMax: 240,
		Factor:   0.4,
	}
	en := tmp.EnhancedCopy(opt)
	defer en.Destroy()

	en.WriteToFile("bad_enhanced1.jpg", leptonica.JFIF_JPEG)
}

func TestEnhancePixLoop(t *testing.T) {
	pix := leptonica.NewPixFromFile("good-eprescr.tif")
	if pix == 0 {
		t.Error("error loading pix from file")
	}
	defer pix.Destroy()

	for i := 10; i < 100; i += 20 {
		opt := leptonica.EnhanceOptions{
			TileX:    10,
			TileY:    10,
			Thresh:   40,
			MinCount: i, //50,
			BgVal:    250,
			SmoothX:  1,
			SmoothY:  1,
			Gamma:    0.5, //float32(i) / 100, //0.6,
			GammaMin: 20,
			GammaMax: 240,
			Factor:   0.4, //float32(i) / 100, //0.4,
		}
		en := pix.EnhancedCopy(opt)
		tmp := en.GetGrayCopy(leptonica.GRAY_CAST_REMOVE_COLORS, leptonica.DefaultGrayOptions)

		tmp.WriteToFile(fmt.Sprintf("20240130173214_370-%d.jpg", i), leptonica.JFIF_JPEG)
		tmp.Destroy()
		en.Destroy()
	}
}
