package web

import (
	//"code.google.com/p/mx3/data"
	//"code.google.com/p/mx3/draw"
	//"code.google.com/p/mx3/engine"
	//"image/jpeg"
	//"log"
	"net/http"
	//"strings"
)

// Render image of quantity.
// Accepts url: /render/name and /render/name/component
func render(w http.ResponseWriter, r *http.Request) {
	panic("uncomment")
	/*url := r.URL.Path[len("/render/"):]
	words := strings.Split(url, "/")
	quant := words[0]
	comp := ""
	if len(words) > 1 {
		comp = words[1]
	}
	h, ok := engine.Quants[quant]
	if !ok {
		err := "render: unknown quantity: " + url
		log.Println(err)
		http.Error(w, err, http.StatusNotFound)
		return
	} else {
		var d *data.Slice
		// TODO: could download only needed component
		engine.InjectAndWait(func() { d = h.Download() })
		if comp != "" && d.NComp() > 1 {
			c := compstr[comp]
			d = d.Comp(c)
		}
		img := draw.Image(d, "auto", "auto")
		jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
	}
	*/
}

var compstr = map[string]int{"x": 2, "y": 1, "z": 0} // also swaps XYZ user space
