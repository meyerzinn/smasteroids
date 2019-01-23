package assets_test

import (
	"github.com/20zinnm/smasteroids/assets"
	"testing"
)

func Test(t *testing.T) {
	if assets.FontInterface == nil {
		t.Fail()
	}
	if assets.FontLabel == nil {
		t.Fail()
	}
	if assets.FontSubtitle == nil {
		t.Fail()
	}
	if assets.FontTitle == nil {
		t.Fail()
	}
}