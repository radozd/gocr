package leptonica

func (boxa Boxa) Draw(draw func(x int, y int, w int, h int)) {
	for i := 0; i < boxaGetCount(boxa); i++ {
		box := boxaGetBox(boxa, i, L_CLONE)

		x, y, w, h := boxGetGeometry(box)
		draw(x, y, w, h)

		boxDestroy(&box)
	}
}

////////////////////////////////////////////////////

func (pix Pix) MaskSquares(thresh int, block int, sqMin int, sqMax int) Pix {
	mask := pixConvertTo1(pix, thresh)

	brick := pixCloseBrick(NullPix, mask, block, block)
	boxes := pixConnCompBB(brick, 8)
	brick.Destroy()

	boxes.Draw(func(x int, y int, w int, h int) {
		diff := 100 * w / h
		if sqMin <= w && w <= sqMax && sqMin <= h && h <= sqMax && 80 < diff && diff < 120 {
			mask.FillRect(x-1, y-1, w+2, h+2, true)
		} else {
			mask.FillRect(x, y, w, h, false)
		}
	})

	boxaDestroy(&boxes)

	return mask
}

func (pix Pix) MaskBars(thresh int, brMin int, brMax int, brWidth int, brHeight int) Pix {
	pixe := pixSobelEdgeFilter(pix, 2 /*L_ALL_EDGES*/)
	pixb := pixConvertTo1(pixe, thresh)
	pixDestroy(&pixe)
	pixInvert(pixb, pixb)

	brick1 := pixCloseBrick(NullPix, pixb, brMax, brMin)
	brick2 := pixOpenBrick(NullPix, pixb, brMax, brMin)
	pixXor(brick2, brick2, brick1)
	pixDestroy(&brick1)
	pixOpenBrick(brick2, brick2, brWidth, brHeight)

	brick1 = pixCloseBrick(NullPix, pixb, brMin, brMax)
	mask := pixOpenBrick(NullPix, pixb, brMin, brMax)
	pixXor(mask, mask, brick1)
	pixDestroy(&brick1)
	pixOpenBrick(mask, mask, brHeight, brWidth)

	pixDestroy(&pixb)

	pixOr(mask, mask, brick2)
	pixDestroy(&brick2)

	boxes := pixConnCompBB(mask, 8)
	boxes.Draw(func(x int, y int, w int, h int) {
		mask.FillRect(x-1, y-1, w+2, h+2, true)
	})
	boxaDestroy(&boxes)

	return mask
}

func (pix Pix) MaskLines(thresh int, lenMin int) Pix {
	pixb := pixConvertTo1(pix, 2*thresh)

	brick1 := pixOpenBrick(NullPix, pixb, lenMin, 1)
	pixDilateBrick(brick1, brick1, 1, 3)

	brick2 := pixOpenBrick(NullPix, pixb, 1, lenMin)
	pixDilateBrick(brick2, brick2, 5, 1)
	pixb.Destroy()

	pixOr(brick1, brick1, brick2)
	brick2.Destroy()

	return brick1
}

func (box Box) overlap(x2 int, y2 int, w2 int, h2 int) bool {
	x1, y1, w1, h1 := boxGetGeometry(box)

	if x1+w1 <= x2 || x2+w2 <= x1 || y1+h1 <= y2 || y2+h2 <= y1 {
		return false
	}
	return true
}

var DefaultBoxWeight = func(w int, h int) int {
	return (w*h + w + h) / 2
}

func (box Box) weightedGeometry() (x int, y int, w int, h int) {
	x, y, w, h = boxGetGeometry(box)

	delta := DefaultBoxWeight(w, h)
	x -= delta
	y -= delta
	w += 2 * delta
	h += 2 * delta

	return
}

func (pix Pix) MaskSpecks(thresh int, max int, weight int) Pix {
	mask := pixConvertTo1(pix, 2*thresh)

	width, height, _ := pixGetDimensions(mask)
	grid_x := (width + 9) / 10
	grid_y := (height + 9) / 10
	grid := [11][11][]int{}

	brick := pixCloseBrick(NullPix, mask, 3, 3)
	boxes := pixConnCompBB(brick, 8)
	pixDestroy(&brick)

	n := boxaGetCount(boxes)
	speck := make([]bool, n)
	weights := make([]int, n)

	for i := 0; i < n; i++ {
		box := boxaGetBox(boxes, i, L_CLONE)

		x, y, w, h := boxGetGeometry(box)
		if w <= max && h <= max {
			speck[i] = true
		} else {
			mask.FillRect(x, y, w, h, false)
			speck[i] = false
		}

		for row := y / grid_y; row <= (y+h)/grid_y; row++ {
			for col := x / grid_x; col <= (x+w)/grid_x; col++ {
				grid[row][col] = append(grid[row][col], i)
			}
		}
		weights[i] = DefaultBoxWeight(w, h)

		boxDestroy(&box)
	}

	//matches := 0
	for i := 0; i < n; i++ {
		if !speck[i] {
			continue
		}

		box_i := boxaGetBox(boxes, i, L_CLONE)
		x, y, w, h := box_i.weightedGeometry()
		if x < 0 {
			x = 0
		}
		if y < 0 {
			y = 0
		}
		if x+w > width {
			w = width - x
		}
		if y+h > height {
			h = height - y
		}

		indices := make(map[int]bool)
		for row := y / grid_y; row <= (y+h)/grid_y; row++ {
			for col := x / grid_x; col <= (x+w)/grid_x; col++ {
				for _, z := range grid[row][col] {
					indices[z] = true
				}
			}
		}

		ok := true
		for j := range indices {
			if i != j {
				box_j := boxaGetBox(boxes, j, L_CLONE)
				o := box_j.overlap(x, y, w, h)
				boxDestroy(&box_j)

				if o {
					if weights[i]*weights[j] > weight {
						ok = false
						break
					}
				}
			}
		}

		if ok {
			mask.FillRect(x, y, w, h, true)
			//speck[i] = false
			//matches++
		} else {
			mask.FillRect(x, y, w, h, false)
		}

		boxDestroy(&box_i)
	}

	boxaDestroy(&boxes)

	//fmt.Printf("matches = %d\n\n", matches)

	return mask
}

func (pix Pix) MaskAll(opt MaskOptions) {
	//t1 := time.Now()
	mask1 := pix.MaskSquares(opt.Thresh, opt.SqrBlock, opt.SqrMin, opt.SqrMax)
	pix.pixSetMasked(mask1, 0xFFFFFFFF)
	//t2 := time.Now()
	mask2 := pix.MaskBars(opt.Thresh, opt.BarMin, opt.BarMax, opt.BarWidth, opt.BarHeight)
	pix.pixSetMasked(mask2, 0xFFFFFFFF)
	//t3 := time.Now()
	mask3 := pix.MaskLines(opt.Thresh, opt.LinMin)
	pix.pixSetMasked(mask3, 0xFFFFFFFF)
	//t4 := time.Now()
	mask4 := pix.MaskSpecks(opt.Thresh, opt.SpMax, opt.SpWeight)
	pix.pixSetMasked(mask4, 0xFFFFFFFF)
	//t5 := time.Now()

	//fmt.Printf("squares = %v\nbars = %v\nlines = %v\nnoise = %v\n\n", t2.Sub(t1), t3.Sub(t2), t4.Sub(t3), t5.Sub(t4))

	mask1.WriteToFile("b_mask_1.png", PNG)
	mask2.WriteToFile("b_mask_2.png", PNG)
	mask3.WriteToFile("b_mask_3.png", PNG)
	mask4.WriteToFile("b_mask_4.png", PNG)

	//pixOr(mask1, mask1, mask2)
	//pixOr(mask1, mask1, mask3)
	//pixOr(mask1, mask1, mask4)
	mask1.Destroy()
	mask2.Destroy()
	mask3.Destroy()
	mask4.Destroy()
}
