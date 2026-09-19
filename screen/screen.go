package screen

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"net"
	"syscall"
)

const (
	Width  = 1404
	Height = 1872
)

const (
	WaveformDU   int32 = 1 // Low Quality (Direct Update, fast, black & white only)
	WaveformGL16 int32 = 3 // Medium Quality (Grayscale Local 16, faster with "ghosting")
	WaveformGC16 int32 = 2 // High Quality (Grayscale 16)

	ioctlWaveformFlag int32 = 0xf000
)

const (
	bytesPerPixel = 2
	frameSize     = Width * Height * bytesPerPixel
	grayscaleSize = Width * Height
	totalSize     = frameSize + grayscaleSize
)

const sockAddr = "/var/run/rm2fb.sock"

// commands
const (
	msgInit        int32 = 0
	msgUpdate      int32 = 2
	msgUpdateBatch int32 = 3
)

var le = binary.LittleEndian

type Region struct {
	X0, Y0, X1, Y1 int
	Waveform       int32
}

type Screen struct {
	conn   *net.UnixConn
	pixels []byte
	canvas *image.Gray16
}

func Open() (*Screen, error) {
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: sockAddr, Net: "unix"})
	if err != nil {
		return nil, fmt.Errorf("connecting to rm2fb server (is rm2fb.service running?): %w", err)
	}

	fd, err := handshake(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	pixels, err := syscall.Mmap(fd, 0, totalSize, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	// The mapping keeps its own reference; the fd is no longer needed.
	syscall.Close(fd)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("mmap framebuffer: %w", err)
	}

	canvas := image.NewGray16(image.Rect(0, 0, Width, Height))
	draw.Draw(canvas, canvas.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	return &Screen{conn: conn, pixels: pixels, canvas: canvas}, nil
}

func (s *Screen) Close() {
	_ = syscall.Munmap(s.pixels)
	s.conn.Close()
}

func (s *Screen) Refresh(x0, y0, x1, y1 int, waveform int32) error {
	// Volcar canvas → pixels solo en la región afectada
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			c := s.canvas.Gray16At(x, y)
			offset := (y*Width + x) * bytesPerPixel
			le.PutUint16(s.pixels[offset:], gray16ToRGB565(c))
		}
	}

	msg := make([]byte, 4+32)
	le.PutUint32(msg[0:], uint32(msgUpdate))
	putUpdateParams(msg[4:], Region{x0, y0, x1, y1, waveform})
	if _, err := s.conn.Write(msg); err != nil {
		return fmt.Errorf("sending update: %w", err)
	}

	// The server acks every update with a bool.
	var ack [1]byte
	if _, err := io.ReadFull(s.conn, ack[:]); err != nil {
		return fmt.Errorf("reading update ack: %w", err)
	}
	if ack[0] == 0 {
		return fmt.Errorf("server rejected update")
	}
	return nil
}

func (s *Screen) RefreshRegions(regions []Region) error {
	if len(regions) == 0 {
		return nil
	}

	for _, r := range regions {
		for y := r.Y0; y < r.Y1; y++ {
			for x := r.X0; x < r.X1; x++ {
				c := s.canvas.Gray16At(x, y)
				offset := (y*Width + x) * bytesPerPixel
				le.PutUint16(s.pixels[offset:], gray16ToRGB565(c))
			}
		}
	}

	msg := make([]byte, 4+4+32*len(regions))
	le.PutUint32(msg[0:], uint32(msgUpdateBatch))
	le.PutUint32(msg[4:], uint32(len(regions)))
	for i, r := range regions {
		putUpdateParams(msg[8+32*i:], r)
	}
	if _, err := s.conn.Write(msg); err != nil {
		return fmt.Errorf("sending update batch: %w", err)
	}

	// The server acks the whole batch with a single bool.
	var ack [1]byte
	if _, err := io.ReadFull(s.conn, ack[:]); err != nil {
		return fmt.Errorf("reading update batch ack: %w", err)
	}
	if ack[0] == 0 {
		return fmt.Errorf("server rejected update batch")
	}
	return nil
}

func (s *Screen) Canvas() *image.Gray16 {
	return s.canvas
}

func handshake(conn *net.UnixConn) (int, error) {
	// Init message: tag + struct { bool ownSwtcon; 3 bytes padding;
	// int32 width, height, pixelFormat }.
	msg := make([]byte, 4+16)
	le.PutUint32(msg[0:], uint32(msgInit))
	// ownSwtcon = false (regular client), followed by padding.
	le.PutUint32(msg[8:], Width)
	le.PutUint32(msg[12:], Height)
	le.PutUint32(msg[16:], 0) // PixelFormat::RGB565
	if _, err := conn.Write(msg); err != nil {
		return 0, fmt.Errorf("sending init: %w", err)
	}

	// The server replies with the granted FbFormat...
	var format [12]byte
	if _, err := io.ReadFull(conn, format[:]); err != nil {
		return 0, fmt.Errorf("reading granted format: %w", err)
	}
	w, h, pf := le.Uint32(format[0:]), le.Uint32(format[4:]), le.Uint32(format[8:])
	if w != Width || h != Height || pf != 0 {
		return 0, fmt.Errorf("unexpected framebuffer format: %dx%d pixelFormat=%d", w, h, pf)
	}

	// ...followed by the framebuffer fd as SCM_RIGHTS ancillary data.
	buf := make([]byte, 1)
	oob := make([]byte, syscall.CmsgSpace(4))
	_, oobn, _, _, err := conn.ReadMsgUnix(buf, oob)
	if err != nil {
		return 0, fmt.Errorf("receiving framebuffer fd: %w", err)
	}
	cmsgs, err := syscall.ParseSocketControlMessage(oob[:oobn])
	if err != nil || len(cmsgs) == 0 {
		return 0, fmt.Errorf("parsing control message: %w", err)
	}
	fds, err := syscall.ParseUnixRights(&cmsgs[0])
	if err != nil || len(fds) == 0 {
		return 0, fmt.Errorf("no framebuffer fd received: %w", err)
	}
	return fds[0], nil
}

// putUpdateParams encodes a Region as wire-format UpdateParams: y1, x1, y2,
// x2, flags, waveform, temperatureOverride, extraMode. Coordinates are
// inclusive on the wire.
func putUpdateParams(buf []byte, r Region) {
	le.PutUint32(buf[0:], uint32(r.Y0))
	le.PutUint32(buf[4:], uint32(r.X0))
	le.PutUint32(buf[8:], uint32(r.Y1-1))
	le.PutUint32(buf[12:], uint32(r.X1-1))
	le.PutUint32(buf[16:], 0) // flags: RM2_FLAG_NONE
	le.PutUint32(buf[20:], uint32(r.Waveform|ioctlWaveformFlag))
	le.PutUint32(buf[24:], 0) // temperatureOverride (float32 0.0)
	le.PutUint32(buf[28:], 0) // extraMode
}

func gray16ToRGB565(c color.Gray16) uint16 {
	gray8 := uint16(c.Y >> 8)
	return (gray8>>3)<<11 | (gray8>>2)<<5 | (gray8 >> 3)
}
