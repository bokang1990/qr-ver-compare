package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"net/http"

	"github.com/skip2/go-qrcode"
)

func generateQRCode(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	if text == "" {
		http.Error(w, "Text query parameter is required", http.StatusBadRequest)
		return
	}

	qrCode1, err := qrcode.New(text, qrcode.Low)
	qrCode2, err := qrcode.New(text, qrcode.Medium)
	qrCode3, err := qrcode.New(text, qrcode.High)
	qrCode4, err := qrcode.New(text, qrcode.Highest)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}
	a := imerge(qrCode1.Image(256), qrCode2.Image(256))
	b := imerge(qrCode3.Image(256), qrCode4.Image(256))

	final := imerge(a, b)

	w.Header().Set("Content-Type", "image/png")
	err = png.Encode(w, final)
	if err != nil {
		http.Error(w, "Failed to write image", http.StatusInternalServerError)
		return
	}

	fmt.Println(len(text), " | ", qrCode1.VersionNumber, ", ", qrCode2.VersionNumber, ", ", qrCode3.VersionNumber, ", ", qrCode4.VersionNumber)
}

func main() {
	http.HandleFunc("/generate", generateQRCode)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func imerge(img1 image.Image, img2 image.Image) image.Image {

	// Create a new image with the combined width and the max height
	outputWidth := img1.Bounds().Dx() + img2.Bounds().Dx()
	outputHeight := img1.Bounds().Dy()
	if img2.Bounds().Dy() > outputHeight {
		outputHeight = img2.Bounds().Dy()
	}

	outputImage := image.NewRGBA(image.Rect(0, 0, outputWidth, outputHeight))

	// Draw img1 onto the output image
	draw.Draw(outputImage, img1.Bounds(), img1, image.Point{}, draw.Src)

	// Draw img2 onto the output image, offset by the width of img1
	offset := image.Pt(img1.Bounds().Dx(), 0)
	draw.Draw(outputImage, img2.Bounds().Add(offset), img2, image.Point{}, draw.Src)

	return outputImage
}
