package main

import "msa-nd-box-go/processing"

func main() {
	processing.InitManager(2,`C:\projects\msa-nd-box-go\bin\processing_worker_runner_go.exe`)
}

