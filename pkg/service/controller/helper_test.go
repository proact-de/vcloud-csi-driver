package controller

import (
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

func TestSizeFromCapacityRange(t *testing.T) {
	testCases := []struct {
		Name          string
		CapacityRange *csi.CapacityRange
		Size          int64
		MaxSize       int64
		OK            bool
	}{
		{
			Name:          "without capacity range",
			CapacityRange: nil,
			Size:          DefaultVolumeSize,
			MaxSize:       0,
			OK:            true,
		},
		{
			Name:          "empty capacity range",
			CapacityRange: &csi.CapacityRange{},
			Size:          DefaultVolumeSize,
			MaxSize:       0,
			OK:            true,
		},
		{
			Name: "with required bytes less than minimum size",
			CapacityRange: &csi.CapacityRange{
				RequiredBytes: MinVolumeSize - 100,
			},
			Size:    MinVolumeSize,
			MaxSize: 0,
			OK:      true,
		},
		{
			Name: "with required bytes exactly minimum size",
			CapacityRange: &csi.CapacityRange{
				RequiredBytes: MinVolumeSize,
			},
			Size:    MinVolumeSize,
			MaxSize: 0,
			OK:      true,
		},
		{
			Name: "with required bytes slightly more than minimum size",
			CapacityRange: &csi.CapacityRange{
				RequiredBytes: MinVolumeSize + 1000,
			},
			Size:    MinVolumeSize + 1000,
			MaxSize: 0,
			OK:      true,
		},
		{
			Name: "with required and limit bytes",
			CapacityRange: &csi.CapacityRange{
				RequiredBytes: MinVolumeSize + 1000,
				LimitBytes:    2 * MinVolumeSize,
			},
			Size:    MinVolumeSize + 1000,
			MaxSize: 2 * MinVolumeSize,
			OK:      true,
		},
		{
			Name: "with required and limit bytes (same value)",
			CapacityRange: &csi.CapacityRange{
				RequiredBytes: MinVolumeSize + 1000,
				LimitBytes:    MinVolumeSize + 1000,
			},
			Size:    MinVolumeSize + 1000,
			MaxSize: MinVolumeSize + 1000,
			OK:      true,
		},
		{
			Name: "with lower limit than required",
			CapacityRange: &csi.CapacityRange{
				RequiredBytes: 2 * MinVolumeSize,
				LimitBytes:    MinVolumeSize,
			},
			Size:    0,
			MaxSize: 0,
			OK:      false,
		},
		{
			Name: "with invalid required bytes",
			CapacityRange: &csi.CapacityRange{
				RequiredBytes: -10,
				LimitBytes:    1,
			},
			Size:    0,
			MaxSize: 0,
			OK:      false,
		},
		{
			Name: "with invalid limit bytes",
			CapacityRange: &csi.CapacityRange{
				RequiredBytes: 1,
				LimitBytes:    -10,
			},
			Size:    0,
			MaxSize: 0,
			OK:      false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			size, err := sizeFromCapacityRange(testCase.CapacityRange)

			if size != testCase.Size {
				t.Errorf("got %q, want %q", size, testCase.Size)
			}

			if err != nil && testCase.OK {
				t.Fatalf("got %s, want nil", err)
			}
		})
	}
}
