// Code generated by go-swagger; DO NOT EDIT.

package types

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"strconv"
)

// Validate validates this HTTP error
func (m *HTTPError) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTitle(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HTTPError) validateCode(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", m.Code); err != nil {
		return err
	}

	if err := validate.MinimumInt("status", "body", int64(*m.Code), 100, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("status", "body", int64(*m.Code), 599, false); err != nil {
		return err
	}

	return nil
}

func (m *HTTPError) validateTitle(formats strfmt.Registry) error {

	if err := validate.Required("title", "body", m.Title); err != nil {
		return err
	}

	return nil
}

func (m *HTTPError) validateType(formats strfmt.Registry) error {

	if err := validate.Required("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HTTPError) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HTTPError) UnmarshalBinary(b []byte) error {
	var res HTTPError
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// Validate validates this HTTP validation error
func (m *HTTPValidationError) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTitle(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateValidationErrors(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HTTPValidationError) validateCode(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", m.Code); err != nil {
		return err
	}

	if err := validate.MinimumInt("status", "body", int64(*m.Code), 100, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("status", "body", int64(*m.Code), 599, false); err != nil {
		return err
	}

	return nil
}

func (m *HTTPValidationError) validateTitle(formats strfmt.Registry) error {

	if err := validate.Required("title", "body", m.Title); err != nil {
		return err
	}

	return nil
}

func (m *HTTPValidationError) validateType(formats strfmt.Registry) error {

	if err := validate.Required("type", "body", m.Type); err != nil {
		return err
	}

	return nil
}

func (m *HTTPValidationError) validateValidationErrors(formats strfmt.Registry) error {

	if err := validate.Required("validationErrors", "body", m.ValidationErrors); err != nil {
		return err
	}

	for i := 0; i < len(m.ValidationErrors); i++ {
		if swag.IsZero(m.ValidationErrors[i]) { // not required
			continue
		}

		if m.ValidationErrors[i] != nil {
			if err := m.ValidationErrors[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("validationErrors" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *HTTPValidationError) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HTTPValidationError) UnmarshalBinary(b []byte) error {
	var res HTTPValidationError
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// Validate validates this HTTP validation error detail
func (m *HTTPValidationErrorDetail) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateError(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIn(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateKey(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HTTPValidationErrorDetail) validateError(formats strfmt.Registry) error {

	if err := validate.Required("error", "body", m.Error); err != nil {
		return err
	}

	return nil
}

func (m *HTTPValidationErrorDetail) validateIn(formats strfmt.Registry) error {

	if err := validate.Required("in", "body", m.In); err != nil {
		return err
	}

	return nil
}

func (m *HTTPValidationErrorDetail) validateKey(formats strfmt.Registry) error {

	if err := validate.Required("key", "body", m.Key); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HTTPValidationErrorDetail) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HTTPValidationErrorDetail) UnmarshalBinary(b []byte) error {
	var res HTTPValidationErrorDetail
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// Validate validates this post change password payload
func (m *PostChangePasswordPayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCurrentPassword(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNewPassword(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostChangePasswordPayload) validateCurrentPassword(formats strfmt.Registry) error {

	if err := validate.Required("currentPassword", "body", m.CurrentPassword); err != nil {
		return err
	}

	if err := validate.MinLength("currentPassword", "body", string(*m.CurrentPassword), 1); err != nil {
		return err
	}

	return nil
}

func (m *PostChangePasswordPayload) validateNewPassword(formats strfmt.Registry) error {

	if err := validate.Required("newPassword", "body", m.NewPassword); err != nil {
		return err
	}

	if err := validate.MinLength("newPassword", "body", string(*m.NewPassword), 1); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostChangePasswordPayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostChangePasswordPayload) UnmarshalBinary(b []byte) error {
	var res PostChangePasswordPayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// Validate validates this post forgot password complete payload
func (m *PostForgotPasswordCompletePayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePassword(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateToken(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostForgotPasswordCompletePayload) validatePassword(formats strfmt.Registry) error {

	if err := validate.Required("password", "body", m.Password); err != nil {
		return err
	}

	if err := validate.MinLength("password", "body", string(*m.Password), 1); err != nil {
		return err
	}

	return nil
}

func (m *PostForgotPasswordCompletePayload) validateToken(formats strfmt.Registry) error {

	if err := validate.Required("token", "body", m.Token); err != nil {
		return err
	}

	if err := validate.FormatOf("token", "body", "uuid4", m.Token.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostForgotPasswordCompletePayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostForgotPasswordCompletePayload) UnmarshalBinary(b []byte) error {
	var res PostForgotPasswordCompletePayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// Validate validates this post forgot password payload
func (m *PostForgotPasswordPayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateUsername(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostForgotPasswordPayload) validateUsername(formats strfmt.Registry) error {

	if err := validate.Required("username", "body", m.Username); err != nil {
		return err
	}

	if err := validate.MinLength("username", "body", string(*m.Username), 1); err != nil {
		return err
	}

	if err := validate.MaxLength("username", "body", string(*m.Username), 255); err != nil {
		return err
	}

	if err := validate.FormatOf("username", "body", "email", m.Username.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostForgotPasswordPayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostForgotPasswordPayload) UnmarshalBinary(b []byte) error {
	var res PostForgotPasswordPayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// Validate validates this post login payload
func (m *PostLoginPayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePassword(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUsername(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostLoginPayload) validatePassword(formats strfmt.Registry) error {

	if err := validate.Required("password", "body", m.Password); err != nil {
		return err
	}

	if err := validate.MinLength("password", "body", string(*m.Password), 1); err != nil {
		return err
	}

	return nil
}

func (m *PostLoginPayload) validateUsername(formats strfmt.Registry) error {

	if err := validate.Required("username", "body", m.Username); err != nil {
		return err
	}

	if err := validate.MinLength("username", "body", string(*m.Username), 1); err != nil {
		return err
	}

	if err := validate.MaxLength("username", "body", string(*m.Username), 255); err != nil {
		return err
	}

	if err := validate.FormatOf("username", "body", "email", m.Username.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostLoginPayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostLoginPayload) UnmarshalBinary(b []byte) error {
	var res PostLoginPayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// Validate validates this post login response
func (m *PostLoginResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAccessToken(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExpiresIn(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRefreshToken(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTokenType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostLoginResponse) validateAccessToken(formats strfmt.Registry) error {

	if err := validate.Required("access_token", "body", m.AccessToken); err != nil {
		return err
	}

	if err := validate.FormatOf("access_token", "body", "uuid4", m.AccessToken.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *PostLoginResponse) validateExpiresIn(formats strfmt.Registry) error {

	if err := validate.Required("expires_in", "body", m.ExpiresIn); err != nil {
		return err
	}

	return nil
}

func (m *PostLoginResponse) validateRefreshToken(formats strfmt.Registry) error {

	if err := validate.Required("refresh_token", "body", m.RefreshToken); err != nil {
		return err
	}

	if err := validate.FormatOf("refresh_token", "body", "uuid4", m.RefreshToken.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *PostLoginResponse) validateTokenType(formats strfmt.Registry) error {

	if err := validate.Required("token_type", "body", m.TokenType); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostLoginResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostLoginResponse) UnmarshalBinary(b []byte) error {
	var res PostLoginResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// Validate validates this post logout payload
func (m *PostLogoutPayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRefreshToken(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostLogoutPayload) validateRefreshToken(formats strfmt.Registry) error {

	if swag.IsZero(m.RefreshToken) { // not required
		return nil
	}

	if err := validate.FormatOf("refresh_token", "body", "uuid4", m.RefreshToken.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostLogoutPayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostLogoutPayload) UnmarshalBinary(b []byte) error {
	var res PostLogoutPayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// Validate validates this post refresh payload
func (m *PostRefreshPayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRefreshToken(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostRefreshPayload) validateRefreshToken(formats strfmt.Registry) error {

	if err := validate.Required("refresh_token", "body", m.RefreshToken); err != nil {
		return err
	}

	if err := validate.FormatOf("refresh_token", "body", "uuid4", m.RefreshToken.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostRefreshPayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostRefreshPayload) UnmarshalBinary(b []byte) error {
	var res PostRefreshPayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// Validate validates this post register payload
func (m *PostRegisterPayload) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePassword(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUsername(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostRegisterPayload) validatePassword(formats strfmt.Registry) error {

	if err := validate.Required("password", "body", m.Password); err != nil {
		return err
	}

	if err := validate.MinLength("password", "body", string(*m.Password), 1); err != nil {
		return err
	}

	return nil
}

func (m *PostRegisterPayload) validateUsername(formats strfmt.Registry) error {

	if err := validate.Required("username", "body", m.Username); err != nil {
		return err
	}

	if err := validate.MinLength("username", "body", string(*m.Username), 1); err != nil {
		return err
	}

	if err := validate.MaxLength("username", "body", string(*m.Username), 255); err != nil {
		return err
	}

	if err := validate.FormatOf("username", "body", "email", m.Username.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostRegisterPayload) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostRegisterPayload) UnmarshalBinary(b []byte) error {
	var res PostRegisterPayload
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}