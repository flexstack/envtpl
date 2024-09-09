// Credit: https://github.com/matoous/go-nanoid
// Credit: https://github.com/ai/nanoid
package nanoid

import (
	"crypto/rand"
	"math"
)

// defaultAlphabet is the alphabet used for ID characters by default.
var (
	AlphabetDefault = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	AlphabetBase64  = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	AlphabetHex     = []rune("0123456789abcdef")
	AlphabetAscii85 = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!#$%&()*+-;<=>?@^_`{|}~")
)

// getMask generates bit mask used to obtain bits from the random bytes that are used to get index of random character
// from the alphabet. Example: if the alphabet has 6 = (110)_2 characters it is sufficient to use mask 7 = (111)_2
func getMask(alphabetSize int) int {
	for i := 1; i <= 8; i++ {
		mask := (2 << uint(i)) - 1
		if mask >= alphabetSize-1 {
			return mask
		}
	}
	return 0
}

// Generate is a low-level function to change alphabet and ID size.
func Generate(alphabet string, size int) string {
	chars := []rune(alphabet)
	mask := getMask(len(chars))
	// estimate how many random bytes we will need for the ID, we might actually need more but this is tradeoff
	// between average case and worst case
	ceilArg := 1.6 * float64(mask*size) / float64(len(alphabet))
	step := int(math.Ceil(ceilArg))
	id := make([]rune, size)
	bytes := make([]byte, step)

	for j := 0; ; {
		_, err := rand.Read(bytes)

		if err != nil {
			panic(err)
		}

		for i := 0; i < step; i++ {
			currByte := bytes[i] & byte(mask)
			if currByte < byte(len(chars)) {
				id[j] = chars[currByte]
				j++
				if j == size {
					return string(id[:size])
				}
			}
		}
	}
}

func New(size int, alphabet ...[]rune) string {
	var chars []rune
	if len(alphabet) == 0 {
		chars = AlphabetDefault
	} else {
		chars = alphabet[0]
	}

	return Generate(string(chars), size)
}
