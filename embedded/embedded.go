package embedded

import "embed"

//go:embed savedata
var ExampleFS embed.FS

// Имя папки внутри embed (важно для копирования)
const ExampleDirName = "savedata"