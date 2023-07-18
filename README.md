# Archive-differ

The archive differ is an application that takes a archived bag in the Bagit
format, and checks the mets file for difference in comparison to another bag.
If there are differences it will display the difference that bags have to
the user in a useful and readable format.


## Features
The first focus is to display the difference and counts of the preservation events.
It will also give users the output of the diff tag if asked for.

The tool may also do other things given time.

## Usage
The archive differ is primarly a command line tool although in the
future if it proves useful, it's feature set may include a text user interface.

## Motivation 
The reason this exists is to test different digital archiving
tools to see if they produce the same types of preservation metadata in a
quick way without resorting to human error and intense labour.

The primary applications it is meant to test is the a3m and archivematica
Archival Information Packages.
