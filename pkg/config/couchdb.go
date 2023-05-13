package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/couchbase/gocb/v2"
)

var defaultBucket *gocb.Bucket
var defaultCluster *gocb.Cluster
var defaultScope *gocb.Scope

type DBQueryParameters = map[string]interface{}

var collections []string = []string{"users", "credentials", "sessions", "authenticators"}

type DBOpError struct {
	Status  int
	Message string
}

func (err *DBOpError) Error() string {
	return err.Message
}

func (err *DBOpError) GetStatus() int {
	return err.Status
}

// NewCouchbaseConnection Creates a new Couchbase connection
func NewCouchbaseConnection() (*gocb.Cluster, error) {
	connectionString := fmt.Sprintf("couchbase://%s", os.Getenv("COUCHBASE_HOST"))

	clusterOpts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: os.Getenv("COUCHBASE_USERNAME"),
			Password: os.Getenv("COUCHBASE_PASSWORD"),
		},
	}

	cluster, err := gocb.Connect(connectionString, clusterOpts)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to connect to Couchbase cluster at %s: %s",
			connectionString, err.Error(),
		)
	}

	return cluster, nil
}

// InitDefaultBucket Creates a new Connection and opens the default bucket
func InitCouchbaseConnection() {

	connection, bucket, err := WaitForConnection()

	if err != nil {
		log.Panicln(err.Error())
	}

	// check or create scope
	scope := CheckOrCreateScope(bucket)
	// check or create collections
	CheckOrCreateCollections(bucket)
	// check or create indexes
	BuildIndexes(connection)

	// create secondary indexds
	// BuildSecondaryIndexes(connection)

	// save these connections:
	defaultCluster = connection
	defaultBucket = bucket
	defaultScope = scope

	log.Println("Initialized database connection.")
}

func GetDefaultBucket() *gocb.Bucket {
	return defaultBucket
}

func GetDefaultCluster() *gocb.Cluster {
	return defaultCluster
}

func GetDefaultScope() *gocb.Scope {
	return defaultScope
}

func ExecuteDBQuery(scope *gocb.Scope, queryString string, params *DBQueryParameters) (*gocb.QueryResult, *DBOpError) {

	results, err := scope.Query(queryString, &gocb.QueryOptions{
		NamedParameters: *params,
		Timeout:         120 * time.Second,
		ScanConsistency: gocb.QueryScanConsistencyRequestPlus,
	})

	if err != nil {
		return nil, &DBOpError{
			Status:  500,
			Message: err.Error(),
		}
	}

	return results, nil
}

func CheckOrCreateScope(bucket *gocb.Bucket) *gocb.Scope {
	collectionMgr := bucket.Collections()
	err := collectionMgr.CreateScope(os.Getenv("COUCHBASE_SCOPE"), &gocb.CreateScopeOptions{
		Timeout: 120 * time.Second,
	})

	if err != nil {
		if errors.Is(err, gocb.ErrScopeExists) {
			log.Printf("Scope %s already exists", os.Getenv("COUCHBASE_SCOPE"))
		} else {
			log.Fatalf("failed to create scope %s: %s", os.Getenv("COUCHBASE_SCOPE"), err.Error())
		}
	}

	log.Printf("Checked scope %s.", os.Getenv("COUCHBASE_SCOPE"))
	return bucket.Scope(os.Getenv("COUCHBASE_SCOPE"))
}

func CheckOrCreateCollections(bucket *gocb.Bucket) {
	collectionMgr := bucket.Collections()

	// create collections one by one:
	for _, collection := range collections {
		err := collectionMgr.CreateCollection(gocb.CollectionSpec{
			Name:      collection,
			ScopeName: os.Getenv("COUCHBASE_SCOPE"),
		}, &gocb.CreateCollectionOptions{
			Timeout: 120 * time.Second,
		})

		if err != nil {
			if errors.Is(err, gocb.ErrCollectionExists) {
				log.Printf("Collection %s already exists.", collection)
			} else {
				log.Fatalf("failed to create collection %s: %s", collection, err.Error())
			}
		}
	}
}

func WaitForConnection() (*gocb.Cluster, *gocb.Bucket, error) {
	MaxDBRetries, _ := strconv.Atoi(os.Getenv("MAX_DB_RETRIES"))

	for retry_count := 0; retry_count < MaxDBRetries; retry_count++ {
		log.Printf("%d. Trying to connect to couchbase DB", retry_count+1)
		connection, err := NewCouchbaseConnection()
		if err != nil {
			log.Println("failed to connect to Couchbase DB, retrying....")
			time.Sleep(120 * time.Second)
		}

		// wait for bucket or retry:
		bucket := connection.Bucket(os.Getenv("COUCHBASE_BUCKET"))
		err = bucket.WaitUntilReady(120*time.Second, nil)
		if err != nil {
			log.Println("failed to open couchbase bucket, retrying....")
		} else {
			return connection, bucket, nil
		}
	}

	return nil, nil, fmt.Errorf("could not connect to couchbase, tried %d times", MaxDBRetries)
}

func BuildIndexes(connection *gocb.Cluster) {

	log.Println("Building all indexes....")

	// create a primary index on these collections:
	for _, collection := range collections {
		for {
			path := fmt.Sprintf("%s`.`%s`.`%s", os.Getenv("COUCHBASE_BUCKET"), os.Getenv("COUCHBASE_SCOPE"), collection)
			name := fmt.Sprintf(
				"primary_%s_%s_%s", os.Getenv("COUCHBASE_BUCKET"),
				os.Getenv("COUCHBASE_SCOPE"), collection,
			)

			log.Printf("Building primary index for %s...", name)

			/*queryString := fmt.Sprintf(
				"CREATE PRIMARY INDEX %s ON %s",
				name, path,
			)

			queryResult, err := connection.Query(queryString, &gocb.QueryOptions{
				Timeout: 120 * time.Second,
			})*/

			qm := connection.QueryIndexes()
			err := qm.CreatePrimaryIndex(path, &gocb.CreatePrimaryQueryIndexOptions{
				Timeout:    120 * time.Second,
				CustomName: name,
			})

			if err != nil && !errors.Is(err, gocb.ErrIndexExists) {
				log.Printf("index creation error, retrying in 5 seconds: %s", err.Error())
				time.Sleep(120 * time.Second)
				continue
			}

			break
		}
	}

	log.Println("Created all indexes. Ready to serve queries.")
}

// func BuildSecondaryIndexes(connection *gocb.Cluster) {

// 	log.Println("Building all secondary indexes....")

// 	var queries []string = []string{
// 		"CREATE INDEX idx ON remaster(users.address)",
// 		"CREATE INDEX adv_orgId_deploymentInfo_contractAddress ON `default`:`remaster`.`default`.`agreements`(`orgId`,`deploymentInfo`.`contractAddress`)",
// 		"CREATE INDEX adv_templateId ON `default`:`remaster`.`default`.`contractData`(`templateId`)",
// 		"CREATE INDEX adv_templateId_abi ON `default`:`remaster`.`default`.`contractData`(`templateId`,`abi`)",
// 		"CREATE INDEX adv_createdBy ON `default`:`remaster`.`default`.`licenseTypes`(`createdBy`)",
// 		"CREATE INDEX adv_agreementId_seller_buyer ON `default`:`remaster`.`default`.`licenses`(`agreementId`,`seller`,`buyer`)",
// 		"CREATE INDEX adv_agreementId_seller_buyer ON `default`:`remaster`.`default`.`licenseAgreements`(`agreementId`,`seller`,`buyer`)",
// 		"CREATE INDEX adv_buyer_agreementId_seller ON `default`:`remaster`.`default`.`licenses`(`buyer`,`agreementId`,`seller`)",
// 		"CREATE INDEX adv_r2Address_tokenId_licenseGlobalId_signatureSeller_signatureBuyer_expiryTimeEpoch ON `default`:`remaster`.`default`.`licenses`(`r2Address`,`tokenId`,`licenseGlobalId`,`signatureSeller`,`signatureBuyer`,`expiryTimeEpoch`)",
// 		"CREATE INDEX adv_seller_agreementId_buyer ON `default`:`remaster`.`default`.`licenses`(`seller`,`agreementId`,`buyer`)",
// 		"CREATE INDEX adv_array_length_array_intersect_tagstags ON `default`:`remaster`.`default`.`licenseContracts`(array_length(array_intersect(`tags`, `tags`)))",
// 	}

// 	// create a secondary index on these collections:
// 	for _, query := range queries {

// 		log.Printf("Building secondary index for %s...", query)

// 		_, err := connection.Query(query, &gocb.QueryOptions{
// 			Timeout: 120 * time.Second,
// 		})

// 		if err != nil {
// 			log.Println("secondary index creation error", err.Error())
// 		}
// 	}

// 	log.Println("Created all secondary indexes. Ready to serve queries.")
// }
