Tetra Processor
========

Tetra Processor is an AWS Lambda application which performs CSV operations in
 the context of Tesseract.

1. run `make` in the `tetra-processor` directory.
2. Zip the output, `tetra-processor/bin/cx-tetra-processor`
3. Upload the zip to the Lambda.
4. Configure S3 as an input trigger for the Lambda.

Features
--------

- CSV Operations (e.g. Slice Rows)
- Rename output files
- Decrypt incoming PGP-encrypted CSV files

Installation
------------

Install `tetra-processor` by running:
    
    go get github.com/arizard/tetra-processor
    
And then to clone the repo:

    cd $GOPATH/src/github.com/arizard
    git clone git@github.com:Arizard/tetra-processor.git

Contribute
----------

- Issue Tracker: https://github.com/arizard/tetra-processor/issues
- Source Code: https://github.com/Arizard/tetra-processor

Support
-------

If you are having issues, please submit them in the issue tracker.