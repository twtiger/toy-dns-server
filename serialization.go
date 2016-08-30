package nameserver

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func flattenBytes(i ...interface{}) (b []byte) {
	for _, e := range i {
		switch k := e.(type) {
		case string:
			b = append(b, []byte(k)...)
		case int:
			b = append(b, byte(k))
		case byte:
			b = append(b, k)
		case []byte:
			b = append(b, k...)
		case uint16:
			b = append(b, []byte{0, byte(k)}...)
		default:
			panic(fmt.Sprintf("cannot flatten: %#v", e))
		}
	}
	return b
}

func (m *message) serialize() ([]byte, error) {
	q, err := m.query.serialize()
	if err != nil {
		return nil, err
	}

	return flattenBytes(serializeHeaders(m.header), q, serializeAnswer(m.answers)), nil
}

func (q *query) serialize() ([]byte, error) {
	l, err := serializeLabels(q.qname)
	if err != nil {
		return nil, err
	}

	return flattenBytes(l, serializeUint16(uint16(q.qtype)), serializeUint16(uint16(q.qclass))), nil
}

func (r *record) serialize() (b []byte) {
	l, _ := serializeLabels(r.Name)

	b = flattenBytes(l, serializeUint16(uint16(r.Type)), serializeUint16(uint16(r.Class)), serializeUint32(uint32(r.TTL)), serializeUint16(r.RDLength), r.RData)
	return
}

func (l label) serialize() (b []byte) {
	b = append(b, byte(len(l)))
	b = append(b, []byte(l)...)
	return
}

func serializeLabels(l []label) ([]byte, error) {
	var b []byte
	if len(l) == 0 {
		return nil, errors.New("no labels to serialize")
	}

	for _, e := range l {
		b = append(b, e.serialize()...)
	}
	b = append(b, 0)
	return b, nil
}

func serializeUint16(i uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return b
}

func serializeUint32(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}

func serializeAnswer(r []*record) (b []byte) {
	for _, e := range r {
		b = append(b, e.serialize()...)
	}
	return
}

func serializeHeaders(h *header) (b []byte) {
	IDinBytes := serializeUint16(h.id)
	b = []byte(IDinBytes)
	b = append(b, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}...)
	return
}
