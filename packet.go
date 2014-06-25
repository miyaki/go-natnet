package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

// contants for natnet packet handling
const (
	MAX_PACKETSIZE = 100000

	NAT_PING                 = 0
	NAT_PINGRESPONSE         = 1
	NAT_REQUEST              = 2
	NAT_RESPONSE             = 3
	NAT_REQUEST_MODELDEF     = 4
	NAT_MODELDEF             = 5
	NAT_REQUEST_FRAMEOFDATA  = 6
	NAT_FRAMEOFDATA          = 7
	NAT_MESSAGESTRING        = 8
	NAT_UNRECOGNIZED_REQUEST = 100
)

// Packet in NatNet Protocol
type Packet struct {
	MessageID int16
	dataLen   uint16
	frame     *MocapFrame
	//	model *Model
}

// Point3f holds 3d point data
type Point3f struct {
	x, y, z float32
}

// Quaternion4f
type Quaternion4f struct {
	qx, qy, qz, qw float32
}

// RigidBody has info for rigid object
type RigidBody struct {
	ID             int32
	p              Point3f
	qx, qy, qz, qw float32
	nRigidMarkers  int32
	MarkerData     []Point3f
	MarkerID       []uint32
	MarkerSize     []float32
	MarkerError    float32
}

// MarkerSet holds 3d point
type MarkerSet struct {
	x, y, z float32
}

//Skelton not yet supported
type Skelton struct {
}

//type LabeledMarker

// MocapFrame conntains frame info
type MocapFrame struct {
	frameNumber         int32
	nMarkerSets         int32
	MarkerSets          []Point3f
	nMarkerUnidentified int32
	MarkerUnidentified  []Point3f
	nRigidBodies        int32
	RigidBodies         []RigidBody
	nSkeltons           int32
	Skeltons            []Skelton
	latency             float32
	timecode            uint32
	//subtimecode       uint32
}

//interface PacketReader

func (rb *RigidBody) decode(reader io.Reader) {
	binary.Read(reader, binary.LittleEndian, &rb.ID)

	binary.Read(reader, binary.LittleEndian, &rb.p.x)
	binary.Read(reader, binary.LittleEndian, &rb.p.y)
	binary.Read(reader, binary.LittleEndian, &rb.p.z)

	binary.Read(reader, binary.LittleEndian, &rb.qx)
	binary.Read(reader, binary.LittleEndian, &rb.qy)
	binary.Read(reader, binary.LittleEndian, &rb.qz)
	binary.Read(reader, binary.LittleEndian, &rb.qw)

	binary.Read(reader, binary.LittleEndian, &rb.nRigidMarkers)

	rb.MarkerData = make([]Point3f, rb.nRigidMarkers)
	for i := range rb.MarkerData {
		rb.MarkerData[i].decode(reader)
	}

	//if Major >=2
	rb.MarkerID = make([]uint32, rb.nRigidMarkers)
	for i := range rb.MarkerID {
		binary.Read(reader, binary.LittleEndian, &rb.MarkerID[i])
	}
	rb.MarkerSize = make([]float32, rb.nRigidMarkers)
	for i := range rb.MarkerSize {
		binary.Read(reader, binary.LittleEndian, &rb.MarkerSize[i])
	}

	binary.Read(reader, binary.LittleEndian, &rb.MarkerError)
}

func (ms *MarkerSet) decode(reader io.Reader) {
	binary.Read(reader, binary.LittleEndian, &ms.x)
	binary.Read(reader, binary.LittleEndian, &ms.y)
	binary.Read(reader, binary.LittleEndian, &ms.z)
}

func (p *Point3f) decode(reader io.Reader) {
	binary.Read(reader, binary.LittleEndian, &p.x)
	binary.Read(reader, binary.LittleEndian, &p.y)
	binary.Read(reader, binary.LittleEndian, &p.z)
}

func (mf *MocapFrame) decode(reader io.Reader) {
	binary.Read(reader, binary.LittleEndian, &mf.frameNumber)

	binary.Read(reader, binary.LittleEndian, &mf.nMarkerSets)

	mf.MarkerSets = make([]Point3f, mf.nMarkerSets)
	for i := range mf.MarkerSets {
		mf.MarkerSets[i].decode(reader)
	}

	binary.Read(reader, binary.LittleEndian, &mf.nMarkerUnidentified)

	mf.MarkerUnidentified = make([]Point3f, mf.nMarkerUnidentified)
	for i := range mf.MarkerUnidentified {
		mf.MarkerUnidentified[i].decode(reader)
	}

	binary.Read(reader, binary.LittleEndian, &mf.nRigidBodies)

	mf.RigidBodies = make([]RigidBody, mf.nRigidBodies)
	for i := range mf.RigidBodies {
		mf.RigidBodies[i].decode(reader)
	}

	// if 2.1
	binary.Read(reader, binary.LittleEndian, &mf.nSkeltons)

	//if 2.3
	//labeled markers

	binary.Read(reader, binary.LittleEndian, &mf.latency)

	//if 2.3
	binary.Read(reader, binary.LittleEndian, &mf.timecode)

	//subtimecode
}

func (p *Packet) decode(reader io.Reader) {
	p.MessageID = int16(0)
	binary.Read(reader, binary.LittleEndian, &p.MessageID)

	p.dataLen = uint16(0)
	binary.Read(reader, binary.LittleEndian, &p.dataLen)

	switch p.MessageID {
	default:
		fmt.Printf("unknown MessageId: %d\n", p.MessageID)
	case NAT_FRAMEOFDATA:
		p.frame = NewFrame()
		p.frame.decode(reader)
	case NAT_MODELDEF:
		//		decodeModel(data)
	}
}

// NewFrame initialize MocapFrame
func NewFrame() (frame *MocapFrame) {
	return &MocapFrame{
		frameNumber: 0,
		nMarkerSets: 0,
	}
}

// NewPacket initalize Packet
func NewPacket() (packet *Packet) {
	return &Packet{
		MessageID: 0,
		dataLen:   0,
	}
}

func (f MocapFrame) Format(fs fmt.State, c rune) {
	fmt.Fprintf(fs, "  FrameNumber: %d\n", f.frameNumber)
	fmt.Fprintf(fs, "  nMarkerSets: %d\n", f.nMarkerSets)
	fmt.Fprintf(fs, "  nMarkerSetsUndef: %d\n", f.nMarkerUnidentified)
	fmt.Fprintf(fs, "  nRigidBodies: %d\n", f.nRigidBodies)
	for _, v := range f.RigidBodies {
		fmt.Fprintf(fs, "  %v\n", v)
	}
	fmt.Fprintf(fs, "  nSkelton: %d\n", f.nSkeltons)
	fmt.Fprintf(fs, "  latency: %f\n", f.latency)
	fmt.Fprintf(fs, "  timecode: %d\n", f.timecode)

	return
}

func (p Packet) Format(fs fmt.State, c rune) {
	fmt.Fprintf(fs, "MessageId: %d\n", p.MessageID)
	fmt.Fprintf(fs, "Length: %d\n", p.dataLen)
	fmt.Fprint(fs, p.frame)

	return
}