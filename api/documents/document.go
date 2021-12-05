package documents

import (
	"context"
	"io"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/a-h/templ"
)

type TemplateData struct {
	Name string
}

func ToPDF(ctx context.Context, component templ.Component, w io.Writer) error {
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		return err
	}

	pdfg.SetOutput(w)
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtml.OrientationPortrait)
	pdfg.Grayscale.Set(false)

	cr, cw := io.Pipe()
	var cerr error
	go func() {
		cerr = component.Render(ctx, cw)
		cw.Close()
	}()

	page := wkhtml.NewPageReader(cr)
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	pdfg.AddPage(page)

	if err = pdfg.Create(); err != nil {
		return err
	}
	return cerr
}
