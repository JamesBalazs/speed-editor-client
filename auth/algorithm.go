// Based on the authentication algorithm by Sylvain Munaut <tnt@246tNt.com> (Apache-2.0)
// https://github.com/smunaut/blackmagic-misc/blob/master/bmd.py
// SPDX-License-Identifier: Apache-2.0

package auth

var (
	AuthEvenTable = []uint64{
		0x3ae1206f97c10bc8,
		0x2a9ab32bebf244c6,
		0x20a6f8b8df9adf0a,
		0xaf80ece52cfc1719,
		0xec2ee2f7414fd151,
		0xb055adfd73344a15,
		0xa63d2e3059001187,
		0x751bf623f42e0dde,
	}

	AuthOddTable = []uint64{
		0x3e22b34f502e7fde,
		0x24656b981875ab1c,
		0xa17f3456df7bf8c3,
		0x6df72e1941aef698,
		0x72226f011e66ab94,
		0x3831a3c606296b42,
		0xfd7ff81881332c89,
		0x61a3f6474ff236c6,
	}

	Mask = uint64(0xa79a63f585d37bf0)
)

func rol8(v uint64) uint64 {
	return ((v << 56) | (v >> 8)) & 0xffffffffffffffff
}

func rol8n(v, n uint64) uint64 {
	for _ = range n {
		v = rol8(v)
	}

	return v
}

func calculateChallengeResponse(challenge uint64) uint64 {
	n := challenge & 7
	v := rol8n(challenge, n)

	var k uint64
	if (v & 1) == ((0x78 >> n) & 1) {
		k = AuthEvenTable[n]
	} else {
		v = v ^ rol8(v)
		k = AuthOddTable[n]
	}

	return v ^ (rol8(v) & Mask) ^ k
}
