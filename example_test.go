// Hi there. Below is some toy code that I've written to improve my
// understanding of how postgres stores and accesses data. You might enjoy
// reading it along with the other code in this repostiory. Or you can start
// with the README to get a more verbose overview of this project :).

package gopg

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	pq "github.com/lib/pq"
)

func TestExample(t *testing.T) {
	// Connect to our database in order to create some sample data and to lookup
	// some internal identifiers that can't be resolved yet by this library.
	dbName := os.Getenv("PGDB")
	opts := []string{
		"sslmode=disable",
		"port=" + os.Getenv("PGPORT"),
		"database=" + dbName,
	}
	sqlDB, err := sql.Open("postgres", strings.Join(opts, " "))
	if err != nil {
		t.Fatal(err)
	}

	// Create a very simple table with an int and text column.
	tableName := "foo"
	sql := `
DROP TABLE IF EXISTS foo;
CREATE table ` + pq.QuoteIdentifier(tableName) + ` AS
SELECT i, 'item-'||i
FROM generate_series(1, 3) i;
`
	if _, err := sqlDB.Exec(sql); err != nil {
		t.Fatal(err)
	}

	// Figure out the oid of our database.
	dbOIDRow := sqlDB.QueryRow("SELECT oid FROM pg_database WHERE datname = $1", dbName)
	var dbOID uint32
	if err := dbOIDRow.Scan(&dbOID); err != nil {
		t.Fatal(err)
	}

	// Get the relfilenode of our table.
	tableRelFileNodeRow := sqlDB.QueryRow("SELECT relfilenode FROM pg_class WHERE relname = $1", tableName)
	var tableRelFileNode uint32
	if err := tableRelFileNodeRow.Scan(&tableRelFileNode); err != nil {
		t.Fatal(err)
	}

	// Open the cluster data directory.
	config := ClusterConfig{Path: os.Getenv("PGDATA")}
	c, err := config.Open()
	if err != nil {
		t.Fatal(err)
	}

	// Open the database by oid.
	db, err := c.DatabaseByOID(dbOID)
	if err != nil {
		t.Fatal(err)
	}

	// Open the table by relfilenode.
	table, err := db.TableByRelFileNode(tableRelFileNode)
	if err != nil {
		t.Fatal(err)
	}

	// Create an iterator to iterate through the pages.
	pageIter, err := table.Pages()
	if err != nil {
		t.Fatal(err)
	}
	defer pageIter.Close()

	// Keep track of the page block number (i) so we can output it. Page block
	// numbers are counted starting from 0 by convention in postgres.
	for i := 0; ; i++ {
		// Fetch a new page, or break the loop once we've read all pages (EOF).
		page, err := pageIter.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		// Extract the header data from the page. You might find the API a bit
		// overkill at this point, and perhaps I should have simply converted the
		// whole page into a Go tree data structure to keep things simple. But I
		// decided to keep the API verbose to highlight the individual steps that
		// are taken for educational purpose.
		pageHeader, err := page.Header()
		if err != nil {
			t.Fatal(err)
		}

		// Print information about our page and its header.
		fmt.Printf("page %d:\n", i)
		fmt.Printf("  header: %+v:\n", pageHeader)

		// Get an iterator for our item identifiers that point to the tuples on the
		// page.
		itemIdentifierIter := page.ItemIdentifiers()

		// Keep track of tuple offsets (j). For some reason the convention is to
		// start counting from 1 for those in postgres.
		for j := 1; ; j++ {
			// Fetch the next item identifier, or break out of the loop once we've
			// seen them all (EOF).
			itemIdentifier, err := itemIdentifierIter.Next()
			if err == io.EOF {
				break
			} else if err != nil {
				t.Fatal(err)
			}

			// Follow the item identifiers pointer in the heap to the tuple it
			// belongs to.
			tuple, err := itemIdentifier.Tuple()
			if err != nil {
				t.Fatal(err)
			}

			// Print information about our tuple.
			fmt.Printf("  tuple %d\n", j)
			fmt.Printf("    item identifier: %+v\n", itemIdentifier)
			fmt.Printf("    tuple header: %+v\n", tuple)
			fmt.Printf("    data: %x\n", tuple.Data())

			// This was fun to implement, and has given me a better intuition for the
			// kind of stuff that is going on inside of postgres when performing a
			// `SELECT * FROM foo` query on a table. Obviously there is a lot more to
			// it, and perhaps I'll end up implementing some additional features :).
		}
	}
}
