package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/shenzhencenter/google-ads-pb/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func listCompaigns(ctx context.Context, asc services.GoogleAdsServiceClient) {

	res, err := asc.SearchStream(ctx, &services.SearchGoogleAdsStreamRequest{
		CustomerId: os.Getenv("CUSTOMER_ID"),
		Query:      "SELECT campaign.id FROM campaign",
	})

	if err != nil {
		log.Fatal(err)
	}

	for {
		response, err := res.Recv()

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Fatal(err)
		}

		for _, r := range response.Results {
			fmt.Println(r)
		}
	}

}

func main() {

	ctx := context.Background()

	headers := metadata.Pairs(
		"authorization", "Bearer "+os.Getenv("ACCESS_TOKEN"),
		"developer-token", os.Getenv("DEVELOPER_TOKEN"),
		"login-customer-id", os.Getenv("CUSTOMER_ID"),
	)
	ctx = metadata.NewOutgoingContext(ctx, headers)

	cred := grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	conn, err := grpc.Dial("googleads.googleapis.com:443", cred)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	listCompaigns(ctx, services.NewGoogleAdsServiceClient(conn))

	// customerServiceClient := services.NewCustomerServiceClient(conn)
	// accessibleCustomers, err := customerServiceClient.ListAccessibleCustomers(
	// 	ctx,
	// 	&services.ListAccessibleCustomersRequest{},
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// for _, customer := range accessibleCustomers.ResourceNames {
	// 	fmt.Println("ResourceName: " + customer)
	// }
}
