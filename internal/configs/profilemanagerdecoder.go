// Code generated by github.com/gokultp/go-envparser. DO NOT EDIT.
package configs

import (
	"os"
	"strconv"
)

func (t *ProfileManager) DecodeEnv() error {
	if _recUrlStr := os.Getenv("PROFILE_MANAGER_URL"); _recUrlStr != "" {
		_recUrl := _recUrlStr
		t.URL = &_recUrl
	}
	if _recInitconnectionsStr := os.Getenv("PROFILE_MANAGER_INIT_CONN"); _recInitconnectionsStr != "" {
		_recInitconnections64, err := strconv.ParseInt(_recInitconnectionsStr, 10, 32)
		if err != nil {
			return err
		}
		_recInitconnections := int(_recInitconnections64)
		t.InitConnections = &_recInitconnections
	}
	if _recConnectioncapacityStr := os.Getenv("PROFILE_MANAGER_CONN_CAPACITY"); _recConnectioncapacityStr != "" {
		_recConnectioncapacity64, err := strconv.ParseInt(_recConnectioncapacityStr, 10, 32)
		if err != nil {
			return err
		}
		_recConnectioncapacity := int(_recConnectioncapacity64)
		t.ConnectionCapacity = &_recConnectioncapacity
	}
	if _recConnectionttlminutesStr := os.Getenv("PROFILE_MANAGER_CONN_TTL"); _recConnectionttlminutesStr != "" {
		_recConnectionttlminutes64, err := strconv.ParseInt(_recConnectionttlminutesStr, 10, 32)
		if err != nil {
			return err
		}
		_recConnectionttlminutes := int(_recConnectionttlminutes64)
		t.ConnectionTTLMinutes = &_recConnectionttlminutes
	}
	return nil
}
