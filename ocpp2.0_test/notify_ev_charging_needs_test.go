package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Tests
func (suite *OcppV2TestSuite) TestACChargingParametersValidation() {
	var testTable = []GenericTestEntry{
		{smartcharging.ACChargingParameters{EnergyAmount: 42, EvMinCurrent: 6, EvMaxCurrent: 42, EvMaxVoltage: 240}, true},
		{smartcharging.ACChargingParameters{}, true},
		{smartcharging.ACChargingParameters{EnergyAmount: 42, EvMinCurrent: 6, EvMaxCurrent: 42, EvMaxVoltage: -1}, false},
		{smartcharging.ACChargingParameters{EnergyAmount: 42, EvMinCurrent: 6, EvMaxCurrent: -1, EvMaxVoltage: 240}, false},
		{smartcharging.ACChargingParameters{EnergyAmount: 42, EvMinCurrent: -1, EvMaxCurrent: 42, EvMaxVoltage: 240}, false},
		{smartcharging.ACChargingParameters{EnergyAmount: -1, EvMinCurrent: 6, EvMaxCurrent: 42, EvMaxVoltage: 240}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}


func (suite *OcppV2TestSuite) TestDCChargingParametersValidation() {
	var testTable = []GenericTestEntry{
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42, EvMaxVoltage: 240, EnergyAmount: newInt(10), EvMaxPower: newInt(42), StateOfCharge: newInt(46), EvEnergyCapacity: newInt(98), FullSoC: newInt(100), BulkSoC: newInt(42)}, true},
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42, EvMaxVoltage: 240}, true},
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42}, true},
		{smartcharging.DCChargingParameters{}, true},
		{smartcharging.DCChargingParameters{EvMaxCurrent: -1, EvMaxVoltage: 240}, false},
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42, EvMaxVoltage: -1}, false},
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42, EvMaxVoltage: 240, EnergyAmount: newInt(-1)}, false},
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42, EvMaxVoltage: 240, EvMaxPower: newInt(-1)}, false},
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42, EvMaxVoltage: 240, StateOfCharge: newInt(-1)}, false},
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42, EvMaxVoltage: 240, StateOfCharge: newInt(101)}, false},
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42, EvMaxVoltage: 240, EvEnergyCapacity: newInt(-1)}, false},
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42, EvMaxVoltage: 240, FullSoC: newInt(-1)}, false},
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42, EvMaxVoltage: 240, FullSoC: newInt(101)}, false},
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42, EvMaxVoltage: 240, BulkSoC: newInt(-1)}, false},
		{smartcharging.DCChargingParameters{EvMaxCurrent: 42, EvMaxVoltage: 240, FullSoC: newInt(101)}, false},
	}
	ExecuteGenericTestTable(suite.T(), testTable)
}

func (suite *OcppV2TestSuite) TestNotifyEVChargingNeedsRequestValidation() {
	acChargingParameters := smartcharging.ACChargingParameters{
		EnergyAmount: 10,
		EvMinCurrent: 0,
		EvMaxCurrent: 42,
		EvMaxVoltage: 240,
	}
	dcChargingParameters := smartcharging.DCChargingParameters{
		EvMaxCurrent:     42,
		EvMaxVoltage:     240,
		EnergyAmount:     newInt(10),
		EvMaxPower:       newInt(42),
		StateOfCharge:    newInt(46),
		EvEnergyCapacity: newInt(98),
		FullSoC:          newInt(100),
		BulkSoC:          newInt(42),
	}
	var requestTable = []GenericTestEntry{
		{smartcharging.NotifyEVChargingNeedsRequest{MaxScheduleTuples: newInt(2), EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC, DepartureTime: types.NewDateTime(time.Now()), ACChargingParameters: &acChargingParameters, DCChargingParameters: &dcChargingParameters}}, true},
		{smartcharging.NotifyEVChargingNeedsRequest{MaxScheduleTuples: newInt(2), EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC, DepartureTime: types.NewDateTime(time.Now()), ACChargingParameters: &acChargingParameters}}, true},
		{smartcharging.NotifyEVChargingNeedsRequest{MaxScheduleTuples: newInt(2), EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC, DepartureTime: types.NewDateTime(time.Now())}}, true},
		{smartcharging.NotifyEVChargingNeedsRequest{MaxScheduleTuples: newInt(2), EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC}}, true},
		{smartcharging.NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC}}, true},
		{smartcharging.NotifyEVChargingNeedsRequest{EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{}}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{EvseID: 1}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC}}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{MaxScheduleTuples: newInt(2), EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: "invalidEnergyTransferMode", DepartureTime: types.NewDateTime(time.Now()), ACChargingParameters: &acChargingParameters, DCChargingParameters: &dcChargingParameters}}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{MaxScheduleTuples: newInt(2), EvseID: -1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC, DepartureTime: types.NewDateTime(time.Now()), ACChargingParameters: &acChargingParameters, DCChargingParameters: &dcChargingParameters}}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{MaxScheduleTuples: newInt(-1), EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC, DepartureTime: types.NewDateTime(time.Now()), ACChargingParameters: &acChargingParameters, DCChargingParameters: &dcChargingParameters}}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{MaxScheduleTuples: newInt(2), EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC, DepartureTime: types.NewDateTime(time.Now()), ACChargingParameters: &acChargingParameters, DCChargingParameters: &smartcharging.DCChargingParameters{EvMaxCurrent: -1, EvMaxVoltage: 240}}}, false},
		{smartcharging.NotifyEVChargingNeedsRequest{MaxScheduleTuples: newInt(2), EvseID: 1, ChargingNeeds: smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC, DepartureTime: types.NewDateTime(time.Now()), ACChargingParameters: &smartcharging.ACChargingParameters{EvMaxCurrent: -1}, DCChargingParameters: &dcChargingParameters}}, false},
	}
	ExecuteGenericTestTable(suite.T(), requestTable)
}

func (suite *OcppV2TestSuite) TestNotifyEVChargingNeedsResponseValidation() {
	t := suite.T()
	var responseTable = []GenericTestEntry{
		{smartcharging.NotifyEVChargingNeedsResponse{Status: smartcharging.NotifyEVChargingNeedsStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "Accepted", AdditionalInfo: "dummyInfo"}}, true},
		{smartcharging.NotifyEVChargingNeedsResponse{Status: smartcharging.NotifyEVChargingNeedsStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "Accepted"}}, true},
		{smartcharging.NotifyEVChargingNeedsResponse{Status: smartcharging.NotifyEVChargingNeedsStatusAccepted}, true},
		{smartcharging.NotifyEVChargingNeedsResponse{}, false},
		{smartcharging.NotifyEVChargingNeedsResponse{Status: "invalidStatus"}, false},
		{smartcharging.NotifyEVChargingNeedsResponse{Status: smartcharging.NotifyEVChargingNeedsStatusAccepted, StatusInfo: &types.StatusInfo{}}, false},
	}
	ExecuteGenericTestTable(t, responseTable)
}

func (suite *OcppV2TestSuite) TestNotifyEVChargingNeedsE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	maxScheduleTuples := newInt(2)
	evseID := 1
	dcChargingParameters := smartcharging.DCChargingParameters{
		EvMaxCurrent:     42,
		EvMaxVoltage:     240,
		StateOfCharge:    newInt(46),
	}
	chargingNeeds := smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC, DepartureTime: types.NewDateTime(time.Now()), DCChargingParameters: &dcChargingParameters}
	status := smartcharging.NotifyEVChargingNeedsStatusAccepted
	statusInfo := types.StatusInfo{ReasonCode: "Accepted", AdditionalInfo: "dummyInfo"}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"maxScheduleTuples":%v,"evseId":%v,"chargingNeeds":{"requestedEnergyTransfer":"%v","departureTime":"%v","dcChargingParameters":{"evMaxCurrent":%v,"evMaxVoltage":%v,"stateOfCharge":%v}}}]`,
		messageId, smartcharging.NotifyEVChargingNeedsFeatureName, *maxScheduleTuples, evseID, chargingNeeds.RequestedEnergyTransfer, chargingNeeds.DepartureTime.FormatTimestamp(), dcChargingParameters.EvMaxCurrent, dcChargingParameters.EvMaxVoltage, *dcChargingParameters.StateOfCharge)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v","additionalInfo":"%v"}}]`,
		messageId, status, statusInfo.ReasonCode, statusInfo.AdditionalInfo)
	response := smartcharging.NewNotifyEVChargingNeedsResponse(status)
	response.StatusInfo = &statusInfo
	channel := NewMockWebSocket(wsId)

	handler := MockCSMSSmartChargingHandler{}
	handler.On("OnNotifyEVChargingNeeds", mock.AnythingOfType("string"), mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*smartcharging.NotifyEVChargingNeedsRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		require.NotNil(t, request.MaxScheduleTuples)
		assert.Equal(t, *maxScheduleTuples, *request.MaxScheduleTuples)
		assert.Equal(t, evseID, request.EvseID)
		assert.Equal(t, chargingNeeds.RequestedEnergyTransfer, request.ChargingNeeds.RequestedEnergyTransfer)
		require.NotNil(t, request.ChargingNeeds.DepartureTime)
		assertDateTimeEquality(t, chargingNeeds.DepartureTime, request.ChargingNeeds.DepartureTime)
		require.NotNil(t, request.ChargingNeeds.DCChargingParameters)
		assert.Equal(t, chargingNeeds.DCChargingParameters.EvMaxCurrent, request.ChargingNeeds.DCChargingParameters.EvMaxCurrent)
		assert.Equal(t, chargingNeeds.DCChargingParameters.EvMaxVoltage, request.ChargingNeeds.DCChargingParameters.EvMaxVoltage)
		require.NotNil(t, request.ChargingNeeds.DCChargingParameters.StateOfCharge)
		assert.Equal(t, *chargingNeeds.DCChargingParameters.StateOfCharge, *request.ChargingNeeds.DCChargingParameters.StateOfCharge)
		require.Nil(t, request.ChargingNeeds.DCChargingParameters.EnergyAmount)
		require.Nil(t, request.ChargingNeeds.DCChargingParameters.BulkSoC)
		require.Nil(t, request.ChargingNeeds.DCChargingParameters.FullSoC)
		require.Nil(t, request.ChargingNeeds.DCChargingParameters.EvMaxPower)
		require.Nil(t, request.ChargingNeeds.DCChargingParameters.EvEnergyCapacity)
		require.Nil(t, request.ChargingNeeds.ACChargingParameters)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	r, err := suite.chargingStation.NotifyEVChargingNeeds(evseID, chargingNeeds, func(request *smartcharging.NotifyEVChargingNeedsRequest) {
		request.MaxScheduleTuples = maxScheduleTuples
	})
	require.Nil(t, err)
	require.NotNil(t, r)
}

func (suite *OcppV2TestSuite) TestNotifyEVChargingNeedsInvalidEndpoint() {
	messageId := defaultMessageId
	maxScheduleTuples := newInt(2)
	evseID := 1
	dcChargingParameters := smartcharging.DCChargingParameters{
		EvMaxCurrent:     42,
		EvMaxVoltage:     240,
		StateOfCharge:    newInt(46),
	}
	chargingNeeds := smartcharging.ChargingNeeds{RequestedEnergyTransfer: smartcharging.EnergyTransferModeDC, DepartureTime: types.NewDateTime(time.Now()), DCChargingParameters: &dcChargingParameters}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"maxScheduleTuples":%v,"evseId":%v,"chargingNeeds":{"requestedEnergyTransfer":"%v","departureTime":"%v","dcChargingParameters":{"evMaxCurrent":%v,"evMaxVoltage":%v,"stateOfCharge":%v}}}]`,
		messageId, smartcharging.NotifyEVChargingNeedsFeatureName, *maxScheduleTuples, evseID, chargingNeeds.RequestedEnergyTransfer, chargingNeeds.DepartureTime.FormatTimestamp(), dcChargingParameters.EvMaxCurrent, dcChargingParameters.EvMaxVoltage, *dcChargingParameters.StateOfCharge)
	request := smartcharging.NewNotifyEVChargingNeedsRequest(evseID, chargingNeeds)
	request.MaxScheduleTuples = maxScheduleTuples
	testUnsupportedRequestFromCentralSystem(suite, request, requestJson, messageId)
}
