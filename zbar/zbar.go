package zbar

type ZBAR_CODETYPE int

const ZBAR_QRCODE ZBAR_CODETYPE = 64
const ZBAR_CODE128 ZBAR_CODETYPE = 128

type Code struct {
	CodeType ZBAR_CODETYPE
	Value    string
	X        int
	Y        int
	Width    int
	Height   int
}
