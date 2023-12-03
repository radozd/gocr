package leptonica

import (
	"fmt"
	"os"
	"testing"
)

func TestConvertGrey(t *testing.T) {
	pix := NewPixFromFile("test.png")
	if pix == 0 {
		t.Error("error loading pix from file")
	}
	defer DestroyPix(pix)

	w, h, d := GetPixDimensions(pix)
	t.Log("img:", w, h, d)

	gray := GetRawGrayData(pix)
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
