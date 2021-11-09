package tests

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/ipld/go-ipld-prime/node/basicnode"
)

func Test_nodeContentEqualsChecker_Check(t *testing.T) {
	someNode := basicnode.Prototype__String{}.NewBuilder().Build()
	nb := basicnode.Prototype__String{}.NewBuilder()
	err := nb.AssignString("fish")
	qt.Assert(t, err, qt.IsNil)
	someOtherNode := nb.Build()

	tests := []struct {
		name    string
		got     interface{}
		want    interface{}
		wantErr string
	}{
		{
			name:    "nilWantIsError",
			got:     "not a node",
			want:    nil,
			wantErr: "got non-nil value",
		},
		{
			name:    "nonNodeAsWantIsError",
			got:     "not a node",
			want:    someNode,
			wantErr: "this checker only supports checking datamodel.Node values",
		},
		{
			name:    "nonNodeAsGotIsError",
			got:     someNode,
			want:    "not a node",
			wantErr: "this checker only supports checking datamodel.Node values",
		},
		{
			name: "nilWantAndGotAreEqual",
			got:  nil,
			want: nil,
		},
		{
			name: "equivalentNodesAreEqual",
			got:  someNode,
			want: someNode,
		},
		{
			name:    "differentNodesAreNotEqual",
			got:     someNode,
			want:    someOtherNode,
			wantErr: "values are not equal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NodeContentEquals.Check(tt.got, []interface{}{tt.want}, nil)
			if tt.wantErr == "" {
				qt.Assert(t, err, qt.IsNil)
			} else {
				qt.Assert(t, err, qt.Not(qt.IsNil))
				qt.Assert(t, err.Error(), qt.Equals, tt.wantErr)
			}
		})
	}
}
