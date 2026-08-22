package models

import (
	"math/rand"
	"sync/atomic"
	"time"
)

// lastNano holds the most recently issued nanosecond value used by
// SortableShortUUID. It guarantees strict monotonicity so that rapid,
// back-to-back calls (e.g. a batch insert of order lines) never collide on the
// same UnixNano() tick — which previously produced duplicate primary keys that
// SQLite silently dropped.
var lastNano int64

const numerals = "0123456789ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"

func ShortUUID(length int) string {
	// seed a new source of randon numbers
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	// get int value of uuid
	uuid := r.Int63()
	numerals := alphabet

	// encode it until we get the length we want
	shortuuid := encode(uuid, len(numerals), numerals)
	for len(shortuuid) < length {
		uuid = r.Int63()
		shortuuid = encode(uuid, len(numerals), numerals)
		// if length is too long, append a second encoded uuid
		if length > 11 { // 11 is the max length using alphabet, probably =)
			shortuuid += encode(r.Int63(), len(numerals), numerals)
		}
	}
	return shortuuid[0:length]
}

// SortableShortUUID generates a sortable, unique, time-ordered string. It is
// safe under rapid concurrent/looped calls: a monotonic counter ensures each
// call yields a strictly larger value than the previous one even when the
// system clock has not advanced between calls.
func SortableShortUUID() string {
	for {
		prev := atomic.LoadInt64(&lastNano)
		next := time.Now().UnixNano()
		if next <= prev {
			next = prev + 1
		}
		if atomic.CompareAndSwapInt64(&lastNano, prev, next) {
			return encode(next, len(alphabet), alphabet)
		}
	}
}

func encode(num int64, base int, numerals string) string {
	var result string
	if numerals == "" {
		numerals = alphabet
	}
	for num > 0 {
		result = string(numerals[num%int64(base)]) + result
		num /= int64(base)
	}
	return result
}
