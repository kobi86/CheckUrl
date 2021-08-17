# Check URL

This tool scans a file with URL's, exports them into a seprate file and then checks if the url's exists. 

## Prerequisites  

- Go version 1.16.6

## Usage 

This program uses config file for several settings: 

- httpslist - File which includes the urls that needs to be scaned 
- outputfile - Output file which will  include only the urls (in case https-list file includes more then urls )
- search - Term to search like http/https/www/etc.. 
- errorCodes - List of codes that will output an error 

