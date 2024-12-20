package gocr

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

func TestEnhancePix(t *testing.T) {
	pix := leptonica.NewPixFromFile("bad.jpg")
	if pix == leptonica.NullPix {
		t.Error("error loading pix from file")
		return
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
	pix := leptonica.NewPixFromFile("new_worse.jpg")
	if pix == leptonica.NullPix {
		t.Error("error loading pix from file")
		return
	}
	defer pix.Destroy()

	for i := 0; i < 10; i++ {
		opt := leptonica.EnhanceOptions{
			TileX:    10, //10,
			TileY:    10, //10,
			Thresh:   40, //40,      //40,
			MinCount: 40, //10,
			BgVal:    250,
			SmoothX:  4,                   //2,
			SmoothY:  4,                   //2,
			Gamma:    float32(i)/10 + 0.5, //0.3, //0.3,
			GammaMin: 20,
			GammaMax: 240,
			Factor:   0.8, //float32(i) / 10, //0.8,
		}
		en := pix.EnhancedCopy(opt)
		deskew := en.GetDeskewedCopy(0)

		var opt2 = leptonica.GrayOptions{
			Saturation: 90,
			ThreshDiff: 90, //110, //40,
			MinDist:    2,
			WhitePoint: 250,
		}
		tmp := deskew.GetGrayCopy(leptonica.GRAY_CAST_REMOVE_COLORS, opt2)
		tmp.WriteToFile(fmt.Sprintf("enhloop-%d.jpg", i), leptonica.JFIF_JPEG)

		tmp.Destroy()
		deskew.Destroy()
		en.Destroy()
	}
}

func TestRect(t *testing.T) {
	pix := leptonica.NewPixFromFile("test.png")
	if pix == leptonica.NullPix {
		t.Error("error loading pix from file")
		return
	}
	defer pix.Destroy()

	w, h, d := pix.GetDimensions()
	t.Log("img:", w, h, d)

	pix.FillRect(10, 10, 10, 50, false)
	pix.FillRect(150, 150, 50, 50, true)

	pix.WriteToFile("rect.jpg", leptonica.JFIF_JPEG)
}

func TestLines(t *testing.T) {
	pix := leptonica.NewPixFromFile("b4affd2d20b7b97dd4fe91679a9e6363.tif")
	if pix == leptonica.NullPix {
		t.Error("error loading pix from file")
		return
	}
	defer pix.Destroy()

	deskew := pix.GetDeskewedCopy(0)
	defer deskew.Destroy()

	var opt2 = leptonica.GrayOptions{
		Saturation: 90,
		ThreshDiff: 90, //110, //40,
		MinDist:    2,
		WhitePoint: 250,
	}
	tmp := deskew.GetGrayCopy(leptonica.GRAY_CAST_REMOVE_COLORS, opt2)
	defer tmp.Destroy()

	tmp.RemoveLines(leptonica.DefaultLineOptions)

	tmp.WriteToFile("lines-no.jpg", leptonica.JFIF_JPEG)
}
