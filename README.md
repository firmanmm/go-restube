# Go Restube
A REST Based Youtube Downloader. It allows you switch between local storage and remote storage. Supported remote storage are Azure Redis Cache and Azure Blob Storage.
```
Note : This is a small project that is meant for tutorial only since this project ignore most security best practice like hashing with secure password hashing algorithm and many more
```

## Running
To get all available paramters please use -h flag
`./go-restube -h`
### Flags
```
-azRA string
    Azure Redis Address for session storage
-azRK string
    Azure Redis Key for session storage
-azSK string
    Azure Storage Key for remote file storage
-azSU string
    Azure Storage Username for remote file storage
```

## Endpoint
### View Video Information
#### Description
Used to obtain video information
#### Endpoint
**GET  /video**
#### URL Parameter
*None*
##### Query Parameter
*url : string : required* = Video URL to be checked
##### Post Parameter :
*None*
##### Header
*Authorization : string : required* = Authorization that will be used to access the API

### Request Video Download
#### Description 
Used to request a video to be downloaded
#### Endpoint
**POST  /video**
#### URL Parameter
*None*
#### Query Parameter 
*None*
#### Post Parameter
*url : string : required* = Video URL to be requested
*mode : int : required* = Video mode to be downloaded (Can be checked on *View Video Information* API)
#### Header
*Authorization : string : required* = Authorization that will be used to access the API

### Downloading Video
#### Description 
Used to download requested video
#### Endpoint
**GET /video/:url**
#### URL Parameter
*url : string : required* = Video URL to be requested
#### Query Parameter 
*None*
#### Post Parameter
*None*
#### Header
*Authorization : string : required* = Authorization that will be used to access the API

### Checking Current Usage
#### Description 
Used to check current user's usage
#### Endpoint
**GET /usage**
#### URL Parameter
*None*
#### Query Parameter 
*None*
#### Post Parameter
*None*
#### Header
*Authorization : string : required* = Authorization that will be used to access the API

### Checking All User's Usage
#### Description 
Used to check all user's usage. Only username named **myadmin** can access it.
#### Endpoint
**GET /usage/all**
#### URL Parameter
*None*
#### Query Parameter 
*None*
#### Post Parameter
*None*
#### Header
*Authorization : string : required* = Authorization that will be used to access the API

### Creating Account
#### Description 
Used to create new account.
#### Endpoint
**POST  /register**
#### URL Parameter
*None*
#### Query Parameter 
*None*
#### Post Parameter
*username : string : required* = Username to be registered
*password : string: required* =  Password to be associated with the username
#### Header
*None*

### Obtaining Authorization
#### Description 
Used to obtain authorization header.
#### Endpoint
**POST  /login**
#### URL Parameter
*None*
#### Query Parameter 
*None*
#### Post Parameter
*username : string : required* = Username to be checked
*password : string: required* =  Password to be checked with the username
#### Header
*None*