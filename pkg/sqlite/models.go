package sqlite

import (
	"github.com/Duneus/go-kudos/pkg/gokudos"
)

type kudos struct {
	Message string
}

func mapKudosToModel(kudos2 gokudos.Kudos) *kudos {
	return &kudos{
		Message: kudos2.Message,
	}
}

func (k *kudos) toModel() gokudos.Kudos {
	return gokudos.Kudos{Message: k.Message}
}
