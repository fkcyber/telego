package telego

import (
	"fmt"
	"sync"
	"testing"

	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/mymmrac/telego/telegoapi"
	mockAPI "github.com/mymmrac/telego/telegoapi/mock"
)

var (
	data      = &telegoapi.RequestData{}
	emptyResp = &telegoapi.Response{
		Ok: true,
	}

	expectedMessage = &Message{
		MessageID: 1,
	}

	testPortStart = 3100
	testPortLock  = sync.Mutex{}
)

func testAddress(t *testing.T) string {
	t.Helper()

	testPortLock.Lock()
	defer testPortLock.Unlock()

	testPortStart++
	return fmt.Sprintf("127.0.0.1:%d", testPortStart)
}

func telegoResponse(t *testing.T, v interface{}) *telegoapi.Response {
	t.Helper()

	bytesData, err := json.Marshal(v)
	require.NoError(t, err)
	return &telegoapi.Response{
		Ok:     true,
		Result: bytesData,
	}
}

type mockedBot struct {
	MockAPICaller          *mockAPI.MockCaller
	MockRequestConstructor *mockAPI.MockRequestConstructor
	Bot                    *Bot
}

func newMockedBot(ctrl *gomock.Controller) mockedBot {
	mb := mockedBot{
		MockAPICaller:          mockAPI.NewMockCaller(ctrl),
		MockRequestConstructor: mockAPI.NewMockRequestConstructor(ctrl),
	}

	//nolint:errcheck
	bot, _ := NewBot(token,
		WithAPICaller(mb.MockAPICaller),
		WithRequestConstructor(mb.MockRequestConstructor),
		WithDiscardLogger(),
		WithWarnings())

	mb.Bot = bot

	return mb
}

type testNamedReade struct{}

func (t testNamedReade) Read(_ []byte) (n int, err error) {
	panic("implement me")
}

func (t testNamedReade) Name() string {
	return "test"
}

var testInputFile = InputFile{
	File: testNamedReade{},
}
