package msgpack

import (
	"io"
	"strconv"
	"unsafe"
)


func unpackArraySane(reader io.Reader, nelems uint) (v []interface{}, n int, err error) {
	retval := make([]interface{}, nelems)
	nbytesread := 0
	var i uint
	for i = 0; i < nelems; i++ {
		v, n, e := UnpackSane(reader)
		nbytesread += n
		if e != nil {
			return nil, nbytesread, e
		}
		retval[i] = v
	}
	return retval, nbytesread, nil
}


func unpackMapSane(reader io.Reader, nelems uint) (v map[interface{}]interface{}, n int, err error) {
	retval := make(map[interface{}]interface{})
	nbytesread := 0
	var i uint
	for i = 0; i < nelems; i++ {
		k, n, e := UnpackSane(reader)
		nbytesread += n
		if e != nil {
			return nil, nbytesread, e
		}
		v, n, e := UnpackSane(reader)
		nbytesread += n
		if e != nil {
			return nil, nbytesread, e
		}
		if str, ok := k.([]uint8); ok {
			retval[string(str)] = v
		} else {
			retval[k] = v
		}
	}
	return retval, nbytesread, nil
}

func unpack_sane(reader io.Reader) (v interface{}, n int, err error) {
	var retval interface{}
	var nbytesread int = 0

	c, e := readByte(reader)
	if e != nil {
		return nil, 0, e
	}
	nbytesread += 1
	if c < 0x80 || c >= 0xe0 {
		retval = int8(c)
	} else if c >= 0x80 && c <= 0x8f {
		retval, n, e = unpackMapSane(reader, uint(c&0xf))
		nbytesread += n
		if e != nil {
			return nil, nbytesread, e
		}
		nbytesread += n
	} else if c >= 0x90 && c <= 0x9f {
		retval, n, e = unpackArraySane(reader, uint(c&0xf))
		nbytesread += n
		if e != nil {
			return nil, nbytesread, e
		}
		nbytesread += n
	} else if c >= 0xa0 && c <= 0xbf {
		data := make([]byte, c&0x1f)
		n, e := reader.Read(data)
		nbytesread += n
		if e != nil {
			return nil, nbytesread, e
		}
		retval = data
	} else {
		switch c {
		case 0xc0:
			retval = nil
		case 0xc2:
			retval = false
		case 0xc3:
			retval = true
		case 0xca:
			data, n, e := readUint32(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval = *(*float32)(unsafe.Pointer(&data))
		case 0xcb:
			data, n, e := readUint64(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval = *(*float64)(unsafe.Pointer(&data))
		case 0xcc:
			data, e := readByte(reader)
			if e != nil {
				return nil, nbytesread, e
			}
			retval = uint8(data)
			nbytesread += 1
		case 0xcd:
			data, n, e := readUint16(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval = data
		case 0xce:
			data, n, e := readUint32(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval = data
		case 0xcf:
			data, n, e := readUint64(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval = data
		case 0xd0:
			data, e := readByte(reader)
			if e != nil {
				return nil, nbytesread, e
			}
			retval = int8(data)
			nbytesread += 1
		case 0xd1:
			data, n, e := readInt16(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval = data
		case 0xd2:
			data, n, e := readInt32(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval = data
		case 0xd3:
			data, n, e := readInt64(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval = data
		case 0xda:
			nbytestoread, n, e := readUint16(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			data := make([]byte, nbytestoread)
			n, e = reader.Read(data)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval = data
		case 0xdb:
			nbytestoread, n, e := readUint32(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			data := make([]byte, nbytestoread)
			n, e = reader.Read(data)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval = data
		case 0xdc:
			nelemstoread, n, e := readUint16(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval, n, e = unpackArraySane(reader, uint(nelemstoread))
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
		case 0xdd:
			nelemstoread, n, e := readUint32(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval, n, e = unpackArraySane(reader, uint(nelemstoread))
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
		case 0xde:
			nelemstoread, n, e := readUint16(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval, n, e = unpackMapSane(reader, uint(nelemstoread))
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
		case 0xdf:
			nelemstoread, n, e := readUint32(reader)
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
			retval, n, e = unpackMapSane(reader, uint(nelemstoread))
			nbytesread += n
			if e != nil {
				return nil, nbytesread, e
			}
		default:
			panic("unsupported code: " + strconv.Itoa(int(c)))
		}
	}
	return retval, nbytesread, nil
}

// Reads a value from the reader, unpack and returns it.
func UnpackSane(reader io.Reader) (v interface{}, n int, err error) {
	return unpack_sane(reader)
}

