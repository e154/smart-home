variable "RELEASE_VERSION" {
  default = ""
}

variable "GO_BUILD_LDFLAGS" {
  default = ""
}

target "goxx" {
  context = "./bin/docker"
}

target "_commons" {
  contexts = {
    goxx = "target:goxx"
  }
  args = {
    RELEASE_VERSION: "${RELEASE_VERSION}"
    GO_BUILD_LDFLAGS: "${GO_BUILD_LDFLAGS}"
  }
}

group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["_commons"]

  labels = {
    "org.opencontainers.image.title": "Smart home"
    "org.opencontainers.image.authors": "Filippov Alex <af@e154.ru>"
    "org.opencontainers.image.description": "Managing iot devices has become easier thanks to the smart-home automation software platform for managing iot devices with a graphical interface and javascript's scripts."
    "org.opencontainers.image.licenses": "GPL-3.0-only"
    "org.opencontainers.image.version": "${RELEASE_VERSION}"
  }

  tags = ["e154/smart-home-server-test:${RELEASE_VERSION}", "e154/smart-home-server-test:latest"]
}

target "image-local" {
  dockerfile = "Dockerfile"
  inherits = ["image"]
  output = ["type=docker"]
}

target "image-all" {
  dockerfile = "Dockerfile"
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm64",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/ppc64le",
    "linux/riscv64",
    "linux/s390x"
  ]
}

target "image-linux-arm64" {
  dockerfile = "Dockerfile"
  inherits = ["image"]
  platforms = [
    "linux/arm64"
  ]
}

target "artifact" {
  inherits = ["_commons"]
  target = "artifact"
  output = ["./dist"]
}

target "artifact-all" {
  dockerfile = "Dockerfile"
  inherits = ["artifact"]
  platforms = [
    "linux/amd64",
    "linux/arm64",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/ppc64le",
    "linux/riscv64",
    "linux/s390x",
    "windows/amd64",
    "windows/arm64"
  ]
}

target "artifact-darwin-arm64" {
  dockerfile = "Dockerfile"
  inherits = ["artifact"]
  platforms = [
   "darwin/arm64"
  ]
}

target "artifact-public" {
  dockerfile = "./bin/docker/Dockerfile.public"
  output = ["./build/public"]
}
