package zbar

// #include <stdlib.h>
import "C"

func (img image) first() symbol {
	return zbar_image_first_symbol(img)
}

func (sym symbol) next() symbol {
	return zbar_symbol_next(sym)
}

func (sym symbol) string() string {
	return zbar_symbol_get_data(sym)
}

func (sym symbol) get_loc_size() int {
	return zbar_symbol_get_loc_size(sym)
}

func (sym symbol) get_loc_x(idx int) int {
	return zbar_symbol_get_loc_x(sym, idx)
}

func (sym symbol) get_loc_y(idx int) int {
	return zbar_symbol_get_loc_y(sym, idx)
}

func (sym symbol) symbol_type() ZBAR_CODETYPE {
	return zbar_symbol_get_type(sym)
}

func (sym symbol) get_rect() (int, int, int, int) {
	x1 := 1000000
	y1 := 1000000
	x2 := 0
	y2 := 0
	for i := 0; i < sym.get_loc_size(); i++ {
		x := sym.get_loc_x(i)
		y := sym.get_loc_y(i)
		if x < x1 {
			x1 = x
		}
		if x > x2 {
			x2 = x
		}
		if y < y1 {
			y1 = y
		}
		if y > y2 {
			y2 = y
		}
	}
	return x1, y1, x2, y2
}

func (sym symbol) each(f func(Code)) {
	s := sym
	for s != nullSymbol {
		x1, y1, x2, y2 := s.get_rect()

		code := Code{
			CodeType: s.symbol_type(),
			Value:    s.string(),
			X:        x1,
			Y:        y1,
			Width:    x2 - x1,
			Height:   y2 - y1,
		}
		f(code)
		s = s.next()
	}
}
