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
	"github.com/Kartikkala/drive_svc/subscription"
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
	driveRepository := drive.NewDriveRepository(app.DB)
	driveSvc := drive.NewDriveService(driveRepository)
	driveSvcWithHooks := drive.NewDriveServiceWithHooks(driveSvc)
	drive.RegisterAllDriveHooks(nc, driveSvcWithHooks)
	eventHandlers := subscription.NewDriveEventHandler(driveSvcWithHooks)
	subscription.AttachEvents(nc, eventHandlers)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("Drive service active and listening...")
	<-ctx.Done()
}
