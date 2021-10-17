package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (app Application) request(req string) (interface{}, error) {
	url := fmt.Sprintf("https://api.vk.com/method/%s&access_token=%s&v=5.131", req, app.token)
	response, err := http.Get(url)
	url = strings.Replace(url, app.token, "your_token", 1)
	app.infoLog.Println("Sending request to " + url)
	if err != nil {
		// app.errorLog.Fatalln(err)
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// app.errorLog.Panicln(err)
		return nil, err
	}
	app.responseLog.Println(string(responseBody))
	defer response.Body.Close()

	var jsonUnwrapper interface{}
	err = json.Unmarshal(responseBody, &jsonUnwrapper)
	if err != nil {
		// app.errorLog.Fatalln(err)
		return nil, err
	}

	app.infoLog.Println("Done")

	return jsonUnwrapper, nil 
}

func (app Application) GetUploadServer() (string) {
	app.infoLog.Println("Getting upload server url")

	response, err := app.request(fmt.Sprintf("docs.getMessagesUploadServer?type=audio_message&peer_id=%s", app.ownerId))
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	data := response.(map[string]interface{})["response"].(map[string]interface{})
	uploadUrl := data["upload_url"].(string)

	return uploadUrl
}

func (app Application) AudioUploader(uploadUrl string) (string) {
        app.infoLog.Println("Uploading audiomessage to server")

	app.infoLog.Println("Opening audio file")
	file, err := os.Open("bassboost/output.ogg")
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	defer file.Close()

	app.infoLog.Println("Formating request body")
	var requsetBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requsetBody)

	fileWriter, err := multiPartWriter.CreateFormFile("file", "output.ogg")
	if err != nil {
		app.errorLog.Fatalln(err)
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	multiPartWriter.Close()

	app.infoLog.Println("Sending request to " + uploadUrl)
	response, _ := http.Post(uploadUrl, multiPartWriter.FormDataContentType(), &requsetBody)

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	response.Body.Close()
	app.responseLog.Println(string(responseBody))

	var jsonUnwrapper interface{}
	err = json.Unmarshal(responseBody, &jsonUnwrapper)
	if err != nil {
		app.errorLog.Fatalln(err)
	}

	audioMessage := jsonUnwrapper.(map[string]interface{})["file"].(string)

	app.infoLog.Println("Done")
	return audioMessage
}

func (app Application) AudioSaver(file string) string {
	app.infoLog.Println("Saving an audiomessage")

	data, err := app.request(fmt.Sprintf("docs.save?file=%s&title=audio", file))
	if err != nil {
		app.errorLog.Fatalln(err)
	}
	audioMessageId := strconv.FormatFloat(data.(map[string]interface{})["response"].(map[string]interface{})["audio_message"].(map[string]interface{})["id"].(float64), 'f', 1, 64)
	audioMessageId = strings.Replace(audioMessageId, ".0", "", 1)

	return audioMessageId
}

func (app Application) MessageSender(audiomessageId string) {
	app.infoLog.Println("Sending a message to id=" + app.id)
	_, err := app.request(fmt.Sprintf("messages.send?peer_id=%s&random_id=%s&attachment=audio_message%s_%s", app.id, audiomessageId, app.ownerId,audiomessageId))
	if err != nil {
		app.errorLog.Fatalln(err)
	}
}
