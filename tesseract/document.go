package tesseract

type TextBlock struct {
	Goodness float32
	Value    string
	X        int
	Y        int
	Width    int
	Height   int
}

func Process(api *Api, level PageIteratorLevel) []TextBlock {
	resIt := api.GetIterator()
	pageIt := resIt.GetPageIterator()
	defer resIt.Delete()

	blocks := make([]TextBlock, 0)

	good := true
	for good {
		if pageIt.IsAtBeginningOf(level) {
			x1, y1, x2, y2 := pageIt.BoundingBox(level)
			text, goodness := resIt.GetUTF8Text(level)
			blocks = append(blocks, TextBlock{
				Goodness: goodness,
				Value:    text,
				X:        x1,
				Y:        y1,
				Width:    x2 - x1,
				Height:   y2 - y1,
			})
		}

		good = pageIt.Next(level)
	}
	return blocks
}
