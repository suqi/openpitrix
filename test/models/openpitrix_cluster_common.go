// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// OpenpitrixClusterCommon openpitrix cluster common
// swagger:model openpitrixClusterCommon
type OpenpitrixClusterCommon struct {

	// advanced actions
	AdvancedActions string `json:"advanced_actions,omitempty"`

	// agent installed
	AgentInstalled bool `json:"agent_installed,omitempty"`

	// backup policy
	BackupPolicy string `json:"backup_policy,omitempty"`

	// backup service
	BackupService string `json:"backup_service,omitempty"`

	// cluster id
	ClusterID string `json:"cluster_id,omitempty"`

	// custom metadata script
	CustomMetadataScript string `json:"custom_metadata_script,omitempty"`

	// custom service
	CustomService string `json:"custom_service,omitempty"`

	// delete snapshot service
	DeleteSnapshotService string `json:"delete_snapshot_service,omitempty"`

	// destroy service
	DestroyService string `json:"destroy_service,omitempty"`

	// health check
	HealthCheck string `json:"health_check,omitempty"`

	// hypervisor
	Hypervisor string `json:"hypervisor,omitempty"`

	// image id
	ImageID string `json:"image_id,omitempty"`

	// incremental backup supported
	IncrementalBackupSupported bool `json:"incremental_backup_supported,omitempty"`

	// init service
	InitService string `json:"init_service,omitempty"`

	// monitor
	Monitor string `json:"monitor,omitempty"`

	// passphraseless
	Passphraseless string `json:"passphraseless,omitempty"`

	// restart service
	RestartService string `json:"restart_service,omitempty"`

	// restore service
	RestoreService string `json:"restore_service,omitempty"`

	// role
	Role string `json:"role,omitempty"`

	// scale in service
	ScaleInService string `json:"scale_in_service,omitempty"`

	// scale out service
	ScaleOutService string `json:"scale_out_service,omitempty"`

	// server id upper bound
	ServerIDUpperBound int64 `json:"server_id_upper_bound,omitempty"`

	// start service
	StartService string `json:"start_service,omitempty"`

	// stop service
	StopService string `json:"stop_service,omitempty"`

	// upgrade service
	UpgradeService string `json:"upgrade_service,omitempty"`

	// vertical scaling policy
	VerticalScalingPolicy string `json:"vertical_scaling_policy,omitempty"`
}

// Validate validates this openpitrix cluster common
func (m *OpenpitrixClusterCommon) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *OpenpitrixClusterCommon) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OpenpitrixClusterCommon) UnmarshalBinary(b []byte) error {
	var res OpenpitrixClusterCommon
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
