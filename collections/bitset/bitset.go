package bitset

import (
	"sync"
)

const wordSize = 64 // bits per word (uint64)

// BitSet represents a set of non-negative integers using bit manipulation
type BitSet struct {
	bits  []uint64
	mutex sync.RWMutex
}

// New creates a new empty BitSet
func New() *BitSet {
	return &BitSet{
		bits: make([]uint64, 1),
	}
}

// WithCapacity creates a new BitSet with initial capacity for the given maximum value
func WithCapacity(maxValue int) *BitSet {
	if maxValue < 0 {
		maxValue = 0
	}
	words := (maxValue / wordSize) + 1
	return &BitSet{
		bits: make([]uint64, words),
	}
}

// Set adds the given value to the set
func (bs *BitSet) Set(value int) {
	if value < 0 {
		return // BitSet only supports non-negative integers
	}

	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	wordIndex := value / wordSize
	bitIndex := uint(value % wordSize)

	// Ensure the slice is large enough
	if wordIndex >= len(bs.bits) {
		newSize := wordIndex + 1
		newBits := make([]uint64, newSize)
		copy(newBits, bs.bits)
		bs.bits = newBits
	}

	// Set the bit
	bs.bits[wordIndex] |= 1 << bitIndex
}

// Clear removes the given value from the set
func (bs *BitSet) Clear(value int) {
	if value < 0 {
		return
	}

	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	wordIndex := value / wordSize
	if wordIndex >= len(bs.bits) {
		return // Value not present
	}

	bitIndex := uint(value % wordSize)

	// Clear the bit
	bs.bits[wordIndex] &^= 1 << bitIndex
}

// Get checks if the given value is in the set
func (bs *BitSet) Get(value int) bool {
	if value < 0 {
		return false
	}

	bs.mutex.RLock()
	defer bs.mutex.RUnlock()

	wordIndex := value / wordSize
	if wordIndex >= len(bs.bits) {
		return false
	}

	bitIndex := uint(value % wordSize)

	// Check if the bit is set
	return (bs.bits[wordIndex] & (1 << bitIndex)) != 0
}

// Size returns the number of elements in the set
func (bs *BitSet) Size() int {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()

	count := 0
	for _, word := range bs.bits {
		count += popCount(word)
	}
	return count
}

// popCount counts the number of set bits in a word (Hamming weight)
func popCount(x uint64) int {
	// Use built-in function if available, otherwise implement manually
	// Go has runtime support for this, but we'll implement a portable version
	count := 0
	for x != 0 {
		count += int(x & 1)
		x >>= 1
	}
	return count
}

// IsEmpty returns true if the set is empty
func (bs *BitSet) IsEmpty() bool {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()

	for _, word := range bs.bits {
		if word != 0 {
			return false
		}
	}
	return true
}

// ClearAll removes all elements from the set
func (bs *BitSet) ClearAll() {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	for i := range bs.bits {
		bs.bits[i] = 0
	}
}

// Union performs bitwise OR with another BitSet and stores the result in this BitSet
func (bs *BitSet) Union(other *BitSet) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	other.mutex.RLock()
	defer other.mutex.RUnlock()

	// Ensure this BitSet is at least as large as the other
	if len(bs.bits) < len(other.bits) {
		newBits := make([]uint64, len(other.bits))
		copy(newBits, bs.bits)
		bs.bits = newBits
	}

	// Perform bitwise OR
	for i, word := range other.bits {
		bs.bits[i] |= word
	}
}

// Intersection performs bitwise AND with another BitSet and stores the result in this BitSet
func (bs *BitSet) Intersection(other *BitSet) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	other.mutex.RLock()
	defer other.mutex.RUnlock()

	// Truncate or extend this BitSet to match the other's size
	if len(bs.bits) > len(other.bits) {
		// Clear extra words
		for i := len(other.bits); i < len(bs.bits); i++ {
			bs.bits[i] = 0
		}
		bs.bits = bs.bits[:len(other.bits)]
	} else if len(bs.bits) < len(other.bits) {
		// Extend with zeros
		newBits := make([]uint64, len(other.bits))
		copy(newBits, bs.bits)
		bs.bits = newBits
	}

	// Perform bitwise AND
	for i := range bs.bits {
		bs.bits[i] &= other.bits[i]
	}
}

// Difference performs bitwise AND NOT with another BitSet and stores the result in this BitSet
func (bs *BitSet) Difference(other *BitSet) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	other.mutex.RLock()
	defer other.mutex.RUnlock()

	// Ensure this BitSet is at least as large as the other
	if len(bs.bits) < len(other.bits) {
		newBits := make([]uint64, len(other.bits))
		copy(newBits, bs.bits)
		bs.bits = newBits
	}

	// Perform bitwise AND NOT (A & ~B)
	for i, word := range other.bits {
		bs.bits[i] &^= word
	}
}

// SymmetricDifference performs XOR with another BitSet and stores the result in this BitSet
func (bs *BitSet) SymmetricDifference(other *BitSet) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	other.mutex.RLock()
	defer other.mutex.RUnlock()

	// Ensure this BitSet is at least as large as the other
	maxLen := max(len(other.bits), len(bs.bits))

	if len(bs.bits) < maxLen {
		newBits := make([]uint64, maxLen)
		copy(newBits, bs.bits)
		bs.bits = newBits
	}

	// Perform bitwise XOR
	for i := range maxLen {
		var otherWord uint64
		if i < len(other.bits) {
			otherWord = other.bits[i]
		}
		bs.bits[i] ^= otherWord
	}
}

// Clone creates a deep copy of the BitSet
func (bs *BitSet) Clone() *BitSet {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()

	newBits := make([]uint64, len(bs.bits))
	copy(newBits, bs.bits)

	return &BitSet{
		bits: newBits,
	}
}

// Equals checks if two BitSets are equal
func (bs *BitSet) Equals(other *BitSet) bool {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()

	other.mutex.RLock()
	defer other.mutex.RUnlock()

	// Find the maximum length needed to compare
	maxLen := len(bs.bits)
	if len(other.bits) > maxLen {
		maxLen = len(other.bits)
	}

	for i := 0; i < maxLen; i++ {
		var thisWord, otherWord uint64
		if i < len(bs.bits) {
			thisWord = bs.bits[i]
		}
		if i < len(other.bits) {
			otherWord = other.bits[i]
		}
		if thisWord != otherWord {
			return false
		}
	}
	return true
}

// ToSlice returns a slice containing all values in the set
func (bs *BitSet) ToSlice() []int {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()

	var result []int
	for wordIndex, word := range bs.bits {
		if word == 0 {
			continue
		}
		for bitIndex := 0; bitIndex < wordSize; bitIndex++ {
			if (word & (1 << uint(bitIndex))) != 0 {
				result = append(result, wordIndex*wordSize+bitIndex)
			}
		}
	}
	return result
}

// ForEach iterates over all values in the set
func (bs *BitSet) ForEach(action func(value int)) {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()

	for wordIndex, word := range bs.bits {
		if word == 0 {
			continue
		}
		for bitIndex := range wordSize {
			if (word & (1 << uint(bitIndex))) != 0 {
				action(wordIndex*wordSize + bitIndex)
			}
		}
	}
}

// NextSetBit returns the next set bit after the given index, or -1 if none exists
func (bs *BitSet) NextSetBit(fromIndex int) int {
	if fromIndex < 0 {
		fromIndex = 0
	}

	bs.mutex.RLock()
	defer bs.mutex.RUnlock()

	startWord := fromIndex / wordSize
	startBit := fromIndex % wordSize

	// Check the starting word from the start bit
	if startWord < len(bs.bits) {
		word := bs.bits[startWord] >> uint(startBit)
		if word != 0 {
			// Find the first set bit in the shifted word
			for i := 0; i < wordSize-startBit; i++ {
				if (word & (1 << uint(i))) != 0 {
					return fromIndex + i
				}
			}
		}
	}

	// Check subsequent words
	for wordIndex := startWord + 1; wordIndex < len(bs.bits); wordIndex++ {
		word := bs.bits[wordIndex]
		if word != 0 {
			// Find the first set bit in this word
			for i := range wordSize {
				if (word & (1 << uint(i))) != 0 {
					return wordIndex*wordSize + i
				}
			}
		}
	}

	return -1
}

// Max returns the maximum value in the set, or -1 if the set is empty
func (bs *BitSet) Max() int {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()

	for i := len(bs.bits) - 1; i >= 0; i-- {
		if bs.bits[i] != 0 {
			word := bs.bits[i]
			// Find the highest set bit in this word
			for bitIndex := wordSize - 1; bitIndex >= 0; bitIndex-- {
				if (word & (1 << uint(bitIndex))) != 0 {
					return i*wordSize + bitIndex
				}
			}
		}
	}
	return -1
}

// Min returns the minimum value in the set, or -1 if the set is empty
func (bs *BitSet) Min() int {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()

	for wordIndex, word := range bs.bits {
		if word != 0 {
			// Find the lowest set bit in this word
			for bitIndex := 0; bitIndex < wordSize; bitIndex++ {
				if (word & (1 << uint(bitIndex))) != 0 {
					return wordIndex*wordSize + bitIndex
				}
			}
		}
	}
	return -1
}
