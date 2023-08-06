package main

import (
	"os"
	"time"
)

type MockFileInfo struct{}

func (m *MockFileInfo) Name() string {
	return "MockFile"
}

func (m *MockFileInfo) Size() int64 {
	return 0
}

func (m *MockFileInfo) Mode() os.FileMode {
	return 0
}

func (m *MockFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (m *MockFileInfo) IsDir() bool {
	return false
}

func (m *MockFileInfo) Sys() interface{} {
	return nil
}
