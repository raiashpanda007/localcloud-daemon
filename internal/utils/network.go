package utils

import "net"

func GetRandomPort() (net.Listener, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	return listener, nil
}
