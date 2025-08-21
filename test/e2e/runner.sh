#!/bin/bash

# Install ginkgo
go install github.com/onsi/ginkgo/v2/ginkgo@latest

# Run Test suites
ginkgo -r test/e2e/tests

echo "=====Successfully ran all the test suites======"
