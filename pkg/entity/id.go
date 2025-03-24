package entity

import "github.com/google/uuid"

type ID uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

func ParseID(id string) (ID, error) {
	parsedID, err := uuid.Parse(id)
	return ID(parsedID), err
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

func IsUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
