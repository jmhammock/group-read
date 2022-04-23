package handlers

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/jmhammock/ereader/cmd/web/application"
	"github.com/jmhammock/ereader/cmd/web/room"
	"github.com/julienschmidt/httprouter"
)

func RoomForm(fs embed.FS) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		t, err := template.ParseFS(
			fs,
			"templates/create-room.page.html",
			"templates/base.layout.html",
		)
		if err != nil {
			return
		}
		t.Execute(w, nil)
	}
}

func CreateRoom(app application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		err := r.ParseForm()
		if err != nil {
			return
		}

		room := &room.Room{}

		err = app.RoomRepository.Create(room)
		if err != nil {
			return
		}

		roomPath := fmt.Sprintf("/room/%s", room.Id)
		http.Redirect(w, r, roomPath, http.StatusTemporaryRedirect)
	}
}

func Room(fs embed.FS, app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		roomId := p.ByName("id")
		// _, err := app.RoomRepository.GetById(roomId)
		// if err != nil {
		// 	if errors.Is(err, sql.ErrNoRows) {
		// 		return
		// 	}
		// 	return
		// }

		t, err := template.ParseFS(
			fs,
			"templates/room.page.html",
			"templates/base.layout.html",
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		data := struct {
			Title string
		}{
			Title: roomId,
		}

		err = t.Execute(w, data)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
