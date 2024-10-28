// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package media

import (
	"path/filepath"
)

var (
	//Default www static file dir
	DefaultHTTPDir = "static_source/media"
)

// ServerHTTPDir
func (obj *StorageST) ServerHTTPDir() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if filepath.Clean(obj.Server.HTTPDir) == "." {
		return DefaultHTTPDir
	}
	return filepath.Clean(obj.Server.HTTPDir)
}

// ServerHTTPDebug read debug options
func (obj *StorageST) ServerHTTPDebug() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPDebug
}

//// ServerLogLevel read debug options
//func (obj *StorageST) ServerLogLevel() logrus.Level {
//	obj.mutex.RLock()
//	defer obj.mutex.RUnlock()
//	return obj.Server.LogLevel
//}

// ServerHTTPDemo read demo options
func (obj *StorageST) ServerHTTPDemo() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPDemo
}

// ServerHTTPLogin read Login options
func (obj *StorageST) ServerHTTPLogin() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPLogin
}

// ServerHTTPPassword read Password options
func (obj *StorageST) ServerHTTPPassword() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPPassword
}

// ServerHTTPPort read HTTP Port options
func (obj *StorageST) ServerHTTPPort() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPPort
}

// ServerRTSPPort read HTTP Port options
func (obj *StorageST) ServerRTSPPort() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.RTSPPort
}

// ServerHTTPS read HTTPS Port options
func (obj *StorageST) ServerHTTPS() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPS
}

// ServerHTTPSPort read HTTPS Port options
func (obj *StorageST) ServerHTTPSPort() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSPort
}

// ServerHTTPSAutoTLSEnable read HTTPS Port options
func (obj *StorageST) ServerHTTPSAutoTLSEnable() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSAutoTLSEnable
}

// ServerHTTPSAutoTLSName read HTTPS Port options
func (obj *StorageST) ServerHTTPSAutoTLSName() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSAutoTLSName
}

// ServerHTTPSCert read HTTPS Cert options
func (obj *StorageST) ServerHTTPSCert() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSCert
}

// ServerHTTPSKey read HTTPS Key options
func (obj *StorageST) ServerHTTPSKey() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSKey
}

// ServerICEServers read ICE servers
func (obj *StorageST) ServerICEServers() []string {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.ICEServers
}

// ServerICEServers read ICE username
func (obj *StorageST) ServerICEUsername() string {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.ICEUsername
}

// ServerICEServers read ICE credential
func (obj *StorageST) ServerICECredential() string {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.ICECredential
}

// ServerTokenEnable read HTTPS Key options
func (obj *StorageST) ServerTokenEnable() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.Token.Enable
}

// ServerTokenBackend read HTTPS Key options
func (obj *StorageST) ServerTokenBackend() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.Token.Backend
}

// ServerWebRTCPortMin read WebRTC Port Min
func (obj *StorageST) ServerWebRTCPortMin() uint16 {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.WebRTCPortMin
}

// ServerWebRTCPortMax read WebRTC Port Max
func (obj *StorageST) ServerWebRTCPortMax() uint16 {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.WebRTCPortMax
}
