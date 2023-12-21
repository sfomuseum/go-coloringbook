package outline

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"
)

type Outline interface {
	Write(context.Context, io.Writer) error
}

type PNGOutline struct {
	Outline
	image image.Image
}

func (o *PNGOutline) Write(ctx context.Context, wr io.Writer) error {
	return png.Encode(wr, o.image)
}

type SVGOutline struct {
	Outline
	svg []byte
}

func (o *SVGOutline) Write(ctx context.Context, wr io.Writer) error {
	_, err := wr.Write(o.svg)
	return err
}

type OutlineOptions struct {
	Contour   *ContourOptions
	Trace     *TraceOptions
	Rasterize *RasterizeOptions
}

func GenerateOutline(ctx context.Context, im image.Image, opts *OutlineOptions) (Outline, error) {

	vtrace_infile, err := os.CreateTemp("", "vtrace.*.png")

	if err != nil {
		return nil, fmt.Errorf("Failed to create vtrace input file, %w", err)
	}

	infile_uri := vtrace_infile.Name()
	defer os.Remove(infile_uri)

	err = png.Encode(vtrace_infile, im)

	if err != nil {
		return nil, fmt.Errorf("Failed to encode image for tracing, %w", err)
	}

	err = vtrace_infile.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close infile after writing, %w", err)
	}

	vtrace_outfile, err := os.CreateTemp("", "vtrace.*.svg")

	if err != nil {
		return nil, fmt.Errorf("Failed to create vtrace outfile file, %w", err)
	}

	err = vtrace_outfile.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close outfile, %w", err)
	}

	outfile_uri := vtrace_outfile.Name()
	defer os.Remove(outfile_uri)

	log.Println("TRACE")

	traced_im, err := Trace(ctx, infile_uri, outfile_uri, opts.Trace, opts.Rasterize)

	if err != nil {
		return nil, fmt.Errorf("Failed to trace image, %w", err)
	}

	log.Println("CONTOUR")

	o, err := Contour(ctx, traced_im, opts.Contour)

	if err != nil {
		return nil, fmt.Errorf("Failed to contour image, %w", err)
	}

	return o, nil
}
