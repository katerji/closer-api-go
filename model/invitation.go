package model

type Invitation struct {
	Id      int  `json:"id"`
	Inviter User `json:"inviter"`
	Contact User `json:"contact"`
}

func (i *Invitation) ToOutput() User {
	return i.Inviter
}
