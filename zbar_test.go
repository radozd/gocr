package gocr

import (
	"testing"

	"github.com/radozd/gocr/leptonica"
	"github.com/radozd/gocr/zbar"
)

func TestZBar(t *testing.T) {
	pix := leptonica.NewPixFromFile("6e772912391b37eec9239ee9b8393756.tif")
	if pix == leptonica.NullPix {
		t.Error("error loading pix from file")
	}
	defer pix.Destroy()

	enhanced := pix.EnhancedCopy(leptonica.DefaultEnhanceOptions)
	defer enhanced.Destroy()

	deskew, _ := enhanced.GetDeskewedCopyAndAngle(0)
	defer deskew.Destroy()

	grayopt := leptonica.GrayOptions{
		Saturation: 150,
		WhitePoint: 250,

		ThreshDiff: 90,
		MinDist:    2,
	}
	gray := deskew.GetGrayCopy(leptonica.GRAY_CAST_REMOVE_COLORS, grayopt)

	scn := zbar.NewScanner()
	defer scn.Destroy()

	codes := scn.Process(gray)
	if len(codes) == 0 {
		t.Error("no codes were found")
	}

	for _, c := range codes {
		t.Log(c)
	}
}
