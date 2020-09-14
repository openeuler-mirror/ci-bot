package cibot

import "testing"

import (
	"testing"
)

func Test_formatDescription(t *testing.T) {
	type args struct {
		user      string
		reviewers []string
		signers   []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test description format",
			args: args{
				user:      "fakeuser",
				reviewers: []string{"@aaa", "@bbb"},
				signers:   []string{"@ccc", "@ddd"},
			},
			want: "From: @fakeuser\nReviewed-by: @aaa,@bbb\nSigned-off-by: @ccc,@ddd\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDescription(tt.args.user, tt.args.reviewers, tt.args.signers); got != tt.want {
				t.Errorf("formatDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}
