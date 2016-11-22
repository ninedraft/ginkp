package frame

// binkp specification: http://ftsc.org/docs/fts-1026.001

import "bytes"
import "errors"
import "io"
import "sync"
import (
	"encoding/binary"
	"ginkp/command"
)

var (
	// ErrFrameIsNotCommand ...
	ErrFrameIsNotCommand = errors.New("frame isn't command")
	// ErrEmptyFrame ...
	ErrEmptyFrame = errors.New("empty frame")
	framePool     = sync.Pool{
		New: func() interface{} {
			return &Frame{}
		},
	}
)

// Frame ...
// Pinkb protocol frame.
// The frame header contains two bytes defining the type and the length (in bytes) of the data that follows the
// header. If the most significant bit of the header is reset, then all the data received with the frame should be
// appended to the current file being received if the file has been opened, otherwise the data should be
// discarded. If the bit is set, the data should be considered as a command changing the protocol state. The first
// data byte in the frame is the command number. The rest of the bytes form an argument. A command
// argument is an arbitrary character set not required to be bounded with a '\0'.
//
//   +­-+-------­­­­­­­+­­­­­­­­--------+--------­­­­­­­------­+
//   |1|  HI  0   LO   1|     DATA     |
//   +-+-------+-­­­­­­­­­­­­­­­­­­­­­­­­­­­­­­-------+-------------­­­­­­­­-+
//   |                  | <- 32k max-> |
//   |<-    2 bytes   ->|              + datablock/msg's argument
//                      +­­­­­­­­ length of the frame without the header ­ 1 byte
//    +­ this is a command
type Frame struct {
	bytes.Buffer
}

// NewFrame ...
// New binkp frame from frame pool.
func NewFrame() *Frame {
	return framePool.Get().(*Frame)
}

// DeleteFrame ...
// Resets and returns frame to frame pool for reusing in the future
func DeleteFrame(frame *Frame) {
	frame.Reset()
	framePool.Put(frame)
}

// IsCommand ...
func (frame *Frame) IsCommand() bool {
	if frame.Len() < 2 {
		frame.prepare(2)
		return false
	}
	return frame.Bytes()[0]&0x80>>7 == 1
}

// UpdateLen ...
// Recalcultes and updates len field.
func (frame *Frame) UpdateLen() {
	fLen := uint16(frame.Len())
	fBytes := frame.Bytes()
	if fLen > 2 {
		binary.BigEndian.PutUint16(fBytes[0:2], (fLen-2)&0x7FF)
	} else {
		frame.Grow(2 - int(fLen))
	}
}

// prepare ...
// Does preparations to write command code.
// Frows frame or truncate if need
func (frame *Frame) prepare(n int) {
	fLen := frame.Len()
	if fLen < n {
		buf := getBytes()
		p := (*buf)[:n]
		_, err := frame.Write(p)
		returnBytes(buf)
		if err != nil {
			panic(err)
		}
	} else if fLen > n {
		frame.Truncate(fLen - n)
	}
}

func (frame *Frame) setCommandFlag(command func(*Frame)) {
	frame.prepare(3)
	frame.Bytes()[0] = frame.Bytes()[0] | 0x80
	command(frame)
}

func (frame *Frame) writeCommandWithArgs(command command.Code, args []byte) {
	fBytes := frame.Bytes()
	frame.prepare(3)
	fBytes[2] = byte(command)
	_, err := frame.Write(args)
	if err != nil {
		panic(err)
	}
	frame.UpdateLen()
}

func (frame *Frame) writeCommand(command command.Code) {
	fBytes := frame.Bytes()
	frame.prepare(3)
	fBytes[2] = byte(command)
	frame.UpdateLen()
}

// SetNUL ...
// The command argument is ignored (and is possibly logged).
// This is the way we transmit the nodelist information,
// the sysop’s name and so on.
// eg, "ZYZ Dima Maloff"
func (frame *Frame) SetNUL(msg string) {
	frame.writeCommandWithArgs(command.M_NUL, []byte(msg))
}

// SetADR ...
// A list of 5D addresses delimited by spaces.
// eg, "2:5047/13@fidonet 2:5047/0@fidonet"
func (frame *Frame) SetADR(addr []string) {
	frame.writeCommand(command.M_ADR)
	for _, address := range addr {
		_, err := frame.WriteString(address + " ")
		assert(err)
	}
}

// SetPWD ...
// A password. After the successful processing of the password received from the remote,
// the binkd server rescans the queue.
// eg, "pAsSwOrD"
func (frame *Frame) SetPWD(password string) {
	frame.writeCommandWithArgs(command.M_PWD, []byte(password))
}

// SetFILE ...
// The properties of the next file. They are delimited by spaces: filename without spaces, size, UNIX­
// time, the offset to transfer the file. All the numbers are decimal. All the data blocks received after that
// relate to this file until the next M_FILE is received. There is no special end­of­file marker since the file
// size is known beforehand. Binkd will append the "excessive" blocks to the current file. We start
// transmitting every new file from the offset 0. On receiving M_GET from the remote system we must
// do the seek operation.
// eg, "config.sys 125 2476327846 0"
func (frame *Frame) SetFILE(file *command.FileConfig) {
	frame.writeCommandWithArgs(command.M_FILE, []byte(file.String()))
}

// SetOK ...
// A reply to the correct password. The binkd client rescans the queue after receiving the message.
// The command argument is ignored.
func (frame *Frame) SetOK() {
	frame.writeCommand(command.M_OK)
}

// SetEOB ...
// End­of­Batch. EOB is transmitted after all the files have been sent. If we are in the EOB state (all the
// files are sent), we get EOB from the remote (no more files for us), we received all acknowledgements
// for all the sent files, we received all the files resent in reply to GET,
// then the session is considered to be successfully completed.
func (frame *Frame) SetEOB() {
	frame.writeCommand(command.M_EOB)
}

// SetGOT ...
// It is sent as an acknowledgement by the system which has received a file after receiving the last
// portion of the file data. The arguments are copies of the FILE command arguments received from the
// remote system except the last one, the offset which should not be returned to the system which sent
// M_FILE. GOT may also be sent during the process of receiving a file; the sending partner should react
// to it with the destructive skip.
// eg, "config.sys 125 2476327846"
func (frame *Frame) SetGOT(file *command.FileConfig) {
	frame.writeCommandWithArgs(command.M_FILE, []byte(file.StringWithoutOff()))
}

// SetERR ...
// A fatal error. The partner who has sent M_ERR aborts the session. The argument contains the text
// explaining the reason and it is logged. Binkd sends M_ERR in response to an incorrect password.
// eg, "Incorrect password"
func (frame *Frame) SetERR(err error) {
	frame.writeCommandWithArgs(command.M_ERR, []byte(err.Error()))
}

// SetBSY ...
// Our system sends it if it is busy. The receiving partner ignores the argument (logs it).
// eg, "Too many servers are running already"
func (frame *Frame) SetBSY(err error) {
	frame.writeCommandWithArgs(command.M_ERR, []byte(err.Error()))
}

// SetGET ...
// M_GET is used as a request to resend a file. The M_GET arguments copy the arguments of the
// M_FILE command which we’d like to see from the remote system. :) Binkd sends it as a response to
// M_FILE if it does not like the offset from which the file transmission has been started by the remote
// system.
// eg, "config.sys 125 2476327846 100"
// At present binkd handles it as follows: according to the first fields (name/size/UNIX­time) it
// determines whether the M_GET argument is the file we currently transmit (or the file has been
// transmitted and we are waiting for M_GOT for it). If this is the case it seeks the specified offset in the
// file and sends M_FILE after that. For the example above M_FILE will have the following arguments:
// "config.sys 125 2476327846 100"
func (frame *Frame) SetGET(file *command.FileConfig) {
	frame.writeCommandWithArgs(command.M_FILE, []byte(file.String()))
}

// SetSKIP ...
// Non destructive skip. An example of the argument line:
// "config.sys 125 2476327846"
func (frame *Frame) SetSKIP(file *command.FileConfig) {
	frame.writeCommandWithArgs(command.M_FILE, []byte(file.StringWithoutOff()))
}

// CommandCode ...
// Returns command.Code or error ErrFrameIsNotCommand,
// if frame doesn't contain command or
// ErrEmptyFrame, if frame is empty (less 3 bytes)
func (frame *Frame) CommandCode() (command.Code, error) {
	if !frame.IsCommand() {
		return 0, ErrFrameIsNotCommand
	}
	fLen := frame.Len()
	fBytes := frame.Bytes()
	if fLen >= 3 {
		return command.Code(fBytes[2]), nil
	}
	return 0, ErrEmptyFrame
}

// Args ...
// Returns string, which represents command args or error
// ErrFrameIsNotCommand, if frame doesn't contain command or
// ErrEmptyFrame if frame is empty (less 3 bytes)
func (frame *Frame) Args() (string, error) {
	fLen := frame.Len()
	fBytes := frame.Bytes()
	if fLen >= 3 {
		if (fBytes[0]&0x80)>>1 != 1 {
			return "", ErrFrameIsNotCommand
		}
		return string(fBytes[3:]), nil
	}
	return "", ErrEmptyFrame
}

// WriteDataTo ...
// Write bytes from frame data section to io.Writer.
// Returns error if occcures.
func (frame *Frame) WriteDataTo(w io.Writer) (int64, error) {
	if frame.Len() < 2 {
		return 0, ErrEmptyFrame
	}
	buf := bytes.NewBuffer(frame.Bytes()[2:])
	n, err := buf.WriteTo(w)
	frame.UpdateLen()
	return n, err
}

// ReadDataFrom ...
// Read bytes from io.Reader to frame data section.
// Returns error if occcures.
func (frame *Frame) ReadDataFrom(r io.Reader) (int64, error) {
	if frame.IsCommand() {
		frame.prepare(3)
	} else {
		frame.prepare(2)
	}
	n, err := frame.ReadFrom(r)
	frame.UpdateLen()
	return n, err
}

// Copy ...
// Makes full frame copy.
func (frame *Frame) Copy() (*Frame, error) {
	copyFrame := NewFrame()
	_, err := copyFrame.ReadFrom(frame)
	return copyFrame, err
}

// Reader ...
// Returns io.Reader, which reads bytes from frame(including header).
// Doesn't copy frame, so all changes over original frame
// will be reflected in io.Reader data
func (frame *Frame) Reader() io.Reader {
	return bytes.NewReader(frame.Bytes())
}

// DataReader ...
// Returns io.Reader, which reads bytes from frame data section
// Doesn't copy frame data, so all changes over original frame
// will be reflected in io.Reader data
func (frame *Frame) DataReader() (io.Reader, error) {
	fLen := frame.Len()
	if fLen >= 3 {
		if frame.IsCommand() {
			return bytes.NewReader(frame.Bytes()[3:]), nil
		}
		return bytes.NewReader(frame.Bytes()[2:]), nil
	}
	return nil, ErrEmptyFrame
}

// DataSize ...
// Returns frame data section size (0 if frame is less then 2 bytes)
func (frame *Frame) DataSize() int {
	fLen := frame.Len()
	if fLen >= 2 {
		return int(binary.BigEndian.Uint16(frame.Bytes()[:2]))
	}
	return 0
}
