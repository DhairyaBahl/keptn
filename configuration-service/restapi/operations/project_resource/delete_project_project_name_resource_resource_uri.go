// Code generated by go-swagger; DO NOT EDIT.

package project_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteProjectProjectNameResourceResourceURIHandlerFunc turns a function with the right signature into a delete project project name resource resource URI handler
type DeleteProjectProjectNameResourceResourceURIHandlerFunc func(DeleteProjectProjectNameResourceResourceURIParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteProjectProjectNameResourceResourceURIHandlerFunc) Handle(params DeleteProjectProjectNameResourceResourceURIParams) middleware.Responder {
	return fn(params)
}

// DeleteProjectProjectNameResourceResourceURIHandler interface for that can handle valid delete project project name resource resource URI params
type DeleteProjectProjectNameResourceResourceURIHandler interface {
	Handle(DeleteProjectProjectNameResourceResourceURIParams) middleware.Responder
}

// NewDeleteProjectProjectNameResourceResourceURI creates a new http.Handler for the delete project project name resource resource URI operation
func NewDeleteProjectProjectNameResourceResourceURI(ctx *middleware.Context, handler DeleteProjectProjectNameResourceResourceURIHandler) *DeleteProjectProjectNameResourceResourceURI {
	return &DeleteProjectProjectNameResourceResourceURI{Context: ctx, Handler: handler}
}

/*DeleteProjectProjectNameResourceResourceURI swagger:route DELETE /project/{projectName}/resource/{resourceURI} Project Resource deleteProjectProjectNameResourceResourceUri

Delete the specified resource

*/
type DeleteProjectProjectNameResourceResourceURI struct {
	Context *middleware.Context
	Handler DeleteProjectProjectNameResourceResourceURIHandler
}

func (o *DeleteProjectProjectNameResourceResourceURI) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteProjectProjectNameResourceResourceURIParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
