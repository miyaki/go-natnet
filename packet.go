package natnet

import (
	"encoding/binary"
	"fmt"
	"io"
)

// contants for natnet packet handling
const (
	maxNAMELENGTH = 256
	maxPACKETSIZE = 100000

	natPING                = 0
	natPINGRESPONSE        = 1
	natREQUEST             = 2
	natRESPONSE            = 3
	natREQUESTMODELDEF     = 4
	natMODELDEF            = 5
	natREQUESTFRAMEOFDATA  = 6
	natFRAMEOFDATA         = 7
	natMESSAGESTRING       = 8
	natUNRECOGNIZEDREQUEST = 100
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

// Quaternion4f holds orientation
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
//type MarkerSet struct {
//	x, y, z float32
//}

//Skelton not yet supported
type Skelton struct {
	ID           int32
	nRigidBodies int32
	RigidBodies  []RigidBody
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

//TODO: interface PacketReader

// Decode for RigidBody
func (rb *RigidBody) Decode(reader io.Reader) {
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
		rb.MarkerData[i].Decode(reader)
	}

	//if VERSION >=2.0
	rb.MarkerID = make([]uint32, rb.nRigidMarkers)
	for i := range rb.MarkerID {
		binary.Read(reader, binary.LittleEndian, &rb.MarkerID[i])
	}
	rb.MarkerSize = make([]float32, rb.nRigidMarkers)
	for i := range rb.MarkerSize {
		binary.Read(reader, binary.LittleEndian, &rb.MarkerSize[i])
	}

	//if VERSION >= 2.6
	{
		//binary.Read(read, binary.LittleEndian, &rb.params)
		//bTrackingValid := rb.params & 0x01
	}

	binary.Read(reader, binary.LittleEndian, &rb.MarkerError)
}

// Decode for Skelton
func (s *Skelton) Decode(reader io.Reader) {

}

/*
func (ms *MarkerSet) Decode(reader io.Reader) {
	binary.Read(reader, binary.LittleEndian, &ms.x)
	binary.Read(reader, binary.LittleEndian, &ms.y)
	binary.Read(reader, binary.LittleEndian, &ms.z)
}
*/

// Decode for Point3f
func (p *Point3f) Decode(reader io.Reader) {
	binary.Read(reader, binary.LittleEndian, &p.x)
	binary.Read(reader, binary.LittleEndian, &p.y)
	binary.Read(reader, binary.LittleEndian, &p.z)
}

// Decode for MocapFrame
func (mf *MocapFrame) Decode(reader io.Reader) {
	binary.Read(reader, binary.LittleEndian, &mf.frameNumber)

	binary.Read(reader, binary.LittleEndian, &mf.nMarkerSets)
	mf.MarkerSets = make([]Point3f, mf.nMarkerSets)
	for i := range mf.MarkerSets {
		mf.MarkerSets[i].Decode(reader)
	}

	binary.Read(reader, binary.LittleEndian, &mf.nMarkerUnidentified)
	mf.MarkerUnidentified = make([]Point3f, mf.nMarkerUnidentified)
	for i := range mf.MarkerUnidentified {
		mf.MarkerUnidentified[i].Decode(reader)
	}

	binary.Read(reader, binary.LittleEndian, &mf.nRigidBodies)
	mf.RigidBodies = make([]RigidBody, mf.nRigidBodies)
	for i := range mf.RigidBodies {
		mf.RigidBodies[i].Decode(reader)
	}

	// if VERSION >= 2.1
	binary.Read(reader, binary.LittleEndian, &mf.nSkeltons)
	mf.Skeltons = make([]Skelton, mf.nSkeltons)
	for i := range mf.Skeltons {
		mf.Skeltons[i].Decode(reader)
	}

	//if VERSION >= 2.3
	//labeled markers
	{
		//id
		//pos
		//size

		//if >= 2.6
		{
			//marker params
			//param
			//binary.Read(reader, binary.LittleEndian, &mf.)

			//bOccluded
			//bPCSolved
			//bModelSolved
		}
	}

	binary.Read(reader, binary.LittleEndian, &mf.latency)

	binary.Read(reader, binary.LittleEndian, &mf.timecode)
	hour := (mf.timecode >> 24) & 0xff
	minute := (mf.timecode >> 16) & 0xff
	second := (mf.timecode >> 8) & 0xff
	frame := mf.timecode & 0xff
	fmt.Println("%d:%d:%d:%d", hour, minute, second, frame)

	//if >= 2.3
	//subtimecodesub
	//binary.Read(reader, binary.LittleEndian, &mf.timecodesub)

	//if >= 2.6
	//timstamp
	//binary.Read(reader, binary.LittleEndian, &mf.timestamp) // int32
	//frameparams
	//binary.Read(reader, binary.LittleEndian, &mf.params) // int16
	//bIsRecording          := mf.params & 0x01
	//bTrackedModeIsChanged := mf.params & 0x02

}

// Decode for Packet
func (p *Packet) Decode(reader io.Reader) {
	p.MessageID = int16(0)
	binary.Read(reader, binary.LittleEndian, &p.MessageID)

	p.dataLen = uint16(0)
	binary.Read(reader, binary.LittleEndian, &p.dataLen)

	switch p.MessageID {
	default:
		fmt.Printf("unknown MessageId: %d\n", p.MessageID)
	case natFRAMEOFDATA:
		p.frame = NewFrame()
		p.frame.Decode(reader)
	case natMODELDEF:
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

// Format for MocapFrame
func (mf MocapFrame) Format(fs fmt.State, c rune) {
	fmt.Fprintf(fs, "  FrameNumber: %d\n", mf.frameNumber)
	fmt.Fprintf(fs, "  nMarkerSets: %d\n", mf.nMarkerSets)
	fmt.Fprintf(fs, "  nMarkerSetsUndef: %d\n", mf.nMarkerUnidentified)
	fmt.Fprintf(fs, "  nRigidBodies: %d\n", mf.nRigidBodies)
	for _, v := range mf.RigidBodies {
		fmt.Fprintf(fs, "  %v\n", v)
	}
	fmt.Fprintf(fs, "  nSkelton: %d\n", mf.nSkeltons)
	fmt.Fprintf(fs, "  latency: %f\n", mf.latency)
	fmt.Fprintf(fs, "  timecode: %d\n", mf.timecode)

	return
}

// Format for Packet
func (p Packet) Format(fs fmt.State, c rune) {
	fmt.Fprintf(fs, "MessageId: %d\n", p.MessageID)
	fmt.Fprintf(fs, "Length: %d\n", p.dataLen)
	fmt.Fprint(fs, p.frame)

	return
}
