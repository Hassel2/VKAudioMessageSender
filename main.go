package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Application struct {
	token string
	path string
	id string
	ownerId string
	infoLog *log.Logger
	errorLog *log.Logger
	responseLog *log.Logger
}

func main() {
	id := flag.String("id", "", "Recevier peer_id")
	path := flag.String("p", "", "Path to audiofile")
	flag.Parse()

	token := os.Getenv("TOKEN");

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	responseLog := log.New(os.Stdout, "RESPONSE\t", log.Ldate|log.Ltime)
	app := Application {
		token: token,
		path: *path,
		id: *id,
		infoLog: infoLog,
		errorLog: errorLog,
		responseLog: responseLog,
	}

	if (app.id == "" || app.id == " ") {
		app.errorLog.Fatalln("Please use -id flag to set the reciever id")
		return
	}
	if (app.token == "") {
		app.errorLog.Fatalln("No access tokeb found use export TOKEN=<your_token> in terminal to set access token")
	}

	app.infoLog.Println("Formatting audio")
	app.FormatAudio()

	resp, _ := app.request("users.get?")
	userIdButFloat := resp.(map[string]interface{})["response"].([]interface{})[0].(map[string]interface{})["id"].(float64)
	userId := strconv.FormatFloat(userIdButFloat, 'f', 1, 64)
	userId = strings.Replace(userId, ".0", "", 1)

	app.ownerId = userId

	uploadUrl := app.GetUploadServer()
	file := app.AudioUploader(uploadUrl)
	fileId := app.AudioSaver(file)
	app.MessageSender(fileId)
}

func (app Application) FormatAudio() {
	command := exec.Command("ffmpeg", "-i", app.path, "-ac", "1", "output.mp3", "-y")
	err := command.Run()
	if err != nil {
		// app.errorLog.Fatalln(err)
	}
	command = exec.Command("ffmpeg", "-i", "output.mp3", "-c:a", "libvorbis", "-q:a", "4", "output.ogg", "-y")
	err = command.Run()
	if err != nil {
		// app.errorLog.Fatalln(err)
	}
}

