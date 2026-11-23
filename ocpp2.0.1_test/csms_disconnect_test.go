package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocpp2 "github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Responses are matched to callbacks by position only. A request that never got its response
// before the charging station disconnected must therefore be canceled on disconnect. Otherwise
// its callback is paired with the first response after a reconnect, and the callback's type
// assertion panics (a SetChargingProfileResponse arriving at a TriggerMessage callback).
func (suite *OcppV2TestSuite) TestCSMSDisconnectCancelsPendingCallbacks() {
	t := suite.T()
	wsId := "test_id"
	channel := NewMockWebSocket(wsId)

	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId})
	suite.csms.SetChargingStationDisconnectedHandler(func(ocpp2.ChargingStationConnection) {})
	suite.csms.Start(8887, "somePath")
	suite.mockWsServer.NewClientHandler(channel)

	// 1. Send a request and lose the station before it answers.
	triggerResult := make(chan error, 1)
	err := suite.csms.TriggerMessage(wsId, func(response *remotecontrol.TriggerMessageResponse, err error) {
		assert.Nil(t, response)
		triggerResult <- err
	}, remotecontrol.MessageTriggerStatusNotification)
	require.NoError(t, err)

	suite.mockWsServer.DisconnectedClientHandler(channel)
	select {
	case err := <-triggerResult:
		require.Error(t, err)
	case <-time.After(time.Second):
		t.Fatal("callback of the unanswered request was not canceled on disconnect")
	}

	// 2. After the station reconnects, a response must reach the callback of its own request.
	suite.mockWsServer.NewClientHandler(channel)
	schedule := types.NewChargingSchedule(1, types.ChargingRateUnitWatts, types.NewChargingSchedulePeriod(0, 200.0))
	profile := types.NewChargingProfile(1, 7, types.ChargingProfilePurposeChargingStationMaxProfile, types.ChargingProfileKindAbsolute, []types.ChargingSchedule{*schedule})
	profileResult := make(chan *smartcharging.SetChargingProfileResponse, 1)
	err = suite.csms.SetChargingProfile(wsId, func(response *smartcharging.SetChargingProfileResponse, err error) {
		assert.NoError(t, err)
		profileResult <- response
	}, 1, profile)
	require.NoError(t, err)

	// Mark the request as pending, otherwise the response is ignored.
	suite.ocppjServer.RequestState.AddPendingRequest(wsId, defaultMessageId, smartcharging.NewSetChargingProfileRequest(1, profile))
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, defaultMessageId, smartcharging.ChargingProfileStatusAccepted)
	require.NoError(t, suite.mockWsServer.MessageHandler(channel, []byte(responseJson)))

	select {
	case response := <-profileResult:
		require.NotNil(t, response)
		assert.Equal(t, smartcharging.ChargingProfileStatusAccepted, response.Status)
	case <-time.After(time.Second):
		t.Fatal("response did not reach the callback of its request")
	}
	// The canceled callback must not be called a second time.
	select {
	case err := <-triggerResult:
		t.Fatalf("canceled callback was called again: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
}
