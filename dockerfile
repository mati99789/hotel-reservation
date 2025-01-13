FROM golang:1.23.1

# Set the working directory to /app
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory 
COPY go.mod go.sum ./

# download and insttall any required Go dependenciesdencies
RUN go mod download

# Copy the entire source code to the working directory
COPY . .

# Build the go application
RUN go build -o main .

# Expose the port specified by the PORT environment variable
EXPOSE 3000

# Set the etry point of the container to the executable
CMD ["./main"]
