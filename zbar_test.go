package gocr

import (
	"testing"

	"github.com/radozd/gocr/leptonica"
	"github.com/radozd/gocr/zbar"
)

func TestZBar(t *testing.T) {
	pix := leptonica.NewPixFromFile("test.png")
	if pix == 0 {
		t.Error("error loading pix from file")
	}
	defer pix.Destroy()

	codes := zbar.Process(pix)
	if len(codes) == 0 {
		t.Error("no codes were found")
	}

	for _, c := range codes {
		t.Log(c)
	}
}
