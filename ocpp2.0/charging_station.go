package ocpp2

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/internal/callbackqueue"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/authorization"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/data"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/display"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/meter"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/security"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/tariffcost"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/transactions"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
)

type chargingStation struct {
	client               *ocppj.Client
	securityHandler      security.ChargingStationHandler
	provisioningHandler  provisioning.ChargingStationHandler
	authorizationHandler authorization.ChargingStationHandler
	localAuthListHandler localauth.ChargingStationHandler
	transactionsHandler  transactions.ChargingStationHandler
	remoteControlHandler remotecontrol.ChargingStationHandler
	availabilityHandler  availability.ChargingStationHandler
	reservationHandler   reservation.ChargingStationHandler
	tariffCostHandler    tariffcost.ChargingStationHandler
	meterHandler         meter.ChargingStationHandler
	smartChargingHandler smartcharging.ChargingStationHandler
	firmwareHandler      firmware.ChargingStationHandler
	iso15118Handler      iso15118.ChargingStationHandler
	diagnosticsHandler   diagnostics.ChargingStationHandler
	displayHandler       display.ChargingStationHandler
	dataHandler          data.ChargingStationHandler
	responseHandler      chan ocpp.Response
	errorHandler         chan error
	callbacks            callbackqueue.CallbackQueue
	stopC                chan struct{}
	errC                 chan error // external error channel
}

func (cs *chargingStation) error(err error) {
	if cs.errC != nil {
		cs.errC <- err
	}
}

// Errors returns a channel for error messages. If it doesn't exist it es created.
func (cs *chargingStation) Errors() <-chan error {
	if cs.errC == nil {
		cs.errC = make(chan error, 1)
	}
	return cs.errC
}

// Callback invoked whenever a queued request is canceled, due to timeout.
// By default, the callback returns a GenericError to the caller, who sent the original request.
func (cs *chargingStation) onRequestTimeout(rID string, action string, request ocpp.Request) {
	err := ocpp.NewError(ocppj.GenericError, "request timed out, no response received from server", rID)
	cs.errorHandler <- err
}

func (cs *chargingStation) BootNotification(reason provisioning.BootReason, model string, vendor string, props ...func(request *provisioning.BootNotificationRequest)) (*provisioning.BootNotificationResponse, error) {
	request := provisioning.NewBootNotificationRequest(reason, model, vendor)
	for _, fn := range props {
		fn(request)
	}
	response, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return response.(*provisioning.BootNotificationResponse), err
	}
}

func (cs *chargingStation) Authorize(idToken string, tokenType types.IdTokenType, props ...func(request *authorization.AuthorizeRequest)) (*authorization.AuthorizeResponse, error) {
	request := authorization.NewAuthorizationRequest(idToken, tokenType)
	for _, fn := range props {
		fn(request)
	}
	response, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return response.(*authorization.AuthorizeResponse), err
	}
}

func (cs *chargingStation) ClearedChargingLimit(chargingLimitSource types.ChargingLimitSourceType, props ...func(request *smartcharging.ClearedChargingLimitRequest)) (*smartcharging.ClearedChargingLimitResponse, error) {
	request := smartcharging.NewClearedChargingLimitRequest(chargingLimitSource)
	for _, fn := range props {
		fn(request)
	}
	response, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return response.(*smartcharging.ClearedChargingLimitResponse), err
	}
}

func (cs *chargingStation) DataTransfer(vendorId string, props ...func(request *data.DataTransferRequest)) (*data.DataTransferResponse, error) {
	request := data.NewDataTransferRequest(vendorId)
	for _, fn := range props {
		fn(request)
	}
	response, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return response.(*data.DataTransferResponse), err
	}
}

func (cs *chargingStation) FirmwareStatusNotification(status firmware.FirmwareStatus, requestID int, props ...func(request *firmware.FirmwareStatusNotificationRequest)) (*firmware.FirmwareStatusNotificationResponse, error) {
	request := firmware.NewFirmwareStatusNotificationRequest(status, requestID)
	for _, fn := range props {
		fn(request)
	}
	response, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return response.(*firmware.FirmwareStatusNotificationResponse), err
	}
}

func (cs *chargingStation) Get15118EVCertificate(schemaVersion string, exiRequest string, props ...func(request *iso15118.Get15118EVCertificateRequest)) (*iso15118.Get15118EVCertificateResponse, error) {
	request := iso15118.NewGet15118EVCertificateRequest(schemaVersion, exiRequest)
	for _, fn := range props {
		fn(request)
	}
	response, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return response.(*iso15118.Get15118EVCertificateResponse), err
	}
}

func (cs *chargingStation) GetCertificateStatus(ocspRequestData types.OCSPRequestDataType, props ...func(request *iso15118.GetCertificateStatusRequest)) (*iso15118.GetCertificateStatusResponse, error) {
	request := iso15118.NewGetCertificateStatusRequest(ocspRequestData)
	for _, fn := range props {
		fn(request)
	}
	response, err := cs.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return response.(*iso15118.GetCertificateStatusResponse), err
	}
}

func (cs *chargingStation) SetSecurityHandler(handler security.ChargingStationHandler) {
	cs.securityHandler = handler
}

func (cs *chargingStation) SetProvisioningHandler(handler provisioning.ChargingStationHandler) {
	cs.provisioningHandler = handler
}

func (cs *chargingStation) SetAuthorizationHandler(handler authorization.ChargingStationHandler) {
	cs.authorizationHandler = handler
}

func (cs *chargingStation) SetLocalAuthListHandler(handler localauth.ChargingStationHandler) {
	cs.localAuthListHandler = handler
}

func (cs *chargingStation) SetTransactionsHandler(handler transactions.ChargingStationHandler) {
	cs.transactionsHandler = handler
}

func (cs *chargingStation) SetRemoteControlHandler(handler remotecontrol.ChargingStationHandler) {
	cs.remoteControlHandler = handler
}

func (cs *chargingStation) SetAvailabilityHandler(handler availability.ChargingStationHandler) {
	cs.availabilityHandler = handler
}

func (cs *chargingStation) SetReservationHandler(handler reservation.ChargingStationHandler) {
	cs.reservationHandler = handler
}

func (cs *chargingStation) SetTariffCostHandler(handler tariffcost.ChargingStationHandler) {
	cs.tariffCostHandler = handler
}

func (cs *chargingStation) SetMeterHandler(handler meter.ChargingStationHandler) {
	cs.meterHandler = handler
}

func (cs *chargingStation) SetSmartChargingHandler(handler smartcharging.ChargingStationHandler) {
	cs.smartChargingHandler = handler
}

func (cs *chargingStation) SetFirmwareHandler(handler firmware.ChargingStationHandler) {
	cs.firmwareHandler = handler
}

func (cs *chargingStation) SetISO15118Handler(handler iso15118.ChargingStationHandler) {
	cs.iso15118Handler = handler
}

func (cs *chargingStation) SetDiagnosticsHandler(handler diagnostics.ChargingStationHandler) {
	cs.diagnosticsHandler = handler
}

func (cs *chargingStation) SetDisplayHandler(handler display.ChargingStationHandler) {
	cs.displayHandler = handler
}

func (cs *chargingStation) SetDataHandler(handler data.ChargingStationHandler) {
	cs.dataHandler = handler
}

func (cs *chargingStation) SendRequest(request ocpp.Request) (ocpp.Response, error) {
	featureName := request.GetFeatureName()
	if _, found := cs.client.GetProfileForFeature(featureName); !found {
		return nil, fmt.Errorf("feature %v is unsupported on charging station (missing profile), cannot send request", featureName)
	}

	// Wraps an asynchronous response
	type asyncResponse struct {
		r ocpp.Response
		e error
	}
	// Create channel and pass it to a callback function, for retrieving asynchronous response
	asyncResponseC := make(chan asyncResponse, 1)
	send := func() error {
		return cs.client.SendRequest(request)
	}
	err := cs.callbacks.TryQueue("main", send, func(confirmation ocpp.Response, err error) {
		asyncResponseC <- asyncResponse{r: confirmation, e: err}
	})
	if err != nil {
		return nil, err
	}
	asyncResult, ok := <-asyncResponseC
	if !ok {
		return nil, fmt.Errorf("internal error while receiving result for %v request", request.GetFeatureName())
	}
	return asyncResult.r, asyncResult.e
}

func (cs *chargingStation) SendRequestAsync(request ocpp.Request, callback func(response ocpp.Response, err error)) error {
	featureName := request.GetFeatureName()
	if _, found := cs.client.GetProfileForFeature(featureName); !found {
		return fmt.Errorf("feature %v is unsupported on charging station (missing profile), cannot send request", featureName)
	}
	switch featureName {
	case authorization.AuthorizeFeatureName, provisioning.BootNotificationFeatureName, smartcharging.ClearedChargingLimitFeatureName, data.DataTransferFeatureName, firmware.FirmwareStatusNotificationFeatureName, iso15118.Get15118EVCertificateFeatureName, iso15118.GetCertificateStatusFeatureName:
		break
	default:
		return fmt.Errorf("unsupported action %v on charging station, cannot send request", featureName)
	}
	// Response will be retrieved asynchronously via asyncHandler
	send := func() error {
		return cs.client.SendRequest(request)
	}
	err := cs.callbacks.TryQueue("main", send, callback)
	return err
}

func (cs *chargingStation) asyncCallbackHandler() {
	for {
		select {
		case confirmation := <-cs.responseHandler:
			// Get and invoke callback
			if callback, ok := cs.callbacks.Dequeue("main"); ok {
				callback(confirmation, nil)
			} else {
				cs.error(fmt.Errorf("no callback available for incoming response %v", confirmation.GetFeatureName()))
			}
		case protoError := <-cs.errorHandler:
			// Get and invoke callback
			if callback, ok := cs.callbacks.Dequeue("main"); ok {
				callback(nil, protoError)
			} else {
				cs.error(fmt.Errorf("no callback available for incoming error %w", protoError))
			}
		case _, _ = <-cs.stopC:
			return
		}
	}
}

func (cs *chargingStation) sendResponse(response ocpp.Response, err error, requestId string) {
	// send error response
	if err != nil {
		err = cs.client.SendError(requestId, ocppj.ProtocolError, err.Error(), nil)
		if err != nil {
			cs.error(fmt.Errorf("replying cs to request %s with 'protocol error': %w", requestId, err))
		}
		return
	}

	if response == nil {
		err = fmt.Errorf("empty response to request %s", requestId)
		cs.error(err)
		return
	}

	// send response
	err = cs.client.SendResponse(requestId, response)
	if err != nil {
		cs.error(fmt.Errorf("replying to request %s with 'protocol error': %w", requestId, err))
	}
}

func (cs *chargingStation) Start(csmsUrl string) error {
	cs.stopC = make(chan struct{}, 1)
	// Async response handler receives incoming responses/errors and triggers callbacks
	err := cs.client.Start(csmsUrl)
	if err == nil {
		go cs.asyncCallbackHandler()
	}
	return err
}

func (cs *chargingStation) Stop() {
	cs.client.Stop()
}

func (cs *chargingStation) notImplementedError(requestId string, action string) {
	err := cs.client.SendError(requestId, ocppj.NotImplemented, fmt.Sprintf("no handler for action %v implemented", action), nil)
	if err != nil {
		cs.error(fmt.Errorf("replying csms to request %v with error: %w", requestId, err))
	}
}

func (cs *chargingStation) notSupportedError(requestId string, action string) {
	err := cs.client.SendError(requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on charging station", action), nil)
	if err != nil {
		cs.error(fmt.Errorf("replying csms to request %s with 'not supported': %w", requestId, err))
	}
}

func (cs *chargingStation) handleIncomingRequest(request ocpp.Request, requestId string, action string) {
	profile, found := cs.client.GetProfileForFeature(action)
	// Check whether action is supported and a listener for it exists
	if !found {
		cs.notImplementedError(requestId, action)
		return
	} else {
		supported := true
		switch profile.Name {
		case authorization.ProfileName:
			if cs.authorizationHandler == nil {
				supported = false
			}
		case availability.ProfileName:
			if cs.availabilityHandler == nil {
				supported = false
			}
		case data.ProfileName:
			if cs.dataHandler == nil {
				supported = false
			}
		case diagnostics.ProfileName:
			if cs.diagnosticsHandler == nil {
				supported = false
			}
		case display.ProfileName:
			if cs.displayHandler == nil {
				supported = false
			}
		case firmware.ProfileName:
			if cs.firmwareHandler == nil {
				supported = false
			}
		case iso15118.ProfileName:
			if cs.iso15118Handler == nil {
				supported = false
			}
		case localauth.ProfileName:
			if cs.localAuthListHandler == nil {
				supported = false
			}
		case meter.ProfileName:
			if cs.meterHandler == nil {
				supported = false
			}
		case provisioning.ProfileName:
			if cs.provisioningHandler == nil {
				supported = false
			}
		case remotecontrol.ProfileName:
			if cs.remoteControlHandler == nil {
				supported = false
			}
		case reservation.ProfileName:
			if cs.reservationHandler == nil {
				supported = false
			}
		case security.ProfileName:
			if cs.securityHandler == nil {
				supported = false
			}
		case smartcharging.ProfileName:
			if cs.smartChargingHandler == nil {
				supported = false
			}
		case tariffcost.ProfileName:
			if cs.tariffCostHandler == nil {
				supported = false
			}
		case transactions.ProfileName:
			if cs.transactionsHandler == nil {
				supported = false
			}
		}
		if !supported {
			cs.notSupportedError(requestId, action)
			return
		}
	}
	// Process request
	var response ocpp.Response = nil
	cs.client.GetProfileForFeature(action)
	var err error = nil
	switch action {
	case reservation.CancelReservationFeatureName:
		response, err = cs.reservationHandler.OnCancelReservation(request.(*reservation.CancelReservationRequest))
	case security.CertificateSignedFeatureName:
		response, err = cs.securityHandler.OnCertificateSigned(request.(*security.CertificateSignedRequest))
	case availability.ChangeAvailabilityFeatureName:
		response, err = cs.availabilityHandler.OnChangeAvailability(request.(*availability.ChangeAvailabilityRequest))
	case authorization.ClearCacheFeatureName:
		response, err = cs.authorizationHandler.OnClearCache(request.(*authorization.ClearCacheRequest))
	case smartcharging.ClearChargingProfileFeatureName:
		response, err = cs.smartChargingHandler.OnClearChargingProfile(request.(*smartcharging.ClearChargingProfileRequest))
	case display.ClearDisplayFeatureName:
		response, err = cs.displayHandler.OnClearDisplay(request.(*display.ClearDisplayRequest))
	case diagnostics.ClearVariableMonitoringFeatureName:
		response, err = cs.diagnosticsHandler.OnClearVariableMonitoring(request.(*diagnostics.ClearVariableMonitoringRequest))
	case tariffcost.CostUpdatedFeatureName:
		response, err = cs.tariffCostHandler.OnCostUpdated(request.(*tariffcost.CostUpdatedRequest))
	case diagnostics.CustomerInformationFeatureName:
		response, err = cs.diagnosticsHandler.OnCustomerInformation(request.(*diagnostics.CustomerInformationRequest))
	case data.DataTransferFeatureName:
		response, err = cs.dataHandler.OnDataTransfer(request.(*data.DataTransferRequest))
	case iso15118.DeleteCertificateFeatureName:
		response, err = cs.iso15118Handler.OnDeleteCertificate(request.(*iso15118.DeleteCertificateRequest))
	case provisioning.GetBaseReportFeatureName:
		response, err = cs.provisioningHandler.OnGetBaseReport(request.(*provisioning.GetBaseReportRequest))
	case smartcharging.GetChargingProfilesFeatureName:
		response, err = cs.smartChargingHandler.OnGetChargingProfiles(request.(*smartcharging.GetChargingProfilesRequest))
	case smartcharging.GetCompositeScheduleFeatureName:
		response, err = cs.smartChargingHandler.OnGetCompositeSchedule(request.(*smartcharging.GetCompositeScheduleRequest))
	case display.GetDisplayMessagesFeatureName:
		response, err = cs.displayHandler.OnGetDisplayMessages(request.(*display.GetDisplayMessagesRequest))
	case iso15118.GetInstalledCertificateIdsFeatureName:
		response, err = cs.iso15118Handler.OnGetInstalledCertificateIds(request.(*iso15118.GetInstalledCertificateIdsRequest))
	case localauth.GetLocalListVersionFeatureName:
		response, err = cs.localAuthListHandler.OnGetLocalListVersion(request.(*localauth.GetLocalListVersionRequest))
	case diagnostics.GetLogFeatureName:
		response, err = cs.diagnosticsHandler.OnGetLog(request.(*diagnostics.GetLogRequest))
	case diagnostics.GetMonitoringReportFeatureName:
		response, err = cs.diagnosticsHandler.OnGetMonitoringReport(request.(*diagnostics.GetMonitoringReportRequest))
	default:
		cs.notSupportedError(requestId, action)
		return
	}
	cs.sendResponse(response, err, requestId)
}
