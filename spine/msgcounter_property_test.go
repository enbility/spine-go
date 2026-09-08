package spine

import (
	"math"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"testing/quick"

	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

// Property 1: msgCounter values are always ascending (except at overflow)
// TC_SPINE_DATA_001: outgoing msgCounter values are assigned in strictly ascending order.
func TestProperty_MsgCounter_AlwaysAscending(t *testing.T) {
	config := &quick.Config{
		MaxCount: 1000,
	}
	
	property := func(numMessages uint16) bool {
		if numMessages == 0 {
			return true
		}
		
		temp := &WriteMessageHandler{}
		sut := NewSender(temp)
		senderImpl := sut.(*Sender)
		
		var prevCounter model.MsgCounterType
		for i := uint16(0); i < numMessages; i++ {
			counter := senderImpl.getMsgCounter()
			
			if i > 0 {
				// Check ascending (allow for overflow)
				if *counter < prevCounter && prevCounter != math.MaxUint64 {
					return false
				}
			}
			prevCounter = *counter
		}
		return true
	}
	
	if err := quick.Check(property, config); err != nil {
		t.Error(err)
	}
}

// Property 2: msgCounter values are unique within a window
func TestProperty_MsgCounter_UniqueInWindow(t *testing.T) {
	config := &quick.Config{
		MaxCount: 100,
	}
	
	property := func(windowSize uint16) bool {
		if windowSize == 0 || windowSize > 10000 {
			windowSize = 1000 // Reasonable window size
		}
		
		temp := &WriteMessageHandler{}
		sut := NewSender(temp)
		senderImpl := sut.(*Sender)
		
		seen := make(map[model.MsgCounterType]bool)
		
		for i := uint16(0); i < windowSize; i++ {
			counter := senderImpl.getMsgCounter()
			if seen[*counter] {
				return false // Duplicate found
			}
			seen[*counter] = true
		}
		return true
	}
	
	if err := quick.Check(property, config); err != nil {
		t.Error(err)
	}
}

// Property 3: Thread safety - concurrent access produces unique counters
func TestProperty_MsgCounter_ThreadSafe(t *testing.T) {
	config := &quick.Config{
		MaxCount: 50,
	}
	
	property := func(numGoroutines, msgsPerGoroutine uint8) bool {
		if numGoroutines == 0 {
			numGoroutines = 10
		}
		if msgsPerGoroutine == 0 {
			msgsPerGoroutine = 10
		}
		
		temp := &WriteMessageHandler{}
		sut := NewSender(temp)
		senderImpl := sut.(*Sender)
		
		totalMessages := int(numGoroutines) * int(msgsPerGoroutine)
		countersChan := make(chan model.MsgCounterType, totalMessages)
		
		var wg sync.WaitGroup
		wg.Add(int(numGoroutines))
		
		for g := uint8(0); g < numGoroutines; g++ {
			go func() {
				defer wg.Done()
				for m := uint8(0); m < msgsPerGoroutine; m++ {
					counter := senderImpl.getMsgCounter()
					countersChan <- *counter
				}
			}()
		}
		
		wg.Wait()
		close(countersChan)
		
		// Check uniqueness
		seen := make(map[model.MsgCounterType]bool)
		count := 0
		for counter := range countersChan {
			if seen[counter] {
				return false // Duplicate found
			}
			seen[counter] = true
			count++
		}
		
		return count == totalMessages
	}
	
	if err := quick.Check(property, config); err != nil {
		t.Error(err)
	}
}

// Property 4: msgCounter never skips backwards (except overflow)
func TestProperty_MsgCounter_NoBackwardSkips(t *testing.T) {
	config := &quick.Config{
		MaxCount: 500,
		Values: func(values []reflect.Value, rand *rand.Rand) {
			// Generate test cases with different starting points
			startingPoint := rand.Uint64()
			numMessages := rand.Intn(100) + 1
			values[0] = reflect.ValueOf(startingPoint)
			values[1] = reflect.ValueOf(numMessages)
		},
	}
	
	property := func(startingPoint uint64, numMessages int) bool {
		temp := &WriteMessageHandler{}
		sut := NewSender(temp)
		senderImpl := sut.(*Sender)
		
		// Set starting point
		senderImpl.msgNum = startingPoint
		
		var prevCounter model.MsgCounterType
		for i := 0; i < numMessages; i++ {
			counter := senderImpl.getMsgCounter()
			
			if i > 0 {
				// Check no backward skips (except at overflow boundary)
				if *counter < prevCounter {
					// This is only valid if we wrapped around from max to 0
					if prevCounter != math.MaxUint64 || *counter != 0 {
						return false
					}
				}
			}
			prevCounter = *counter
		}
		return true
	}
	
	if err := quick.Check(property, config); err != nil {
		t.Error(err)
	}
}

// Property 5: Overflow behavior - max+1 becomes 0
func TestProperty_MsgCounter_OverflowBehavior(t *testing.T) {
	// Direct test since we need specific values near overflow
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)
	senderImpl := sut.(*Sender)
	
	testCases := []uint64{
		math.MaxUint64 - 10,
		math.MaxUint64 - 5,
		math.MaxUint64 - 2,
		math.MaxUint64 - 1,
	}
	
	for _, startValue := range testCases {
		// Reset sender with new starting value
		senderImpl.msgNum = startValue
		
		// Generate counters until we cross the overflow boundary
		var counters []model.MsgCounterType
		for i := 0; i < 15; i++ {
			counter := senderImpl.getMsgCounter()
			counters = append(counters, *counter)
		}
		
		// Find the overflow point
		overflowFound := false
		for i := 1; i < len(counters); i++ {
			if counters[i] < counters[i-1] {
				// Overflow detected
				assert.Equal(t, model.MsgCounterType(math.MaxUint64), counters[i-1], 
					"Counter before overflow should be max value")
				assert.Equal(t, model.MsgCounterType(0), counters[i], 
					"Counter after overflow should be 0")
				overflowFound = true
				break
			}
		}
		
		if startValue >= math.MaxUint64-14 {
			assert.True(t, overflowFound, "Overflow should have been detected for start value %d", startValue)
		}
	}
}

// Property 6: Starting value is always 1 for new sender
func TestProperty_MsgCounter_InitialValue(t *testing.T) {
	config := &quick.Config{
		MaxCount: 100,
	}
	
	property := func(iterations uint8) bool {
		if iterations == 0 {
			iterations = 1
		}
		
		for i := uint8(0); i < iterations; i++ {
			temp := &WriteMessageHandler{}
			sut := NewSender(temp)
			senderImpl := sut.(*Sender)
			
			counter := senderImpl.getMsgCounter()
			if *counter != 1 {
				return false
			}
		}
		return true
	}
	
	if err := quick.Check(property, config); err != nil {
		t.Error(err)
	}
}

// Property 7: Gap sizes are reasonable (implementation allows skipping)
func TestProperty_MsgCounter_ReasonableGaps(t *testing.T) {
	temp := &WriteMessageHandler{}
	sut := NewSender(temp)
	senderImpl := sut.(*Sender)
	
	const numMessages = 1000
	var prevCounter model.MsgCounterType
	
	for i := 0; i < numMessages; i++ {
		counter := senderImpl.getMsgCounter()
		
		if i > 0 {
			gap := *counter - prevCounter
			// With atomic increment, gap should always be 1
			assert.Equal(t, model.MsgCounterType(1), gap, 
				"Gap between consecutive counters should be 1")
		}
		prevCounter = *counter
	}
}