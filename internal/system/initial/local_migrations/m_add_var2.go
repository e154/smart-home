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

package local_migrations

import (
	"context"

	"github.com/e154/smart-home/pkg/adaptors"
)

type MigrationAddVar2 struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationAddVar2(adaptors *adaptors.Adaptors) *MigrationAddVar2 {
	return &MigrationAddVar2{
		adaptors: adaptors,
	}
}

func (n *MigrationAddVar2) Up(ctx context.Context) error {

	AddVariableIfNotExist(n.adaptors, ctx, "certPublic", `-----BEGIN CERTIFICATE-----
MIIC+TCCAeGgAwIBAgIQTGmiSrBW3EqxA38tD9oLzTANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMB4XDTI0MDIyMDE5MjEzN1oXDTI1MDIxOTE5MjEz
N1owEjEQMA4GA1UEChMHQWNtZSBDbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBAKvuAyWSqCkamVi+qNb3S1J3Ar0LE90/g7kOUbf9uJf7dvB87khgjVOH
R3mGRaQboSb528KaOkmIeOz+aLqKJd9/H/5u8yylXDBiCO3dK9c9PZBIb5REMfke
ZqsdQuaoowCrcMeZ0PJP8RQ3fDh6COugS/m6JCfpCN6JbVpbiLlTqLefnscdxHAA
btxEBTh1gP7BOZnQHChrbwzs5E0AR/UWtqm8L4ADCFvNaDYXddd83lMCRypiFsXI
fc9JuLyKcu44zgdJ/ha44vEYdjrOuETd8DI/wqDnk1m1Kvmr8PznvN1LmJs9NBoW
oWxfq2Vv4CBOKur7D6tsJLULoW+U188CAwEAAaNLMEkwDgYDVR0PAQH/BAQDAgWg
MBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwFAYDVR0RBA0wC4IJ
bG9jYWxob3N0MA0GCSqGSIb3DQEBCwUAA4IBAQBiy/8KPt9yfaIAk2gEpvJxbuZt
Gcqd/WMpYdm3YdsLHPNma+mGUBwF05xjfjL3Li4DyS+TCpWJbsrfJvA2bBJ1Rzzo
ZEjP2IAkW0A8yTFTozd2SDzdWX807RJY8Cy/ebT7r3VD7Si6Dxm1MZif5oC/Ma3x
6WDqhGKQNLi0+dnYUAIPIEXyGDiB7rc3yC5pznzq6ELv64eCs+JnJU7U1xesFXNW
VU0GW3f/PxHFK57mRHtvYPCY+49//FUZvipSEHCT6x2vhr/yThUA8fkklQejvFf4
V8TvJlKnOe9TnhVwHEKgEkwrWj6RlxNTn8OoZh3WEoM5NNXhsk/Xa5iUIZnM
-----END CERTIFICATE-----
`)
	AddVariableIfNotExist(n.adaptors, ctx, "certKey", `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCr7gMlkqgpGplY
vqjW90tSdwK9CxPdP4O5DlG3/biX+3bwfO5IYI1Th0d5hkWkG6Em+dvCmjpJiHjs
/mi6iiXffx/+bvMspVwwYgjt3SvXPT2QSG+URDH5HmarHULmqKMAq3DHmdDyT/EU
N3w4egjroEv5uiQn6QjeiW1aW4i5U6i3n57HHcRwAG7cRAU4dYD+wTmZ0Bwoa28M
7ORNAEf1FrapvC+AAwhbzWg2F3XXfN5TAkcqYhbFyH3PSbi8inLuOM4HSf4WuOLx
GHY6zrhE3fAyP8Kg55NZtSr5q/D857zdS5ibPTQaFqFsX6tlb+AgTirq+w+rbCS1
C6FvlNfPAgMBAAECggEAcupXOBnaQ/7/WA23lFceBTR+pBRvZoY5aMtlW3E+nHb7
fKpEKiQ+0gGtiFBy48mD4SVH+b5UDyokiWNSZLxJrCSwIcPOzZyJDd240iPuVaMd
Lv77dUJPlI75WI3qVXmJ2by9WOw6eHtuS3D6mlUW+UbfAT+lQvfDcdqxOJ/NtvBW
rBl6t7xvoSdMIOItSBSrlxkSf/iS6lCkfjl9/RnctbgpRKxNOa9ygoVMAtv8QdXz
OAtK91/EpaxA/DXtv9uVj94nuMHVI/KaxYsEKZxUhdv7WtBsqya/SFA0EbgRKkCb
/AcIlnu81buHiULSGN77OKaRLDY8YSyfutKjmMQO4QKBgQDCrdz7deQBBW/ylzLh
Fog/ThlhVxKhtaPGCn3iXceWPW0NsW7c6C6PaYSafsW2QoXYJIPVhcqpsaiHsoR/
D5/2j8s03LeAfJ2hHvuT644qP3tzWW+zdLNPa+r8MGCpuJ0dVHpXqEZsMqMMquaP
PIAi1717nKBUBueGgESR2HaiAwKBgQDiFbp2hSaNIDt0z7+XCjwP2hlUgKu95agJ
y9xFDsH51mXxdOo3js11T78HaMXpxDhVOwJdhfOrylGfJdLXsU2ScJf1psmsHI0E
Ms/80LXY7lsF+41SREkmaM1EKja2mOHJy8w/WcAJR9UII6peaWdPvdwce6+EdA1M
FiPaV18PRQKBgFlCHe+tPby6IXm4mTtaeV2NEVXv9jrubQiABveix7+6qiV9FLd0
POTEHGg5d6z5EyTmmQttLF5hdkPBNb6MC2ugwiqaiAgBGsSkz/QiLCxyrsdUWkKN
Xykq8bJu++LVIjQwZ2eJX8B79wX31X4W3ykagWYLifb8s2qXpNi5xqrBAoGBANFv
vX9j05WywUFx3pMl6QPqT8LHO9G7uNHyZwgWJr2xzlhCrWJwSIafLGDJv2XR3zmd
Y5CNsXB4COq8WDh1yF2bLjpjmrl07XFxPNl/4qWmEO5W3NDT70vElDrMPEklES0m
PPyGwO4X/BIBMCcrz4fAYvKD5dF2zIaqj+YJQEddAoGAIYgsbABIuLZrbI6RVG5J
P7zr8OvFQjhY21XipQ43ktOr6z4dCPeXqXo0CFoSgIM7Iy6STD5iTR1NVb7j7s0a
NCmNwBvrKDpDuDsgubbAzhYMJ+LZA8oRJgVGsJDS9jEh+skIWfV/nZFq/omZvl+u
WIxWVkOlrXbrK5IM5YogidY=
-----END PRIVATE KEY-----
`)
	return nil
}
