package application

import (
	"github.com/jmhammock/ereader/cmd/server/wsroom"
	"github.com/jmhammock/ereader/internal/models"
)

type Application struct {
	FamilyModel   *models.FamilyModel
	UserModel     *models.UserModel
	UserRoleModel *models.UserRoleModel
	RoomModel     *models.RoomModel
	BookModel     *models.BookModel
	WSRooms       map[string]*wsroom.WSRoom
}
