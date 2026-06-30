package errors

import (
	"errors"
	"testing"
)

func TestKindOf(t *testing.T) {
	tests := []struct {
		err  error
		want Kind
	}{
		{err: ErrReadOnly, want: KindReadOnly},
		{err: ErrUnsupportedResource, want: KindUnsupportedResource},
		{err: ErrInvalidParent, want: KindInvalidParent},
		{err: NewNotFound("track", "1"), want: KindNotFound},
		{err: NewReadOnly("mock"), want: KindReadOnly},
		{err: errors.New("other"), want: KindUnknown},
		{err: &Error{Kind: KindAlreadyExists, Msg: "exists"}, want: KindAlreadyExists},
	}

	for _, tt := range tests {
		if got := KindOf(tt.err); got != tt.want {
			t.Errorf("KindOf(%v) = %v, want %v", tt.err, got, tt.want)
		}
	}
}
