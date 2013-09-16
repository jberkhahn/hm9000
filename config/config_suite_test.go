package config

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "testing"
)

func TestBootstrap(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Config Suite")
}
