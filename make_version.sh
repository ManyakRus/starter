# версия приложения из git заполняется в файл: version.txt
# образец:
# v1.0.4-23-gf3bbaf4 2024-10-14 14:43:55 +0300
# v1.0.61 2024-10-14 14:25:20 +0300
# git commit TAG + git commit HASH + git commit date and time

echo $(git describe --always --tags $(git rev-parse HEAD)) $(git show --no-patch --format=%ci) >./pkg/version/version.txt
