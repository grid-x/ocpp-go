package main

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocppj"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/authorization"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/data"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/meter"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/security"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/transactions"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"

	"github.com/sirupsen/logrus"
)

const (
	listenPort = 8090
	listenPath = "/ocpp/CP_X"

	heartBeatInterval = 5 * time.Second
)

type handler struct {
	csms ocpp2.CSMS
}

func main() {
	log := logrus.WithField("component", "server")

	h := &handler{}

	ocppj.SetLogger(log)
	log.Logger.SetLevel(logrus.DebugLevel)

	h.csms = ocpp2.NewCSMS(nil, nil)

	h.csms.SetNewChargingStationHandler(func(cp ocpp2.ChargingStationConnection) {
		fmt.Println("NewChargingStationHandler: ", cp)
	})
	h.csms.SetChargingStationDisconnectedHandler(func(cp ocpp2.ChargingStationConnection) {
		fmt.Println("ChargingStationDisconnectedHandler: ", cp)

	})
	h.csms.SetAuthorizationHandler(h)
	h.csms.SetAvailabilityHandler(h)
	h.csms.SetDataHandler(h)
	h.csms.SetDiagnosticsHandler(h)
	h.csms.SetDisplayHandler(h)
	h.csms.SetFirmwareHandler(h)
	h.csms.SetISO15118Handler(h)
	h.csms.SetLocalAuthListHandler(h)
	h.csms.SetMeterHandler(h)
	h.csms.SetProvisioningHandler(h)
	h.csms.SetRemoteControlHandler(h)
	h.csms.SetReservationHandler(h)
	h.csms.SetSecurityHandler(h)
	h.csms.SetSmartChargingHandler(h)
	h.csms.SetTariffCostHandler(h)
	h.csms.SetTransactionsHandler(h)

	// Start the OCPP 2 Server
	fmt.Printf("Starting OCPP Server on port: %d, path: %s \n", listenPort, listenPath)
	go h.csms.Start(listenPort, listenPath)

	select {
	case ocppErr := <-h.csms.Errors():
		fmt.Println("ERROR: ", ocppErr)
	}
}

func (h *handler) OnAuthorize(chargingStationID string, request *authorization.AuthorizeRequest) (confirmation *authorization.AuthorizeResponse, err error) {
	fmt.Println("OnAuthorize: ", chargingStationID, request)
	return nil, nil
}

// OnHeartbeat is called on the CSMS whenever a HeartbeatResponse is received from a charging station.
func (h *handler) OnHeartbeat(chargingStationID string, request *availability.HeartbeatRequest) (response *availability.HeartbeatResponse, err error) {
	fmt.Println("OnHeartbeat: ", chargingStationID, request)

	return availability.NewHeartbeatResponse(*types.NewDateTime(time.Now())), nil
}

// OnStatusNotification is called on the CSMS whenever a StatusNotificationRequest is received from a charging station.
func (h *handler) OnStatusNotification(chargingStationID string, request *availability.StatusNotificationRequest) (response *availability.StatusNotificationResponse, err error) {
	fmt.Println("OnStatusNotification: ", chargingStationID, request)

	return availability.NewStatusNotificationResponse(), nil
}

// OnDataTransfer is called on the CSMS whenever a DataTransferRequest is received from a charging station.
func (h *handler) OnDataTransfer(chargingStationID string, request *data.DataTransferRequest) (confirmation *data.DataTransferResponse, err error) {
	fmt.Println("OnDataTransfer: ", chargingStationID, request)

	return data.NewDataTransferResponse(data.DataTransferStatusUnknownVendorId), nil
}

// OnLogStatusNotification is called on the CSMS whenever a LogStatusNotificationRequest is received from a Charging Station.
func (h *handler) OnLogStatusNotification(chargingStationID string, request *diagnostics.LogStatusNotificationRequest) (response *diagnostics.LogStatusNotificationResponse, err error) {
	fmt.Println("OnLogStatusNotification: ", chargingStationID, request)

	return diagnostics.NewLogStatusNotificationResponse(), nil
}

// OnNotifyCustomerInformation is called on the CSMS whenever a NotifyCustomerInformationRequest is received from a Charging Station.
func (h *handler) OnNotifyCustomerInformation(chargingStationID string, request *diagnostics.NotifyCustomerInformationRequest) (response *diagnostics.NotifyCustomerInformationResponse, err error) {
	fmt.Println("OnNotifyCustomerInformation: ", chargingStationID, request)

	return diagnostics.NewNotifyCustomerInformationResponse(), nil
}

// OnNotifyEvent is called on the CSMS whenever a NotifyEventRequest is received from a Charging Station.
func (h *handler) OnNotifyEvent(chargingStationID string, request *diagnostics.NotifyEventRequest) (response *diagnostics.NotifyEventResponse, err error) {
	fmt.Println("OnNotifyEvent: ", chargingStationID, request)

	return diagnostics.NewNotifyEventResponse(), nil
}

// OnNotifyMonitoringReport is called on the CSMS whenever a NotifyMonitoringReportRequest is received from a Charging Station.
func (h *handler) OnNotifyMonitoringReport(chargingStationID string, request *diagnostics.NotifyMonitoringReportRequest) (response *diagnostics.NotifyMonitoringReportResponse, err error) {
	fmt.Println("OnNotifyMonitoringReport: ", chargingStationID, request)

	return diagnostics.NewNotifyMonitoringReportResponse(), nil
}

// OnNotifyDisplayMessages is called on the CSMS whenever a NotifyDisplayMessagesRequest is received from a Charging Station.
func (h *handler) OnNotifyDisplayMessages(chargingStationID string, request *display.NotifyDisplayMessagesRequest) (response *display.NotifyDisplayMessagesResponse, err error) {
	fmt.Println("OnNotifyDisplayMessages: ", chargingStationID, request)

	return display.NewNotifyDisplayMessagesResponse(), nil
}

// OnFirmwareStatusNotification is called on the CSMS whenever a FirmwareStatusNotificationRequest is received from a charging station.
func (h *handler) OnFirmwareStatusNotification(chargingStationID string, request *firmware.FirmwareStatusNotificationRequest) (response *firmware.FirmwareStatusNotificationResponse, err error) {
	fmt.Println("OnFirmwareStatusNotification: ", chargingStationID, request)

	return firmware.NewFirmwareStatusNotificationResponse(), nil
}

// OnPublishFirmwareStatusNotification is called on the CSMS whenever a PublishFirmwareStatusNotificationRequest is received from a local controller.
func (h *handler) OnPublishFirmwareStatusNotification(chargingStationID string, request *firmware.PublishFirmwareStatusNotificationRequest) (response *firmware.PublishFirmwareStatusNotificationResponse, err error) {
	fmt.Println("OnPublishFirmwareStatusNotification: ", chargingStationID, request)

	return firmware.NewPublishFirmwareStatusNotificationResponse(), nil
}

// OnGet15118EVCertificate is called on the CSMS whenever a Get15118EVCertificateRequest is received from a charging station.
func (h *handler) OnGet15118EVCertificate(chargingStationID string, request *iso15118.Get15118EVCertificateRequest) (response *iso15118.Get15118EVCertificateResponse, err error) {
	fmt.Println("OnGet15118EVCertificate: ", chargingStationID, request)
	return nil, nil
}

// OnGetCertificateStatus is called on the CSMS whenever a GetCertificateStatusRequest is received from a charging station.
func (h *handler) OnGetCertificateStatus(chargingStationID string, request *iso15118.GetCertificateStatusRequest) (response *iso15118.GetCertificateStatusResponse, err error) {
	fmt.Println("OnGetCertificateStatus: ", chargingStationID, request)
	return nil, nil
}

// OnMeterValues is called on the CSMS whenever a MeterValuesRequest is received from a charging station.
func (h *handler) OnMeterValues(chargingStationID string, request *meter.MeterValuesRequest) (response *meter.MeterValuesResponse, err error) {
	fmt.Println("OnMeterValues: ", chargingStationID, request)
	return meter.NewMeterValuesResponse(), nil
}

// OnBootNotification is called on the CSMS whenever a BootNotificationRequest is received from a charging station.
func (h *handler) OnBootNotification(chargingStationID string, request *provisioning.BootNotificationRequest) (confirmation *provisioning.BootNotificationResponse, err error) {
	fmt.Println("OnBootNotification: ", chargingStationID, request)

	fmt.Println("Send chargingProfile")
	go func() {
		profile := &types.ChargingProfile{
			ID:                     1,
			StackLevel:             2,
			ChargingProfilePurpose: types.ChargingProfilePurposeTxProfile,
			ChargingProfileKind:    types.ChargingProfileKindRelative,
			ChargingSchedule: []types.ChargingSchedule{
				{
					ChargingSchedulePeriod: []types.ChargingSchedulePeriod{
						{
							Limit:       11000,
							StartPeriod: 0,
						},
					},
					ID:               42,
					ChargingRateUnit: types.ChargingRateUnitWatts,
				},
			},
		}
		evseID := 1
		// req := smartcharging.NewSetChargingProfileRequest(evseID, profile)
		time.Sleep(time.Second * 5)
		h.csms.SetChargingProfile(chargingStationID, func(scpr *smartcharging.SetChargingProfileResponse, e error) {
			fmt.Println("Response receiiived ", scpr)
		}, evseID, profile)
	}()

	return provisioning.NewBootNotificationResponse(types.NewDateTime(time.Now()), int(heartBeatInterval), provisioning.RegistrationStatusAccepted), nil
}

// OnNotifyReport is called on the CSMS whenever a NotifyReportRequest is received from a charging station.
func (h *handler) OnNotifyReport(chargingStationID string, request *provisioning.NotifyReportRequest) (confirmation *provisioning.NotifyReportResponse, err error) {
	fmt.Println("OnNotifyReport: ", chargingStationID, request)
	return provisioning.NewNotifyReportResponse(), nil
}

// OnReservationStatusUpdate is called on the CSMS whenever a ReservationStatusUpdateRequest is received from a charging station.
func (h *handler) OnReservationStatusUpdate(chargingStationID string, request *reservation.ReservationStatusUpdateRequest) (resp *reservation.ReservationStatusUpdateResponse, err error) {
	fmt.Println("OnReservationStatusUpdate: ", chargingStationID, request)
	return reservation.NewReservationStatusUpdateResponse(), nil
}

// OnSecurityEventNotification is called on the CSMS whenever a SecurityEventNotificationRequest is received from a charging station.
func (h *handler) OnSecurityEventNotification(chargingStationID string, request *security.SecurityEventNotificationRequest) (response *security.SecurityEventNotificationResponse, err error) {
	fmt.Println("OnSecurityEventNotification: ", chargingStationID, request)
	return security.NewSecurityEventNotificationResponse(), nil
}

// OnSignCertificate is called on the CSMS whenever a SignCertificateRequest is received from a charging station.
func (h *handler) OnSignCertificate(chargingStationID string, request *security.SignCertificateRequest) (response *security.SignCertificateResponse, err error) {
	fmt.Println("OnSignCertificate: ", chargingStationID, request)
	return nil, nil
}

// OnClearedChargingLimit is called on the CSMS whenever a ClearedChargingLimitRequest is received from a charging station.
func (h *handler) OnClearedChargingLimit(chargingStationID string, request *smartcharging.ClearedChargingLimitRequest) (response *smartcharging.ClearedChargingLimitResponse, err error) {
	fmt.Println("OnClearedChargingLimit: ", chargingStationID, request)
	return smartcharging.NewClearedChargingLimitResponse(), nil
}

// OnNotifyChargingLimit is called on the CSMS whenever a NotifyChargingLimitRequest is received from a charging station.
func (h *handler) OnNotifyChargingLimit(chargingStationID string, request *smartcharging.NotifyChargingLimitRequest) (response *smartcharging.NotifyChargingLimitResponse, err error) {
	fmt.Println("OnNotifyChargingLimit: ", chargingStationID, request)
	return smartcharging.NewNotifyChargingLimitResponse(), nil
}

// OnNotifyEVChargingNeeds is called on the CSMS whenever a NotifyEVChargingNeedsRequest is received from a charging station.
func (h *handler) OnNotifyEVChargingNeeds(chargingStationID string, request *smartcharging.NotifyEVChargingNeedsRequest) (response *smartcharging.NotifyEVChargingNeedsResponse, err error) {
	fmt.Println("OnNotifyEVChargingNeeds: ", chargingStationID, request)

	return smartcharging.NewNotifyEVChargingNeedsResponse(smartcharging.EVChargingNeedsStatusAccepted), nil
}

// OnNotifyEVChargingSchedule is called on the CSMS whenever a NotifyEVChargingScheduleRequest is received from a charging station.
func (h *handler) OnNotifyEVChargingSchedule(chargingStationID string, request *smartcharging.NotifyEVChargingScheduleRequest) (response *smartcharging.NotifyEVChargingScheduleResponse, err error) {
	fmt.Println("OnNotifyEVChargingSchedule: ", chargingStationID, request)
	return nil, nil
}

// OnReportChargingProfiles is called on the CSMS whenever a ReportChargingProfilesRequest is received from a charging station.
func (h *handler) OnReportChargingProfiles(chargingStationID string, request *smartcharging.ReportChargingProfilesRequest) (reponse *smartcharging.ReportChargingProfilesResponse, err error) {
	fmt.Println("OnReportChargingProfiles: ", chargingStationID, request)
	return smartcharging.NewReportChargingProfilesResponse(), nil
}

// OnTransactionEvent is called on the CSMS whenever a TransactionEventRequest is received from a charging station.
func (h *handler) OnTransactionEvent(chargingStationID string, request *transactions.TransactionEventRequest) (response *transactions.TransactionEventResponse, err error) {
	fmt.Println("OnTransactionEvent: ", chargingStationID, request)

	return transactions.NewTransactionEventResponse(), nil
}
