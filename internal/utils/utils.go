package utils

import "fmt"

func FormatSize(size int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	var i int
	floatSize := float64(size)

	for i = 0; floatSize >= 1024 && i < len(units)-1; i++ {
		floatSize /= 1024
	}

	return fmt.Sprintf("%.2f %s", floatSize, units[i])

}
