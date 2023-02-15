![Go Version](https://img.shields.io/github/go-mod/go-version/antfie/scan_health)
![Docker Image Size](https://img.shields.io/docker/image-size/antfie/scan_health/latest)
![Downloads](https://img.shields.io/github/downloads/antfie/scan_health/total)

# Veracode Scan Health 🏥

This is an unofficial Veracode product. It does not come with any support or warrenty.

Use this console tool to see the health of Veracode Static Analysis (SAST) scans and get some suggestions to improve scan performance anf flaw quality. The scans must have completed for the tool to work.

## Usage

We recommend you configure a Veracode API credentials file as documented here: https://docs.veracode.com/r/c_configure_api_cred_file.

Alternatively you can use environment variables (`VERACODE_API_KEY_ID` and `VERACODE_API_KEY_SECRET`) or CLI flags (`-vid` and `-vkey`) to authenticate with the Veracode APIs.

```
./scan_health -h
Scan Health v1.x
Copyright © Veracode, Inc. 2023. All Rights Reserved.
This is an unofficial Veracode product. It does not come with any support or warrenty.

Usage of scan_health:
  -sast string
        Veracode Platform URL or build ID for health review"
  -region string
        Veracode Region [global, us, eu]
  -vid string
        Veracode API ID - See https://docs.veracode.com/r/t_create_api_creds
  -vkey string
        Veracode API key - See https://docs.veracode.com/r/t_create_api_creds
```

## Example Usage

```
./scan_health -sast https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview:75603:793744:22132159:22103486:22119136::::5000002
```

If you know the build IDs you can use them instead of URLs like so:

```
./scan_health -sast 22132159
```

Using Docker:

```
docker run -t -v "$HOME/.veracode:/.veracode" antfie/scan_health -sast https://analysiscenter.veracode.com/auth/index.jsp#StaticOverview:75603:793744:22132159:22103486:22119136::::5000002
```

## Development 🛠️

### Running

```
go run *.go
```

### Compiling

```
./release.sh
```

# Bug Reports 🐞

If you find a bug, please file an Issue right here in GitHub, and I will try to resolve it in a timely manner.