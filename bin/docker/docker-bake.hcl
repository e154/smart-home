variable "GO_VERSION" {
  default = null
}

target "_common" {
  args = {
    GO_VERSION = GO_VERSION
    BUILDKIT_CONTEXT_KEEP_GIT_DIR = 1
  }
}

// Special target: https://github.com/docker/metadata-action#bake-definition
target "docker-metadata-action" {
  tags = ["goxx:local"]
}

group "default" {
  targets = ["image-local"]
}

target "base" {
  inherits = ["_common"]
  target = "base"
  output = ["type=cacheonly"]
}

target "image" {
  inherits = ["_common", "docker-metadata-action"]
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
}
