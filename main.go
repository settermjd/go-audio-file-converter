package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"golang.org/x/exp/maps"
)

type App struct {
	SourceFormat string
	formats      map[string]string
}

func NewApp() App {
	a := App{
		formats: map[string]string{
			"audio/mpeg":               "mp3",
			"application/octet-stream": "aac",
			"application/ogg":          "ogg",
			"audio/wave":               "wav",
		},
	}

	return a
}

// uploadFile uploads an audio file in one of the accepted audio file formats
// and stores it in the temporary data directory
func (a *App) uploadFile(r *http.Request) (*os.File, error) {
	// Attempt to open the uploaded file
	file, _, err := r.FormFile("file-upload")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer file.Close()

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := os.CreateTemp("data/tmp/audio", "upload-*")
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	allowedMimeTypes := maps.Keys(a.formats)
	index := slices.IndexFunc(allowedMimeTypes, func(mimetype string) bool {
		return mimetype == http.DetectContentType(fileBytes)
	})
	if index == -1 {
		err := errors.New("uploaded file was not an mp3, ogg, acc, or wav file")
		log.Print(err)
		return nil, err
	}
	a.SourceFormat = http.DetectContentType(fileBytes)

	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	return tempFile, nil
}

type PageData struct {
	ConvertToFormat, ConvertFromFormat, ConvertedFile string
	Error                                             error
}

func home(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		ConvertToFormat: "mp3",
		Error:           nil,
	}

	a := NewApp()

	if r.Method == "POST" {
		// Parse a multi-part form allowing for files to be uploaded to a maximum size of 10MB.
		r.ParseMultipartForm(10 << 20)

		data.ConvertToFormat = r.FormValue("format")
		log.Printf("converting audio file to %s format.", strings.ToUpper(data.ConvertToFormat))

		// Should handle the fact that the function returned false
		uploadedFile, err := a.uploadFile(r)
		if err != nil {
			log.Print(err.Error())
			data.Error = err
		}
		defer os.Remove(uploadedFile.Name())

		log.Print(uploadedFile.Name())
		data.ConvertedFile = fmt.Sprintf("static/audio/download.%s", data.ConvertToFormat)
		file, err := os.Create(data.ConvertedFile)
		if err != nil {
			log.Print(err.Error())
			data.Error = err
		}
		log.Printf("converted audio file is: %s", data.ConvertedFile)
		defer file.Close()
		err = ffmpeg.Input(uploadedFile.Name(), ffmpeg.KwArgs{"f": a.formats[a.SourceFormat]}).
			Output(file.Name(), ffmpeg.KwArgs{"f": data.ConvertToFormat}).
			OverWriteOutput().
			ErrorToStdOut().
			Run()

		if err != nil {
			log.Print(err.Error())
			data.Error = err
		}
	}

	ts, err := template.
		New("home.tmpl").
		Funcs(template.FuncMap{
			"upper": strings.ToUpper,
		}).
		ParseFiles("./templates/home.tmpl")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
