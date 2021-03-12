/*
Package network contains functionality for working GCLoud networks API resources

Example to List Network

	allPages, err := availablenetworks.List(networkClient).AllPages()
	if err != nil {
		panic(err)
	}

	allNetworks, err := availablenetworks.ExtractNetworks(allPages)
	if err != nil {
		panic(err)
	}

	for _, network := range allNetworks {
		fmt.Printf("%+v", network)
	}

*/
package availablenetworks
