package controllers

import (
	"github.com/revel/revel"
	"github.com/toshi3221/theta_v2"
	"github.com/toshi3221/theta_v2/command"
	"io"
	"io/ioutil"
	"net/http"
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

func (c Theta) ImageList() revel.Result {
	client, _ := theta_v2.NewClient("http://192.168.1.1")

	listImagesCommand := new(command.ListImagesCommand)
	entryCount, includeThumb := 10, false
	listImagesCommand.Parameters.EntryCount = &entryCount
	listImagesCommand.Parameters.IncludeThumb = &includeThumb

	client.CommandExecute(listImagesCommand)
	entries := *listImagesCommand.Results.Entries

	return c.Render(entries)
}

type JpegResponse string

func (jr JpegResponse) Apply(req *revel.Request, resp *revel.Response) {
	resp.WriteHeader(http.StatusOK, "image/jpeg")
	resp.Out.Write([]byte(jr))
}

func (c Theta) ImageThumbnail(uri string) revel.Result {
	client, _ := theta_v2.NewClient("http://192.168.1.1")

	getImageCommand := new(command.GetImageCommand)
	imageType := "thumb"
	getImageCommand.Parameters.Type = &imageType
	getImageCommand.Parameters.FileUri = &uri
	response, _ := client.CommandExecute(getImageCommand)
	revel.INFO.Println("  fileUri:", uri)
	http_body := response.Results.(io.ReadCloser)
	defer http_body.Close()
	body, _ := ioutil.ReadAll(http_body)
	return JpegResponse(body)
}
