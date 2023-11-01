package network

import (
	"encoding/json"
	"errors"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
)

func (s *Subsystem) getRoutesHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return GetRouteTable()
	}, "route table entries")
}

func (s *Subsystem) getRouteHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return nil, errors.New("this function has not been migrated yet")
	}, "")
}

func (s *Subsystem) createRouteHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		var row RouteTableRow

		// Transform the body into a RouteTableRow
		if err = json.Unmarshal(in.Body, &row); err != nil {
			return nil, err
		}

		// Create the route
		if err = CreateRoute(row); err != nil {
			return nil, err
		}

		// Return the route table
		return GetRouteTable()
	}, "route has been created")
}

func (s *Subsystem) updateRouteHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		var row RouteTableRow

		// Transform the body into a RouteTableRow
		if err = json.Unmarshal(in.Body, &row); err != nil {
			return nil, err
		}

		// Update the route
		if err = UpdateRoute(row); err != nil {
			return nil, err
		}

		// Return the route table
		return GetRouteTable()
	}, "route has been updated")
}

func (s *Subsystem) deleteRouteHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		var row RouteTableRow

		// Transform the body into a RouteTableRow
		if err = json.Unmarshal(in.Body, &row); err != nil {
			return nil, err
		}

		// Delete the route
		if err = DeleteRoute(row); err != nil {
			return nil, err
		}

		return GetRouteTable()
	}, "route has been deleted")
}

func (s *Subsystem) getRoutesUnifiedHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return nil, errors.New("this function has not been migrated yet")
	}, "")
}

func (s *Subsystem) getRouteUnifiedHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return nil, errors.New("this function has not been migrated yet")
	}, "")
}

func (s *Subsystem) createRouteUnifiedHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return nil, errors.New("this function has not been migrated yet")
	}, "")
}

func (s *Subsystem) updateRouteUnifiedHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return nil, errors.New("this function has not been migrated yet")
	}, "")
}
func (s *Subsystem) deleteRouteUnifiedHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return nil, errors.New("this function has not been migrated yet")
	}, "")
}
