package models

type Room struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Rooms []Room

type RoomModel struct {
	db DB
}

func NewRoomModel(db DB) *RoomModel {
	return &RoomModel{
		db: db,
	}
}

func (r RoomModel) Get(limit, offset uint32) (*Rooms, error) {
	q := `SELECT *
		FROM rooms r
		LIMIT $1
		OFFSET $2;`

	rows, err := r.db.Query(q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms Rooms
	for rows.Next() {
		var room Room
		err := rows.Scan(&room)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &rooms, nil
}

func (r RoomModel) GetById(id string) (*Room, error) {
	q := `SELECT *
		FROM rooms r
		WHERE r.id = $1;`

	var room *Room
	err := r.db.QueryRow(q, id).Scan(&room)

	return room, err
}

func (r RoomModel) Insert(name string) (string, error) {
	q := `INSERT INTO rooms (name) VALUES($1)
		RETURNING id;`

	var id string
	err := r.db.QueryRow(q, name).Scan(&id)

	return id, err
}

func (r RoomModel) Update(room *Room) error {
	q := `UPDATE rooms
		SET name = $1
		WHERE id = $2;`

	_, err := r.db.Execute(q, room.Name, room.Id)
	return err
}

func (r RoomModel) Delete(id string) error {
	q := `DELETE FROM rooms
		WHERE id = $1;`

	_, err := r.db.Execute(q, id)
	return err
}
