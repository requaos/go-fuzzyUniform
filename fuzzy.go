package fuzzy

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unsafe"
)

const (
	blockBase           = 36
	byteSize            = 8
	blockHashModulo     = math.MaxInt32
	blockInnerSeparator = "/"
	factorSeparator     = ":"
	blocksSeparator     = "-"

	Similarity = 1 + iota
	ReverseSimilarity
	Maximum
	Minimum
	ArithmeticMean
	GeometricMean
)

var (
	ErrInvalidFactor     = errors.New("factor must be odd and greater than two")
	ErrNilData           = errors.New("data is nil")
	ErrEmptyHashString   = errors.New("hash string is empty")
	ErrInvalidHashString = errors.New("hash string does not fit the format factor:blocks")

	factorWithSepMaxChars = len(fmt.Sprint(blockHashModulo)) + len(factorSeparator)
	blockWithSepMaxChars  = blockMaxChars + len(blocksSeparator)
	blockIntMaxChars      = len(strconv.FormatInt(blockHashModulo, blockBase))
	blockMaxChars         = 2*blockIntMaxChars + len(blockInnerSeparator)
)

type SimilarityType int

type UniformFuzzyHash struct {
	factor   int
	dataSize int
	blocks   []UniformFuzzyHashBlock
}

// Compute is the main algorithm computation.
//
// factor is the relation between data length and the hash mean number of blocks.
// Must be greater than 2 and must be odd.
func Compute(data []byte, factor int) (*UniformFuzzyHash, error) {
	if data == nil {
		return nil, ErrNilData
	}
	err := checkFactor(factor)
	if err != nil {
		return nil, err
	}

	u := &UniformFuzzyHash{}

	// Set hash attributes
	u.factor = factor
	u.dataSize = len(data)
	u.blocks = []UniformFuzzyHashBlock{}

	// Size in bytes of the rolling window.
	// Size in bytes of factor + 5.
	windowSize := int(unsafe.Sizeof(factor)) + 5

	// Window size shifter.
	// Used to extract old data from the window.
	// (2 ^ (8 * windowSize)) % factor.
	windowSizeShifter := shiftBytesMod(windowSize, factor)

	// Window hash match value to produce a block.
	// Any number between 0 and factor - 1 should be valid.
	windowHashMatchValue := factor - 1

	// Rolling window hash.
	windowHash := int64(0)

	// Block hash.
	blockHash := int64(0)

	// Block starting position (0 based).
	blockStartingBytePosition := 0

	// Hash computation.
	for i := 0; i < len(data); i++ {

		// Unsigned datum.
		datum := int(data[i])

		// Window hash shift, new datum addition and old datum extraction.
		if i < windowSize {
			windowHash = ((windowHash << byteSize) + int64(datum)) % int64(factor)
		} else {
			windowHash = ((windowHash << byteSize) + int64(datum) - int64(int(data[i-windowSize])*windowSizeShifter)) % int64(factor)
		}

		// Due to the subtraction, the modulo result might be negative.
		if windowHash < 0 {
			windowHash += int64(factor)
		}

		// Block hash shift and new datum addition.
		blockHash = ((blockHash << byteSize) + int64(datum)) % blockHashModulo

		// Possible window hash match (block production).
		// Match is only checked if the initial window has already been computed.
		// Last data byte always produces a block.
		if (windowHash == int64(windowHashMatchValue) && i >= (windowSize-1)) || (i == (len(data) - 1)) {

			// New block addition.
			u.blocks = append(u.blocks, UniformFuzzyHashBlock{int(blockHash), blockStartingBytePosition, i})

			// Block hash reset.
			blockHash = 0

			// Next block starting byte position.
			blockStartingBytePosition = i + 1
		}
	}
	return u, nil
}

func checkFactor(factor int) error {
	if factor&1 != 1 || factor < 3 {
		return ErrInvalidFactor
	}
	return nil
}

func shiftBytesMod(bytesShift int, modulo int) int {
	ret := int64(1)
	for i := 0; i < bytesShift; i++ {
		ret = (ret << byteSize) % int64(modulo)
	}
	return int(ret)
}

func (u *UniformFuzzyHash) String() string {
	sb := strings.Builder{}
	sb.Grow(factorWithSepMaxChars + blockWithSepMaxChars + len(u.blocks))
	sb.WriteString(fmt.Sprint(u.factor))
	sb.WriteString(factorSeparator)
	for i := range u.blocks {
		if i != 0 {
			sb.WriteString(blocksSeparator)
		}
		sb.WriteString(u.blocks[i].String())
	}
	return sb.String()
}

func ParseHash(hashString string) (*UniformFuzzyHash, error) {
	if hashString == "" {
		return nil, ErrEmptyHashString
	}

	u := &UniformFuzzyHash{}

	splitIndex := strings.LastIndex(hashString, factorSeparator)
	if splitIndex < 0 {
		return nil, ErrInvalidHashString
	}

	factorString := hashString[:splitIndex]
	if factorString == "" {
		return nil, ErrInvalidHashString
	}

	i, err := strconv.ParseInt(factorString, 10, 64)
	if err != nil {
		return nil, err
	}

	err = checkFactor(int(i))
	if err != nil {
		return nil, err
	}
	u.factor = int(i)

	u.blocks = []UniformFuzzyHashBlock{}

	blockNumber := 0
	blockStartingBytePosition := 0

	blocksString := hashString[splitIndex+1:]
	if blocksString == "" {
		return nil, ErrInvalidHashString
	}

	splitIndex = 0
	lastSplitIndex := 0
	for splitIndex >= 0 {
		splitIndex = strings.Index(blocksString[lastSplitIndex:], blocksSeparator)
		blockString := ""
		if splitIndex >= 0 {
			blockString = blocksString[lastSplitIndex:splitIndex]
		} else {
			blockString = blocksString[lastSplitIndex:]
		}
		lastSplitIndex = splitIndex + len(blocksSeparator)

		hashBlock, err := ParseHashBlock(blockString, blockStartingBytePosition)
		if err != nil {
			return nil, err
		}
		u.blocks = append(u.blocks, *hashBlock)
		blockNumber++
		blockStartingBytePosition = hashBlock.blockEndingBytePosition + 1
	}

	u.dataSize = blockStartingBytePosition

	return u, nil
}

func Compare(a, b *UniformFuzzyHash) (float64, error) {
	if a == nil || b == nil {
		return 0, errors.New("hash cannot be nil")
	}

	// If the pointers are the same then it is the same object
	if a == b {
		return 1, nil
	}

	if a.factor != b.factor {
		return 0, errors.New("hash factors cannot be different")
	}

	if len(a.blocks) == 0 || len(b.blocks) == 0 {
		return 0, nil
	}

	sizeSum := int64(0)
	bString := b.String()
	for i := range a.blocks {
		if strings.Contains(bString, a.blocks[i].String()) {
			sizeSum += a.blocks[i].blockSize()
		}
	}

	return float64(sizeSum) / float64(a.dataSize), nil
}

func CompareSimilarity(a, b *UniformFuzzyHash, similarityType SimilarityType) (float64, error) {
	switch similarityType {
	case Similarity:
		return Compare(a, b)
	case ReverseSimilarity:
		return Compare(b, a)
	default:
		similarity, err := Compare(a, b)
		if err != nil {
			return 0, err
		}
		reverse, err := Compare(b, a)
		if err != nil {
			return 0, err
		}
		switch similarityType {
		case Maximum:
			if similarity > reverse {
				return similarity, nil
			} else {
				return reverse, nil
			}
		case Minimum:
			if similarity < reverse {
				return similarity, nil
			} else {
				return reverse, nil
			}
		case ArithmeticMean:
			return (similarity + reverse) / 2, nil
		case GeometricMean:
			return math.Sqrt(similarity * reverse), nil
		default:
			return similarity, nil
		}
	}
}
