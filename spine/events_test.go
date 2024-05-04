package spine

import (
	"sync"
	"testing"
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestEventsSuite(t *testing.T) {
	suite.Run(t, new(EventsTestSuite))
}

type TestDummy struct {
}

func (s *TestDummy) HandleEvent(event api.EventPayload) {}

type EventsTestSuite struct {
	suite.Suite

	mux sync.Mutex

	handlerInvoked bool
}

func (s *EventsTestSuite) BeforeTest(suiteName, testName string) {
	s.setHandlerInvoked(false)
}

func (s *EventsTestSuite) setHandlerInvoked(value bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.handlerInvoked = value
}

func (s *EventsTestSuite) isHandlerInvoked() bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.handlerInvoked
}

func (s *EventsTestSuite) HandleEvent(event api.EventPayload) {
	s.setHandlerInvoked(true)
}

func (s *EventsTestSuite) Test_Un_Subscribe() {
	err := Events.Subscribe(s)
	assert.Nil(s.T(), err)

	err = Events.Subscribe(s)
	assert.Nil(s.T(), err)

	testDummy := &TestDummy{}
	err = Events.Subscribe(testDummy)
	assert.Nil(s.T(), err)

	err = Events.Unsubscribe(s)
	assert.Nil(s.T(), err)

	err = Events.Unsubscribe(s)
	assert.Nil(s.T(), err)

	err = Events.Unsubscribe(testDummy)
	assert.Nil(s.T(), err)
}

func (s *EventsTestSuite) Test_Publish_Core() {
	err := Events.subscribe(api.EventHandlerLevelCore, s)
	assert.Nil(s.T(), err)

	Events.Publish(api.EventPayload{})

	assert.True(s.T(), s.isHandlerInvoked())

	err = Events.Unsubscribe(s)
	assert.Nil(s.T(), err)
}

func (s *EventsTestSuite) Test_Publish_Application() {
	err := Events.Subscribe(s)
	assert.Nil(s.T(), err)

	Events.Publish(api.EventPayload{})

	time.Sleep(time.Millisecond * 200)
	assert.True(s.T(), s.isHandlerInvoked())

	err = Events.Unsubscribe(s)
	assert.Nil(s.T(), err)
}
