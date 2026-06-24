package leptonica

func (boxa Boxa) Draw(draw func(x int, y int, w int, h int)) {
	if boxa.p == nil {
		return
	}

	for i := 0; i < boxaGetCount(boxa); i++ {
		if box := boxaGetBox(boxa, i, L_CLONE); box.p != nil {
			x, y, w, h := boxGetGeometry(box)
			draw(x, y, w, h)

			boxDestroy(&box)
		}
	}
}

////////////////////////////////////////////////////

func (pix Pix) MaskSquares(thresh int, block int, sqMin int, sqMax int) Pix {
	mask := pixConvertTo1(pix, thresh)
	if mask == NullPix {
		return NullPix
	}

	brick := pixCloseBrick(NullPix, mask, block, block)
	if brick == NullPix {
		mask.Destroy()
		return NullPix
	}

	boxes := pixConnCompBB(brick, 8)
	brick.Destroy()

	if boxes.p == nil {
		mask.Destroy()
		return NullPix
	}

	boxes.Draw(func(x int, y int, w int, h int) {
		if h == 0 {
			return
		}

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
	if pixe == NullPix {
		return NullPix
	}

	pixb := pixConvertTo1(pixe, thresh)
	pixDestroy(&pixe)
	if pixb == NullPix {
		return NullPix
	}
	pixInvert(pixb, pixb)

	brick1 := pixCloseBrick(NullPix, pixb, brMax, brMin)
	if brick1 == NullPix {
		pixDestroy(&pixb)
		return NullPix
	}

	brick2 := pixOpenBrick(NullPix, pixb, brMax, brMin)
	if brick2 == NullPix {
		pixDestroy(&brick1)
		pixDestroy(&pixb)
		return NullPix
	}

	pixXor(brick2, brick2, brick1)
	pixDestroy(&brick1)
	pixOpenBrick(brick2, brick2, brWidth, brHeight)

	brick1 = pixCloseBrick(NullPix, pixb, brMin, brMax)
	if brick1 == NullPix {
		pixDestroy(&brick2)
		pixDestroy(&pixb)
		return NullPix
	}

	mask := pixOpenBrick(NullPix, pixb, brMin, brMax)
	if mask == NullPix {
		pixDestroy(&brick1)
		pixDestroy(&brick2)
		pixDestroy(&pixb)
		return NullPix
	}

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
	if pixb == NullPix {
		return NullPix
	}

	brick1 := pixOpenBrick(NullPix, pixb, lenMin, 1)
	if brick1 == NullPix {
		pixb.Destroy()
		return NullPix
	}

	pixDilateBrick(brick1, brick1, 1, 3)

	brick2 := pixOpenBrick(NullPix, pixb, 1, lenMin)
	if brick2 == NullPix {
		pixDestroy(&brick1)
		pixb.Destroy()
		return NullPix
	}

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
	if mask == NullPix {
		return NullPix
	}

	width, height, _ := pixGetDimensions(mask)
	if width == 0 || height == 0 {
		pixDestroy(&mask)
		return NullPix
	}

	grid_x := (width + 9) / 10
	grid_y := (height + 9) / 10
	if grid_x == 0 || grid_y == 0 {
		pixDestroy(&mask)
		return NullPix
	}
	grid := [11][11][]int{}

	brick := pixCloseBrick(NullPix, mask, 3, 3)
	if brick == NullPix {
		pixDestroy(&mask)
		return NullPix
	}

	boxes := pixConnCompBB(brick, 8)
	pixDestroy(&brick)
	if boxes.p == nil {
		pixDestroy(&mask)
		return NullPix
	}

	n := boxaGetCount(boxes)
	speck := make([]bool, n)
	weights := make([]int, n)

	for i := 0; i < n; i++ {
		box := boxaGetBox(boxes, i, L_CLONE)
		if box.p == nil {
			continue
		}

		x, y, w, h := boxGetGeometry(box)
		if w <= 0 || h <= 0 {
			boxDestroy(&box)
			continue
		}

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
		if box_i.p == nil {
			continue
		}

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
				if box_j.p == nil {
					continue
				}

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
	if opt.Thresh == 0 {
		return
	}

	mask1 := pix.MaskSquares(opt.Thresh, opt.SqrBlock, opt.SqrMin, opt.SqrMax)
	if mask1 != NullPix {
		pix.pixSetMasked(mask1, 0xFFFFFFFF)
		//mask1.WriteToFile("b_mask_1.png", IFF_PNG)
		mask1.Destroy()
	}
	mask2 := pix.MaskBars(opt.Thresh, opt.BarMin, opt.BarMax, opt.BarWidth, opt.BarHeight)
	if mask2 != NullPix {
		pix.pixSetMasked(mask2, 0xFFFFFFFF)
		//mask2.WriteToFile("b_mask_2.png", IFF_PNG)
		mask2.Destroy()
	}
	mask3 := pix.MaskLines(opt.Thresh, opt.LinMin)
	if mask3 != NullPix {
		pix.pixSetMasked(mask3, 0xFFFFFFFF)
		//mask3.WriteToFile("b_mask_3.png", IFF_PNG)
		mask3.Destroy()
	}
	mask4 := pix.MaskSpecks(opt.Thresh, opt.SpMax, opt.SpWeight)
	if mask4 != NullPix {
		pix.pixSetMasked(mask4, 0xFFFFFFFF)
		//mask4.WriteToFile("b_mask_4.png", IFF_PNG)
		mask4.Destroy()
	}
}
