# Using base image of Ubuntu
FROM ubuntu AS base

# Copy the SQLite3 from local directory to docker image
ADD  ./database/sqlite3/sqlite3 sqlite3

# Run the SQLite inside docker container when docker image is at run
CMD [ "./sqlite3" ]