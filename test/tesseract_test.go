package main

import (
	"testing"

	"github.com/radozd/gocr/leptonica"
	"github.com/radozd/gocr/tesseract"
)

func TestOcr(t *testing.T) {
	pix := leptonica.NewPixFromFile("test.png")
	if pix == leptonica.NullPix {
		t.Error("error loading pix")
		return
	}
	defer pix.Destroy()

	tess := tesseract.NewApi("fra")
	defer tess.Close()

	tess.SetPageSegMode(tesseract.PSM_AUTO_OSD)
	tess.SetVariable("preserve_interword_spaces", "1")

	tess.SetImagePix(pix)
	tess.Recognize()

	t.Log(tess.GetPageOrientation())
	t.Log("\n" + tess.Text())

	resIt := tess.GetIterator()
	pageIt := resIt.GetPageIterator()
	defer resIt.Delete()

	good := true
	for good {
		if pageIt.IsAtBeginningOf(tesseract.RIL_BLOCK) {
			text, goodness := resIt.GetUTF8Text(tesseract.RIL_BLOCK)
			t.Log(text, goodness, "\n")
		}

		good = pageIt.Next(tesseract.RIL_BLOCK)
	}
}
