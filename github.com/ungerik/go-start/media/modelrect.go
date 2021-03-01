package media

import (
	"github.com/ungerik/go-start/model"
	"image"
)

type ModelRect struct {
	MinX model.Int
	MinY model.Int
	MaxX model.Int
	MaxY model.Int
}

func (self *ModelRect) Rectangle() image.Rectangle {
	return image.Rect(self.MinX.GetInt(), self.MinY.GetInt(), self.MaxX.GetInt(), self.MaxY.GetInt())
}

func (self *ModelRect) SetRectangle(r image.Rectangle) {
	self.MinX = model.Int(r.Min.X)
	self.MinY = model.Int(r.Min.Y)
	self.MaxX = model.Int(r.Max.X)
	self.MaxY = model.Int(r.Max.Y)
}
