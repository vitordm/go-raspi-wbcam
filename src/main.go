package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vitordm/go-raspi-wbcam/src/utils"
	"gocv.io/x/gocv"
)

var (
	subjectFrameWebcam = utils.ChainByteSubject{}
	upgrader           = websocket.Upgrader{}
)

func startReadWebcam(webcam *gocv.VideoCapture) {
	for {
		frame := gocv.NewMat()
		webcam.Read(&frame)
		// Converte o frame para JPEG
		buf, _ := gocv.IMEncode(".jpg", frame)

		subjectFrameWebcam.Notify(buf.GetBytes())

		time.Sleep(5 * time.Second)
	}
}

func main() {
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		panic(err)
	}
	defer webcam.Close()

	go startReadWebcam(webcam)

	router := gin.Default()
	router.Static("/teste-video", "../public")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/video", func(ctx *gin.Context) {

		conn, _ := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		defer conn.Close()

		observer := make(chan []byte)

		subjectFrameWebcam.Attach(observer)

		// Loop infinito para enviar frames em tempo real
		for {

			select {
			case x := <-ctx.Done():
				fmt.Println(x)
				fmt.Println("Conexão perdida")
				subjectFrameWebcam.Detach(observer)
				return
			case frame := <-observer:
				if err := conn.WriteMessage(websocket.BinaryMessage, frame); err != nil {
					return
				}

			}
		}

	})

	router.Run(":8080")

}
