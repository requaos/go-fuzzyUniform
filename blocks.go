package fuzzy

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidBlockString    = errors.New("block string does not fit the format of blockHash/blockSize")
	ErrUnparsableBlockString = errors.New("block hash/size is not able to be parsed")
)

type UniformFuzzyHashBlock struct {
	blockHash, blockStartingBytePosition, blockEndingBytePosition int
}

func (b *UniformFuzzyHashBlock) String() string {
	sb := strings.Builder{}
	sb.Grow(blockMaxChars)
	sb.WriteString(strconv.FormatInt(int64(b.blockHash), blockBase))
	sb.WriteString(blocksSeparator)
	sb.WriteString(strconv.FormatInt(b.blockSize(), blockBase))
	return sb.String()
}

func (b *UniformFuzzyHashBlock) blockSize() int64 {
	return int64(b.blockEndingBytePosition) - int64(b.blockStartingBytePosition) + 1
}

func ParseHashBlock(blockString string, blockStartingBytePosition int) (*UniformFuzzyHashBlock, error) {
	b := &UniformFuzzyHashBlock{}

	splitIndex := strings.LastIndex(blockString, blockInnerSeparator)
	if splitIndex < 0 {
		return nil, ErrInvalidBlockString
	}

	blockHashString := blockString[:splitIndex]
	if blockHashString == "" {
		return nil, ErrInvalidBlockString
	}

	i, err := strconv.ParseInt(blockHashString, blockBase, 64)
	if err != nil {
		return nil, err
	}

	if i < 0 || i >= blockHashModulo {
		return nil, ErrUnparsableBlockString
	}
	b.blockHash = int(i)

	blockHashSize := blockString[splitIndex+1:]
	if blockHashSize == "" {
		return nil, ErrInvalidBlockString
	}

	blockSize, err := strconv.ParseInt(blockHashSize, blockBase, 64)
	if err != nil {
		return nil, err
	}

	if blockSize <= 0 {
		return nil, ErrUnparsableBlockString
	}

	b.blockStartingBytePosition = blockStartingBytePosition
	b.blockEndingBytePosition = blockStartingBytePosition + int(blockSize) - 1

	return b, nil
}

func (b *UniformFuzzyHashBlock) isEqual(o *UniformFuzzyHashBlock) bool {
	if o == nil {
		return false
	}
	if b.blockHash != o.blockHash {
		return false
	}
	if b.blockSize() != o.blockSize() {
		return false
	}
	return true
}
