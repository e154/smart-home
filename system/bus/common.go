// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package bus

// TopicMatch 返回topic和topic filter是否
//
// TopicMatch returns whether the topic and topic filter is matched.
func TopicMatch(topic []byte, topicFilter []byte) bool {
	var spos int
	var tpos int
	var sublen int
	var topiclen int
	var multilevelWildcard bool //是否是多层的通配符
	sublen = len(topicFilter)
	topiclen = len(topic)
	if sublen == 0 || topiclen == 0 {
		return false
	}
	if (topicFilter[0] == '$' && topic[0] != '$') || (topic[0] == '$' && topicFilter[0] != '$') {
		return false
	}
	for {
		//e.g. ([]byte("foo/bar"),[]byte("foo/+/#")
		if spos < sublen && tpos <= topiclen {
			if tpos != topiclen && topicFilter[spos] == topic[tpos] { // sublen是订阅 topiclen是发布,首字母匹配
				if tpos == topiclen-1 { //遍历到topic的最后一个字节
					/* Check for e.g. foo matching foo/# */
					if spos == sublen-3 && topicFilter[spos+1] == '/' && topicFilter[spos+2] == '#' {
						return true
					}
				}
				spos++
				tpos++
				if spos == sublen && tpos == topiclen { //长度相等，内容相同，匹配
					return true
				} else if tpos == topiclen && spos == sublen-1 && topicFilter[spos] == '+' {
					//订阅topic比发布topic多一个字节，并且多出来的内容是+ ,比如: sub: foo/+ ,topic: foo/
					if spos > 0 && topicFilter[spos-1] != '/' {
						return false
					}
					return true
				}
			} else {
				if topicFilter[spos] == '+' { //sub 和 topic 内容不匹配了
					spos++
					for { //找到topic的下一个主题分割符
						if tpos < topiclen && topic[tpos] != '/' {
							tpos++
						} else {
							break
						}
					}
					if tpos == topiclen && spos == sublen { //都遍历完了,返回true
						return true
					}
				} else if topicFilter[spos] == '#' {
					return true
				} else {
					/* Check for e.g. foo/bar matching foo/+/# */
					if spos > 0 && spos+2 == sublen && tpos == topiclen && topicFilter[spos-1] == '+' && topicFilter[spos] == '/' && topicFilter[spos+1] == '#' {
						return true
					}
					return false
				}
			}
		} else {
			break
		}
	}
	if !multilevelWildcard && (tpos < topiclen || spos < sublen) {
		return false
	}
	return false
}
