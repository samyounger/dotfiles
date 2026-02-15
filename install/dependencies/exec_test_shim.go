package dependencies

var (
	origExecLookPath = execLookPath
	origExecCommand  = execCommand
)

func resetExecMocks() {
	execLookPath = origExecLookPath
	execCommand = origExecCommand
}
