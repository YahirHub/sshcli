package paths

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// ToLocal convierte rutas estilo Unix/MSYS2 a rutas nativas de Windows.
func ToLocal(p string) string {
	if p == "" {
		return ""
	}
	if runtime.GOOS == "windows" {
		if len(p) >= 3 && p[0] == '/' && p[2] == '/' {
			drive := strings.ToUpper(string(p[1]))
			p = drive + ":" + p[2:]
		}
		return filepath.FromSlash(p)
	}
	return filepath.Clean(p)
}

// ToRemote limpia mutilaciones de Git Bash y asegura formato Linux absoluto.
func ToRemote(p string) string {
	if p == "" || p == "/" || p == "." {
		return "/"
	}

	// 1. Normalizar slashes
	res := strings.ReplaceAll(p, "\\", "/")

	// 2. Si estamos en Windows, limpiar interferencias
	if runtime.GOOS == "windows" {
		// Quitar letra de unidad (C:/...)
		if len(res) >= 2 && res[1] == ':' {
			res = res[2:]
		}

		// Lista de prefijos de mutilación conocidos
		prefixes := []string{
			"/Program Files/Git",
			"/Program Files (x86)/Git",
			"/usr/bin",
			"/mingw64",
			"/mingw32",
			"/usr",
		}

		for _, pre := range prefixes {
			if strings.HasPrefix(res, pre) {
				res = strings.TrimPrefix(res, pre)
				break
			}
		}
	}

	// 3. Normalizar con path (Linux)
	res = path.Clean(res)

	// 4. Si quedó vacío tras limpiar, es la raíz
	if res == "" || res == "." {
		return "/"
	}

	// 5. Garantizar que sea absoluta si no es relativa explícita
	if !strings.HasPrefix(res, "/") && !strings.HasPrefix(res, ".") {
		res = "/" + res
	}

	return res
}