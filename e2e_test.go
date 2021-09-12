package dalgo_firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/strongo/dalgo"
	"github.com/strongo/validation"
	"log"
	"os"
	"strings"
	"testing"
)

const E2ETestKind = "E2ETest"

type TestData struct {
	StringProp  string
	IntegerProp int
}

func (v TestData) Validate() error {
	if strings.TrimSpace(v.StringProp) == "" {
		return validation.NewErrRecordIsMissingRequiredField("StringProp")
	}
	if v.IntegerProp < 0 {
		return validation.NewErrBadRecordFieldValue("IntegerProp", fmt.Sprintf("should be > 0, got: %v", v.IntegerProp))
	}
	return nil
}

func TestEndToEnd(t *testing.T) {

	firestoreProjectID := os.Getenv("DALGO_E2E_PROJECT_ID")

	if firestoreProjectID == "" {
		firestoreProjectID = "DALGO_E2E"
		//t.Fatalf("Environment variable DALGO_E2E_PROJECT_ID is not set")
	}
	log.Println("Firestore Project ID:", firestoreProjectID)
	log.Println("ENV: GOOGLE_APPLICATION_CREDENTIALS:", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, firestoreProjectID)
	if err != nil {
		t.Fatalf("failed to create Firestore client: %v", err)
	}
	db := NewDatabase(client)
	key := dalgo.NewRecordKey(dalgo.RecordRef{Kind: E2ETestKind, ID: "r1"})
	t.Run("get", func(t *testing.T) {
		data := TestData{
			StringProp:  "str1",
			IntegerProp: 1,
		}
		record := dalgo.NewRecord(key, &data)
		if err = db.Get(ctx, record); err != nil {
			if dalgo.IsNotFound(err) {
				if err = db.Delete(ctx, record.Key()); err != nil {
					t.Fatalf("failed to delete: %v", err)
				}
			} else {
				t.Errorf("unexpected error: %v", err)
			}
		}
	})
	t.Run("create", func(t *testing.T) {
		t.Run("with_predefined_id", func(t *testing.T) {
			data := TestData{
				StringProp:  "str1",
				IntegerProp: 1,
			}
			record := dalgo.NewRecord(key, &data)
			err := db.Insert(ctx, record, dalgo.NewInsertOptions())
			if err != nil {
				t.Errorf("got unexpected error: %v", err)
			}
		})
	})
	t.Run("delete", func(t *testing.T) {
		if err := db.Delete(ctx, key); err != nil {
			t.Errorf("Failed to delete: %v", err)
		}
	})
}