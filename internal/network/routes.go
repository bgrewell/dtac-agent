package network

import (
	"encoding/json"
	"errors"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
)

func (s *Subsystem) getRoutesHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		rt, err := GetRouteTable()
		if err != nil {
			return nil, err
		}

		return json.Marshal(rt)
	}, "route table entries")
}

func (s *Subsystem) getRouteHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		return nil, errors.New("this function has not been migrated yet")
	}, "")
}

func (s *Subsystem) createRouteHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
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
		rt, err := GetRouteTable()
		if err != nil {
			return nil, err
		}

		return json.Marshal(rt)
	}, "route has been created")
}

func (s *Subsystem) updateRouteHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
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
		rt, err := GetRouteTable()
		if err != nil {
			return nil, err
		}

		return json.Marshal(rt)
	}, "route has been updated")
}

func (s *Subsystem) deleteRouteHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		var row RouteTableRow

		// Transform the body into a RouteTableRow
		if err = json.Unmarshal(in.Body, &row); err != nil {
			return nil, err
		}

		// Delete the route
		if err = DeleteRoute(row); err != nil {
			return nil, err
		}

		rt, err := GetRouteTable()
		if err != nil {
			return nil, err
		}

		return json.Marshal(rt)
	}, "route has been deleted")
}

func (s *Subsystem) getRoutesUnifiedHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		return nil, errors.New("this function has not been migrated yet")
	}, "")
}

func (s *Subsystem) getRouteUnifiedHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		return nil, errors.New("this function has not been migrated yet")
	}, "")
}

func (s *Subsystem) createRouteUnifiedHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		return nil, errors.New("this function has not been migrated yet")
	}, "")
}

func (s *Subsystem) updateRouteUnifiedHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		return nil, errors.New("this function has not been migrated yet")
	}, "")
}
func (s *Subsystem) deleteRouteUnifiedHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		return nil, errors.New("this function has not been migrated yet")
	}, "")
}
