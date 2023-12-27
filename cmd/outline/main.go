package main

import (
	"context"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/sfomuseum/go-coloringbook/outline"
	"github.com/sfomuseum/go-flags/flagset"
)

func main() {

	var contour_iterations int
	var contour_scale float64
	var contour_format string
	var contour_smoothing bool
	var smoothing_close bool
	var smoothing_iterations int

	var vtracer_precision int
	var vtracer_speckle int

	var use_batik bool
	var path_batik string

	var infile string
	var outfile string

	fs := flagset.NewFlagSet("coloringbook")

	fs.IntVar(&contour_iterations, "contour-iteration", 8, "The number of iterations to perform generating an image contour.")
	fs.Float64Var(&contour_scale, "contour-scale", 1.0, "The scale factor of the final image contour.")
	fs.StringVar(&contour_format, "contour-format", "png", "The format for the final image contour. Valid options are: png, svg.")
	fs.BoolVar(&contour_smoothing, "contour-smoothing", false, "...")
	fs.IntVar(&smoothing_iterations, "contour-smoothing-iterations", 6, "...")
	fs.BoolVar(&smoothing_close, "contour-smoothing-close", true, "...")

	fs.IntVar(&vtracer_precision, "vtracer-precision", 6, "Number of significant bits (color precision) to use in an RGB channel.")
	fs.IntVar(&vtracer_speckle, "vtracer-speckle", 8, "Discard patches smaller than X px in size")

	fs.BoolVar(&use_batik, "use-batik", true, "Use the Java Batik SVG raterizer.")
	fs.StringVar(&path_batik, "path-batik", "/usr/local/src/batik-1.17/batik-rasterizer-1.17.jar", "The path to the Java Batik SVG raterizer JAR file.")

	fs.StringVar(&infile, "infile", "", "The path to the image you want to generate an outline for.")
	fs.StringVar(&outfile, "outfile", "", "The path to the final image that has been outlined.")

	flagset.Parse(fs)

	ctx := context.Background()

	contour_opts := &outline.ContourOptions{
		Iterations:          contour_iterations,
		Scale:               contour_scale,
		Format:              contour_format,
		Smoothing:           contour_smoothing,
		SmoothingClose:      smoothing_close,
		SmoothingIterations: smoothing_iterations,
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
