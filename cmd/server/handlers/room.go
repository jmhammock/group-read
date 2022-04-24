package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmhammock/ereader/cmd/server/application"
	"github.com/jmhammock/ereader/cmd/server/helpers"
	"github.com/jmhammock/ereader/internal/models"
	"github.com/julienschmidt/httprouter"
)

func GetRooms(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		limit, err := helpers.StrToUint32(r.URL.Query().Get("limit"))
		if err != nil {
			limit = 10
		}
		offset, err := helpers.StrToUint32(r.URL.Query().Get("offset"))
		if err != nil {
			offset = 0
		}

		rooms, err := app.RoomModel.Get(limit, offset)
		if err != nil {
			return
		}

		resp, err := json.Marshal(&rooms)
		if err != nil {
			return
		}
		w.Write(resp)
	}
}

func GetRoom(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		if id == "" {
			return
		}

		room, err := app.RoomModel.GetById(id)
		if err != nil {
			return
		}

		resp, err := json.Marshal(&room)
		w.Write(resp)
	}
}

func CreateRoom(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var room *models.Room
		err := json.NewDecoder(r.Body).Decode(&room)
		if err != nil {
			return
		}

		id, err := app.RoomModel.Insert(room.Name)
		if err != nil {
			return
		}

		room.Id = id

		resp, err := json.Marshal(&room)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(resp)
	}
}

func UpdateRoom(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		var room *models.Room
		err := json.NewDecoder(r.Body).Decode(&room)
		if err != nil {
			return
		}

		if id != room.Id {
			return
		}

		err = app.RoomModel.Update(room)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteRoom(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		var room *models.Room
		err := json.NewDecoder(r.Body).Decode(&room)
		if err != nil {
			return
		}

		if id != room.Id {
			return
		}

		err = app.RoomModel.Delete(id)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
