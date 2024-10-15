package version

import _ "embed"

// Version - версия приложения из git, заполняется при компиляции программы
// из файла version.txt
// для обновления версии запустите
// make_version.sh
//
//go:embed version.txt
var Version string
