package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/hertz-contrib/websocket"
)

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
)

const (
	Channel uint32 = iota
	Room
	Broadcast
)

//const (
//	// size
//	_packSize      = 4
//	_headerSize    = 2
//	_verSize       = 2
//	_opSize        = 4
//	_seqSize       = 4
//	_heartSize     = 4
//	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
//	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
//	// offset
//	_packOffset   = 0
//	_headerOffset = _packOffset + _packSize
//	_verOffset    = _headerOffset + _headerSize
//	_opOffset     = _verOffset + _verSize
//	_seqOffset    = _opOffset + _opSize
//	_heartOffset  = _seqOffset + _seqSize
//)

const (
	_packetSize    = 4
	_versionSize   = 2
	_magicSize     = 4
	_operationSize = 4
	_sequenceSize  = 4
	_headerSize    = _versionSize + _magicSize + _operationSize + _sequenceSize
	_rawHeaderSize = _packetSize + _headerSize

	_packetOffset    = 0
	_versionOffset   = _packetOffset + _packetSize
	_magicOffset     = _versionOffset + _versionSize
	_operationOffset = _magicOffset + _magicSize
	_sequenceOffset  = _operationOffset + _operationSize
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

var (
	// ProtoReady proto ready
	ProtoReady = &Message{Operation: Ready}
	// ProtoFinish proto finish
	ProtoFinish = &Message{Operation: Finish}
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (m *Message) ToBytes() []byte {
	var packetLen = _rawHeaderSize + uint32(len(m.Data))
	var data = make([]byte, packetLen)
	binary.BigEndian.PutUint32(data[_packetOffset:], packetLen)
	binary.BigEndian.PutUint16(data[_versionOffset:], uint16(m.Version))
	binary.BigEndian.PutUint32(data[_magicOffset:], m.MagicNum)
	binary.BigEndian.PutUint32(data[_operationOffset:], m.Operation)
	binary.BigEndian.PutUint32(data[_sequenceOffset:], m.Sequence)
	if m.Data != nil {
		copy(data[_rawHeaderSize:], m.Data)
	}
	return data
}

func (m *Message) FromBytes(data []byte) (err error) {
	if len(data) < _rawHeaderSize {
		return fmt.Errorf("data's length is less than raw header size")
	}
	var packetLen uint32
	packetLen = binary.BigEndian.Uint32(data[_packetOffset:])
	if len(data) != int(packetLen) {
		return fmt.Errorf("data's length is not equal to packet length")
	}
	m.Version = uint32(binary.BigEndian.Uint16(data[_versionOffset:]))
	m.MagicNum = binary.BigEndian.Uint32(data[_magicOffset:])
	m.Operation = binary.BigEndian.Uint32(data[_operationOffset:])
	m.Sequence = binary.BigEndian.Uint32(data[_sequenceOffset:])
	if len(data[_rawHeaderSize:]) > 0 {
		//m.Data = make([]byte, len(data[_rawHeaderSize:]))
		//copy(m.Data, data[_rawHeaderSize:])
		m.Data = data[_rawHeaderSize:]
	}
	return nil
}

// WriteTo write a proto to bytes writer.
//func (m *Message) WriteTo(b *bytes.Writer) {
//	var (
//		packLen = _rawHeaderSize + int32(len(m.Body))
//		buf     = b.Peek(_rawHeaderSize)
//	)
//	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
//	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
//	binary.BigEndian.PutInt16(buf[_verOffset:], int16(m.Ver))
//	binary.BigEndian.PutInt32(buf[_opOffset:], m.Op)
//	binary.BigEndian.PutInt32(buf[_seqOffset:], m.Seq)
//	if m.Body != nil {
//		b.Write(m.Body)
//	}
//}

// ReadTCP read a proto from TCP reader.
//func (m *Message) ReadTCP(rr *bufio.Reader) (err error) {
//	var (
//		bodyLen   int
//		headerLen int16
//		packLen   int32
//		buf       []byte
//	)
//	if buf, err = rr.Pop(_rawHeaderSize); err != nil {
//		return
//	}
//	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
//	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
//	m.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
//	m.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
//	m.Seq = binary.BigEndian.Int32(buf[_seqOffset:])
//	if packLen > _maxPackSize {
//		return ErrProtoPackLen
//	}
//	if headerLen != _rawHeaderSize {
//		return ErrProtoHeaderLen
//	}
//	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
//		m.Body, err = rr.Pop(bodyLen)
//	} else {
//		m.Body = nil
//	}
//	return
//}

// WriteTCP write a proto to TCP writer.
//func (m *Message) WriteTCP(wr *bufio.Writer) (err error) {
//	var (
//		buf     []byte
//		packLen int32
//	)
//	if m.Op == OpRaw {
//		// write without buffer, job concact proto into raw buffer
//		_, err = wr.WriteRaw(m.Body)
//		return
//	}
//	packLen = _rawHeaderSize + int32(len(m.Body))
//	if buf, err = wr.Peek(_rawHeaderSize); err != nil {
//		return
//	}
//	binary.BigEndian.PutInt32(buf[_packOffset:], packLen)
//	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
//	binary.BigEndian.PutInt16(buf[_verOffset:], int16(m.Ver))
//	binary.BigEndian.PutInt32(buf[_opOffset:], m.Op)
//	binary.BigEndian.PutInt32(buf[_seqOffset:], m.Seq)
//	if m.Body != nil {
//		_, err = wr.Write(m.Body)
//	}
//	return
//}

// WriteTCPHeart write TCP heartbeat with room online.
//func (m *Message) WriteTCPHeart(wr *bufio.Writer, online int32) (err error) {
//	var (
//		buf     []byte
//		packLen int
//	)
//	packLen = _rawHeaderSize + _heartSize
//	if buf, err = wr.Peek(packLen); err != nil {
//		return
//	}
//	// header
//	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
//	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
//	binary.BigEndian.PutInt16(buf[_verOffset:], int16(m.Ver))
//	binary.BigEndian.PutInt32(buf[_opOffset:], m.Op)
//	binary.BigEndian.PutInt32(buf[_seqOffset:], m.Seq)
//	// body
//	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
//	return
//}

type WebsocketReadWriter interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

// ReadWebsocket read a proto from websocket connection.
func (m *Message) ReadWebsocket(conn WebsocketReadWriter) (err error) {
	var (
		packetLen uint32
		data      []byte
	)
	if _, data, err = conn.ReadMessage(); err != nil {
		return
	}
	if len(data) < _rawHeaderSize {
		return fmt.Errorf("data's length is less than raw header size")
	}
	packetLen = binary.BigEndian.Uint32(data[_packetOffset:_versionOffset])
	if len(data) != int(packetLen) {
		return fmt.Errorf("data's length is not equal to packet length")
	}
	m.Version = uint32(binary.BigEndian.Uint16(data[_versionOffset:_magicOffset]))
	m.MagicNum = binary.BigEndian.Uint32(data[_magicOffset:_operationOffset])
	m.Operation = binary.BigEndian.Uint32(data[_operationOffset:_sequenceOffset])
	m.Sequence = binary.BigEndian.Uint32(data[_sequenceOffset:])
	if len(data[_rawHeaderSize:]) > 0 {
		//zlog.Infoln("read data")
		//copy(m.Data, data[_rawHeaderSize:])
		//zlog.Infof("data:%s", string(m.Data))
		m.Data = data[_rawHeaderSize:]
	}
	return nil
}

// WriteWebsocket write a proto to websocket connection.
func (m *Message) WriteWebsocket(conn WebsocketReadWriter) (err error) {
	var packetLen = _rawHeaderSize + uint32(len(m.Data))
	var data = make([]byte, packetLen)
	binary.BigEndian.PutUint32(data[_packetOffset:], packetLen)
	binary.BigEndian.PutUint16(data[_versionOffset:], uint16(m.Version))
	binary.BigEndian.PutUint32(data[_magicOffset:], m.MagicNum)
	binary.BigEndian.PutUint32(data[_operationOffset:], m.Operation)
	binary.BigEndian.PutUint32(data[_sequenceOffset:], m.Sequence)
	if m.Data != nil {
		//zlog.Infoln("write data")
		copy(data[_rawHeaderSize:], m.Data)
	}
	return conn.WriteMessage(websocket.BinaryMessage, data)
}

// WriteWebsocketHeart write websocket heartbeat with room online.
//func (m *Message) WriteWebsocketHeart(wr *websocket.Conn, online int32) (err error) {
//	var (
//		buf     []byte
//		packLen int
//	)
//	packLen = _rawHeaderSize + _heartSize
//	// websocket header
//	if err = wr.WriteHeader(websocket.BinaryMessage, packLen); err != nil {
//		return
//	}
//	if buf, err = wr.Peek(packLen); err != nil {
//		return
//	}
//	// proto header
//	binary.BigEndian.PutInt32(buf[_packOffset:], int32(packLen))
//	binary.BigEndian.PutInt16(buf[_headerOffset:], int16(_rawHeaderSize))
//	binary.BigEndian.PutInt16(buf[_verOffset:], int16(m.Ver))
//	binary.BigEndian.PutInt32(buf[_opOffset:], m.Op)
//	binary.BigEndian.PutInt32(buf[_seqOffset:], m.Seq)
//	// proto body
//	binary.BigEndian.PutInt32(buf[_heartOffset:], online)
//	return
//}
