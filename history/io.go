package history

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

func Load(filename string, initHistory bool) (*Accounts, error) {
	reader, err := os.Open(filename)
	if os.IsNotExist(err) && initHistory {
		return &Accounts{nil, &sync.Mutex{}}, nil
	}
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	result, err := LoadFrom(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %w", filename, err)
	}
	return result, err
}

func LoadFrom(reader io.Reader) (*Accounts, error) {
	decoder := yaml.NewDecoder(reader)
	var result []Account
	for {
		var account Account
		err := decoder.Decode(&account)
		if err == io.EOF {
			return &Accounts{sortAccounts(result), &sync.Mutex{}}, nil
		}
		if err != nil {
			return nil, err
		}
		sortHistory(account.History)
		result = append(result, account)
	}
}

func (a *Accounts) Save(filename string) error {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		bak := strings.Join([]string{filename, "bak"}, ".")
		copy(filename, bak)
	}

	writer, err := os.Create(filename)
	if err != nil {
		return err
	}
	encoder := yaml.NewEncoder(writer)

	a.lock.Lock()
	defer a.lock.Unlock()

	for _, account := range a.accounts {
		err := encoder.Encode(account)
		if err != nil {
			return fmt.Errorf("could not write entry to %s: %w", filename, err)
		}
	}
	return encoder.Close()
}

func copy(from, to string) {
	in, err := os.Open(from)
	if err != nil {
		fmt.Printf("Could not write backup: %v\n", err)
		return
	}
	defer in.Close()
	out, err := os.Create(to)
	if err != nil {
		fmt.Printf("Could not write backup: %v\n", err)
		return
	}
	_, err = io.Copy(out, in)
	if err != nil {
		fmt.Printf("Could not write backup: %v\n", err)
	}
}
