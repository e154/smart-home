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
	"os"
	"path"
	"strings"

	"github.com/e154/smart-home/internal/common"
	. "github.com/e154/smart-home/internal/system/initial/assertions"
	"github.com/e154/smart-home/pkg/adaptors"
	m "github.com/e154/smart-home/pkg/models"
)

type MigrationImages struct {
	adaptors *adaptors.Adaptors
	dir      string
}

func NewMigrationImages(adaptors *adaptors.Adaptors, dir string) *MigrationImages {
	if dir == "" {
		dir = "./"
	}
	return &MigrationImages{
		adaptors: adaptors,
		dir:      dir,
	}
}

func (i *MigrationImages) Up(ctx context.Context) (err error) {

	imageList := map[string]*m.Image{
		"button_v1_off": {
			Image:    "30d2f4116a09fd14b49c266985db8109.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     2518,
			Name:     "button_v1_off.svg",
		},
		"button_v1_refresh": {
			Image:    "86486ca5d086aafd5724d61251b94bba.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3212,
			Name:     "button_v1_refresh.svg",
		},
		"lamp_v1_r": {
			Image:    "2d4a761241e24a77725287180656b466.svg",
			MimeType: "text/xml; charset=utf-8",
			Size:     2261,
			Name:     "lamp_v1_r.svg",
		},
		"socket_v1_b": {
			Image:    "bef910d70c56f38b22cea0c00d92d8cc.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     7326,
			Name:     "socket_v1_b.svg",
		},
		"button_v1_on": {
			Image:    "7c145f62dcaf8da2a9eb43f2b23ea2b1.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     2398,
			Name:     "button_v1_on.svg",
		},
		"socket_v1_def": {
			Image:    "4c28edf0700531731df43ed055ebf56d.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     7326,
			Name:     "socket_v1_def.svg",
		},
		"socket_v1_r": {
			Image:    "e91e461f7c9a800eed5a074101d3e5a5.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     7326,
			Name:     "socket_v1_r.svg",
		},
		"lamp_v1_def": {
			Image:    "91e93ee7e7734654083dee0a5cbe55e9.svg",
			MimeType: "text/xml; charset=utf-8",
			Size:     2266,
			Name:     "lamp_v1_def.svg",
		},
		"socket_v1_g": {
			Image:    "4819b36056dfa786f5856fa45e9a3151.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     7326,
			Name:     "socket_v1_g.svg",
		},
		"lamp_v1_y": {
			Image:    "c1c5ec4e75bb6ec33f5f8cfd87b0090e.svg",
			MimeType: "text/xml; charset=utf-8",
			Size:     2261,
			Name:     "lamp_v1_y.svg",
		},
		"socket_v2_b": {
			Image:    "c813ac54bb4dd6b99499d097eda67310.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3060,
			Name:     "socket_v2_b.svg",
		},
		"socket_v2_def": {
			Image:    "f0ea38f2b388dc2bb2566f6efc7731b0.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3060,
			Name:     "socket_v2_def.svg",
		},
		"socket_v2_g": {
			Image:    "fa6b42c81056069d03857cfbb2cf95eb.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3060,
			Name:     "socket_v2_g.svg",
		},
		"socket_v2_r": {
			Image:    "e565f191030491cfdc39ad728559c18f.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3060,
			Name:     "socket_v2_r.svg",
		},
		"socket_v3_b": {
			Image:    "297d56426098a53091fb8f91aabe3cd7.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     2718,
			Name:     "socket_v3_b.svg",
		},
		"socket_v3_def": {
			Image:    "becf0f8f635061c143acb4329f744615.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     2718,
			Name:     "socket_v3_def.svg",
		},
		"socket_v3_g": {
			Image:    "850bf4da00cb9de85e1442695230a127.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     2718,
			Name:     "socket_v3_g.svg",
		},
		"socket_v3_r": {
			Image:    "434514389e95cab6d684b978378055d5.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     2718,
			Name:     "socket_v3_r.svg",
		},
		"map-schematic-original": {
			Image:    "9384f1f6f9c2f4bf00fbc6debaae9b26.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     195108,
			Name:     "map-schematic-original.svg",
		},
		"temp_v1_r": {
			Image:    "688d2d752252de21c9d62a643c37ea40.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3980,
			Name:     "temp_v1_r.svg",
		},
		"temp_v1_y": {
			Image:    "8b2f46785aa3bdf7a6a487fc89a0f99e.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3980,
			Name:     "temp_v1_y.svg",
		},
		"temp_v1_def": {
			Image:    "655d491beafaefce2117cb2012dc674a.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3980,
			Name:     "temp_v1_def.svg",
		},
		"temp_v1_original": {
			Image:    "e8dee745788685f9f86e611cf5758cab.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     3770,
			Name:     "temp_v1_original.svg",
		},
		"fan_v1_r": {
			Image:    "eaf1c68959341c466fac68363f21cbbe.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     6238,
			Name:     "fan_v1_r.svg",
		},
		"fan_v1_y": {
			Image:    "33a5d5e7290e0f37a4c160cdbd0b5f23.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     6238,
			Name:     "fan_v1_y.svg",
		},
		"fan_v1_def": {
			Image:    "fd64ec639417d88e37b1c2cc167bcafc.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     6238,
			Name:     "fan_v1_def.svg",
		},
		"fan_v1_original": {
			Image:    "b4820c5939fe6b042888c922dfd1bada.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     5799,
			Name:     "fan_v1_original.svg",
		},
		"door_v1_closed": {
			Image:    "2e3e5c74775360e0274576ba6c83f044.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     509,
			Name:     "door_v1_closed.svg",
		},
		"door_v1_closed_r": {
			Image:    "221f451b426188a2df987163a2ab5715.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     508,
			Name:     "door_v1_closed_r.svg",
		},
		"door_v1_closed_def": {
			Image:    "dd2a735b71b2899869e36c54f140b3fa.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     509,
			Name:     "door_v1_closed_def.svg",
		},
		"door_v1_opened1": {
			Image:    "74cb4de3f70bb7a7e5d651ee6a23bffc.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     725,
			Name:     "door_v1_opened1.svg",
		},
		"door_v1_opened2": {
			Image:    "def90d2778eb6e4465f5808889e2a92c.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     862,
			Name:     "door_v1_opened2.svg",
		},
		"door_v1_opened3": {
			Image:    "fe7c9ecdbbdedc99ab16070da52251a4.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     1179,
			Name:     "door_v1_opened3.svg",
		},
		"md_v1_def": {
			Image:    "3f7482861152f6bf9de3940aa031e7bf.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     1686,
			Name:     "motion_detection_v1_def.svg",
		},
		"md_v1_original": {
			Image:    "763b23acd999eb268deec7320c4b0b88.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     1689,
			Name:     "motion_detection_v1_original.svg",
		},
		"md_v1_r": {
			Image:    "5d6d372b84cb75b80e6447cbc5cecb72.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     1685,
			Name:     "motion_detection_v1_r.svg",
		},
		"md_v1_o": {
			Image:    "0fcdb4e0857adeb71f699b18ce22e403.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     1711,
			Name:     "motion_detection_v1_o.svg",
		},
		"md_v1_y": {
			Image:    "55a011dcd81772b9752470471842d365.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     1689,
			Name:     "motion_detection_v1_y.svg",
		},
		"md_v2_def": {
			Image:    "fe4fad32a33ef2448debab3cf2dc4c6f.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     1388,
			Name:     "motion_detection_v2_def.svg",
		},
		"md_v2_original": {
			Image:    "73125c0ca60b84b647fe5f7bc2833432.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     1391,
			Name:     "motion_detection_v2_original.svg",
		},
		"md_v2_r": {
			Image:    "82799f02881efe45f2a5211ab357e2e2.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     1387,
			Name:     "motion_detection_v2_r.svg",
		},
		"md_v2_o": {
			Image:    "58b1695473e44277c0306a726a601aef.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     1391,
			Name:     "motion_detection_v2_o.svg",
		},
		"md_v2_y": {
			Image:    "aab3542bb028c61a4343bfc4e8c92daf.svg",
			MimeType: "text/html; charset=utf-8",
			Size:     1391,
			Name:     "motion_detection_v2_y.svg",
		},
	}

	var subDir string
	for _, image := range imageList {
		if _, err = i.adaptors.Image.GetByImageName(ctx, image.Image); err == nil {
			continue
		}

		image.Id, err = i.adaptors.Image.Add(ctx, image)
		So(err, ShouldBeNil)

		fullPath := common.GetFullPath(image.Image)
		to := path.Join(i.dir, fullPath, image.Image)

		//log.Infof("create dir %s", path.Join(i.dir, fullPath))
		_ = os.MkdirAll(path.Join(i.dir, fullPath), os.ModePerm)

		if exist := common.FileExist(to); !exist {

			switch {
			case strings.Contains(image.Name, "button"):
				subDir = "buttons"
			case strings.Contains(image.Name, "lamp"):
				subDir = "lamp"
			case strings.Contains(image.Name, "socket"):
				subDir = "socket"
			case strings.Contains(image.Name, "temp"):
				subDir = "temp"
			case strings.Contains(image.Name, "fan"):
				subDir = "fan"
			case strings.Contains(image.Name, "map"):
				subDir = "map"
			case strings.Contains(image.Name, "door"):
				subDir = "door"
			case strings.Contains(image.Name, "motion_detection"):
				subDir = "motion_detection"
			}

			from := path.Join("data", "icons", subDir, image.Name)
			//log.Infof("copy %s --> %s", from, to)
			common.CopyFile(from, to)
		}
	}

	return nil
}
