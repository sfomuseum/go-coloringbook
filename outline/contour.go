package outline

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"image"
	"strings"

	"github.com/fogleman/colormap"
	"github.com/fogleman/contourmap"
	"github.com/fogleman/gg"
)

type ContourFormat uint8

const (
	PNG ContourFormat = iota
	SVG
)

type ContourOptions struct {
	Iterations int
	Scale      float64
	Format     string
}

func Contour(ctx context.Context, im image.Image, opts *ContourOptions) (Outline, error) {

	switch strings.ToUpper(opts.Format) {
	case "PNG":
		return ContourPNG(ctx, im, opts)
	case "SVG":
		return ContourSVG(ctx, im, opts)
	default:
		return nil, fmt.Errorf("Invalid contour format")
	}
}

func ContourSVG(ctx context.Context, im image.Image, opts *ContourOptions) (Outline, error) {

	iterations := opts.Iterations
	scale := opts.Scale

	m := contourmap.FromImage(im).Closed()
	z0 := m.Min
	z1 := m.Max

	w := int(float64(m.W) * scale)
	h := int(float64(m.H) * scale)

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	fmt.Fprintf(wr, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, w, h, w, h)

	for i := 0; i < iterations; i++ {

		t := float64(i) / (float64(iterations) - 1)
		z := z0 + (z1-z0)*t
		contours := m.Contours(z + 1e-9)

		z = z * float64(i)

		for _, c := range contours {

			fmt.Fprintf(wr, `<path stroke="%s" stroke-width="%02f" stroke-opacity="1" fill-opacity="0" d="M`, "#000000", z)

			for i, p := range c {

				if i > 0 {
					fmt.Fprintf(wr, `L`)
				}

				fmt.Fprintf(wr, `%d,%d`, int(p.X), int(p.Y))
			}

			fmt.Fprintf(wr, `Z"></path>`)
		}
	}

	fmt.Fprintf(wr, `</svg>`)

	o := &SVGOutline{
		svg: buf.Bytes(),
	}

	return o, nil
}

func ContourImage(ctx context.Context, im image.Image, opts *ContourOptions) (image.Image, error) {

	iterations := opts.Iterations
	scale := opts.Scale

	m := contourmap.FromImage(im).Closed()
	z0 := m.Min
	z1 := m.Max

	w := int(float64(m.W) * scale)
	h := int(float64(m.H) * scale)

	dc := gg.NewContext(w, h)
	dc.SetRGB(1, 1, 1)
	dc.SetColor(colormap.ParseColor("FFFFFF"))
	dc.Clear()
	dc.Scale(scale, scale)

	for i := 0; i < iterations; i++ {

		t := float64(i) / (float64(iterations) - 1)
		z := z0 + (z1-z0)*t
		contours := m.Contours(z + 1e-9)

		// Do line smoothing here?

		for _, c := range contours {

			dc.NewSubPath()

			for _, p := range c {
				dc.LineTo(p.X, p.Y)
			}
		}

		dc.SetRGB(0, 0, 0)

		// z = 1.0
		z = z * float64(i)

		dc.SetLineWidth(z)
		dc.Stroke()
	}

	return dc.Image(), nil
}

func ContourPNG(ctx context.Context, im image.Image, opts *ContourOptions) (Outline, error) {

	new_im, err := ContourImage(ctx, im, opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive raster, %w", err)
	}

	o := &PNGOutline{
		image: new_im,
	}

	return o, nil
}
