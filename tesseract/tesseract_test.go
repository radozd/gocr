package tesseract

import (
	"testing"

	"github.com/radozd/gocr/leptonica"
)

func TestOcr(t *testing.T) {
	pix := leptonica.NewPixFromFile("test.png")
	if pix == 0 {
		t.Error("error loading pix")
	}
	defer leptonica.DestroyPix(pix)

	tess := NewApi("fra")
	defer tess.Close()

	tess.SetPageSegMode(1) // PSM_AUTO_OSD = Automatic page segmentation with orientation and script detection. (OSD)
	tess.SetVariable("preserve_interword_spaces", "1")

	tess.SetImagePix(pix)
	tess.Recognize()

	t.Log("\n" + tess.Text())

	resIt := tess.GetIterator()
	pageIt := resIt.GetPageIterator()
	defer resIt.Delete()

	good := true
	for good {
		if pageIt.IsAtBeginningOf(RIL_BLOCK) {
			text, goodness := resIt.GetUTF8Text(RIL_BLOCK)
			t.Log(text, goodness, "\n")
		}

		good = pageIt.Next(RIL_BLOCK)
	}
}
