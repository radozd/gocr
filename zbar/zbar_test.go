package zbar

import (
	"testing"

	"github.com/radozd/gocr/leptonica"
)

func TestConvertGrey(t *testing.T) {
	pix := leptonica.NewPixFromFile("test.png")
	if pix == 0 {
		t.Error("error loading pix from file")
	}
	defer leptonica.DestroyPix(pix)

	codes := Process(pix)
	if len(codes) == 0 {
		t.Error("no codes were found")
	}

	for _, c := range codes {
		t.Log(c)
	}
}
