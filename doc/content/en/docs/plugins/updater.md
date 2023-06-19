
---
title: "Updater"
linkTitle: "updater"
date: 2021-10-20
description: >
  
---

The "updated" plugin allows the system to check for updates and provides information about the latest version. This is useful for users who want to stay up to date with the latest updates and implement new features or bug fixes. It displays the following parameters:

1. `latest_version`: The number of the latest available version of the **Smart Home** system.
2. `latest_version_time`: The date and time when the latest version was released.
3. `latest_download_url`: The URL where the latest version of the system can be downloaded.
4. `last_check`: The date and time of the last update check.

Here's an example of using the "updated" plugin:

```javascript
const entity = entityManager.getEntity('updater.updater')
const updateInfo = entity.getAttributes()

console.log("Latest version:", updateInfo.latest_version);
console.log("Latest version time:", updateInfo.latest_version_time);
console.log("Latest download URL:", updateInfo.latest_download_url);
console.log("Last check:", updateInfo.last_check);
```

This code retrieves the attributes of the "updated" entity and logs the latest version, release time, download URL, and the last update check time.
