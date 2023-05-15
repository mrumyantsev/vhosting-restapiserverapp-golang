package video

import (
	"app/pkg/user"
	"github.com/gin-gonic/gin"
)

type Video struct {
	Id         int    `json:"id"         db:"id"`
	Url        string `json:"url"        db:"url"`
	File       string `json:"file"       db:"file"`
	CreateDate string `json:"createDate" db:"create_date"`
	InfoId     int    `json:"infoId"     db:"info_id"`
	UserId     int    `json:"userId"     db:"user_id"`
}

type VideoCommon interface {
	CreateVideo(vid *Video) error
	GetVideo(id int) (*Video, error)
	GetAllVideos(urlparams *user.Pagin) (map[int]*Video, error)
	PartiallyUpdateVideo(vid *Video) error
	DeleteVideo(id int) error

	IsVideoExists(id int) (bool, error)
}

type VideoUseCase interface {
	VideoCommon

	BindJSONVideo(ctx *gin.Context) (*Video, error)
	IsRequiredEmpty(url, filename string) bool
	AtoiRequestedId(ctx *gin.Context) (int, error)
}

type VideoRepository interface {
	VideoCommon
}
