package idgen

import "github.com/google/uuid"

func MakeUUID() string {
	newUUID := uuid.New()
	return newUUID.String()
}
