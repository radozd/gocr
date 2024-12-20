package main

import (
	"testing"

	"github.com/radozd/gocr/leptonica"
	"github.com/radozd/gocr/zbar"
)

func TestZBar(t *testing.T) {
	pix := leptonica.NewPixFromFile("barcodes1.tif")
	if pix == leptonica.NullPix {
		t.Error("error loading pix from file")
	}
	defer pix.Destroy()

	enhanced := pix.EnhancedCopy(leptonica.DefaultEnhanceOptions)
	defer enhanced.Destroy()

	deskew, _ := enhanced.GetDeskewedCopyAndAngle(0)
	defer deskew.Destroy()

	gray := deskew.GetGrayCopy(leptonica.GRAY_CAST_REMOVE_COLORS, leptonica.DefaultGrayOptions)

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
