package bloomfilter

import (
	"encoding/binary"
	"errors"
	"hash/fnv"
	"math"
	"sync"
)

var ErrIncompatibleFilters = errors.New("bloom filters have different sizes or number of hash functions")

// BloomFilter represents a probabilistic data structure for membership testing
type BloomFilter struct {
	bitSet    []uint64
	numHashes int
	mutex     sync.RWMutex
}

// New creates a new BloomFilter with the given expected number of elements and false positive rate
func New(expectedElements int, falsePositiveRate float64) *BloomFilter {
	if expectedElements <= 0 {
		expectedElements = 1
	}
	if falsePositiveRate <= 0 || falsePositiveRate >= 1 {
		falsePositiveRate = 0.01 // default 1% false positive rate
	}

	// Calculate optimal bit array size: m = -n * ln(p) / (ln(2))^2
	bits := uint64(-float64(expectedElements) * math.Log(falsePositiveRate) / (math.Ln2 * math.Ln2))

	// Calculate optimal number of hash functions: k = (m/n) * ln(2)
	numHashes := min(max(int(float64(bits)*math.Ln2/float64(expectedElements)), 1), 32)

	// Ensure bits is at least 1
	if bits == 0 {
		bits = 1
	}

	// Create bit array (using uint64 words)
	numWords := (bits + 63) / 64
	bitSet := make([]uint64, numWords)

	return &BloomFilter{
		bitSet:    bitSet,
		numHashes: numHashes,
	}
}

// WithSize creates a new BloomFilter with explicit bit array size and number of hash functions
func WithSize(bitArraySize uint64, numHashes int) *BloomFilter {
	if bitArraySize == 0 {
		bitArraySize = 1
	}
	if numHashes <= 0 {
		numHashes = 1
	}
	if numHashes > 32 {
		numHashes = 32
	}

	numWords := (bitArraySize + 63) / 64
	bitSet := make([]uint64, numWords)

	return &BloomFilter{
		bitSet:    bitSet,
		numHashes: numHashes,
	}
}

// Put adds an element to the BloomFilter
func (bf *BloomFilter) Put(data []byte) {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	hashes := bf.getHashes(data)
	for _, hash := range hashes {
		wordIndex := hash / 64
		bitIndex := hash % 64

		if wordIndex < uint64(len(bf.bitSet)) {
			bf.bitSet[wordIndex] |= 1 << bitIndex
		}
	}
}

// MightContain checks if the element might be in the set
// Returns true if the element might be in the set (could be false positive)
// Returns false if the element is definitely not in the set
func (bf *BloomFilter) MightContain(data []byte) bool {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	hashes := bf.getHashes(data)
	for _, hash := range hashes {
		wordIndex := hash / 64
		bitIndex := hash % 64

		if wordIndex >= uint64(len(bf.bitSet)) {
			return false
		}

		if bf.bitSet[wordIndex]&(1<<bitIndex) == 0 {
			return false
		}
	}
	return true
}

// getHashes generates multiple hash values for the given data
func (bf *BloomFilter) getHashes(data []byte) []uint64 {
	hashes := make([]uint64, bf.numHashes)

	// Generate multiple hashes using different salts with FNV
	for i := 0; i < bf.numHashes; i++ {
		h := fnv.New64a()
		h.Write(data)
		// Add salt by writing the hash function index as 8 bytes
		var salt [8]byte
		binary.LittleEndian.PutUint64(salt[:], uint64(i))
		h.Write(salt[:])

		hash := h.Sum64()
		hashes[i] = hash % (uint64(len(bf.bitSet)) * 64)
	}

	return hashes
}

// Size returns the approximate number of elements in the filter
// Note: This is an estimate and may not be accurate due to hash collisions
func (bf *BloomFilter) Size() int {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	// Count set bits
	setBits := 0
	for _, word := range bf.bitSet {
		setBits += popCount(word)
	}

	// Estimate number of elements using the formula:
	// n = -m * ln(1 - X/m) / k
	// where X is the number of set bits, m is total bits, k is number of hash functions
	totalBits := uint64(len(bf.bitSet)) * 64
	if setBits == 0 {
		return 0
	}
	if setBits >= int(totalBits) {
		// Filter is saturated, return a large number
		return int(totalBits)
	}

	X := float64(setBits)
	m := float64(totalBits)
	k := float64(bf.numHashes)

	// Avoid division by zero and log of negative numbers
	if X >= m {
		return int(m)
	}

	n := -m * math.Log(1.0-X/m) / k
	if n < 0 {
		return 0
	}
	return int(n)
}

// IsEmpty returns true if no elements have been added (all bits are 0)
func (bf *BloomFilter) IsEmpty() bool {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	for _, word := range bf.bitSet {
		if word != 0 {
			return false
		}
	}
	return true
}

// Clear removes all elements from the filter
func (bf *BloomFilter) Clear() {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	for i := range bf.bitSet {
		bf.bitSet[i] = 0
	}
}

// Clone creates a deep copy of the BloomFilter
func (bf *BloomFilter) Clone() *BloomFilter {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	newBitSet := make([]uint64, len(bf.bitSet))
	copy(newBitSet, bf.bitSet)

	return &BloomFilter{
		bitSet:    newBitSet,
		numHashes: bf.numHashes,
	}
}

// Merge combines this BloomFilter with another one
// Both filters must have the same size and number of hash functions
func (bf *BloomFilter) Merge(other *BloomFilter) error {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	other.mutex.RLock()
	defer other.mutex.RUnlock()

	if len(bf.bitSet) != len(other.bitSet) || bf.numHashes != other.numHashes {
		return ErrIncompatibleFilters
	}

	for i, word := range other.bitSet {
		bf.bitSet[i] |= word
	}

	return nil
}

// popCount counts the number of set bits in a word (Hamming weight)
func popCount(x uint64) int {
	count := 0
	for x != 0 {
		count += int(x & 1)
		x >>= 1
	}
	return count
}

// FalsePositiveRate returns the current estimated false positive rate
func (bf *BloomFilter) FalsePositiveRate() float64 {
	bf.mutex.RLock()
	defer bf.mutex.RUnlock()

	totalBits := uint64(len(bf.bitSet)) * 64
	setBits := 0
	for _, word := range bf.bitSet {
		setBits += popCount(word)
	}

	if setBits == 0 {
		return 0.0
	}

	X := float64(setBits)
	m := float64(totalBits)
	k := float64(bf.numHashes)

	// False positive rate: (1 - e^(-k*n/m))^k
	// But we can estimate it as: (X/m)^k
	rate := math.Pow(X/m, k)
	if rate > 1.0 {
		rate = 1.0
	}
	return rate
}

// Capacity returns the total number of bits in the filter
func (bf *BloomFilter) Capacity() uint64 {
	return uint64(len(bf.bitSet)) * 64
}

// NumHashes returns the number of hash functions used
func (bf *BloomFilter) NumHashes() int {
	return bf.numHashes
}
