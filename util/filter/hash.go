package filter


const (
	RS_HASH  = iota
	JS_HASH
	ELF_HASH
	BKD_HASH
	AP_HASH
	DJB_HASH
	SDB_HASH
	PJW_HASH
)
func GetHash(n int, bytes []byte) []int32 {
	var ret = make([]int32, 0)
	for i := 0; i < n; i++ {
		ret = append(ret, hash(i, bytes))
	}
	return ret
}

func hash(i int, bytes []byte) int32 {
	var ret int32
	switch i {
	case RS_HASH:
		ret = rsHash(bytes)
	case JS_HASH:
		ret = jsHash(bytes)
	case ELF_HASH:
		ret = elfHash(bytes)
	case BKD_HASH:
		ret = bkdRHash(bytes)
	case AP_HASH:
		ret = apHash(bytes)
	case DJB_HASH:
		ret = djbHash(bytes)
	case SDB_HASH:
		ret = sdbMHash(bytes)
	case PJW_HASH:
		ret = pjwHash(bytes)
	}
	return ret
}

func rsHash(bytes []byte) int32 {
	var (
		hash  = 0
		magic = 63689
		n     = len(bytes)
	)
	for i := 0; i < n; i++ {
		hash = (hash * magic) + int(bytes[i])
		magic = magic * 378551
	}
	return int32(hash)
}

func jsHash(bytes []byte) int32 {
	var hash = 1315423911
	for i := 0; i < len(bytes); i++ {
		hash ^= (hash << 5) + int(bytes[i]) + (hash >> 2);
	}
	return int32(hash)
}

func elfHash(bytes []byte) int32 {
	var (
		x    = 0
		hash = 0
	)
	for i := 0; i < len(bytes); i++ {
		hash = (hash << 4) + int(bytes[i])
		x = hash & 0xF0000000
		if x != 0 {
			hash ^= x >> 24
			hash &= ^x
		}
	}
	return int32(hash)
}

func bkdRHash(bytes []byte) int32 {
	var (
		seed = 131
		hash = 0
	)
	for i := 0; i < len(bytes); i++ {
		hash = (hash * seed) + int(bytes[i])
	}
	return int32(hash)
}

func apHash(bytes []byte) int32 {
	var hash = 0
	for i := 0; i < len(bytes); i++ {
		if i&1 == 0 {
			hash ^= hash<<7 ^ int(bytes[i]) ^ hash>>3
		} else {
			hash ^= ^(hash<<11 ^ int(bytes[i]) ^ hash>>5)
		}
	}
	return int32(hash)
}

func djbHash(bytes []byte) int32 {
	var hash = 5381
	for i := 0; i < len(bytes); i++ {
		hash = ((hash << 5) + hash) + int(bytes[i])
	}
	return int32(hash)
}

func sdbMHash(bytes []byte) int32 {
	var hash = 0
	for i := 0; i < len(bytes); i++ {
		hash = int(bytes[i]) + (hash << 6) + (hash << 16) - hash
	}
	return int32(hash)
}

func pjwHash(bytes []byte) int32 {
	var (
		t             uint64 = 0
		hash          uint64 = 0
		bitsInteger   uint64 = 4 * 8
		threeQuarters uint64 = (bitsInteger * 3) / 4
		oneEighth     uint64 = bitsInteger / 8
		highBits      uint64 = 0xFFFFFFFF << (bitsInteger - oneEighth)
	)
	for i := 0; i < len(bytes); i++ {
		hash = (hash << oneEighth) + uint64(bytes[i]);
		if t = hash & highBits; t != 0 {
			hash = (hash ^ (t >> threeQuarters)) & (^highBits)
		}
	}
	return int32(hash)
}
