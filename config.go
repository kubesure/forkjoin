package main

//Config defines configuration for routing
type config struct {
	checker Checker
}

var configuration = []config{
	{
		checker: &policechecker{},
	},
	/*{
		checker: &centralbankchecker{},
	},
	{
		checker: &creditratingchecker{},
	},*/
}
