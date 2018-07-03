package server

//func newClient(url string) (*elastic.Client, error) {
//	ctx := context.Background()
//	var err error
//
//	client, err := elastic.NewClient()
//	if err != nil {
//		return nil, err
//	}
//
//	info, code, err := client.Ping(url).Do(ctx)
//	if err != nil {
//		return nil, err
//	}
//	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
//
//	return client, nil
//}
