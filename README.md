# Tetra Processor

A Lambda written in Go as part of the Tesseract ecosystem. Pre-processes csv files which are dropped into an S3 bucket.

## Deploying

1. run `make` in the `tetra-processor` directory.
2. Zip the output at `tetra-processor/bin/cx-tetra-processor`
3. Upload the output to the Lambda.