package main

import (
	"log/slog"
	"os"
	"sync"

	"github.com/google/uuid"

	"github.com/NordSecurity-Interviews/BE-PatrykPasterny/client/client"
)

var clientsToSetup = []client.Customer{
	{
		ClientUUID: uuid.MustParse("8212d8ba-74d1-49af-8a84-6d6c392ec71c"),
		Longitude:  73.4,
		Latitude:   45.4,
		City:       "Ottawa",
		Height:     50000.0,
		Width:      55000.0,
	},
	{
		ClientUUID: uuid.MustParse("897737a8-77f1-4f53-8a51-6f9edaee6ed9"),
		Longitude:  73.5,
		Latitude:   45.5,
		City:       "Ottawa",
		Height:     50000.0,
		Width:      55000.0,
	},
	{
		ClientUUID: uuid.MustParse("4443822a-530c-43b9-a1ed-80cdf47a3cb3"),
		Longitude:  65.5,
		Latitude:   30.5,
		City:       "Montreal",
		Height:     50000.0,
		Width:      55000.0,
	},
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	customerClient := client.NewCustomerClient(logger)

	wg := sync.WaitGroup{}

	logger.Info("Setting up clients.")

	for i := range clientsToSetup {
		go func(customer client.Customer, group *sync.WaitGroup) {
			defer group.Done()

			customerClient.SimulateUser(&customer)
		}(clientsToSetup[i], &wg)
		wg.Add(1)
	}

	wg.Wait()

	logger.Info("Clients finished their rides.")
}
