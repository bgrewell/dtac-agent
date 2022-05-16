package remote

// Remote module is designed to show how you could have locally compiled logic that is set up to connect to a remote
// process (locally on the same system or over the network).
// This example shows how you can retrieve a message stored in the remote process. To see this example in action
// compile the 'main.go' file in the ./process folder using the command `go build -o bin/remote_process modules/remote/process/main.go`
// then execute the process with `bin/remote_process sUpEr_S3crEt_MesSag3` once this is done you can go to
// http://<ip>:8080/ext/remote/ and see the message
