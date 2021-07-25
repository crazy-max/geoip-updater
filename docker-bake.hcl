// Go version
variable "GO_VERSION" {
  default = "1.16"
}

// GitHub reference as defined in GitHub Actions (eg. refs/head/master))
variable "GITHUB_REF" {
  default = ""
}

target "go-version" {
  args = {
    GO_VERSION = GO_VERSION
  }
}

// Special target: https://github.com/docker/metadata-action#bake-definition
target "docker-metadata-action" {
  tags = ["crazymax/geoip-updater:local"]
}

group "default" {
  targets = ["image-local"]
}

group "validate" {
  targets = ["lint", "vendor-validate"]
}

target "lint" {
  inherits = ["go-version"]
  dockerfile = "./hack/lint.Dockerfile"
  target = "lint"
  output = ["type=cacheonly"]
}

target "vendor-validate" {
  inherits = ["go-version"]
  dockerfile = "./hack/vendor.Dockerfile"
  target = "validate"
  output = ["type=cacheonly"]
}

target "vendor-update" {
  inherits = ["go-version"]
  dockerfile = "./hack/vendor.Dockerfile"
  target = "update"
  output = ["."]
}

target "test" {
  inherits = ["go-version"]
  dockerfile = "./hack/test.Dockerfile"
  target = "test-coverage"
  output = ["."]
}

target "docs" {
  dockerfile = "./hack/docs.Dockerfile"
  target = "release"
  output = ["./site"]
}

target "artifact" {
  args = {
    GIT_REF = GITHUB_REF
  }
  inherits = ["go-version"]
  target = "artifacts"
  output = ["./dist"]
}

target "artifact-all" {
  inherits = ["artifact"]
  platforms = [
    "darwin/amd64",
    "darwin/arm64",
    "freebsd/386",
    "freebsd/amd64",
    "linux/386",
    "linux/amd64",
    "linux/arm/v5",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
    "linux/ppc64le",
    "linux/riscv64",
    "linux/s390x",
    "windows/386",
    "windows/amd64"
  ]
}

target "image" {
  inherits = ["go-version", "docker-metadata-action"]
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/386",
    "linux/amd64",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
    "linux/ppc64le"
  ]
}
