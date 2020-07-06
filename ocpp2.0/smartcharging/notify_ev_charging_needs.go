package smartcharging

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Notify Charging Limit (CS -> CSMS) --------------------

const NotifyEVChargingNeedsFeatureName = "NotifyEVChargingNeeds"

type EnergyTransferMode string

const (
	EnergyTransferModeDC            EnergyTransferMode = "DC"              // DC charging.
	EnergyTransferModeACSinglePhase EnergyTransferMode = "AC_single_phase" // AC single phase charging according to IEC 62196.
	EnergyTransferModeACTwoPhase    EnergyTransferMode = "AC_two_phase"    // AC two phase charging according to IEC 62196.
	EnergyTransferModeACThreePhase  EnergyTransferMode = "AC_three_phase"  // AC three phase charging according to IEC 62196.
)

func isValidEnergyTransferMode(fl validator.FieldLevel) bool {
	status := EnergyTransferMode(fl.Field().String())
	switch status {
	case EnergyTransferModeDC, EnergyTransferModeACSinglePhase, EnergyTransferModeACTwoPhase, EnergyTransferModeACThreePhase:
		return true
	default:
		return false
	}
}

type ACChargingParameters struct {
	EnergyAmount int `json:"energyAmount" validate:"gte=0"` // Amount of energy requested (in Wh). This includes energy required for preconditioning.
	EvMinCurrent int `json:"evMinCurrent" validate:"gte=0"` // Minimum current (amps) supported by the electric vehicle (per phase).
	EvMaxCurrent int `json:"evMaxCurrent" validate:"gte=0"` // Maximum current (amps) supported by the electric vehicle (per phase). Includes cable capacity.
	EvMaxVoltage int `json:"evMaxVoltage" validate:"gte=0"` // Maximum voltage supported by the electric vehicle.
}

type DCChargingParameters struct {
	EvMaxCurrent     int  `json:"evMaxCurrent" validate:"gte=0"`                              // Maximum current (amps) supported by the electric vehicle. Includes cable capacity.
	EvMaxVoltage     int  `json:"evMaxVoltage" validate:"gte=0"`                              // Maximum voltage supported by the electric vehicle.
	EnergyAmount     *int `json:"energyAmount,omitempty" validate:"omitempty,gte=0"`          // Amount of energy requested (in Wh). This inludes energy required for preconditioning.
	EvMaxPower       *int `json:"evMaxPower,omitempty" validate:"omitempty,gte=0"`            // Maximum power (in W) supported by the electric vehicle. Required for DC charging.
	StateOfCharge    *int `json:"stateOfCharge,omitempty" validate:"omitempty,min=0,max=100"` // Energy available in the battery (in percent of the battery capacity).
	EvEnergyCapacity *int `json:"energyCapacity,omitempty" validate:"omitempty,gte=0"`        // Capacity of the electric vehicle battery (in Wh).
	FullSoC          *int `json:"fullSoC,omitempty" validate:"omitempty,min=0,max=100"`       // Percentage of SoC at which the EV considers the battery fully charged. (possible values: 0 - 100).
	BulkSoC          *int `json:"bulkSoC,omitempty" validate:"omitempty,min=0,max=100"`       // Percentage of SoC at which the EV considers a fast charging process to end. (possible values: 0 - 100).
}

// ChargingNeeds contains the charging needs requested from the EV.
type ChargingNeeds struct {
	RequestedEnergyTransfer EnergyTransferMode    `json:"requestedEnergyTransfer" validate:"required,energyTransferMode"` // Mode of energy transfer requested by the EV.
	DepartureTime           *types.DateTime       `json:"departureTime,omitempty" validate:"omitempty"`                   // Estimated departure time of the EV.
	ACChargingParameters    *ACChargingParameters `json:"acChargingParameters,omitempty" validate:"omitempty"`            // EV AC charging parameters.
	DCChargingParameters    *DCChargingParameters `json:"dcChargingParameters,omitempty" validate:"omitempty"`            // EV DC charging parameters.
}

// NotifyEVChargingNeedsStatus is used within a NotifyEVChargingNeedsResponse.
type NotifyEVChargingNeedsStatus string

const (
	NotifyEVChargingNeedsStatusAccepted   NotifyEVChargingNeedsStatus = "Accepted"   // A schedule will be provided momentarily.
	NotifyEVChargingNeedsStatusRejected   NotifyEVChargingNeedsStatus = "Rejected"   // Service not available.
	NotifyEVChargingNeedsStatusProcessing NotifyEVChargingNeedsStatus = "Processing" // The CSMS is gathering information to provide a schedule.
)

func isValidNotifyEVChargingNeedsStatus(fl validator.FieldLevel) bool {
	status := NotifyEVChargingNeedsStatus(fl.Field().String())
	switch status {
	case NotifyEVChargingNeedsStatusAccepted, NotifyEVChargingNeedsStatusRejected, NotifyEVChargingNeedsStatusProcessing:
		return true
	default:
		return false
	}
}

// The field definition of the NotifyEVChargingNeeds request payload sent by the Charging Station to the CSMS.
type NotifyEVChargingNeedsRequest struct {
	MaxScheduleTuples *int          `json:"maxScheduleTuples,omitempty" validate:"omitempty,gte=0"` // Contains the maximum schedule tuples the car supports per schedule.
	EvseID            int           `json:"evseId" validate:"gt=0"`                                 // Defines the EVSE and connector to which the EV is connected. EvseId may not be 0.
	ChargingNeeds     ChargingNeeds `json:"chargingNeeds" validate:"required"`                      // The characteristics of the energy delivery required.
}

// This field definition of the NotifyEVChargingNeeds response payload, sent by the CSMS to the Charging Station in response to a NotifyEVChargingNeedsRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type NotifyEVChargingNeedsResponse struct {
	Status     NotifyEVChargingNeedsStatus `json:"status" validate:"required,notifyEvChargingNeedsStatus"` // Returns whether the CSMS has been able to process the message successfully. It does not imply that the evChargingNeeds can be met with the current charging profile.
	StatusInfo *types.StatusInfo           `json:"statusInfo,omitempty" validate:"omitempty"`              // Detailed status information.
}

// The EV sends a ChargeParameterDiscoveryReq message to the Charging Station.
// Whenever an EV sends a ChargeParameterDiscoveryReq with with charging needs parameters,
// the Charging Station then sends this information in a NotifyEVChargingNeedsRequest to the CSMS.
// The CSMS responds with a NotifyEVChargingNeedsResponse to the Charging Station.
//
// After responding, the CSMS calculates new charging schedule, that tries to accomodate the EV charging needs and
// still fits within the schedule boundaries imposed by other parameters.
// The CSMS will asynchronously send a SetChargingProfileRequest message with the updated schedule to the Charging Station.
type NotifyEVChargingNeedsFeature struct{}

func (f NotifyEVChargingNeedsFeature) GetFeatureName() string {
	return NotifyEVChargingNeedsFeatureName
}

func (f NotifyEVChargingNeedsFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyEVChargingNeedsRequest{})
}

func (f NotifyEVChargingNeedsFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyEVChargingNeedsResponse{})
}

func (r NotifyEVChargingNeedsRequest) GetFeatureName() string {
	return NotifyEVChargingNeedsFeatureName
}

func (c NotifyEVChargingNeedsResponse) GetFeatureName() string {
	return NotifyEVChargingNeedsFeatureName
}

// Creates a new NotifyEVChargingNeedsRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifyEVChargingNeedsRequest(evseID int, chargingNeeds ChargingNeeds) *NotifyEVChargingNeedsRequest {
	return &NotifyEVChargingNeedsRequest{EvseID: evseID, ChargingNeeds: chargingNeeds}
}

// Creates a new NotifyEVChargingNeedsResponse, containing all required fields. Optional fields may be set afterwards.
func NewNotifyEVChargingNeedsResponse(status NotifyEVChargingNeedsStatus) *NotifyEVChargingNeedsResponse {
	return &NotifyEVChargingNeedsResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("energyTransferMode", isValidEnergyTransferMode)
	_ = types.Validate.RegisterValidation("notifyEvChargingNeedsStatus", isValidNotifyEVChargingNeedsStatus)
}
