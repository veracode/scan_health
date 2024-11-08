![GitHub License](https://img.shields.io/github/license/veracode/scan_health)
![Go Version](https://img.shields.io/github/go-mod/go-version/veracode/scan_health)
![Docker Image Size](https://img.shields.io/docker/image-size/antfie/scan_health/latest)
![Downloads](https://img.shields.io/github/downloads/veracode/scan_health/total)
[![Go Report Card](https://goreportcard.com/badge/github.com/veracode/scan_health/v2)](https://goreportcard.com/report/github.com/github.com/veracode/scan_health/v2)
![Docker Pulls](https://img.shields.io/docker/pulls/antfie/scan_health)

# Veracode Scan Health Tool 🏥

**Note this tool is not an official Veracode product. It comes with no support or warranty.**

Use this console tool to gain insight into the health of a Veracode scan. The report includes any problems found as well as suggestions to improve scan performance and flaw quality. At present this tool only works on Static Analysis (SAST) scans.

The tool has the following modes of operation:

* **Health Report** - For when we want to know if a given SAST scan has been configured correctly
* **Scan Compare** - For when we want to compare two SAST scans for troubleshooting issues

## Key Features ✅

* Uses the [Veracode APIs](https://docs.veracode.com/r/Veracode_APIs) and standard Veracode [API Credentials File](https://docs.veracode.com/r/c_api_credentials3).
* Works in Veracode Commerical (US) and European regions.
* Easy to understand configuration issues and recommendations.
* Outputs a summary to the console and also optionally a full listing in JSON.
* Returns exit code 1 if there are any serious configuration issues found when using `-error-on-high-severity` flag. 
* Set `-previous-scan=true` to run checks against the previous scan. We may make this the default at some point in the
  future.

## Installation 📦

The easiest way to install this is to download the binary from GitHub, but there are other options.

### From GitHub 

From the [Releases](https://github.com/veracode/scan_health/releases) page you can download the latest version of this standalone executable for Windows, Linux and macOS. Note that at this time we are unable to code-sign these binaries.

### Using Docker 🐳

You can also use the [Docker image](https://hub.docker.com/r/antfie/scan_health):

```bash
docker pull antfie/scan_health
```

### Compile Manually

You will need the latest version of [Go](https://go.dev/dl/) and then run the following commands to get the code and compile it:

```bash
git clone https://github.com/veracode/scan_health
cd scan_health
./scripts/build.sh
```

Note it is best to point to the latest release commit as the main branch can sometimes have have work in progress.

In the `dist` folder you should find the compiled binaries.

## Configuration 🛠️

This tool makes use of the [Veracode APIs](https://docs.veracode.com/r/Veracode_APIs) ([listed below](#outbound-api-calls)). You will need a Veracode [API Credentials File](https://docs.veracode.com/r/c_api_credentials3) and the [Reviewer or Security Lead role](https://docs.veracode.com/r/c_API_roles_details#results-api) for this tool to work.

Follow the instructions on how to configure a Veracode [API Credentials File](https://docs.veracode.com/r/c_api_credentials3) if you do not already have this set up.

Alternatively you can use
environment
variables (`VERACODE_API_KEY_ID` and `VERACODE_API_KEY_SECRET`) or the CLI flags (`-vid` and `-vkey`) to authenticate
with the Veracode APIs.

### Using Docker 🐳

If you are using Docker you can mount the API credentials file in the container like this:

```bash
docker run -t -v "$HOME/.veracode:/.veracode" antfie/scan_health -action helth -sast https://analysiscenter.veracode.com/auth/index.jsp#...
```

## Using the Tool

Simply running the tool with `-h` should provide helpful usage instructions like so:

```bash
Scan Health v2.34
Copyright © Veracode, Inc. 2024. All Rights Reserved.
This is an unofficial Veracode product. It does not come with any support or warranty.

Usage of ./dist/scan_health-mac-amd64:
  -cache
    	Enable caching of API responses (useful for development)
  -error-on-high-severity
    	Return a non-zero exit code if any high severity issues are found
  ...
```

### Health Report

The most common use of this tool is to get a health report for a given scan.

```bash
./scan_health -action health -sast https://analysiscenter.veracode.com/auth/index.jsp#...
```

Here is an example report:

```
Scan Health v2.32
Copyright © Veracode, Inc. 2023. All Rights Reserved.
This is an unofficial Veracode product. It does not come with any support or warranty.

Inspecting SAST build id = 22132159 in the commercial region
Warning: This is not the latest scan

Scan Summary
============
Business unit:      Cyber Security
Application:        Secure File Transfer
Sandbox:            Development
Scan name:          15 Nov 2022 Static
Review Modules URL: https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:75603:793744:22132159:22103486:22119136::::5000002
Triage Flaws URL:   https://analysiscenter.veracode.com/auth/index.jsp#ReviewResultsStaticFlaws:75603:793744:22132159:22103486:22119136::::5000002
Files uploaded:     127
Total modules:      142
Modules selected:   1
Engine version:     20221021161836 (Release notes: https://docs.veracode.com/updates/r/c_all_static)
Submitted:          2022-11-15 12:28:29 +0000 GMT (225d 2h 16m 5.453014s ago)
Published:          2022-11-15 12:37:33 +0000 GMT (225d 2h 7m 1.453023s ago)
Duration:           9m 4s
Latest app scan:    2023-01-31 08:52:10 +0000 GMT (148d 5h 52m 24.453035s ago)

Flaw Summary
============
Total (less fixed):         49
Fixed (no longer reported): 0
Policy affecting:           35
Mitigated:                  2
Open affecting policy:      33
Open not affecting policy:  14

Modules Selected for Analysis
=============================
* app.war

Issues
======
❌ A Java module was identified that contained no compiled Java classes: "common-0.1.0-SNAPSHOT.jar". Sometimes this can indicate the file contained only Java source code which is not processed by Veracode.
⚠️ 2 unnecessary files were uploaded: ".DS_Store", "._.DS_Store".
⚠️ 2 JavaScript modules were not selected for analysis: "JS files within common-0.1.0-SNAPSHOT.jar", "JS files within auth-0.1.0-SNAPSHOT.jar".
⚠️ There have not been recent scans of this application The application was last scanned on 2023-01-31 08:52:10 +0000 GMT which was 148d 5h 52m 24.452762s ago. It is not uncommon for new flaws to be reported over time because Veracode is continuously improving their products, and because new SCA vulnerabilities are reported every day, and this could impact the application.

Recommendations
===============
💡 Follow the packaging instructions to keep the upload as small as possible to improve upload and scan times.
💡 Veracode requires the Java application to be compiled into a JAR, WAR or EAR file as per the packaging instructions.
💡 Veracode extracts JavaScript modules from the upload. Consider selecting the appropriate "JS files within ..." modules for analysis in order to cover the JavaScript risk from these components.
💡 Under-selection of first party modules affects results quality. Ensure the correct entry points have been selected as recommended and refer to this article: https://community.veracode.com/s/article/What-are-Modules-and-how-do-my-results-change-based-on-what-I-select.
💡 Regular scanning, preferably via automation will allow the application team to respond faster to any new issues reported.
💡 Consider scheduling a consultation to review the packaging and module configuration: https://docs.veracode.com/r/t_schedule_consultation.
```

The tool also outputs JSON if the `-format json` flag is set. Note that unlike the console this JSON will not contain summaries but all the details. An example can be seen below:

```json
{
  "health_tool": {
    "report_date": "2023-06-29T12:20:49.769824+01:00",
    "version": "2.1",
    "region": "commercial"
  },
  "scan": {
    "account_id": 75603,
    "business_unit": "Cyber Security",
    "application_id": 793744,
    "application_name": "Secure File Transfer",
    "scan_name": "30 Jan 2023 Static",
    "sandbox_id": 5139524,
    "sandbox_name": "Dart and Flutter",
    "build_id": 23664244,
    "review_modules_url": "https://analysiscenter.veracode.com/auth/index.jsp#AnalyzeAppModuleList:75603:793744:23664244:23635472:23651122::::5139524",
    "triage_flaws_url": "https://analysiscenter.veracode.com/auth/index.jsp#ReviewResultsStaticFlaws:75603:793744:23664244:23635472:23651122::::5139524",
    "engine_version": "20230123231701",
    "submitted_date": "2023-01-30T12:17:19Z",
    "published_data": "2023-01-30T12:27:38Z",
    "scan_duration": 619000000000,
    "analysis_size": 157513561,
    "is_latest_scan": false
  },
  "flaws": {
    "total": 5,
    "total_affecting_policy": 5,
    "open_affecting_policy": 5
  },
  "uploaded_files": [
    {
      "id": 10552845755,
      "name": "wonderous.ipa",
      "status": "Uploaded",
      "md5": "502f1a945eb85a0d09da917bba5d7de0",
      "is_ignored": false,
      "is_third_party": false
    }
  ],
  "modules": [
    {
      "name": "wonderous.ipa",
      "compiler": "DART",
      "operating_system": "Dart",
      "architecture": "DART",
      "is_ignored": false,
      "is_third_party": false,
      "is_dependency": false,
      "is_selected": true,
      "has_fatal_errors": false,
      "flaw_summary": {}
    },
    {
      "id": 1756110687,
      "name": "wonderous.ipa",
      "is_ignored": false,
      "is_third_party": false,
      "is_dependency": false,
      "is_selected": false,
      "has_fatal_errors": false,
      "status": "OK",
      "platform": "DART / Dart / DART",
      "size": "150MB",
      "md5": "502f1a945eb85a0d09da917bba5d7de0",
      "issues": [
        "No supporting files or PDB files"
      ],
      "flaw_summary": {}
    }
  ],
  "issues": [
    {
      "description": "There have not been recent scans of this application The application was last scanned on 2023-01-31 08:52:10 +0000 GMT which was 149d 2h 28m 42.711006s ago. It is not uncommon for new flaws to be reported over time because Veracode is continuously improving their products, and because new SCA vulnerabilities are reported every day, and this could impact the application.",
      "severity": "medium"
    }
  ],
  "recommendations": [
    "Regular scanning, preferably via automation will allow the application team to respond faster to any new issues reported.",
    "Consider scheduling a consultation to review the packaging and module configuration: https://docs.veracode.com/r/t_schedule_consultation."
  ],
  "last_activity": "2023-01-31T08:52:10Z"
}
```

### Scan Compare

Scans can be compared like so:

```bash
./scan_health -action compare -a https://analysiscenter.veracode.com/auth/index.jsp#... -b https://analysiscenter.veracode.com/auth/index.jsp#...
```

## Different ways to run

If you know the build IDs you can use them instead of URLs if preferred, like so:

```sh
./scan_health -sast 22132159
```

Please note that this will only work if the scan has a Detailed Report, i.e. the scan has finished. It is preferable to use scan URLs over build IDs. 

### Using a zsh helper alias

If you use the zsh shell you might want to have a helper alias. To do that add this to your `~/.zshrc` file:

```sh
alias vsh='f() { /path/to/scan_health-mac-arm64 -sast "$1" };f'
```

Then you can simply run:

```sh
vsh https://analysiscenter.veracode.com/auth/index.jsp#...
```

## Development 🪚

You may want to use the `-cache=true` flag to speed up development by caching API responses. It is not recommended to use caching in production as sensitive data may be left in files.

### Compiling

Runningt he build script will also run the tests:

```sh
./scripts/build.sh
```

## Bug Reports 🐞

If you find a bug, please file an Issue right here in GitHub, and I will try to resolve it in a timely manner.

## API Calls

This tool makes the following read-only API requests for data:

| API                                                                         | Notes / API Documentation                                  |
|-----------------------------------------------------------------------------|------------------------------------------------------------|
| <https://github.com/veracode/scan_health/releases/latest>                     | To determine if a new version of the tool is available.    |
| <https://analysiscenter.veracode.com/api/3.0/getmaintenancescheduleinfo.do> | <https://docs.veracode.com/r/r_getmaintenancescheduleinfo> |
| <https://analysiscenter.veracode.com/api/5.0/detailedreport.do>             | <https://docs.veracode.com/r/r_detailedreport>             |
| <https://analysiscenter.veracode.com/api/5.0/getappinfo.do>                 | <https://docs.veracode.com/r/r_getappinfo>                 |
| <https://analysiscenter.veracode.com/api/5.0/getbuildinfo.do>               | <https://docs.veracode.com/r/r_getbuildinfo>               |
| <https://analysiscenter.veracode.com/api/5.0/getbuildlist.do>               | https://docs.veracode.com/r/r_getbuildlist                 |
| <https://analysiscenter.veracode.com/api/5.0/getfilelist.do>                | <https://docs.veracode.com/r/r_getfilelist>                |
| <https://analysiscenter.veracode.com/api/5.0/getprescanresults.do>          | <https://docs.veracode.com/r/r_getprescanresults>          |
| <https://analysiscenter.veracode.com/api/5.0/getsandboxlist.do>             | <https://docs.veracode.com/r/r_getsandboxlist>             |
