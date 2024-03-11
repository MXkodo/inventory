package qrcode

import (
	"os"

	qrcode "github.com/skip2/go-qrcode"
)

func Generate() {
	 var png []byte
  	 png, err := qrcode.Encode("https://astu.org", qrcode.Medium, 256)
	 if err != nil {
		return
	 }
	 file,_:= os.Create("qrcode.png")
	 defer file.Close()
	 file.Write(png)

}