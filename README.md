
# codice-fuzzcale

[![asciicast](https://asciinema.org/a/428249.svg)](https://asciinema.org/a/428249)

codice-fuzzcale is a CLI tool written in Go for generating Italian fiscal codes (codici fiscali) for pentesting Italian web applications such as portals, banking apps etc. 

The fiscal code is a unique identifier for an individual, and is composed of personal data such as name, surname, date of birth, place of birth and sex. 

codice-fuzzcale will generate a list of all possible fiscal codes, based on incomplete information provided by the user and write the list to an outfile which can then be used in further pentesting activities. 

For example, if I know the name, surname, sex and place of birth of an individual, but not their date of birth, I can enter those data to the tool and it will generate a list of all possible codes, which I can then use in further penetration testing on the target application. 

Please note, this tool is intended only for use in target-approved ethical hacking activities, to help raise application security and privacy levels. 
