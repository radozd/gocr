package main

import (
	"errors"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/radozd/gocr/leptonica"
	"golang.org/x/exp/maps"
)

var originalImage leptonica.Pix

var paramsEnhance map[string]*int = map[string]*int{
	"tilex":    &leptonica.DefaultEnhanceOptions.TileX,
	"tiley":    &leptonica.DefaultEnhanceOptions.TileY,
	"thresh1":  &leptonica.DefaultEnhanceOptions.Thresh,
	"minc":     &leptonica.DefaultEnhanceOptions.MinCount,
	"white1":   &leptonica.DefaultEnhanceOptions.WhitePoint,
	"smoothx":  &leptonica.DefaultEnhanceOptions.SmoothX,
	"smoothy":  &leptonica.DefaultEnhanceOptions.SmoothY,
	"gamma":    &leptonica.DefaultEnhanceOptions.Gamma,
	"gammamin": &leptonica.DefaultEnhanceOptions.GammaMin,
	"gammamax": &leptonica.DefaultEnhanceOptions.GammaMax,
	"factor":   &leptonica.DefaultEnhanceOptions.Factor,
	"border":   &leptonica.DefaultEnhanceOptions.RemoveBorders,
}

var paramsGray map[string]*int = map[string]*int{
	"sat":     &leptonica.DefaultGrayOptions.Saturation,
	"white2":  &leptonica.DefaultGrayOptions.WhitePoint,
	"thresh2": &leptonica.DefaultGrayOptions.ThreshDiff,
	"mindist": &leptonica.DefaultGrayOptions.MinDist,
}

var paramsMask map[string]*int = map[string]*int{
	"thresh3":  &leptonica.DefaultMaskOptions.Thresh,
	"sqblock":  &leptonica.DefaultMaskOptions.SqrBlock,
	"sqmin":    &leptonica.DefaultMaskOptions.SqrMin,
	"sqmax":    &leptonica.DefaultMaskOptions.SqrMax,
	"brmin":    &leptonica.DefaultMaskOptions.BarMin,
	"brmax":    &leptonica.DefaultMaskOptions.BarMax,
	"brwidth":  &leptonica.DefaultMaskOptions.BarWidth,
	"brheight": &leptonica.DefaultMaskOptions.BarHeight,
	"lnmin":    &leptonica.DefaultMaskOptions.LinMin,
	"spmax":    &leptonica.DefaultMaskOptions.SpMax,
	"spweight": &leptonica.DefaultMaskOptions.SpWeight,
}

func main() {
	pix := leptonica.NewPixFromFile("__e2abcbf44eb6794003774813e2bae73f.tif") //"test1.tif")
	//pix := leptonica.NewPixFromFile("test1.tif")
	defer pix.Destroy()

	originalImage = pix //pix.GetDeskewedCopy(0)

	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/image", serveImage)
	http.HandleFunc("/process", processImage)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			var err error
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			log.Print(err)
			log.Print(string(debug.Stack()[:]))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	sliderInt := func(params map[string]*int, name string, min int, max int) string {
		smin := strconv.Itoa(min)
		smax := strconv.Itoa(max)
		value := -1
		if v, ok := params[name]; ok {
			value = *v
		} else {
			panic("unknown name: " + name)
		}
		sval := strconv.Itoa(value)
		return "<div><label for='" + name + "'>" + name + ":</label>" +
			"<input type='range' id='" + name + "' min='" + smin + "' max='" + smax + "' step='1' value='" + sval + "' oninput='updateImage(); this.nextElementSibling.value=this.value'>" +
			"<output name='" + name + "_value' for='" + name + "'>" + sval + "</output>" +
			"</div>"
	}
	readVal := func(name string) string {
		return "const " + name + " = document.getElementById('" + name + "').value;"
	}
	readVals := func(names []string) string {
		sb := strings.Builder{}
		for _, n := range names {
			sb.WriteString("\n" + readVal(n))
		}
		return sb.String()
	}
	query := func(names []string) string {
		sb := strings.Builder{}
		for _, n := range names {
			sb.WriteString(n + "=${" + n + "}&")
		}
		return strings.TrimSuffix(sb.String(), "&")
	}

	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Image Adjuster</title>
		<style>
			html, body { margin: 0; padding: 0; }
			body { display: flex; flex-direction: row; background-color: lightgrey; }
			nav { margin: 1em; position: fixed; background-color: lightgrey; transform-origin: top left; }
			label { width: 5em; display: inline-block; }
			#image { margin-left: auto; margin-right: auto; max-height: 100vh; width: auto; }
		</style>

		<script>
			let _showTimer = null;

			function processedUrl() {` +
		readVals(maps.Keys(paramsEnhance)) + readVals(maps.Keys(paramsGray)) + readVals(maps.Keys(paramsMask)) +
		`return ` + "`/process?" + query(maps.Keys(paramsEnhance)) + query(maps.Keys(paramsGray)) + query(maps.Keys(paramsMask)) + "`;" +
		`}

			function updateImage() {
				if( _showTimer !== null ) {
					clearTimeout(_showTimer);
				}

				_showTimer = setTimeout(function() {
					document.getElementById('image').onmouseout();
				}, 200);
			}

			window.onload = function() {
				const image = document.getElementById('image');
				image.src = "/image";
				image.onmouseover = function() { this.src='/image'; };
				image.onmouseout = function() { this.src=processedUrl(); };

				_showTimer = setTimeout(function() {
					document.getElementById('image').onmouseout();
				}, 1000);

				window.addEventListener('scroll', function(e) {
					const el = document.getElementById('side');
					const scale = window.innerWidth/document.documentElement.clientWidth;
					el.style["transform"] = "scale(" + scale + ")";
					el.style.left = window.pageXOffset + 'px';
					el.style.top = window.pageYOffset + 'px';
				})
			}
		</script>
	</head>
	<body>
		<nav id='side'>
			<fieldset>
				<legend>Enhance</legend>` +
		sliderInt(paramsEnhance, "tilex", 0, 40) +
		sliderInt(paramsEnhance, "tiley", 0, 40) +
		sliderInt(paramsEnhance, "thresh1", 1, 128) +
		sliderInt(paramsEnhance, "minc", 5, 100) +
		sliderInt(paramsEnhance, "white1", 1, 254) +
		sliderInt(paramsEnhance, "smoothx", 0, 20) +
		sliderInt(paramsEnhance, "smoothy", 0, 20) +
		sliderInt(paramsEnhance, "gamma", 0, 100) +
		sliderInt(paramsEnhance, "gammamin", 1, 254) +
		sliderInt(paramsEnhance, "gammamax", 1, 254) +
		sliderInt(paramsEnhance, "factor", 0, 100) +
		sliderInt(paramsEnhance, "border", 1, 255) + `
			</fieldset>
			<fieldset>
				<legend>Gray</legend>` +
		sliderInt(paramsGray, "sat", 1, 255) +
		sliderInt(paramsGray, "white2", 1, 254) +
		sliderInt(paramsGray, "thresh2", 1, 128) +
		sliderInt(paramsGray, "mindist", 1, 20) + `
			</fieldset>
			<fieldset>
				<legend>Mask</legend>` +
		sliderInt(paramsMask, "thresh3", 0, 128) +
		sliderInt(paramsMask, "sqblock", 1, 30) +
		sliderInt(paramsMask, "sqmin", 30, 150) +
		sliderInt(paramsMask, "sqmax", 160, 500) +
		sliderInt(paramsMask, "brmin", 1, 10) +
		sliderInt(paramsMask, "brmax", 10, 30) +
		sliderInt(paramsMask, "brwidth", 20, 80) +
		sliderInt(paramsMask, "brheight", 20, 80) +
		sliderInt(paramsMask, "lnmin", 20, 200) +
		sliderInt(paramsMask, "spmax", 2, 20) +
		sliderInt(paramsMask, "spweight", 1, 200) + `
			</fieldset>
		</nav>
		<img id="image" src="/image" alt="Predefined Image">
	</body>
	</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func pixToJpeg(pix leptonica.Pix) []byte {
	var res []byte
	var err error
	if res, err = pix.WriteToMem(leptonica.JFIF_JPEG); err != nil {
		return nil
	}
	return res
}

func serveImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(pixToJpeg(originalImage))
}

func processImage(w http.ResponseWriter, r *http.Request) {
	getInt := func(params map[string]*int, name string) {
		if n, err := strconv.Atoi(r.URL.Query().Get(name)); err == nil {
			if pi, ok := params[name]; ok {
				*pi = n
			}
		}
	}
	getInts := func(params map[string]*int) {
		for _, n := range maps.Keys(params) {
			getInt(params, n)
		}
	}

	getInts(paramsEnhance)
	getInts(paramsGray)
	getInts(paramsMask)

	w.Header().Set("Content-Type", "image/jpeg")

	tmp := originalImage.Copy()
	defer tmp.Destroy()

	en := tmp.EnhancedCopy(leptonica.DefaultEnhanceOptions)
	defer en.Destroy()

	gray := en.GetGrayCopy(leptonica.GRAY_CAST_REMOVE_COLORS, leptonica.DefaultGrayOptions)
	defer gray.Destroy()

	deskew, _ := gray.GetDeskewedCopyAndAngle(0)
	defer deskew.Destroy()

	deskew.MaskAll(leptonica.DefaultMaskOptions)
	w.Write(pixToJpeg(deskew))
}
