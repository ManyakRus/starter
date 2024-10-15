echo $(git describe --tags $(git rev-parse HEAD)) $(git show --no-patch --format=%ci) >./pkg/version/version.txt
