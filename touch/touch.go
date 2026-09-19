package touch

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const touchDevice = "/dev/input/event2"

const MaxTouchSlots = 10

// Linux input event types.
const (
	evSyn = 0x00
	evAbs = 0x03
)

// EV_SYN codes.
const synReport = 0x00

// EV_ABS codes (Multitouch Protocol Type B).
const (
	absMtSlot       = 0x2f
	absMtTrackingID = 0x39
	absMtPositionX  = 0x35
	absMtPositionY  = 0x36
	absMtPressure   = 0x3a
)

type inputEvent struct {
	Sec   uint32
	Usec  uint32
	Type  uint16
	Code  uint16
	Value int32
}

// TouchPoint holds the latest known state of a single finger contact.
type TouchPoint struct {
	// ID is the kernel tracking id for this contact, or -1 if the slot is
	// not currently in use (no finger down).
	ID int32

	// X and Y in digitizer coordinates (see TouchMaxX / TouchMaxY).
	X, Y int32

	// Pressure reported by the touch controller, if any (0 otherwise).
	Pressure int32
}

type TouchState struct {
	Points [MaxTouchSlots]TouchPoint
}

// Touch reads finger touch events from the touchscreen digitizer and emits
// one TouchState per input frame (SYN_REPORT).
type Touch struct {
	file *os.File
	slot int
}

func Open() (*Touch, error) {
	f, err := os.Open(touchDevice)
	if err != nil {
		return nil, fmt.Errorf("opening touch device %s: %w", touchDevice, err)
	}
	t := &Touch{file: f}
	// Slots not yet touched by an ABS_MT_TRACKING_ID event are unused.
	return t, nil
}

// Close releases the device file descriptor.
func (t *Touch) Close() {
	t.file.Close()
}

// Read blocks until the next SYN_REPORT frame is complete and returns the
// updated TouchState. Call it in a loop to receive a continuous stream of
// state snapshots. Returns an error (typically io.EOF or io.ErrUnexpectedEOF)
// when the device is closed or an I/O failure occurs.
func (t *Touch) Read(state *TouchState) error {
	var ev inputEvent
	for {
		if err := binary.Read(t.file, binary.LittleEndian, &ev); err != nil {
			if err == io.EOF {
				return io.EOF
			}
			return fmt.Errorf("reading input event: %w", err)
		}

		switch ev.Type {
		case evAbs:
			switch ev.Code {
			case absMtSlot:
				t.slot = int(ev.Value)
			case absMtTrackingID:
				if t.slot >= 0 && t.slot < MaxTouchSlots {
					state.Points[t.slot].ID = ev.Value
				}
			case absMtPositionX:
				if t.slot >= 0 && t.slot < MaxTouchSlots {
					state.Points[t.slot].X = ev.Value
				}
			case absMtPositionY:
				if t.slot >= 0 && t.slot < MaxTouchSlots {
					state.Points[t.slot].Y = ev.Value
				}
			case absMtPressure:
				if t.slot >= 0 && t.slot < MaxTouchSlots {
					state.Points[t.slot].Pressure = ev.Value
				}
			}

		case evSyn:
			if ev.Code == synReport {
				// Frame complete — return the accumulated state.
				return nil
			}
		}
	}
}
