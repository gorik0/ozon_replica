package models

import uuid "github.com/satori/go.uuid"

//easyjson:json
type Profile struct {
	Id           uuid.UUID `json:"id"`
	Login        string    `json:"login"`
	Description  string    `json:"description,omitempty"`
	ImgSrc       string    `json:"img"`
	Phone        string    `json:"phone"`
	PasswordHash []byte    `json:"-"`
}

//
//func (p *Profile) LogValue() slog.Value {
//	return slog.GroupValue(
//		slog.String("id", p.Id.String()),
//		slog.String("login", p.Login),
//	)
//}
