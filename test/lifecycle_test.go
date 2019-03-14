package test

// +build lifecycle

import (
	"flag"
	"testing"

	v2 "github.com/openservicebrokerapi/osb-checker/autogenerated/models"
	. "github.com/openservicebrokerapi/osb-checker/config"
	"github.com/openservicebrokerapi/osb-checker/test/common"
	"github.com/satori/go.uuid"
)

var (
	configFile string
)

func init() {
	// Parse some configuration fields from command line.
	flag.StringVar(&configFile, "config-file", "configs/config_mock.yaml", "Please specify the config file of service broker you want to test!")
	flag.Parse()

	if err := Load(configFile); err != nil {
		panic(err)
	}
}

func TestLifeCycle(t *testing.T) {
	for _, svc := range CONF.Services {
		instanceID := uuid.NewV4().String()
		bindingID := uuid.NewV4().String()
		serviceID, organizationGUID, spaceGUID :=
			svc.ServiceID, svc.OrganizationGUID, svc.SpaceGUID

		for _, operation := range svc.Operations {
			switch operation.Type {
			case "provision":
				currentPlanID := operation.PlanID
				req := &v2.ServiceInstanceProvisionRequest{
					ServiceID:        &serviceID,
					PlanID:           &currentPlanID,
					OrganizationGUID: &organizationGUID,
					SpaceGUID:        &spaceGUID,
					Parameters:       operation.Parameters,
				}

				common.TestProvision(t, instanceID, req, operation.Async)
				break
			case "update":
				currentPlanID := operation.PlanID
				req := &v2.ServiceInstanceUpdateRequest{
					ServiceID:  &serviceID,
					PlanID:     currentPlanID,
					Parameters: operation.Parameters,
				}

				common.TestUpdateInstance(t, instanceID, req, operation.Async)
				break
			case "bind":
				currentPlanID := operation.PlanID
				req := &v2.ServiceBindingRequest{
					ServiceID:  &serviceID,
					PlanID:     &currentPlanID,
					Parameters: operation.Parameters,
				}

				common.TestBind(t, instanceID, bindingID, req, operation.Async)
				break
			}
		}
	}
}
