package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kartikkala/drive_svc/config"
	"github.com/Kartikkala/drive_svc/drive"
	"github.com/nats-io/nats.go"
)

func main() {
	app, err := config.NewApp()
	if err != nil {
		fmt.Println(err.Error())
	}

	nc, err := nats.Connect(app.Cfg.NATS.URL)

	if err != nil {
		log.Println("Error in NATS server connection...", err)
		return
	}
	driveSvc := drive.NewDriveService(app.DB)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	drive.AttachEvents(ctx, nc, driveSvc)
	log.Println("Drive service active and listening...")
	<-ctx.Done()
}
