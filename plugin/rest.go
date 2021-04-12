package plugin

import (
	"errors"
	gplug "github.com/hashicorp/go-plugin"
	"net/rpc"
)

type HttpRoute struct {
	Method string
	Path   string
}

type Rest interface {
	// Routes() is used to get a list of all the routes that the plugin exposes
	Routes() (routes []*HttpRoute, err error)
	// Get() is used to retrieve an entity
	Get(path string, reqHeaders map[string]string) (status int, respHeaders map[string]string, respBody []byte, err error)
	// Head() is similar to Get() but doesn't get the body. It is helpful to check if a resource has changed
	Head(path string, reqHeaders map[string]string) (status int, respHeaders map[string]string, respBody []byte, err error)
	// Put() is used to update an entity
	Put(path string, reqHeaders map[string]string, reqBody []byte) (status int, respHeaders map[string]string, respBody []byte, err error)
	// Post() is used to create an entity
	Post(path string, reqHeaders map[string]string, reqBody []byte) (status int, respHeaders map[string]string, respBody []byte, err error)
	// Patch() is used to update part of an entity
	Patch(path string, reqHeaders map[string]string, reqBody []byte) (status int, respHeaders map[string]string, respBody []byte, err error)
	// Delete() is used to remove entity
	Delete(path string, reqHeaders map[string]string) (status int, respHeaders map[string]string, respBody []byte, err error)
}

/* CLIENT SIDE RPC INTERFACE */
type RestRPC struct{ client *rpc.Client }

func (r *RestRPC) Routes() (routes []*HttpRoute, err error) {
	return nil, errors.New("this method is not yet implemented")
}

func (r *RestRPC) Get(path string, reqHeaders map[string]string) (status int, respHeaders map[string]string, respBody []byte, err error) {
	return 0, nil, nil, errors.New("this method is not yet implemented")
}

func (r *RestRPC) Head(path string, reqHeaders map[string]string) (status int, respHeaders map[string]string, respBody []byte, err error) {
	return 0, nil, nil, errors.New("this method is not yet implemented")
}

func (r *RestRPC) Put(path string, reqHeaders map[string]string, reqBody []byte) (status int, respHeaders map[string]string, respBody []byte, err error) {
	return 0, nil, nil, errors.New("this method is not yet implemented")
}

func (r *RestRPC) Post(path string, reqHeaders map[string]string, reqBody []byte) (status int, respHeaders map[string]string, respBody []byte, err error) {
	return 0, nil, nil, errors.New("this method is not yet implemented")
}

func (r *RestRPC) Patch(path string, reqHeaders map[string]string, reqBody []byte) (status int, respHeaders map[string]string, respBody []byte, err error) {
	return 0, nil, nil, errors.New("this method is not yet implemented")
}

func (r *RestRPC) Delete(path string, reqHeaders map[string]string) (status int, respHeaders map[string]string, respBody []byte, err error) {
	return 0, nil, nil, errors.New("this method is not yet implemented")
}

/* SERVER SIDE RPC INTERFACE */
type RestRPCServer struct {
	Impl Rest
}

type RestPlugin struct {
	Impl Rest
}

func (p *RestPlugin) Server(broker *gplug.MuxBroker) (interface{}, error) {
	return &RestRPCServer{Impl: p.Impl}, nil
}

func (RestPlugin) Client(b *gplug.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RestRPC{client: c}, nil
}


// TODO: Need to have System-API have a way to call plugins, then call Routes() to get the routes.
// TODO: Once we have the routes we need to dynamically configure them
// TODO: When the dynamic routes are called it should grab the path, headers and body of the request and send it to the plugin via the matching HTTP method
// TODO: The plugin should update the status, the headers (or return them unmodified), any payload and any error details