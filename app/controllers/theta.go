package controllers

import (
	"github.com/revel/revel"
	"github.com/toshi3221/theta_v2"
	"github.com/toshi3221/theta_v2/command"
)

type Theta struct {
	*revel.Controller
}

func (c Theta) TakePicture() revel.Result {
	client, _ := theta_v2.NewClient("http://192.168.1.1")

	// camera.startSession
	revel.INFO.Println("camera.startSession:")
	startSessionCommand := new(command.StartSessionCommand)
	client.CommandExecute(startSessionCommand)

	revel.INFO.Println("  sessionId:", *startSessionCommand.Results.SessionId)
	sessionId := startSessionCommand.Results.SessionId

	// camera.takePicture
	revel.INFO.Println("camera.takePicture")
	takePictureCommand := new(command.TakePictureCommand)
	takePictureCommand.Parameters.SessionId = sessionId
	takePictureCommandResponse, _ := client.CommandExecute(takePictureCommand)

	// camera.closeSession
	revel.INFO.Println("camera.closeSession:")
	closeSessionCommand := new(command.CloseSessionCommand)
	closeSessionCommand.Parameters.SessionId = sessionId
	client.CommandExecute(closeSessionCommand)

	c.Response.Status = 202
	return c.RenderJson(takePictureCommandResponse)
}
