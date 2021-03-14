package cont

import (
	"context"
	"fmt"
	"reflect"

	"firebase.google.com/go/v4/auth"
	firebaseauth "firebase.google.com/go/v4/auth"
)

type ContextKey string

const (
	OperationNameContextKey = ContextKey("OperationName")
	IDTokenClaimsContextKey = ContextKey("IDTokenClaims")
)

func newContextKeyNotFoundErr(key ContextKey) *ErrContextKeyNotFound {
	return &ErrContextKeyNotFound{string(key)}
}

type ErrContextKeyNotFound struct {
	Key string
}

func (e *ErrContextKeyNotFound) Error() string {
	return fmt.Sprintf("context key %s is not found", e.Key)
}

func newContextUnexpectedTypeErr(key ContextKey, expected string, got interface{}) *ErrContextUnexpectedType {
	return &ErrContextUnexpectedType{
		string(key), expected, got,
	}
}

type ErrContextUnexpectedType struct {
	Key      string
	Expected string
	Got      interface{}
}

func (e *ErrContextUnexpectedType) Error() string {
	return fmt.Sprintf("expected type '%s' on context key %s, found '%s'", e.Expected, e.Key, reflect.TypeOf(e.Got).String())
}

func GetOperationName(ctx context.Context) (string, error) {
	val := ctx.Value(OperationNameContextKey)
	if val == nil {
		return "", newContextKeyNotFoundErr(OperationNameContextKey)
	}

	name, ok := val.(string)
	if !ok {
		return "", newContextUnexpectedTypeErr(OperationNameContextKey, "string", val)
	}

	return name, nil
}

func WithOperationName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, OperationNameContextKey, name)
}

func GetIDTokenClaims(ctx context.Context) (*firebaseauth.Token, error) {
	val := ctx.Value(IDTokenClaimsContextKey)
	if val == nil {
		return nil, newContextKeyNotFoundErr(IDTokenClaimsContextKey)
	}

	claims, ok := val.(*firebaseauth.Token)
	if !ok {
		return nil, newContextUnexpectedTypeErr(IDTokenClaimsContextKey, "*firebase.google.com/go/v4/auth.Token", val)
	}

	return claims, nil
}

func WithIDTokenClaims(ctx context.Context, claims *auth.Token) context.Context {
	return context.WithValue(ctx, IDTokenClaimsContextKey, claims)
}
