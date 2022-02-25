# main.go README

Implements the page number change made in host.go. 
Search automatically starts on page 1, and increments to display results on next page on (Y) input
Any other input other than (Y) will end the search. 

Search results expanded to include IP, Port, City, Country, and Organization related to the search. 

Also removed Host Data dump for clarity in final output. 

to use:
go build main.go
SHODAN_API_KEY={YOUR KEY} ./main {search term or filter} 

Ex. SHODAN_API_KEY={YOUR KEY} ./main apple 
    SHODAN_API_KEY={YOUR KEY} ./main city:paris