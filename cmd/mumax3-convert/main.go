/*
mumax3-convert converts mumax3 output files to various formats and images.
It also provides basic manipulations like data rescale etc.


Usage

Command-line flags must always preceed the input files:
	mumax3-convert [flags] files
For a overview of flags, run:
	mumax3-convert -help
Example: convert all .ovf files to PNG:
	mumax3-convert -png *.ovf
Example: resize data to a 32 x 32 x 1 mesh, normalize vectors to unit length and convert the result to OOMMF binary output:
	mumax3-convert -resize 32x32x1 -normalize -ovf binary file.ovf
Example: convert all .ovf files to VTK binary saving only the X component. Also output to JPEG in the meanwhile:
	mumax3-convert -comp 0 -vtk binary -jpg *.ovf
Example: convert legacy .dump files to .ovf:
	mumax3-convert -ovf2 *.dump
Example: cut out a piece of the data between min:max. max is exclusive bound. bounds can be omitted, default to 0 lower bound or maximum upper bound
	mumax3-convert -xrange 50:100 -yrange :100 file.ovf
Example: select the bottom layer
	mumax3-convert -zrange :1 file.ovf

Output file names are automatically assigned.
*/
package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/mumax/3/data"
	"github.com/mumax/3/draw"
	"github.com/mumax/3/dump"
	"github.com/mumax/3/oommf"
	"github.com/mumax/3/util"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	flag_comp      = flag.String("comp", "", "Select a component of vector data. (0,1,2 or x,y,z)")
	flag_show      = flag.Bool("show", false, "Human-readible output to stdout")
	flag_format    = flag.String("f", "%v", "Printf format string")
	flag_png       = flag.Bool("png", false, "PNG output")
	flag_jpeg      = flag.Bool("jpg", false, "JPEG output")
	flag_gif       = flag.Bool("gif", false, "GIF output")
	flag_svg       = flag.Bool("svg", false, "SVG output")
	flag_svgz      = flag.Bool("svgz", false, "SVGZ output (compressed)")
	flag_gnuplot   = flag.Bool("gplot", false, "Gnuplot-compatible output")
	flag_ovf1      = flag.String("ovf", "", `"text" or "binary" OVF1 output`)
	flag_omf       = flag.String("omf", "", `"text" or "binary" OVF1 output`)
	flag_ovf2      = flag.String("ovf2", "", `"text" or "binary" OVF2 output`)
	flag_vtk       = flag.String("vtk", "", `"ascii" or "binary" VTK output`)
	flag_dump      = flag.Bool("dump", false, `output in dump format`)
	flag_csv       = flag.Bool("csv", false, `output in CSV format`)
	flag_json      = flag.Bool("json", false, `output in JSON format`)
	flag_min       = flag.String("min", "auto", `Minimum of color scale: "auto" or value.`)
	flag_max       = flag.String("max", "auto", `Maximum of color scale: "auto" or value.`)
	flag_normalize = flag.Bool("normalize", false, `Normalize vector data to unit length`)
	flag_normpeak  = flag.Bool("normpeak", false, `Scale vector data, maximum to unit length`)
	flag_resize    = flag.String("resize", "", "Resize. E.g.: 128x128x4")
	flag_cropx     = flag.String("xrange", "", "Crop x range min:max (both optional, max=exclusive)")
	flag_cropy     = flag.String("yrange", "", "Crop y range min:max (both optional, max=exclusive)")
	flag_cropz     = flag.String("zrange", "", "Crop z range min:max (both optional, max=exclusive)")
	flag_dir       = flag.String("o", "", "Save all output in this directory")
	flag_arrows    = flag.Int("arrows", 0, "Arrow size for vector bitmap image output")
)

var que chan task
var wg sync.WaitGroup

type task struct {
	*data.Slice
	info  data.Meta
	fname string
}

func main() {
	log.SetFlags(0)
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("no input files")
	}

	// start many worker goroutines taking tasks from que
	runtime.GOMAXPROCS(runtime.NumCPU())
	ncpu := runtime.GOMAXPROCS(-1)
	que = make(chan task, ncpu)
	if ncpu == 0 {
		ncpu = 1
	}
	for i := 0; i < ncpu; i++ {
		go work()
	}

	// politely try to make the output directory
	if *flag_dir != "" {
		_ = os.Mkdir(*flag_dir, 0777)
	}

	// read all input files and put them in the task que
	for _, fname := range flag.Args() {
		log.Println(fname)

		var slice *data.Slice
		var info data.Meta
		var err error

		switch path.Ext(fname) {
		default:
			log.Println("skipping unsupported type", path.Ext(fname))
			continue
		case ".ovf", ".omf", ".ovf2":
			slice, info, err = oommf.ReadFile(fname)
		case ".dump":
			slice, info, err = dump.ReadFile(fname)
		}

		if err != nil {
			log.Println(err)
			continue
		}
		wg.Add(1)
		outfname := util.NoExt(fname)
		if *flag_dir != "" {
			outfname = *flag_dir + "/" + path.Base(outfname)
		}
		que <- task{slice, info, outfname}
	}

	// wait for work to finish
	wg.Wait()
}

func work() {
	for task := range que {
		process(task.Slice, task.info, task.fname)
		wg.Done()
	}
}

func open(fname string) (*os.File, io.Writer) {
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	util.FatalErr(err)
	return f, bufio.NewWriter(f)
}

func process(f *data.Slice, info data.Meta, name string) {
	preprocess(f)

	haveOutput := false

	if *flag_jpeg {
		draw.RenderFile(name+".jpg", f, *flag_min, *flag_max, *flag_arrows)
		haveOutput = true
	}

	if *flag_png {
		draw.RenderFile(name+".png", f, *flag_min, *flag_max, *flag_arrows)
		haveOutput = true
	}

	if *flag_gif {
		draw.RenderFile(name+".gif", f, *flag_min, *flag_max, *flag_arrows)
		haveOutput = true
	}

	if *flag_svg {
		out, bufout := open(name + ".svg")
		defer out.Close()
		draw.SVG(bufout, f.Vectors())
		haveOutput = true
	}

	if *flag_svgz {
		out1, _ := open(name + ".svgz")
		defer out1.Close()
		out2 := gzip.NewWriter(out1)
		defer out2.Close()
		draw.SVG(out2, f.Vectors())
		haveOutput = true
	}

	if *flag_gnuplot {
		out, bufout := open(name + ".gplot")
		defer out.Close()
		dumpGnuplot(bufout, f, info)
		haveOutput = true
	}

	if *flag_ovf1 != "" {
		out, bufout := open(name + ".ovf")
		defer out.Close()
		oommf.WriteOVF1(bufout, f, info, *flag_ovf1)
		haveOutput = true
	}

	if *flag_omf != "" {
		out, bufout := open(name + ".omf")
		defer out.Close()
		oommf.WriteOVF1(bufout, f, info, *flag_omf)
		haveOutput = true
	}

	if *flag_ovf2 != "" {
		out, bufout := open(name + ".ovf")
		defer out.Close()
		oommf.WriteOVF2(bufout, f, info, *flag_ovf2)
		haveOutput = true
	}

	if *flag_vtk != "" {
		out, bufout := open(name + ".vts") // vts is the official extension for VTK files containing StructuredGrid data
		defer out.Close()
		dumpVTK(bufout, f, info, *flag_vtk)
		haveOutput = true
	}

	if *flag_csv {
		out, bufout := open(name + ".csv")
		defer out.Close()
		dumpCSV(bufout, f)
		haveOutput = true
	}

	if *flag_json {
		out, bufout := open(name + ".json")
		defer out.Close()
		dumpJSON(bufout, f)
		haveOutput = true
	}

	if *flag_dump {
		dump.MustWriteFile(name+".dump", f, info)
		haveOutput = true
	}

	if !haveOutput || *flag_show {
		fmt.Println(info)
		util.Fprintf(os.Stdout, *flag_format, f.Tensors())
		haveOutput = true
	}

}

func preprocess(f *data.Slice) {
	if *flag_normalize {
		normalize(f, 1)
	}
	if *flag_normpeak {
		normpeak(f)
	}
	if *flag_comp != "" {
		*f = *f.Comp(parseComp(*flag_comp))
	}
	crop(f)
	if *flag_resize != "" {
		resize(f, *flag_resize)
	}
}

func parseComp(c string) int {
	if i, err := strconv.Atoi(c); err == nil {
		return i
	}
	switch c {
	default:
		log.Fatal("illegal component:", c, "(need x, y or z)")
		panic(0)
	case "x", "X":
		return 0
	case "y", "Y":
		return 1
	case "z", "Z":
		return 2
	}
}

func crop(f *data.Slice) {
	N := f.Size()
	// default ranges
	x1, x2 := 0, N[X]
	y1, y2 := 0, N[Y]
	z1, z2 := 0, N[Z]
	havework := false

	if *flag_cropz != "" {
		z1, z2 = parseRange(*flag_cropz, N[Z])
		havework = true
	}
	if *flag_cropy != "" {
		y1, y2 = parseRange(*flag_cropy, N[Y])
		havework = true
	}
	if *flag_cropx != "" {
		x1, x2 = parseRange(*flag_cropx, N[X])
		havework = true
	}

	if havework {
		*f = *data.Crop(f, x1, x2, y1, y2, z1, z2)
	}
}

func parseRange(r string, max int) (int, int) {
	a, b := 0, max
	spl := strings.Split(r, ":")
	if len(spl) != 2 {
		log.Fatal("range needs min:max syntax, have:", r)
	}
	if spl[0] != "" {
		a = atoi(spl[0])
	}
	if spl[1] != "" {
		b = atoi(spl[1])
	}
	return a, b
}

func atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		panic(err)
	}
	return i
}

const (
	X = data.X
	Y = data.Y
	Z = data.Z
)
