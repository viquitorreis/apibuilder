package main

	import (
		"github.com/tupatech/tupa"
	)

	func main() {
		server := tupa.NewAPIServer(":6969", nil)
		server.New()
	}
	