package domain

import "github.com/music-tribe/uuid"

type ImageFile struct {
	Id          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	Filepath    string    `json:"filepath"`
}
