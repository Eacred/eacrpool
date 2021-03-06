package pool

import (
	"os"
	"testing"

	bolt "github.com/coreos/bbolt"
	"github.com/Eacred/eacrd/chaincfg"
	"github.com/Eacred/eacrd/dcrutil"
)

var (
	// testDB represents the database used in testing.
	testDB = "testdb"
	// Account X address.
	xAddr = "SsWKp7wtdTZYabYFYSc9cnxhwFEjA5g4pFc"
	// Account X id.
	xID = ""
	// Account Y address.
	yAddr = "Ssp7J7TUmi5iPhoQnWYNGQbeGhu6V3otJcS"
	// Account Y id.
	yID = ""
	// Pool fee address.
	poolFeeAddrs, _ = dcrutil.DecodeAddress(
		"SsnbEmxCVXskgTHXvf3rEa17NA39qQuGHwQ",
		chaincfg.SimNetParams())
)

// setupDB initializes the pool database.
func setupDB() (*bolt.DB, error) {
	os.Remove(testDB)
	db, err := openDB(testDB)
	if err != nil {
		return nil, err
	}
	err = createBuckets(db)
	if err != nil {
		return nil, err
	}
	err = upgradeDB(db)
	if err != nil {
		return nil, err
	}

	xID, err = AccountID(xAddr, chaincfg.SimNetParams())
	if err != nil {
		return nil, err
	}
	yID, err = AccountID(yAddr, chaincfg.SimNetParams())
	if err != nil {
		return nil, err
	}
	_, err = persistAccount(db, xAddr, chaincfg.SimNetParams())
	if err != nil {
		return nil, err
	}
	_, err = persistAccount(db, yAddr, chaincfg.SimNetParams())
	if err != nil {
		return nil, err
	}
	return db, err
}

// teardownDB closes the connection to the db and deletes the db file.
func teardownDB(db *bolt.DB, dbPath string) error {
	err := db.Close()
	if err != nil {
		return err
	}
	err = os.Remove(dbPath)
	if err != nil {
		return err
	}
	return nil
}

// TestPool runs all pool related tests.
func TestPool(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	td := func() {
		err = teardownDB(db, testDB)
		if err != nil {
			t.Fatalf("teardown error: %v", err)
		}
	}

	defer td()

	testInitDB(t)
	testDatabase(t, db)
	testAcceptedWork(t, db)
	testAccount(t, db)
	testJob(t, db)
	testShares(t, db)
	testLimiter(t)
	testSharePercentages(t)
	testCalculatePoolTarget(t)
	testGeneratePaymentDetails(t, db)
	testArchivedPaymentsFiltering(t, db)
	testAccountPayments(t, db)
	testDifficulty(t)
	testEndpoint(t, db)
	testClient(t, db)
	testPaymentMgr(t, db)
	testChainState(t, db)
	testHub(t, db)
}
