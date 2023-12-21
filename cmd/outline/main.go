package main

import (
	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-coloringbook/outline"
)

func main() {

	var contour_iterations int
	var contour_scale float64
	var contour_format string

	var vtracer_precision int
	var vtracer_speckle int

	var use_batik bool
	var path_batik string

	var infile string
	var outfile string

	fs := flagset.NewFlagSet("coloringbook")

	fs.IntVar(&contour_iterations, "contour-iteration", 8, "...")
	fs.Float64Var(&contour_scale, "contour-scale", 1.0, "...")
	fs.StringVar(&contour_format, "contour-format", "png", "...")

	fs.IntVar(&vtracer_precision, "vtracer-precision", 6, "...")
	fs.IntVar(&vtracer_speckle, "vtracer-speckle", 8, "...")

	fs.BoolVar(&use_batik, "use-batik", true, "...")
	fs.StringVar(&path_batik, "path-batik", "/usr/local/src/batik-1.17/batik-rasterizer-1.17.jar", "...")

	fs.StringVar(&infile, "infile", "", "...")
	fs.StringVar(&outfile, "outfile", "", "...")

	flagset.Parse(fs)

	ctx := context.Background()

	contour_opts := &outline.ContourOptions{
		Iterations: contour_iterations,
		Scale:      contour_scale,
		Format:     contour_format,
	}

	trace_opts := &outline.TraceOptions{
		Precision: vtracer_precision,
		Speckle:   vtracer_speckle,
	}

	raster_opts := &outline.RasterizeOptions{
		UseBatik:  use_batik,
		PathBatik: path_batik,
	}

	outline_opts := &outline.OutlineOptions{
		Contour:   contour_opts,
		Trace:     trace_opts,
		Rasterize: raster_opts,
	}

	r, err := os.Open(infile)

	if err != nil {
		log.Fatalf("Failed to open %s for reading, %v", infile, err)
	}

	defer r.Close()

	im, _, err := image.Decode(r)

	if err != nil {
		log.Fatalf("Failed to decode %s, %v", infile, err)
	}

	outline, err := outline.GenerateOutline(ctx, im, outline_opts)

	if err != nil {
		log.Fatalf("Failed to generate outline, %v", err)
	}

	wr, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("Failed to open %s for writing, %v", outfile, err)
	}

	err = outline.Write(ctx, wr)

	if err != nil {
		log.Fatalf("Failed to encode %s, %v", outfile, err)
	}

	err = wr.Close()

	if err != nil {
		log.Fatalf("Failed to close %s after writing, %v", outfile, err)
	}
}
