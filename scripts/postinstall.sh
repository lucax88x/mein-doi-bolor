#!/bin/bash

echo "⚙️ Post Install"

# Correct versions (go -> golangci-lint)
declare -a version_list
version_list=(
  "1.24 2.2.1"
)

# Installed versions
go_version=$(go version | awk '{print $3}' | sed 's/go//' | awk -F. '{print $1"."$2}')
lint_version=$(golangci-lint --version | head -n1 | awk '{print $4}' | sed 's/v//')

# Expected golangci-lint version
expected_lint_version=""
for v in "${version_list[@]}"; do
    go_ver=$(echo "$v" | awk '{print $1}')
    lint_ver=$(echo "$v" | awk '{print $2}')
    if [[ "$go_ver" == "$go_version" ]]; then
        expected_lint_version="$lint_ver"
        break
    fi
done

echo "🔹 go            $go_version"
echo "🔸 golangci-lint $lint_version"
echo ""

# Check version compatibility
version_ge() {
    printf '%s\n%s\n' "$2" "$1" | sort -V | head -n1 | grep -q "$2"
}
if [[ -n "$expected_lint_version" ]]; then
    if version_ge "$lint_version" "$expected_lint_version"; then
        echo "✅  The installed go and golangci-lint versions are compatible"
    else
        echo "❌  The installed go and golangci-lint versions are NOT compatible!"
        echo "   ↳ go $go_version requires at least golangci-lint $expected_lint_version"
        echo ""
        echo "ℹ️ Use the command below to install the correct version of golangci-lint"
        echo "   ↳ curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v$expected_lint_version"
        exit 1
    fi
else
    echo "⚠️  The installed go version ($go_version) is not covered by the script, please update the postinstall.sh script"
    exit 2
fi
