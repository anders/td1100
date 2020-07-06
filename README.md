# td1100

Tool to configure Datalogic TD1100 handheld barcode readers. It outputs a list (PDF) of special barcodes to change the
behavior of the scanner.

The configuration format is: `<FNC3> + $ + cmd[, cmd] + \r` (code 128 barcode).

The default list does the following:

1. Restore to EU factory settings (P,Ae,P)
2. Simulate USB keyboard (P,HA35,P - maybe it could be done as $HA35, don't remember)
3. Toggle programming (P)
4. Enable USB sleep mode (CUSSE01)
5. Enable constant reading (CSNRM04)
6. Enable UK Plessey barcode support (CPLEN01)
7. Enable UK Plessey checksum calculation (CPLCC01)
8. Include UK Plessey checksum (CPLCT01)
9. Disable EAN13 (C3BEN00)
10. Disable EAN8 (C8BEN00)
11. Toggle programming (P)

## Usage
````
$ go build
$ ./td1100 -help
Usage of ./td1100:
  -list string
       path to JSON list of barcodes (optional)
  -output string
       path to PDF output
  -writeList
       write default list to file
$ ./td1100 -output codes.pdf
````

You can modify the list by first running `-writeList list.json`, then run the tool with `-list list.json`.

## License
MIT license (see LICENSE.txt).

