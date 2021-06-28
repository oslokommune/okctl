package fetch_test

/*
func TestVerifer(t *testing.T) {
	testCases := []struct {
		name        string
		stager      *fetch.Stager
		expect      interface{}
		expectError bool
	}{
		{
			name: "Nop pipeline",
			stager: &fetch.Stager{
				Storage:      fetch.NewEphemeralStorage(),
				Fetcher:      fetch.NewStaticFetcher([]byte("hi there")),
				Verifier:     fetch.NewNoopVerifier(),
				Decompressor: fetch.NewNoopDecompressor(),
			},
			expect: "hi there",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			m := make(map[digest.Type]string)

			fetch.NewVerifier()
		})
	}
}
*/
