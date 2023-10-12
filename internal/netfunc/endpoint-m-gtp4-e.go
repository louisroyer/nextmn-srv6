// Copyright 2023 Louis Royer and the NextMN-SRv6 contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package netfunc

import (
	"fmt"
	"net/netip"

	"github.com/google/gopacket/layers"
	gopacket_srv6 "github.com/louisroyer/gopacket-srv6"
	"github.com/nextmn/srv6/internal/constants"
	"github.com/nextmn/srv6/internal/mup"
)

type EndpointMGTP4E struct {
	BaseHandler
}

func NewEndpointMGTP4E(prefix netip.Prefix, ttl uint8, hopLimit uint8) *EndpointMGTP4E {
	return &EndpointMGTP4E{
		BaseHandler: NewBaseHandler(prefix, ttl, hopLimit),
	}
}

// Get IPv6 Destination Address Fields from Packet
func (e EndpointMGTP4E) ipv6DAFields(p *Packet) (*mup.EndMGTP4EIPv6DstFields, error) {
	layerIPv6 := p.Layer(layers.LayerTypeIPv6)
	if layerIPv6 == nil {
		return nil, fmt.Errorf("Malformed IPv6 packet")
	}
	// get destination address
	dstSlice := layerIPv6.(*layers.IPv6).NetworkFlow().Dst().Raw()
	if dst, err := mup.NewEndMGTP4EIPv6DstFields(dstSlice, uint(e.Prefix().Bits())); err != nil {
		return nil, err
	} else {
		return dst, nil
	}
}

// Get IPv6 Source Address Fields from Packet
func (e EndpointMGTP4E) ipv6SAFields(p *Packet) (*mup.EndMGTP4EIPv6SrcFields, error) {
	layerIPv6 := p.Layer(layers.LayerTypeIPv6)
	if layerIPv6 == nil {
		return nil, fmt.Errorf("Malformed IPv6 packet")
	}
	// get destination address
	srcSlice := layerIPv6.(*layers.IPv6).NetworkFlow().Src().Raw()
	if src, err := mup.NewEndMGTP4EIPv6SrcFields(srcSlice); err != nil {
		return nil, err
	} else {
		return src, nil
	}
}

// Handle a packet
func (e EndpointMGTP4E) Handle(packet []byte) ([]byte, error) {
	pqt, err := NewIPv6Packet(packet)
	if err != nil {
		return nil, err
	}
	if err := e.CheckDAInPrefixRange(pqt); err != nil {
		return nil, err
	}

	// SRH is optionnal (unless the endpoint is configured to accept only packet with HMAC TLV)
	if layerSRH := pqt.Layer(gopacket_srv6.LayerTypeIPv6Routing); layerSRH != nil {
		srh := layerSRH.(*gopacket_srv6.IPv6Routing)
		// RFC 9433 section 6.6. End.M.GTP4.E
		// S01. When an SRH is processed {
		// S02.   If (Segments Left != 0) {
		// S03.      Send an ICMP Parameter Problem to the Source Address with
		//              Code 0 (Erroneous header field encountered) and
		//              Pointer set to the Segments Left field,
		//              interrupt packet processing, and discard the packet.
		// S04.   }
		if srh.SegmentsLeft != 0 {
			// TODO: Send ICMP response
			return nil, fmt.Errorf("Segments Left is not zero")
		}
		// TODO: check HMAC

		// S05.   Proceed to process the next header in the packet
		// S06. }
	} //TODO: else if HMAC -> error: no SRH

	// S01. Store the IPv6 DA and SA in buffer memory
	ipv6SA, err := e.ipv6SAFields(pqt)
	if err != nil {
		return nil, err
	}
	ipv6DA, err := e.ipv6DAFields(pqt)
	if err != nil {
		return nil, err
	}

	// S02. Pop the IPv6 header and all its extension headers
	payload, err := pqt.PopIPv6Headers()
	if err != nil {
		return nil, err
	}

	// S03. Push a new IPv4 header with a UDP/GTP-U header
	// S04. Set the outer IPv4 SA and DA (from buffer memory)
	// S05. Set the outer Total Length, DSCP, Time To Live, and
	//      Next Header fields
	ipv4 := layers.IPv4{
		// IPv4
		Version: 4,
		// Next Header: UDP
		Protocol: layers.IPProtocolUDP,
		// Fragmentation is inefficient and should be avoided (TS 129.281 section 4.2.2)
		// It is recommended to set the default inner MTU size instead.
		Flags: layers.IPv4DontFragment,
		// Destination IP from buffer
		SrcIP: ipv6SA.IPv4(),
		// Source IP from buffer
		DstIP: ipv6DA.IPv4(),
		// TOS = DSCP + ECN
		// We copy the QFI into the DSCP Field
		TOS: ipv6DA.QFI() << 2,
		// TTL from tun config
		TTL: e.TTL(),
		// other fields are initialized at zero
		// cheksum, and length are computed at serialization

	}

	udp := layers.UDP{
		// Source Port
		SrcPort: ipv6SA.UDPPortNumber(),
		SrcPort: constants.GTPU_PORT_INT,
		// cheksum, and length are computed at serialization
	}

	// S06.    Set the GTP-U TEID (from buffer memory)
	pduSessionContainerLength := 0 // FIXME
	gtpu := layers.GTPv1U{
		// Version should always be set to 1
		Version: 1,
		// TS 128281:
		// > This bit is used as a protocol discriminator between
		// > GTP (when PT is '1') and GTP' (whenPT is '0').
		ProtocolType: 1,
		// We use extension header "PDU Session Container"
		ExtensionHeaderFlag: true,
		GTPExtensionHeaders: nil, // FIXME
		// TS 128281:
		// > Since the use of Sequence Numbers is optional for G-PDUs, the PGW,
		// > SGW, ePDG, eNodeB and TWAN should set the flag to '0'.
		SequenceNumberFlag: false,
		// message type: G-PDU
		MessageType: constants.GTPU_MESSAGE_TYPE_GPDU,
		TEID:        ipv6DA.PDUSessionID(),
		// TS 128281:
		// > This field indicates the length in octets of the payload, i.e. the rest of the packet following the mandatory
		// > part of the GTP header (that is the first 8 octets). The Sequence Number, the N-PDU Number or any Extension
		// > headers shall be considered to be part of the payload, i.e. included in the length count
		MessageLength: uint16(len(payload.LayerContents()) + pduSessionContainerLength),
	}
	// create buffer for the packet
	//buf := gopacket.NewSerializeBuffer()
	// initialize buffer with the payload
	// Initial content of the buffer : [ ]
	// Updated content of the buffer : [ PDU ]
	//err = gopacket.Payload(pdu).SerializeTo(buf, gopacket.SerializeOptions{
	//	FixLengths:       true,
	//	ComputeChecksums: true,
	//})

	// S07.    Submit the packet to the egress IPv4 FIB lookup for
	//            transmission to the new destination

	// extract TEID from destination address
	// destination address is formed as follow : [ SID (netsize bits) + IPv4 DA (only if ipv4) + ArgsMobSession ]
	//	dstarray := dst.As16()
	//	offset := 0
	//	if s.gtpIPVersion == 4 {
	//		offset = 32 / 8
	//	}
	// TODO: check segments left = 1, and if not send ICMP Parameter Problem to the Source Address (code 0, pointer to SegemntsLeft field), and drop the packet
	//	args, err := mup.ParseArgsMobSession(dstarray[(s.netsize/8)+offset:])
	//	if err != nil {
	//		return err
	//	}
	//	teid := args.PDUSessionID()
	// retrieve nextGTPNode (SHR[0])

	//	nextGTPNode := ""
	//	if s.gtpIPVersion == 6 {
	//		// workaround: enforce use of gopacket_srv6 functions
	//		shr := gopacket.NewPacket(pqt.Layers()[1].LayerContents(), gopacket_srv6.LayerTypeIPv6Routing, gopacket.Default).Layers()[0].(*gopacket_srv6.IPv6Routing)
	//		log.Println("layer type", pqt.Layers()[1].LayerType())
	//		log.Println("RoutingType", shr.RoutingType)
	//		log.Println("LastEntry:", shr.LastEntry)
	//		log.Println("sourceRoutingIPs len:", len(shr.SourceRoutingIPs))
	//		log.Println("sourceRoutingIPs[0]:", shr.SourceRoutingIPs[0])
	//		nextGTPNode = fmt.Sprintf("[%s]:%s", shr.SourceRoutingIPs[0].String(), GTPU_PORT)
	//	} else {
	//		// IPv4
	//		ip_arr := dstarray[s.netsize/8 : (s.netsize/8)+4]
	//		ipv4_address := net.IPv4(ip_arr[0], ip_arr[1], ip_arr[2], ip_arr[3])
	//		nextGTPNode = fmt.Sprintf("%s:%s", ipv4_address, GTPU_PORT)
	//	}
	//	raddr, err := net.ResolveUDPAddr("udp", nextGTPNode)
	//	if err != nil {
	//		log.Println("Error while resolving ", nextGTPNode, "(remote node)")
	//		return nil
	//	}
	// retrieve payload
	//			pdu := pqt.Layers()[2].LayerContents() // We expect the packet to contains the following layers [ IPv6 Header (0) + IPv6Routing Ext Header (1) + PDU (2) ]
	//			// Search for existing Uconn with this peer and use it
	//			if s.uConn[nextGTPNode] == nil {
	//				// Start uConn with this peer
	//				ch := make(chan bool)
	//				go s.StartUconn(ch, nextGTPNode, raddr)
	//				_ = <-ch
	//			}
	//			s.uConn[nextGTPNode].WriteToGTP(teid, pdu, raddr)

	// create gopacket
	return nil, fmt.Errorf("TODO")
}
