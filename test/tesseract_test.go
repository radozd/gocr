package main

import (
	"strings"
	"testing"

	"github.com/radozd/gocr/leptonica"
	"github.com/radozd/gocr/tesseract"
)

func TestOcr(t *testing.T) {
	pix := leptonica.NewPixFromFile("cda94bbc460be2a16fb7e719292e8593.tif") //("test.png")
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
	//t.Log("\n" + tess.Text())

	resIt := tess.GetIterator()
	pageIt := resIt.AsPageIterator()
	defer resIt.Delete()

	level := tesseract.RIL_BLOCK

	good := true
	for good {
		//if pageIt.IsAtBeginningOf(level) {
		text, goodness := resIt.GetUTF8Text(level)
		if s := strings.TrimSpace(text); s != "" {
			l, tt, r, b := pageIt.BoundingBox(level)
			sparsity := float32((r-l)*(b-tt)) / float32(len(s)) / 1000
			if sparsity < 4 {
				t.Logf("%s\ngoodness=%f\nsparsity=%.3f\n", text, goodness, sparsity)
			} else {
				t.Logf("****************************************\n")

				it2 := resIt.Copy()
				pit2 := it2.AsPageIterator()
				good2 := true
				for good2 {
					text, goodness := it2.GetUTF8Text(tesseract.RIL_TEXTLINE)
					if s := strings.TrimSpace(text); s != "" {
						l, tt, r, b = pit2.BoundingBox(tesseract.RIL_TEXTLINE)
						sparsity = float32((r-l)*(b-tt)) / float32(len(s)) / 1000
						t.Logf("%s\ngoodness=%f\nsparsity=%.3f\n", text, goodness, sparsity)
					}
					good2 = pit2.Next(tesseract.RIL_TEXTLINE)
					if pit2.IsAtBeginningOf(level) {
						break
					}
				}
				it2.Delete()
			}
		}
		//}

		good = pageIt.Next(level)
	}
}
