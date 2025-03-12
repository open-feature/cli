package requests

import (
	"fmt"
	"io"
	"net/http"

	"github.com/open-feature/cli/internal/flagset"
)

func FetchFlags(flagSourceUrl string, authToken string) (flagset.Flagset, error) {
	flags := flagset.Flagset{}
	req, err := http.NewRequest("GET", flagSourceUrl, nil)
	if err != nil {
		return flags, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return flags, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return flags, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return flags, fmt.Errorf("Received error response from flag source: %s", string(body))
	}

	loadedFlags, err := flagset.LoadFromSourceFlags(body)
	if err != nil {
		return flags, err
	}
	flags.Flags = *loadedFlags


	return flags, nil
}
