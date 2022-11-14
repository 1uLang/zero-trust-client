package utils

import "github.com/shirou/gopsutil/host"

func GetDeviceId() (string, error) {
	i, err := host.Info()
	if err != nil {
		return "", err
	}
	return i.HostID, nil
}
