package entities

import "time"

type Entity[ID comparable] interface {
	GetID() ID
	SetID(id ID)
	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time)
}
