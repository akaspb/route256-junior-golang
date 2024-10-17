package event

import (
	"math"
	"time"
)

type Factory struct {
	idGenerator        IDGenerator
	timestampGenerator TimestampGenerator
}

type IDGenerator interface {
	Generate() int64
}

type TimestampGenerator interface {
	Generate() time.Time
}

func NewFactory(
	idGenerator IDGenerator,
	timestampGenerator TimestampGenerator,
) *Factory {
	return &Factory{
		idGenerator:        idGenerator,
		timestampGenerator: timestampGenerator,
	}
}

func NewDefaultFactory(startIDValue int64) *Factory {
	return NewFactory(NewSerialIDGen(startIDValue), NewClockGen())
}

func (f *Factory) Create(eventType EventType, operationSerialised []byte) Event {
	return Event{
		ID:        f.idGenerator.Generate(),
		EventType: eventType,
		Operation: operationSerialised,
		Timestamp: f.timestampGenerator.Generate(),
	}
}

type SerialIDGen struct {
	currentID int64
}

func NewSerialIDGen(startIDValue int64) *SerialIDGen {
	return &SerialIDGen{
		currentID: startIDValue,
	}
}

func (g *SerialIDGen) Generate() int64 {
	id := g.currentID

	if id < math.MaxInt64 {
		g.currentID++
	} else {
		g.currentID = 1
	}

	return id
}

type ClockGen struct{}

func NewClockGen() *ClockGen {
	return &ClockGen{}
}

func (g *ClockGen) Generate() time.Time {
	return time.Now()
}
