package gocr

import (
	"fmt"
	"os"
	"testing"

	"github.com/radozd/gocr/leptonica"
)

func TestConvertGrey(t *testing.T) {
	pix := leptonica.NewPixFromFile("test.png")
	if pix == 0 {
		t.Error("error loading pix from file")
	}
	defer leptonica.DestroyPix(pix)

	w, h, d := leptonica.GetPixDimensions(pix)
	t.Log("img:", w, h, d)

	gray := leptonica.GetRawGrayData(pix)
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
