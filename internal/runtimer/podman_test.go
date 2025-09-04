package runtimer

// type bindingTest struct {
// 	artifactDirPath string
// 	imageCacheDir   string
// 	sock            string
// 	tempDirPath     string
// 	runRoot         string
// 	crioRoot        string
// 	conn            context.Context
// }
//
// func newBindingTest() *bindingTest {
// 	tmpPath, _ := createTempDirInTempDir()
// 	b := bindingTest{
// 		crioRoot:        filepath.Join(tmpPath, "crio"),
// 		runRoot:         filepath.Join(tmpPath, "run"),
// 		artifactDirPath: "",
// 		imageCacheDir:   "",
// 		sock:            fmt.Sprintf("unix://%s", filepath.Join(tmpPath, "api.sock")),
// 		tempDirPath:     tmpPath,
// 	}
// 	return &b
// }
//
// func TestConnector(t *testing.T) {
// 	_, err := NewConnector("/foo/socket")
// 	if err != nil {
// 		t.Fatal("Socket path could not be started")
// 	}
// }
