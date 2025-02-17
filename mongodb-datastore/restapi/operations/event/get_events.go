// Code generated by go-swagger; DO NOT EDIT.

package event

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/keptn/keptn/mongodb-datastore/models"
)

// GetEventsHandlerFunc turns a function with the right signature into a get events handler
type GetEventsHandlerFunc func(GetEventsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetEventsHandlerFunc) Handle(params GetEventsParams) middleware.Responder {
	return fn(params)
}

// GetEventsHandler interface for that can handle valid get events params
type GetEventsHandler interface {
	Handle(GetEventsParams) middleware.Responder
}

// NewGetEvents creates a new http.Handler for the get events operation
func NewGetEvents(ctx *middleware.Context, handler GetEventsHandler) *GetEvents {
	return &GetEvents{Context: ctx, Handler: handler}
}

/* GetEvents swagger:route GET /event event getEvents

Gets events from the data store, either keptnContext or project must be specified

*/
type GetEvents struct {
	Context *middleware.Context
	Handler GetEventsHandler
}

func (o *GetEvents) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetEventsParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetEventsOKBody get events o k body
//
// swagger:model GetEventsOKBody
type GetEventsOKBody struct {

	// events
	Events []*models.KeptnContextExtendedCE `json:"events"`

	// Pointer to the next page
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of the returned page
	PageSize int64 `json:"pageSize,omitempty"`

	// Total number of events
	TotalCount int64 `json:"totalCount,omitempty"`
}

// Validate validates this get events o k body
func (o *GetEventsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateEvents(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetEventsOKBody) validateEvents(formats strfmt.Registry) error {
	if swag.IsZero(o.Events) { // not required
		return nil
	}

	for i := 0; i < len(o.Events); i++ {
		if swag.IsZero(o.Events[i]) { // not required
			continue
		}

		if o.Events[i] != nil {
			if err := o.Events[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getEventsOK" + "." + "events" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getEventsOK" + "." + "events" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get events o k body based on the context it is used
func (o *GetEventsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateEvents(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetEventsOKBody) contextValidateEvents(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Events); i++ {

		if o.Events[i] != nil {
			if err := o.Events[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getEventsOK" + "." + "events" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getEventsOK" + "." + "events" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetEventsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetEventsOKBody) UnmarshalBinary(b []byte) error {
	var res GetEventsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
