package main

var configuration = []config{
	{
		c: &policechecker{},
	},
	/*{
		checker: &centralbankchecker{},
	},
	/*{
		checker: &creditratingchecker{},
	},*/
}

func (p *policechecker) check(pc *prospectcompany) <-chan result {
	return nil
}
