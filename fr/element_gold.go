// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fr

import "math/bits"

// Arithmetic Modulo ϕ₆(2³²) = q

// Goldilocks element
type GElement [1]uint64

const (
	gq  uint64 = 18446744069414584321 // ϕ₆(2³²)
	lsb uint64 = 4294967295           // 2³²-1
	msb uint64 = 18446744069414584320 // 2³²-1 << 32
)

// smallerThanModulus returns true if z < q
// This is not constant time
func (z *GElement) smallerThanModulus() bool {
	return z[0] < gq
}

// Sub z = x - y (mod q)
func (z *GElement) Sub(x, y *GElement) *GElement {
	var b uint64
	z[0], b = bits.Sub64(x[0], y[0], 0)
	if b != 0 {
		z[0] += gq
	}
	return z
}

// Add z = x + y (mod q)
func (z *GElement) Add(x, y *GElement) *GElement {

	var carry uint64
	z[0], carry = bits.Add64(x[0], y[0], 0)
	if carry != 0 || z[0] >= q {
		z[0] -= q
	}
	return z
}

// Double z = x + x (mod q), aka Lsh 1
func (z *GElement) Double(x *GElement) *GElement {
	if x[0]&(1<<63) == (1 << 63) {
		// if highest bit is set, then we have a carry to x + x, we shift and subtract q
		z[0] = (x[0] << 1) - q
	} else {
		// highest bit is not set, but x + x can still be >= q
		z[0] = (x[0] << 1)
		if z[0] >= q {
			z[0] -= q
		}
	}
	return z
}

// SetZero z = 0
func (z *GElement) SetZero() *GElement {
	z[0] = 0
	return z
}

// SetOne z = 1 (in Montgomery form)
func (z *GElement) SetOne() *GElement {
	z[0] = 1
	return z
}

// IsZero returns z == 0
func (z *GElement) IsZero() bool {
	return (z[0]) == 0
}

// NewElement returns a new Element from a uint64 value
//
// it is equivalent to
//
//	var v Element
//	v.SetUint64(...)
func GNewElement(v uint64) GElement {
	z := GElement{v}
	return z
}

// SetUint64 sets z to v and returns z
func (z *GElement) SetUint64(v uint64) *GElement {
	//  sets z LSB to v (non-Montgomery form) and convert z to Montgomery form
	*z = GElement{v}
	return z
}

// SetInt64 sets z to v and returns z
func (z *GElement) SetInt64(v int64) *GElement {

	// absolute value of v
	m := v >> 63
	z.SetUint64(uint64((v ^ m) - m))

	if m != 0 {
		// v is negative
		z.Neg(z)
	}

	return z
}

// Neg z = q - x
func (z *GElement) Neg(x *GElement) *GElement {
	if x.IsZero() {
		z.SetZero()
		return z
	}
	z[0] = q - x[0]
	return z
}

// Inverse z = x⁻¹ (mod q)
//
// if x == 0, sets and returns z = x
func (z *GElement) Inverse(x *GElement) *GElement {

	// Algorithm 16 in "Efficient Software-Implementation of Finite Fields with Applications to Cryptography"

	if x.IsZero() {
		z.SetZero()
		return z
	}

	var r, s, u, v uint64
	u = gq
	s = 1
	r = 0
	v = x[0]

	var carry, borrow uint64

	for (u != 1) && (v != 1) {
		for v&1 == 0 {
			v >>= 1
			if s&1 == 0 {
				s >>= 1
			} else {
				s, carry = bits.Add64(s, q, 0)
				s >>= 1
				if carry != 0 {
					s |= (1 << 63)
				}
			}
		}
		for u&1 == 0 {
			u >>= 1
			if r&1 == 0 {
				r >>= 1
			} else {
				r, carry = bits.Add64(r, q, 0)
				r >>= 1
				if carry != 0 {
					r |= (1 << 63)
				}
			}
		}
		if v >= u {
			v -= u
			s, borrow = bits.Sub64(s, r, 0)
			if borrow == 1 {
				s += q
			}
		} else {
			u -= v
			r, borrow = bits.Sub64(r, s, 0)
			if borrow == 1 {
				r += q
			}
		}
	}

	if u == 1 {
		z[0] = r
	} else {
		z[0] = s
	}

	return z
}

// // Mul
// func (z *GElement) Mul(x, y  *GElement) *GElement {

// }
