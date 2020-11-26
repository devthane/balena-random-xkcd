package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type comicJson struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	a := app.New()
	w := a.NewWindow("xkcd")
	w.SetFullScreen(true)

	resp, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
		log.Fatal(err)
	}

	originalComicData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	_ = resp.Body.Close()

	origComic := new(comicJson)
	err = json.Unmarshal(originalComicData, origComic)
	if err != nil {
		log.Fatal(err)
	}

	comicResource, err := fyne.LoadResourceFromURLString(origComic.Img)
	if err != nil {
		log.Fatalln(err)
	}
	img := canvas.NewImageFromResource(comicResource)
	img.FillMode = canvas.ImageFillContain

	//TODO: wordwrap the alt text if its going to overflow.
	card := widget.NewCard(
		origComic.Title,
		origComic.Alt,
		img,
	)

	w.SetContent(card)

	go func() {
		// TODO: Make the tick interval configurable via an environment variable.
		for range time.Tick(time.Second * 30) {
			randomComicIndex := rand.Intn(origComic.Num) + 1

			newComicResp, err := http.Get(fmt.Sprintf("https://xkcd.com/%d/info.0.json", randomComicIndex))
			if err != nil {
				log.Println(err)
				continue
			}

			newComicJDat, err := ioutil.ReadAll(newComicResp.Body)
			if err != nil {
				log.Println(err)
				continue
			}
			newComicResp.Body.Close()

			newComic := new(comicJson)
			err = json.Unmarshal(newComicJDat, newComic)
			if err != nil {
				log.Println(err)
				continue
			}

			newComicResource, err := fyne.LoadResourceFromURLString(newComic.Img)
			if err != nil {
				log.Println(err)
				continue
			}

			card.SetTitle(newComic.Title)
			card.SetSubTitle(newComic.Alt)
			img.Resource = newComicResource

			card.Refresh()
		}
	}()

	w.ShowAndRun()
}
