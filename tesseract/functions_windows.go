package tesseract

// #include <stdlib.h>
import "C"
import "syscall"

// Function signatures from Tesseract DLL
var (
	tessDll = syscall.NewLazyDLL("tesseract50.dll")

	tessCreate         = tessDll.NewProc("TessBaseAPICreate")
	tessBaseAPIEnd     = tessDll.NewProc("TessBaseAPIEnd")
	tessDelete         = tessDll.NewProc("TessBaseAPIDelete")
	tessInit2          = tessDll.NewProc("TessBaseAPIInit2")
	tessSetPageSegMode = tessDll.NewProc("TessBaseAPISetPageSegMode")
	tessSetImage2      = tessDll.NewProc("TessBaseAPISetImage2")

	/* Utility */
	tessSetVariable  = tessDll.NewProc("TessBaseAPISetVariable")
	tessFreeUTF8Text = tessDll.NewProc("TessDeleteText")

	/* Whole text */
	tessRecognize   = tessDll.NewProc("TessBaseAPIRecognize")
	tessGetUTF8Text = tessDll.NewProc("TessBaseAPIGetUTF8Text")
	tessGetHOCRText = tessDll.NewProc("TessBaseAPIGetHOCRText")

	/* Result iterator */
	tessGetIterator                   = tessDll.NewProc("TessBaseAPIGetIterator")            //TessResultIterator *(TessBaseAPI *handle)
	tessResultIteratorGetPageIterator = tessDll.NewProc("TessResultIteratorGetPageIterator") //TessPageIterator *(TessResultIterator *handle);

	tessResultIteratorDelete      = tessDll.NewProc("TessResultIteratorDelete")      // void(TessResultIterator *handle);
	tessResultIteratorNext        = tessDll.NewProc("TessResultIteratorNext")        // BOOL(TessResultIterator *handle, TessPageIteratorLevel level)
	tessResultIteratorGetUTF8Text = tessDll.NewProc("TessResultIteratorGetUTF8Text") // char*(const TessResultIterator *handle, TessPageIteratorLevel level)
	tessResultIteratorConfidence  = tessDll.NewProc("TessResultIteratorConfidence")  // float(const TessResultIterator *handle, TessPageIteratorLevel level)

	/* Page iterator */
	tessPageIteratorDelete           = tessDll.NewProc("TessPageIteratorDelete")           // void(TessPageIterator *handle)
	tessPageIteratorBegin            = tessDll.NewProc("TessPageIteratorBegin")            // void(TessPageIterator *handle)
	tessPageIteratorNext             = tessDll.NewProc("TessPageIteratorNext")             // BOOL (TessPageIterator *handle, TessPageIteratorLevel level)
	tessPageIteratorIsAtBeginningOf  = tessDll.NewProc("TessPageIteratorIsAtBeginningOf")  // BOOL(const TessPageIterator *handle, TessPageIteratorLevel level)
	tessPageIteratorIsAtFinalElement = tessDll.NewProc("TessPageIteratorIsAtFinalElement") //BOOL(const TessPageIterator *handle, TessPageIteratorLevel level, TessPageIteratorLevel element)
	tessPageIteratorBoundingBox      = tessDll.NewProc("TessPageIteratorBoundingBox")      //BOOL(const TessPageIterator *handle, TessPageIteratorLevel level, int *left, int *top, int *right, int *bottom)
)
