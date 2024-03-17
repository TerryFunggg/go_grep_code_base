package grep

import (
    "testing"
)


func TestGetFilePath(t *testing.T) {
    targetPath := "~/code"
    path := GetFilePath(targetPath)

    if path != targetPath {
        t.Fatalf(`GetFilePath("%q") want match for %q`, path, targetPath)
    }
}
