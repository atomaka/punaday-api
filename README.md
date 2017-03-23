# punaday-api

JSON version of http://www.punoftheday.com/

## Description

Server that takes Pun of the Day puns and converts it into JSON.  Currently
living as a serverless application in AWS Lamda. It is frontended by AWS
Cloudfront because AWS Cloudformation does not yet have resources for AWS API
Gateway domains.  This script will also create a DNS record to point at the
AWS Cloudfront Distribution.

## Usage

* `cp .env.sample .env.production`
  * `AWS_CERTIFICATE`: arn to aws certificate that you generated for a domain
  * `DNS_ZONE`: Route 53 zone name (ex. `example.com.`)
  * `DOMAIN`: Route 53 domain (ex. `puns.example.com`)
* Set AWS credentials via environment or credential file
* `make`
