### Augmenting MMDB files with your own data using Go

[MaxMind DB](https://github.com/maxmind/MaxMind-DB/blob/master/MaxMind-DB-spec.md) (or MMDB) files facilitate the storage and retrieval of data in connection with IPs and IP ranges, making queries for such data very fast and easy to perform. While MMDB files are usable on a variety of platforms and in a number of different programming languages, this article will focus on building MMDB files using the [Go programming language](https://golang.org/).

MaxMind offers several prebuilt MMDB files, like the free [GeoLite2 Country](https://dev.maxmind.com/geoip/geoip2/geolite2/) MMDB file. For many situations these MMDB files are useful enough as is. If, however, you have your own data associated with IP ranges, you can create hybrid MMDB files, augmenting existing MMDB contents with your own data. In this article, we're going to add details about a fictional company's IP ranges to the GeoLite2 Country MMDB file. We'll be building a new MMDB file, one that contains both MaxMind's and our fictional company's data.

If you don't need any of the MaxMind data, but you still want to create a fast, easy-to-query database keyed on IPs and IP ranges, you can consult this example code showing [how to create an MMDB file from scratch](https://github.com/maxmind/mmdbwriter/blob/master/examples/asn-writer/main.go).

### Prerequisites

- you must have [`git`](https://git-scm.com/downloads) installed in order to clone the code and install the dependencies, and it must be in your `$PATH`
- [Go 1.14](https://golang.org/dl/) or later must be installed, and `go` must be in your `$PATH`
- the [`mmdbinspect`](https://github.com/maxmind/mmdbinspect) tool must be installed, and be in your `$PATH`
- a copy of the [GeoLite2 Country](https://dev.maxmind.com/geoip/geoip2/geolite2/) database must be in your working directory
- your working directory (which can be located under any parent directory) must be named `mmdb-from-go-blogpost` (if you clone the code using the instructions below, this directory will be created in the appropriate place)
- a basic understanding of [Go](https://gobyexample.com/) and of [IP addresses](https://en.wikipedia.org/wiki/IP_address) and [CIDR notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing#CIDR_notation) will be helpful, but allowances have been made for the intrepid explorer for whom these concepts are novel!

### AcmeCorp's data

For the purposes of this tutorial, I have mocked up some data for a fictional company, AcmeCorp. This method can be adapted for your own real data, as long as that data maps to IPs or IP ranges.

AcmeCorp has three departments:
- SRE, whose IPs come from the 56.0.0.0/16 range,
- Development, whose IPs come from the 56.1.0.0/16 range, and
- Management, whose IPs come from the 56.2.0.0/16 range.

Members of the SRE department have access to all three of AcmeCorp's environments, `development`, `staging`, and `production`. Members of the Development and Management departments have access to the `development` and `staging` environments (but not to `production`).

The GeoLite2 Country MMDB file has a map, or no record, associated with every IP range.

For each of the AcmeCorp ranges, we're going to make sure a map record does get returned when we query for the ranges or any subranges of them, and that the map will contain `AcmeCorp.Environments` and `AcmeCorp.DeptName` keys, as well as any keys that were contained in the original GeoLite2 Country record, if one existed. More on this later.

### The steps we're going to take

We're going to [write some Go code](https://github.com/maxmind/mmdb-from-go-blogpost/blob/master/main.go) that makes use of the MaxMind [`mmdbwriter`](https://pkg.go.dev/github.com/maxmind/mmdbwriter) Go module to:

1. Load the GeoLite2 Country MaxMind DB.
   - We will take a pathname to the MMDB file and call [`mmdbwriter.Load()`](https://pkg.go.dev/github.com/maxmind/mmdbwriter?tab=doc#Load) on it, yielding `writer`, an [`*mmdbwriter.Tree`](https://pkg.go.dev/github.com/maxmind/mmdbwriter?tab=doc#Tree).
2. Add our own internal department data to the appropriate IP ranges.
	- We will call [`writer.InsertFunc()`](https://pkg.go.dev/github.com/maxmind/mmdbwriter?tab=doc#Tree.InsertFunc) once for each department's IP range.
3. Write the augmented database to a new MMDB file.
	- We will call [`writer.WriteTo()`](https://pkg.go.dev/github.com/maxmind/mmdbwriter?tab=doc#Tree.WriteTo).
4. Look up the new data in the augmented database to confirm our additions.
	- We will use the [`mmdbinspect`](https://github.com/maxmind/mmdbinspect) tool to see our new data in the MMDB file we've built and compare a few ranges in it to those in the old GeoLite2 Country MMDB file.

The full code is presented in the next section. Let's dive in!

### The code, explained

The repo for this tutorial is [available on GitHub](https://github.com/maxmind/mmdb-from-go-blogpost). You can clone it locally and `cd` into the repo dir by running the following in a terminal window:

```bash
me@myhost:~/dev $ git clone https://github.com/maxmind/mmdb-from-go-blogpost.git
me@myhost:~/dev $ cd mmdb-from-go-blogpost
```

Now I’m going to break down the contents of `main.go` from the repo, the code that will perform steps 1-3 of the tutorial. If you prefer to read the code directly, you can skip to the next section.

```go
package main
```
All Go programs begin with a `package main`, indicating that this file will contain a `main` function, the start of our program's execution. This program is no exception.

```go
import (
	"log"
	"net"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)
```
Most programs have a list of `import`ed packages next. In our case, the list of packages imported include some from the standard library: [`log`](https://golang.org/pkg/log/), which we use for outputting in error scenarios; [`net`](https://golang.org/pkg/net/), for the `net.ParseCIDR` function and the `net.IPNet` type, which we use when inserting new data into the MMDB tree; and [`os`](https://golang.org/pkg/os/), which we use when creating a new file into which we will write the MMDB tree. We also import some packages from MaxMind's [`mmdbwriter`](https://github.com/maxmind/mmdbwriter/) repo, which are designed specifically for building MMDB files and for working with MMDB trees -- you'll see how we use those below.

```go
func main() {
	// Load the database we wish to augment.
	writer, err := mmdbwriter.Load("GeoLite2-Country.mmdb", mmdbwriter.Options{})
	if err != nil {
		log.Fatal(err)
	}
```
Now we're at the start of the program execution. We begin by loading the existing database, `GeoLite2-Country.mmdb`, that we're going to augment.

```go
	// Define and insert the new data.
	_, sreNet, err := net.ParseCIDR("56.0.0.0/16")
	if err != nil {
		log.Fatal(err)
	}
```
Having loaded the existing GeoLite2 Country database, we begin defining the data we wish to augment it with. The second return value of the [`net.ParseCIDR()`](https://golang.org/pkg/net/#ParseCIDR) function is of type [`*net.IPNet`](https://golang.org/pkg/net/#IPNet), which is what we need for the first parameter for our upcoming [`writer.InsertFunc()`](https://pkg.go.dev/github.com/maxmind/mmdbwriter?tab=doc#Tree.InsertFunc) call, so we use `net.ParseCIDR` to go from the `string`-literal CIDR form `"56.0.0.0/16"` to the desired `*net.IPnet`.

```go
	sreData := mmdbtype.Map{
		"AcmeCorp.Environments": mmdbtype.Slice{
			mmdbtype.String("development"),
			mmdbtype.String("staging"),
			mmdbtype.String("production"),
		},
		"AcmeCorp.DeptName": mmdbtype.String("SRE"),
	}
```
Next, we're introduced to the [`mmdbtype.DataType`](https://pkg.go.dev/github.com/maxmind/mmdbwriter/mmdbtype?tab=doc#DataType) [interface](https://gobyexample.com/interfaces): Every piece of data that is attached to an IP or IP range in an MMDB file conforms to this interface.

While the MaxMind DB spec does not require it, all of the MMDB files built by MaxMind have either a map as a record or no record attached to each IP range; i.e. if a record exists for an IP range, its outermost value is a map, which for our purposes is a [`mmdbtype.Map`](https://pkg.go.dev/github.com/maxmind/mmdbwriter/mmdbtype?tab=doc#Map), which satisfies the `mmdbtype.DataType` interface.

\[An aside: If you look at the output of running the `mmdbinspect -db GeoLite2-Country.mmdb 56.0.0.1` command in your terminal, examining the `$.[0].Records[0].Record` [JSONPath](https://goessner.net/articles/JsonPath/) (i.e. the sole record, stripped of its wrappers), then you'll see that it is a JSON Object, which as expected corresponds to the `mmdbtype.Map` type.\]

We're going to take advantage of this, adding the two previously mentioned key/value pairs previously to all the existing maps for those IP ranges, where they exist, or creating a map with just those key/value pairs, where one didn't exist. The keys are `AcmeCorp.Environments` (whose value will be an [`mmdbtype.Slice`](https://pkg.go.dev/github.com/maxmind/mmdbwriter/mmdbtype?tab=doc#Slice), a [slice](https://gobyexample.com/slices) of [`mmdbtype.String`](https://pkg.go.dev/github.com/maxmind/mmdbwriter/mmdbtype?tab=doc#String)s containing the allowed environment strings for an IP range), and `AcmeCorp.DeptName` (whose value will be an [`mmdbtype.String`](https://pkg.go.dev/github.com/maxmind/mmdbwriter/mmdbtype?tab=doc#String) that is the name of the department for an IP range).

```go
	if err := writer.InsertFunc(sreNet, inserter.TopLevelMergeWith(sreData)); err != nil {
		log.Fatal(err)
	}
```
Now that we've got our data, we insert it into the MMDB tree using the [`inserter.TopLevelMergeWith`](https://pkg.go.dev/github.com/maxmind/mmdbwriter/inserter?tab=doc#TopLevelMergeWith) strategy, leaving us with an MMDB tree with the AcmeCorp SRE IPs in the 56.0.0.0/16 range, whose maps contain the new environment and department name keys in addition to whatever GeoLite2 Country data they returned previously. (Note that we carefully picked non-clashing, top-level keys; no key in the GeoLite2 Country data starts with `AcmeCorp.`)

What happens if there is an IP for which no record exists? With the `inserter.TopLevelMergeWith` strategy, this IP will also happily take our new top-level keys as well.

```go
	_, devNet, err := net.ParseCIDR("56.1.0.0/16")
	if err != nil {
		log.Fatal(err)
	}
	devData := mmdbtype.Map{
		"AcmeCorp.Environments": mmdbtype.Slice{
			mmdbtype.String("development"),
			mmdbtype.String("staging"),
		},
		"AcmeCorp.DeptName": mmdbtype.String("Development"),
	}
	if err := writer.InsertFunc(devNet, inserter.TopLevelMergeWith(devData)); err != nil {
		log.Fatal(err)
	}

	_, mgmtNet, err := net.ParseCIDR("56.2.0.0/16")
	if err != nil {
		log.Fatal(err)
	}
	mgmtData := mmdbtype.Map{
		"AcmeCorp.Environments": mmdbtype.Slice{
			mmdbtype.String("development"),
			mmdbtype.String("staging"),
		},
		"AcmeCorp.DeptName": mmdbtype.String("Management"),
	}
	if err := writer.InsertFunc(mgmtNet, inserter.TopLevelMergeWith(mgmtData)); err != nil {
		log.Fatal(err)
	}
```
We repeat the process for the Development and Management departments, taking care to update the range itself, the list of environments, and the department name as we go.

```go
	// Write the newly augmented DB to the filesystem.
	fh, err := os.Create("GeoLite2-Country-with-Department-Data.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}
```
Finally we write the new database to disk.

### Building the code and running it

So that's our code! Now we build the program and run it. On my 2015-model laptop it takes under 10 seconds to run.

```bash
me@myhost:~/dev/mmdb-from-go-blogpost $ go build
me@myhost:~/dev/mmdb-from-go-blogpost $ ./mmdb-from-go-blogpost 
```

This will have built the augmented database. Finally, we will compare some IP address and range queries on the original and augmented database using the [`mmdbinspect`](https://github.com/maxmind/mmdbinspect) tool.

```bash
me@myhost:~/dev/mmdb-from-go-blogpost $ mmdbinspect -db GeoLite2-Country.mmdb -db GeoLite2-Country-with-Department-Data.mmdb 56.0.0.1 56.1.0.0/24 56.2.0.54 56.3.0.1 | less
```

The [output](https://gist.github.com/nchelluri/ad079300b92a634bc4b36249b77f3893) from this command, elided here for brevity, shows us that the `AcmeCorp.Environments` and `AcmeCorp.DeptName` keys are not present in the original MMDB file at all and that they are present in the augmented MMDB file when expected. The 56.3.0.1 IP address remains identical across both databases (without any AcmeCorp fields) as a control.

And that's it! You've now built yourself a GeoLite2 Country MMDB file augmented with custom data.

Feel free to open an issue in the [repo](https://github.com/maxmind/mmdb-from-go-blogpost/issues) if you have any questions or just want to tell us what you've created.
