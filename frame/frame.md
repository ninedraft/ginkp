# frame
--
    import "ginkp/frame"


## Usage

```go
var (
	ErrFrameIsNotCommand = errors.New("frame isn't command")
	ErrEmptyFrame        = errors.New("empty frame")
)
```

#### func  DeleteFrame

```go
func DeleteFrame(frame *Frame)
```

#### type Frame

```go
type Frame struct {
	bytes.Buffer
}
```


#### func  NewFrame

```go
func NewFrame() *Frame
```

#### func (*Frame) Args

```go
func (frame *Frame) Args() (string, error)
```
Args ... Returns string, which represents command args or error
ErrFrameIsNotCommand, if frame doesn't contain command or ErrEmptyFrame if frame
is empty (less 3 bytes)

#### func (*Frame) CommandCode

```go
func (frame *Frame) CommandCode() (CommandCode, error)
```
CommandCode ... Returns CommandCode or error ErrFrameIsNotCommand, if frame
doesn't contain command or ErrEmptyFrame, if frame is empty (less 3 bytes)

#### func (*Frame) Copy

```go
func (frame *Frame) Copy() (*Frame, error)
```
Copy ... Makes full frame copy.

#### func (*Frame) DataReader

```go
func (frame *Frame) DataReader() (io.Reader, error)
```
DataReader ... Returns io.Reader, which reads bytes from frame data section
Doesn't copy frame data, so all changes over original frame will be reflected in
io.Reader data

#### func (*Frame) DataSize

```go
func (frame *Frame) DataSize() int
```
DataSize ... Returns frame data section size (0 if frame is less then 3 bytes)

#### func (*Frame) IsCommand

```go
func (frame *Frame) IsCommand() bool
```

#### func (*Frame) M_ADR

```go
func (frame *Frame) M_ADR(addr []string)
```
M_ADR ... A list of 5D addresses delimited by spaces. eg, "2:5047/13@fidonet
2:5047/0@fidonet"

#### func (*Frame) M_BSY

```go
func (frame *Frame) M_BSY(err error)
```
M_BSY ... Our system sends it if it is busy. The receiving partner ignores the
argument (logs it). eg, "Too many servers are running already"

#### func (*Frame) M_EOB

```go
func (frame *Frame) M_EOB()
```
M_EOB ... End­of­Batch. EOB is transmitted after all the files have been sent.
If we are in the EOB state (all the files are sent), we get EOB from the remote
(no more files for us), we received all acknowledgements for all the sent files,
we received all the files resent in reply to GET, then the session is considered
to be successfully completed.

#### func (*Frame) M_ERR

```go
func (frame *Frame) M_ERR(err error)
```
M_ERR ... A fatal error. The partner who has sent M_ERR aborts the session. The
argument contains the text explaining the reason and it is logged. Binkd sends
M_ERR in response to an incorrect password. eg, "Incorrect password"

#### func (*Frame) M_FILE

```go
func (frame *Frame) M_FILE(file *FileConfig)
```
M_FILE ... The properties of the next file. They are delimited by spaces:
filename without spaces, size, UNIX­ time, the offset to transfer the file. All
the numbers are decimal. All the data blocks received after that relate to this
file until the next M_FILE is received. There is no special end­of­file marker
since the file size is known beforehand. Binkd will append the "excessive"
blocks to the current file. We start transmitting every new file from the offset
0. On receiving M_GET from the remote system we must do the seek operation. eg,
"config.sys 125 2476327846 0"

#### func (*Frame) M_GET

```go
func (frame *Frame) M_GET(file *FileConfig)
```
M_GET ... M_GET is used as a request to resend a file. The M_GET arguments copy
the arguments of the M_FILE command which we’d like to see from the remote
system. :) Binkd sends it as a response to M_FILE if it does not like the offset
from which the file transmission has been started by the remote system. eg,
"config.sys 125 2476327846 100" At present binkd handles it as follows:
according to the first fields (name/size/UNIX­time) it determines whether the
M_GET argument is the file we currently transmit (or the file has been
transmitted and we are waiting for M_GOT for it). If this is the case it seeks
the specified offset in the file and sends M_FILE after that. For the example
above M_FILE will have the following arguments: "config.sys 125 2476327846 100"

#### func (*Frame) M_GOT

```go
func (frame *Frame) M_GOT(file *FileConfig)
```
M_GOT ... It is sent as an acknowledgement by the system which has received a
file after receiving the last portion of the file data. The arguments are copies
of the FILE command arguments received from the remote system except the last
one, the offset which should not be returned to the system which sent M_FILE.
GOT may also be sent during the process of receiving a file; the sending partner
should react to it with the destructive skip. eg, "config.sys 125 2476327846"

#### func (*Frame) M_NUL

```go
func (frame *Frame) M_NUL(msg string)
```
M_NUL ... The command argument is ignored (and is possibly logged). This is the
way we transmit the nodelist information, the sysop’s name and so on. eg, "ZYZ
Dima Maloff"

#### func (*Frame) M_OK

```go
func (frame *Frame) M_OK()
```
M_OK ... A reply to the correct password. The binkd client rescans the queue
after receiving the message. The command argument is ignored.

#### func (*Frame) M_PWD

```go
func (frame *Frame) M_PWD(password string)
```
M_PWD ... A password. After the successful processing of the password received
from the remote, the binkd server rescans the queue. eg, "pAsSwOrD"

#### func (*Frame) M_SKIP

```go
func (frame *Frame) M_SKIP(file *FileConfig)
```
M_SKIP ... Non destructive skip. An example of the argument line: "config.sys
125 2476327846"

#### func (*Frame) ReadDataFrom

```go
func (frame *Frame) ReadDataFrom(r io.Reader) (int64, error)
```
ReadDataFrom ... Read bytes from io.Reader to frame data section. Returns error
if occcures.

#### func (*Frame) Reader

```go
func (frame *Frame) Reader() io.Reader
```
Reader ... Returns io.Reader, which reads bytes from frame(including header).
Doesn't copy frame, so all changes over original frame will be reflected in
io.Reader data

#### func (*Frame) UpdateLen

```go
func (frame *Frame) UpdateLen()
```

#### func (*Frame) WriteDataTo

```go
func (frame *Frame) WriteDataTo(w io.Writer) (int64, error)
```
WriteDataTo ... Write bytes from frame data section to io.Writer. Returns error
if occcures.
