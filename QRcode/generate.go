package qrcode

import (
	"os"

	qrcode "github.com/skip2/go-qrcode"
)

func Generate(url string) {
	 var png []byte
  	 png, err := qrcode.Encode(url, qrcode.Medium, 256)
	 if err != nil {
		return
	 }
	 file,_:= os.Create("qrcode.png")
	 defer file.Close()
	 file.Write(png)

}