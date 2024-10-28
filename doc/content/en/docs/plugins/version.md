---
title: "Version"
linkTitle: "version"
date: 2021-10-20
description: >
  
---

The "version" plugin in the system provides information about the version and build of the project. It displays the
following parameters:

1. `version`: The project version.
2. `revision`: The revision identifier (commit) in the version control system.
3. `revision_url`: The URL where you can view the details of the revision.
4. `generated`: The date and time of the project build generation.
5. `developers`: The list of developers involved in the project.
6. `build_num`: The project's build number.
7. `docker_image`: The name of the Docker image associated with the project.
8. `go_version`: The version of the Go programming language used for project development.

The "version" plugin provides access to this information, which is useful for tracking versions and monitoring the
development process.

Here's an example of using the "version" plugin:

```javascript
const versionInfo = EntityGetAttributes('version.version')
console.log("Project version:", versionInfo.version);
console.log("Revision:", versionInfo.revision);
console.log("Revision URL:", versionInfo.revision_url);
console.log("Generated:", versionInfo.generated);
console.log("Developers:", versionInfo.developers);
console.log("Build number:", versionInfo.build_num);
console.log("Docker image:", versionInfo.docker_image);
console.log("Go version:", versionInfo.go_version);
```

This plugin allows you to obtain up-to-date information about the project version and use it for various purposes such
as displaying the version in the user interface or logging.
