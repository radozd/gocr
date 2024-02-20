//go:build windows

package zbar

// #include <stdlib.h>
import "C"
import "unsafe"

type symbol uintptr

func (img image) first() symbol {
	sym, _, _ := zbar_image_first_symbol.Call(uintptr(img))
	return symbol(sym)
}

func (sym symbol) next() symbol {
	n, _, _ := zbar_symbol_next.Call(uintptr(sym))
	return symbol(n)
}

func (sym symbol) string() string {
	text, _, _ := zbar_symbol_get_data.Call(uintptr(sym))
	if text == 0 {
		return ""
	}
	return C.GoString((*C.char)(unsafe.Pointer(text)))
}

func (sym symbol) get_loc_size() int {
	res, _, _ := zbar_symbol_get_loc_size.Call(uintptr(sym))
	return int(res)
}

func (sym symbol) get_loc_x(idx int) int {
	res, _, _ := zbar_symbol_get_loc_x.Call(uintptr(sym), uintptr(idx))
	return int(res)
}

func (sym symbol) get_loc_y(idx int) int {
	res, _, _ := zbar_symbol_get_loc_y.Call(uintptr(sym), uintptr(idx))
	return int(res)
}

func (sym symbol) symbol_type() ZBAR_CODETYPE {
	t, _, _ := zbar_symbol_get_type.Call(uintptr(sym))
	return ZBAR_CODETYPE(t)
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
	for s != 0 {
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
