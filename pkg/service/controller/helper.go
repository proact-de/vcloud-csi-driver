package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

const (
	_   = iota
	kiB = 1 << (10 * iota)
	miB
	giB
	tiB
)

func sizeFromCapacityRange(capacityRange *csi.CapacityRange) (int64, error) {
	if capacityRange == nil {
		return DefaultVolumeSize, nil
	}

	limitBytes := capacityRange.GetLimitBytes()
	limitSet := 0 < limitBytes

	requiredBytes := capacityRange.GetRequiredBytes()
	requiredSet := 0 < requiredBytes

	if !requiredSet && !limitSet {
		return DefaultVolumeSize, nil
	}

	if requiredSet && limitSet && limitBytes < requiredBytes {
		return 0, fmt.Errorf(
			"limit %v can not be less than required size%v ",
			formatHumanReadableSize(limitBytes),
			formatHumanReadableSize(requiredBytes),
		)
	}

	if requiredSet && !limitSet && requiredBytes < MinVolumeSize {
		return 0, fmt.Errorf(
			"required %v can not be less than minimum size %v",
			formatHumanReadableSize(requiredBytes),
			formatHumanReadableSize(MinVolumeSize),
		)
	}

	if limitSet && limitBytes < MinVolumeSize {
		return 0, fmt.Errorf(
			"limit %v can not be less than minimum size %v",
			formatHumanReadableSize(limitBytes),
			formatHumanReadableSize(MinVolumeSize),
		)
	}

	if requiredSet && requiredBytes > MaxVolumeSize {
		return 0, fmt.Errorf(
			"required %v can not exceed maximum size %v",
			formatHumanReadableSize(requiredBytes),
			formatHumanReadableSize(MaxVolumeSize),
		)
	}

	if !requiredSet && limitSet && limitBytes > MaxVolumeSize {
		return 0, fmt.Errorf(
			"limit %v can not exceed maximum size %v",
			formatHumanReadableSize(limitBytes),
			formatHumanReadableSize(MaxVolumeSize),
		)
	}

	if requiredSet && limitSet && requiredBytes == limitBytes {
		return requiredBytes, nil
	}

	if requiredSet {
		return requiredBytes, nil
	}

	if limitSet {
		return limitBytes, nil
	}

	return DefaultVolumeSize, nil
}

func isCapabilitySupported(cap *csi.VolumeCapability) bool {
	if cap.AccessMode == nil {
		return false
	}

	switch cap.AccessMode.Mode {
	case csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER:
		return true
	default:
		return false
	}
}

func formatHumanReadableSize(input int64) string {
	output := float64(input)
	unit := ""

	switch {
	case input >= tiB:
		output = output / tiB
		unit = "Ti"
	case input >= giB:
		output = output / giB
		unit = "Gi"
	case input >= miB:
		output = output / miB
		unit = "Mi"
	case input >= kiB:
		output = output / kiB
		unit = "Ki"
	case input == 0:
		return "0"
	}

	return strings.TrimSuffix(
		strconv.FormatFloat(
			output,
			'f',
			1,
			64,
		),
		".0",
	) + unit
}
