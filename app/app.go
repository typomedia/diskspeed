package app

type Application struct {
	Name        string
	Version     string
	Author      string
	Description string
}

var App = Application{
	Name:        "diskspeed",
	Version:     "0.1.0",
	Author:      "Brian Cunnie, Brendan Cunnie, Philipp Speck",
	Description: "Disk speed benchmarking tool",
}
