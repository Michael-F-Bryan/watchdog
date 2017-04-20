package main

import "io"

// ## Login Handshake
//
//
// The client sends a packet containing their username and password (and one extra
// unknown field), where each is encoded as a UTF-16 string, prefixed with the
// string length.
//
// Between each field is a delimiter of `0x01 0x00 0x00 0x00`.
// The first 4 bytes indicate the total length of the rest of the packet as a
// 32-bit integer.
//
// Client sends:
//
// Hex                                                Ascii             Field
// ===                                                =====             =====
//
// 48 00 00 00                                        H...             packet len
// 01 00 00 00                                        ....             delimiter
//
// 1a 00 00 00                                        ....             user length
// 4d 00 69 00 63 00 68 00                            M.i.c.h.         username
// 61 00 65 00 6c 00 5f 00                            a.e.l._.         ""
// 42 00 72 00 79 00 61 00                            B.r.y.a.         ""
// 6e 00                                              n.               ""
// 01 00 00 00                                        ....             delimiter
//
// 06 00 00 00                                        ....             pass length
// 6d 00 69 00 63 00                                  m.i.c.           password
// 01 00 00 00                                        ....             delimiter
//
// 00 00 00 00  36 00 00 00 32 45 0a 85               .... 6...2E..    ?
// 00 00 00 00                                        ....             ?
// 01 00 00 00                                        ....             delimiter

type PdmPacket struct {
	Username string
	Password string
	Payload  []byte
}

func NewPdmPacket(username, password string) PdmPacket {
	return PdmPacket{
		Username: username,
		Password: password,
		Payload:  []byte{0x36, 0x00, 0x00, 0x00, 0x32, 0x45, 0x0a, 0x85},
	}
}

func (pkt PdmPacket) Serialize(w io.Writer) error {
	return nil
}
