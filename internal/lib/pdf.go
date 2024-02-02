package lib

import (
	"context"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func htmlToPdf(html []byte, outputPath string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	tmpHTML, err := os.CreateTemp(cwd, "temp_*.html")
	if err != nil {
		return err
	}
	defer tmpHTML.Close()
	defer os.Remove(tmpHTML.Name())

	_, err = tmpHTML.Write(html)
	if err != nil {
		return err
	}
	defer tmpHTML.Close()

	var buf []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate("file:///"+tmpHTML.Name()),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().
				WithDisplayHeaderFooter(false).
				WithLandscape(false).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return err
	}

	err = os.WriteFile(outputPath, buf, 0o644)
	if err != nil {
		return err
	}

	return nil
}
