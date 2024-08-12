// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package api

import (
	"encoding/json"
	"fmt"
	"time"
)

// ChatRoom defines model for ChatRoom.
type ChatRoom struct {
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	Id        *string    `json:"id,omitempty"`
	Name      string     `json:"name"`
}

// CreateChatRoomRequest defines model for CreateChatRoomRequest.
type CreateChatRoomRequest struct {
	Name string `json:"name"`
}

// CreateChatRoomResponse defines model for CreateChatRoomResponse.
type CreateChatRoomResponse struct {
	ChatRoom *ChatRoom `json:"chatRoom,omitempty"`
}

// GoogleProtobufAny Contains an arbitrary serialized message along with a @type that describes the type of the serialized message.
type GoogleProtobufAny struct {
	// Type The type of the serialized message.
	Type                 *string                `json:"@type,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// Status The `Status` type defines a logical error model that is suitable for different programming environments, including REST APIs and RPC APIs. It is used by [gRPC](https://github.com/grpc). Each `Status` message contains three pieces of data: error code, error message, and error details. You can find out more about this error model and how to work with it in the [API Design Guide](https://cloud.google.com/apis/design/errors).
type Status struct {
	// Code The status code, which should be an enum value of [google.rpc.Code][google.rpc.Code].
	Code *int32 `json:"code,omitempty"`

	// Details A list of messages that carry the error details.  There is a common set of message types for APIs to use.
	Details *[]GoogleProtobufAny `json:"details,omitempty"`

	// Message A developer-facing error message, which should be in English. Any user-facing error message should be localized and sent in the [google.rpc.Status.details][google.rpc.Status.details] field, or localized by the client.
	Message *string `json:"message,omitempty"`
}

// ChatServiceCreateChatRoomJSONRequestBody defines body for ChatServiceCreateChatRoom for application/json ContentType.
type ChatServiceCreateChatRoomJSONRequestBody = CreateChatRoomRequest

// Getter for additional properties for GoogleProtobufAny. Returns the specified
// element and whether it was found
func (a GoogleProtobufAny) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for GoogleProtobufAny
func (a *GoogleProtobufAny) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for GoogleProtobufAny to handle AdditionalProperties
func (a *GoogleProtobufAny) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["@type"]; found {
		err = json.Unmarshal(raw, &a.Type)
		if err != nil {
			return fmt.Errorf("error reading '@type': %w", err)
		}
		delete(object, "@type")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for GoogleProtobufAny to handle AdditionalProperties
func (a GoogleProtobufAny) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	if a.Type != nil {
		object["@type"], err = json.Marshal(a.Type)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '@type': %w", err)
		}
	}

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}
