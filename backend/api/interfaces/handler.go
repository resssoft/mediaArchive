package interfaces

import "github.com/resssoft/mediaArchive/models"

type ItemHandler interface {
	ProcessItem(item models.Item) models.Item
	Check(item models.Item) bool
}
