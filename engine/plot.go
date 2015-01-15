package engine

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/mumax/3/httpfs"
)

func (g *guistate) servePlot(w http.ResponseWriter, r *http.Request) {
	a := g.StringValue("usingx")
	b := g.StringValue("usingy")

	cmd := "gnuplot"
	args := []string{"-e", fmt.Sprintf(`set format x "%%g"; set key off; set format y "%%g"; set term svg size 480,320 fsize 10; plot "-" u %v:%v w li; set output;exit;`, a, b)}
	excmd := exec.Command(cmd, args...)
	stdin, _ := excmd.StdinPipe()
	stdout, _ := excmd.StdoutPipe()
	data, _ := httpfs.Read(fmt.Sprintf(`%vtable.txt`, OD()))
	err := excmd.Start()
	defer excmd.Wait()
	stdin.Write(data)
	stdin.Close()
	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		w.Write(emptyIMG())
		g.Set("plotErr", string(out))
		return
	} else {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(out)
		g.Set("plotErr", "")
	}
}

var empty_img []byte

// empty image to show if there's no plot...
func emptyIMG() []byte {
	if empty_img == nil {
		o := bytes.NewBuffer(nil)
		png.Encode(o, image.NewNRGBA(image.Rect(0, 0, 4, 4)))
		empty_img = o.Bytes()
	}
	return empty_img
}
