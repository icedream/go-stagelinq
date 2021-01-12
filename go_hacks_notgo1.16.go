//+build !go1.16

package stagelinq

func checkErrIsNetClosed(err error) bool {
	return err != nil && err.Error() == "use of closed network connection"
}
