package svg

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func RasterizeAsPNG(ctx context.Context, r io.Reader, wr io.Writer) error {

	im, err := Rasterize(ctx, r)

	if err != nil {
		fmt.Errorf("Failed to rasterize image, %w", err)
	}

	err = png.Encode(wr, im)

	if err != nil {
		return fmt.Errorf("Failed to encode image, %w", err)
	}

	return nil
}

func Rasterize(ctx context.Context, r io.Reader) (image.Image, error) {

	icon, err := oksvg.ReadIconStream(r, oksvg.StrictErrorMode)

	if err != nil {
		return nil, fmt.Errorf("Failed to read stream, %w", err)
	}

	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	scanner := rasterx.NewScannerGV(w, h, img, img.Bounds())
	raster := rasterx.NewDasher(w, h, scanner)

	icon.Draw(raster, 1.0)

	return img, nil

}
