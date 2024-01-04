
---
title: "HDD"
linkTitle: "hdd"
date: 2021-10-20
description: >
  
---

{{< figure src="/smart-home/hdd1.png" width="300" >}}

&nbsp;

&nbsp;

In the **Smart Home** system, there is a "hdd" plugin that provides the display of hard disk parameters. This plugin allows you to obtain information about various characteristics of the hard disk. Here are some of these parameters:

1. `path`: The `path` parameter contains the path to the mount point of the hard disk.

2. `fstype`: The `fstype` parameter indicates the file system type used on the hard disk.

3. `total`: The `total` parameter displays the total size of the hard disk in bytes.

4. `free`: The `free` parameter shows the amount of free space on the hard disk in bytes.

5. `used`: The `used` parameter indicates the used space on the hard disk in bytes.

6. `used_percent`: The `used_percent` parameter shows the percentage of space used on the hard disk.

7. `inodes_total`: The `inodes_total` parameter displays the total number of inodes on the hard disk.

8. `inodes_used`: The `inodes_used` parameter indicates the number of used inodes on the hard disk.

9. `inodes_free`: The `inodes_free` parameter shows the number of free inodes on the hard disk.

10. `inodes_used_percent`: The `inodes_used_percent` parameter indicates the percentage of inodes used on the hard disk.

Additionally, the "hdd" plugin has a `mount_point` settings option that allows you to specify the mount point to display parameters for a specific hard disk.

Here's an example of using the "hdd" plugin to retrieve hard disk parameters:

```javascript
const hddParams = EntityGetAttributes('hdd.hdd1')
const hddSettings = EntityGetSettings('hdd.hdd1')
console.log(hddSettings.mount_point);
console.log(hddParams.path);
console.log(hddParams.fstype);
console.log(hddParams.total);
console.log(hddParams.free);
console.log(hddParams.used);
console.log(hddParams.used_percent);
console.log(hddParams.inodes_total);
console.log(hddParams.inodes_used);
console.log(hddParams.inodes_free);
console.log(hddParams.inodes_used_percent);
```

You can use these parameters to display information about the hard disk in your **Smart Home** project.
